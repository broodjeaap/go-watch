package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/andybalholm/cascadia"
	"github.com/antchfx/htmlquery"
	"github.com/tidwall/gjson"
	"golang.org/x/net/html"
	"gorm.io/gorm"
)

func processFilters(filters []Filter, db *gorm.DB, urlCache map[string]string, useCache bool, setCache bool) {
	processedMap := make(map[uint]bool, len(filters))
	for len(filters) > 0 {
		filter := &filters[0]
		filters = filters[1:]
		var allParentsProcessed = true
		for _, parent := range filter.Parents {
			if _, contains := processedMap[parent.ID]; !contains {
				allParentsProcessed = false
				break
			}
		}
		if !allParentsProcessed {
			filters = append(filters, *filter)
			continue
		}
		getFilterResult(filter, db, urlCache, useCache, setCache)
		processedMap[filter.ID] = true
	}
}

func getFilterResult(filter *Filter, db *gorm.DB, urlCache map[string]string, useCache bool, setCache bool) {
	switch {
	case filter.Type == "gurl":
		{
			getFilterResultURL(filter, urlCache, useCache, setCache)
		}
	case filter.Type == "gurls":
		{
			getFilterResultURL(filter, urlCache, useCache, setCache)
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
			storeFilterResult(filter, db)
		}
	case filter.Type == "condition":
		{
			switch filter.Var1 {
			case "diff":
				{
					getFilterResultConditionDiff(filter, db)
				}
			case "lowerl":
				{
					getFilterResultConditionLowerLast(filter, db)
				}
			case "lowest":
				{
					getFilterResultConditionLowest(filter, db)
				}
			case "lowert":
				{
					getFilterResultConditionLowerThan(filter)
				}
			case "higherl":
				{
					getFilterResultConditionHigherLast(filter, db)
				}
			case "highest":
				{
					getFilterResultConditionHighest(filter, db)
				}
			case "highert":
				{
					getFilterResultConditionHigherThan(filter)
				}
			}
		}
	default:
		log.Println("getFilterResult called with filter.Type == ", filter.Type)
	}
}

func getFilterResultURL(filter *Filter, urlCache map[string]string, useCache bool, setCache bool) {
	url := filter.Var1
	val, exists := urlCache[url]
	if useCache && exists {
		filter.Results = append(filter.Results, val)
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Could not fetch url", url)
		log.Println("Reason:", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Could not fetch url", url)
		log.Println("Reason:", err)
	}
	str := string(body)
	filter.Results = append(filter.Results, str)
	if setCache {
		urlCache[url] = str
	}
}

func getFilterResultXPath(filter *Filter) {
	if filter.Parents == nil {
		log.Println("Filter", filter.Name, "called without parents for", filter.Type)
		return
	}
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			doc, err := htmlquery.Parse(strings.NewReader(result))
			if err != nil {
				log.Print(err)
				continue
			}
			nodes, _ := htmlquery.QueryAll(doc, filter.Var1)
			for _, node := range nodes {
				var b bytes.Buffer
				html.Render(&b, node)
				filter.Results = append(filter.Results, html.UnescapeString(b.String()))
			}
		}
	}
}

func getFilterResultJSON(filter *Filter) {
	if filter.Parents == nil {
		log.Println("Filter", filter.Name, "called without parent for", filter.Type)
		return
	}
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			for _, match := range gjson.Get(result, filter.Var1).Array() {
				filter.Results = append(filter.Results, match.String())
			}
		}
	}
}

func getFilterResultCSS(filter *Filter) {
	if filter.Parents == nil {
		log.Println("Filter", filter.Name, "called without parent for", filter.Type)
		return
	}
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			doc, err := html.Parse(strings.NewReader(result))
			if err != nil {
				log.Print(err)
				continue
			}
			sel, err := cascadia.Parse(filter.Var1)
			if err != nil {
				log.Print(err)
				continue
			}
			for _, node := range cascadia.QueryAll(doc, sel) {
				var b bytes.Buffer
				html.Render(&b, node)
				filter.Results = append(filter.Results, html.UnescapeString(b.String()))
			}
		}
	}
}

func getFilterResultReplace(filter *Filter) {
	if filter.Parents == nil {
		log.Println("Filter", filter.Name, "called without parent for", filter.Type)
		return
	}
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			r, err := regexp.Compile(filter.Var1)
			if err != nil {
				log.Print(err)
				continue
			}
			if filter.Var2 == nil {
				filter.Results = append(filter.Results, r.ReplaceAllString(result, ""))
			} else {
				filter.Results = append(filter.Results, r.ReplaceAllString(result, *filter.Var2))
			}
		}
	}
}

func getFilterResultMatch(filter *Filter) {
	if filter.Parents == nil {
		log.Println("Filter", filter.Name, "called without parent for", filter.Type)
		return
	}
	r, err := regexp.Compile(filter.Var1)
	if err != nil {
		log.Print(err)
		return
	}
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			for _, str := range r.FindAllString(result, -1) {
				filter.Results = append(filter.Results, str)
			}
		}
	}
}

