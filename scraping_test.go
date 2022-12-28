package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
				Parents: []*Filter{
					{Results: []string{HTML_STRING}},
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
				Parents: []*Filter{
					{Results: []string{JSON_STRING}},
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
				Parents: []*Filter{
					{Results: []string{HTML_STRING}},
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

		// regex tests
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
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%s", test.Query)
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
			if (filter.Results != nil && test.Want != nil) && !reflect.DeepEqual(test.Want, filter.Results) {
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
			if (filter.Results != nil && test.Want != nil) && !reflect.DeepEqual(test.Want, filter.Results) {
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
			if (filter.Results != nil && test.Want != nil) && !reflect.DeepEqual(test.Want, filter.Results) {
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
			if (filter.Results != nil && test.Want != nil) && !reflect.DeepEqual(test.Want, filter.Results) {
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
			if (filter.Results != nil && test.Want != nil) && !reflect.DeepEqual(test.Want, filter.Results) {
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
			if (filter.Results != nil && test.Want != nil) && !reflect.DeepEqual(test.Want, filter.Results) {
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
			if (filter.Results != nil && test.Want != nil) && !reflect.DeepEqual(test.Want, filter.Results) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}

func getTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open("./test.db"))
	db.AutoMigrate(&Watch{}, &Filter{}, &FilterConnection{}, &FilterOutput{})
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
		WatchID uint
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
			if (filter.Results != nil && test.Want != nil) && !reflect.DeepEqual(test.Want, filter.Results) {
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
		WatchID uint
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
					WatchID: 1,
					Name:    testName,
					Value:   "A",
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
					Value:   "3",
					Time:    time1,
				},
				{
					WatchID: 2,
					Name:    testName,
					Value:   "2",
					Time:    time2,
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
					Time:    time1,
				},
				{
					WatchID: 4,
					Name:    testName,
					Value:   "1",
					Time:    time2,
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
			getFilterResultConditionLowerLast(
				&filter,
				db,
			)
			if (filter.Results != nil && test.Want != nil) && !reflect.DeepEqual(test.Want, filter.Results) {
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
		WatchID uint
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
			if (filter.Results != nil && test.Want != nil) && !reflect.DeepEqual(test.Want, filter.Results) {
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
		{[]string{"2", "3", "4"}, "1", []string{"1"}},
		{[]string{"A", "3", "4"}, "2", []string{"2"}},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%s", test.Input)
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parents: []*Filter{
					{
						Results: test.Input,
						Var2:    &test.Threshold,
					},
				},
			}
			getFilterResultConditionLowerThan(
				&filter,
			)
			if (filter.Results != nil && test.Want != nil) && !reflect.DeepEqual(test.Want, filter.Results) {
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
		WatchID uint
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
					WatchID: 1,
					Name:    testName,
					Value:   "A",
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
					Value:   "3",
					Time:    time1,
				},
				{
					WatchID: 2,
					Name:    testName,
					Value:   "2",
					Time:    time2,
				},
			},
			2,
			[]string{"3"},
			[]string{"3"},
		},
		{
			[]FilterOutput{
				{
					WatchID: 3,
					Name:    testName,
					Value:   "2",
				},
			},
			3,
			[]string{"1"},
			[]string{},
		},
		{
			[]FilterOutput{
				{
					WatchID: 4,
					Name:    testName,
					Value:   "1",
					Time:    time1,
				},
				{
					WatchID: 4,
					Name:    testName,
					Value:   "3",
					Time:    time2,
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
			getFilterResultConditionHigherLast(
				&filter,
				db,
			)
			if (filter.Results != nil && test.Want != nil) && !reflect.DeepEqual(test.Want, filter.Results) {
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
		WatchID uint
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
					WatchID: 1,
					Name:    testName,
					Value:   "A",
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
					Value:   "1",
				},
				{
					WatchID: 2,
					Name:    testName,
					Value:   "2",
				},
			},
			2,
			[]string{"3"},
			[]string{"3"},
		},
		{
			[]FilterOutput{
				{
					WatchID: 3,
					Name:    testName,
					Value:   "2",
				},
			},
			3,
			[]string{"1"},
			[]string{},
		},
		{
			[]FilterOutput{
				{
					WatchID: 4,
					Name:    testName,
					Value:   "1",
				},
				{
					WatchID: 4,
					Name:    testName,
					Value:   "3",
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
			getFilterResultConditionHighest(
				&filter,
				db,
			)
			if (filter.Results != nil && test.Want != nil) && !reflect.DeepEqual(test.Want, filter.Results) {
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
		{[]string{"1", "2", "3"}, "4", []string{"4"}},
		{[]string{"1", "2", "3", "A"}, "4", []string{"4"}},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%s", test.Input)
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parents: []*Filter{
					{
						Results: test.Input,
						Var2:    &test.Threshold,
					},
				},
			}
			getFilterResultConditionLowerThan(
				&filter,
			)
			if (filter.Results != nil && test.Want != nil) && !reflect.DeepEqual(test.Want, filter.Results) {
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
			var wantMap map[string]bool
			wantMap = make(map[string]bool)
			for _, want := range test.Want {
				wantMap[want] = true
			}
			var resultMap map[string]bool
			resultMap = make(map[string]bool)
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
		testname := fmt.Sprintf("%s", test.Name)
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
			log.Println(test.Want)
			log.Println(filter.Results)
			if (filter.Results != nil && test.Want != nil) && !reflect.DeepEqual(test.Want, filter.Results) {
				t.Errorf("Got %s, want %s", filter.Results, test.Want)
			}
		})
	}
}

func TestFilterLuaLibs(t *testing.T) {
	regex := `
	sdfsdfsefs
local regexp = require("regexp")
local inspect = require("inspect")

-- regexp.match(regexp, data)
local result, err = regexp.match("hello", "hello world")
if err then error(err) end
if not(result==true) then error("regexp.match()") end`
	var tests = []struct {
		Name  string
		Input []string
		Lua   string
		Want  []string
	}{
		{"Regex", []string{}, regex, []string{"trudde"}},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%s", test.Name)
		t.Run(testname, func(t *testing.T) {
			filter := Filter{
				Parents: []*Filter{
					{
						Results: test.Input,
						Var1:    test.Lua,
					},
				},
			}
			getFilterResultLua(
				&filter,
			)
			if len(filter.Logs) > 0 {
				t.Errorf("Lua error: %s", filter.Logs)
			}
			if (filter.Results != nil && test.Want != nil) && !reflect.DeepEqual(test.Want, filter.Results) {
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
	processFilters(filters, nil, nil, false, nil)

	if !reflect.DeepEqual(filter1.Results, []string{helloWorld}) {
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
	processFilters(filters, nil, nil, false, nil)

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
	processFilters(filters, nil, nil, false, nil)

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
	processFilters(filters, nil, nil, true, nil)

	if !reflect.DeepEqual(minFilter.Results, []string{"100.000000"}) {
		t.Errorf("%s did not match '100'", minFilter.Results)
	}
	if !reflect.DeepEqual(maxFilter.Results, []string{"400.000000"}) {
		t.Errorf("%s did not match '400'", maxFilter.Results)
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
	processFilters(filters, nil, nil, false, nil)

	if len(filter.Logs) == 0 {
		t.Errorf("Expected error message in filter log, found empty log: %s", filter.Logs)
	}
}
