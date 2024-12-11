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

func TestStringExpression(t *testing.T) {
	cases := []struct {
		input    string
		expected any
	}{
		{`"hello"`, "hello"},
		{`"he" + "llo"`, "hello"},
		{`"he" + ""`, "he"},
		{`"" + ""`, ""},
		{`"x" == "x"`, true},
		{`"x" == "y"`, false},
	}
	for _, tc := range cases {
		prog := expectParse(t, tc.input)
		got := expectEval(t, prog)
		expectLiteral(t, got, tc.expected)
	}
}

func TestBuiltinFunction(t *testing.T) {
	cases := []struct {
		input    string
		expected any
	}{
		{`len("1234")`, 4},
		{`len("ab" + "cd")`, 4},
		// {`len("")`, 0}, // TODO :S
		{`len(2)`, fmt.Errorf("type error: expected STRING but got INTEGER")},
	}
	for _, tc := range cases {
		prog := expectParse(t, tc.input)
		got := expectEval(t, prog)
		expectLiteral(t, got, tc.expected)
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
		{"let a = 10; let b = a > 7; let c = if b { 99 } else { 98 }", 99},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			prog := expectParse(t, tc.input)
			got := expectEval(t, prog)
			expectLiteral(t, got, tc.expected)
		})
	}
}

func TestFunctionLiteral(t *testing.T) {
	cases := []struct {
		input  string
		pcount int
		stmts  []string
	}{
		{"fn(x) { return 2 }", 1, []string{"return 2"}},
		{"fn(x, y) { x + y }", 2, []string{"(x + y)"}},
		{"fn() { true }", 0, []string{"true"}},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			prog := expectParse(t, tc.input)
			got := expectEval(t, prog)
			f, ok := got.(*object.Function)
			if !ok {
				t.Fatalf("not a *object.Function - got %T %q", got, got)
			}
			if n := len(f.Params); n != tc.pcount {
				t.Fatalf("expected %d parameter, got %d", tc.pcount, n)
			}
			if n := len(f.Body.Statements); n != len(tc.stmts) {
				t.Fatalf("expected %d statements, got %d", len(tc.stmts), n)
			}
			for i, expected := range tc.stmts {
				got := f.Body.Statements[i].String()
				if got != expected {
					t.Fatalf("expected  %q got %q", expected, got)
				}
			}
		})
	}
}

func TestFunctionLiteralError(t *testing.T) {
	t.Run("repeated arguments", func(t *testing.T) {
		got := expectEval(t, expectParse(t, "fn(x, x) { x + x }"))
		expectErrorMessage(t, got, `repeated argument "x"`)
	})
}

func TestFunctionApplication(t *testing.T) {
	cases := []struct {
		input    string
		expected any
	}{
		{"let plus = fn(x) { return x + 2 }; plus(2)", 4},
		{"fn(x) { return x + 2 }(2)", 4},
		{"fn(x) { x + 2 }(2)", 4},
		{"let add = fn(x, y) { return x + y }; add(1, 2)", 3},
		{"fn(x, y) { return x + y }(1, 2)", 3},
		{"fn(x, y) { x + y }(1, 2)", 3},
		{"let apply = fn(f, in) { f(in) }; apply(fn(x) { x + 2 }, 2)", 4},
		{"let a = 1; fn(x) { x + a }(1)", 2},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			prog := expectParse(t, tc.input)
			got := expectEval(t, prog)
			expectLiteral(t, got, tc.expected)
		})
	}
}

func TestArrayIndexing(t *testing.T) {
	cases := []struct {
		input    string
		expected any
	}{
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][true]", fmt.Errorf("object type mismatch")},
		{"[1, 2, 3][4]", fmt.Errorf("List index out of range: 4 > 3")},
		{"let index = 2; [1, 2, 3][index]", 3},
		{"[1, 2, 3][1 - 1 - 1 + 1]", 1},
		{"[1, 2, 3][-123]", fmt.Errorf("negative indices not allowed")},
		{"[2+2][0]", 4},
		{"rest([1, 2, 3])", []any{2, 3}},
		{"first([1, 2, 3])", 1},
		{"first([])", nil},
		{"push([1, 2], 3, 4)", []any{1, 2, 3, 4}},
		{"push([1])", []any{1}},
		{"push([])", []any{}},
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
	return eval.Eval(prog, object.NewEnvironment())
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
	case string:
		expectStringLiteral(t, got, e)
	case error:
		expectErrorMessage(t, got, e.Error())
	case []any:
		arr, ok := got.(*object.Array)
		if !ok {
			t.Fatalf("not an array, got %T", got)
		}
		if len(arr.Elems) != len(e) {
			t.Fatalf("length mismatch - got %d expected %d elements", len(arr.Elems), len(e))
		}
		for i, exp := range e {
			got := arr.Elems[i]
			expectLiteral(t, got, exp)
		}
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
func expectStringLiteral(t *testing.T, got object.Object, expected string) {
	t.Helper()
	v, ok := got.(*object.String)
	if !ok {
		t.Fatalf("expected *object.String, got %T (%q)", got, got)
	}
	if v.Value != expected {
		t.Fatalf("value mismatch: expected %q got %q", expected, v.Value)
	}
}

func expectErrorMessage(t *testing.T, got object.Object, expected string) {
	t.Helper()
	v, ok := got.(*object.Error)
	if !ok {
		t.Fatalf("expected *object.Error, got %T (%q)", got, got)
	}
	if v.Message != expected {
		t.Fatalf("value mismatch: expected %q got %q", expected, v.Message)
	}
}
