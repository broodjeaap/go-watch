package web

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	. "github.com/broodjeaap/go-watch/models"
)

const HTML_STRING = `<html>
<head>
	<title>title</title>
</head>
<body>
	<table class="product-table" id="product-table">
		<caption class="h3" id="table-caption" data="data">product-table-caption</caption>
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
<div id="empty-div"></div>
<div id="multiple-children-div"><div id="first-child"></div><div id="second-child"></div></div>
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

func DeepEqualStringSlice(a []string, b []string) bool {
	return len(a) == len(b) && (len(a) == 0 || reflect.DeepEqual(a, b))
}

func TestFilterXPathNode(t *testing.T) {
	var2 := "node"
	var tests = []struct {
		Query string
		Want  []string
	}{
		{"//title", []string{"<title>title</title>"}},
		{"//table[@id='product-table']/caption", []string{`<caption class="h3" id="table-caption" data="data">product-table-caption</caption>`}},
		{"//table[@id='product-table']//tr//td[last()]", []string{`<td class="price">100</td>`, `<td class="price">200</td>`, `<td class="price">300</td>`, `<td class="price">400</td>`}},
		{"//td[@class='price']", []string{`<td class="price">100</td>`, `<td class="price">200</td>`, `<td class="price">300</td>`, `<td class="price">400</td>`}},
		{"//table[@id='product-table']//tr//td[2]", []string{`<td class="stock">10</td>`, `<td class="stock">20</td>`, `<td class="stock">30</td>`, `<td class="stock">40</td>`}},
		{"//td[@class='stock']", []string{`<td class="stock">10</td>`, `<td class="stock">20</td>`, `<td class="stock">30</td>`, `<td class="stock">40</td>`}},
		{"//div[@id='empty-div']", []string{`<div id="empty-div"></div>`}},
		{"//div[@id='multiple-children-div']", []string{`<div id="multiple-children-div"><div id="first-child"></div><div id="second-child"></div></div>`}},
		{"//div[@id='does-not-exist']", []string{}},
	}

	for _, test := range tests {
		testname := test.Query
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parents: []*Filter{
					{Results: []string{HTML_STRING}},
				},
				Var1: test.Query,
				Var2: &var2,
			}
			getFilterResultXPath(
				&filter,
			)
			if !DeepEqualStringSlice(filter.Results, test.Want) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}

func TestFilterXPathInnerHTML(t *testing.T) {
	var2 := "inner"
	var tests = []struct {
		Query string
		Want  []string
	}{
		{"//title", []string{"title"}},
		{"//table[@id='product-table']/caption", []string{`product-table-caption`}},
		{"//table[@id='product-table']//tr//td[last()]", []string{`100`, `200`, `300`, `400`}},
		{"//td[@class='price']", []string{`100`, `200`, `300`, `400`}},
		{"//table[@id='product-table']//tr//td[2]", []string{`10`, `20`, `30`, `40`}},
		{"//td[@class='stock']", []string{`10`, `20`, `30`, `40`}},
		{"//div[@id='empty-div']", []string{}},
		{"//div[@id='multiple-children-div']", []string{`<div id="first-child"></div><div id="second-child"></div>`}},
	}

	for _, test := range tests {
		testname := test.Query
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parents: []*Filter{
					{Results: []string{HTML_STRING}},
				},
				Var1: test.Query,
				Var2: &var2,
			}
			getFilterResultXPath(
				&filter,
			)
			if !DeepEqualStringSlice(filter.Results, test.Want) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}

func TestFilterXPathAttributes(t *testing.T) {
	var2 := "attr"
	var tests = []struct {
		Query string
		Want  []string
	}{
		{"//title", []string{}},
		{"//table[@id='product-table']/caption", []string{`class="h3"`, `id="table-caption"`, `data="data"`}},
		{"//table[@id='product-table']//tr//td[last()]", []string{`class="price"`, `class="price"`, `class="price"`, `class="price"`}},
		{"//td[@class='price']", []string{`class="price"`, `class="price"`, `class="price"`, `class="price"`}},
		{"//table[@id='product-table']//tr//td[2]", []string{`class="stock"`, `class="stock"`, `class="stock"`, `class="stock"`}},
		{"//td[@class='stock']", []string{`class="stock"`, `class="stock"`, `class="stock"`, `class="stock"`}},
		{"//*[@class='does-not-exists']", []string{}},
	}

	for _, test := range tests {
		testname := test.Query
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parents: []*Filter{
					{Results: []string{HTML_STRING}},
				},
				Var1: test.Query,
				Var2: &var2,
			}
			getFilterResultXPath(
				&filter,
			)
			if !DeepEqualStringSlice(filter.Results, test.Want) {
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
		{"does.not.exist", []string{}},
	}

	for _, test := range tests {
		testname := test.Query
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parents: []*Filter{
					{Results: []string{JSON_STRING}},
				},
				Var1: test.Query,
			}
			getFilterResultJSON(
				&filter,
			)
			if !DeepEqualStringSlice(filter.Results, test.Want) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}

func TestFilterCSSNode(t *testing.T) {
	var2 := "node"
	var tests = []struct {
		Query string
		Want  []string
	}{
		{"title", []string{"<title>title</title>"}},
		{".product-table tr td:last-child", []string{`<td class="price">100</td>`, `<td class="price">200</td>`, `<td class="price">300</td>`, `<td class="price">400</td>`}},
		{".price", []string{`<td class="price">100</td>`, `<td class="price">200</td>`, `<td class="price">300</td>`, `<td class="price">400</td>`}},
		{".product-table tr td:nth-child(2)", []string{`<td class="stock">10</td>`, `<td class="stock">20</td>`, `<td class="stock">30</td>`, `<td class="stock">40</td>`}},
		{".stock", []string{`<td class="stock">10</td>`, `<td class="stock">20</td>`, `<td class="stock">30</td>`, `<td class="stock">40</td>`}},
		{"#empty-div", []string{`<div id="empty-div"></div>`}},
		{"#multiple-children-div", []string{`<div id="multiple-children-div"><div id="first-child"></div><div id="second-child"></div></div>`}},
	}

	for _, test := range tests {
		testname := test.Query
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parents: []*Filter{
					{Results: []string{HTML_STRING}},
				},
				Var1: test.Query,
				Var2: &var2,
			}
			getFilterResultCSS(
				&filter,
			)
			if !DeepEqualStringSlice(filter.Results, test.Want) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}

func TestFilterCSSInnerHTML(t *testing.T) {
	var2 := "inner"
	var tests = []struct {
		Query string
		Want  []string
	}{
		{"title", []string{"title"}},
		{".product-table tr td:last-child", []string{`100`, `200`, `300`, `400`}},
		{".price", []string{`100`, `200`, `300`, `400`}},
		{".product-table tr td:nth-child(2)", []string{`10`, `20`, `30`, `40`}},
		{".stock", []string{`10`, `20`, `30`, `40`}},
		{"#empty-div]", []string{}},
		{"#multiple-children-div", []string{`<div id="first-child"></div><div id="second-child"></div>`}},
	}

	for _, test := range tests {
		testname := test.Query
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parents: []*Filter{
					{Results: []string{HTML_STRING}},
				},
				Var1: test.Query,
				Var2: &var2,
			}
			getFilterResultCSS(
				&filter,
			)
			if !DeepEqualStringSlice(filter.Results, test.Want) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}

func TestFilterCSSAttributes(t *testing.T) {
	var2 := "attr"
	var tests = []struct {
		Query string
		Want  []string
	}{
		{"title", []string{}},
		{"#table-caption", []string{`class="h3"`, `id="table-caption"`, `data="data"`}},
		{".product-table tr td:last-child", []string{`class="price"`, `class="price"`, `class="price"`, `class="price"`}},
		{".price", []string{`class="price"`, `class="price"`, `class="price"`, `class="price"`}},
		{".product-table tr td:nth-child(2)", []string{`class="stock"`, `class="stock"`, `class="stock"`, `class="stock"`}},
		{".stock", []string{`class="stock"`, `class="stock"`, `class="stock"`, `class="stock"`}},
		{".does-not-exists", []string{}},
	}

	for _, test := range tests {
		testname := test.Query
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parents: []*Filter{
					{Results: []string{HTML_STRING}},
				},
				Var1: test.Query,
				Var2: &var2,
			}
			getFilterResultCSS(
				&filter,
			)
			if !DeepEqualStringSlice(filter.Results, test.Want) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}

func TestFilterReplace(t *testing.T) {
	var tests = []struct {
		Input string
		Var1  string
		Var2  string
		Want  string
	}{
		// remove tests
		{"0123456789", "0", "", "123456789"},
		{"0123456789", "9", "", "012345678"},
		{"0123456789", "3456", "", "012789"},
		{"0123456789_0123456789", "3456", "", "012789_012789"},
		{"世界日本語", "世", "", "界日本語"},
		{"世界日本語", "語", "", "世界日本"},
		{"世界日_世界日_世界日", "界", "", "世日_世日_世日"},

		// replace tests
		{"0123456789", "0", "a", "a123456789"},
		{"0123456789", "9", "b", "012345678b"},
		{"0123456789", "3456", "abcd", "012abcd789"},
		{"0123456789_0123456789", "3456", "abcd", "012abcd789_012abcd789"},
		{"世界日本語", "世", "本語", "本語界日本語"},
		{"世界日本語", "語", "日", "世界日本日"},
		{"世界日_世界日_世界日", "界", "語", "世語日_世語日_世語日"},

		// regex remove tests
		{"0123456789", "0[0-9]{2}", "", "3456789"},
		{"0123456789", "[0-9]{2}9", "", "0123456"},
		{"0123456789", "[0-9]+", "", ""},
		{"Hello World!", "[eo]", "", "Hll Wrld!"},
		{"世界日本語", "[界本]", "", "世日語"},

		// regex replace tests
		{"<span>123.65</span>", "<span>([0-9.]+)</span>", "$1", "123.65"},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%s '%s' '%s'", test.Input, test.Var1, test.Var2)
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parents: []*Filter{
					{Results: []string{test.Input}},
				},
				Var1: test.Var1,
				Var2: &test.Var2,
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
		{"", "日本", []string{}},
	}

	for _, test := range tests {
		testname := test.Query
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parents: []*Filter{
					{Results: []string{test.Input}},
				},
				Var1: test.Query,
			}
			getFilterResultMatch(
				&filter,
			)
			if !DeepEqualStringSlice(filter.Results, test.Want) {
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
				Parents: []*Filter{
					{Results: []string{test.Input}},
				},
				Var1: test.Query,
			}
			getFilterResultSubstring(
				&filter,
			)
			if len(filter.Results) > 0 && filter.Results[0] != test.Want {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}

func TestFilterSubstringOutOfBounds(t *testing.T) {
	var tests = []struct {
		Input string
		Query string
	}{
		{"01234", ":-6"},
		{"01234", "-6:"},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%s %s", test.Input, test.Query)
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parents: []*Filter{
					{Results: []string{test.Input}},
				},
				Var1: test.Query,
			}
			getFilterResultSubstring(
				&filter,
			)
			if len(filter.Logs) == 0 {
				t.Errorf("No log message, expected one for OoB")
			}
		})
	}
}

func TestFilterContains(t *testing.T) {
	var tests = []struct {
		Input  []string
		Query  string
		Invert string
		Want   []string
	}{
		{[]string{"some text", "other text"}, "some", "false", []string{"some text"}},
		{[]string{"some text", "other text"}, "some", "true", []string{"other text"}},
		{[]string{"some text", "other text"}, "needle", "false", []string{}},
		{[]string{"some text", "other text"}, "needle", "true", []string{"some text", "other text"}},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%s %s %s", test.Input, test.Query, test.Invert)
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parents: []*Filter{
					{Results: test.Input},
				},
				Var1: test.Query,
				Var2: &test.Invert,
			}
			getFilterResultContains(
				&filter,
			)
			if !DeepEqualStringSlice(filter.Results, test.Want) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}

func TestFilterSum(t *testing.T) {
	var tests = []struct {
		Input []string
		Want  []string
	}{
		{[]string{"1"}, []string{"1.000000"}},
		{[]string{"1"}, []string{"1.000000"}},
		{[]string{"1", "1", "A"}, []string{"2.000000"}},
		{[]string{"1", "A", "B", "1"}, []string{"2.000000"}},
		{[]string{}, []string{"0.000000"}},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%s", test.Input)
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parents: []*Filter{
					{Results: test.Input},
				},
			}
			getFilterResultSum(
				&filter,
			)
			if !DeepEqualStringSlice(filter.Results, test.Want) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}

func TestFilterMin(t *testing.T) {
	var tests = []struct {
		Input []string
		Want  []string
	}{
		{[]string{"1"}, []string{"1.000000"}},
		{[]string{"10000"}, []string{"10000.000000"}},
		{[]string{"1", "2", "3", "4"}, []string{"1.000000"}},
		{[]string{"2000000", "100000", "A"}, []string{"100000.000000"}},
		{[]string{"1", "A", "B", "2"}, []string{"1.000000"}},
		{[]string{"1.1", "0.1", "10"}, []string{"0.100000"}},
		{[]string{}, []string{}},
		{[]string{"A"}, []string{}},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%s", test.Input)
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parents: []*Filter{
					{Results: test.Input},
				},
			}
			getFilterResultMin(
				&filter,
			)
			if !DeepEqualStringSlice(filter.Results, test.Want) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}

func TestFilterMax(t *testing.T) {
	var tests = []struct {
		Input []string
		Want  []string
	}{
		{[]string{"1"}, []string{"1.000000"}},
		{[]string{"10000"}, []string{"10000.000000"}},
		{[]string{"1", "2", "3", "4"}, []string{"4.000000"}},
		{[]string{"200000", "100000", "A"}, []string{"200000.000000"}},
		{[]string{"1", "A", "B", "2"}, []string{"2.000000"}},
		{[]string{"1.1", "0.1", "10"}, []string{"10.000000"}},
		{[]string{}, []string{}},
		{[]string{"A"}, []string{}},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%s", test.Input)
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parents: []*Filter{
					{Results: test.Input},
				},
			}
			getFilterResultMax(
				&filter,
			)
			if !DeepEqualStringSlice(filter.Results, test.Want) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}

func TestFilterAverage(t *testing.T) {
	var tests = []struct {
		Input []string
		Want  []string
	}{
		{[]string{"1"}, []string{"1.000000"}},
		{[]string{"10000"}, []string{"10000.000000"}},
		{[]string{"1", "2", "3", "4"}, []string{"2.500000"}},
		{[]string{"200000", "100000", "A"}, []string{"150000.000000"}},
		{[]string{"1", "A", "B", "2"}, []string{"1.500000"}},
		{[]string{"3.5", "5.5", "1.75", "1.25"}, []string{"3.000000"}},
		{[]string{}, []string{}},
		{[]string{"A"}, []string{}},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%s", test.Input)
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parents: []*Filter{
					{Results: test.Input},
				},
			}
			getFilterResultAverage(
				&filter,
			)
			if !DeepEqualStringSlice(filter.Results, test.Want) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}

func TestFilterCount(t *testing.T) {
	var tests = []struct {
		Input []string
		Want  []string
	}{
		{[]string{"1"}, []string{"1"}},
		{[]string{"10000"}, []string{"1"}},
		{[]string{"1", "2", "3", "4"}, []string{"4"}},
		{[]string{"200000", "100000", "A"}, []string{"3"}},
		{[]string{"1", "A", "B", "2"}, []string{"4"}},
		{[]string{"3.5", "5.5", "1.75", "1.25"}, []string{"4"}},
		{[]string{}, []string{"0"}},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%s", test.Input)
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parents: []*Filter{
					{Results: test.Input},
				},
			}
			getFilterResultCount(
				&filter,
			)
			if !DeepEqualStringSlice(filter.Results, test.Want) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}

func TestFilterRound(t *testing.T) {
	var tests = []struct {
		Input []string
		Query string
		Want  []string
	}{
		{[]string{"1.123456789"}, "0", []string{"1.000000"}},
		{[]string{"1.123456789"}, "1", []string{"1.100000"}},
		{[]string{"1.123456789"}, "2", []string{"1.120000"}},
		{[]string{"1.123456789"}, "3", []string{"1.123000"}},
		{[]string{"1.123456789"}, "4", []string{"1.123500"}},
		{[]string{"1.123456789"}, "5", []string{"1.123460"}},
		{[]string{"1.123456789"}, "6", []string{"1.123457"}},
		{[]string{"1.123456789"}, "7", []string{"1.123457"}},
		{[]string{"1.123456789"}, "8", []string{"1.123457"}},
		{[]string{"1.123456789"}, "9", []string{"1.123457"}},
		{[]string{"A"}, "9", []string{}},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%s %s", test.Input, test.Query)
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parents: []*Filter{
					{Results: test.Input},
				},
				Var2: &test.Query,
			}
			getFilterResultRound(
				&filter,
			)
			if !DeepEqualStringSlice(filter.Results, test.Want) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}

func getTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open("./test.db"))
	db.AutoMigrate(&Watch{}, &Filter{}, &FilterConnection{}, &FilterOutput{}, &ExpectFail{})
	return db
}

func TestConditionDiff(t *testing.T) {
	db := getTestDB()
	const timeLayout = "2006-01-02"
	time1, err := time.Parse(timeLayout, "2000-01-01")
	if err != nil {
		t.Error("Can't parse time")
	}
	time2, err := time.Parse(timeLayout, "2001-01-01")
	if err != nil {
		t.Error("Can't parse time")
	}
	testName := "Test"
	var tests = []struct {
		dbInput []FilterOutput
		WatchID WatchID
		Input   []string
		Want    []string
	}{
		{
			[]FilterOutput{
				{
					WatchID: 1,
					Name:    testName,
					Value:   "Last",
				},
			},
			1,
			[]string{"New"},
			[]string{"New"},
		},
		{
			[]FilterOutput{
				{
					WatchID: 2,
					Name:    testName,
					Value:   "Previous",
					Time:    time1,
				},
				{
					WatchID: 2,
					Name:    testName,
					Value:   "Last",
					Time:    time2,
				},
			},
			2,
			[]string{"New"},
			[]string{"New"},
		},
		{
			[]FilterOutput{
				{
					WatchID: 3,
					Name:    testName,
					Value:   "Same",
				},
			},
			3,
			[]string{"Same"},
			[]string{},
		},
		{
			[]FilterOutput{
				{
					WatchID: 4,
					Name:    testName,
					Value:   "Previous",
					Time:    time1,
				},
				{
					WatchID: 4,
					Name:    testName,
					Value:   "Same",
					Time:    time2,
				},
			},
			4,
			[]string{"Same"},
			[]string{},
		},
	}
	for _, test := range tests {
		testname := fmt.Sprintf("%s", test.Input)
		t.Run(testname, func(t *testing.T) {
			db.Create(&test.dbInput)
			filter := Filter{
				WatchID: test.WatchID,
				Name:    "Test",
				Var2:    &testName,
				Parents: []*Filter{
					{Results: test.Input},
				},
			}
			getFilterResultConditionDiff(
				&filter,
				db,
			)
			if !DeepEqualStringSlice(filter.Results, test.Want) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
	err = os.Remove("./test.db")
	if err != nil {
		log.Println("Could not remove test db:", err)
	}
}

func TestConditionLowerLast(t *testing.T) {
	db := getTestDB()
	const timeLayout = "2006-01-02"
	time1, err := time.Parse(timeLayout, "2000-01-01")
	if err != nil {
		t.Error("Can't parse time")
	}
	time2, err := time.Parse(timeLayout, "2001-01-01")
	if err != nil {
		t.Error("Can't parse time")
	}
	testName := "Test"
	var tests = []struct {
		dbInput []FilterOutput
		WatchID WatchID
		Input   []string
		Want    []string
	}{
		{
			[]FilterOutput{
				{
					WatchID: 1,
					Name:    testName,
					Value:   "2",
				},
			},
			1,
			[]string{"1"},
			[]string{"1"},
		},
		{
			[]FilterOutput{
				{
					WatchID: 2,
					Name:    testName,
					Value:   "A",
				},
			},
			2,
			[]string{"1"},
			[]string{"1"},
		},
		{
			[]FilterOutput{
				{
					WatchID: 3,
					Name:    testName,
					Value:   "3",
					Time:    time1,
				},
				{
					WatchID: 3,
					Name:    testName,
					Value:   "2",
					Time:    time2,
				},
			},
			3,
			[]string{"1"},
			[]string{"1"},
		},
		{
			[]FilterOutput{
				{
					WatchID: 4,
					Name:    testName,
					Value:   "1",
				},
			},
			4,
			[]string{"2"},
			[]string{},
		},
		{
			[]FilterOutput{
				{
					WatchID: 5,
					Name:    testName,
					Value:   "3",
					Time:    time1,
				},
				{
					WatchID: 5,
					Name:    testName,
					Value:   "1",
					Time:    time2,
				},
			},
			5,
			[]string{"2"},
			[]string{},
		},
	}
	for _, test := range tests {
		testname := fmt.Sprintf("%s", test.Input)
		t.Run(testname, func(t *testing.T) {
			db.Create(&test.dbInput)
			filter := Filter{
				WatchID: test.WatchID,
				Name:    testName,
				Var2:    &testName,
				Parents: []*Filter{
					{Results: test.Input},
				},
			}
			getFilterResultConditionLowerLast(
				&filter,
				db,
			)
			if !DeepEqualStringSlice(filter.Results, test.Want) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
	err = os.Remove("./test.db")
	if err != nil {
		log.Println("Could not remove test db:", err)
	}
}

func TestConditionLowest(t *testing.T) {
	db := getTestDB()
	testName := "Test"
	var tests = []struct {
		dbInput []FilterOutput
		WatchID WatchID
		Input   []string
		Want    []string
	}{
		{
			[]FilterOutput{
				{
					WatchID: 1,
					Name:    testName,
					Value:   "5",
				},
			},
			1,
			[]string{"4"},
			[]string{"4"},
		},
		{
			[]FilterOutput{
				{
					WatchID: 1,
					Name:    testName,
					Value:   "A",
				},
			},
			1,
			[]string{"4"},
			[]string{"4"},
		},
		{
			[]FilterOutput{
				{
					WatchID: 2,
					Name:    testName,
					Value:   "3",
				},
				{
					WatchID: 2,
					Name:    testName,
					Value:   "2",
				},
			},
			2,
			[]string{"1"},
			[]string{"1"},
		},
		{
			[]FilterOutput{
				{
					WatchID: 3,
					Name:    testName,
					Value:   "1",
				},
			},
			3,
			[]string{"2"},
			[]string{},
		},
		{
			[]FilterOutput{
				{
					WatchID: 4,
					Name:    testName,
					Value:   "3",
				},
				{
					WatchID: 4,
					Name:    testName,
					Value:   "1",
				},
			},
			4,
			[]string{"2"},
			[]string{},
		},
	}
	for _, test := range tests {
		testname := fmt.Sprintf("%s", test.Input)
		t.Run(testname, func(t *testing.T) {
			db.Create(&test.dbInput)
			filter := Filter{
				WatchID: test.WatchID,
				Name:    testName,
				Var2:    &testName,
				Parents: []*Filter{
					{Results: test.Input},
				},
			}
			getFilterResultConditionLowest(
				&filter,
				db,
			)
			if !DeepEqualStringSlice(filter.Results, test.Want) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
	err := os.Remove("./test.db")
	if err != nil {
		log.Println("Could not remove test db:", err)
	}
}

func TestFilterLowerThan(t *testing.T) {
	var tests = []struct {
		Input     []string
		Threshold string
		Want      []string
	}{
		{[]string{"1"}, "2", []string{"1"}},
		{[]string{"2"}, "1", []string{}},
		{[]string{"1"}, "1", []string{}},
		{[]string{"2", "3", "4"}, "3", []string{"2"}},
		{[]string{"A", "3", "4"}, "2", []string{}},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%s", test.Input)
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parents: []*Filter{
					{
						Results: test.Input,
					},
				},
				Var2: &test.Threshold,
			}
			getFilterResultConditionLowerThan(
				&filter,
			)
			if !DeepEqualStringSlice(filter.Results, test.Want) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}

func TestConditionHigherLast(t *testing.T) {
	db := getTestDB()
	const timeLayout = "2006-01-02"
	time1, err := time.Parse(timeLayout, "2000-01-01")
	if err != nil {
		t.Error("Can't parse time")
	}
	time2, err := time.Parse(timeLayout, "2001-01-01")
	if err != nil {
		t.Error("Can't parse time")
	}
	testName := "Test"
	var tests = []struct {
		dbInput []FilterOutput
		WatchID WatchID
		Input   []string
		Want    []string
	}{
		{
			[]FilterOutput{
				{
					WatchID: 1,
					Name:    testName,
					Value:   "1",
					Time:    time1,
				},
			},
			1,
			[]string{"2"},
			[]string{"2"},
		},
		{
			[]FilterOutput{
				{
					WatchID: 2,
					Name:    testName,
					Value:   "A",
				},
			},
			2,
			[]string{"1"},
			[]string{"1"},
		},
		{
			[]FilterOutput{
				{
					WatchID: 3,
					Name:    testName,
					Value:   "3",
					Time:    time1,
				},
				{
					WatchID: 3,
					Name:    testName,
					Value:   "2",
					Time:    time2,
				},
			},
			3,
			[]string{"3"},
			[]string{"3"},
		},
		{
			[]FilterOutput{
				{
					WatchID: 4,
					Name:    testName,
					Value:   "2",
				},
			},
			4,
			[]string{"1"},
			[]string{},
		},
		{
			[]FilterOutput{
				{
					WatchID: 5,
					Name:    testName,
					Value:   "1",
					Time:    time1,
				},
				{
					WatchID: 5,
					Name:    testName,
					Value:   "3",
					Time:    time2,
				},
			},
			5,
			[]string{"2"},
			[]string{},
		},
	}
	for _, test := range tests {
		testname := fmt.Sprintf("%s", test.Input)
		t.Run(testname, func(t *testing.T) {
			db.Create(&test.dbInput)
			filter := Filter{
				WatchID: test.WatchID,
				Name:    testName,
				Var2:    &testName,
				Parents: []*Filter{
					{Results: test.Input},
				},
			}
			getFilterResultConditionHigherLast(
				&filter,
				db,
			)
			if !DeepEqualStringSlice(filter.Results, test.Want) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
	err = os.Remove("./test.db")
	if err != nil {
		log.Println("Could not remove test db:", err)
	}
}

func TestConditionHighest(t *testing.T) {
	db := getTestDB()
	testName := "Test"
	var tests = []struct {
		dbInput []FilterOutput
		WatchID WatchID
		Input   []string
		Want    []string
	}{
		{
			[]FilterOutput{
				{
					WatchID: 1,
					Name:    testName,
					Value:   "1",
				},
			},
			1,
			[]string{"2"},
			[]string{"2"},
		},
		{
			[]FilterOutput{
				{
					WatchID: 2,
					Name:    testName,
					Value:   "A",
				},
			},
			2,
			[]string{"1"},
			[]string{"1"},
		},
		{
			[]FilterOutput{
				{
					WatchID: 3,
					Name:    testName,
					Value:   "1",
				},
				{
					WatchID: 3,
					Name:    testName,
					Value:   "2",
				},
			},
			3,
			[]string{"3"},
			[]string{"3"},
		},
		{
			[]FilterOutput{
				{
					WatchID: 4,
					Name:    testName,
					Value:   "2",
				},
			},
			4,
			[]string{"1"},
			[]string{},
		},
		{
			[]FilterOutput{
				{
					WatchID: 5,
					Name:    testName,
					Value:   "1",
				},
				{
					WatchID: 5,
					Name:    testName,
					Value:   "3",
				},
			},
			5,
			[]string{"2"},
			[]string{},
		},
	}
	for _, test := range tests {
		testname := fmt.Sprintf("%s", test.Input)
		t.Run(testname, func(t *testing.T) {
			db.Create(&test.dbInput)
			filter := Filter{
				WatchID: test.WatchID,
				Name:    testName,
				Var2:    &testName,
				Parents: []*Filter{
					{Results: test.Input},
				},
			}
			getFilterResultConditionHighest(
				&filter,
				db,
			)
			if !DeepEqualStringSlice(filter.Results, test.Want) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
	err := os.Remove("./test.db")
	if err != nil {
		log.Println("Could not remove test db:", err)
	}
}

func TestFilterHigherThan(t *testing.T) {
	var tests = []struct {
		Input     []string
		Threshold string
		Want      []string
	}{
		{[]string{"2"}, "1", []string{"2"}},
		{[]string{"1"}, "2", []string{}},
		{[]string{"1"}, "1", []string{}},
		{[]string{"1", "2", "3"}, "1", []string{"2", "3"}},
		{[]string{"1", "2", "3", "A"}, "4", []string{}},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%s", test.Input)
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parents: []*Filter{
					{
						Results: test.Input,
					},
				},
				Var2: &test.Threshold,
			}
			getFilterResultConditionHigherThan(
				&filter,
			)
			if !DeepEqualStringSlice(filter.Results, test.Want) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}

func TestFilterUnique(t *testing.T) {
	var tests = []struct {
		Input []string
		Want  []string
	}{
		{[]string{"1"}, []string{"1"}},
		{[]string{"1", "2"}, []string{"1", "2"}},
		{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "1"}, []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}},
		{[]string{"1", "1"}, []string{"1"}},
		{[]string{}, []string{}},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%s", test.Input)
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parents: []*Filter{
					{
						Results: test.Input,
					},
				},
			}
			getFilterResultUnique(
				&filter,
			)
			wantMap := make(map[string]bool)
			for _, want := range test.Want {
				wantMap[want] = true
			}
			resultMap := make(map[string]bool)
			for _, result := range filter.Results {
				resultMap[result] = true
			}
			if !reflect.DeepEqual(wantMap, resultMap) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}

func TestFilterLua(t *testing.T) {
	passAll := `
