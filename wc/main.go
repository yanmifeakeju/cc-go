package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/yanmifeakeju/cc-go/wc/wc"
)

func main() {
	countBytes := flag.Bool("c", false, "count the number of bytes in a file")
	countLines := flag.Bool("l", false, "count the number of lines in a file")
	countWords := flag.Bool("w", false, "count the number of words in a file")
	countChars := flag.Bool("m", false, "count the number of characters")

	flag.Parse()

	// Characters is not in default
	if !*countBytes && !*countLines && !*countWords {
		*countBytes = true
		*countLines = true
		*countWords = true
	}

	opts := wc.WCopts{
		CountBytes: *countBytes,
		CountLines: *countLines,
		CountWords: *countWords,
		CountChars: *countChars,
	}

	if flag.NArg() == 0 {
		if err := processInput(os.Stdin, "stdin", opts); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing stdin: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Process each file provided as argument
	for _, filename := range flag.Args() {
		f, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening %s: %v\n", filename, err)
			continue
		}

		if err := processInput(f, filename, opts); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", filename, err)
		}
		f.Close()
	}
}

func processInput(r io.Reader, name string, opts wc.WCopts) error {
	c := wc.WC{Reader: r}
	counts, err := c.Count(opts)
	if err != nil {
		return err
	}

	// Build output string based on requested counts
	var output string
	if opts.CountLines {
		output += fmt.Sprintf("%d ", counts.Lines)
	}
	if opts.CountWords {
		output += fmt.Sprintf("%d ", counts.Words)
	}
	if opts.CountBytes {
		output += fmt.Sprintf("%d ", counts.Bytes)
	}

	if opts.CountChars {
		output += fmt.Sprintf("%d ", counts.Characters)
	}

	if name != "stdin" {
		output += name
	}

	fmt.Println(output)
	return nil
}
