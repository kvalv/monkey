package token

type Type string

const (
	IDENT   Type = "IDENT"
	INT     Type = "INT"
	ILLEGAL Type = "ILLEGAL"
	COMMA   Type = ","

	POPEN  Type = "("
	PCLOSE Type = ")"
	LBRACK Type = "{"
	RBRACK Type = "}"

	FUNC   Type = "fn"
	RETURN Type = "return"
	LET    Type = "let"

	ASSIGN    Type = "="
	MUL       Type = "*"
	DIV       Type = "/"
	PLUS      Type = "+"
	MINUS     Type = "-"
	BANG      Type = "!"
	NEQ       Type = "!="
	SEMICOLON Type = "SEMICOLON"

	EOF Type = "EOF"
)

type Token struct {
	Type    Type
	Literal string
}