for i,input in pairs(inputs) do
   	table.insert(outputs, input)
end`
	lessThanFour := `
for i,input in pairs(inputs) do
	if tonumber(input) < 4 then
		table.insert(outputs, input)
	end
end`
	concat := `table.insert(outputs, table.concat(inputs, ","))`
	var tests = []struct {
		Name  string
		Input []string
		Lua   string
		Want  []string
	}{
		{"Pass all", []string{"1", "2", "3", "4", "5"}, passAll, []string{"1", "2", "3", "4", "5"}},
		{"Less than four", []string{"1", "2", "3", "4", "5"}, lessThanFour, []string{"1", "2", "3"}},
		{"Concat", []string{"1", "2", "3", "4", "5"}, concat, []string{"1,2,3,4,5"}},
	}

	for _, test := range tests {
		testname := test.Name
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Var1: test.Lua,
				Parents: []*Filter{
					{
						Results: test.Input,
					},
				},
			}
			getFilterResultLua(
				&filter,
			)
			if len(filter.Logs) > 0 {
				t.Errorf("Lua error: %s", filter.Logs)
			}
			if !DeepEqualStringSlice(filter.Results, test.Want) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}

func TestFilterLuaLibs(t *testing.T) {
	regex := `
local regexp = require("regexp")
local inspect = require("inspect")

