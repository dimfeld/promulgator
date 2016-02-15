package jiraclient

import (
	"testing"
)

func TestParseCommon(t *testing.T) {
	cases := map[string]parsedCommand{
		"a b c": {"a", "b", "c"},
		"apple bee cow dog elephant fox":   {"apple", "bee", "cow dog elephant fox"},
		"apple  bee cow dog elephant fox":  {"apple", "bee", "cow dog elephant fox"},
		"apple BEE  cow dog elephant fox":  {"apple", "BEE", "cow dog elephant fox"},
		"apple  bee  cow DOG elephant fox": {"apple", "bee", "cow DOG elephant fox"},
		"apple bee":                        {"apple", "bee", ""},
		"apple  bee":                       {"apple", "bee", ""},
		"APPLE":                            {"APPLE", "", ""},
		"apple bee        cow\ndog\nelephant": {"apple", "bee", "cow\ndog\nelephant"},
		"": {"", "", ""},
	}

	for input, expected := range cases {
		t.Logf("Testing %s", input)
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
