package eval

import (
	"crypto/md5"
	"fmt"

	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/object"
)

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	hash := &object.Hash{
		Pairs: make(map[string]object.Pair),
	}
	for k, v := range node.Pairs {
		key, value := Eval(k, env), Eval(v, env)
		if object.IsError(key) {
			return key
		}
		if object.IsError(value) {
			return value
		}
		// for now we'll just naively use md5sum on the string representation. what could go wrong
		hkey := computeStringHash(key)
		hash.Pairs[hkey] = object.Pair{Key: key, Value: value}
	}
	return hash
}

func computeStringHash(obj object.Object) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(obj.String())))
}
