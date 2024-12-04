package main

import (
	"os"

	"github.com/kvalv/monkey/repl"
)

func main() {
	repl.Start(os.Stdout, os.Stdin)
}
