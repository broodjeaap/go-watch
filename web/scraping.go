package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/andybalholm/cascadia"
	"github.com/antchfx/htmlquery"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	lualibs "github.com/vadv/gopher-lua-libs"
	lua "github.com/yuin/gopher-lua"
	"golang.org/x/net/html"
	"gorm.io/gorm"

	. "github.com/broodjeaap/go-watch/models"
)

func ProcessFilters(filters []Filter, web *Web, watch *Watch, debug bool, scheduleID *uint) {
	processedMap := make(map[uint]bool, len(filters))
	if scheduleID != nil {
		processedMap[*scheduleID] = true
	}

	for _, filter := range filters {
		if scheduleID != nil && filter.ID == *scheduleID {
			log.Println(fmt.Sprintf("Scheduled Watch for '%s', triggered by schedule '%s'", watch.Name, filter.Name))
		}
	}

	// check if there are multiple 'cron' filters connected to a single filter
	// this filter can't be run, only one filter is ever triggered,
	// the other prevents it from running, allParentsProcessed is always false
	// just warn the user when this happens and return
	for i := range filters {
		filter := &filters[i]
		cronParentCount := 0
		for _, parent := range filter.Parents {
			if parent.Type == "cron" {
				cronParentCount++
			}
		}
		if cronParentCount > 1 {
			filter.Log("Multiple schedules on the same filter is not supported!")
			return
		}
	}

	currentFilters := make([]*Filter, 0, len(filters))
	for i := range filters {
		filter := &filters[i]
		currentFilters = append(currentFilters, filter)
	}

	// collect 'store' and 'notify' filters so we can process them separately  at the end
	storeFilters := make([]*Filter, 0, 5)
	notifyFilters := make([]*Filter, 0, 5)

	for {
		nextFilters := make([]*Filter, 0, len(currentFilters))
		for i := range currentFilters {
			filter := currentFilters[i]
			if filter.Type == "store" {
				storeFilters = append(storeFilters, filter)
				processedMap[filter.ID] = true
				continue
			}
			if debug && filter.Type == "cron" {
				if filter.Var2 != nil && *filter.Var2 == "no" {
					filter.Log("Schedule is disabled")
				}
				processedMap[filter.ID] = true
				getCronDebugResult(filter)
				continue
			}
			if filter.Type == "echo" {
				getFilterResultEcho(filter)
				processedMap[filter.ID] = true
				continue
			}
			if debug && filter.Type == "notify" {
				notifyFilters = append(notifyFilters, filter)
			}
			if len(filter.Parents) == 0 && !debug {
				continue
			}
			var allParentsProcessed = true
			for _, parent := range filter.Parents {
				if _, contains := processedMap[parent.ID]; !contains {
					allParentsProcessed = false
					break
				}
			}
			if !allParentsProcessed {
				nextFilters = append(nextFilters, filter)
				continue
			}
			getFilterResult(filters, filter, watch, web, debug)
			processedMap[filter.ID] = true
		}
		if len(nextFilters) == 0 {
			break
		}
		currentFilters = nextFilters
	}

	// process the store filters last
	for _, storeFilter := range storeFilters {
		getFilterResult(filters, storeFilter, watch, web, debug)
	}

	// process the notify filters when editing, so it still logs and can put the message in results
	if debug {
		for _, nFilter := range notifyFilters {
			if len(nFilter.Results) == 0 {
				notifyFilter(filters, nFilter, watch, web, debug)
				log.Println(nFilter.Results)
			}
		}
	}
}

