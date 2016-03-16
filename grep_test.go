package grep

import (
	"testing"
	"strconv"
)

func BenchmarkGrep(b *testing.B) {
	for i := 0; i < b.N; i++ {
		total := 0
		test := NewFileGrep("/proc/self/smaps", "Pss:\\s+(\\d+)\\skB")
		for _, result := range(test.Search().Result) {
			size := result.RegExp.FindStringSubmatch(string(result.Line))
			i, _ := strconv.Atoi(size[1])
			total += i
		}
		b.Logf("Total %d kb\n", total)
	}
}

func TestGrepPattern(t *testing.T) {

	test := NewFileGrep("/proc/self/*", "Pss:\\s+(\\d+)\\skB")
	if (len(test.Files) <= 0) {
		t.Fail()
	} else {
		t.Logf("Found %d files", len(test.Files))
	}

}