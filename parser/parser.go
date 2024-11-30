package parser

import (
	"fmt"
	"strconv"

	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/lex"
	"github.com/kvalv/monkey/token"
)

type PrefixFn func() ast.Expression
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
	p.prefixFns[token.TRUE] = p.parseBoolean
	p.prefixFns[token.FALSE] = p.parseBoolean
	p.prefixFns[token.POPEN] = p.parseGroupExpression
	p.prefixFns[token.IF] = p.parseIfExpression

	p.infixFns[token.EQ] = p.parseInfixExpression
	p.infixFns[token.PLUS] = p.parseInfixExpression
	p.infixFns[token.MINUS] = p.parseInfixExpression
	p.infixFns[token.MUL] = p.parseInfixExpression
	p.infixFns[token.DIV] = p.parseInfixExpression
	p.infixFns[token.GT] = p.parseInfixExpression
	p.infixFns[token.Lt] = p.parseInfixExpression
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

func (p *Parser) parseIdentifier() ast.Expression {
	var out ast.Identifier
	defer trace("parseIdentifier", p.curr)(&out)
	if p.curr.Type != token.IDENT {
		p.errExpected(token.IDENT)
		return nil
	}
	out.Token = p.curr
	out.Value = p.curr.Literal
	return &out
}

func (p *Parser) parseLetStatement(precedence int) *ast.LetStatement {
	var stmt ast.LetStatement
	defer trace("parseLetStatement", p.curr)(&stmt)
	var ok bool
	if stmt.Token, ok = p.parseToken(token.LET); !ok {
		return nil
	}
	lhs := p.parseIdentifier()
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

func (p *Parser) parseGroupExpression() ast.Expression {
	defer trace("parseGroupExpression", p.curr)(nil)
	p.advance()
	exp := p.parseExpression(LOWEST)
	if p.next.Type != token.PCLOSE {
		p.errExpected(token.PCLOSE)
		return nil
	}
	p.advance()
	return exp
}
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	expr := &ast.ExpressionStatement{Token: p.curr}
	defer trace("parseExpressionStatement", p.curr)(expr)
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
	var out ast.Statement
	defer trace("parseStatement", p.curr)(out)
	switch p.curr.Type {
	case token.LET:
		out = p.parseLetStatement(LOWEST)
	default:
		out = p.parseExpressionStatement()
	}
	return out
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := ast.PrefixExpression{Token: p.curr, Op: p.curr.Literal}
	defer trace("parsePrefixExpression", p.curr)(&exp)
	p.advance()
	rhs := p.parseExpression(PREFIX)
	if rhs == nil {
		return nil
	}
	exp.Rhs = rhs
	return &exp
}

func (p *Parser) parseInfixExpression(precedence int, lhs ast.Expression) ast.Expression {
	var exp ast.InfixExpression
	defer trace("parseInfixExpression", p.curr)(&exp)
	exp = ast.InfixExpression{Token: p.curr, Op: p.curr.Literal, Lhs: lhs}
	p.advance()
	if exp.Rhs = p.parseExpression(precedence); exp.Rhs == nil {
		return nil
	}
	return &exp
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	var expr ast.Expression
	defer trace("parseExpression", p.curr)(expr)
	fn, ok := p.prefixFns[p.curr.Type]
	if !ok {
		p.errorf("prefixFn not found for type %v", p.curr.Type)
		return nil
	}
	expr = fn()

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

func (p *Parser) parseIfExpression() ast.Expression {
	var out ast.IfExpression
	defer trace("parseIfExpression", p.curr)(&out)
	out.Token = p.curr
	p.advance()
	if out.Cond = p.parseExpression(LOWEST); out.Cond == nil {
		return nil
	}
	p.advance()
	if out.Then = p.parseBlockStatement(); out.Then == nil {
		return nil
	}
	if p.next.Type == token.ELSE {
		p.advance()
		p.advance()
		if out.Else = p.parseBlockStatement(); out.Else == nil {
			return nil
		}
	}
	return &out
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	out := ast.BlockStatement{Statements: []ast.Statement{}}
	defer trace("parseBlockStatement", p.curr)(&out)
	out.Token = p.curr
	if p.curr.Type != token.LBRACK {
		p.errExpected(token.LBRACK)
		return nil
	}
	p.advance()
	for !(p.currIsType(token.EOF) || p.currIsType(token.RBRACK)) {
		stmt := p.parseStatement()
		if stmt != nil {
			out.Statements = append(out.Statements, stmt)
		} else {
			return nil
		}
	}
	if p.curr.Type != token.RBRACK {
		p.errExpected(token.RBRACK)
	}
	return &out
}

func (p *Parser) parseBoolean() ast.Expression {
	var out ast.Boolean
	defer trace("parseBoolean", p.curr)(&out)
	if !(p.curr.Type == token.FALSE || p.curr.Type == token.TRUE) {
		p.errExpected(token.FALSE, token.TRUE)
		return nil
	}
	out.Type = p.curr.Type
	out.Value = p.curr.Literal == "true"
	return &out
}

func (p *Parser) parseNumber() ast.Expression {
	var out ast.Number
	defer trace("parseNumber", p.curr)(&out)
	if p.curr.Type != token.INT {
		p.errExpected(token.INT)
		return nil
	}
	value, err := strconv.Atoi(p.curr.Literal)
	if err != nil {
		p.errorf("parseNumber: failed to parse as %q as number", p.curr.Literal)
		return nil
	}
	out.Token = p.curr
	out.Value = value
	return &out
}
