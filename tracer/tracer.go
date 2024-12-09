package tracer

import (
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/kvalv/monkey/ast"
)

type Tracer struct {
	log.Logger
	level int
}

func New(w io.Writer) *Tracer {
	return &Tracer{
		Logger: *log.New(w, "", 0),
	}
}

func (t *Tracer) Trace(name string) func(n ast.Node) {
	start := time.Now()
	indent := strings.Repeat(" ", t.level*2)
	t.level++
	t.Logger.Printf("%sBEGIN %s", indent, name)
	return func(n ast.Node) {
		elapsed := time.Since(start)
		line := fmt.Sprintf("%sEND %s [%s]", indent, name, elapsed.String())
		if n != nil {
			line = fmt.Sprintf("%s -> %s", line, n.String())
		}
		t.Logger.Printf(line)
		t.level--
	}
}
