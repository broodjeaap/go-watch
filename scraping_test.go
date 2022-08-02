package main

import (
	"fmt"
	"testing"
)

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
