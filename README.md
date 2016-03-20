## File Grep

A simple library that i use for grepping/searching files.  

It is similar as:

```
cat file | egrep pattern 
```

But this is fully written in go and uses goroutines for parsing files.

small example:

```
grep := NewFileGrep("Pss:\\s+(\\d+)\\skB", "/proc/self/smaps")
for _, result := range(grep.Search().Result) {
    size := result.RegExp.FindStringSubmatch(string(result.Line))
    i, _ := strconv.Atoi(size[1])
    total += i
}
fmt.Sprintf("Total %d kb\n", total)
```