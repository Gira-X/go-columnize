package columnize

import "testing"

func checkColumnize(expect string, data interface{}, opts Options, t *testing.T) {
	got := Format(data, opts)
	if expect != got {
		t.Errorf("got:\n%s\nwant:\n%s\n", got, expect)
	}
}

func verifyEqualityOfFormatFunctions(data []string, opts Options, t *testing.T) {
	formatOutput := Format(data, opts)
	formatStringOutput := FormatStringList(data, opts)
	if formatOutput != formatStringOutput {
		t.Errorf("Format() and FormatStringList() don't return an equal result:\n" +
			"Format():\n%v\n" +
			"FormatStringList():\n%v", formatOutput, formatStringOutput)
	}
}

func checkIntArrays(opts Options, t *testing.T) {
	testData := []int{1, 2, 3, 4}
	checkColumnize("1  3\n2  4\n", testData, opts, t)

	opts.ArrangeVertical = false
	checkColumnize("1  2\n3  4\n", testData, opts, t)

	opts.ArrangeArray = true
	checkColumnize("[1, 2,\n 3, 4,\n ]\n", testData, opts, t)

	opts.DisplayWidth = 8
	opts.CellFmt = "%02d"
	checkColumnize("[01, 02,\n 03, 04,\n ]\n", testData, opts, t)
}

func TestColumnize(t *testing.T) {
	opts := DefaultOptions()

	data := []string{"1", "2", "3"}
	opts.DisplayWidth = 10
	opts.ColSep = ", "
	checkColumnize("1, 2, 3\n", data, opts, t)
	verifyEqualityOfFormatFunctions(data, opts, t)

	data = []string{"1", "2", "3", "4"}
	opts.ColSep = "  "
	opts.DisplayWidth = 4
	checkColumnize("1  3\n2  4\n", data, opts, t)
	verifyEqualityOfFormatFunctions(data, opts, t)

	checkIntArrays(opts, t)

}
