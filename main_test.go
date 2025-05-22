package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gregoryv/golden"
)

func Test_main(t *testing.T) {
	var buf bytes.Buffer
	DefaultOutput = &buf

	dir := t.TempDir()
	filename := filepath.Join(dir, "text")
	err := os.WriteFile(filename, []byte(`
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

`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	os.Args = []string{"test", "-i", filename}
	main()

	result := strings.TrimSpace(buf.String())
	golden.AssertWith(t, result, "testdata/found.txt")
}
