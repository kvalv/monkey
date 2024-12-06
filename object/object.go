package object

import (
	"fmt"
)

type Type string

const (
	INTEGER_OBJ = "integer"
	BOOLEAN_OBJ = "boolean"
	NULL_OBJ    = "null"
	RETURN_OBJ  = "return"
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

type (
	Integer struct{ Value int64 }
	Boolean struct{ Value bool }
	Null    struct{}
	Return  struct{ Object }
)

func (i *Integer) Type() Type     { return INTEGER_OBJ }
func (i *Integer) String() string { return fmt.Sprintf("%d", i.Value) }

func (b *Boolean) Type() Type     { return BOOLEAN_OBJ }
func (b *Boolean) String() string { return fmt.Sprintf("%t", b.Value) }

func (b *Null) Type() Type     { return NULL_OBJ }
func (b *Null) String() string { return NULL_OBJ }

func (r *Return) Type() Type     { return RETURN_OBJ }
func (r *Return) String() string { return r.Object.String() }
