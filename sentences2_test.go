package main

import (
	"bytes"
	"io"
	"unicode"
)

func sentences2(w io.Writer, r io.Reader) {
	next := capitalLetter
	for next != nil {
		next = next(w, r)
	}
}

var buf = &bytes.Buffer{}
var p = make([]byte, 1)

func capitalLetter(w io.Writer, r io.Reader) parseFn {
	for {
		_, err := r.Read(p)
		if err != nil {
			return nil
		}

		if unicode.IsUpper(rune(p[0])) {
			buf.Write(p)
			return end
		}
	}
	return nil
}

func end(w io.Writer, r io.Reader) parseFn {
	var lastNewline bool
	for {
		_, err := r.Read(p)
		if err != nil {
			return nil
		}

		switch p[0] {
		case '.':
			buf.Write(p)
			return space

		case '?', '!':
			buf.Write(p)
			buf.WriteString("\n")
			io.Copy(w, buf)
			return capitalLetter

		case '\n':
			if lastNewline {
				buf.Truncate(0)
			} else {
				lastNewline = true
				buf.WriteString(" ")
			}

		default:
			buf.Write(p)
		}
	}
	return nil
}

func space(w io.Writer, r io.Reader) parseFn {
	for {
		_, err := r.Read(p)
		if err != nil {
			return nil
		}

		switch p[0] {
		case ' ', '\n', '\t':
			buf.WriteString("\n")
			io.Copy(w, buf)
			return capitalLetter

		default:
			return end
		}
	}
}

type parseFn func(w io.Writer, r io.Reader) parseFn