func getFilterResult(filters []Filter, filter *Filter, watch *Watch, web *Web, debug bool) {
	switch {
	case filter.Type == "gurl":
		{
			getFilterResultURL(filter, web.urlCache, debug)
		}
	case filter.Type == "gurls":
		{
			getFilterResultURLs(filter, web.urlCache, debug)
		}
	case filter.Type == "xpath":
		{
			getFilterResultXPath(filter)
		}
	case filter.Type == "json":
		{
			getFilterResultJSON(filter)
		}
	case filter.Type == "css":
		{
			getFilterResultCSS(filter)
		}
	case filter.Type == "replace":
		{
			getFilterResultReplace(filter)
		}
	case filter.Type == "match":
		{
			getFilterResultMatch(filter)
		}
	case filter.Type == "substring":
		{
			getFilterResultSubstring(filter)
		}
	case filter.Type == "contains":
		{
			getFilterResultContains(filter)
		}
	case filter.Type == "math":
		{
			switch {
			case filter.Var1 == "sum":
				{
					getFilterResultSum(filter)
				}
			case filter.Var1 == "min":
				{
					getFilterResultMin(filter)
				}
			case filter.Var1 == "max":
				{
					getFilterResultMax(filter)
				}
			case filter.Var1 == "average":
				{
					getFilterResultAverage(filter)
				}
			case filter.Var1 == "count":
				{
					getFilterResultCount(filter)
				}
			case filter.Var1 == "round":
				{
					getFilterResultRound(filter)
				}
			}
		}
	case filter.Type == "store":
		{
			storeFilterResult(filter, web.db, debug)
		}
	case filter.Type == "notify":
		{
			notifyFilter(filters, filter, watch, web, debug)
		}
	case filter.Type == "lua":
		{
			getFilterResultLua(filter)
		}
	case filter.Type == "cron":
		{

		}
	case filter.Type == "condition":
		{
			switch filter.Var1 {
			case "diff":
				{
					getFilterResultConditionDiff(filter, web.db)
				}
			case "lowerl":
				{
					getFilterResultConditionLowerLast(filter, web.db)
				}
			case "lowest":
				{
					getFilterResultConditionLowest(filter, web.db)
				}
			case "lowert":
				{
					getFilterResultConditionLowerThan(filter)
				}
			case "higherl":
				{
					getFilterResultConditionHigherLast(filter, web.db)
				}
			case "highest":
				{
					getFilterResultConditionHighest(filter, web.db)
				}
			case "highert":
				{
					getFilterResultConditionHigherThan(filter)
				}
			}
		}
	case filter.Type == "brow":
		{
			switch filter.Var1 {
			case "gurl":
				{
					getFilterResultBrowserlessURL(filter, web.urlCache, debug)
				}
			case "gurls":
				{
					getFilterResultBrowserlessURLs(filter, web.urlCache, debug)
				}
			case "func":
				{
					getBrowserlessFunctionResult(filter)
				}
			case "funcs":
				{
					getBrowserlessFunctionResults(filter)
				}
			}
		}
	case filter.Type == "echo":
		{
			getFilterResultEcho(filter)
		}
	default:
		filter.Log("getFilterResult called with filter.Type == '", filter.Type, "'")
	}
}

func getFilterResultURL(filter *Filter, urlCache map[string]string, debug bool) {
	fetchURL := filter.Var1
	val, exists := urlCache[fetchURL]
	if debug && exists {
		filter.Results = append(filter.Results, val)
		return
	}
	str, err := getURLContent(filter, fetchURL)
	if err != nil {
		log.Println("Could not fetch url: ", fetchURL, " - ", err)
		filter.Log("Could not fetch url: ", fetchURL, " - ", err)
		return
	}
	filter.Results = append(filter.Results, str)
	if debug {
		urlCache[fetchURL] = str
	}
}

func getFilterResultURLs(filter *Filter, urlCache map[string]string, debug bool) {
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			fetchURL := result
			val, exists := urlCache[fetchURL]
			if debug && exists {
				filter.Results = append(filter.Results, val)
				continue
			}

			str, err := getURLContent(filter, fetchURL)
			if err != nil {
				log.Println("Could not fetch url: ", fetchURL, " - ", err)
				filter.Log("Could not fetch url: ", fetchURL, " - ", err)
				continue
			}
			filter.Results = append(filter.Results, str)
			if debug {
				urlCache[fetchURL] = str
			}
		}
	}
}

