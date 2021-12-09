package lexer

import (
	"monkey/token"
	"testing"
)

type ExpectedToken struct {
	expectedType    token.TokenType
	expectedLiteral string
}

func TestNextToken(t *testing.T) {
	input := `let five = 5;
   let ten = 10;
   
   let add = fn(x, y) {
      x + y
   };
   
   let result = add(five, ten);

   !-/*5;
   5 < 10 > 5;

   if (5 < 10) {
      return true;
   } else {
      return false;
   }

   10 == 10;
   10 != 9;

   let s = "monkey";
   let ns = "";

   false || true;
   true && false;

   [1, 2];
   a[0];

   {"foo": "bar"};
   `
	
	tests := []struct{
      expectedType token.TokenType
      expectedLiteral string
   }{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
      {token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
      {token.BANG, "!"},
      {token.MINUS, "-"},
      {token.SLASH, "/"},
      {token.ASTERISK, "*"},
      {token.INT, "5"},
      {token.SEMICOLON, ";"},
      {token.INT, "5"},
      {token.LT, "<"},
      {token.INT, "10"},
      {token.GT, ">"},
      {token.INT, "5"},
      {token.SEMICOLON, ";"},
      {token.IF, "if"},
      {token.LPAREN, "("},
      {token.INT, "5"},
      {token.LT, "<"},
      {token.INT, "10"},
      {token.RPAREN, ")"},
      {token.LBRACE, "{"},
      {token.RETURN, "return"},
      {token.TRUE, "true"},
      {token.SEMICOLON, ";"},
      {token.RBRACE, "}"},
      {token.ELSE, "else"},
      {token.LBRACE, "{"},
      {token.RETURN, "return"},
      {token.FALSE, "false"},
      {token.SEMICOLON, ";"},
      {token.RBRACE, "}"},
      {token.INT, "10"},
      {token.EQ, "=="},
      {token.INT, "10"},
      {token.SEMICOLON, ";"},
      {token.INT, "10"},
      {token.NOT_EQ, "!="},
      {token.INT, "9"},
      {token.SEMICOLON, ";"},
      {token.LET, "let"},
      {token.IDENT, "s"},
      {token.ASSIGN, "="},
      {token.STRING, "monkey"},
      {token.SEMICOLON, ";"},
      {token.LET, "let"},
      {token.IDENT, "ns"},
      {token.ASSIGN, "="},
      {token.STRING, ""},
      {token.SEMICOLON, ";"},
      {token.FALSE, "false"},
      {token.OR, "||"},
      {token.TRUE, "true"},
      {token.SEMICOLON, ";"},      
      {token.TRUE, "true"},
      {token.AND, "&&"},
      {token.FALSE, "false"},
      {token.SEMICOLON, ";"},
      {token.LBRACKET, "["},
      {token.INT, "1"},
      {token.COMMA, ","},
      {token.INT, "2"},
      {token.RBRACKET, "]"},
      {token.SEMICOLON, ";"},
      {token.IDENT, "a"},
      {token.LBRACKET, "["},
      {token.INT, "0"},
      {token.RBRACKET, "]"},
      {token.SEMICOLON, ";"},
      {token.LBRACE, "{"},
      {token.STRING, "foo"},
      {token.COLON, ":"},
      {token.STRING, "bar"},
      {token.RBRACE, "}"},
      {token.SEMICOLON, ";"},
      {token.EOF, ""},
	}

   l := New(input)

	for i, ttok := range tests {
		tok := l.NextToken()

		if tok.Type != ttok.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, ttok.expectedType, tok.Type)
		}

		if tok.Literal != ttok.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, ttok.expectedLiteral, tok.Literal)
		}
	}
}
