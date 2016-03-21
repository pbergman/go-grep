## GoGrep

A simple library that i use for grepping/searching files, similar to grep but this is fully written in go.

small example:

```
if grep, errors := NewGoGrep("Pss:\\s+(\\d+)\\skB", "/proc/self/smaps"); errors.HasErrors() {
    panic(errors) 
} else {
    defer grep.Close()
    for _, result := range(grep.Search().Result) {
        size := result.RegExp.FindStringSubmatch(string(result.Line))
        i, _ := strconv.Atoi(size[1])
        total += i
    }
    fmt.Sprintf("Total %d kb\n", total)
}
```

for more example see tests.