package token

type TokenType string

type Token struct {
   Type TokenType
   Literal string
}

const (
   ILLEGAL = "ILLEGAL"
   EOF = "EOF"

   // Identifier
   IDENT = "IDENT"

   // Literal
   INT = "INT"

   // Operator
   ASSIGN = "ASSIGN"
   PLUS = "PLUS"

   // Delimiter
   COMMA = "COMMA"
   SEMICOLON = "SEMICOLON"

   LPAREN = "("
   RPAREN = ")"
   LBRACE = "{"
   RBRACE = "}"

   // Keyword
   FUNC = "FUNC"
   LET = "LET"
)

