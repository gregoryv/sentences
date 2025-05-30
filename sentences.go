package main

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
)

type Command struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

func (c *Command) Run() {
	next := capitalLetter
	for next != nil {
		next = next(c.Out, bufio.NewReader(c.In))
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
			buf.WriteByte(b)
			return end
		}
	}
}

type parseFn func(w io.Writer, r *bufio.Reader) parseFn
