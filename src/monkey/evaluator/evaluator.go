package evaluator

import (
   "fmt"
   "monkey/ast"
   "monkey/object"
)

var (
   NULL = &object.Null{}
   TRUE = &object.Boolean{Value: true}
   FALSE = &object.Boolean{Value: false}
)

/*
 * Tree-Walking Interpreter
 *    ~ recursively interpret AST "on the fly", without any preprocessing or compilation step.
 */
func Eval(node ast.Node, env *object.Environment) object.Object {
   switch node := node.(type) {
      // Statements
      case *ast.Program:
         return evalProgram(node, env)
      case *ast.LetStatement:
         value := Eval(node.Value, env) // note: lexical scope
         if isError(value) {
            return value
         }
         env.Set(node.Name.Value, value) // note: identifier added to function's environment
      case *ast.ReturnStatement:
         value := Eval(node.ReturnValue, env)
         if isError(value) {
            return value
         }
         return &object.ReturnValue{Value: value}
      case *ast.ExpressionStatement:
         return Eval(node.Expression, env)
      case *ast.BlockStatement:
         return evalBlockStatement(node, env)
      // Expressions
      case *ast.PrefixExpression:
         right := Eval(node.Right, env)
         if isError(right) {
            return right
         }
         return evalPrefixExpression(node.Operator, right)
      case *ast.InfixExpression:
         left := Eval(node.Left, env)
         if isError(left) {
            return left
         }
         right := Eval(node.Right, env)
         if isError(right) {
            return right
         }
         return evalInfixExpression(node.Operator, left, right)
      case *ast.IfExpression:
         return evalIfExpression(node, env)
      case *ast.FunctionLiteral:
         params := node.Parameters
         body := node.Body
         return &object.Function{Parameters: params, Body: body, Env: env}
      case *ast.CallExpression:
         function := Eval(node.Function, env)
         if isError(function) {
            return function
         }
         args := evalExpressions(node.Arguments, env)
         if len(args) == 1 && isError(args[0]) {
            return args[0]
         }
         return applyFunction(function, args) 
      case *ast.Identifier:
         return evalIdentifier(node, env)
      case *ast.IntegerLiteral:
         return &object.Integer{Value: node.Value} // self-evaluating expression
      case *ast.StringLiteral:
         return &object.String{Value: node.Value}  // self-evaluating expression
      case *ast.Boolean:
         return nativeBoolToBoolObject(node.Value) // self-evaluating expression
      case *ast.ArrayLiteral:
         elements := evalExpressions(node.Elements, env)
         if len(elements) == 1 && isError(elements[0]) {
            return elements[0]
         }
         return &object.Array{Elements: elements}
      case *ast.IndexExpression:
         left := Eval(node.Left, env)
         if isError(left) {
            return left
         }
         index := Eval(node.Index, env)
         if isError(index) {
            return index
         }
         return evalIndexExpression(left, index)
      case *ast.HashLiteral:
         return evalHashLiteral(node, env)
   }
   return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
   var result object.Object
   for _, stmt := range program.Statements {
      result = Eval(stmt, env)
      switch result := result.(type) {
         case *object.ReturnValue:
            return result.Value
         case *object.Error: 
            return result
      }
   }
   return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
   var result object.Object

   for _, stmt := range block.Statements {
      result = Eval(stmt, env)

      switch result := result.(type) {
         case *object.ReturnValue:
            return result
         case *object.Error:
            return result
      }
   }

   return result
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
   var result []object.Object

   for _, exp := range exps {
      evaluated := Eval(exp, env)
      if isError(evaluated) {
         return []object.Object{evaluated}
      }
      result = append(result, evaluated)
   }

   return result
}

func evalPrefixExpression(op string, right object.Object) object.Object {
   switch op {
      case "!":
         return evalBangOperatorExpression(right)
      case "-":
         return evalMinusPrefixOperatorExpression(right)
      default:
         return newError("unknown operator: %s%s", op, right.Type())
   }
}

func evalInfixExpression(op string, left, right object.Object) object.Object {
   switch {
      case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
         return evalInfixIntegerExpression(op, left, right)
      case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
         return evalInfixStringExpression(op, left, right)
      case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
         return evalInfixBooleanExpression(op, left, right)
      case left.Type() != right.Type():
         return newError("type mismatch: %s %s %s", left.Type(), op, right.Type())
      default:
         return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
   }
}

func evalInfixIntegerExpression(op string, left, right object.Object) object.Object {
   leftVal := left.(*object.Integer).Value
   rightVal := right.(*object.Integer).Value 
   
   switch op {
      case "+":
         return &object.Integer{Value: leftVal + rightVal}
      case "-":
         return &object.Integer{Value: leftVal - rightVal}
      case "*":
         return &object.Integer{Value: leftVal * rightVal}
      case "/":
         return &object.Integer{Value: leftVal / rightVal}
      case "<":
         return nativeBoolToBoolObject(leftVal < rightVal)
      case ">":
         return nativeBoolToBoolObject(leftVal > rightVal)
      case "==":
         return nativeBoolToBoolObject(leftVal == rightVal)
      case "!=":
         return nativeBoolToBoolObject(leftVal != rightVal)
      default:
         return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
   }
}

