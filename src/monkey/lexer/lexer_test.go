   package lexer

   import (
      "testing"
      "monkey/token"
   )

   func TestNextToken(t *testing.T) {
      input := `=+(){},;`

      tests := []struct {
         expectedType    token.TokenType
         expectedLiteral string
      }{
         {token.ASSIGN, "="},
         {token.PLUS, "+"},
         {token.LPAREN, "("},
         {token.RPAREN, ")"},
         {token.LBRACE, "{"},
         {token.RBRACE, "}"},
         {token.COMMA, ","},
         {token.SEMICOLON, ";"},
         {token.EOF, ""},
   }

   lex := New(input)

   for i, ttok := range tests {
      tok := lex.NextToken()

      if tok.Type != ttok.expectedType {
         t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", 
            i, ttok.expectedType, tok.Type)
      }

      if tok.Literal != ttok.expectedLiteral {
         t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", 
            i, ttok.expectedLiteral, tok.Literal)
      }
   }
}
