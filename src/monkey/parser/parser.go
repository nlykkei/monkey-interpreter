package parser

import (
   "fmt"
   "strconv"
   "monkey/ast"
   "monkey/lexer"
   "monkey/token"
)

const (
   _ int = iota
   LOWEST         
   EQUALS         // ==
   OR             // ||
   AND            // &&
   LESSGREATER    // < or >
   SUM            // +
   PRODUCT        // *
   PREFIX         // -x or !x
   CALL           // myFunc()
   INDEX          // a[1]
)

var precedences = map[token.TokenType]int{
   token.RPAREN:     LOWEST,
   token.EQ:         EQUALS,
   token.NOT_EQ:     EQUALS,
   token.OR:         OR,
   token.AND:        AND,
   token.LT:         LESSGREATER,
   token.GT:         LESSGREATER,
   token.PLUS:       SUM,
   token.MINUS:      SUM,
   token.SLASH:      PRODUCT,
   token.ASTERISK:   PRODUCT,
   token.LPAREN:     CALL,
   token.LBRACKET:   INDEX, 
}

/*
 * Parsing Function Protocol:
 *    ~ parsing function call: p.curToken is the first token of the statement/expression
 *    ~ parsing function return: p.curToken is the last token of the statement/expression       
 */
type Parser struct {
   l *lexer.Lexer
   curToken token.Token
   peekToken token.Token
   errors []string

   prefixParseFns map[token.TokenType]prefixParseFn
   infixParseFns map[token.TokenType]infixParseFn
}


/*
 * Pratt parsing: token types are associated with up to two parsing functions: 
 *     ~ prefix parse function
 *     ~ infix parse function
 */
type (
   prefixParseFn func() ast.Expression
   infixParseFn func(ast.Expression) ast.Expression // left expression
)

func (p *Parser) registerPrefixFn(tokenType token.TokenType, fn prefixParseFn) {
   p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfixFn(tokenType token.TokenType, fn infixParseFn) {
   p.infixParseFns[tokenType] = fn
}

func New(l *lexer.Lexer) *Parser {
   p := &Parser{l: l, errors: []string{}}

   // prefix parse functions
   p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
   p.registerPrefixFn(token.IDENT, p.parseIdentifier)
   p.registerPrefixFn(token.INT, p.parseIntegerLiteral)
   p.registerPrefixFn(token.STRING, p.parseStringLiteral)
   p.registerPrefixFn(token.TRUE, p.parseBoolean)
   p.registerPrefixFn(token.FALSE, p.parseBoolean)
   p.registerPrefixFn(token.BANG, p.parsePrefixExpression)
   p.registerPrefixFn(token.MINUS, p.parsePrefixExpression)
   p.registerPrefixFn(token.LPAREN, p.parseGroupedExpression) // token.LPAREN
   p.registerPrefixFn(token.IF, p.parseIfExpression)
   p.registerPrefixFn(token.FUNCTION, p.parseFuncLiteral)
   p.registerPrefixFn(token.LBRACKET, p.parseArrayLiteral)
   p.registerPrefixFn(token.LBRACE, p.parseHashLiteral)

   // infix parse functions
   p.infixParseFns = make(map[token.TokenType]infixParseFn)
   p.registerInfixFn(token.PLUS, p.parseInfixExpression)
   p.registerInfixFn(token.MINUS, p.parseInfixExpression)
   p.registerInfixFn(token.SLASH, p.parseInfixExpression)
   p.registerInfixFn(token.ASTERISK, p.parseInfixExpression)
   p.registerInfixFn(token.EQ, p.parseInfixExpression)
   p.registerInfixFn(token.NOT_EQ, p.parseInfixExpression)
   p.registerInfixFn(token.OR, p.parseInfixExpression)
   p.registerInfixFn(token.AND, p.parseInfixExpression)
   p.registerInfixFn(token.LT, p.parseInfixExpression)
   p.registerInfixFn(token.GT, p.parseInfixExpression)
   p.registerInfixFn(token.LPAREN, p.parseCallExpression) // token.LPAREN
   p.registerInfixFn(token.LBRACKET, p.parseIndexExpression) 

   // initialize p.curToken and p.peekToken
   p.nextToken() 
   p.nextToken()

   return p
}

func (p *Parser) Errors() []string {
   return p.errors
}

func (p *Parser) nextToken() token.Token {
   p.curToken = p.peekToken
   p.peekToken = p.l.NextToken()
   return p.curToken
}

func (p *Parser) ParseProgram() *ast.Program {
//   defer untrace(trace("ParseProgram"))

   prog := &ast.Program{}
   prog.Statements = []ast.Statement{}

   for !p.curTokenIs(token.EOF) {
      stmt := p.parseStatement()
      if stmt != nil {
         prog.Statements = append(prog.Statements, stmt)
      }
      p.nextToken() // consume (optional) token.SEMICOLON 
   }

   return prog
}

func (p *Parser) parseStatement() ast.Statement {
//   defer untrace(trace("parseStatement"))

   switch p.curToken.Type {
      case token.LET:
         return p.parseLetStatement()
      case token.RETURN:
         return p.parseReturnStatement()
      default:
         return p.parseExpressionStatement()
   }
}

func (p *Parser) parseLetStatement() ast.Statement {
//   defer untrace(trace("parseLetStatement"))

   stmt := &ast.LetStatement{Token: p.curToken}

   if !p.expectPeek(token.IDENT) {
      return nil
   }

   stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

   if !p.expectPeek(token.ASSIGN) {
      return nil
   }

   p.nextToken() // consume token.ASSIGN

   stmt.Value = p.parseExpression(LOWEST)
   
   if p.peekTokenIs(token.SEMICOLON) {
      p.nextToken()
   }

   return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
//   defer untrace(trace("parseReturnStatement"))

   stmt := &ast.ReturnStatement{Token: p.curToken}

   p.nextToken()

   stmt.ReturnValue = p.parseExpression(LOWEST)

   if p.peekTokenIs(token.SEMICOLON) {
      p.nextToken()
   }

   return stmt
}

func (p *Parser) parseExpressionStatement() ast.Statement {
//   defer untrace(trace("parseExpressionStatement"))

   stmt := &ast.ExpressionStatement{Token: p.curToken}

   stmt.Expression = p.parseExpression(LOWEST) // LOWEST parses everything left-to-right

   if p.peekTokenIs(token.SEMICOLON) { 
      p.nextToken() // align p.curToken with (optional) token.SEMICOLON 
   }

   return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
//   defer untrace(trace("parseBlockStatement"))

   bs := &ast.BlockStatement{Token: p.curToken}
   bs.Statements = []ast.Statement{}

   p.nextToken() // consume token.LBRACE

   for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
      stmt := p.parseStatement()
      if stmt != nil {
         bs.Statements = append(bs.Statements, stmt)
      }
      p.nextToken()
   }

   return bs
}

