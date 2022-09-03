package main

import (
	"fmt"
	"reflect"
	"testing"
)

const HTML_STRING = `<html>
<head>
	<title>title</title>
</head>
<body>
	<table class="product-table" id="product-table">
		<caption>product-table-caption</caption>
		<thead>
			<tr>
				<th>Name</th>
				<th>Stock</th>
				<th>Price</th>
			</tr>
		</thead>
		<tbody>
			<tr><td class="product">product1</td><td class="stock">10</td><td class="price">100</td></tr>
			<tr><td class="product">product2</td><td class="stock">20</td><td class="price">200</td></tr>
			<tr><td class="product">product3</td><td class="stock">30</td><td class="price">300</td></tr>
			<tr><td class="product">product4</td><td class="stock">40</td><td class="price">400</td></tr>
		</tbody>
	</table>
</body>
</html>`

const JSON_STRING = `{
	"date": "1970-01-01",
	"products": [
		{"name": "product1", "stock": 10, "price": 100},
		{"name": "product2", "stock": 20, "price": 200},
		{"name": "product3", "stock": 30, "price": 300},
		{"name": "product4", "stock": 40, "price": 400}
	]
}`

func TestFilterXPath(t *testing.T) {
	var tests = []struct {
		Query string
		Want  []string
	}{
		{"//title", []string{"<title>title</title>"}},
		{"//table[@id='product-table']//tr//td[last()]", []string{`<td class="price">100</td>`, `<td class="price">200</td>`, `<td class="price">300</td>`, `<td class="price">400</td>`}},
		{"//td[@class='price']", []string{`<td class="price">100</td>`, `<td class="price">200</td>`, `<td class="price">300</td>`, `<td class="price">400</td>`}},
		{"//table[@id='product-table']//tr//td[2]", []string{`<td class="stock">10</td>`, `<td class="stock">20</td>`, `<td class="stock">30</td>`, `<td class="stock">40</td>`}},
		{"//td[@class='stock']", []string{`<td class="stock">10</td>`, `<td class="stock">20</td>`, `<td class="stock">30</td>`, `<td class="stock">40</td>`}},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%s", test.Query)
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parent: &Filter{
					Results: []string{HTML_STRING},
				},
				Var1: test.Query,
			}
			getFilterResultXPath(
				&filter,
			)
			if !reflect.DeepEqual(test.Want, filter.Results) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}

func TestFilterJSON(t *testing.T) {
	var tests = []struct {
		Query string
		Want  []string
	}{
		{"date", []string{"1970-01-01"}},
		{"products.#.name", []string{"product1", "product2", "product3", "product4"}},
		{"products.#.stock", []string{"10", "20", "30", "40"}},
		{"products.#.price", []string{"100", "200", "300", "400"}},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%s", test.Query)
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parent: &Filter{
					Results: []string{JSON_STRING},
				},
				Var1: test.Query,
			}
			getFilterResultJSON(
				&filter,
			)
			if !reflect.DeepEqual(test.Want, filter.Results) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}

func TestFilterCSS(t *testing.T) {
	var tests = []struct {
		Query string
		Want  []string
	}{
		{"title", []string{"<title>title</title>"}},
		{".product-table tr td:last-child", []string{`<td class="price">100</td>`, `<td class="price">200</td>`, `<td class="price">300</td>`, `<td class="price">400</td>`}},
		{".price", []string{`<td class="price">100</td>`, `<td class="price">200</td>`, `<td class="price">300</td>`, `<td class="price">400</td>`}},
		{".product-table tr td:nth-child(2)", []string{`<td class="stock">10</td>`, `<td class="stock">20</td>`, `<td class="stock">30</td>`, `<td class="stock">40</td>`}},
		{".stock", []string{`<td class="stock">10</td>`, `<td class="stock">20</td>`, `<td class="stock">30</td>`, `<td class="stock">40</td>`}},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%s", test.Query)
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parent: &Filter{
					Results: []string{HTML_STRING},
				},
				Var1: test.Query,
			}
			getFilterResultCSS(
				&filter,
			)
			if !reflect.DeepEqual(test.Want, filter.Results) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}

func TestFilterReplace(t *testing.T) {
	var tests = []struct {
		Input string
		Query string
		Want  string
	}{
		{"0123456789", "0", "123456789"},
		{"0123456789", "9", "012345678"},
		{"0123456789", "3456", "012789"},
		{"0123456789_0123456789", "3456", "012789_012789"},
		{"世界日本語", "世", "界日本語"},
		{"世界日本語", "語", "世界日本"},
		{"世界日_世界日_世界日", "界", "世日_世日_世日"},
		// TODO add replace tests
		// TODO add regex tests
		// TODO add regex replace tests
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%s %s", test.Input, test.Query)
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parent: &Filter{
					Results: []string{test.Input},
				},
				Var1: test.Query,
			}
			getFilterResultReplace(
				&filter,
			)
			if filter.Results[0] != test.Want {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}

func TestFilterMatch(t *testing.T) {
	var tests = []struct {
		Input string
		Query string
		Want  []string
	}{
		{"0123456789", "0123|789", []string{"0123", "789"}},
		{"0123456789", "abc|321", []string{}},
		{"0123456789", "[0-9]{3}", []string{"012", "345", "678"}},
		{"世界日本語", "日本", []string{"日本"}},
		{"世界日本語_世界日本語_世界日本語", "日本", []string{"日本", "日本", "日本"}},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%s", test.Query)
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parent: &Filter{
					Results: []string{test.Input},
				},
				Var1: test.Query,
			}
			getFilterResultMatch(
				&filter,
			)
			// len() thing cuz filterResults == nil and test.Want == [], same thing but not really...
			if !(len(filter.Results) == 0 && len(test.Want) == 0) && !reflect.DeepEqual(filter.Results, test.Want) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}

func TestFilterSubstring(t *testing.T) {
	var tests = []struct {
		Input string
		Query string
		Want  string
	}{
		{"0123456789", "0", "0"},
		{"0123456789", "9", "9"},
		{"0123456789", "0,9", "09"},
		{"0123456789", "0:3", "012"},
		{"0123456789", ":3", "012"},
		{"0123456789", "3:", "3456789"},
		{"0123456789", "-3:", "789"},
		{"0123456789", ":-3", "0123456"},
		{"0123456789", ":-1", "012345678"},

		{"0123456789", "0,3,7,9", "0379"},
		{"0123456789", "0:3,7,9", "01279"},

		{"世界", "1", "界"},
		{"世界日本語", ":3", "世界日"},
		{"世界日本語", ":-1", "世界日本"},
		{"世界日本語", "-1:", "語"},

		{"0123456789", "A", "0123456789"},
		{"0123456789", "1:A", "0123456789"},
		{"0123456789", "A:1", "0123456789"},
		{"0123456789", "A:B", "0123456789"},
		{"0123456789", "A:B:C", "0123456789"},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%s %s", test.Input, test.Query)
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parent: &Filter{
					Results: []string{test.Input},
				},
				Var1: test.Query,
			}
			getFilterResultSubstring(
				&filter,
			)
			if filter.Results[0] != test.Want {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}
