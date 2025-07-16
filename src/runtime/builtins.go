package runtime

import (
	"bufio"
	"finescript/src/lexer"
	"finescript/src/parser"
	"fmt"
	"os"
	"strings"
)

func handleArgs(argsCount int, paramCount int) {
	if argsCount > paramCount {
		panic("There are more arguments than parameters for the function.")
	} else if argsCount < paramCount {
		panic("there are fewer arguments than parameters for the function.")
	}
}

func Sprintf(args []RuntimeVal, env Environment) RuntimeVal {
	return StringVal{
		Value: sprint_format(args),
	}
}

func Print(args []RuntimeVal, env Environment) RuntimeVal {
	print(sprint_format(args))
	return NullVal{}
}

func Println(args []RuntimeVal, env Environment) RuntimeVal {
	println(sprint_format(args))
	return NullVal{}
}

func Int(args []RuntimeVal, env Environment) RuntimeVal {
	handleArgs(len(args), 1)
	return ToInt(args[0])
}

func Float(args []RuntimeVal, env Environment) RuntimeVal {
	handleArgs(len(args), 1)
	return ToFloat(args[0])
}

func String(args []RuntimeVal, env Environment) RuntimeVal {
	handleArgs(len(args), 1)
	return ToString(args[0])
}

func Bool(args []RuntimeVal, env Environment) RuntimeVal {
	handleArgs(len(args), 1)
	return ToBool(args[0])
}

func Input(args []RuntimeVal, env Environment) RuntimeVal {
	if len(args) == 0 {
		args = []RuntimeVal{
			StringVal{
				Value: "",
			},
		}
	}
	handleArgs(len(args), 1)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(ToString(args[0]).Value)
	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return StringVal{
		Value: strings.TrimSpace(text),
	}
}

func Eval(args []RuntimeVal, env Environment) RuntimeVal {
	handleArgs(len(args), 1)
	if _, ok := args[0].(StringVal); ok {
		tokens, errs := lexer.Tokenize(args[0].(StringVal).Value)
		if len(errs) > 0 {
			panic(strings.Join(errs, "\n"))
		}
		ast, errs := parser.Parse(tokens, args[0].(StringVal).Value)
		if len(errs) > 0 {
			panic(strings.Join(errs, "\n"))
		}
		result := EvaluateStmt(ast, GlobalEnv())
		return result
	}

	panic("String required for eval function")
}

// func nativeLen(args []RuntimeVal, env Environment) RuntimeVal {
// 	handleArgs(len(args), 1)
// 	return IntVal{
// 		Value: int64(len(args[0].(ArrayVal).Elements)),
// 	}
// }

func Exit(args []RuntimeVal, env Environment) RuntimeVal {
	os.Exit(0)
	return NullVal{}
}
