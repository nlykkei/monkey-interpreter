package lexer

import "monkey/token"

type Lexer struct {
   input        string // source code
   position     int    // input position (current char)
   readPosition int    // reading position (next char)
   ch           byte   // current char
}

func New(input string) *Lexer {
   lex := &Lexer{input: input}
   lex.readChar()
   return lex
}

func (lex *Lexer) NextToken() token.Token {
   var tok token.Token

   switch lex.ch {
      case '=':
         tok = newToken(token.ASSIGN, lex.ch)
      case '+':
         tok = newToken(token.PLUS, lex.ch)
      case ',':
         tok = newToken(token.COMMA, lex.ch)
      case ';':
         tok = newToken(token.SEMICOLON, lex.ch)
      case '(':
         tok = newToken(token.LPAREN, lex.ch)
      case ')':
         tok = newToken(token.RPAREN, lex.ch)
      case '{':
         tok = newToken(token.LBRACE, lex.ch)
      case '}':
         tok = newToken(token.RBRACE, lex.ch)
      case 0:
         tok.Type = token.EOF
         tok.Literal = ""
   }

   lex.readChar()
   return tok
}

func (lex *Lexer) readChar() {
   if lex.readPosition >= len(lex.input) {
      lex.ch = 0
   } else {
      lex.ch = lex.input[lex.readPosition]
   }
   
   lex.position = lex.readPosition
   lex.readPosition += 1
}

func newToken(tokType token.TokenType, ch byte) token.Token {
   return token.Token{Type: tokType, Literal: string(ch)}
 }
