package token

type TokenType string

const (
	IDENT   TokenType = "IDENT"
	INT     TokenType = "INT"
	ILLEGAL TokenType = "ILLEGAL"
	COMMA   TokenType = ","

	POPEN  TokenType = "("
	PCLOSE TokenType = ")"
	LBRACK TokenType = "{"
	RBRACK TokenType = "}"

	FUNC   TokenType = "fn"
	RETURN TokenType = "return"

	EQ    TokenType = "EQ"
	PLUS  TokenType = "+"
	MINUS TokenType = "-"
	BANG  TokenType = "!"
	NEQ   TokenType = "!="

	EOF TokenType = "EOF"
)

type Token struct {
	Ttype   TokenType
	Literal string
}
