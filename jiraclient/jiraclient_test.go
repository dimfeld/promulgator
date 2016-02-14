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
		actual, err := parseCommon(input)
		if err != nil {
			t.Errorf("%s: saw unexpected error %s", input, err.Error())
		}
		if actual != expected {
			t.Errorf("%s: saw %v, expected %v", input, actual, expected)
		}
	}
}

func TestParseCommonErrors(t *testing.T) {
	cases := []string{
		"apple /../application-properties cow dog",
		"apple / cow",
		"apple /",
		"apple bee#cow",
		"apple bee?page=1",
		"apple bee&cow",
	}

	for _, input := range cases {
		_, err := parseCommon(input)
		if err == nil {
			t.Errorf("%s: expected error but saw none", input)
		}
	}
}