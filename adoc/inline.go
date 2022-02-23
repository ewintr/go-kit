package adoc

type InlineElement interface {
	Text() string
}

type PlainText string

func (pt PlainText) Text() string { return string(pt) }

type StrongText string

func (st StrongText) Text() string { return string(st) }

type EmpText string

func (et EmpText) Text() string { return string(et) }

type StrongEmpText string

func (set StrongEmpText) Text() string { return string(set) }

type Link struct {
	url   string
	title string
}

func NewLink(url, title string) Link {
	return Link{
		url:   url,
		title: title,
	}
}

func (l Link) URL() string  { return l.url }
func (l Link) Text() string { return l.title }

type CodeText string

func (ct CodeText) Text() string { return string(ct) }
