package parser

import (
	"fmt"
	"log"
	"strconv"

	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/lex"
	"github.com/kvalv/monkey/token"
)

type Parser struct {
	l          *lex.Lex
	curr, next token.Token
	errs       []error
}

func New(input string) *Parser {
	p := &Parser{
		l: lex.New(input),
	}
	p.advance() // populate curr and next
	p.advance()
	return p
}
func (p *Parser) advance() {
	if p.next.Type == token.EOF {
		p.curr = p.next
		return
	}
	p.curr = p.next
	log.Printf("advance: curr is now %+v", p.curr)
	p.next = p.l.NextToken()
}

// appends an error to the error list
func (p *Parser) errorf(format string, a ...any) { p.errs = append(p.errs, fmt.Errorf(format, a...)) }
func (p *Parser) errExpected(tp token.Type) {
	p.errorf("Parse(): expected %v but got %v", tp, p.curr.Type)
}

func (p *Parser) currIsType(tp token.Type) bool { return p.curr.Type == tp }
func (p *Parser) nextIsType(tp token.Type) bool { return p.next.Type == tp }

func (p *Parser) Parse() (*ast.Program, []error) {
	prog := &ast.Program{
		Statements: []ast.Statement{},
	}
	for !p.currIsType(token.EOF) {
		if stmt := p.parseStatement(); stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		} else {
			break
		}
	}
	return prog, p.errs
}

// advances if the current token type matches ttype. Otherwise, it does not advance, and returns false.
// The first value is the current token
func (p *Parser) parseToken(ttype token.Type) (token.Token, bool) {
	if !p.currIsType(ttype) {
		p.errExpected(ttype)
		return p.curr, false
	}
	tk := p.curr
	p.advance()
	return tk, true
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	iden, ok := p.parseToken(token.IDENT)
	if !ok {
		return nil
	}
	return &ast.Identifier{Token: iden, Value: iden.Literal}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	var (
		stmt ast.LetStatement
		ok   bool
	)
	if stmt.Token, ok = p.parseToken(token.LET); !ok {
		return nil
	}
	if stmt.Lhs = p.parseIdentifier(); stmt.Lhs == nil {
		return nil
	}
	if _, ok = p.parseToken(token.ASSIGN); !ok {
		return nil
	}
	if stmt.Rhs = p.parseExpression(); stmt.Rhs == nil {
		return nil
	}
	if _, ok = p.parseToken(token.SEMICOLON); !ok {
		return nil
	}
	return &stmt
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curr.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		p.errorf("parseStatement: unexpected token: %v", p.curr.Type)
		return nil
	}
}

func (p *Parser) parseExpression() ast.Expression {
	switch p.curr.Type {
	case token.IDENT:
		return p.parseIdentifier()
	case token.INT:
		return p.parseNumber()
	default:
		p.errorf("parseExpression: unexpected token: %v", p.curr.Type)
		return nil
	}
}

func (p *Parser) parseNumber() *ast.Number {
	token, ok := p.parseToken(token.INT)
	if !ok {
		return nil
	}
	value, err := strconv.Atoi(token.Literal)
	if err != nil {
		p.errorf("parseNumber: failed to parse as value: %v", err)
		return nil
	}
	return &ast.Number{Token: token, Value: value}
}
