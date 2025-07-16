package runtime

import (
	"finescript/src/ast"
	"fmt"
)

func EvaluateStmt(node ast.Stmt, env Environment) RuntimeVal {
	switch stmt := node.(type) {
	case ast.Program:
		return evalProgram(stmt, env)
	case ast.BlockStmt:
		return evalBlockStmt(stmt, env)
	case ast.VarDeclStmt:
		return env.declareVar(stmt.Name, evaluateExpr(stmt.Value, env), stmt.IsConstant)
	case ast.FunDeclStmt:
		return env.declareVar(stmt.Name, FunctionVal{
			Name:           stmt.Name,
			Params:         stmt.Params,
			Body:           stmt.Body,
			ReturnType:     stmt.ReturnType,
			DeclarationEnv: env,
		}, true)
	case ast.TypeAliasDecl:
		resolveType(stmt.Type, env)
		return env.declareVar(stmt.Name, TypeAliasVal{
			Name: stmt.Name,
			Type: resolveType(stmt.Type, env),
		}, true)
	case ast.IfStmt:
		return evalIfStmt(stmt, env)
	case ast.ExprStmt:
		return evaluateExpr(stmt.Expr, env)
	default:
		panic("Unknown Stmt")
	}
}

func evaluateExpr(node ast.Expr, env Environment) RuntimeVal {
	switch expr := node.(type) {
	case ast.Identifier:
		return env.lookupVar(expr.Name).Value
	case ast.IntLiteral:
		return IntVal{
			Value: expr.Value,
		}
	case ast.FloatLiteral:
		return FloatVal{
			Value: expr.Value,
		}
	case ast.StringLiteral:
		return StringVal{
			Value: expr.Value,
		}
	case ast.BoolLiteral:
		return BoolVal{
			Value: expr.Value,
		}
	// case ast.ArrayLiteral:
	// 	result := make([]RuntimeVal, 0)
	// 	for _, elem := range expr.Elements {
	// 		result = append(result, evaluateExpr(elem, env))
	// 	}
	// 	return ArrayVal{
	// 		Elements: result,
	// 	}
	case ast.BinaryExpr:
		return evalBinaryExpr(expr, env)
	case ast.UnaryExpr:
		return evalUnaryExpr(expr, env)
	case ast.AssignExpr:
		return evalAssignExpr(expr, env)
	case ast.CallExpr:
		return evalCallExpr(expr, env)
	default:
		panic(fmt.Sprintf("Unknown Expr at %s", expr.Pos().String()))
	}
}