func getURLContent(filter *Filter, fetchURL string) (string, error) {
	var httpClient *http.Client
	if viper.IsSet("proxy.url") {
		proxyUrl, err := url.Parse(viper.GetString("proxy.url"))
		if err != nil {
			return "", err
		}
		httpClient = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	} else {
		httpClient = &http.Client{}
	}
	resp, err := httpClient.Get(fetchURL)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func getFilterResultBrowserlessURL(filter *Filter, urlCache map[string]string, debug bool) {
	if filter.Var2 == nil {
		filter.Log("filter.Var2 == nil")
		return
	}
	fetchURL := *filter.Var2
	val, exists := urlCache["b"+fetchURL]
	if debug && exists {
		filter.Results = append(filter.Results, val)
		return
	}
	str, err := getBrowserlessURLContent(filter, fetchURL)
	if err != nil {
		log.Println("Could not fetch url: ", fetchURL, " - ", err)
		filter.Log("Could not fetch url: ", fetchURL, " - ", err)
		return
	}
	filter.Results = append(filter.Results, str)
	if debug {
		urlCache["b"+fetchURL] = str
	}
}

func getFilterResultBrowserlessURLs(filter *Filter, urlCache map[string]string, debug bool) {
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			fetchURL := result
			val, exists := urlCache["b"+fetchURL]
			if debug && exists {
				filter.Results = append(filter.Results, val)
				continue
			}

			str, err := getBrowserlessURLContent(filter, fetchURL)
			if err != nil {
				log.Println("Could not fetch url: ", fetchURL, " - ", err)
				filter.Log("Could not fetch url: ", fetchURL, " - ", err)
				continue
			}
			filter.Results = append(filter.Results, str)
			if debug {
				urlCache["b"+fetchURL] = str
			}
		}
	}
}

func getBrowserlessURLContent(filter *Filter, fetchURL string) (string, error) {
	if !viper.IsSet("browserless.url") {
		return "", errors.New("browserless.url not set")
	}
	browserlessURL := viper.GetString("browserless.url")
	data := struct {
		URL string `json:"url"`
	}{
		URL: fetchURL,
	}
	jsn, err := json.Marshal(data)
	if err != nil {
		log.Println("Could not marshal url:", err)
		filter.Log("Could not marshal url:", err)
		return "", err
	}
	browserlessURL = browserlessURL + "/content"
	resp, err := http.Post(browserlessURL, "application/json", bytes.NewBuffer(jsn))
	if err != nil {
		log.Println("Could not get browserless response content:", err)
		filter.Log("Could not get browserless response content:", err)
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Could not fetch url through browserless: ", fetchURL, " - ", err)
		filter.Log("Could not fetch url through browserless: ", fetchURL, " - ", err)
		return "", err
	}
	return string(body), nil
}

func getBrowserlessFunctionResult(filter *Filter) {
	result, err := getBrowserlessFunctionContent(filter, "")
	if err != nil {
		log.Println(err)
		filter.Log(err)
		return
	}
	filter.Results = append(filter.Results, result)
}

func getBrowserlessFunctionResults(filter *Filter) {
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			fetchURL := result

			str, err := getBrowserlessFunctionContent(filter, fetchURL)
			if err != nil {
				log.Println(err)
				filter.Log(err)
				continue
			}
			filter.Results = append(filter.Results, str)
		}
	}
}

type BrowserlessContext struct {
	Result string `json:"result"`
}