/*
 * Top Down Operator Precedence (Pratt) parser 
 *    ~ precedence: right-binding power (current operator)
 *    ~ peekPrecedence: left-binding power (next operator)
 */
func (p *Parser) parseExpression(precedence int) ast.Expression { 
//   defer untrace(trace("parseExpression"))

   prefix := p.prefixParseFns[p.curToken.Type] 

   if prefix == nil {
      msg := fmt.Sprintf("parseExpression: found no prefix parse function for %s (%s)", p.curToken.Type, p.curToken.Position.String())
      p.errors = append(p.errors, msg)
      return nil
   }

   leftExp := prefix() 

   for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() { // token.LPAREN, token.EOF etc. have LOWEST precedence (unspecified)
      infix := p.infixParseFns[p.peekToken.Type]

      if infix == nil {
         return leftExp
      }

      p.nextToken()
      leftExp = infix(leftExp)
   }
   
   return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
//  defer untrace(trace("parsePrefixExpression"))

  pe := &ast.PrefixExpression{
     Token: p.curToken,
     Operator: p.curToken.Literal,
  }

  p.nextToken()
  pe.Right = p.parseExpression(PREFIX)
  return pe
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
//   defer untrace(trace("parseInfixExpression"))

   ie := &ast.InfixExpression{
      Token: p.curToken,
      Operator: p.curToken.Literal,
      Left: left,
   }

   precedence := p.curPrecedence()
   p.nextToken()
   ie.Right = p.parseExpression(precedence)

   return ie
}

func (p *Parser) parseIfExpression() ast.Expression {
//   defer untrace(trace("parseIfExpression"))
   
   ie := &ast.IfExpression{Token: p.curToken}

   if !p.expectPeek(token.LPAREN) {
      return nil
   }

   p.nextToken()  // consume token.LPAREN
   ie.Condition = p.parseExpression(LOWEST)

   if !p.expectPeek(token.RPAREN) {
      return nil
   }

   if !p.expectPeek(token.LBRACE) {
      return nil
   }

   ie.Consequence = p.parseBlockStatement()

   if p.peekTokenIs(token.ELSE) {
      p.nextToken() // consume token.RBRACE

      if !p.expectPeek(token.LBRACE) {
         return nil 
      }

      ie.Alternative = p.parseBlockStatement()
   }

   return ie
}

func (p *Parser) parseIdentifier() ast.Expression {
//   defer untrace(trace("parseIdentifier"))

   return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
//   defer untrace(trace("parseIntegerLiteral"))

   il := &ast.IntegerLiteral{Token: p.curToken}

   val, err := strconv.ParseInt(p.curToken.Literal, 0, 64)

   if err != nil {
      msg := fmt.Sprintf("parseIdentifier: could not parse %s as integer (%s)", p.curToken.Literal, p.curToken.Position.String())
      p.errors = append(p.errors, msg)
      return nil
   }

   il.Value = val

   return il
}

