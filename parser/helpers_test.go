package parser_test

import (
	"testing"

	"github.com/kvalv/monkey/ast"
)

func expectLiteral(t *testing.T, got ast.Expression, value any) {
	switch v := value.(type) {
	case int:
		expectNumberLiteral(t, got, v)
	case string:
		expectIdentifierLiteral(t, got, v)
	default:
		t.Fatalf("unknown type %T", v)
	}
}

func expectNumberLiteral(t *testing.T, got ast.Expression, value int) {
	e, ok := got.(*ast.Number)
	if !ok {
		t.Fatalf("expected number - got %T", got)
	}
	if e.Value != value {
		t.Fatalf("value mismatch: expected %d but got %d", value, e.Value)
	}
}

func expectIdentifierLiteral(t *testing.T, got ast.Expression, value string) {
	e, ok := got.(*ast.Identifier)
	if !ok {
		t.Fatalf("expected identifier, got %T", got)
	}
	if e.Value != value {
		t.Fatalf("value mismatch: expected %s but got %s", value, e.Value)
	}
}
func expectInfixExpression(t *testing.T, got ast.Expression, lhs any, op string, rhs any) {
	exp, ok := got.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("not an infix expression: %T", got)
	}

	expectLiteral(t, exp.Lhs, lhs)
	if exp.Op != op {
		t.Fatalf("operator mismatch: got %v, want %v", exp.Op, op)
	}
	expectLiteral(t, exp.Rhs, rhs)
}
func expectLetStatement(t *testing.T, got ast.Statement, name string) {
	if got.TokenLiteral() != "let" {
		t.Fatalf("TokenLiteral mismatch: expected 'let' got %q", got.TokenLiteral())
	}

	stmt, ok := got.(*ast.LetStatement)
	if !ok {
		t.Fatal("not a let statement")
	}
	if stmt.Lhs.Value != name {
		t.Fatalf("value mismatch: expected %q, got %q", name, stmt.Lhs.Value)
	}
	// we won't check the rhs for now
}
