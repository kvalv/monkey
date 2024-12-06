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
	Output      io.Writer
	level       int
	initialized bool
}

func (t *Tracer) Trace(name string) func(n ast.Node) {
	if !t.initialized {
		t.init()
	}
	start := time.Now()
	indent := strings.Repeat(" ", t.level*2)
	t.level++
	log.Printf("%sBEGIN %s", indent, name)
	return func(n ast.Node) {
		elapsed := time.Since(start)
		line := fmt.Sprintf("%sEND %s [%s]", indent, name, elapsed.String())
		if n != nil {
			line = fmt.Sprintf("%s -> %s", line, n.String())
		}
		log.Printf(line)
		t.level--
	}
}

func (t *Tracer) init() {
	t.Logger.SetFlags(0) // doesn't really belong here but whatever
	if t.Output != nil {
		t.Logger.SetOutput(t.Output)
	}
	t.initialized = true
}
