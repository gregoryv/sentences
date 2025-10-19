package sentences

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode"
)

func ParseString(out io.Writer, v string) {
	Parse(out, strings.NewReader(v))
}

func Parse(out io.Writer, in io.Reader) {
	cmd := Command{
		In:  bufio.NewReader(in),
		Out: out,
	}
	cmd.Run()
}

type Command struct {
	In  *bufio.Reader
	Out io.Writer
}

func (c *Command) Run() {
	next := capitalLetter
	for next != nil {
		next = next(c.Out, c.In)
	}
}

var buf = &bytes.Buffer{}

func capitalLetter(w io.Writer, r *bufio.Reader) parseFn {
	for {
		r, width, err := r.ReadRune()
		if err != nil {
			return nil
		}

		if unicode.IsUpper(r) {
			buf.Truncate(0)
			buf.WriteRune(r)
			return endOfSentence
		}

		if width == 1 {
			b := string(r)[0]
			switch b {
			case '(', '[', '{':
				return endOf(opposite[b], capitalLetter)
			}
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
		case '(', '[', '{':
			buf.WriteByte(b)
			return endOf(opposite[b], endOfSentence)

		case '.':
			buf.WriteByte(b)
			return spaceAfterDot

		case '?', '!':
			buf.WriteByte(b)
			buf.WriteString("\n")
			writeSentence(w, buf)
			return capitalLetter

		case '\n':
			scanEmpty(r)
			if lastNewline {
				// no delimiter but two new lines would mean the
				// sentence is wrongly formatted or e.g. there is
				// and indented example or something like that.
				// we skip it nontheless
				return capitalLetter
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

func endOf(endChar byte, next parseFn) parseFn {
	return func(w io.Writer, r *bufio.Reader) parseFn {
		for {
			b, err := r.ReadByte()
			if err != nil {
				return nil
			}
			switch b {
			case endChar:
				buf.WriteByte(b)
				return next

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
			writeSentence(w, buf)
			return capitalLetter

		default:
			buf.WriteByte(b)
			return endOfSentence
		}
	}
}

func writeSentence(w io.Writer, buf *bytes.Buffer) {
	line := buf.Bytes()
	defer buf.Reset()
	if buf.Len() <= 5 { // there are no sentences this short
		return
	}

	if !bytes.Contains(line[:len(line)-1], []byte{' '}) {
		return
	}
	io.Copy(w, buf)
}

var opposite = map[byte]byte{
	'(': ')',
	'[': ']',
	'{': '}',
}

type parseFn func(w io.Writer, r *bufio.Reader) parseFn
