package parser_test

import (
	"testing"

	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/parser"
)

func TestParsePrefix(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{input: "!3;", want: "(!3)"},
		{input: "-foo", want: "(-foo)"},
	}
	for i, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			p := parser.New(tc.input)
			prog, errs := p.Parse()
			if len(errs) > 0 {
				t.Fatalf("got %d errors: %+v", len(errs), errs)
			}
			got := prog.Statements[0].String()
			if got != tc.want {
				t.Fatalf("[%d]: mismatch: got %q, want %q", i, got, tc.want)
			}
		})
	}
}

func TestParseInfixExpression(t *testing.T) {
	cases := []struct {
		input string
		lhs   any
		op    string
		rhs   any
	}{
		{"true == false", true, "==", false},
		{"false > true", false, ">", true},
		{"true != false", true, "!=", false},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			p := parser.New(tc.input)
			prog, errs := p.Parse()
			if len(errs) > 0 {
				t.Fatalf("got %d errors: %+v", len(errs), errs)
			}
			got := prog.Statements[0]
			exprStmt, ok := got.(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("expected expression statement, got %T", exprStmt)
			}
			expectInfixExpression(t, exprStmt.Expr, tc.lhs, tc.op, tc.rhs)
		})
	}
}

func TestParseExpression(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{input: "1 + 2 + 3", want: "((1 + 2) + 3)"},
		{input: "1 + 2 * 3", want: "(1 + (2 * 3))"},
	}
	for i, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			p := parser.New(tc.input)
			prog, errs := p.Parse()
			if len(errs) > 0 {
				t.Fatalf("got %d errors: %+v", len(errs), errs)
			}
			if n := len(prog.Statements); n != 1 {
				t.Fatalf("expected 1 statement, got %d: prog.Statements=%+v", n, prog.Statements)
			}
			got := prog.Statements[0].String()
			if got != tc.want {
				t.Fatalf("[%d]: mismatch: got %q, want %q", i, got, tc.want)
			}
		})
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{{"true", true}, {"false", false}}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			p := parser.New(tc.input)
			prog, err := p.Parse()
			if err != nil {
				t.Fatalf("got error %v", err)
			}
			if n := len(prog.Statements); n != 1 {
				t.Fatalf("expected 1 statement but got %d", n)
			}
			stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("expected ExpressionStatement got %T", stmt)
			}
			expectLiteral(t, stmt.Expr, tc.expected)
		})
	}
}

func TestParseString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello"`, "hello"},
		{`"hello world"`, "hello world"},
		{`""`, ""},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			p := parser.New(tc.input)
			prog, err := p.Parse()
			if err != nil {
				t.Fatalf("got error %v", err)
			}
			if n := len(prog.Statements); n != 1 {
				t.Fatalf("expected 1 statement but got %d", n)
			}
			stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("expected ExpressionStatement got %T", stmt)
			}
			expectString(t, stmt.Expr, tc.expected)
		})
	}
}

func TestPrefixParse(t *testing.T) {
	p := parser.New("3")
	prog, err := p.Parse()
	if err != nil {
		t.Fatalf("got error %v", err)
	}
	if n := len(prog.Statements); n != 1 {
		t.Fatalf("expected 1 statement but got %d", n)
	}
	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement got %T", prog.Statements[0])
	}
	expectLiteral(t, stmt.Expr, 3)
}

func TestReturnExpression(t *testing.T) {
	cases := []struct {
		input    string
		expected any
	}{
		{"return 4", 4},
		{"return true", true},
		{"return 12", 12},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			prog, err := parser.New(tc.input).Parse()
			if err != nil {
				t.Fatalf("got error %v", err)
			}
			if n := len(prog.Statements); n != 1 {
				t.Fatalf("expected 1 statement but got %d", n)
			}
			stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("expected ExpressionStatement got %T", prog.Statements[0])
			}
			rexpr, ok := stmt.Expr.(*ast.ReturnExpression)
			if !ok {
				t.Fatalf("expected ast.ReturnExpression got %T", rexpr)
			}
			expectLiteral(t, rexpr.Value, tc.expected)
		})
	}
}

