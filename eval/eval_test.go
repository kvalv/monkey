package eval_test

import (
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
			got := eval.Eval(prog)
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
		got := eval.Eval(prog)
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
			got := eval.Eval(prog)
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
			got := eval.Eval(prog)
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

func expectLiteral(t *testing.T, got object.Object, expected any) {
	t.Helper()
	if expected == nil {
		if got != object.NULL {
			t.Fatalf("expected null, got %T %+v", got, got)
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
	default:
		t.Fatalf("unexpected type: %T", expected)
	}
}

func expectIntegerLiteral(t *testing.T, got object.Object, expected int64) {
	t.Helper()
	v, ok := got.(*object.Integer)
	if !ok {
		t.Fatalf("expected *object.Integer, got %T", got)
	}
	if v.Value != expected {
		t.Fatalf("value mismatch: expected %d got %d", expected, v.Value)
	}
}
func expectBooleanLiteral(t *testing.T, got object.Object, expected bool) {
	t.Helper()
	v, ok := got.(*object.Boolean)
	if !ok {
		t.Fatalf("expected *object.Boolean, got %T", got)
	}
	if v.Value != expected {
		t.Fatalf("value mismatch: expected %t got %t", expected, v.Value)
	}
}
