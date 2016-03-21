package gogrep

import (
	"strconv"
	"testing"
	"io/ioutil"
	"os"
	"fmt"
)

func BenchmarkGoGrep(b *testing.B) {
	test, _ := NewGoGrep("Pss:\\s+(\\d+)\\skB", "/proc/self/smaps")
	for i := 0; i < b.N; i++ {
		// Rest pointer for each iteration
		test.Reset()
		total := 0
		for _, result := range test.Search().Result {
			size := result.RegExp.FindStringSubmatch(string(result.Line))
			i, _ := strconv.Atoi(size[1])
			total += i
		}
		b.Logf("Total %d kb\n", total)
	}
}


func ExampleGoGrep() {
	content := []byte(`Lorem ipsum dolor sit amet, consectetur adipiscing elit. In vel augue a urna tempor pellentesque. Curabitur in odio non est ornare dapibus. Integer viverra ornare elit, sed sodales magna vehicula sit amet. Ut id dapibus dui. Maecenas ornare, sem et scelerisque bibendum, urna purus gravida urna, vel condimentum justo purus a purus. Maecenas sit amet sem eros. Ut sem libero, pellentesque vel vehicula ac, ultricies non nisi. Aliquam faucibus urna id lorem dictum, et tincidunt ligula fermentum.
	Donec elementum pharetra arcu, vel pulvinar lacus condimentum vitae. Pellentesque sed dolor finibus, dignissim nisl et, scelerisque orci. In quis feugiat orci. Phasellus pellentesque metus diam, id tincidunt augue luctus non. In sit amet dui turpis. Mauris imperdiet ligula at nibh tristique, sed faucibus turpis blandit. In fringilla erat turpis, ac cursus erat mollis sed.
	Pellentesque eleifend mattis egestas. In euismod lorem placerat nulla tempor, eu venenatis tellus facilisis. Aenean sodales consequat mollis. Morbi nec turpis sit amet metus vulputate lobortis. Donec accumsan erat nec lectus placerat tristique. Nullam suscipit iaculis eros, non viverra sem laoreet sit amet. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. In nisl sapien, luctus id vestibulum nec, dignissim sit amet ante. Pellentesque eu erat diam. Aliquam et varius augue, quis sodales diam.
	Nam scelerisque dictum nisi ut sodales. In tempus augue sed sapien convallis, aliquet semper mauris luctus. Suspendisse gravida nisi ut risus fringilla, vulputate malesuada enim imperdiet. Proin eget feugiat massa. Maecenas in mi viverra, dignissim purus ut, tincidunt nisi. Suspendisse bibendum nec orci ac rutrum. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Pellentesque egestas leo ut leo sollicitudin, non luctus nunc mattis. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer aliquet pulvinar velit, id cursus metus lobortis vel. Maecenas fringilla lobortis ante, id tincidunt justo. Vivamus maximus, ex sit amet mollis fermentum, tellus mauris elementum tortor, non varius dolor urna et nunc.
	In in cursus risus, vitae cursus lacus. Mauris eget velit fringilla, porta mauris et, volutpat metus. Vestibulum velit tortor, hendrerit id pharetra quis, venenatis in lacus. Curabitur arcu quam, congue nec sapien id, accumsan porta dolor. Nulla gravida ex non ipsum tincidunt, vitae pulvinar justo maximus. Praesent vitae enim at augue blandit interdum eget at nisl. Duis finibus facilisis ante ac volutpat. In ullamcorper elit a nunc lobortis vestibulum. Proin eu purus odio. Curabitur rutrum pretium ligula, id blandit neque. Nunc in elementum dolor. Quisque mi lacus, commodo at placerat at, imperdiet vitae arcu.`)
	tmpfile, _ := ioutil.TempFile("", "lorem")
	defer os.Remove(tmpfile.Name())
	tmpfile.Write(content);
	tmpfile.Close()
	test, _ := NewGoGrep("que dictum nisi u", tmpfile.Name())
	for _, result := range test.Search().Result {
		match := result.RegExp.FindSubmatchIndex(result.Line)
		offset := 0
		length := len(result.Line)
		if len(result.Line) > match[1] {
			length = match[1]
		}
		if 0 < match[0] {
			offset = match[0]
		}
		fmt.Printf("Line: %d Match: \"%s\"\n", result.LineNumber, result.Line[offset:length])
	}
	test.Close()
	lorem, _ := NewGoGrep("(?i)lorem", tmpfile.Name())
	for _, result := range lorem.Search().Result {
		matches := result.RegExp.FindAllSubmatchIndex(result.Line, -1)
		for _, match := range matches {
			offset := 0
			length := len(result.Line)
			if len(result.Line) > match[1] + 4 {
				length = match[1] + 4
			}
			if 0 < match[0] - 4 {
				offset = match[0] - 4
			}
			fmt.Printf("Line: %d Match: \"%s\"\n", result.LineNumber, result.Line[offset:length])
		}
	}
	lorem.Close()
	// Output:
	// Line: 4 Match: "que dictum nisi u"
	// Line: 1 Match: "Lorem ips"
	// Line: 1 Match: " id lorem dic"
	// Line: 3 Match: "mod lorem pla"
	// Line: 4 Match: "is. Lorem ips"
}