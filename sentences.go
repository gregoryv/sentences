package sentences

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
)

type Command struct {
	In  io.Reader
	Out io.Writer
}

func (c *Command) Run() {
	next := capitalLetter
	for next != nil {
		next = next(c.Out, bufio.NewReader(c.In))
	}
}

var (
	buf = &bytes.Buffer{}
	p   = make([]byte, 1)
)

func capitalLetter(w io.Writer, r *bufio.Reader) parseFn {
	for {
		r, _, err := r.ReadRune()
		if err != nil {
			return nil
		}

		if unicode.IsUpper(r) {
			buf.WriteRune(r)
			return endOfSentence
		}
	}
}

func endOfSentence(w io.Writer, r *bufio.Reader) parseFn {
	var lastNewline bool
	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil
		}

		switch b {
		case '(':
			buf.WriteByte(b)
			return endOf(')')

		case '.':
			buf.WriteByte(b)
			return spaceAfterDot

		case '?', '!':
			buf.WriteByte(b)
			buf.WriteString("\n")
			io.Copy(w, buf)
			return capitalLetter

		case '\n':
			scanEmpty(r)
			if lastNewline {
				// no delimiter but two new lines would mean the
				// sentence is wrongly formatted or e.g. there is
				// and indented example or something like that.
				// we skip it nontheless
				buf.Truncate(0)
			} else {
				lastNewline = true
				buf.WriteByte(' ')
			}

		default:
			lastNewline = false
			buf.WriteByte(b)
		}
	}
}

func endOf(endChar byte) parseFn {
	return func(w io.Writer, r *bufio.Reader) parseFn {
		for {
			b, err := r.ReadByte()
			if err != nil {
				return nil
			}
			switch b {
			case endChar:
				buf.WriteByte(b)
				return endOfSentence

			case '\n':
				buf.WriteByte(' ')
				scanEmpty(r)
			default:
				buf.WriteByte(b)
			}
		}
	}
}

func scanEmpty(r *bufio.Reader) {
	for {
		b, err := r.ReadByte()
		if b == ' ' {
			continue
		}
		if err != nil {
			return
		}
		r.UnreadByte()
		return
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
			return endOfSentence
		}
	}
}

type parseFn func(w io.Writer, r *bufio.Reader) parseFn
