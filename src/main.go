package main

import (
	"bufio"
	"finescript/src/lexer"
	"finescript/src/parser"
	"finescript/src/runtime"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sanity-io/litter"
	"github.com/spf13/cobra"
)

var (
	showTokens,
	showAST,
	showResult,
	showTime bool
)

var rootCmd = &cobra.Command{
	Use:   "finescript",
	Short: "A simple programming language.",
	Long:  "He is fine!",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		env := runtime.GlobalEnv()

		for {
			print("> ")
			text, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}

			source := strings.TrimSpace(text)
			tokens, errs := lexer.Tokenize(source)
			if len(errs) > 0 {
				panic(strings.Join(errs, "\n"))
			}
			ast, errs := parser.Parse(tokens, source)
			if len(errs) > 0 {
				panic(strings.Join(errs, "\n"))
			}
			result := runtime.EvaluateStmt(ast, env)
			println(fmt.Sprintf("\n%s", runtime.Format(result)))
		}
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run program from file.",
	Long:  "Specify the path to your software file and enjoy!",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		startReadFile := time.Now()
		sourceBytes, err := os.ReadFile(args[0])
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}
		source := string(sourceBytes)
		durationReadFile := time.Since(startReadFile)

		startLexer := time.Now()
		tokens, errs := lexer.Tokenize(source)
		if len(errs) > 0 {
			panic(strings.Join(errs, "\n"))
		}
		durationLexer := time.Since(startLexer)

		startParser := time.Now()
		ast, errs := parser.Parse(tokens, source)
		if len(errs) > 0 {
			panic(strings.Join(errs, "\n"))
		}
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

func main() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().BoolVarP(&showTokens, "show-tokens", "t", false, "Enables program tokens visibility")
	runCmd.PersistentFlags().BoolVarP(&showAST, "show-ast", "a", false, "Enables program AST visibility")
	runCmd.PersistentFlags().BoolVarP(&showResult, "show-result", "r", false, "Enables program result visibility")
	runCmd.PersistentFlags().BoolVarP(&showTime, "show-time", "s", false, "Enables program execute time visibility")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
