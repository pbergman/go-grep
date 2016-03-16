package grep

import (
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

// MatchResult is basic struct for the a result
type MatchResult struct {
	File       string
	LineNumber uint64
	Line       []byte
	RegExp     *regexp.Regexp
}
