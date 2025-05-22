package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gregoryv/golden"
)

func Test_sentences(t *testing.T) {
	var buf bytes.Buffer

	r := strings.NewReader(`
# Some header

This is a sentence. This is another
sentence. And another one.

- not a good one.

What do you mean?

you are wrong!
No, but you are!

As seen in this example

   1 2 3

Sentence starting after a newline.

`)

	sentences(&buf, r)

	result := strings.TrimSpace(buf.String())
	golden.AssertWith(t, result, "testdata/found.txt")
}
