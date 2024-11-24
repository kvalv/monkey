package parser

import (
	"fmt"
	"strconv"

	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/lex"
	"github.com/kvalv/monkey/token"
)

type PrefixFn func(precedence int) ast.Expression
type InfixFn func(precedence int, lhs ast.Expression) ast.Expression

type Parser struct {
	l          *lex.Lex
	curr, next token.Token
	errs       []error

	prefixFns map[token.Type]PrefixFn
	infixFns  map[token.Type]InfixFn
}

func New(input string) *Parser {
	p := &Parser{
		l:         lex.New(input),
		prefixFns: make(map[token.Type]PrefixFn),
		infixFns:  make(map[token.Type]InfixFn),
	}
	p.advance() // populate curr and next
	p.advance()
	p.prefixFns[token.BANG] = p.parsePrefixExpression
	p.prefixFns[token.MINUS] = p.parsePrefixExpression
	p.prefixFns[token.INT] = p.parseNumber
	p.prefixFns[token.IDENT] = p.parseIdentifier

	p.infixFns[token.PLUS] = p.parseInfixExpression
	p.infixFns[token.MINUS] = p.parseInfixExpression
	p.infixFns[token.MUL] = p.parseInfixExpression
	p.infixFns[token.DIV] = p.parseInfixExpression
	return p
}
func (p *Parser) advance() {
	if p.next.Type == token.EOF {
		p.curr = p.next
		return
	}
	p.curr = p.next
	p.next = p.l.NextToken()
}

// appends an error to the error list
func (p *Parser) errorf(format string, a ...any) { p.errs = append(p.errs, fmt.Errorf(format, a...)) }
func (p *Parser) errExpected(tp ...token.Type) {
	if len(tp) == 1 {
		p.errorf("Parse(): expected %v but got %v", tp[0], p.curr.Type)
	} else if len(tp) > 1 {
		p.errorf("Parse(): expected one of %s, but got %v", tp, p.curr.Type)
	}
}

func (p *Parser) currIsType(tp ...token.Type) bool {
	for _, t := range tp {
		if p.curr.Type == t {
			return true
		}
	}
	return false
}
func (p *Parser) nextIsType(tp ...token.Type) bool {
	for _, t := range tp {
		if p.next.Type == t {
			return true
		}
	}
	return false
}

func (p *Parser) Parse() (*ast.Program, []error) {
	prog := &ast.Program{
		Statements: []ast.Statement{},
	}
	for !p.currIsType(token.EOF) {
		if len(p.errs) > 0 { // hack - remove
			break
		}
		stmt := p.parseStatement()
		if stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		} else {
			break
		}
	}
	return prog, p.errs
}

// advances if the current token type matches ttype. Otherwise, it does not advance, and returns false.
// The first value is the current token
func (p *Parser) parseToken(ttype ...token.Type) (token.Token, bool) {
	if !p.currIsType(ttype...) {
		p.errExpected(ttype...)
		return p.curr, false
	}
	tk := p.curr
	p.advance()
	return tk, true
}

func (p *Parser) parseIdentifier(precedence int) ast.Expression {
	if p.curr.Type != token.IDENT {
		p.errExpected(token.IDENT)
		return nil
	}
	return &ast.Identifier{Token: p.curr, Value: p.curr.Literal}
}

func (p *Parser) parseLetStatement(precedence int) *ast.LetStatement {
	var (
		stmt ast.LetStatement
		ok   bool
	)
	if stmt.Token, ok = p.parseToken(token.LET); !ok {
		return nil
	}
	lhs := p.parseIdentifier(precedence)
	if lhs == nil {
		return nil
	}
	stmt.Lhs = lhs.(*ast.Identifier) // also probably not ideal..
	if _, ok = p.parseToken(token.ASSIGN); !ok {
		return nil
	}
	if stmt.Rhs = p.parseExpression(LOWEST); stmt.Rhs == nil {
		return nil
	}
	if _, ok = p.parseToken(token.SEMICOLON); !ok {
		return nil
	}
	return &stmt
}
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	expr := &ast.ExpressionStatement{Token: p.curr}
	exp := p.parseExpression(LOWEST)
	if exp == nil {
		return nil
	}
	expr.Expr = exp
	p.advance()
	if p.curr.Type == token.SEMICOLON {
		// semicolons are optional and ignored in expression statements
		// ... so we can omit semicolons in the repl
		p.advance()
	}
	return expr
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curr.Type {
	case token.LET:
		return p.parseLetStatement(LOWEST)
	default:
		got := p.parseExpressionStatement()
		return got
	}
}

func (p *Parser) parsePrefixExpression(precedence int) ast.Expression {
	exp := ast.PrefixExpression{Token: p.curr, Op: p.curr.Literal}
	p.advance()
	rhs := p.parseExpression(precedence)
	if rhs == nil {
		return nil
	}
	exp.Rhs = rhs
	return &exp
}

func (p *Parser) parseInfixExpression(precedence int, lhs ast.Expression) ast.Expression {
	exp := ast.InfixExpression{Token: p.curr, Op: p.curr.Literal, Lhs: lhs}
	p.advance()
	if exp.Rhs = p.parseExpression(precedence); exp.Rhs == nil {
		return nil
	}
	return &exp
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	fn, ok := p.prefixFns[p.curr.Type]
	if !ok {
		p.errorf("prefixFn not found for type %v", p.curr.Type)
		return nil
	}
	expr := fn(precedence)

	for p.next.Type != token.SEMICOLON && precedence < tokenPrecedence(p.next.Type) {
		p.advance()
		fn, ok := p.infixFns[p.curr.Type]
		if !ok {
			p.errorf("infixFn not found for type %v", p.next.Type)
			return nil
		}
		expr = fn(tokenPrecedence(p.curr.Type), expr)
	}

	return expr
}

func (p *Parser) parseNumber(precedence int) ast.Expression {
	if p.curr.Type != token.INT {
		p.errExpected(token.INT)
		return nil
	}
	value, err := strconv.Atoi(p.curr.Literal)
	if err != nil {
		p.errorf("parseNumber: failed to parse as %q as number", p.curr.Literal)
		return nil
	}
	return &ast.Number{Token: p.curr, Value: value}
}
