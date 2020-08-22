Forked for those simple improvements:

* Fixed all code warnings reported by [Goland](https://www.jetbrains.com/go/)
* Reformatted code
* Renamed `columnize.Columnize()` to `columnize.Format()` and `columnize.FormatStringList()` 
to fit in with Golang's naming of things, as well as the now named `Options` type
* Improved `columnize_test.go` and `columnize-demo/main.go` a bit and adjusted to new naming