func getBrowserlessFunctionContent(filter *Filter, result string) (string, error) {
	if !viper.IsSet("browserless.url") {

		return "", errors.New("browserless.url not set")
	}
	browserlessURL := viper.GetString("browserless.url")
	if filter.Var2 == nil {
		return "", errors.New("filter.Var2 == nil")
	}
	code := *filter.Var2
	data := struct {
		Code     string             `json:"code"`
		Context  BrowserlessContext `json:"context"`
		Detached bool               `json:"detached"`
	}{
		Code:     code,
		Context:  BrowserlessContext{Result: result},
		Detached: false,
	}
	jsn, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	browserlessURL = browserlessURL + "/function"
	resp, err := http.Post(browserlessURL, "application/json", bytes.NewBuffer(jsn))
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func getFilterResultXPath(filter *Filter) {
	selectType := "node"
	if filter.Var2 != nil {
		selectType = *filter.Var2
	}
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			doc, err := htmlquery.Parse(strings.NewReader(result))
			if err != nil {
				filter.Log(err)
				continue
			}
			nodes, _ := htmlquery.QueryAll(doc, filter.Var1)
			for _, node := range nodes {
				switch selectType {
				case "inner":
					{
						if node.FirstChild == nil {
							continue
						}

						var result bytes.Buffer
						for child := node.FirstChild; child != nil; child = child.NextSibling {
							var b bytes.Buffer
							html.Render(&b, child)
							result.WriteString(b.String())
						}
						filter.Results = append(filter.Results, html.UnescapeString(result.String()))
						break
					}
				case "attr":
					{
						for _, attr := range node.Attr {
							result := fmt.Sprintf("%s=\"%s\"", attr.Key, attr.Val)
							filter.Results = append(filter.Results, html.UnescapeString(result))
						}
						break
					}
				default:
					{
						var b bytes.Buffer
						html.Render(&b, node)
						filter.Results = append(filter.Results, html.UnescapeString(b.String()))
						break
					}
				}
			}
		}
	}
}

func getFilterResultJSON(filter *Filter) {
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			for _, match := range gjson.Get(result, filter.Var1).Array() {
				filter.Results = append(filter.Results, match.String())
			}
		}
	}
}

func getFilterResultCSS(filter *Filter) {
	selectType := "node"
	if filter.Var2 != nil {
		selectType = *filter.Var2
	}
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			doc, err := html.Parse(strings.NewReader(result))
			if err != nil {
				filter.Log(err)
				continue
			}
			sel, err := cascadia.Parse(filter.Var1)
			if err != nil {
				filter.Log(err)
				continue
			}
			for _, node := range cascadia.QueryAll(doc, sel) {
				switch selectType {
				case "inner":
					{
						if node.FirstChild == nil {
							continue
						}

						var result bytes.Buffer
						for child := node.FirstChild; child != nil; child = child.NextSibling {
							var b bytes.Buffer
							html.Render(&b, child)
							result.WriteString(b.String())
						}
						filter.Results = append(filter.Results, html.UnescapeString(result.String()))
						break
					}
				case "attr":
					{
						for _, attr := range node.Attr {
							result := fmt.Sprintf("%s=\"%s\"", attr.Key, attr.Val)
							filter.Results = append(filter.Results, html.UnescapeString(result))
						}
						break
					}
				default:
					{
						var b bytes.Buffer
						html.Render(&b, node)
						filter.Results = append(filter.Results, html.UnescapeString(b.String()))
						break
					}
				}
			}
		}
	}
}

func getFilterResultReplace(filter *Filter) {
	r, err := regexp.Compile(filter.Var1)
	if err != nil {
		filter.Log(err)
		return
	}
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			if filter.Var2 == nil {
				filter.Results = append(filter.Results, r.ReplaceAllString(result, ""))
			} else {
				filter.Results = append(filter.Results, r.ReplaceAllString(result, *filter.Var2))
			}
		}
	}
}

func getFilterResultMatch(filter *Filter) {
	r, err := regexp.Compile(filter.Var1)
	if err != nil {
		filter.Log("Could not compile regex: ", err)
		return
	}
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			filter.Results = append(filter.Results, r.FindAllString(result, -1)...)
		}
	}
}

