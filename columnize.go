// Copyright 2013 Rocky Bernstein.
//
// The functions toStringArrayFromIndexable and toStringArray are
// from Carlos Castillo. Thanks Carlos!
// http://play.golang.org/p/bxdcIj6ueH
//
// Adapted from the routine of the same name from Ruby.

package columnize

import (
	"fmt"
	"reflect"
)

// Use DefaultOptions() or ArrayOptions() to retrieve an object with sane defaults.
type Options struct {
	ArrangeVertical bool
	ArrayPrefix     string
	ArraySuffix     string
	CellFmt         string
	ColSep          string
	DisplayWidth    int
	LinePrefix      string
	LineSuffix      string
	LJustify        bool
}

func DefaultOptions() Options {
	return Options{
		ArrangeVertical: true,
		ArrayPrefix:     "",
		ArraySuffix:     "",
		CellFmt:         "",
		ColSep:          "  ",
		DisplayWidth:    80,
		LinePrefix:      "",
		LineSuffix:      "\n",
		LJustify:        true,
	}
}

func ArrayOptions() Options {
	return Options{
		ArrangeVertical: false,
		ArrayPrefix:     "[",
		ArraySuffix:     "]",
		CellFmt:         `"%v"`,
		ColSep:          ", ",
		DisplayWidth:    80,
		LinePrefix:      "",
		LineSuffix:      ", ",
		LJustify:        true,
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

/*
toStringSliceFromIndexable runs fmt.Sprint on each of the passed elements in x.
If optFmt is given, that is the format string passed to fmt.Sprintf.

This function assumes that x is a slice or an array.
No checking on or error is thrown if this is not the case.
*/
func toStringSliceFromIndexable(x interface{}, optFmt string) []string {
	v := reflect.ValueOf(x)
	out := make([]string, v.Len())
	for i := range out {
		if optFmt == "" {
			out[i] = fmt.Sprint(v.Index(i).Interface())
		} else {
			out[i] = fmt.Sprintf(optFmt, v.Index(i).Interface())
		}
	}
	return out
}

/*
toStringSlice uses toStringSliceFromIndexable if the passed data is a slice or array
or just converts the passed data via fmt.Sprintf.
*/
func toStringSlice(x interface{}, optFmt string) []string {
	v := reflect.ValueOf(x)
	if v.Kind() != reflect.Array && v.Kind() != reflect.Slice {
		// Not an array or slice, so run fmt.Sprint and turn that
		// single item into a slice.
		if optFmt == "" {
			return []string{fmt.Sprint(x)}
		}
		return []string{fmt.Sprintf(optFmt, x)}
	}
	return toStringSliceFromIndexable(x, optFmt)
}

/*
Format() returns a string from an array with embedded newlines formatted so
that when printed the columns are aligned.

For example:, for a line width of 4 characters (arranged vertically):

	a = [] int{1,2,4,4}
	columnize.Format(a) => '1  3\n2  4\n'

Arranged horizontally:

	opts := columnize.Default_options()
	opts.arrange_vertical = false
	columnize.Format(a) =>  '1  2\n3  4\n'

Each column is only as wide as necessary.
By default, columns are separated by two spaces.
*/
func Format(list interface{}, opts Options) string {
	l := toStringSlice(list, opts.CellFmt)
	return FormatStringList(l, opts)
}

// FormatStringList works like Format(), but only accepts a string slice
func FormatStringList(list []string, opts Options) string {

	var prefix string
	if len(opts.ArrayPrefix) == 0 {
		prefix = opts.LinePrefix
	} else {
		prefix = opts.ArrayPrefix
	}
	if len(list) == 0 {
		result :=
			fmt.Sprintf("%s%s",
				prefix, opts.ArraySuffix)
		return result
	}

	if len(list) == 1 {
		result :=
			fmt.Sprintf("%s%s%s",
				prefix, list[0], opts.ArraySuffix)
		return result
	}

	if opts.DisplayWidth-len(opts.LinePrefix) < 4 {
		opts.DisplayWidth = len(opts.LinePrefix) + 4
	} else {
		opts.DisplayWidth -= len(opts.LinePrefix)
	}

	var ncols, nrows int

	if opts.ArrangeVertical {
		arrayIndex := func(numRows, row, col int) int {
			return numRows*col + row
		}
		var colwidths []int
		// Try every row count from 1 upwards
		for nrows = 1; nrows < len(list); nrows++ {
			ncols = (len(list) + nrows - 1) / nrows
			totwidth := -len(opts.ColSep)
			colwidths = make([]int, 0)

			for col := 0; col < ncols; col++ {
				// get max column width for this column
				colwidth := 0
				for row := 0; row < nrows; row++ {
					i := arrayIndex(nrows, row, col)
					if i >= len(list) {
						break
					}
					colwidth = max(len(list[i]), colwidth)
				}
				colwidths = append(colwidths, colwidth)
				totwidth += colwidth + len(opts.ColSep)
				if totwidth > opts.DisplayWidth {
					ncols = col
					break
				}
			}
			if totwidth <= opts.DisplayWidth {
				break
			}
		}
		if ncols < 1 {
			ncols = 1
		}
		if ncols == 1 {
			nrows = len(list)
		}
		// The smallest number of rows computed and the max widths for
		// each column has been obtained.
		// Now we just have to format each of the rows.
		s := ""
		for row := 0; row < nrows; row++ {
			texts := make([]string, 0)
			for col := 0; col < ncols; col++ {
				var x string
				i := arrayIndex(nrows, row, col)
				if i >= len(list) {
					x = ""
				} else {
					x = list[i]
				}
				texts = append(texts, x)
			}
			// texts.pop while !texts.empty? and texts[-1] == ''
			if len(texts) > 0 {
				for col := 0; col < len(texts); col++ {
					if ncols != 1 {
						var fmtStr string
						if opts.LJustify {
							fmtStr = fmt.Sprintf("%%%ds", -colwidths[col])
							texts[col] = fmt.Sprintf(fmtStr, texts[col])
						} else {
							fmtStr = fmt.Sprintf("%%%ds", colwidths[col])
							texts[col] = fmt.Sprintf(fmtStr, texts[col])
						}
					}
				}
				line := opts.LinePrefix
				for i := 0; i < len(texts)-1; i++ {
					line += fmt.Sprintf("%s%s", texts[i], opts.ColSep)
				}
				if len(texts) > 0 {
					line += fmt.Sprintf("%s%s", texts[len(texts)-1], opts.LineSuffix)
				}
				s += line
			}
		}
		return s
	} else {
		var colwidths []int
		arrayIndex := func(ncols, row, col int) int {
			return ncols*(row-1) + col
		}
		// Assign to make enlarge scope of loop variables.
		var totalWidth, i, roundedSize int
		var ncols, nrows int
		// Try every column count from size downwards.
		for ncols = len(list); ncols >= 1; ncols-- {
			// Try every row count from 1 upwards
			minRows := (len(list) + ncols - 1) / ncols
			for nrows = minRows; nrows <= (len(list)); nrows++ {
				roundedSize = nrows * ncols
				colwidths = make([]int, 0)
				totalWidth = -len(opts.ColSep)
				var colwidth, row int
				for col := 0; col < ncols; col++ {
					// get max column width for this column
					for row = 1; row <= nrows; row++ {
						i = arrayIndex(ncols, row, col)
						if i >= len(list) {
							break
						}
						colwidth = max(colwidth, len(list[i]))
					}
					colwidths = append(colwidths, colwidth)
					totalWidth += colwidth + len(opts.ColSep)
					if totalWidth > opts.DisplayWidth {
						break
					}
				}
				if totalWidth <= opts.DisplayWidth {
					// Found the right nrows and ncols
					nrows = row
					break
				} else {
					if totalWidth > opts.DisplayWidth {
						// Need to reduce ncols
						break
					}
				}
			}
			if totalWidth <= opts.DisplayWidth && i >= roundedSize-1 {
				break
			}
		}
		if ncols < 1 {
			ncols = 1
		}
		if ncols == 1 {
			nrows = len(list)
		}
		// The smallest number of rows computed and the max widths for
		// each column has been obtained.  Now we just have to format
		// each of the rows.
		s := ""
		var prefix string
		if len(opts.ArrayPrefix) == 0 {
			prefix = opts.LinePrefix
		} else {
			prefix = opts.ArrayPrefix
		}
		for row := 1; row <= nrows; row++ {
			texts := make([]string, 0)
			for col := 0; col < ncols; col++ {
				var x string
				i = arrayIndex(ncols, row, col)
				if i >= len(list) {
					break
				} else {
					x = list[i]
				}
				texts = append(texts, x)
			}
			for col := 0; col < len(texts); col++ {
				if ncols != 1 {
					var fmtStr string
					if opts.LJustify {
						fmtStr = fmt.Sprintf("%%%ds", -colwidths[col])
						texts[col] = fmt.Sprintf(fmtStr, texts[col])
					} else {
						fmtStr = fmt.Sprintf("%%%ds", colwidths[col])
						texts[col] = fmt.Sprintf(fmtStr, texts[col])
					}
				}
			}
			line := prefix
			for i := 0; i < len(texts)-1; i++ {
				line += fmt.Sprintf("%s%s", texts[i], opts.ColSep)
			}
			if len(texts) > 0 {
				line += fmt.Sprintf("%s%s", texts[len(texts)-1], opts.LineSuffix)
			}
			s += line
			prefix = opts.LinePrefix
		}
		s += opts.ArraySuffix
		return s
	}
}
