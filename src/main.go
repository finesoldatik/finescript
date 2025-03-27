package main

import (
	"finescript/src/lexer"
	"finescript/src/parser"
	"finescript/src/runtime"
	"fmt"
	"os"
	"time"

	"github.com/sanity-io/litter"
)

func main() {
	sourceBytes, _ := os.ReadFile("../examples/test.fs")
	source := string(sourceBytes)

	startLexer := time.Now()
	tokens := lexer.Tokenize(source)
	durationLexer := time.Since(startLexer)

	startParser := time.Now()
	ast := parser.Parse(tokens, source)
	fmt.Println("\nAST:==================================")
	litter.Dump(ast)
	durationParser := time.Since(startParser)

	startInterpreter := time.Now()
	env := runtime.GlobalEnv()
	fmt.Println("\nRESULT:===============================")
	result := runtime.EvaluateStmt(ast, env)
	print("\n\nFINAL VALUE:==========================\n")
	println(runtime.Format(result))
	durationInterpreter := time.Since(startInterpreter)

	fmt.Println("\nTIME:=================================")
	fmt.Printf("Duration Lexer: %s\nDuration Parser: %s\nDuration Interpreter: %s\n", durationLexer, durationParser, durationInterpreter)
}