func getFilterResultSubstring(filter *Filter) {
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			substrings := strings.Split(filter.Var1, ",")
			var sb strings.Builder
			asRunes := []rune(result)

			for _, substring := range substrings {
				if strings.Contains(substring, ":") {
					from_to := strings.Split(substring, ":")
					if len(from_to) != 2 {
						filter.Log("Missing value in range: '", substring, "'")
						return
					}
					fromStr := from_to[0]
					var hasFrom bool = true
					if fromStr == "" {
						hasFrom = false
					}
					from64, err := strconv.ParseInt(fromStr, 10, 32)
					var from = int(from64)
					if hasFrom && err != nil {
						filter.Log("Could not parse left side of: '", substring, "'")
						return
					} else if from < 0 {
						from = len(asRunes) + from
					}
					if from < 0 {
						filter.Log("Out of bounds:", from_to)
						continue
					}

					toStr := from_to[1]
					var hasTo bool = true
					if toStr == "" {
						hasTo = false
					}
					to64, err := strconv.ParseInt(toStr, 10, 32)
					var to = int(to64)
					if hasTo && err != nil {
						filter.Log("Could not parse right side of: '", substring, "'")
						return
					} else if to < 0 {
						to = len(asRunes) + to
					}
					if to < 0 {
						filter.Log("Out of bounds:", from_to)
						continue
					}
					if hasFrom && hasTo {
						_, err := sb.WriteString(string(asRunes[from:to]))
						if err != nil {
							filter.Log("Could not substring: ", err)
						}
					} else if hasFrom {
						sb.WriteString(string(asRunes[from:]))
					} else if hasTo {
						sb.WriteString(string(asRunes[:to]))
					}
				} else {
					pos, err := strconv.ParseInt(substring, 10, 32)
					if err != nil || pos < 0 {
						filter.Log("Could not parse: '", substring, "'")
						return
					}
					sb.WriteRune(asRunes[pos])
				}
			}
			filter.Results = append(filter.Results, sb.String())
		}
	}
}

func getFilterResultContains(filter *Filter) {
	r, err := regexp.Compile(filter.Var1)
	invert, err := strconv.ParseBool(*filter.Var2)
	if err != nil {
		invert = false
	}

	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			contains := r.MatchString(result)
			if contains && !invert {
				filter.Results = append(filter.Results, result)
			} else if !contains && invert {
				filter.Results = append(filter.Results, result)
			}
		}
	}
}

func getFilterResultSum(filter *Filter) {
	var sum float64 = 0.0
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			if number, err := strconv.ParseFloat(result, 64); err == nil {
				sum += number
			} else {
				if len(result) > 50 {
					filter.Log("Could not convert value, with length ", len(result), ", to number")
				} else {
					filter.Log("Could not convert value, '", result, "', to number")
				}
			}
		}
	}
	filter.Results = append(filter.Results, fmt.Sprintf("%f", sum))
}

func getFilterResultMin(filter *Filter) {
	var min = math.MaxFloat64
	var setMin = false
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			if number, err := strconv.ParseFloat(result, 64); err == nil {
				if number < min {
					min = number
					setMin = true
				}
			} else {
				if len(result) > 50 {
					filter.Log("Could not convert value, '", result, "', to number")
					//filter.Log("Could not convert value, with length ", len(result), ", to number")
				} else {
					filter.Log("Could not convert value, '", result, "', to number")
				}
			}
		}
	}

	if setMin {
		filter.Results = append(filter.Results, fmt.Sprintf("%f", min))
	}
}

