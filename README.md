Forked for those simple improvements:

* Fixed all code warnings reported by [Goland](https://www.jetbrains.com/go/) apart 
from the `Indexing may panic because of 'nil' slice` warnings on the `colwidths` slice
  * Those warnings were surpressed by `//goland:noinspection GoNilness` because the code 
  is certainly correct
* Reformatted code
* Renamed `columnize.Columnize()` to `columnize.Format()` and `columnize.FormatStringList()` 
to fit in with Golang's naming of things, as well as the, now named, `Options` type
* Improved `columnize_test.go` a bit
