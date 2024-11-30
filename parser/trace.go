package parser

import (
	"log"
	"strings"

	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/token"
)

var level int

var verbose bool = true

func trace(name string, tk token.Token) func(n ast.Node) {
	log.SetFlags(0) // doesn't really belong here but whatever
	indent := strings.Repeat(" ", level*2)
	level++
	if verbose {
		log.Printf("%sBEGIN %s %s", indent, name, tk.Literal)
	} else {
		log.Printf("%sBEGIN %s %s", indent, name, "")
	}
	return func(n ast.Node) {
		if n != nil {
			log.Printf("%sEND %s -> %s", indent, name, n.String())
		} else {
			log.Printf("%sEND %s", indent, name)
		}
		level--
	}
}
