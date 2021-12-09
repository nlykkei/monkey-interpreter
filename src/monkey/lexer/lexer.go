package lexer

import (
   "errors"
	"monkey/token"
   "fmt"
   "os"
)

type Lexer struct {
	input          string               // source code
	index          int                  // input index (current char)
	readIndex      int                  // read index (next char)
   ch             byte                 // current char
   position       token.SourcePosition // source position
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
   // initialize l.index, l.readIndex, and l.ch
	l.readChar() 
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

   l.skipWhitespace()

	switch l.ch {
	case '=':
      if l.peekChar() == '=' {
         tok = l.makeTwoCharToken(token.EQ)
      } else {
		   tok = l.newToken(token.ASSIGN)
      }
   case '&':
      if l.peekChar() == '&' {
         tok = l.makeTwoCharToken(token.AND)
      } else {
			tok = l.newToken(token.ILLEGAL) 
      }
   case '|':
      if l.peekChar() == '|' {
         tok = l.makeTwoCharToken(token.OR)
      } else {
         tok = l.newToken(token.ILLEGAL)
      }
	case '+':
		tok = l.newToken(token.PLUS)
   case '-':
      tok = l.newToken(token.MINUS)
   case '!':
      if l.peekChar() == '=' {
         tok = l.makeTwoCharToken(token.NOT_EQ)
      } else {
         tok = l.newToken(token.BANG)
      }
   case '*':
      tok = l.newToken(token.ASTERISK)
   case '/':
      tok = l.newToken(token.SLASH)
   case '<':
      tok = l.newToken(token.LT)
   case '>':
      tok = l.newToken(token.GT) 
	case ',':
		tok = l.newToken(token.COMMA)
	case ';':
		tok = l.newToken(token.SEMICOLON)
   case ':':
      tok = l.newToken(token.COLON)
	case '(':
		tok = l.newToken(token.LPAREN)
	case ')':
		tok = l.newToken(token.RPAREN)
	case '{':
		tok = l.newToken(token.LBRACE)
	case '}':
		tok = l.newToken(token.RBRACE)
   case '[':
      tok = l.newToken(token.LBRACKET)
   case ']':
      tok = l.newToken(token.RBRACKET)
   case '"':
      position := l.position 
      str, err := l.readString()
      if err == nil {
         tok = token.Token{Type: token.STRING, Literal: str, Position: position}
      } else {
         tok = token.Token{Type: token.ILLEGAL, Literal: str, Position: position}
         fmt.Fprintf(os.Stderr, "%s (%s)", err, position.String())
      }
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if isLetter(l.ch) {     
         // keyword or identifier
         position := l.position
			ident := l.readIdentifier() 
         return token.Token{Type: token.LookupIdent(ident), Literal: ident, Position: position} 
      } else if isDigit(l.ch) {
         // integer
         position := l.position
         num := l.readNumber()
         return token.Token{Type: token.INT, Literal: num, Position: position}
      } else {
         // error
			tok = l.newToken(token.ILLEGAL) 
		}
	}

	l.readChar() 
	return tok
}

func (l *Lexer) readChar() {
	if l.readIndex >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readIndex]
      l.position.Char += 1
	}
	l.index = l.readIndex
	l.readIndex += 1
}

func (l *Lexer) peekChar() byte {
  if l.readIndex >= len(l.input) {
      return 0 
  } else {
      return l.input[l.readIndex] 
  }
}

func (l *Lexer) readLiteral(pred func (byte) bool) string {
	index := l.index
	for pred(l.ch) {
		l.readChar()
	}
	return l.input[index:l.index]
}

func (l *Lexer) readIdentifier() string {
   return l.readLiteral(isLetter)
}

func isLetter(ch byte) bool {
	return ('a' <= ch) && (ch <= 'z') || ('A' <= ch) && (ch <= 'Z') || ch == '_'
}

func (l *Lexer) readNumber() string {
   return l.readLiteral(isDigit)
}

func isDigit(ch byte) bool {
   return ('0' <= ch) && (ch <= '9')
      
}

func (l *Lexer) readString() (string, error) {
   index := l.index + 1 // don't include "s in token
   for {
      l.readChar()
      switch (l.ch) {
         case '"':
            return l.input[index:l.index], nil
         case 0:
            return l.input[index:l.index], errors.New("readString: could not tokenize non-terminated string") 
      }
   }
}

func (l *Lexer) newToken(tt token.TokenType) token.Token {
	return token.Token{Type: tt, Literal: string(l.ch), Position: l.position}
}

func (l *Lexer) makeTwoCharToken(tt token.TokenType) token.Token {
   ch := l.ch
   position := l.position
   l.readChar() // consume second char
   return token.Token{Type: tt, Literal: string(ch) + string(l.ch), Position: position}
}

func (l *Lexer) skipWhitespace() {
   for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
      if l.ch == '\n' {
         l.position.Line += 1
         l.position.Char = 0 
      } 
      l.readChar()
   }
}
