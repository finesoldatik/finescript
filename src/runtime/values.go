package runtime

import "finescript/src/ast"

type RuntimeVal interface {
	runtime_val()
}

type IntVal struct {
	Value int64
}

func (r IntVal) runtime_val() {}

type FloatVal struct {
	Value float64
}

func (r FloatVal) runtime_val() {}

type StringVal struct {
	Value string
}

func (r StringVal) runtime_val() {}

type BoolVal struct {
	Value bool
}

func (r BoolVal) runtime_val() {}

type NullVal struct{}

func (r NullVal) runtime_val() {}

// type ArrayVal struct {
// 	Elements []RuntimeVal
// }

// func (r ArrayVal) runtime_val() {}

// type ObjectVal struct {
// 	Elements map[string]RuntimeVal
// }

// func (r ObjectVal) runtime_val() {}

type FunctionVal struct {
	Name           string
	Params         []ast.Param
	Body           []ast.Stmt
	DeclarationEnv Environment
}

func (r FunctionVal) runtime_val() {}

type FunctionCall = func(args []RuntimeVal, env Environment) RuntimeVal

type NativeFnVal struct {
	Name string
	Call FunctionCall
}

func (r NativeFnVal) runtime_val() {}

type TypeAliasVal struct {
	Name string
	Type ast.Type
}

func (r TypeAliasVal) runtime_val() {}
