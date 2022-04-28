package replacer

import (
	"regexp"
	"strings"
)

type Matcher interface {
	match(s string) bool
	replace(s string) string
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

func (rm *RegexpMatcher) match(s string) bool {
	if rm.fromRe == nil {
		rm.fromRe = regexp.MustCompile(rm.from)
	}
	return rm.fromRe.Match([]byte(s))
}

func (rm *RegexpMatcher) replace(s string) string {
	if rm.fromRe == nil {
		rm.fromRe = regexp.MustCompile(rm.from)
	}
	return string(rm.fromRe.ReplaceAll([]byte(s), []byte(rm.to)))
}
