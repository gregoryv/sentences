package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_main(t *testing.T) {
	var buf bytes.Buffer
	DefaultOutput = &buf

	dir := t.TempDir()
	filename := filepath.Join(dir, "text")
	err := os.WriteFile(filename, []byte(`
# some header

This is a sentence. This is another
sentence. And another one.

- not a good one.

What do you mean?

you are wrong!
No, but you are!

`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	os.Args = []string{"test", "-i", filename}
	main()

	result := strings.TrimSpace(buf.String())
	got := strings.Split(result, "\n")
	if len(got) != 5 {
		for _, s := range got {
			t.Logf("%q", s)
		}
		t.Fail()
	}
}
