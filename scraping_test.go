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
			want := []string{}
			getFilterResultXPath(
				HTML_STRING,
				&Filter{
					From: test.Query,
				},
				&want,
			)
			if !reflect.DeepEqual(test.Want, want) {
				t.Errorf("Got %s, want %s", want, test.Want)
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
			want := []string{}
			getFilterResultCSS(
				HTML_STRING,
				&Filter{
					From: test.Query,
				},
				&want,
			)
			if !reflect.DeepEqual(test.Want, want) {
				t.Errorf("Got %s, want %s", want, test.Want)
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
			want := []string{test.Want}
			getFilterResultSubstring(
				test.Input,
				&Filter{
					From: test.Query,
				},
				&want,
			)
			if want[0] != test.Want {
				t.Errorf("Got %s, want %s", want[0], test.Want)
			}
		})
	}
}
