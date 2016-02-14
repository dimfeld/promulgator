package jiraclient

import (
	"testing"
)

func TestParseCommon(t *testing.T) {
	cases := map[string]parsedCommand{
		"a b c": {"a", "B", "c"},
		"apple bee cow dog elephant fox": {"apple", "BEE", "cow dog elephant fox"},
		"apple bee":                      {"apple", "BEE", ""},
		"apple":                          {"apple", "", ""},
		"apple bee cow\ndog\nelephant": {"apple", "BEE", "cow\ndog\nelephant"},
		"": {"", "", ""},
	}

	for input, expected := range cases {
		actual := parseCommon(input)
		if actual != expected {
			t.Errorf("%s: saw %v, expected %v", input, actual, expected)
		}
	}
}
