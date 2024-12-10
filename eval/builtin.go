package eval

import "github.com/kvalv/monkey/object"

var builtin map[string]*object.Builtin = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.Errorf("len() accepts 1 argument, got %d", len(args))
			}
			value := args[0]
			s, ok := value.(*object.String)
			if !ok {
				return object.Errorf("type error: expected %s but got %s", object.STRING_OBJ, value.Type())
			}
			return &object.Integer{Value: int64(len(s.Value))}
		},
	},
}
