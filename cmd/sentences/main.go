package sentences

import (
	"bufio"
	"log"
	"os"

	"github.com/gregoryv/sentences"
)

func main() {
	log.SetFlags(0)
	cmd := sentences.Command{
		In:  bufio.NewReader(os.Stdin),
		Out: os.Stdout,
		Err: os.Stderr,
	}
	cmd.Run()
}
