package replacer

import (
	"regexp"
	"strings"

	"github.com/fatih/color"
)

var (
	red = color.New(color.FgRed, color.Bold).SprintFunc()
)

type Matcher interface {
	match(s string) bool
	replace(s string) string
	colorizeFrom(s string) string
}

type StringMatcher struct {
	from string
	to   string
}

type RegexpMatcher struct {
	from   string
	to     string
	fromRe *regexp.Regexp
}

func (sm *StringMatcher) match(s string) bool {
	return strings.Contains(s, sm.from)
}

func (sm *StringMatcher) replace(s string) string {
	return strings.ReplaceAll(s, sm.from, sm.to)
}

func (sm *StringMatcher) colorizeFrom(s string) string {
	return strings.ReplaceAll(s, sm.from, red(sm.from))
}

func (rm *RegexpMatcher) match(s string) bool {
	if rm.fromRe == nil {
		rm.fromRe = regexp.MustCompile(rm.from)
	}
	return rm.fromRe.MatchString(s)
}

func (rm *RegexpMatcher) replace(s string) string {
	if rm.fromRe == nil {
		rm.fromRe = regexp.MustCompile(rm.from)
	}
	return string(rm.fromRe.ReplaceAllString(s, rm.to))
}

func (rm *RegexpMatcher) colorizeFrom(s string) string {
	if rm.fromRe == nil {
		rm.fromRe = regexp.MustCompile(rm.from)
	}
	return string(rm.fromRe.ReplaceAllString(s, red(rm.from)))
}
