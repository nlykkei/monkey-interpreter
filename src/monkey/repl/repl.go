package repl

import (
   "bufio"
   "fmt"
   "io"
   "monkey/lexer"
   "monkey/token"
   "monkey/ast"
   "monkey/parser"
   "monkey/object"
   "monkey/evaluator"
)

const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

const PROMPT = "> "

func Start(in io.Reader, out io.Writer) {
   scanner := bufio.NewScanner(in)
   env := object.NewEnvironment()

   for {
      fmt.Printf(PROMPT)
      scanned := scanner.Scan()
      if !scanned {
         return
      }


      line := scanner.Text()
      l := lexer.New(line)
      p := parser.New(l)

      printTokens(line)

      prog := p.ParseProgram()
      if len(p.Errors()) != 0 {
         printParserErrors(out, p.Errors())
         continue
      }

      printAST(prog, out)

      eval := evaluator.Eval(prog, env)
      if eval != nil {
         io.WriteString(out, eval.Inspect())
         io.WriteString(out, "\n")
      }
   }
}

func printTokens(line string) {
   l := lexer.New(line)
   for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
      fmt.Printf("%s\n", tok)
   }
}

func printAST(prog *ast.Program, out io.Writer) {
   io.WriteString(out, prog.String())
   io.WriteString(out, "\n")
}

func printParserErrors(out io.Writer, errors []string) {
   io.WriteString(out, MONKEY_FACE)
   io.WriteString(out, "Whoops! We ran into some monkey business here!\n")
   io.WriteString(out, "parser errors:\n")
   for _, msg := range errors {
      io.WriteString(out, "\t" + msg + "\n")
   }
}


