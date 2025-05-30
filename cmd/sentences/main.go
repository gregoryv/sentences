package main

import (
	"bufio"
	"flag"
	"log"
	"os"

	"github.com/gregoryv/sentences"
)

const usage = `Usage: sentences < file.txt

Parses plain text from stdin and outputs sentences one by
one on stdout.

Indentation and non recognized sentences are removed.
The primary goal of this tool is to parse sentences from
RFC like documents.`

func main() {
	flag.Usage = func() {
		log.SetFlags(0)
		log.Println(usage)
	}
	flag.Parse()
	cmd := sentences.Command{
		In:  bufio.NewReader(os.Stdin),
		Out: os.Stdout,
	}
	cmd.Run()
}
