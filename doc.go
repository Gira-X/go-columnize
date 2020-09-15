/*

Package columnize formats a slice or array into a single string with embedded
newlines.
On printing the string, the columns are aligned.

Summary

Return a string from an array with embedded newlines formatted
so that when printed the columns are aligned.

	data := []int{1, 2, 3, 4}
	opts := columnize.DefaultOptions()
	opts.DisplayWidth = 10
	fmt.Println(columnize.Format(data, opts))
	# prints "1  3\n2  4\n"


Options:
	ArrangeArray bool	format string as a go array
	ArrangeVertical bool	format entries top down and left to right
				Otherwise left to right and top to bottom
	ArrayPrefix string	Add this to the beginning of the string
	ArraySuffix string	Add this to the end of the string
	CellFmt			A format specify for formatting each item each array
				item to a string
	ColSep string		Add this string between columns
	DisplayWidth int	the maximum line display width
	LinePrefix string	Add this prefix for each line
	LineSuffix string	Add this suffix for each line
	LJustify bool		whether to left justify text instead of right justify

Examples

	data := []int{1, 2, 3, 4}
	opts := columnize.DefaultOptions()
	fmt.Println(columnize.Format(data, opts))

	opts.DisplayWidth = 8
	opts.CellFmt = "%02d"
	fmt.Println(columnize.Format(data, opts))

	opts.ArrangeArray = true
	opts.CellFmt = ""
	fmt.Println(columnize.Format(data, opts))

Author

	Rocky Bernstein	<rocky@gnu.org>

	Also available in Python (columnize), and Perl
	(Array::Format) and Ruby (columnize)

	Copyright 2013 Rocky Bernstein.

*/
package columnize