func evalInfixStringExpression(op string, left, right object.Object) object.Object {
   leftVal := left.(*object.String).Value
   rightVal := right.(*object.String).Value

   switch op {
      case "+":
         return &object.String{Value: leftVal + rightVal}
      default:
         return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
   }
}

func evalInfixBooleanExpression(op string, left, right object.Object) object.Object {
   leftVal := left.(*object.Boolean).Value
   rightVal := right.(*object.Boolean).Value

   switch op {
      case "==":
         return nativeBoolToBoolObject(leftVal == rightVal)
      case "!=":
         return nativeBoolToBoolObject(leftVal != rightVal)
      case "&&":
         return nativeBoolToBoolObject(leftVal && rightVal)
      case "||":
         return nativeBoolToBoolObject(leftVal || rightVal)
      default:
         return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
   }
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
   condition := Eval(ie.Condition, env)
   if isError(condition) {
      return condition
   }
   if isTruthy(condition) {
      return Eval(ie.Consequence, env)
   } else if ie.Alternative != nil {
      return Eval(ie.Alternative, env)
   } else {
      return NULL
   }
}

// "truthy": any value not null or false
func isTruthy(obj object.Object) bool {
   switch obj {
      case TRUE:
         return true
      case FALSE:
         return false
      case NULL:
         return false
      default:
         return true
   }
}

func evalBangOperatorExpression(right object.Object) object.Object {
  switch right {
      case TRUE:
         return FALSE
      case FALSE:
         return TRUE
      case NULL:
         return TRUE
      default:
         return FALSE // !5
  }
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
   if right.Type() != object.INTEGER_OBJ {
      return newError("unknown operator: -%s", right.Type())
   }

   value := right.(*object.Integer).Value
   return &object.Integer{Value: -value}
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
   switch fn := fn.(type) {
   case *object.Function:
      extendedEnv := extendFunctionEnv(fn, args)
      evaluated := Eval(fn.Body, extendedEnv)
      return unwrapReturnValue(evaluated) // implicit return (last statement)
   case *object.Builtin:
      return fn.Fn(args...)
   default:
         return newError("not a function: %s", fn.Type())
   }
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
   env := object.NewExtendedEnvironment(fn.Env)
   for paramIdx, param := range fn.Parameters {
      env.Set(param.Value, args[paramIdx])
   }
   return env
}

func unwrapReturnValue(obj object.Object) object.Object {
   if returnValue, ok := obj.(*object.ReturnValue); ok {
      return returnValue.Value
   }
   return obj
}

func evalIndexExpression(left, index object.Object) object.Object {
   switch {
      case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
         return evalArrayIndexExpression(left, index)
      case left.Type() == object.HASH_OBJ:
         return evalHashIndexExpression(left, index)
      default:
         return newError("index operator not supported: %s", left.Type())
   }
}

func evalArrayIndexExpression(left, index object.Object) object.Object {
   array := left.(*object.Array)
   idx := index.(*object.Integer).Value
   max := int64(len(array.Elements) - 1)
   if idx < 0 || idx > max {
      return NULL
   }
   return array.Elements[idx]
}

func evalHashIndexExpression(left, index object.Object) object.Object {
   hash := left.(*object.Hash)
   key, ok := index.(object.Hashable)
   if !ok {
      return newError("unusable as hash key: %s", index.Type())
   }
   pair, ok := hash.Pairs[key.HashKey()]
   if !ok {
      return NULL
   }
   return pair.Value
}

func evalIdentifier(id *ast.Identifier, env *object.Environment) object.Object {
   if value, ok := env.Get(id.Value); ok {
      return value
   }
   if builtin, ok := builtins[id.Value]; ok {
      return builtin
   }
   return newError("identifier not found: " + id.Value)
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
   pairs := make(map[object.HashKey]object.HashPair)
   for nodeKey, nodeValue := range node.Pairs {
      key := Eval(nodeKey, env)
      if isError(key) {
         return key
      }
      hashKey, ok := key.(object.Hashable) // cast to interface
      if !ok {
         return newError("unusable as hash key: %s", key.Type())
      }
      value := Eval(nodeValue, env)
      if isError(value) {
         return value
      }
      hashed := hashKey.HashKey()
      pairs[hashed] = object.HashPair{Key: key, Value: value}
   }
   return &object.Hash{Pairs: pairs}
}

func nativeBoolToBoolObject(input bool) *object.Boolean {
   if input {
      return TRUE
   } 
   return FALSE
}

func newError(format string, a ...interface{}) *object.Error {
   return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
   if _, ok := obj.(*object.Error); ok {
      return true
   }
   return false
}
