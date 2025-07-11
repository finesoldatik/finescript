package runtime

import (
	"finescript/src/ast"
)

func evalProgram(stmt ast.Program, env Environment) RuntimeVal {
	var lastEvaluated RuntimeVal = NullVal{}
	for _, bodyStmt := range stmt.Body {
		lastEvaluated = EvaluateStmt(bodyStmt, env)
	}
	return lastEvaluated
}

func evalBlockStmt(stmt ast.BlockStmt, env Environment) RuntimeVal {
	var lastEvaluated RuntimeVal = NullVal{}
	scope := Environment{
		parent:    &env,
		variables: make(map[string]variable),
	}

	for _, bodyStmt := range stmt.Body {
		lastEvaluated = EvaluateStmt(bodyStmt, scope)
	}
	return lastEvaluated
}

func evalIfStmt(stmt ast.IfStmt, env Environment) RuntimeVal {
	condition := ToBool(evaluateExpr(stmt.Condition, env)).Value

	if condition {
		scope := Environment{
			parent:    &env,
			variables: make(map[string]variable),
		}

		for _, consequentStmt := range stmt.Consequent {
			EvaluateStmt(consequentStmt, scope)
		}
	} else {
		scope := Environment{
			parent:    &env,
			variables: make(map[string]variable),
		}

		for _, alternateStmt := range stmt.Alternate {
			EvaluateStmt(alternateStmt, scope)
		}
	}

	return NullVal{}
}
