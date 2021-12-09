package ast

import (
   "fmt"
   "bytes"
   "strings"
   "monkey/token"
)

type Node interface {
   TokenLiteral() string
   String() string
}

type Statement interface {
   Node
   statementNode()
}

type Expression interface {
   Node
   expressionNode()
}

// <statement>*
type Program struct { // root node
   Statements []Statement
}

func (p *Program) TokenLiteral() string {
   if len(p.Statements) > 0 {
      return p.Statements[0].TokenLiteral()
   } else {
      return ""
   }
}

func (p *Program) String() string {
   var out bytes.Buffer

   length := len(p.Statements)
   switch length { 
   case 0:
      out.WriteString("")
   case 1:
      out.WriteString(p.Statements[0].String())
   default:
      for i, s := range p.Statements {
         out.WriteString(s.String())
         if i < length - 1 {
            out.WriteString("; ")
         }
      }
   }

   return out.String()
}

// let <identifier> = <expression>;
type LetStatement struct {
   Token token.Token // token.LET
   Name *Identifier
   Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

func (ls *LetStatement) String() string {
   var out bytes.Buffer

   out.WriteString(ls.TokenLiteral() + " " + ls.Name.String() + " = ")
   
   if ls.Value != nil {
      out.WriteString(ls.Value.String())
   }

   return out.String()
}

// return <expression>;
type ReturnStatement struct {
   Token token.Token // token.RETURN
   ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

func (rs *ReturnStatement) String() string {
   var out bytes.Buffer

   out.WriteString(rs.TokenLiteral() + " ")

   if rs.ReturnValue != nil {
      out.WriteString(rs.ReturnValue.String())
   }

   return out.String()
}

// <expression>[;]
type ExpressionStatement struct {
   Token token.Token
   Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

func (es *ExpressionStatement) String() string {
   if es.Expression != nil {
      return es.Expression.String() 
   }
   return ""
}

// { <statement>* }
type BlockStatement struct {
   Token token.Token // "{" token
   Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
   var out bytes.Buffer

   length := len(bs.Statements)
   switch length {
   case 0:
      out.WriteString("{ }")
   case 1:
      out.WriteString("{ " + bs.Statements[0].String() + " }") 
   default:
      out.WriteString("{ ")
      for i, s := range bs.Statements {
         out.WriteString(s.String())
         if i < length - 1 {
            out.WriteString("; ")
         }
      }
      out.WriteString(" }")
   }

   return out.String()
}

// [_A-Za-z]+
type Identifier struct {
   Token token.Token // token.IDENT
   Value string
}

func (id *Identifier) expressionNode() {}
func (id *Identifier) TokenLiteral() string { return id.Token.Literal }
func (id *Identifier) String() string { return id.Value }

// "[^"]*"
type StringLiteral struct {
   Token token.Token
   Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string { return sl.Token.Literal }

// [0-9]+
type IntegerLiteral struct {
   Token token.Token
   Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string { return il.Token.Literal }

// true|false
type Boolean struct {
   Token token.Token 
   Value bool
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string { return b.Token.Literal }

// <prefix-operator> <expression>
type PrefixExpression struct {
   Token token.Token
   Operator string
   Right Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
   var out bytes.Buffer

   out.WriteString("(")
   out.WriteString(pe.Operator)
   out.WriteString(pe.Right.String())
   out.WriteString(")")

   return out.String()
}

// <expression> <infix-operator> <expression>
type InfixExpression struct {
   Token token.Token
   Operator string
   Left Expression
   Right Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
   var out bytes.Buffer

   out.WriteString("(")
   out.WriteString(ie.Left.String())
   out.WriteString(" " + ie.Operator + " ")
   out.WriteString(ie.Right.String())
   out.WriteString(")")

   return out.String()
}

// if (<condition>) <consequence> [else <alernative>]
type IfExpression struct {
   Token token.Token // "if" token
   Condition Expression
   Consequence *BlockStatement
   Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
   var out bytes.Buffer

   out.WriteString("if (")
   out.WriteString(ie.Condition.String())
   out.WriteString(") ")
   out.WriteString(ie.Consequence.String())
   
   if ie.Alternative != nil {
      out.WriteString("else ")
      out.WriteString(ie.Alternative.String())
   }

   return out.String()
}

// fn (<parameter-list>) <block-statement>
type FunctionLiteral struct {
   Token token.Token // "fn" token
   Parameters []*Identifier
   Body *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
   var out bytes.Buffer

   params := []string{} 

   for _, p := range fl.Parameters {
      params = append(params, p.String())

   }

   out.WriteString(fl.TokenLiteral())
   out.WriteString("(")
   out.WriteString(strings.Join(params, ", "))
   out.WriteString(") ")
   out.WriteString(fl.Body.String())

   return out.String()
}

// <expression>(<parameter-list>)
type CallExpression struct {
   Token token.Token // "(" token
   Function Expression
   Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
   var out bytes.Buffer
   
   args := []string{}
   for _, a := range ce.Arguments {
      args = append(args, a.String())
   }

   out.WriteString(ce.Function.String())
   out.WriteString("(")
   out.WriteString(strings.Join(args, ", "))
   out.WriteString(")")
   return out.String()
}

// [<expression-list>]
type ArrayLiteral struct {
   Token token.Token // token.LBRACKET
   Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
   var out bytes.Buffer

   elements := []string{}
   for _, el := range al.Elements {
      elements = append(elements, el.String())
   }

   out.WriteString("[")
   out.WriteString(strings.Join(elements, ", "))
   out.WriteString("]")
   return out.String()
}

// {<expression>:<expression>,...}
type HashLiteral struct {
   Token token.Token // token.LBRACE
   Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode() {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
   var out bytes.Buffer

   pairs := []string{}
   for key, value := range hl.Pairs {
      pairs = append(pairs, fmt.Sprintf("%s:%s", key.String(), value.String()))
   }

   out.WriteString("{")
   out.WriteString(strings.Join(pairs, ", "))
   out.WriteString("}")
   return out.String()
}


// <expression>[<expression>]
type IndexExpression struct {
   Token token.Token // token.IDENT
   Left Expression
   Index Expression
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
   var out bytes.Buffer 
   out.WriteString("(")
   out.WriteString(ie.Left.String())
   out.WriteString("[")
   out.WriteString(ie.Index.String())
   out.WriteString("])")
   return out.String()
}
