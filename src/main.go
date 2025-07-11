package main

import (
	"finescript/src/lexer"
	"finescript/src/parser"
	"finescript/src/runtime"
	"fmt"
	"os"
	"time"

	"github.com/sanity-io/litter"
	"github.com/spf13/cobra"
)

var (
	showTokens bool
	showAST    bool
	showResult bool
	showTime   bool
	filePath   string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "finescript",
		Short: "A simple programming language.",
		Long:  "He is fine!",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				filePath = args[0]
			}

			startReadFile := time.Now()
			sourceBytes, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Error reading file: %v\n", err)
				os.Exit(1)
			}
			source := string(sourceBytes)
			durationReadFile := time.Since(startReadFile)

			startLexer := time.Now()
			tokens := lexer.Tokenize(source)
			durationLexer := time.Since(startLexer)

			startParser := time.Now()
			ast := parser.Parse(tokens, source)
			durationParser := time.Since(startParser)

			startInterpreter := time.Now()
			if showTokens || showAST || showResult || showTime {
				println("RUNTIME:===============================")
			}
			result := runtime.EvaluateStmt(ast, runtime.GlobalEnv())
			println()
			durationInterpreter := time.Since(startInterpreter)

			if showTokens {
				println("\nTOKENS:===============================")
				litter.Dump(tokens)
			}
			if showAST {
				println("\nAST:==================================")
				litter.Dump(ast)
			}
			if showResult {
				println("\nRESULT:===============================")
				println(runtime.Format(result))
			}
			if showTime {
				println("\nTIME:=================================")
				fmt.Printf("Duration Read File: %s\nDuration Lexer: %s\nDuration Parser: %s\nDuration Interpreter: %s\n",
					durationReadFile, durationLexer, durationParser, durationInterpreter)
			}
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&showTokens, "show-tokens", "t", false, "Enables program tokens visibility")
	rootCmd.PersistentFlags().BoolVarP(&showAST, "show-ast", "a", false, "Enables program AST visibility")
	rootCmd.PersistentFlags().BoolVarP(&showResult, "show-result", "r", false, "Enables program result visibility")
	rootCmd.PersistentFlags().BoolVarP(&showTime, "show-time", "s", false, "Enables program execute time visibility")

	rootCmd.Args = cobra.ExactArgs(1)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
