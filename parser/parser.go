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
func (p *Parser) parseToken(ttype ...token.Type) (token.Token, bool) {
	if !p.currIsType(ttype...) {
		p.errExpected(ttype...)
		return p.curr, false
	}
	tk := p.curr
	p.advance()
	return tk, true
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	if p.curr.Type != token.IDENT {
		p.errExpected(token.IDENT)
		return nil
	}
	return &ast.Identifier{Token: p.curr, Value: p.curr.Literal}
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
	p.errorf("Not yet implemented")
	return nil
}

func (p *Parser) parseNumber() *ast.Number {
	if p.curr.Type != token.IDENT {
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
