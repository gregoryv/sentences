package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/gregoryv/golden"
)

func Test(t *testing.T) {
	data := []byte(`
# Some header

This is a sentence. This is another
sentence. And another one.

- not a good one.

What do you mean?

you are wrong!
No, but you are!

Software v1.1 is great!

As seen in this example

   1 2 3

Sentence starting after a newline.

Requirement links SHOULD(#R7) be written within parenthesis and start with
'#R' followed by a number.



`)

	t.Run("", func(t *testing.T) {
		var buf bytes.Buffer
		r := bytes.NewReader(data)
		sentences(&buf, r)

		result := strings.TrimSpace(buf.String())
		golden.AssertWith(t, result, "testdata/found.txt")
	})

	t.Run("", func(t *testing.T) {
		var buf bytes.Buffer
		r := bufio.NewReader(bytes.NewReader(data))
		sentences2(&buf, r)

		result := strings.TrimSpace(buf.String())
		golden.AssertWith(t, result, "testdata/found2.txt")
	})
}

func Benchmark(b *testing.B) {
	data, _ := os.ReadFile("testdata/rfc2616.txt")

	b.Run("reader", func(b *testing.B) {
		r := bytes.NewReader(data)
		p := make([]byte, 1024)
		for b.Loop() {
		inner:
			for {
				_, err := r.Read(p)
				if err != nil {
					break inner
				}
			}
			r.Reset(data)
		}
	})

	b.Run("scanner", func(b *testing.B) {
		r := bytes.NewReader(data)
		for b.Loop() {
			s := bufio.NewScanner(r)
			for s.Scan() {
				_ = s.Bytes()
			}
			r.Reset(data)
		}
	})

	b.Run("", func(b *testing.B) {
		r := bytes.NewReader(data)
		for b.Loop() {
			sentences(ioutil.Discard, r)
			r.Reset(data)
		}
	})

	b.Run("", func(b *testing.B) {
		rdata := bytes.NewReader(data)
		r := bufio.NewReader(rdata)
		for b.Loop() {
			sentences2(ioutil.Discard, r)
			rdata.Reset(data)
			r.Reset(rdata)
		}
	})

}
