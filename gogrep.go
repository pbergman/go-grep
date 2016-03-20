package gogrep

import (
	"os"
	"regexp"
	"sync"
	"bytes"
	"io"
	"sort"
)

// GoGrep basic struct for for the file search
type GoGrep struct {
	Files        []*os.File     // stack of files to search
	Pattern      *regexp.Regexp // search pattern
	MaxWorkers   uint64         // amount of concurrent workers to start for processing lines
	MaxOpenFiles uint64         // concurrent open files to process, if value = 0 then no limit is set
}

// getMaxOpenFiles will return the set maxOpenFiles, if value is 0 it will be set to size of files slice
func (f *GoGrep) getMaxOpenFiles() uint64 {
	if limit := f.MaxOpenFiles; limit <= 0 {
		return uint64(len(f.Files))
	} else {
		return limit
	}
}

// Search the files for the given pattern
func (f *GoGrep) Search() *MatchResults {
	var workers sync.WaitGroup
	var readers sync.WaitGroup
	result := new(MatchResults)
	files := make(chan *os.File)
	queue := make(chan *struct {
		line  uint64
		file  string
		bytes []byte
	})
	for i := uint64(0); i < f.MaxWorkers; i++ {
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
				iteration := uint64(0)
				for {
					iteration++
					if line, err := readLine(file); err != nil {
						if err == io.EOF {
							break
						} else {
							panic(err)
						}
					} else {
						queue <- &struct {
							line  uint64
							file  string
							bytes []byte
						}{iteration, file.Name(), line}
					}

				}
			}
			readers.Done()
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

type line []byte

// readLine is a internal function that will read a byte till \n
// and than return the line, on error or eof it will return a error
func readLine(file *os.File) (line, error) {
	buffer := new(bytes.Buffer)
	bytes := make([]byte, 1)
	for {
		if _, err := file.Read(bytes); err != nil {
			return nil, err
		} else {
			if bytes[0] == '\n' {
				break
			}
			buffer.Write(bytes)
		}
	}
	return line(buffer.Bytes()), nil
}


func NewGoGrep(pattern string, files ...string) (*GoGrep, error) {
	instance := &GoGrep{
		Files:        make([]*os.File, len(files) - 1),
		Pattern:      regexp.MustCompile(pattern),
		MaxWorkers:   20,
		MaxOpenFiles: 0,
	}
	for _, file := range files {
		if stat, err := os.Stat(file); err != nil {
			return nil, err
		} else {
			if stat.Mode().IsRegular() {
				if of, err := os.Open(file); err != nil {
					return nil, err
				} else {
					instance.Files = append(instance.Files, of)
				}
			}
		}
	}
	return instance, nil
}

// Close all open files
func (f *GoGrep) Close() error {
	for _, file := range f.Files {
		if err := file.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Reset the current pointers of open files
func (f *GoGrep) Reset() error {
	for _, file := range f.Files {
		if _, err := file.Seek(0,0); err != nil {
			return err
		}
	}
	return nil
}
