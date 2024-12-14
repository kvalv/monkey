package parser

import "github.com/kvalv/monkey/token"

const (
	_ int = iota
	LOWEST
	ASSIGN
	EQ
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	ARRAY_INDEX
	FUNCTION_CALL
)

var lookup map[token.Type]int = map[token.Type]int{
	token.ASSIGN: ASSIGN,
	token.EQ:     EQ,
	token.NEQ:    EQ,
	token.BANG:   PREFIX,
	token.DIV:    PRODUCT,
	token.MINUS:  SUM,
	token.MUL:    PRODUCT,
	token.PLUS:   SUM,
	token.GT:     LESSGREATER,
	token.Lt:     LESSGREATER,
	token.POPEN:  FUNCTION_CALL,
	token.SOPEN:  ARRAY_INDEX,
}

func tokenPrecedence(ttype token.Type) int {
	v, ok := lookup[ttype]
	if !ok {
		return LOWEST
	}
	return v
}
