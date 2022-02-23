package adoc_test

import (
	"testing"

	"ewintr.nl/go-kit/adoc"
	"ewintr.nl/go-kit/test"
)

func TestParagraph(t *testing.T) {
	for _, tc := range []struct {
		name     string
		elements []adoc.InlineElement
		exp      string
	}{
		{
			name:     "empty",
			elements: []adoc.InlineElement{},
		},
		{
			name: "one",
			elements: []adoc.InlineElement{
				adoc.PlainText("one"),
			},
			exp: "one",
		},
		{
			name: "many",
			elements: []adoc.InlineElement{
				adoc.PlainText("one"),
				adoc.PlainText("two"),
				adoc.PlainText("three"),
			},
			exp: "one two three",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := adoc.Paragraph(tc.elements)
			test.Equals(t, tc.exp, p.Text())
		})
	}
}

func TestBlockSimple(t *testing.T) {
	text := "text"
	for _, tc := range []struct {
		name    string
		element adoc.BlockElement
	}{
		{
			name:    "subtitle",
			element: adoc.SubTitle(text),
		},
		{
			name:    "subsubtitle",
			element: adoc.SubSubTitle(text),
		},
		{
			name:    "code block",
			element: adoc.CodeBlock(text),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			test.Equals(t, text, tc.element.Text())
		})
	}
}

func TestListItem(t *testing.T) {
	for _, tc := range []struct {
		name     string
		elements []adoc.InlineElement
		exp      string
	}{
		{
			name: "empty",
			exp:  "* ",
		},
		{
			name: "one",
			elements: []adoc.InlineElement{
				adoc.PlainText("one"),
			},
			exp: "* one",
		},
		{
			name: "many",
			elements: []adoc.InlineElement{
				adoc.PlainText("one"),
				adoc.PlainText("two"),
				adoc.PlainText("three"),
			},
			exp: "* onetwothree",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			li := adoc.ListItem(tc.elements)
			test.Equals(t, tc.exp, li.Text())
		})
	}
}

func TestList(t *testing.T) {
	for _, tc := range []struct {
		name     string
		elements []adoc.ListItem
		exp      string
	}{
		{
			name: "empty",
		},
		{
			name: "one",
			elements: []adoc.ListItem{
				{adoc.PlainText("one")},
			},
			exp: "* one",
		},
		{
			name: "many",
			elements: []adoc.ListItem{
				{adoc.PlainText("one")},
				{adoc.PlainText("two")},
				{adoc.PlainText("three")},
			},
			exp: "* one\n* two\n* three",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			l := adoc.List(tc.elements)
			test.Equals(t, tc.exp, l.Text())
		})
	}
}
