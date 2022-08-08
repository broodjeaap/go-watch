package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/andybalholm/cascadia"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

func getGroupResult(group *FilterGroup) []string {
	resp, err := http.Get(group.URL.URL)
	if err != nil {
		log.Print("Something went wrong loading", group.URL.URL)
		return []string{}
	}
	defer resp.Body.Close()
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("Something went wrong loading ", group.URL.URL)
		return []string{}
	}
	resultStrings := []string{string(html)}
	newStrings := []string{}
	for _, filter := range group.Filters {
		for _, resultString := range resultStrings {
			getFilterResult(resultString, &filter, &newStrings)
		}
		log.Println(len(resultStrings), len(newStrings))
		resultStrings = newStrings
		newStrings = nil
	}
	return resultStrings
}

func getFilterResult(s string, filter *Filter, newStrings *[]string) {
	switch {
	case filter.Type == "xpath":
		{
			getFilterResultXPath(s, filter, newStrings)
		}
	case filter.Type == "css":
		{
			getFilterResultCSS(s, filter, newStrings)
		}
	case filter.Type == "replace":
		{
			getFilterResultReplace(s, filter, newStrings)
		}
	case filter.Type == "regex":
		{
			getFilterResultMatch(s, filter, newStrings)
		}
	case filter.Type == "substring":
		{
			getFilterResultSubstring(s, filter, newStrings)
		}
	default:

	}
}

func getFilterResultXPath(s string, filter *Filter, newStrings *[]string) {
	doc, err := htmlquery.Parse(strings.NewReader(s))
	if err != nil {
		log.Print(err)
		return
	}
	nodes, _ := htmlquery.QueryAll(doc, filter.From)
	for _, node := range nodes {
		var b bytes.Buffer
		html.Render(&b, node)
		*newStrings = append(*newStrings, html.UnescapeString(b.String()))
	}
}

func getFilterResultCSS(s string, filter *Filter, newStrings *[]string) {
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		log.Print(err)
		return
	}
	sel, err := cascadia.Parse(filter.From)
	if err != nil {
		log.Print(err)
		return
	}
	for _, node := range cascadia.QueryAll(doc, sel) {
		var b bytes.Buffer
		html.Render(&b, node)
		*newStrings = append(*newStrings, html.UnescapeString(b.String()))
	}
}

func getFilterResultReplace(s string, filter *Filter, newStrings *[]string) {
	regex, err := regexp.Compile(filter.From)
	if err != nil {
		log.Print(err)
		return
	}
	*newStrings = append(*newStrings, regex.ReplaceAllString(filter.From, filter.To))
}

func getFilterResultMatch(s string, filter *Filter, newStrings *[]string) {
	r, err := regexp.Compile(filter.From)
	if err != nil {
		log.Print(err)
		return
	}
	*newStrings = append(*newStrings, r.ReplaceAllString(s, filter.To))
}

func getFilterResultSubstring(s string, filter *Filter, newStrings *[]string) {
	substrings := strings.Split(filter.From, ",")
	var sb strings.Builder
	asRunes := []rune(s)

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
	*newStrings = append(*newStrings, sb.String())
}