func getFilterResultMax(filter *Filter) {
	var max = -math.MaxFloat64
	var setMax = false
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			if number, err := strconv.ParseFloat(result, 64); err == nil {
				if number > max {
					max = number
					setMax = true
				}
			} else {
				if len(result) > 50 {
					filter.Log("Could not convert value, with length ", len(result), ", to number")
				} else {
					filter.Log("Could not convert value, '", result, "', to number")
				}
			}
		}
	}

	if setMax {
		filter.Results = append(filter.Results, fmt.Sprintf("%f", max))
	}
}

func getFilterResultAverage(filter *Filter) {
	var sum float64 = 0.0
	var count float64 = 0.0
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			if number, err := strconv.ParseFloat(result, 64); err == nil {
				sum += number
				count++
			} else {
				if len(result) > 50 {
					filter.Log("Could not convert value, with length ", len(result), ", to number")
				} else {
					filter.Log("Could not convert value, '", result, "', to number")
				}
			}
		}
	}
	if count > 0 {
		filter.Results = append(filter.Results, fmt.Sprintf("%f", sum/count))
	}
}

func getFilterResultCount(filter *Filter) {
	var count = 0
	for _, parent := range filter.Parents {
		count += len(parent.Results)
	}
	filter.Results = append(filter.Results, fmt.Sprintf("%d", count))
}

// https://gosamples.dev/round-float/
func roundFloat(val float64, precision uint) float64 {
	if precision == 0 {
		math.Round(val)
	}
	ratio := math.Pow(10, float64(precision))
	rounded := math.Round(val*ratio) / ratio
	return rounded
}

func getFilterResultRound(filter *Filter) {
	var decimals int64 = 0
	if filter.Var2 != nil {
		d, err := strconv.ParseInt(*filter.Var2, 10, 32)
		if err != nil {
			decimals = 0
		} else {
			decimals = d
		}
	}

	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			if number, err := strconv.ParseFloat(result, 64); err == nil {
				rounded := roundFloat(number, uint(decimals))
				filter.Results = append(filter.Results, fmt.Sprintf("%f", rounded))
			} else {
				if len(result) > 50 {
					filter.Log("Could not convert value, with length ", len(result), ", to number")
				} else {
					filter.Log("Could not convert value, '", result, "', to number")
				}
			}
		}
	}
}

func storeFilterResult(filter *Filter, db *gorm.DB, debug bool) {
	if debug {
		return
	}

	filterOutputs := make([]FilterOutput, 0, 1)
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			filterOutputs = append(filterOutputs, FilterOutput{
				WatchID: filter.WatchID,
				Name:    filter.Name,
				Time:    time.Now(),
				Value:   result,
			})
		}
	}
	if len(filterOutputs) > 0 {
		tx := db.Create(&filterOutputs)
		if tx.Error != nil {
			log.Println("Could not store value:", tx.Error)
		}
	}
}

func getFilterResultConditionDiff(filter *Filter, db *gorm.DB) {
	var previousOutput FilterOutput
	db.Model(&FilterOutput{}).Limit(1).Order("time desc").Where("watch_id = ? AND name = ?", filter.WatchID, filter.Var2).Find(&previousOutput)
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			if previousOutput.WatchID == 0 {
				filter.Results = append(filter.Results, result)
			} else if previousOutput.Value != result {
				filter.Results = append(filter.Results, result)
			}
		}
	}
}

func getFilterResultConditionLowerLast(filter *Filter, db *gorm.DB) {
	var previousOutput FilterOutput
	db.Model(&FilterOutput{}).Order("time desc").Where("watch_id = ? AND name = ?", filter.WatchID, filter.Var2).Limit(1).Find(&previousOutput)
	lastValue, lastValueErr := strconv.ParseFloat(previousOutput.Value, 64)
	if lastValueErr != nil {
		filter.Log("Could not convert previous value to number all will pass: '", previousOutput.Value, "'")
	}
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			if previousOutput.WatchID == 0 {
				filter.Results = append(filter.Results, result)
			} else {
				number, err := strconv.ParseFloat(result, 64)
				if err != nil {
					if len(result) > 50 {
						filter.Log("Could not convert value, with length ", len(result), ", to number")
					} else {
						filter.Log("Could not convert value, '", result, "', to number")
					}
					continue
				}
				if lastValueErr != nil || number < lastValue {
					filter.Results = append(filter.Results, result)
				}
			}
		}
	}
}

