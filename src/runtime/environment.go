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
		Call: nativePrint,
	}, true)

	env.declareVar("println", NativeFnVal{
		Name: "println",
		Call: nativePrintln,
	}, true)

	env.declareVar("sprintf", NativeFnVal{
		Name: "sprintf",
		Call: nativeSprintf,
	}, true)

	env.declareVar("int", NativeFnVal{
		Name: "int",
		Call: nativeInt,
	}, true)

	env.declareVar("float", NativeFnVal{
		Name: "float",
		Call: nativeFloat,
	}, true)

	env.declareVar("string", NativeFnVal{
		Name: "string",
		Call: nativeString,
	}, true)

	env.declareVar("bool", NativeFnVal{
		Name: "bool",
		Call: nativeBool,
	}, true)

	env.declareVar("input", NativeFnVal{
		Name: "input",
		Call: nativeInput,
	}, true)

	// env.declareVar("len", NativeFnVal{
	// 	Name: "len",
	// 	Call: nativeLen,
	// }, true)

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

func (env Environment) declareVar(varname string, value RuntimeVal, isConstant bool) RuntimeVal {
	if _, exists := env.variables[varname]; exists {
		panic("Variable exists")
	}

	env.variables[varname] = variable{
		IsConstant: isConstant,
		Value:      value,
	}

	return value
}

func (env Environment) assignVar(varname string, value RuntimeVal) RuntimeVal {
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

func (env Environment) lookupVar(varname string) variable {
	newEnv := env.resolve(varname)
	return newEnv.variables[varname]
}

func (env Environment) resolve(varname string) Environment {
	if _, exists := env.variables[varname]; exists {
		return env
	}

	if env.parent == nil {
		panic(fmt.Sprintf("Cannot resolve \"%s\" as it does not exist.", varname))
	}

	return env.parent.resolve(varname)
}
