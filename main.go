package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

)

var DefaultOutput io.Writer = os.Stdout

func main() {

	var useIndex bool
	flag.BoolVar(&useIndex, "index", false, "")
	flag.BoolVar(&useIndex, "i", false, "")

	flag.Parse()

	files := flag.Args()
	log.SetFlags(0)

	// ----------------------------------------

	for _, filename := range files {

		fh, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}

		scanner := bufio.NewScanner(fh)
		// Set the split function for the scanning operation.
		scanner.Split(ScanSentences)
		w := DefaultOutput
		var i int
		for scanner.Scan() {
			i++
			line := scanner.Text()
			line = strings.ReplaceAll(line, "\n", " ")
			line = strings.TrimSpace(line)
			if len(line) > 1 { // one character followed by ., ? or !
				if useIndex {
					fmt.Fprintln(w, i, line)
				} else {
					fmt.Fprintln(w, line)
				}
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		fh.Close()
	}
}

// ScanSentence is a split function for a Scanner that returns sentence.
func ScanSentences(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip until first capital letter is found
	start := 0

	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if unicode.IsUpper(r) {
			break
		}
	}
	// Scan until ., marking end of word.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		switch r {
		case '.', '?', '!':
			return i + width, data[start : i+width], nil
			break
		}
	}
	// Request more data.
	return start, nil, nil
}
