package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func bytesToKB(bytes int64) float64 {
	const KB = 1024
	return float64(bytes) / float64(KB)
}

func printTabularData(header string, dataFormat string, data ...interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 10, 1, 5, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, header)

	fs := dataFormat + "\n"
	fmt.Fprintf(w, fs, data...)
}
