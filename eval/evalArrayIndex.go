package eval

import (
	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/object"
)

func evalArrayIndex(arr *ast.ArrayIndex, env *object.Environment) object.Object {

	indexObj := Eval(arr.Index, env)
	if object.IsError(indexObj) {
		return indexObj
	}

	obj := Eval(arr.Array, env)
	if object.IsError(obj) {
		return obj
	}

	// we're either dealing with arrays or indexes. Let's check arrays first
	if arrayObj, ok := obj.(*object.Array); ok {
		intIndex, ok := indexObj.(*object.Integer)
		if !ok {
			return object.ErrorExpected(object.INTEGER_OBJ)
		}
		n := int(intIndex.Value)
		if n >= len(arrayObj.Elems) {
			return object.Errorf("List index out of range: %d > %d", n, len(arrayObj.Elems))
		}
		if n < 0 {
			return object.Errorf("negative indices not allowed")
		}
		return arrayObj.Elems[n]
	}

	// Is it a hash?
	if hashObj, ok := obj.(*object.Hash); ok {
		pair, ok := hashObj.Pairs[computeStringHash(indexObj)]
		if !ok {
			return object.NULL
		}
		return pair.Value
	}

	return object.Errorf("indexing is only supported for arrays or hashes")
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