func getFilterResultConditionLowest(filter *Filter, db *gorm.DB) {
	var previousOutputs []FilterOutput
	db.Model(&FilterOutput{}).Where("watch_id = ? AND name = ?", filter.WatchID, filter.Name).Find(&previousOutputs)
	lowest := math.MaxFloat64
	for _, previousOutput := range previousOutputs {
		number, err := strconv.ParseFloat(previousOutput.Value, 64)
		if err != nil {
			continue
		}
		if number < lowest {
			lowest = number
		}
	}

	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			number, err := strconv.ParseFloat(result, 64)
			if err != nil {
				if len(result) > 50 {
					filter.Log("Could not convert value, with length ", len(result), ", to number")
				} else {
					filter.Log("Could not convert value, '", result, "', to number")
				}
				continue
			}
			if number < lowest {
				filter.Results = append(filter.Results, result)
			}
		}
	}
}

func getFilterResultConditionLowerThan(filter *Filter) {
	if filter.Var2 == nil {
		filter.Log("No threshold given")
		return
	}
	threshold, err := strconv.ParseFloat(*filter.Var2, 64)
	if err != nil {
		filter.Log("Could not convert convert threshold to number: '", *filter.Var2, "'")
		return
	}
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			number, err := strconv.ParseFloat(result, 64)
			if err != nil {
				if len(result) > 50 {
					filter.Log("Could not convert value, with length ", len(result), ", to number")
				} else {
					filter.Log("Could not convert value, '", result, "', to number")
				}
				continue
			}
			if number < threshold {
				filter.Results = append(filter.Results, result)
			}
		}
	}
}

func getFilterResultConditionHigherLast(filter *Filter, db *gorm.DB) {
	var previousOutput FilterOutput
	db.Model(&FilterOutput{}).Order("time desc").Where("watch_id = ? AND name = ?", filter.WatchID, filter.Var2).Limit(1).Find(&previousOutput)
	lastValue, lastValueErr := strconv.ParseFloat(previousOutput.Value, 64)
	if lastValueErr != nil {
		filter.Log("Could not convert previous value to number all will pass: '", previousOutput.Value, "'")
	}
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			if previousOutput.WatchID == 0 {
				filter.Results = append(filter.Results, result)
			} else {
				number, err := strconv.ParseFloat(result, 64)
				if err != nil {
					if len(result) > 50 {
						filter.Log("Could not convert value, with length ", len(result), ", to number")
					} else {
						filter.Log("Could not convert value, '", result, "', to number")
					}
					continue
				}
				if lastValueErr != nil || number > lastValue {
					filter.Results = append(filter.Results, result)
				}
			}
		}
	}
}

func getFilterResultConditionHighest(filter *Filter, db *gorm.DB) {
	var previousOutputs []FilterOutput
	db.Model(&FilterOutput{}).Where("watch_id = ? AND name = ?", filter.WatchID, filter.Var2).Find(&previousOutputs)
	highest := -math.MaxFloat64
	if previousOutputs != nil {
		for _, previousOutput := range previousOutputs {
			number, err := strconv.ParseFloat(previousOutput.Value, 64)
			if err != nil {
				continue
			}
			if number > highest {
				highest = number
			}
		}
	}

	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			number, err := strconv.ParseFloat(result, 64)
			if err != nil {
				if len(result) > 50 {
					filter.Log("Could not convert value, with length ", len(result), ", to number")
				} else {
					filter.Log("Could not convert value, '", result, "', to number")
				}
				continue
			}
			if number > highest {
				filter.Results = append(filter.Results, result)
			}
		}
	}
}

