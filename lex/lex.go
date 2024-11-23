package lex

import (
	"github.com/kvalv/monkey/token"
)

type Lex struct {
	input   string
	ch      byte // current
	pos     int
	peekPos int
}

func New(input string) *Lex {
	return &Lex{input: input}
}
func (l *Lex) NextToken() token.Token {

	l.skipWhitespace()

	c := l.nextChar()
	switch c {
	case 0:
		return token.Token{Ttype: token.EOF, Literal: ""}
	case '=', '+', '-', ',', '(', ')', '{', '}':
		return token.Token{Ttype: builtins[string(c)], Literal: string(c)}
	case '!':
		if l.peek() == '=' {
			l.nextChar()
			return token.Token{Ttype: token.NEQ, Literal: "!="}
		}
		return token.Token{Ttype: token.BANG, Literal: "!"}
	default:
		if isLetter(c) {
			l.goBack()
			word := l.takeWhile(isLetter)
			return token.Token{Ttype: lookupIdentifier(word), Literal: word}
		}
		if isDigit(c) {
			l.goBack()
			word := l.takeWhile(isDigit)
			return token.Token{Ttype: token.INT, Literal: word}
		}
		return token.Token{Ttype: token.ILLEGAL}
	}
}

func (l *Lex) goBack() {
	if l.pos > 0 {
		l.pos -= 1
		l.peekPos -= 1
	}
}
func (l *Lex) peek() byte {
	if l.peekPos >= len(l.input) {
		return 0
	}
	return l.input[l.peekPos]
}

func (l *Lex) skipWhitespace() {
	for isWhitespace(l.peek()) {
		l.nextChar()
	}
}
func (l *Lex) takeWhile(pred func(c byte) bool) string {
	start := l.pos
	end := l.pos
	for pred(l.ch) {
		end = l.pos
		l.nextChar()
	}
	return l.input[start:end]
}

func (l *Lex) nextChar() byte {
	if l.peekPos >= len(l.input) {
		l.peekPos = 0
		l.ch = 0
	} else {
		l.ch = l.input[l.pos]
		l.pos += 1
		l.peekPos += 1
	}
	return l.ch
}

var builtins map[string]token.TokenType = map[string]token.TokenType{
	"=": token.EQ,
	"+": token.PLUS,
	"-": token.MINUS,
	",": token.COMMA,
	"(": token.POPEN,
	")": token.PCLOSE,
	"{": token.LBRACK,
	"}": token.RBRACK,

	"fn":     token.FUNC,
	"return": token.RETURN,
}

func lookupIdentifier(ident string) token.TokenType {
	if ttype, ok := builtins[ident]; ok {
		return ttype
	}
	return token.IDENT
}

func isWhitespace(c byte) bool { return c == ' ' || c == '\n' }
func isLetter(c byte) bool     { return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') }
func isDigit(c byte) bool      { return c >= '0' && c <= '9' }
