package adoc_test

import (
	"fmt"
	"testing"
	"time"

	"ewintr.nl/go-kit/adoc"
	"ewintr.nl/go-kit/test"
)

func TestNew(t *testing.T) {
	one := "one"
	two := "two"
	three := "three"
	ptOne := adoc.PlainText(one)
	ptTwo := adoc.PlainText(two)
	ptThree := adoc.PlainText(three)
	for _, tc := range []struct {
		name  string
		input string
		exp   *adoc.ADoc
	}{
		{
			name: "empty",
			exp: &adoc.ADoc{
				Tags:     []adoc.Tag{},
				Language: adoc.LANGUAGE_EN,
			},
		},
		{
			name:  "title",
			input: "= Title",
			exp: &adoc.ADoc{
				Title:    "Title",
				Tags:     []adoc.Tag{},
				Language: adoc.LANGUAGE_EN,
			},
		},
		{
			name:  "header",
			input: "= Title\nT. Test\n2020-10-27\n:language:\tnl",
			exp: &adoc.ADoc{
				Title:    "Title",
				Author:   "T. Test",
				Language: adoc.LANGUAGE_NL,
				Tags:     []adoc.Tag{},
				Date:     time.Date(2020, time.October, 27, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:  "paragraphs",
			input: fmt.Sprintf("%s\n\n%s\n\n%s", one, two, three),
			exp: &adoc.ADoc{
				Tags:     []adoc.Tag{},
				Language: adoc.LANGUAGE_EN,
				Content: []adoc.BlockElement{
					adoc.Paragraph([]adoc.InlineElement{ptOne}),
					adoc.Paragraph([]adoc.InlineElement{ptTwo}),
					adoc.Paragraph([]adoc.InlineElement{ptThree}),
				},
			},
		},
		{
			name:  "subtitle",
			input: "== Subtitle",
			exp: &adoc.ADoc{
				Tags:     []adoc.Tag{},
				Language: adoc.LANGUAGE_EN,
				Content: []adoc.BlockElement{
					adoc.SubTitle("Subtitle"),
				},
			},
		},
		{
			name:  "code block",
			input: "----\nsome code\nmore code\n----",
			exp: &adoc.ADoc{
				Tags:     []adoc.Tag{},
				Language: adoc.LANGUAGE_EN,
				Content: []adoc.BlockElement{
					adoc.CodeBlock("some code\nmore code"),
				},
			},
		},
		{
			name:  "code block with empty lines",
			input: "----\nsome code\n\nmore code\n----",
			exp: &adoc.ADoc{
				Tags:     []adoc.Tag{},
				Language: adoc.LANGUAGE_EN,
				Content: []adoc.BlockElement{
					adoc.CodeBlock("some code\n\nmore code"),
				},
			},
		},
		{
			name:  "list",
			input: "* item 1\n* item 2\n* *item 3*\n",
			exp: &adoc.ADoc{
				Tags:     []adoc.Tag{},
				Language: adoc.LANGUAGE_EN,
				Content: []adoc.BlockElement{
					adoc.List{
						adoc.ListItem([]adoc.InlineElement{adoc.PlainText("item 1")}),
						adoc.ListItem([]adoc.InlineElement{adoc.PlainText("item 2")}),
						adoc.ListItem([]adoc.InlineElement{adoc.StrongText("item 3")}),
					},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			act := adoc.New(tc.input)

			test.Equals(t, tc.exp, act)
		})
	}
}

func TestParseInline(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input string
		exp   []adoc.InlineElement
	}{{
		name: "empty",
	},
		{
			name:  "plain",
			input: "some test text",
			exp: []adoc.InlineElement{
				adoc.PlainText("some test text")},
		},
		{
			name:  "strong",
			input: "*some strong text*",
			exp: []adoc.InlineElement{
				adoc.StrongText("some strong text"),
			},
		},
		{
			name:  "strong in plain",
			input: "some *strong* text",
			exp: []adoc.InlineElement{
				adoc.PlainText("some "),
				adoc.StrongText("strong"),
				adoc.PlainText(" text"),
			},
		},
		{
			name:  "emphasis",
			input: "_some emphasized text_",
			exp: []adoc.InlineElement{
				adoc.EmpText("some emphasized text"),
			},
		},
		{
			name:  "emphasis in plain",
			input: "some _emphasized_ text",
			exp: []adoc.InlineElement{
				adoc.PlainText("some "),
				adoc.EmpText("emphasized"),
				adoc.PlainText(" text"),
			},
		},
		{
			name:  "emp and strong in plain",
			input: "some _*special*_ text",
			exp: []adoc.InlineElement{
				adoc.PlainText("some "),
				adoc.StrongEmpText("special"),
				adoc.PlainText(" text"),
			},
		},
		{
			name:  "link",
			input: "a link[title] somewhere",
			exp: []adoc.InlineElement{
				adoc.PlainText("a "),
				adoc.NewLink("link", "title"),
				adoc.PlainText(" somewhere"),
			},
		},
		{
			name:  "code",
			input: "`command`",
			exp: []adoc.InlineElement{
				adoc.CodeText("command"),
			},
		},
		{
			name:  "code in plain",
			input: "some `code` in text",
			exp: []adoc.InlineElement{
				adoc.PlainText("some "),
				adoc.CodeText("code"),
				adoc.PlainText(" in text"),
			},
		},
		{
			name:  "link with underscore",
			input: "https://example.com/some_url[some url]",
			exp: []adoc.InlineElement{
				adoc.NewLink("https://example.com/some_url", "some url"),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			act := adoc.ParseInline(tc.input)

			test.Equals(t, tc.exp, act)
		})
	}
}
