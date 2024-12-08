package eval_test

import (
	"fmt"
	"testing"

	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/eval"
	"github.com/kvalv/monkey/object"
	"github.com/kvalv/monkey/parser"
)

func TestIntegerExpression(t *testing.T) {
	cases := []struct {
		input    string
		expected int64
	}{
		{"2", 2},
		{"3", 3},
		{"-3", -3},
		{"2 + 2", 4},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2* 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			prog := expectParse(t, tc.input)
			got := expectEval(t, prog)
			expectIntegerLiteral(t, got, tc.expected)
		})
	}
}

func TestBooleanExpression(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"!false", true},
		{"!true", false},
		{"!!true", true},
		{"1 > 2", false},
		{"3 > 2", true},
		{"true != false", true},
		{"false != false", false},
		{"(1 < 2) == false", false},
		{"(1 < 2) == true", true},
		{"(1 > 2) == false", true},
		{"(1 > 2) == true", false},
	}
	for _, tc := range cases {
		prog := expectParse(t, tc.input)
		got := expectEval(t, prog)
		expectBooleanLiteral(t, got, tc.expected)
	}
}

func TestIfExpression(t *testing.T) {
	cases := []struct {
		input    string
		expected any
	}{
		{"if (3 > 2) { 4 } else { 5 }", 4},
		{"if (3 < 2) { 4 } else { 5 }", 5},
		{"if (3 > 2) { true } else { false }", true},
		{"if (3 < 2) { true } else { false }", false},
		{"if (false) { 1 }", nil},
		{"if (true) { 1 }", 1},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			prog := expectParse(t, tc.input)
			got := expectEval(t, prog)
			expectLiteral(t, got, tc.expected)
		})
	}
}

func TestReturnStatement(t *testing.T) {
	cases := []struct {
		input    string
		expected any
	}{
		{"3; return 4; 5;", 4},
		{"3; return 4; return 5; 6", 4},
		{"3; 4; return 5;", 5},
		{`if true { if true { return 2; } return 1; }`, 2}}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			prog := expectParse(t, tc.input)
			got := expectEval(t, prog)
			expectLiteral(t, got, tc.expected)
		})
	}
}

func TestErrorHandling(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"5 + true", "type mismatch: INTEGER + BOOLEAN"},
		{"true > false", "unknown operator: BOOLEAN > BOOLEAN"},
		{"if true { if true { return 2 + false } }", "type mismatch: INTEGER + BOOLEAN"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			prog := expectParse(t, tc.input)
			got := expectEval(t, prog)
			expectErrorMessage(t, got, tc.expected)
		})
	}
}

func TestLetStatement(t *testing.T) {
	cases := []struct {
		input    string
		expected any
	}{
		{"let x = 5; x", 5},
		{"let a = 1; a + 1", 2},
		{"let a = b", fmt.Errorf(`identifier 'b' not defined`)},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			prog := expectParse(t, tc.input)
			got := expectEval(t, prog)
			expectLiteral(t, got, tc.expected)
		})
	}
}

func expectParse(t *testing.T, input string) *ast.Program {
	t.Helper()
	prog, errs := parser.New(input).Parse()
	if len(errs) > 0 {
		t.Error("failed to parse program:\n")
		for _, err := range errs {
			t.Errorf("\t%s\n", err)
		}
		t.FailNow()
	}
	return prog
}

func expectEval(t *testing.T, prog *ast.Program) object.Object {
	t.Helper()
	return eval.Eval(prog, object.New())
}
func expectLiteral(t *testing.T, got object.Object, expected any) {
	t.Helper()
	if expected == nil {
		if got != object.NULL {
			t.Fatalf("expected null, got %T (%q)", got, got)
		}
		return
	}
	switch e := expected.(type) {
	case int:
		expectIntegerLiteral(t, got, int64(e))
	case int64:
		expectIntegerLiteral(t, got, (e))
	case bool:
		expectBooleanLiteral(t, got, e)
	case error:
		expectErrorMessage(t, got, e.Error())
	default:
		t.Fatalf("unexpected type: %T %q", got, got)
	}
}

func expectIntegerLiteral(t *testing.T, got object.Object, expected int64) {
	t.Helper()
	v, ok := got.(*object.Integer)
	if !ok {
		t.Fatalf("expected *object.Integer, got %T (%q)", got, got)
	}
	if v.Value != expected {
		t.Fatalf("value mismatch: expected %d got %d", expected, v.Value)
	}
}
func expectBooleanLiteral(t *testing.T, got object.Object, expected bool) {
	t.Helper()
	v, ok := got.(*object.Boolean)
	if !ok {
		t.Fatalf("expected *object.Boolean, got %T (%q)", got, got)
	}
	if v.Value != expected {
		t.Fatalf("value mismatch: expected %t got %t", expected, v.Value)
	}
}

func expectErrorMessage(t *testing.T, got object.Object, expected string) {
	t.Helper()
	v, ok := got.(*object.Error)
	if !ok {
		t.Fatalf("expected *object.Error, got %T (%q)", got, got)
	}
	if v.Message != expected {
		t.Fatalf("value mismatch: expected %s got %s", expected, v.Message)
	}
}
