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

	"github.com/andybalholm/cascadia"
	"github.com/antchfx/htmlquery"
	"github.com/tidwall/gjson"
	"golang.org/x/net/html"
)

func fillFilterResults(filters []Filter) {
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
		getFilterResult(filter)
		processedMap[filter.ID] = true
	}
}

func getFilterResult(filter *Filter) {
	switch {
	case filter.Type == "gurl":
		{
			getFilterResultURL(filter)
		}
	case filter.Type == "gurls":
		{
			getFilterResultURL(filter)
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
	case filter.Type == "min":
		{
			getFilterResultMin(filter)
		}
	case filter.Type == "max":
		{
			getFilterResultMax(filter)
		}
	case filter.Type == "average":
		{
			getFilterResultAverage(filter)
		}
	case filter.Type == "count":
		{
			getFilterResultCount(filter)
		}
	default:
		log.Println("getFilterResult called with filter.Type == ", filter.Type)
	}
}

func getFilterResultURL(filter *Filter) {
	url := filter.Var1
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
	filter.Results = append(filter.Results, string(body))
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
	filter.Results = append(filter.Results, fmt.Sprintf("%f", sum/count))
}

func getFilterResultCount(filter *Filter) {
	var count = 0
	for _, parent := range filter.Parents {
		count += len(parent.Children)
	}
	log.Println(fmt.Sprintf("%d", count))
	filter.Results = append(filter.Results, fmt.Sprintf("%d", count))
}
