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

	EQ        Type = "=="
	ASSIGN    Type = "="
	MUL       Type = "*"
	DIV       Type = "/"
	GT        Type = ">"
	Lt        Type = "<"
	PLUS      Type = "+"
	MINUS     Type = "-"
	BANG      Type = "!"
	NEQ       Type = "!="
	SEMICOLON Type = "SEMICOLON"
	TRUE      Type = "true"
	FALSE     Type = "false"

	EOF Type = "EOF"
)

type Token struct {
	Type    Type
	Literal string
}
