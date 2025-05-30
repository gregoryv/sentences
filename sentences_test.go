package sentences

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

   This one is indented
   and spread over many
   lines.

Using parenthesis (e.g.
HTTP 1.1. ).

Using square [HTTP
1.1. ].

Using curly { HTTP
1.1. } .

   Clarified which error code should be used for inbound server failures
   (e.g. DNS failures). (Section 10.5.5).

Incomplete sentence
`)

	var buf bytes.Buffer
	cmd := Command{
		Out: &buf,
		In:  bufio.NewReader(bytes.NewReader(data)),
	}
	cmd.Run()

	result := strings.TrimSpace(buf.String())
	golden.AssertWith(t, result, "testdata/found.txt")

	cases := []string{
		"One sentence.",
		"no sentences",
		"Incomplete (HTTP 1.1",
	}
	for _, v := range cases {
		t.Run(v, func(t *testing.T) {
			cmd := Command{
				Out: ioutil.Discard,
				In:  bufio.NewReader(strings.NewReader(v)),
			}
			cmd.Run()
		})
	}
}

func Benchmark(b *testing.B) {
	data, _ := os.ReadFile("testdata/rfc2616.txt")
	rdata := bytes.NewReader(data)
	r := bufio.NewReader(rdata)
	cmd := Command{
		In:  r,
		Out: ioutil.Discard,
	}
	for b.Loop() {
		cmd.Run()
		rdata.Reset(data)
		r.Reset(rdata)
	}
}
