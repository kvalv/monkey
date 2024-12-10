package parser

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/lex"
	"github.com/kvalv/monkey/token"
	"github.com/kvalv/monkey/tracer"
)

type PrefixFn func() ast.Expression
type InfixFn func(precedence int, lhs ast.Expression) ast.Expression

type Parser struct {
	l          *lex.Lex
	curr, next token.Token
	errs       []error
	tracer     *tracer.Tracer

	prefixFns map[token.Type]PrefixFn
	infixFns  map[token.Type]InfixFn
}

type parseOpt func(p *Parser)

func EnableTracing() parseOpt {
	return func(p *Parser) { p.tracer = tracer.New(os.Stdout) }
}

func New(input string, opts ...parseOpt) *Parser {
	p := &Parser{
		l:         lex.New(input),
		prefixFns: make(map[token.Type]PrefixFn),
		infixFns:  make(map[token.Type]InfixFn),
		tracer:    tracer.New(io.Discard),
	}
	p.advance() // populate curr and next
	p.advance()
	p.prefixFns[token.BANG] = p.parsePrefixExpression
	p.prefixFns[token.MINUS] = p.parsePrefixExpression
	p.prefixFns[token.INT] = p.parseNumber
	p.prefixFns[token.STRING] = p.parseString
	p.prefixFns[token.IDENT] = p.parseIdentifier
	p.prefixFns[token.TRUE] = p.parseBoolean
	p.prefixFns[token.FALSE] = p.parseBoolean
	p.prefixFns[token.POPEN] = p.parseGroupExpression
	p.prefixFns[token.IF] = p.parseIfExpression
	p.prefixFns[token.FUNC] = p.parseFunctionLiteral
	p.prefixFns[token.RETURN] = p.parseReturnExpression

	p.infixFns[token.NEQ] = p.parseInfixExpression
	p.infixFns[token.EQ] = p.parseInfixExpression
	p.infixFns[token.PLUS] = p.parseInfixExpression
	p.infixFns[token.MINUS] = p.parseInfixExpression
	p.infixFns[token.MUL] = p.parseInfixExpression
	p.infixFns[token.DIV] = p.parseInfixExpression
	p.infixFns[token.GT] = p.parseInfixExpression
	p.infixFns[token.Lt] = p.parseInfixExpression
	p.infixFns[token.POPEN] = p.parseCallExpression // todo: function call
	return p
}
func (p *Parser) advance() {
	if p.next.Type == token.EOF {
		p.curr = p.next
		return
	}
	p.curr = p.next
	p.next = p.l.Next()
}

// appends an error to the error list
func (p *Parser) errorf(format string, a ...any) { p.errs = append(p.errs, fmt.Errorf(format, a...)) }
func (p *Parser) errExpected(tp ...token.Type) {
	if len(tp) == 1 {
		p.errorf("Parse(): expected %v but got %v at %d..%d", tp[0], p.curr.Type, p.curr.Start, p.curr.End)
	} else if len(tp) > 1 {
		p.errorf("Parse(): expected one of %s, but got %v at %d..%d", tp, p.curr.Type, p.curr.Start, p.curr.End)
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
			p.advance()
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
	defer p.tracer.Trace("parseIdentifier")(&out)
	if p.curr.Type != token.IDENT {
		p.errExpected(token.IDENT)
		return nil
	}
	out.Token = p.curr
	out.Value = p.curr.Literal
	return &out
}

func (p *Parser) parseLetStatement(precedence int) *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curr}
	p.advance()
	defer p.tracer.Trace("parseLetStatement")(stmt)
	lhs := p.parseIdentifier()
	if lhs == nil {
		return nil
	}
	stmt.Lhs = lhs.(*ast.Identifier) // also probably not ideal..
	p.advance()
	if _, ok := p.parseToken(token.ASSIGN); !ok {
		return nil
	}
	if stmt.Rhs = p.parseExpression(LOWEST); stmt.Rhs == nil {
		return nil
	}
	if p.nextIsType(token.SEMICOLON) {
		p.advance()
	}
	return stmt
}

func (p *Parser) parseGroupExpression() ast.Expression {
	defer p.tracer.Trace("parseGroupExpression")(nil)
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
	defer p.tracer.Trace("parseExpressionStatement")(expr)
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
	defer p.tracer.Trace("parseStatement")(out)
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
	defer p.tracer.Trace("parsePrefixExpression")(&exp)
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
	defer p.tracer.Trace("parseInfixExpression")(&exp)
	exp = ast.InfixExpression{Token: p.curr, Op: p.curr.Literal, Lhs: lhs}
	p.advance()
	if exp.Rhs = p.parseExpression(precedence); exp.Rhs == nil {
		return nil
	}
	return &exp
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	var expr ast.Expression
	defer p.tracer.Trace("parseExpression")(expr)
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
	defer p.tracer.Trace("parseIfExpression")(&out)
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
	defer p.tracer.Trace("parseBlockStatement")(&out)
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
	defer p.tracer.Trace("parseBoolean")(&out)
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
	defer p.tracer.Trace("parseNumber")(&out)
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

func (p *Parser) parseString() ast.Expression {
	var out ast.String
	defer p.tracer.Trace("parseString")(&out)
	if p.curr.Type != token.STRING {
		p.errExpected(token.STRING)
		return nil
	}
	end := len(p.curr.Literal) - 1
	stripped := p.curr.Literal[1:end]
	out.Token = p.curr
	out.Value = stripped
	return &out
}

func (p *Parser) parseParamList() []ast.Identifier {
	out := []ast.Identifier{}
	if p.curr.Type != token.POPEN {
		p.errExpected(token.POPEN)
		return nil
	}
	p.advance()

	for {
		switch p.curr.Type {
		case token.IDENT:
			out = append(out, *p.parseIdentifier().(*ast.Identifier))
			p.advance()
		case token.COMMA:
			p.advance()
		case token.PCLOSE:
			return out
		default:
			p.errExpected(token.IDENT, token.COMMA)
			return nil
		}
	}
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	out := ast.FunctionLiteral{Token: p.curr}
	defer p.tracer.Trace("parseFunctionLiteral")(&out)
	if p.curr.Type != token.FUNC {
		p.errExpected(token.FUNC)
		return nil
	}
	p.advance()
	if out.Params = p.parseParamList(); out.Params == nil {
		return nil
	}
	p.advance()
	if out.Body = p.parseBlockStatement(); out.Body == nil {
		return nil
	}
	return &out
}
func (p *Parser) parseReturnExpression() ast.Expression {
	out := &ast.ReturnExpression{Token: p.curr}
	defer p.tracer.Trace("parseReturnExpression")(out)
	p.advance()
	if out.Value = p.parseExpression(LOWEST); out.Value == nil {
		return nil
	}
	return out
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}
	if p.next.Type == token.PCLOSE {
		p.advance()
		return args
	}
	p.advance()
	args = append(args, p.parseExpression(LOWEST))

	for p.next.Type == token.COMMA {
		p.advance()
		p.advance()
		args = append(args, p.parseExpression(LOWEST))
	}

	if p.next.Type != token.PCLOSE {
		p.errExpected(token.PCLOSE)
		return nil
	}
	p.advance()
	return args
}

func (p *Parser) parseCallExpression(precedence int, left ast.Expression) ast.Expression {
	out := &ast.CallExpression{}
	defer p.tracer.Trace("parseCallExpression")(out)
	out.Function = left
	out.Params = p.parseCallArguments()
	if out.Params == nil {
		panic("params nil")
	}
	return out
}
