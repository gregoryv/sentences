package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)
	c := Command{
		In:  bufio.NewReader(os.Stdin),
		Out: os.Stdout,
		Err: os.Stderr,
	}
	c.Run()
}