func getFilterResultSubstring(filter *Filter) {
	if filter.Parents == nil {
		log.Println("Filter", filter.Name, "called without parent for", filter.Type)
		return
	}
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			substrings := strings.Split(filter.Var1, ",")
			var sb strings.Builder
			asRunes := []rune(result)

			for _, substring := range substrings {
				if strings.Contains(substring, ":") {
					from_to := strings.Split(substring, ":")
					if len(from_to) != 2 {
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
						return
					} else if from < 0 {
						from = len(asRunes) + from
					}
					toStr := from_to[1]
					var hasTo bool = true
					if toStr == "" {
						hasTo = false
					}
					to64, err := strconv.ParseInt(toStr, 10, 32)
					var to = int(to64)
					if hasTo && err != nil {
						return
					} else if to < 0 {
						to = len(asRunes) + to
					}
					if hasFrom && hasTo {
						sb.WriteString(string(asRunes[from:to]))
					} else if hasFrom {
						sb.WriteString(string(asRunes[from:]))
					} else if hasTo {
						sb.WriteString(string(asRunes[:to]))
					}
				} else {
					pos, err := strconv.ParseInt(substring, 10, 32)
					if err != nil || pos < 0 {
						return
					}
					sb.WriteRune(asRunes[pos])
				}
			}
			filter.Results = append(filter.Results, sb.String())
		}
	}
}

func getFilterResultSum(filter *Filter) {
	var sum float64 = 0.0
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			if number, err := strconv.ParseFloat(result, 64); err == nil {
				sum += number
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
	log.Println(val, precision, ratio, rounded)
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
			}
		}
	}
}

func storeFilterResult(filter *Filter, db *gorm.DB) {
	var previousOutput FilterOutput
	db.Model(&FilterOutput{}).Order("time desc").Where("watch_id = ? AND name = ?", filter.WatchID, filter.Name).Limit(1).Find(&previousOutput)

	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			if previousOutput.WatchID == 0 {
				previousOutput.Name = filter.Name
				previousOutput.Time = time.Now()
				previousOutput.Value = result
				previousOutput.WatchID = filter.WatchID
				db.Create(&previousOutput)
			} else {
				previousOutput.Time = time.Now()
				previousOutput.ID = 0
				previousOutput.Value = result
				db.Create(&previousOutput)
			}
		}
	}
}

func getFilterResultConditionDiff(filter *Filter, db *gorm.DB) {
	var previousOutput FilterOutput
	db.Model(&FilterOutput{}).Order("time desc").Where("watch_id = ? AND name = ?", filter.WatchID, filter.Name).Limit(1).Find(&previousOutput)
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
	db.Model(&FilterOutput{}).Order("time desc").Where("watch_id = ? AND name = ?", filter.WatchID, filter.Name).Limit(1).Find(&previousOutput)
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			if previousOutput.WatchID == 0 {
				filter.Results = append(filter.Results, result)
			} else {
				lastValue, err := strconv.ParseFloat(previousOutput.Value, 64)
				if err != nil {
					log.Println("Could not convert previous value to number:", previousOutput.Value)
					continue
				}
				number, err := strconv.ParseFloat(result, 64)
				if err != nil {
					log.Println("Could not convert new value to number:", result)
					continue
				}
				if number < lastValue {
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
	if previousOutputs != nil {
		for _, previousOutput := range previousOutputs {
			number, err := strconv.ParseFloat(previousOutput.Value, 64)
			if err != nil {
				log.Println("Could not convert result to number:", previousOutput.Value)
				continue
			}
			if number < lowest {
				lowest = number
			}
		}
		return
	}

	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			number, err := strconv.ParseFloat(result, 64)
			if err != nil {
				log.Println("Could not convert result to number:", result)
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
		log.Println("No threshold given for Lower Than Filter")
		return
	}
	threshold, err := strconv.ParseFloat(*filter.Var2, 64)
	if err != nil {
		log.Println("Could not convert convert threshold to number:", *filter.Var2)
		return
	}
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			number, err := strconv.ParseFloat(result, 64)
			if err != nil {
				log.Println("Could not convert new value to number:", result)
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
	db.Model(&FilterOutput{}).Order("time desc").Where("watch_id = ? AND name = ?", filter.WatchID, filter.Name).Limit(1).Find(&previousOutput)
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			if previousOutput.WatchID == 0 {
				filter.Results = append(filter.Results, result)
			} else {
				lastValue, err := strconv.ParseFloat(previousOutput.Value, 64)
				if err != nil {
					log.Println("Could not convert previous value to number:", previousOutput.Value)
					continue
				}
				number, err := strconv.ParseFloat(result, 64)
				if err != nil {
					log.Println("Could not convert new value to number:", result)
					continue
				}
				if number > lastValue {
					filter.Results = append(filter.Results, result)
				}
			}
		}
	}
}

func getFilterResultConditionHighest(filter *Filter, db *gorm.DB) {
	var previousOutputs []FilterOutput
	db.Model(&FilterOutput{}).Where("watch_id = ? AND name = ?", filter.WatchID, filter.Name).Find(&previousOutputs)
	highest := math.MaxFloat64
	if previousOutputs != nil {
		for _, previousOutput := range previousOutputs {
			number, err := strconv.ParseFloat(previousOutput.Value, 64)
			if err != nil {
				log.Println("Could not convert result to number:", previousOutput.Value)
				continue
			}
			if number > highest {
				highest = number
			}
		}
		return
	}

	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			number, err := strconv.ParseFloat(result, 64)
			if err != nil {
				log.Println("Could not convert result to number:", result)
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
		log.Println("No threshold given for Higher Than Filter")
		return
	}
	threshold, err := strconv.ParseFloat(*filter.Var2, 64)
	if err != nil {
		log.Println("Could not convert convert threshold to number:", *filter.Var2)
		return
	}
	for _, parent := range filter.Parents {
		for _, result := range parent.Results {
			number, err := strconv.ParseFloat(result, 64)
			if err != nil {
				log.Println("Could not convert new value to number:", result)
				continue
			}
			if number > threshold {
				filter.Results = append(filter.Results, result)
			}
		}
	}
}

func notifyFilter(filter *Filter, db *gorm.DB) {

}
