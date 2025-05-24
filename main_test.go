package main

import (
	"bytes"
	"io/ioutil"
	"os"
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

func Benchmark(b *testing.B) {
	data, _ := os.ReadFile("testdata/rfc2616.txt")

	r := bytes.NewReader(data)
	for b.Loop() {
		sentences(ioutil.Discard, r)
		r.Reset(data)

	}
}
