package object

import (
	"fmt"
	"strings"

	"github.com/kvalv/monkey/ast"
)

type Type string

const (
	INTEGER_OBJ  = "INTEGER"
	BOOLEAN_OBJ  = "BOOLEAN"
	NULL_OBJ     = "NULL"
	RETURN_OBJ   = "RETURN"
	ERROR_OBJ    = "ERROR"
	FUNCTION_OBJ = "FUNCTION"
	STRING_OBJ   = "STRING"
	BUILTIN_OBJ  = "BUILTIN"
	ARRAY_OBJ    = "ARRAY"
)

var (
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
	NULL  = &Null{}
)

type Object interface {
	Type() Type
	String() string
}

type BuiltinFunction func(args ...Object) Object

type (
	Integer  struct{ Value int64 }
	Boolean  struct{ Value bool }
	Null     struct{}
	Return   struct{ Object }
	Error    struct{ Message string }
	Function struct {
		Env    *Environment
		Params []ast.Identifier
		Body   *ast.BlockStatement
	}
	String  struct{ Value string }
	Builtin struct{ Fn BuiltinFunction }
	Array   struct{ Elems []Object }
)

func (i *Integer) Type() Type     { return INTEGER_OBJ }
func (i *Integer) String() string { return fmt.Sprintf("%d", i.Value) }

func (b *Boolean) Type() Type     { return BOOLEAN_OBJ }
func (b *Boolean) String() string { return fmt.Sprintf("%t", b.Value) }

func (b *Null) Type() Type     { return NULL_OBJ }
func (b *Null) String() string { return NULL_OBJ }

func (r *Return) Type() Type     { return RETURN_OBJ }
func (r *Return) String() string { return fmt.Sprintf("return %s", r.Object.String()) }

func (e *Error) Type() Type                 { return ERROR_OBJ }
func (e *Error) String() string             { return fmt.Sprintf("error: %s", e.Message) }
func Errorf(format string, a ...any) *Error { return &Error{Message: fmt.Sprintf(format, a...)} }
func IsError(o Object) bool                 { return o.Type() == ERROR_OBJ }

func (f *Function) Type() Type { return FUNCTION_OBJ }
func (f *Function) String() string {
	var params []string
	for _, p := range f.Params {
		params = append(params, p.String())
	}
	indent2 := func(s string) string { return strings.Replace(s, "\n", "\n  ", -1) }
	return fmt.Sprintf("fn(%s) {\n%s\n}",
		strings.Join(params, ", "),
		indent2(f.Body.String()),
	)
}

func (s *String) Type() Type { return STRING_OBJ }
func (s *String) String() string {
	return s.Value
}

func (b *Builtin) Type() Type { return BUILTIN_OBJ }
func (b *Builtin) String() string {
	return "builtin function"
}

func (a *Array) Type() Type { return ARRAY_OBJ }
func (a *Array) String() string {
	var elems []string
	for _, e := range a.Elems {
		elems = append(elems, e.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(elems, ", "))
}
