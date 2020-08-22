// Copyright 2013 Rocky Bernstein.
//
// Adapted from the routine of the same name from Ruby.

package columnize

import (
	"fmt"
	"reflect"
)

type OptsT struct {
	ArrangeArray    bool
	ArrangeVertical bool
	ArrayPrefix     string
	ArraySuffix     string
	CellFmt         string
	ColSep          string
	DisplayWidth    int
	LinePrefix      string
	LineSuffix      string
	LJustify        bool
	TermAdjust      bool
}

func DefaultOptions() OptsT {
	opts := OptsT{
		ArrangeArray:    false,
		ArrangeVertical: true,
		ArrayPrefix:     "",
		ArraySuffix:     "",
		CellFmt:         "",
		ColSep:          "  ",
		DisplayWidth:    80,
		LinePrefix:      "",
		LineSuffix:      "\n",
		LJustify:        true,
		TermAdjust:      false,
	}
	return opts
}

type KeyValuePairT struct {
	Field string
	Value interface{}
}

func SetOptions(pairs ...KeyValuePairT) OptsT {
	opts := DefaultOptions()
	for _, pair := range pairs {
		switch pair.Field {
		case "ArrangeArray":
			v, _ := pair.Value.(bool)
			opts.ArrangeArray = v
		case "ArrangeVertical":
			v, _ := pair.Value.(bool)
			opts.ArrangeVertical = v
		case "ArrayPrefix":
			v, _ := pair.Value.(string)
			opts.ArrayPrefix = v
		case "ArraySuffix":
			v, _ := pair.Value.(string)
			opts.ArraySuffix = v
		case "CellFmt":
			v, _ := pair.Value.(string)
			opts.CellFmt = v
		case "ColSep":
			v, _ := pair.Value.(string)
			opts.ColSep = v
		case "DisplayWidth":
			v, _ := pair.Value.(int)
			opts.DisplayWidth = v
		case "LJustify":
			v, _ := pair.Value.(bool)
			opts.LJustify = v
		case "TermAdjust":
			v, _ := pair.Value.(bool)
			opts.TermAdjust = v
		}
	}
	return opts
}

// Return the length of string cell
func CellSize(cell string) int {
	return len(cell)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// The following routines ToStringArrayFromIndexable and ToStringArray are
// from Carlos Castillo. Thanks Carlos!
// http://play.golang.org/p/bxdcIj6ueH

/*
   ToStringSliceFromIndexable(slice_or_array, [format_string]) => [] string

Run fmt.Sprint on each of the elemnts in slice_or_array. If
format_string is given, that is the format string passed to fmt.Sprintf.

This routine assumes slice_or_array is a value which has a length,
and can be indexed (a slice/array). No checking on or error is thrown
if this is not the case.
*/
func ToStringSliceFromIndexable(x interface{}, optFmt ...string) []string {
	v := reflect.ValueOf(x)
	out := make([]string, v.Len())
	for i := range out {
		if 0 == len(optFmt) {
			out[i] = fmt.Sprint(v.Index(i).Interface())
		} else {
			out[i] = fmt.Sprintf(optFmt[0], v.Index(i).Interface())
		}
	}
	return out
}

/*
ToStringSlice(data, [format_string]) => [] string

If data is a slice or array, runs
ToStringSliceFromIndexable. Otherwise, data is put into a slice and
ToStringSliceFromIndexable is called on that slice of one element.

*/
func ToStringSlice(x interface{}, optFmt ...string) []string {
	v := reflect.ValueOf(x)
	if v.Kind() != reflect.Array && v.Kind() != reflect.Slice {
		// Not an array or slice, so run fmt.Sprint and turn that
		// single item into a slice.
		if 0 == len(optFmt) {
			return []string{fmt.Sprint(x)}
		} else {
			return []string{fmt.Sprintf(optFmt[0], x)}
		}
	}
	if 0 == len(optFmt) {
		return ToStringSliceFromIndexable(x)
	}
	return ToStringSliceFromIndexable(x, optFmt[0])
}

/*
 Return a string from an array with embedded newlines formatted so
 that when printed the columns are aligned.

 For example:, for a line width of 4 characters (arranged vertically):

	a = [] int{1,2,4,4}
	columnize.Columnize(a) => '1  3\n2  4\n'

Arranged horizontally:

	opts := columnize.Default_options()
	opts.arrange_vertical = false
	columnize.Columnize(a) =>  '1  2\n3  4\n'

Each column is only as wide as necessary.  By default, columns are
separated by two spaces
*/
func Columnize(list interface{}, opts OptsT) string {
	var l []string
	if opts.CellFmt != "" {
		l = ToStringSlice(list, opts.CellFmt)
	} else {
		l = ToStringSlice(list)
	}
	return ColumnizeS(l, opts)
}

// Like Columnize but we are already passed a slice of string
func ColumnizeS(list []string, opts OptsT) string {

	if opts.ArrangeArray {
		opts.ArrayPrefix = "["
		opts.ArraySuffix = "]\n"
		opts.LinePrefix = " "
		opts.LineSuffix = ",\n"
		opts.ColSep = ", "
		opts.ArrangeVertical = false
	}

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
	var colwidths []int
	if opts.ArrangeVertical {
		arrayIndex := func(numRows, row, col int) int {
			return numRows*col + row
		}
		// Try every row count from 1 upwards
		for nrows = 1; nrows < len(list); nrows++ {
			ncols = (len(list) + nrows - 1) / nrows
			colwidths = make([]int, 0)
			totwidth := -len(opts.ColSep)

			for col := 0; col < ncols; col++ {
				// get max column width for this column
				colwidth := 0
				for row := 0; row < nrows; row++ {
					i := arrayIndex(nrows, row, col)
					if i >= len(list) {
						break
					}
					colwidth = max(CellSize(list[i]), colwidth)
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
		// each column has been obtained.  Now we just have to format
		// each of the rows.
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
				colwidths = []int{}
				totalWidth = -len(opts.ColSep)
				var colwidth, row int
				for col := 0; col < ncols; col++ {
					// get max column width for this column
					for row = 1; row <= nrows; row++ {
						i = arrayIndex(ncols, row, col)
						if i >= len(list) {
							break
						}
						colwidth = max(colwidth, CellSize(list[i]))
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