func TestFunctionLiteral(t *testing.T) {
	tests := []struct{ input, expected string }{
		{"fn () { 2 }", "fn() {2}"},
		{"fn (x) { x + 2 }", "fn(x) {(x + 2)}"},
		{"fn (x, y) { x + y }", "fn(x, y) {(x + y)}"},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			p := parser.New(tc.input)
			prog, err := p.Parse()
			if err != nil {
				t.Fatalf("got error %v", err)
			}
			if n := len(prog.Statements); n != 1 {
				t.Fatalf("expected 1 statement but got %d", n)
			}
			stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("expected FunctionLiteral got %T", prog.Statements[0])
			}
			f, ok := stmt.Expr.(*ast.FunctionLiteral)
			if !ok {
				t.Fatalf("expected *FunctionLiteral got %T", stmt)
			}
			got := f.String()
			if got != tc.expected {
				t.Fatalf("string mismatch: expected %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestCallExpression(t *testing.T) {
	input := "concat(1, 2, a + b)"
	prog, err := parser.New(input).Parse()
	if err != nil {
		t.Fatalf("got error %v", err)
	}
	callExp := expectAstNode[*ast.CallExpression](t, prog)
	expectLiteral(t, callExp.Function, "concat")
	if n := len(callExp.Params); n != 3 {
		t.Fatalf("expected 3 params, got %d", n)
	}
	expectLiteral(t, callExp.Params[0], 1)
	expectLiteral(t, callExp.Params[1], 2)
	expectInfixExpression(t, callExp.Params[2], "a", "+", "b")
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	cases := []struct {
		input, expected string
	}{
		{"-a + b", "((-a) + b)"},
		{"!-a", "(!(-a))"},
		{"a+b+c", "((a + b) + c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b * c + d / e -f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{" 3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"(3 + 4) * 5", "((3 + 4) * 5)"},
		{"3 + (4 + 5)", "(3 + (4 + 5))"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			p := parser.New(tc.input)
			prog, err := p.Parse()
			if err != nil {
				t.Fatalf("got error %v", err)
			}
			got := prog.String()
			if got != tc.expected {
				t.Fatalf("expected %q got %q", tc.expected, got)
			}
		})
	}
}

func TestLetStatement(t *testing.T) {
	cases := []struct {
		input, expected string
	}{
		{"let x = 5", "let x = 5"},
		{"let y = true", "let y = true"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			p := parser.New(tc.input)
			prog, err := p.Parse()
			if err != nil {
				t.Fatalf("got error %v", err)
			}
			got := prog.String()
			if got != tc.expected {
				t.Fatalf("expected %q got %q", tc.expected, got)
			}
		})
	}
}

func TestIfExpression(t *testing.T) {
	p := parser.New("if true { 2 } else { 3 }")
	prog, err := p.Parse()
	if err != nil {
		t.Fatalf("got error %v", err)
	}
	if n := len(prog.Statements); n != 1 {
		t.Fatalf("expected 1 statement but got %d", n)
	}

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement got %T", prog.Statements[0])
	}
	expStmt, ok := stmt.Expr.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expected *ast.IfExpression got %T", stmt)
	}
	expectLiteral(t, expStmt.Cond, true)

	{ // lhs
		if n := len(expStmt.Then.Statements); n != 1 {
			t.Fatalf("expected 1 statement, got %d", n)
		}
		lhs, ok := expStmt.Then.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("expected expression statement got %T", expStmt.Then.Statements[0])
		}
		expectLiteral(t, lhs.Expr, 2)
	}

	{ // rhs
		if n := len(expStmt.Else.Statements); n != 1 {
			t.Fatalf("expected 1 statement, got %d", n)
		}
		rhs, ok := expStmt.Else.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("expected expression statement got %T", expStmt.Else.Statements[0])
		}
		expectLiteral(t, rhs.Expr, 3)
	}
}
