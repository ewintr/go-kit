package adoc

import (
	"fmt"
	"strings"
)

type BlockElement interface {
	Text() string
}

type Paragraph []InlineElement

func (p Paragraph) Text() string {
	var text []string
	for _, ie := range p {
		text = append(text, ie.Text())
	}

	return strings.Join(text, " ")
}

type SubTitle string

func (st SubTitle) Text() string { return string(st) }

type SubSubTitle string

func (st SubSubTitle) Text() string { return string(st) }

type CodeBlock string

func (cb CodeBlock) Text() string { return string(cb) }

type Image struct {
	Src string
	Alt string
}

func (i Image) Text() string {
	return i.Alt
}

type ListItem []InlineElement

func (li ListItem) Text() string {
	var text []string
	for _, ie := range li {
		text = append(text, ie.Text())
	}

	return fmt.Sprintf("%s%s", LISTITEM_PREFIX, strings.Join(text, ""))
}

type List []ListItem

func (l List) Text() string {
	var items []string
	for _, item := range l {
		items = append(items, item.Text())
	}

	return strings.Join(items, "\n")
}
