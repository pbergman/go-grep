package grep

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"sync"
)

// FileGrep basic struct for for the file search
type FileGrep struct {
	Files        []string       // stack of files to search
	Pattern      *regexp.Regexp // search pattern
	Concurrent   uint64         // amount of workers to start for processing lines
	MaxOpenFiles uint64         // Allowed open files, <= 0 then no limit is set
	lock         sync.Mutex     // internal lock for synchronous appending files
}

// addFile internal function for adding
func (f *FileGrep) addFile(file string) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.Files = append(f.Files, file)
}

func (f *FileGrep) getMaxOpenFiles() uint64 {
	if limit := f.MaxOpenFiles; limit <= 0 {
		return uint64(len(f.Files))
	} else {
		return limit
	}
}

// Search will search in a asynchronous way the registerd files
func (f *FileGrep) Search() *MatchResults {
	var readers sync.WaitGroup
	var workers sync.WaitGroup
	result := new(MatchResults)
	files := make(chan string)
	queue := make(chan *struct {
		line  uint64
		file  string
		bytes []byte
	})
	for i := uint64(0); i < f.Concurrent; i++ {
		workers.Add(1)
		go func() {
			for item := range queue {
				if f.Pattern.Match(item.bytes) {
					result.add(&MatchResult{
						File:       item.file,
						LineNumber: item.line,
						Line:       item.bytes,
						RegExp:     f.Pattern,
					})
				}
			}
			workers.Done()
		}()
	}
	for i := uint64(0); i < f.getMaxOpenFiles(); i++ {
		readers.Add(1)
		go func() {
			for file := range files {
				currentLine := uint64(0)
				openFile, err := os.Open(file)
				checkError(err)
				defer openFile.Close()
				reader := bufio.NewReader(openFile)
				for {
					if line, err := reader.ReadBytes('\n'); err != nil {
						if err == io.EOF {
							break
						} else {
							panic(err)
						}
					} else {
						currentLine++
						queue <- &struct {
							line  uint64
							file  string
							bytes []byte
						}{currentLine, file, line}
					}
				}
				readers.Done()
				openFile.Close()
			}
		}()
	}
	for _, file := range f.Files {
		files <- file
	}
	close(files)
	readers.Wait()
	close(queue)
	workers.Wait()
	sort.Sort(ResultSort(result.Result))
	return result
}

// NewFileGrep will create a new FileGrep instance and check given patterns and file
// give file can be a file, folder or pattern (see filepath.Glob). If file is folder
// or pattern resolves to folders it will search all the folder for files (non-recursive)
func NewFileGrep(file_path, pattern string) *FileGrep {
	var readers sync.WaitGroup
	instance := &FileGrep{
		Files:        make([]string, 0),
		Pattern:      regexp.MustCompile(pattern),
		Concurrent:   20,
		MaxOpenFiles: 0,
	}
	// check for given patterns , like "*.txt"
	files, err := filepath.Glob(file_path)
	checkError(err)
	for _, file := range files {
		readers.Add(1)
		go func(file string) {
			fileStat, err := os.Stat(file)
			checkError(err)
			if fileStat.IsDir() {
				openFolder, err := os.Open(file)
				defer openFolder.Close()
				checkError(err)
				foundFile, err := openFolder.Readdir(-1)
				for _, foundFileInfo := range foundFile {
					if !foundFileInfo.IsDir() {
						instance.addFile(file + foundFileInfo.Name())
					}
				}
				openFolder.Close()
			} else {
				instance.addFile(file)
			}

			readers.Done()
		}(file)
	}
	readers.Wait()
	return instance
}
