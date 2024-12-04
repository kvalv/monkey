package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/kvalv/monkey/parser"
)

func Start(w io.Writer, r io.Reader) {
	sc := bufio.NewScanner(r)
	fmt.Fprintf(w, "> ")
	for sc.Scan() {
		line := sc.Text()
		p := parser.New(line)
		prog, errs := p.Parse()
		if len(errs) > 0 {
			for _, err := range errs {
				fmt.Fprintf(w, "ERROR: %v", err)
			}
			continue
		}
		fmt.Fprintf(w, "\n%s", prog)
		fmt.Fprintf(w, "\n> ")
	}
}
