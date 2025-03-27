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

// func evalLoopStmt(stmt ast.LoopStmt, env Environment) RuntimeVal {
// 	state := 0
// 	for {
// 		scope := Environment{
// 			parent:    &env,
// 			variables: make(map[string]variable),
// 		}

// 		for _, bodyStmt := range stmt.Body {
// 			switch bodyStmt.(type) {
// 			case ast.ContinueStmt:
// 				state = 1
// 			case ast.BreakStmt:
// 				state = 2
// 			default:
// 				EvaluateStmt(bodyStmt, scope)
// 			}
// 		}

// 		if state == 1 {
// 			continue
// 		} else if state == 2 {
// 			break
// 		}
// 	}

// 	return NullVal{}
// }
