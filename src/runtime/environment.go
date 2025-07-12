package runtime

import (
	"fmt"
)

func GlobalEnv() Environment {
	env := Environment{
		parent:    nil,
		variables: make(map[string]variable),
	}

	env.declareVar("print", NativeFnVal{
		Name: "print",
		Call: Print,
	}, true)

	env.declareVar("println", NativeFnVal{
		Name: "println",
		Call: Println,
	}, true)

	env.declareVar("sprintf", NativeFnVal{
		Name: "sprintf",
		Call: Sprintf,
	}, true)

	env.declareVar("int", NativeFnVal{
		Name: "int",
		Call: Int,
	}, true)

	env.declareVar("float", NativeFnVal{
		Name: "float",
		Call: Float,
	}, true)

	env.declareVar("string", NativeFnVal{
		Name: "string",
		Call: String,
	}, true)

	env.declareVar("bool", NativeFnVal{
		Name: "bool",
		Call: Bool,
	}, true)

	env.declareVar("input", NativeFnVal{
		Name: "input",
		Call: Input,
	}, true)

	env.declareVar("eval", NativeFnVal{
		Name: "eval",
		Call: Eval,
	}, true)

	// env.declareVar("len", NativeFnVal{
	// 	Name: "len",
	// 	Call: nativeLen,
	// }, true)

	env.declareVar("exit", NativeFnVal{
		Name: "exit",
		Call: Exit,
	}, true)

	return env
}

type variable struct {
	IsConstant bool
	Value      RuntimeVal
}

type Environment struct {
	parent    *Environment
	variables map[string]variable
}

func (env *Environment) declareVar(varname string, value RuntimeVal, isConstant bool) RuntimeVal {
	if _, exists := env.variables[varname]; exists {
		panic("Variable exists")
	}

	env.variables[varname] = variable{
		IsConstant: isConstant,
		Value:      value,
	}

	return value
}

func (env *Environment) assignVar(varname string, value RuntimeVal) RuntimeVal {
	newEnv := env.resolve(varname)

	if newEnv.variables[varname].IsConstant {
		panic(fmt.Sprintf("Cannot reasign to variable \"%s\" as it was declared constant.", varname))
	}

	newEnv.variables[varname] = variable{
		IsConstant: false,
		Value:      value,
	}

	return value
}

func (env *Environment) lookupVar(varname string) variable {
	newEnv := env.resolve(varname)
	return newEnv.variables[varname]
}

func (env *Environment) resolve(varname string) *Environment {
	if _, exists := env.variables[varname]; exists {
		return env
	}

	if env.parent == nil {
		panic(fmt.Sprintf("Cannot resolve \"%s\" as it does not exist.", varname))
	}

	return env.parent.resolve(varname)
}