func (p *Parser) parseStringLiteral() ast.Expression {
   return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
//   defer untrace(trace("parseBoolean"))

   return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)} // only called for token.TRUE and token.FALSE
}

func (p *Parser) parseGroupedExpression() ast.Expression {
   p.nextToken() // consume token.LPAREN
   
   exp := p.parseExpression(LOWEST) // parse until matching token.RPAREN

   if !p.expectPeek(token.RPAREN) {
      return nil
   }

   return exp
}

func (p *Parser) parseFuncLiteral() ast.Expression {
//   defer untrace(trace("parseFuncLiteral"))

   fl := &ast.FunctionLiteral{Token: p.curToken}

   if !p.expectPeek(token.LPAREN) {
      return nil
   }
   
   fl.Parameters = p.parseFuncParameters() // (x, y, ...)

   if !p.expectPeek(token.LBRACE) {
      return nil
   }

   fl.Body = p.parseBlockStatement()

   return fl
}

func (p *Parser) parseFuncParameters() []*ast.Identifier {
//   defer untrace(trace("parseFuncParameters"))

   ids := []*ast.Identifier{}

   if p.peekTokenIs(token.RPAREN) { // ()
      p.nextToken() // consume token.LPAREN
      return ids
   }

   p.nextToken()

   id := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
   ids = append(ids, id) 

   for p.peekTokenIs(token.COMMA) {
      p.nextToken() // consume token.IDENT
      p.nextToken() // consume token.COMMA
      id = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
      ids = append(ids, id)
   }

   if !p.expectPeek(token.RPAREN) {
      return nil
   }

   return ids
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
//   defer untrace(trace("parseCallExpression"))

   exp := &ast.CallExpression{Token: p.curToken, Function: function}
   exp.Arguments = p.parseExpressionList(token.RPAREN) // (x, y, ...)
   return exp
}

func (p *Parser) parseArrayLiteral() ast.Expression {
   array := &ast.ArrayLiteral{Token: p.curToken}
   array.Elements = p.parseExpressionList(token.RBRACKET)
   return array
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
//   defer untrace(trace("parseExpressionList"))

   list := []ast.Expression{}

   if p.peekTokenIs(end) {
      p.nextToken() // align p.curToken with end token
      return list
   }

   p.nextToken() // consume start token
   list = append(list, p.parseExpression(LOWEST))

   for p.peekTokenIs(token.COMMA) {
      p.nextToken() 
      p.nextToken() // consume token.COMMA
      list = append(list, p.parseExpression(LOWEST))
   }
   
   if !p.expectPeek(end) { // align p.curToken with end token
      return nil
   } 

   return list
}

func (p *Parser) parseHashLiteral() ast.Expression {
   hash := &ast.HashLiteral{Token: p.curToken}
   hash.Pairs = make(map[ast.Expression]ast.Expression)
   for !p.peekTokenIs(token.RBRACE) {
      p.nextToken()
      key := p.parseExpression(LOWEST)
      if !p.expectPeek(token.COLON) {
         return nil
      }
      p.nextToken() // consume token.COLON
      value := p.parseExpression(LOWEST)
      hash.Pairs[key] = value
      if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) { // expect token.LBRACE or token.COMMA
         return nil
      }
   }
   if !p.expectPeek(token.RBRACE) {
      return nil
   }
   return hash
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
//   defer untrace(trace("parseIndexExpression"))

   ie := &ast.IndexExpression{Token: p.curToken, Left: left}
   p.nextToken() // consume token.LBRACKET
   ie.Index = p.parseExpression(LOWEST)
   if !p.expectPeek(token.RBRACKET) {
      return nil
   }
   return ie
}

func (p *Parser) curTokenIs(tt token.TokenType) bool {
   return p.curToken.Type == tt
}

func (p *Parser) peekTokenIs(tt token.TokenType) bool {
   return p.peekToken.Type == tt
}

func (p *Parser) expectPeek(tt token.TokenType) bool {
   if p.peekTokenIs(tt) {
      p.nextToken()
      return true
   } else {
      msg := fmt.Sprintf("expectPeek: wrong peek token type. expected=%q, got=%q (%s)", tt, p.peekToken.Type, p.peekToken.Position.String())
      p.errors = append(p.errors, msg)
      return false
   }
}

func (p *Parser) curPrecedence() int {
   if p, ok := precedences[p.curToken.Type]; ok {
      return p
   }

   return LOWEST
}

func (p *Parser) peekPrecedence() int {   
   if p, ok := precedences[p.peekToken.Type]; ok {
      return p
   }

   return LOWEST
}
