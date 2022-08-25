package main

import (
	"bytes"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/andybalholm/cascadia"
	"github.com/antchfx/htmlquery"
	"github.com/tidwall/gjson"
	"golang.org/x/net/html"
)

func getFilterResults(filter *Filter) {
	getFilterResult(filter)
	for _, filter := range filter.Filters {
		getFilterResults(&filter)
	}
}

func getFilterResult(filter *Filter) {
	switch {
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
	default:

	}
}

func getFilterResultXPath(filter *Filter) {
	if filter.Parent == nil {
		log.Println("Filter", filter.Name, "called without parent for", filter.Type)
		return
	}
	for _, result := range filter.Parent.results {
		doc, err := htmlquery.Parse(strings.NewReader(result))
		if err != nil {
			log.Print(err)
			continue
		}
		nodes, _ := htmlquery.QueryAll(doc, filter.Var1)
		for _, node := range nodes {
			var b bytes.Buffer
			html.Render(&b, node)
			filter.results = append(filter.results, html.UnescapeString(b.String()))
		}
	}
}

func getFilterResultJSON(filter *Filter) {
	if filter.Parent == nil {
		log.Println("Filter", filter.Name, "called without parent for", filter.Type)
		return
	}
	for _, result := range filter.Parent.results {
		for _, match := range gjson.Get(result, filter.Var1).Array() {
			filter.results = append(filter.results, match.String())
		}
	}
}

func getFilterResultCSS(filter *Filter) {
	if filter.Parent == nil {
		log.Println("Filter", filter.Name, "called without parent for", filter.Type)
		return
	}
	for _, result := range filter.results {
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
			filter.results = append(filter.results, html.UnescapeString(b.String()))
		}
	}
}

func getFilterResultReplace(filter *Filter) {
	if filter.Parent == nil {
		log.Println("Filter", filter.Name, "called without parent for", filter.Type)
		return
	}
	for _, result := range filter.results {
		r, err := regexp.Compile(filter.Var1)
		if err != nil {
			log.Print(err)
			continue
		}
		if filter.Var2 == nil {
			filter.results = append(filter.results, r.ReplaceAllString(result, ""))
		} else {
			filter.results = append(filter.results, r.ReplaceAllString(result, *filter.Var2))
		}
	}
}

func getFilterResultMatch(filter *Filter) {
	if filter.Parent == nil {
		log.Println("Filter", filter.Name, "called without parent for", filter.Type)
		return
	}
	for _, result := range filter.results {
		r, err := regexp.Compile(filter.Var1)
		if err != nil {
			log.Print(err)
			continue
		}
		for _, str := range r.FindAllString(result, -1) {
			filter.results = append(filter.results, str)
		}
	}
}

func getFilterResultSubstring(filter *Filter) {
	if filter.Parent == nil {
		log.Println("Filter", filter.Name, "called without parent for", filter.Type)
		return
	}
	for _, result := range filter.results {
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
		filter.results = append(filter.results, sb.String())
	}
}
