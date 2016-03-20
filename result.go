package gogrep

import (
	"fmt"
	"regexp"
	"sync"
)

// MatchResults is basic struct for a collection of results
type MatchResults struct {
	Result []*MatchResult
	lock   sync.Mutex
}

// add is a internal function for adding a result
func (r *MatchResults) add(m *MatchResult) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.Result = append(r.Result, m)
}

// String implementation
func (f *MatchResults) String() string {
	var s string
	for _, result := range f.Result {
		match := result.RegExp.FindSubmatchIndex(result.Line)
		offset := 0
		length := len(result.Line)
		if len(result.Line) > match[1] {
			length = match[1]
		}
		if 0 < match[0] {
			offset = match[0]
		}
		s += fmt.Sprintf(
			"%s(%d) %s\n",
			result.File,
			result.LineNumber,
			result.Line[offset:length],
		)
	}
	return s
}

// MatchResult is basic struct for the a result
type MatchResult struct {
	File       string
	LineNumber uint64
	Line       []byte
	RegExp     *regexp.Regexp
}
