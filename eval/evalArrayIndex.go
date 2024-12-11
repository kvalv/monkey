package eval

import (
	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/object"
)

func evalArrayIndex(arr *ast.ArrayIndex, env *object.Environment) object.Object {
	arrayObj, err := evalTo[*object.Array](arr.Array, env)
	if err != nil {
		return err
	}

	indexObj, err := evalTo[*object.Integer](arr.Index, env)
	if err != nil {
		return err
	}

	n := int(indexObj.Value)
	if n < 0 {
		return object.Errorf("negative indices not allowed")
	}

	if n >= len(arrayObj.Elems) {
		return object.Errorf("List index out of range: %d > %d", n, len(arrayObj.Elems))
	}
	return arrayObj.Elems[n]
}

func evalTo[T object.Object](node ast.Expression, env *object.Environment) (T, *object.Error) {
	var empty T // this is actually nil

	got := Eval(node, env)
	if object.IsError(got) {
		return empty, got.(*object.Error)
	}

	parsed, ok := got.(T)
	if !ok {
		e := object.Errorf("object type mismatch")
		return empty, e
	}

	return parsed, nil
}
