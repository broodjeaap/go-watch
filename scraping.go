package main

import (
	"bytes"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

func getQueryResult(query *Query) []string {
	doc, err := htmlquery.LoadURL(query.URL.URL)
	if err != nil {
		log.Print("Something went wrong loading loading", query.URL.URL)
	}
	nodes, _ := htmlquery.QueryAll(doc, query.Query)
	nodeStrings := make([]string, len(nodes))
	for i, node := range nodes {
		var b bytes.Buffer
		html.Render(&b, node)
		nodeStrings[i] = b.String()
	}
	for _, filter := range query.Filters {
		for i, nodeString := range nodeStrings {
			nodeStrings[i] = getFilterResult(nodeString, &filter)
		}
	}
	return nodeStrings
}

func getFilterResult(s string, filter *Filter) string {
	switch {
	case filter.Type == "replace":
		{
			return getFilterResultReplace(s, filter)
		}
	case filter.Type == "regex":
		{
			return getFilterResultRegex(s, filter)
		}
	case filter.Type == "substring":
		{
			return getFilterResultSubstring(s, filter)
		}
	default:
		return s
	}
}

func getFilterResultReplace(s string, filter *Filter) string {
	return strings.ReplaceAll(s, filter.From, filter.To)
}

func getFilterResultRegex(s string, filter *Filter) string {
	regex, err := regexp.Compile(filter.From)
	if err != nil {
		return s
	}
	return regex.ReplaceAllString(s, filter.To)
}

func getFilterResultSubstring(s string, filter *Filter) string {
	substrings := strings.Split(filter.From, ",")
	var sb strings.Builder
	asRunes := []rune(s)

	for _, substring := range substrings {
		if strings.Contains(substring, "-") {
			from_to := strings.Split(substring, "-")
			if len(from_to) != 2 {
				return s
			}
			from, err := strconv.ParseInt(from_to[0], 10, 32)
			if err != nil {
				return s
			}
			to, err := strconv.ParseInt(from_to[1], 10, 32)
			if err != nil {
				return s
			}
			sb.WriteString(string(asRunes[from:to]))
		} else {
			pos, err := strconv.ParseInt(substring, 10, 32)
			if err != nil {
				return s
			}
			sb.WriteRune(asRunes[pos])
		}
	}
	return sb.String()
}
