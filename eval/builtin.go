package eval

import "github.com/kvalv/monkey/object"

var builtin map[string]*object.Builtin = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.Errorf("len() accepts 1 argument, got %d", len(args))
			}
			switch obj := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(obj.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(obj.Elems))}
			default:
				return object.Errorf("len() not supported for objects of type %s", obj.Type())
			}
		},
	},
	"first": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.Errorf("first() accepts 1 argument, got %d", len(args))
			}
			switch obj := args[0].(type) {
			case *object.Array:
				if len(obj.Elems) == 0 {
					return object.NULL
				}
				return obj.Elems[0]
			default:
				return object.Errorf("len() not supported for objects of type %s", obj.Type())
			}
		},
	},
	"last": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.Errorf("last() accepts 1 argument, got %d", len(args))
			}
			switch obj := args[0].(type) {
			case *object.Array:
				if n := len(obj.Elems); n == 0 {
					return object.NULL
				} else {
					return obj.Elems[n-1]
				}
			default:
				return object.Errorf("len() not supported for objects of type %s", obj.Type())
			}
		},
	},
	"rest": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.Errorf("rest() accepts 1 argument, got %d", len(args))
			}
			switch obj := args[0].(type) {
			case *object.Array:
				if n := len(obj.Elems); n == 0 {
					return &object.Array{}
				} else {
					res := &object.Array{Elems: make([]object.Object, n-1, n-1)}
					copy(res.Elems, obj.Elems[1:])
					return res
				}
			default:
				return object.Errorf("len() not supported for objects of type %s", obj.Type())
			}
		},
	},
	"push": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) <= 0 {
				return object.Errorf("push() array missing")
			}
			switch obj := args[0].(type) {
			case *object.Array:
				old, ok := args[0].(*object.Array)
				if !ok {
					return object.Errorf("push(): first argument is not an array")
				}
				n := len(args) - 1
				m := len(old.Elems) + n
				res := &object.Array{Elems: make([]object.Object, m, m)}
				copy(res.Elems, old.Elems)
				for i, obj := range args[1:] {
					res.Elems[n+i] = obj
				}
				return res
			default:
				return object.Errorf("len() not supported for objects of type %s", obj.Type())
			}
		},
	},
}
