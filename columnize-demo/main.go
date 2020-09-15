package main

import (
	"columnize"
	"fmt"
)

func main() {
	opts := columnize.DefaultOptions()

	list := make([]string, 0)
	fmt.Println(columnize.Format(list, opts))

	data := []string{
		"one", "two", "three",
		"for", "five", "six",
		"seven", "eight", "nine",
		"ten", "eleven", "twelve",
		"thirteen", "fourteen", "fifteen",
		"sixteen", "seventeen", "eightteen",
		"nineteen", "twenty", "twentyone",
		"twentytwo", "twentythree", "twentyfour",
		"twentyfive", "twentysix", "twentyseven",
	}

	opts.ArrangeVertical = true
	fmt.Println(columnize.FormatStringList(data, opts))
	opts.DisplayWidth = 50
	opts.LJustify = false
	fmt.Println(columnize.FormatStringList(data, opts))

	opts.ArrangeVertical = true
	fmt.Println(columnize.FormatStringList(data, opts))
	opts.DisplayWidth = 80
	opts.LJustify = true
	fmt.Println(columnize.FormatStringList(data, opts))
	fmt.Println("----------------")

	a := []int{31, 4, 1, 59, 2, 6, 5, 3}
	opts.ArrangeVertical = false
	opts.LJustify = false

	opts.DisplayWidth = 8
	fmt.Println(columnize.Format(a, opts))
	fmt.Println("----------------")

	fmt.Println(columnize.Format(a, columnize.ArrayOptions()))
}