func getFilterResultConditionHigherThan(filter *Filter) {
	if filter.Var2 == nil {
		filter.Log("No threshold given for Higher Than Filter")
		return
	}
	threshold, err := strconv.ParseFloat(*filter.Var2, 64)
	if err != nil {
		filter.Log("Could not convert convert threshold to number: '", *filter.Var2, "'")
		return
	}
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			number, err := strconv.ParseFloat(result, 64)
			if err != nil {
				if len(result) > 50 {
					filter.Log("Could not convert value, with length ", len(result), ", to number")
				} else {
					filter.Log("Could not convert value, '", result, "', to number")
				}
				continue
			}
			if number > threshold {
				filter.Results = append(filter.Results, result)
			}
		}
	}
}

func getFilterResultUnique(filter *Filter) {
	valueMap := make(map[string]bool)

	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			valueMap[result] = true
		}
	}

	for value := range valueMap {
		filter.Results = append(filter.Results, value)
	}
}

func notifyFilter(filters []Filter, filter *Filter, watch *Watch, web *Web, debug bool) {
	haveResults := false
	for _, parent := range filter.Parents {
		if len(parent.Results) > 0 {
			haveResults = true
		}
	}
	if !debug && !haveResults {
		return
	}
	tmpl, err := template.New("notify").Parse(filter.Var1)
	if err != nil {
		filter.Log("Could not parse template: ", err)
		log.Println("Could not parse template: ", err)
		return
	}

	dataMap := make(map[string]any, 20)
	for _, f := range filters {
		dataMap[f.Name] = template.HTML(strings.Join(f.Results, ", "))
		dataMap[f.Name+"_Type"] = f.Type
		dataMap[f.Name+"_Var1"] = f.Var1
		dataMap[f.Name+"_Var2"] = f.Var2
	}

	dataMap["WatchName"] = template.HTML(watch.Name)

	var buffer bytes.Buffer
	tmpl.Execute(&buffer, dataMap)
	if debug {
		filter.Results = append(filter.Results, buffer.String())
	} else {
		notifier := filter.Var2
		if notifier == nil {
			return
		}
		web.notify(*notifier, buffer.String())
	}

}

func TriggerSchedule(watchID uint, web *Web, scheduleID *uint) {
	var watch *Watch
	web.db.Model(&Watch{}).First(&watch, watchID)

	var filters []Filter
	web.db.Model(&Filter{}).Where("watch_id = ?", watch.ID).Find(&filters)

	var connections []FilterConnection
	web.db.Model(&FilterConnection{}).Where("watch_id = ?", watch.ID).Find(&connections)

	buildFilterTree(filters, connections)
	ProcessFilters(filters, web, watch, false, scheduleID)
}

func getCronDebugResult(filter *Filter) {
	_, err := cron.ParseStandard(filter.Var1)
	if err != nil {
		filter.Log(err)
	}
}

func getFilterResultLua(filter *Filter) {
	L := lua.NewState()
	defer L.Close()

	lualibs.Preload(L)

	inputs := L.CreateTable(10, 0)
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			inputs.Append(lua.LString(result))
		}
	}
	L.SetGlobal("inputs", inputs)

	outputs := L.CreateTable(10, 0)
	L.SetGlobal("outputs", outputs)

	logs := L.CreateTable(10, 0)
	L.SetGlobal("logs", logs)

	err := L.DoString(filter.Var1)
	if err != nil {
		filter.Log(err)
		return
	}
	outputs.ForEach(
		func(key lua.LValue, value lua.LValue) {
			filter.Results = append(filter.Results, value.String())
		},
	)
	logs.ForEach(
		func(key lua.LValue, value lua.LValue) {
			filter.Logs = append(filter.Logs, value.String())
		},
	)
}

func getFilterResultEcho(filter *Filter) {
	filter.Results = append(filter.Results, filter.Var1)
}
