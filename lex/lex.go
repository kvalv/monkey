package lex

import (
	"fmt"

	"github.com/kvalv/monkey/token"
)

func New(input string) *Lex {
	return &Lex{input: input}
}

type Lex struct {
	input string
	pos   int
}

func (l *Lex) create(tp token.Type, lit string) token.Token {
	if tp == token.EOF {
		return token.Token{
			Type: tp, Literal: "", Span: token.Span{
				Start: len(l.input),
				End:   len(l.input),
			},
		}
	}
	end := l.pos + 1
	start := end - len(lit)
	return token.Token{Type: tp, Literal: lit, Span: token.Span{Start: start, End: end}}
}

func (l *Lex) curr() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}
func (l *Lex) peek() byte {
	if l.pos+1 >= len(l.input) {
		return 0
	}
	return l.input[l.pos+1]
}
func (l *Lex) advance() {
	if l.pos+1 > len(l.input) {
		l.pos = len(l.input)
	} else {
		l.pos++
	}
}

func (l *Lex) Next() (t token.Token) {
	defer func() {
	}()
	if l.curr() == 0 {
		return l.create(token.EOF, "")
	}
	defer l.advance() // advance one byte so we're at the start of next token when we're done here

	l.takeWhile(isWhitespace, true) // ... we want to skip whitespace

	c := l.curr()
	if tp, ok := builtins[string(c)]; ok {
		// all single tokens should match here; =, +, -, (, ...
		return l.create(tp, string(c))
	}
	if c == '=' {
		if l.peek() == '=' {
			l.advance()
			return l.create(token.EQ, "==")
		}
		return l.create(token.ASSIGN, "=")
	}
	if c == '!' {
		if l.peek() == '=' {
			l.advance()
			return l.create(token.NEQ, "!=")
		}
		return l.create(token.BANG, "!")
	}
	if c == '"' {
		l.advance()
		word := l.takeWhile(func(b byte) bool { return b != '"' }, true)
		return l.create(token.STRING, fmt.Sprintf(`"%s`, word))
	}
	// yeah otherwise we'll check for longer tokens: digits and letters
	if isLetter(c) {
		word := l.takeWhile(isLetter, false)
		if typ, ok := builtins[word]; ok {
			// it's a special keyword, such as "if" or "return"
			return l.create(typ, word)
		}
		return l.create(token.IDENT, word)
	}
	if isDigit(c) {
		// being an identifier is OK too!
		return l.create(token.INT, l.takeWhile(isDigit, false))
	}

	return l.create(token.ILLEGAL, string(c))
}

// takewhile consumes until the peek token evaluates to false.
// It also consumes the current token if consume is true
func (l *Lex) takeWhile(pred func(c byte) bool, consume bool) string {
	// starts at current, stops at the end
	if !pred(l.curr()) {
		return ""
	}
	start := l.pos
	for {
		c := l.peek()
		if c == 0 {
			break
		}
		if !pred(c) {
			if consume {
				l.advance()
			}
			break
		}
		l.advance()
	}
	end := l.pos
	return l.input[start : end+1]
}

var builtins map[string]token.Type = map[string]token.Type{
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
