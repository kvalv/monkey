package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/kvalv/monkey/eval"
	"github.com/kvalv/monkey/object"
	"github.com/kvalv/monkey/parser"
)

func Start(w io.Writer, r io.Reader) {
	sc := bufio.NewScanner(r)
	fmt.Fprintf(w, "> ")
	env := object.NewEnvironment()
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
		res := eval.Eval(prog, env)
		fmt.Fprintf(w, "\n%s", res)
		fmt.Fprintf(w, "\n> ")
	}
}
