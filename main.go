package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)
	sentences(os.Stdout, bufio.NewReader(os.Stdin))
}
