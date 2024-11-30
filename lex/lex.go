package lex

import (
	"github.com/kvalv/monkey/token"
)

type Lex struct {
	input   string
	ch      byte // current
	pos     int
	peekPos int
	eof     bool
}

func New(input string) *Lex {
	return &Lex{input: input}
}

func (l *Lex) create(tp token.Type, lit string) token.Token {
	end := l.pos
	start := end - len(lit)
	return token.Token{Type: tp, Literal: lit, Span: token.Span{Start: start, End: end}}
}

func (l *Lex) NextToken() token.Token {

	l.skipWhitespace()

	c := l.nextChar()
	switch c {
	case 0:
		return l.create(token.EOF, "")
	case '*', '/', '+', '-', ',', '(', ')', '{', '}', ';', '>', '<':
		return l.create(builtins[string(c)], string(c))
	case '=':
		if l.peek() == '=' {
			l.nextChar()
			return l.create(token.EQ, "==")
		}
		return l.create(token.ASSIGN, "=")
	case '!':
		if l.peek() == '=' {
			l.nextChar()
			return l.create(token.NEQ, "!=")
		}
		return l.create(token.BANG, "!")
	default:
		if isLetter(c) {
			l.goBack()
			word := l.takeWhile(isLetter)
			return l.create(lookupIdentifier(word), word)
		}
		if isDigit(c) {
			l.goBack()
			word := l.takeWhile(isDigit)
			return l.create(token.INT, word)
		}
		return l.create(token.ILLEGAL, string(c))
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
	if !l.eof {
		l.goBack()
	}
	return l.input[start:end]
}

func (l *Lex) nextChar() byte {
	if l.eof {
		return 0
	}
	if l.peekPos >= len(l.input) {
		l.eof = true
		l.ch = 0
		return 0
	} else {
		l.ch = l.input[l.pos]
		l.pos += 1
		l.peekPos += 1
	}
	return l.ch
}

var builtins map[string]token.Type = map[string]token.Type{
	"=": token.ASSIGN,
	"+": token.PLUS,
	"-": token.MINUS,
	"*": token.MUL,
	"/": token.DIV,
	",": token.COMMA,
	"(": token.POPEN,
	")": token.PCLOSE,
	"{": token.LBRACK,
	"}": token.RBRACK,
	";": token.SEMICOLON,
	">": token.GT,
	"<": token.Lt,

	"if":     token.IF,
	"else":   token.ELSE,
	"let":    token.LET,
	"fn":     token.FUNC,
	"return": token.RETURN,
	"true":   token.TRUE,
	"false":  token.FALSE,
}

func lookupIdentifier(ident string) token.Type {
	if ttype, ok := builtins[ident]; ok {
		return ttype
	}
	return token.IDENT
}

func isWhitespace(c byte) bool { return c == ' ' || c == '\n' }
func isLetter(c byte) bool     { return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') }
func isDigit(c byte) bool      { return c >= '0' && c <= '9' }
