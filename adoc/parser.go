package adoc

import (
	"strings"
	"time"
)

const (
	FIELD_DELIMITER     = ":"
	TITLE_PREFIX        = "= "
	SUBTITLE_PREFIX     = "== "
	SUBSUBTITLE_PREFIX  = "=== "
	PARAGRAPH_SEPARATOR = "\n\n"
	LINE_SEPARATOR      = "\n"
	CODE_PREFIX         = "----\n"
	CODE_SUFFIX         = "\n----"
	LISTITEM_PREFIX     = "* "
	IMAGE_PREFIX        = "image::" // end with [alt text]
)

func New(text string) *ADoc {
	doc := &ADoc{
		Language: LANGUAGE_EN,
		Tags:     []Tag{},
	}

	// split up blocks
	var pars []string
	for _, s := range strings.Split(text, PARAGRAPH_SEPARATOR) {
		if s == "" {

			continue
		}

		pars = append(pars, s)
	}

	// actually, code blocks are allowed to have empty lines
	var blocks []string
	var inCode bool
	var currentBlock string
	for _, par := range pars {
		switch {
		case strings.HasPrefix(par, CODE_PREFIX) && strings.HasSuffix(par, CODE_SUFFIX):
			blocks = append(blocks, par)
		case !inCode && strings.HasPrefix(par, CODE_PREFIX):
			inCode = true
			currentBlock = par
		case inCode && !strings.HasSuffix(par, CODE_SUFFIX):
			currentBlock += PARAGRAPH_SEPARATOR + par
		case inCode && strings.HasSuffix(par, CODE_SUFFIX):
			currentBlock += PARAGRAPH_SEPARATOR + par
			blocks = append(blocks, currentBlock)
			inCode = false
			currentBlock = ""
		default:
			blocks = append(blocks, par)
		}
	}

	// interpret the blocks
	for i, p := range blocks {
		switch {
		case i == 0 && strings.HasPrefix(p, TITLE_PREFIX):
			ParseHeader(p, doc)
		case strings.HasPrefix(p, SUBTITLE_PREFIX):
			p = strings.TrimSpace(p)
			s := strings.Split(p, SUBTITLE_PREFIX)
			if len(s) == 1 || s[1] == "" {

				continue
			}
			doc.Content = append(doc.Content, SubTitle(s[1]))
		case strings.HasPrefix(p, SUBSUBTITLE_PREFIX):
			p = strings.TrimSpace(p)
			s := strings.Split(p, SUBSUBTITLE_PREFIX)
			if len(s) == 1 || s[1] == "" {

				continue
			}
			doc.Content = append(doc.Content, SubSubTitle(s[1]))
		case isCodeBlock(p):
			doc.Content = append(doc.Content, parseCodeBlock(p))
		case strings.HasPrefix(p, LISTITEM_PREFIX):
			p = strings.TrimSpace(p)
			var items []ListItem
			for i, ti := range strings.Split(p, LISTITEM_PREFIX) {
				if i > 0 {
					inline := ParseInline(strings.TrimSpace(ti))
					items = append(items, ListItem(inline))
				}
			}
			doc.Content = append(doc.Content, List(items))
		case strings.HasPrefix(p, IMAGE_PREFIX):
			doc.Content = append(doc.Content, parseImage(p))
		default:
			p = strings.TrimSpace(p)
			doc.Content = append(doc.Content, Paragraph(ParseInline(p)))
		}
	}

	return doc
}

func isCodeBlock(par string) bool {
	return strings.HasPrefix(par, CODE_PREFIX) && strings.HasSuffix(par, CODE_SUFFIX)
}

func parseCodeBlock(par string) CodeBlock {
	ss := strings.Split(par, "\n")
	ss = ss[1 : len(ss)-1]
	content := strings.Join(ss, "\n")

	return CodeBlock(content)
}

