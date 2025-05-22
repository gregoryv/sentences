package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

func main() {
	log.SetFlags(0)
	sentences(os.Stdout, os.Stdin)
}

func sentences(w io.Writer, r io.Reader) {
	scanner := bufio.NewScanner(r)
	// Set the split function for the scanning operation.
	scanner.Split(ScanSentences)

	var i int
	for scanner.Scan() {
		i++
		line := scanner.Text()
		if i := strings.LastIndex(line, "\n\n"); i > -1 {
			// found an empty line, this is normal after headings
			line = line[i+2:]
		}
		line = strings.ReplaceAll(line, "\n", " ")
		line = strings.TrimSpace(line)
		if len(line) > 1 { // one character followed by ., ? or !
			fmt.Fprintln(w, line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
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
