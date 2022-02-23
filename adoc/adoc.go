package adoc

import (
	"strings"
	"time"
)

const (
	LANGUAGE_EN      = Language("en")
	LANGUAGE_NL      = Language("nl")
	LANGUAGE_UNKNOWN = Language("unknown")
)

type Language string

func NewLanguage(ln string) Language {
	switch strings.ToLower(ln) {
	case "nl":
		return LANGUAGE_NL
	case "en":
		return LANGUAGE_EN
	default:
		return LANGUAGE_UNKNOWN
	}
}

type Tag string

type ADoc struct {
	Title    string
	Author   string
	Language Language
	Public   bool
	Path     string
	Date     time.Time
	Updated  time.Time
	Tags     []Tag
	Content  []BlockElement
}
