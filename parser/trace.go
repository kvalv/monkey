package parser

import (
	"io"
	"log"
	"strings"

	"github.com/kvalv/monkey/ast"
)

var (
	enabled bool = true // whether we enable logging
	level   int
)

func trace(name string) func(n ast.Node) {
	log.SetFlags(0) // doesn't really belong here but whatever
	if !enabled {
		log.SetOutput(io.Discard)
	}
	indent := strings.Repeat(" ", level*2)
	level++
	log.Printf("%sBEGIN %s", indent, name)
	return func(n ast.Node) {
		if n != nil {
			log.Printf("%sEND %s -> %s", indent, name, n.String())
		} else {
			log.Printf("%sEND %s", indent, name)
		}
		level--
	}
}
