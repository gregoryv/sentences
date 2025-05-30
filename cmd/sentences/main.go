package sentences

import (
	"bufio"
	"os"

	"github.com/gregoryv/sentences"
)

func main() {
	cmd := sentences.Command{
		In:  bufio.NewReader(os.Stdin),
		Out: os.Stdout,
	}
	cmd.Run()
}
