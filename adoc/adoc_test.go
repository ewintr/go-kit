package adoc_test

import (
	"testing"

	"ewintr.nl/go-kit/adoc"
	"ewintr.nl/go-kit/test"
)

func TestNewLanguage(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input string
		exp   adoc.Language
	}{
		{
			name: "empty",
			exp:  adoc.LANGUAGE_UNKNOWN,
		},
		{
			name:  "dutch lower",
			input: "nl",
			exp:   adoc.LANGUAGE_NL,
		},
		{
			name:  "dutch upper",
			input: "NL",
			exp:   adoc.LANGUAGE_NL,
		},
		{
			name:  "english lower",
			input: "en",
			exp:   adoc.LANGUAGE_EN,
		},
		{
			name:  "english upper",
			input: "EN",
			exp:   adoc.LANGUAGE_EN,
		},
		{
			name:  "unknown",
			input: "something",
			exp:   adoc.LANGUAGE_UNKNOWN,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			act := adoc.NewLanguage(tc.input)
			test.Equals(t, tc.exp, act)
		})

	}
}
