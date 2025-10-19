package sentences

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/gregoryv/golden"
)

func Test(t *testing.T) {
	data := `
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

Över landskap.
I am.

      Notes:

      1) Additional CRLFs may precede the first boundary string in the
         entity.


Incomplete sentence
`

	var buf bytes.Buffer
	ParseString(&buf, data)

	result := strings.TrimSpace(buf.String())
	golden.AssertWith(t, result, "testdata/found.txt")

	cases := []string{
		"One sentence.",
		"no sentences",
		"Incomplete (HTTP 1.1",
		"I.\n",
		"I_am.\n",
	}
	for _, v := range cases {
		t.Run(v, func(t *testing.T) {
			ParseString(io.Discard, v)
		})
	}
}

func Benchmark(b *testing.B) {
	data, _ := os.ReadFile("testdata/rfc2616.txt")
	rdata := bytes.NewReader(data)
	r := bufio.NewReader(rdata)
	cmd := Command{
		In:  r,
		Out: io.Discard,
	}
	for b.Loop() {
		cmd.Run()
		rdata.Reset(data)
		r.Reset(rdata)
	}
}