func parseImage(img string) Image {
	ss := strings.TrimPrefix(img, IMAGE_PREFIX)
	if strings.Index(ss, "[") >= 0 && strings.HasSuffix(ss, "]") {
		parts := strings.Split(ss, "[")
		return Image{
			Src: parts[0],
			Alt: strings.TrimSuffix(parts[1], "]"),
		}
	}

	return Image{
		Src: ss,
	}
}

func ParseHeader(text string, doc *ADoc) {
	text = strings.TrimSpace(text)
	lines := strings.Split(text, LINE_SEPARATOR)
	for i, l := range lines {
		switch {
		case i == 0:
			s := strings.Split(l, TITLE_PREFIX)
			doc.Title = s[1]
		case isDate(l):
			date, _ := time.Parse("2006-01-02", l)
			doc.Date = date
		case isField(l):
			name, value := parseField(l)
			switch name {
			case "language":
				doc.Language = NewLanguage(value)
			case "updated":
				updateDate, _ := time.Parse("2006-01-02", value)
				doc.Updated = updateDate
			case "tags":
				for _, tag := range strings.Split(value, ",") {
					doc.Tags = append(doc.Tags, Tag(strings.TrimSpace(tag)))
				}
			}
		default:
			doc.Author = l
		}
	}
}

func isField(line string) bool {
	parts := strings.Split(line, FIELD_DELIMITER)
	if len(parts) != 3 || parts[0] != "" {
		return false
	}

	return true
}

func parseField(line string) (string, string) {
	parts := strings.Split(line, FIELD_DELIMITER)
	if len(parts) != 3 {
		return "", ""
	}
	trimmedParts := []string{}
	for _, part := range parts {
		trimmedParts = append(trimmedParts, strings.TrimSpace(part))
	}
	if trimmedParts[0] != "" {
		return "", ""
	}

	return trimmedParts[1], trimmedParts[2]
}

func isDate(text string) bool {
	if _, err := time.Parse("2006-01-02", text); err == nil {
		return true
	}

	return false
}

func ParseInline(text string) []InlineElement {
	var e []InlineElement

	ss := strings.Split(text, "")
	var buffer, curWord, prevChar string
	var strong, emp, code, linkTitle bool
	wordStart := true
	for _, s := range ss {
		switch {
		case (s == "_" && wordStart) || (s == "_" && emp):
			e = addElement(e, buffer+curWord, strong, emp, code)
			emp = !emp
			buffer = ""
			curWord = ""
		case s == "*":
			e = addElement(e, buffer+curWord, strong, emp, code)
			buffer = ""
			curWord = ""
			strong = !strong
		case s == "`":
			e = addElement(e, buffer+curWord, strong, emp, code)
			code = !code
			buffer = ""
			curWord = ""
		case s == "[" && prevChar != "":
			e = addElement(e, buffer, strong, emp, code)
			linkTitle = true
			curWord += s
		case s == "]" && linkTitle:
			e = addLink(e, curWord)
			buffer = ""
			curWord = ""
			linkTitle = false
		case s == " " && !linkTitle:
			buffer += curWord + " "
			curWord = ""
		default:
			curWord += s
		}
		prevChar = s
		wordStart = false
		if prevChar == " " {
			wordStart = true
		}
	}
	if len(buffer+curWord) > 0 {
		e = addElement(e, buffer+curWord, false, false, false)
	}

	return e
}

func addLink(ies []InlineElement, linkText string) []InlineElement {
	ss := strings.Split(linkText, "[")
	if len(ss) < 2 {
		ss = append(ss, "ERROR")
	}

	return append(ies, Link{url: ss[0], title: ss[1]})
}

func addElement(ies []InlineElement, text string, strong, emp, code bool) []InlineElement {
	if len(text) == 0 {
		return ies
	}

	var ne InlineElement
	switch {
	case code:
		ne = CodeText(text)
	case strong && emp:
		ne = StrongEmpText(text)
	case strong && !emp:
		ne = StrongText(text)
	case !strong && emp:
		ne = EmpText(text)
	default:
		ne = PlainText(text)
	}

	return append(ies, ne)
}
