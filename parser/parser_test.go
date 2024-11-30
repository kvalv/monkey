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
		t.Fatalf("expected ExpressionStatement got %T", stmt)
	}
	expectLiteral(t, stmt.Expr, 3)
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
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
	}
	for _, tc := range tests {
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
