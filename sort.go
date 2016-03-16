package grep

// ResultSort custom implementation of sort.Interface for sorting the results
type ResultSort []*MatchResult

// Len @sort.Interface
func (s ResultSort) Len() int {
	return len(s)
}

// Swap @sort.Interface
func (s ResultSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less will sort file names and line numbers so
// the are grouped by filename and line number
func (s ResultSort) Less(i, j int) bool {
	if s[i].File == s[j].File {
		return s[i].LineNumber < s[j].LineNumber
	} else {
		return s[i].File < s[j].File
	}
}
