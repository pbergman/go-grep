## File Grep

A simple library that i use for grepping/searching config files.  

It is similar as:

```
cat file | egrep pattern 
```

But this is fully written in go and uses goroutines for opening files
and parsing files.

small example:

```
grep := NewFileGrep("/proc/self/smaps", "Pss:\\s+(\\d+)\\skB")
for _, result := range(grep.Search().Result) {
    size := result.RegExp.FindStringSubmatch(string(result.Line))
    i, _ := strconv.Atoi(size[1])
    total += i
}
fmt.Sprintf("Total %d kb\n", total)
```

the file or directory can be a pattern like /a/b/*/d or /a/b/*.csv
see match and glob from the filepath package.