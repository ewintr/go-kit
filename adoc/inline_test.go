package adoc_test

import (
	"testing"

	"ewintr.nl/go-kit/adoc"
	"ewintr.nl/go-kit/test"
)

func TestInlineSimple(t *testing.T) {
	text := "text"
	for _, tc := range []struct {
		name    string
		element adoc.InlineElement
	}{
		{
			name:    "plain text",
			element: adoc.PlainText(text),
		},
		{
			name:    "strong",
			element: adoc.StrongText(text),
		},
		{
			name:    "emphasis",
			element: adoc.EmpText(text),
		},
		{
			name:    "strong emphasis",
			element: adoc.StrongEmpText(text),
		},
		{
			name:    "code",
			element: adoc.CodeText(text),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			test.Equals(t, text, tc.element.Text())
		})
	}
}

func TextLink(t *testing.T) {
	url := "url"
	title := "title"
	l := adoc.NewLink(url, title)

	test.Equals(t, url, l.URL())
	test.Equals(t, title, l.Text())
}
