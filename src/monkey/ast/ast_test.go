package ast

import (
   "monkey/token"
   "testing"
)

func TestString(t *testing.T) {
   prog := &Program {
      Statements: []Statement{
         &LetStatement{
            Token: token.Token{Type: token.LET, Literal: "let"},
            Name: &Identifier{
               Token: token.Token{Type: token.IDENT, Literal: "myVar"},
               Value: "myVar",
            },
            Value: &Identifier{
               Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
               Value: "anotherVar",
            },
         },
      },
   }

   expected := "let myVar = anotherVar" 
   if prog.String() != expected {
      t.Errorf("prog.String() wrong. got=%q, want=%q", prog.String(), expected)
   }
}
