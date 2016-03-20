package gogrep

import (
	"testing"
	"strconv"
)

func BenchmarkGoGrep(b *testing.B) {

	test, _ := NewGoGrep("Pss:\\s+(\\d+)\\skB", "/proc/self/smaps")

	for i := 0; i < b.N; i++ {
		// Rest pointer for each iteration
		test.Reset()
		total := 0
		for _, result := range(test.Search().Result) {
			size := result.RegExp.FindStringSubmatch(string(result.Line))
			i, _ := strconv.Atoi(size[1])
			total += i
		}
		b.Logf("Total %d kb\n", total)
	}
}