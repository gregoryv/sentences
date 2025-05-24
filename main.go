package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
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
		line := scanner.Bytes()
		if i := bytes.LastIndex(line, doubleNL); i > -1 {
			// found an empty line, this is normal after headings
			line = line[i+2:]
		}
		line = bytes.ReplaceAll(line, oneNL, oneSpace)
		line = bytes.TrimSpace(line)
		if len(line) > 1 { // one character followed by ., ? or !
			w.Write(line)
			w.Write(oneNL)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

var (
	nl       = byte('\n')
	oneNL    = []byte{nl}
	oneSpace = []byte{' '}
	doubleNL = []byte{nl, nl}
)

// ScanSentence is a split function for a Scanner that returns
// sentence.
func ScanSentences(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip until first capital letter is found
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if unicode.IsUpper(r) {
			break
		}
	} // capital letter found

	// find what looks like end of sentence.
	var width int
	for i := start; i < len(data); i += width {
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
