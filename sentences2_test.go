package main

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
)

func sentences2(w io.Writer, r *bufio.Reader) {
	next := capitalLetter
	for next != nil {
		next = next(w, r)
	}
}

var buf = &bytes.Buffer{}
var p = make([]byte, 1)

func capitalLetter(w io.Writer, r *bufio.Reader) parseFn {
	for {
		r, _, err := r.ReadRune()
		if err != nil {
			return nil
		}

		if unicode.IsUpper(r) {
			buf.WriteRune(r)
			return end
		}
	}
	return nil
}

func end(w io.Writer, r *bufio.Reader) parseFn {
	var lastNewline bool
	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil
		}

		switch b {
		case '.':
			buf.WriteByte(b)
			return spaceAfterDot

		case '?', '!':
			buf.WriteByte(b)
			buf.WriteString("\n")
			io.Copy(w, buf)
			return capitalLetter

		case '\n':
			if lastNewline {
				// no end delimiter but two new lines would mean
				buf.Truncate(0)
			} else {
				lastNewline = true
				buf.WriteByte(' ')
			}

		default:
			buf.WriteByte(b)
		}
	}
	return nil
}

func spaceAfterDot(w io.Writer, r *bufio.Reader) parseFn {
	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil
		}

		switch b {
		case ' ', '\n', '\t':
			buf.WriteByte('\n')
			io.Copy(w, buf)
			return capitalLetter

		default:
			return end
		}
	}
}

type parseFn func(w io.Writer, r *bufio.Reader) parseFn