-- regexp.match(regexp, data)
local result, err = regexp.match("hello", "hello world")
table.insert(logs, err)
if err then error(err) end
if not(result==true) then error("regexp.match()") end
table.insert(outputs, result)`
	var tests = []struct {
		Name  string
		Input []string
		Lua   string
		Want  []string
	}{
		{"Regex", []string{}, regex, []string{"true"}},
	}

	for _, test := range tests {
		testname := test.Name
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Var1: test.Lua,
				Parents: []*Filter{
					{
						Results: test.Input,
					},
				},
			}
			getFilterResultLua(
				&filter,
			)
			if len(filter.Logs) > 0 {
				t.Errorf("Lua error: %s", filter.Logs)
			}
			if !DeepEqualStringSlice(filter.Results, test.Want) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}

func TestFilterLuaLogs(t *testing.T) {
	script := `table.insert(logs, "test")`

	filter := Filter{
		Var1: script,
	}

	getFilterResultLua(&filter)

	if len(filter.Logs) == 0 {
		t.Error("Nothing in logs, expected 'test'")
	}

	if filter.Logs[0] != "test" {
		t.Errorf("Unexpected log message: '%s'", filter.Logs[0])
	}
}

func TestEchoFilter(t *testing.T) {
	helloWorld := "Hello World!"
	filters := []Filter{
		{
			ID:   0,
			Name: "Echo",
			Type: "echo",
			Var1: helloWorld,
		},
	}
	filter1 := &filters[0]
	connections := []FilterConnection{}
	buildFilterTree(filters, connections)
	ProcessFilters(filters, nil, nil, false, nil)

	if !DeepEqualStringSlice(filter1.Results, []string{helloWorld}) {
		t.Errorf("%s did not match %s", helloWorld, filter1.Results)
	}
}

func TestSimpleWatch(t *testing.T) {
	filters := []Filter{
		{
			ID:   0,
			Name: "Echo",
			Type: "echo",
			Var1: HTML_STRING,
		},
		{
			ID:   1,
			Name: "XPath",
			Type: "xpath",
			Var1: "//td[@class='price']",
		},
		{
			ID:   2,
			Name: "Replace",
			Type: "replace",
			Var1: "[^0-9]",
		},
		{
			ID:   3,
			Name: "Min",
			Type: "math",
			Var1: "min",
		},
		{
			ID:   4,
			Name: "Max",
			Type: "math",
			Var1: "max",
		},
	}

	minFilter := &filters[3]
	maxFilter := &filters[4]

	connections := []FilterConnection{
		{
			OutputID: 0,
			InputID:  1,
		},
		{
			OutputID: 1,
			InputID:  2,
		},
		{
			OutputID: 2,
			InputID:  3,
		},
		{
			OutputID: 2,
			InputID:  4,
		},
	}

	buildFilterTree(filters, connections)
	ProcessFilters(filters, nil, nil, false, nil)

	if !reflect.DeepEqual(minFilter.Results, []string{"100.000000"}) {
		t.Errorf("%s did not match '100'", minFilter.Results)
	}
	if !reflect.DeepEqual(maxFilter.Results, []string{"400.000000"}) {
		t.Errorf("%s did not match '400'", maxFilter.Results)
	}
}
func TestSimpleIDOrderWatch(t *testing.T) {
	filters := []Filter{
		{
			ID:   7,
			Name: "Echo",
			Type: "echo",
			Var1: HTML_STRING,
		},
		{
			ID:   5,
			Name: "XPath",
			Type: "xpath",
			Var1: "//td[@class='price']",
		},
		{
			ID:   9,
			Name: "Replace",
			Type: "replace",
			Var1: "[^0-9]",
		},
		{
			ID:   15,
			Name: "Min",
			Type: "math",
			Var1: "min",
		},
		{
			ID:   1,
			Name: "Max",
			Type: "math",
			Var1: "max",
		},
	}

	minFilter := &filters[3]
	maxFilter := &filters[4]

	connections := []FilterConnection{
		{
			OutputID: 7,
			InputID:  5,
		},
		{
			OutputID: 5,
			InputID:  9,
		},
		{
			OutputID: 9,
			InputID:  15,
		},
		{
			OutputID: 9,
			InputID:  1,
		},
	}

	buildFilterTree(filters, connections)
	ProcessFilters(filters, nil, nil, false, nil)

	if !reflect.DeepEqual(minFilter.Results, []string{"100.000000"}) {
		t.Errorf("%s did not match '100'", minFilter.Results)
	}
	if !reflect.DeepEqual(maxFilter.Results, []string{"400.000000"}) {
		t.Errorf("%s did not match '400'", maxFilter.Results)
	}
}
func TestSimpleDebugWatch(t *testing.T) {
	filters := []Filter{
		{
			ID:   0,
			Name: "Echo",
			Type: "echo",
			Var1: HTML_STRING,
		},
		{
			ID:   1,
			Name: "XPath",
			Type: "xpath",
			Var1: "//td[@class='price']",
		},
		{
			ID:   2,
			Name: "Replace",
			Type: "replace",
			Var1: "[^0-9]",
		},
		{
			ID:   3,
			Name: "Min",
			Type: "math",
			Var1: "min",
		},
		{
			ID:   4,
			Name: "Max",
			Type: "math",
			Var1: "max",
		},
	}

	minFilter := &filters[3]
	maxFilter := &filters[4]

	connections := []FilterConnection{
		{
			OutputID: 0,
			InputID:  1,
		},
		{
			OutputID: 1,
			InputID:  2,
		},
		{
			OutputID: 2,
			InputID:  3,
		},
		{
			OutputID: 2,
			InputID:  4,
		},
	}

	buildFilterTree(filters, connections)
	ProcessFilters(filters, nil, nil, true, nil)

	if !reflect.DeepEqual(minFilter.Results, []string{"100.000000"}) {
		t.Errorf("%s did not match '100'", minFilter.Results)
	}
	if !reflect.DeepEqual(maxFilter.Results, []string{"400.000000"}) {
		t.Errorf("%s did not match '400'", maxFilter.Results)
	}
}

func TestSimpleTriggeredWatch(t *testing.T) {
	db := getTestDB()
	watch := Watch{
		Name: "Test",
	}
	db.Create(&watch)
	filters := []Filter{
		{
			WatchID: watch.ID,
			Name:    "Schedule",
			Type:    "cron",
		},
		{
			WatchID: watch.ID,
			Name:    "Echo",
			Type:    "echo",
			Var1:    HTML_STRING,
		},
		{
			WatchID: watch.ID,
			Name:    "XPath",
			Type:    "xpath",
			Var1:    "//td[@class='price']",
		},
		{
			WatchID: watch.ID,
			Name:    "Replace",
			Type:    "replace",
			Var1:    "[^0-9]",
		},
		{
			WatchID: watch.ID,
			Name:    "Min",
			Type:    "math",
			Var1:    "min",
		},
		{
			WatchID: watch.ID,
			Name:    "Minimum",
			Type:    "store",
		},
		{
			WatchID: watch.ID,
			Name:    "Max",
			Type:    "math",
			Var1:    "max",
		},
		{
			WatchID: watch.ID,
			Name:    "Maximum",
			Type:    "store",
		},
	}
	db.Create(&filters)
	scheduleFilter := &filters[0]
	echoFilter := &filters[1]
	xpathFilter := &filters[2]
	replaceFilter := &filters[3]
	minFilter := &filters[4]
	storeMinFilter := &filters[5]
	maxFilter := &filters[6]
	storeMaxFilter := &filters[7]

	connections := []FilterConnection{
		{
			WatchID:  watch.ID,
			OutputID: scheduleFilter.ID,
			InputID:  echoFilter.ID,
		},
		{
			WatchID:  watch.ID,
			OutputID: echoFilter.ID,
			InputID:  xpathFilter.ID,
		},
		{
			WatchID:  watch.ID,
			OutputID: xpathFilter.ID,
			InputID:  replaceFilter.ID,
		},
		{
			WatchID:  watch.ID,
			OutputID: replaceFilter.ID,
			InputID:  minFilter.ID,
		},
		{
			WatchID:  watch.ID,
			OutputID: minFilter.ID,
			InputID:  storeMinFilter.ID,
		},
		{
			WatchID:  watch.ID,
			OutputID: replaceFilter.ID,
			InputID:  maxFilter.ID,
		},
		{
			WatchID:  watch.ID,
			OutputID: maxFilter.ID,
			InputID:  storeMaxFilter.ID,
		},
	}
	db.Create(&connections)

	TriggerSchedule(watch.ID, &Web{db: db}, &scheduleFilter.ID)

	var filterOutputs []FilterOutput
	db.Model(&FilterOutput{}).Find(&filterOutputs, fmt.Sprintf("watch_id = %d", watch.ID))

	for _, filterOutput := range filterOutputs {
		if filterOutput.Name == "Maximum" {
			if filterOutput.Value != "400.000000" {
				t.Errorf("Minimum filter value 400.000000 != %s", filterOutput.Value)
			}
		} else if filterOutput.Name == "Minimum" {
			if filterOutput.Value != "100.000000" {
				t.Errorf("Minimum filter value 100.000000 != %s", filterOutput.Value)
			}
		} else {
			t.Errorf("Unknown filter name: %s", filterOutput.Name)
		}
	}
	err := os.Remove("./test.db")
	if err != nil {
		log.Println("Could not remove test db:", err)
	}
}

func TestDontAllowMultipleCronOnSingleFilter(t *testing.T) {
	filters := []Filter{
		{
			ID:   0,
			Name: "Cron1",
			Type: "cron",
			Var1: "@every 1s",
		},
		{
			ID:   1,
			Name: "Cron2",
			Type: "cron",
			Var1: "@every 1s",
		},
		{
			ID:   2,
			Name: "Filter",
			Type: "echo",
			Var1: "-",
		},
	}

	filter := &filters[2]

	connections := []FilterConnection{
		{
			OutputID: 0,
			InputID:  2,
		},
		{
			OutputID: 0,
			InputID:  2,
		},
	}

	buildFilterTree(filters, connections)
	ProcessFilters(filters, nil, nil, false, nil)

	if len(filter.Logs) == 0 {
		t.Errorf("Expected error message in filter log, found empty log: %s", filter.Logs)
	}
}

func TestWatchWithExpectNotTriggering(t *testing.T) {
	db := getTestDB()
	filters := []Filter{
		{
			ID:   0,
			Name: "Echo",
			Type: "echo",
			Var1: HTML_STRING,
		},
		{
			ID:   1,
			Name: "XPath",
			Type: "xpath",
			Var1: "//td[@class='price']",
		},
		{
			ID:   2,
			Name: "Expect",
			Type: "expect",
			Var1: "1",
		},
	}

	expectFilter := &filters[2]

	connections := []FilterConnection{
		{
			OutputID: 0,
			InputID:  1,
		},
		{
			OutputID: 1,
			InputID:  2,
		},
	}

	buildFilterTree(filters, connections)
	ProcessFilters(filters, &Web{db: db}, nil, false, nil)

	if len(expectFilter.Results) != 0 {
		t.Error("Expect has results, should be empty:", expectFilter.Results)
	}

	err := os.Remove("./test.db")
	if err != nil {
		log.Println("Could not remove test db:", err)
	}
}

func TestWatchWithExpectTriggering(t *testing.T) {
	db := getTestDB()
	filters := []Filter{
		{
			ID:   0,
			Name: "Echo",
			Type: "echo",
			Var1: HTML_STRING,
		},
		{
			ID:   1,
			Name: "XPath",
			Type: "xpath",
			Var1: "//div[@class='price']",
		},
		{
			ID:   2,
			Name: "Expect",
			Type: "expect",
			Var1: "1",
		},
	}

	expectFilter := &filters[2]

	connections := []FilterConnection{
		{
			OutputID: 0,
			InputID:  1,
		},
		{
			OutputID: 1,
			InputID:  2,
		},
	}

	buildFilterTree(filters, connections)
	ProcessFilters(filters, &Web{db: db}, nil, false, nil)

	if len(expectFilter.Results) != 1 {
		t.Error("Expect has no results, should have 'expected'")
	}

	err := os.Remove("./test.db")
	if err != nil {
		log.Println("Could not remove test db:", err)
	}
}
func TestWatchWithExpect3Triggering(t *testing.T) {
	db := getTestDB()
	filters := []Filter{
		{
			ID:   0,
			Name: "Echo",
			Type: "echo",
			Var1: HTML_STRING,
		},
		{
			ID:   1,
			Name: "XPath",
			Type: "xpath",
			Var1: "//div[@class='price']",
		},
		{
			ID:   2,
			Name: "Expect",
			Type: "expect",
			Var1: "3",
		},
	}

	expectFilter := &filters[2]

	connections := []FilterConnection{
		{
			OutputID: 0,
			InputID:  1,
		},
		{
			OutputID: 1,
			InputID:  2,
		},
	}

	buildFilterTree(filters, connections)

	ProcessFilters(filters, &Web{db: db}, nil, false, nil)

	if len(expectFilter.Results) != 0 {
		t.Error("Expect has results, should be empty:", expectFilter.Results)
	}

	ProcessFilters(filters, &Web{db: db}, nil, false, nil)

	if len(expectFilter.Results) != 0 {
		t.Error("Expect has results, should be empty:", expectFilter.Results)
	}

	ProcessFilters(filters, &Web{db: db}, nil, false, nil)

	if len(expectFilter.Results) != 1 {
		t.Error("Expect has no results, should have 'expected'")
	}

	err := os.Remove("./test.db")
	if err != nil {
		log.Println("Could not remove test db:", err)
	}
}

func TestWatchWithExpectNotTriggeringDB(t *testing.T) {
	db := getTestDB()
	watch := Watch{
		Name: "Test",
	}
	db.Create(&watch)
	filters := []Filter{
		{
			WatchID: watch.ID,
			Name:    "Schedule",
			Type:    "cron",
		},
		{
			WatchID: watch.ID,
			Name:    "Echo",
			Type:    "echo",
			Var1:    HTML_STRING,
		},
		{
			WatchID: watch.ID,
			Name:    "XPath",
			Type:    "xpath",
			Var1:    "//td[@class='price']",
		},
		{
			WatchID: watch.ID,
			Name:    "Expect",
			Type:    "expect",
			Var1:    "1",
		},
	}
	db.Create(&filters)
	scheduleFilter := &filters[0]
	echoFilter := &filters[1]
	xpathFilter := &filters[2]
	expectFilter := &filters[3]

	connections := []FilterConnection{
		{
			WatchID:  watch.ID,
			OutputID: scheduleFilter.ID,
			InputID:  echoFilter.ID,
		},
		{
			WatchID:  watch.ID,
			OutputID: echoFilter.ID,
			InputID:  xpathFilter.ID,
		},
		{
			WatchID:  watch.ID,
			OutputID: xpathFilter.ID,
			InputID:  expectFilter.ID,
		},
	}
	db.Create(&connections)

	TriggerSchedule(watch.ID, &Web{db: db}, &scheduleFilter.ID)

	var expectFails []ExpectFail
	db.Model(&ExpectFail{}).Find(&expectFails, "watch_id = ?", watch.ID)
	if len(expectFails) > 0 {
		t.Errorf("Found ExpectFail values expected none!")
	}
	err := os.Remove("./test.db")
	if err != nil {
		log.Println("Could not remove test db:", err)
	}
}
func TestWatchWithExpectTriggeringDB(t *testing.T) {
	db := getTestDB()
	watch := Watch{
		Name: "Test",
	}
	db.Create(&watch)
	filters := []Filter{
		{
			WatchID: watch.ID,
			Name:    "Schedule",
			Type:    "cron",
		},
		{
			WatchID: watch.ID,
			Name:    "Echo",
			Type:    "echo",
			Var1:    HTML_STRING,
		},
		{
			WatchID: watch.ID,
			Name:    "XPath",
			Type:    "xpath",
			Var1:    "//div[@class='price']",
		},
		{
			WatchID: watch.ID,
			Name:    "Expect",
			Type:    "expect",
			Var1:    "1",
		},
	}
	db.Create(&filters)
	scheduleFilter := &filters[0]
	echoFilter := &filters[1]
	xpathFilter := &filters[2]
	expectFilter := &filters[3]

	connections := []FilterConnection{
		{
			WatchID:  watch.ID,
			OutputID: scheduleFilter.ID,
			InputID:  echoFilter.ID,
		},
		{
			WatchID:  watch.ID,
			OutputID: echoFilter.ID,
			InputID:  xpathFilter.ID,
		},
		{
			WatchID:  watch.ID,
			OutputID: xpathFilter.ID,
			InputID:  expectFilter.ID,
		},
	}
	db.Create(&connections)

	TriggerSchedule(watch.ID, &Web{db: db}, &scheduleFilter.ID)

	var expectFails []ExpectFail
	db.Model(&ExpectFail{}).Find(&expectFails, "watch_id = ?", watch.ID)
	if len(expectFails) != 1 {
		t.Errorf("Found no ExpectFail values expected 1!")
	}
	err := os.Remove("./test.db")
	if err != nil {
		log.Println("Could not remove test db:", err)
	}
}
func TestWatchWithExpect3TriggeringDB(t *testing.T) {
	db := getTestDB()
	watch := Watch{
		Name: "Test",
	}
	db.Create(&watch)
	filters := []Filter{
		{
			WatchID: watch.ID,
			Name:    "Schedule",
			Type:    "cron",
		},
		{
			WatchID: watch.ID,
			Name:    "Echo",
			Type:    "echo",
			Var1:    HTML_STRING,
		},
		{
			WatchID: watch.ID,
			Name:    "XPath",
			Type:    "xpath",
			Var1:    "//div[@class='price']",
		},
		{
			WatchID: watch.ID,
			Name:    "Expect",
			Type:    "expect",
			Var1:    "3",
		},
	}
	db.Create(&filters)
	scheduleFilter := &filters[0]
	echoFilter := &filters[1]
	xpathFilter := &filters[2]
	expectFilter := &filters[3]

	connections := []FilterConnection{
		{
			WatchID:  watch.ID,
			OutputID: scheduleFilter.ID,
			InputID:  echoFilter.ID,
		},
		{
			WatchID:  watch.ID,
			OutputID: echoFilter.ID,
			InputID:  xpathFilter.ID,
		},
		{
			WatchID:  watch.ID,
			OutputID: xpathFilter.ID,
			InputID:  expectFilter.ID,
		},
	}
	db.Create(&connections)

	var expectFails []ExpectFail
	TriggerSchedule(watch.ID, &Web{db: db}, &scheduleFilter.ID)

	db.Model(&ExpectFail{}).Find(&expectFails, "watch_id = ?", watch.ID)
	if len(expectFails) != 1 {
		t.Errorf("Found %d ExpectFail values, expected 1!", len(expectFails))
		log.Println(expectFails)
	}

	TriggerSchedule(watch.ID, &Web{db: db}, &scheduleFilter.ID)

	db.Model(&ExpectFail{}).Find(&expectFails, "watch_id = ?", watch.ID)
	if len(expectFails) != 2 {
		t.Errorf("Found %d ExpectFail values, expected 2!", len(expectFails))
		log.Println(expectFails)
	}

	TriggerSchedule(watch.ID, &Web{db: db}, &scheduleFilter.ID)

	db.Model(&ExpectFail{}).Find(&expectFails, "watch_id = ?", watch.ID)
	if len(expectFails) != 3 {
		t.Errorf("Found %d ExpectFail values, expected 3! (1)", len(expectFails))
		log.Println(expectFails)
	}
	TriggerSchedule(watch.ID, &Web{db: db}, &scheduleFilter.ID)

	db.Model(&ExpectFail{}).Find(&expectFails, "watch_id = ?", watch.ID)
	if len(expectFails) != 3 {
		t.Errorf("Found %d ExpectFail values, expected 3! (2)", len(expectFails))
		log.Println(expectFails)
	}

	err := os.Remove("./test.db")
	if err != nil {
		log.Println("Could not remove test db:", err)
	}
}
