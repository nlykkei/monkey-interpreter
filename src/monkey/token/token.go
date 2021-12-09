package token

import (
   "fmt"
)

type TokenType string

type Token struct {
	Type     TokenType
	Literal  string
   Position SourcePosition
}

type SourcePosition struct {
   Line int
   Char int
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifier
	IDENT = "IDENT"

	// Literal
	INT = "INT"
   STRING = "STRING"
   BOOL = "BOOL"

	// Operator
	ASSIGN = "="
	PLUS   = "+"
   MINUS = "-"
   BANG = "!"
   ASTERISK = "*"
   SLASH = "/"

   LT = "<"
   GT = ">"

   EQ = "=="
   NOT_EQ = "!="

   AND = "&&"
   OR = "||"

	// Delimiter
	COMMA     = "COMMA"
	SEMICOLON = "SEMICOLON"
   COLON     = "COLON"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
   LBRACKET = "["
   RBRACKET = "]"

	// Keyword
	FUNCTION = "FUNCTION"
   LET  = "LET"

   IF = "IF"
   ELSE = "ELSE" 
	RETURN = "RETURN"

   TRUE = "TRUE"
   FALSE = "FALSE"
)

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
   "if": IF,
   "else": ELSE,
   "return": RETURN,
   "true": TRUE,
   "false": FALSE,
}

func LookupIdent(ident string) TokenType {
	if tt, ok := keywords[ident]; ok {
		return tt
	}
	return IDENT
}

func (tok Token) String() string {
   return fmt.Sprintf("token{type: %v, literal: %q}", tok.Type, tok.Literal)
}

func (sp *SourcePosition) String() string {
   return fmt.Sprintf("position{line: %d, char: %d}", sp.Line, sp.Char)
}
