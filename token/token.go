package token

type Type string

const (
	IDENT   Type = "IDENT"
	INT     Type = "INT"
	ILLEGAL Type = "ILLEGAL"
	STRING  Type = "STRING"
	COMMA   Type = ","
	COLON   Type = ":"

	POPEN  Type = "("
	PCLOSE Type = ")"
	LBRACK Type = "{"
	RBRACK Type = "}"
	SOPEN  Type = "["
	SCLOSE Type = "]"

	FUNC   Type = "fn"
	RETURN Type = "return"
	LET    Type = "let"
	IF     Type = "if"
	ELSE   Type = "else"

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

type Span struct {
	Start, End int
}

type Token struct {
	Type    Type
	Literal string
	// Span marks the position from start (inclusive) to end (exclusive)
	Span
}
