package runtime

import (
	"finescript/src/ast"
	"finescript/src/lexer"
)

func evalLogicalOperations(leftVal RuntimeVal, rightVal RuntimeVal, Op lexer.Token) RuntimeVal {
	switch Op.Kind {
	case lexer.AND:
		return BoolVal{
			Value: ToBool(leftVal).Value && ToBool(rightVal).Value,
		}
	case lexer.OR:
		return BoolVal{
			Value: ToBool(leftVal).Value || ToBool(rightVal).Value,
		}
	default:
		panic("Unknown Binary Operator")
	}
}

func evalComparisonOperations(leftVal RuntimeVal, rightVal RuntimeVal, Op lexer.Token) RuntimeVal {
	switch Op.Kind {
	case lexer.EQUALS:
		return BoolVal{
			Value: ToBool(leftVal).Value == ToBool(rightVal).Value,
		}
	case lexer.NOT_EQUALS:
		return BoolVal{
			Value: ToBool(leftVal).Value != ToBool(rightVal).Value,
		}
	case lexer.LESS:
		return BoolVal{
			Value: ToFloat(leftVal).Value < ToFloat(rightVal).Value,
		}
	case lexer.GREATER:
		return BoolVal{
			Value: ToFloat(leftVal).Value > ToFloat(rightVal).Value,
		}
	case lexer.LESS_EQUALS:
		return BoolVal{
			Value: ToFloat(leftVal).Value <= ToFloat(rightVal).Value,
		}
	case lexer.GREATER_EQUALS:
		return BoolVal{
			Value: ToFloat(leftVal).Value >= ToFloat(rightVal).Value,
		}
	default:
		return evalLogicalOperations(leftVal, rightVal, Op)
	}
}

func evalArithmetiсOperations(leftVal RuntimeVal, rightVal RuntimeVal, Op lexer.Token) RuntimeVal {
	switch Op.Kind {
	case lexer.PLUS:
		switch leftType := leftVal.(type) {
		case IntVal:
			return FloatVal{
				Value: ToFloat(leftType).Value + ToFloat(rightVal).Value,
			}
		case FloatVal:
			return FloatVal{
				Value: leftType.Value + ToFloat(rightVal).Value,
			}
		case StringVal:
			return StringVal{
				Value: leftType.Value + ToString(rightVal).Value,
			}
		case BoolVal:
			return IntVal{
				Value: ToInt(leftType).Value + ToInt(rightVal).Value,
			}
		default:
			panic("These types cannot be added to each other.")
		}

	case lexer.MINUS:
		switch leftType := leftVal.(type) {
		case IntVal:
			return FloatVal{
				Value: ToFloat(leftType).Value - ToFloat(rightVal).Value,
			}
		case FloatVal:
			return FloatVal{
				Value: leftType.Value - ToFloat(rightVal).Value,
			}
		default:
			panic("These types cannot be subtracted from each other.")
		}

	case lexer.STAR:
		switch leftType := leftVal.(type) {
		case IntVal:
			return FloatVal{
				Value: ToFloat(leftType).Value * ToFloat(rightVal).Value,
			}
		case FloatVal:
			return FloatVal{
				Value: leftType.Value * ToFloat(rightVal).Value,
			}
		case StringVal:
			switch rightType := rightVal.(type) {
			case IntVal:
				result := ""
				for range rightType.Value {
					result += leftType.Value
				}
				return StringVal{
					Value: result,
				}
			default:
				panic("These types cannot be multiplied by each other.")
			}
		default:
			panic("These types cannot be multiplied by each other.")
		}

	case lexer.SLASH:
		switch leftType := leftVal.(type) {
		case IntVal:
			return FloatVal{
				Value: ToFloat(leftType).Value / ToFloat(rightVal).Value,
			}
		case FloatVal:
			return FloatVal{
				Value: leftType.Value / ToFloat(rightVal).Value,
			}
		default:
			panic("These types cannot be a divided from each other.")
		}

	case lexer.PERCENT:
		switch leftType := leftVal.(type) {
		case IntVal:
			return IntVal{
				Value: leftType.Value % ToInt(rightVal).Value,
			}
		case FloatVal:
			return IntVal{
				Value: int64(leftType.Value) % ToInt(rightVal).Value,
			}
		default:
			panic("These types cannot be separated from each other using a remainder.")
		}
	default:
		return evalComparisonOperations(leftVal, rightVal, Op)
	}
}

func evalBinaryExpr(expr ast.BinaryExpr, env Environment) RuntimeVal {
	leftVal := evaluateExpr(expr.Left, env)
	rightVal := evaluateExpr(expr.Right, env)

	return evalArithmetiсOperations(leftVal, rightVal, expr.Op)
}

func evalUnaryExpr(expr ast.UnaryExpr, env Environment) RuntimeVal {
	value := evaluateExpr(expr.Expr, env)
	switch expr.Op.Kind {
	case lexer.MINUS:
		return FloatVal{
			Value: -ToFloat(value).Value,
		}
	case lexer.NOT:
		return BoolVal{
			Value: !ToBool(value).Value,
		}
	case lexer.PLUS_PLUS:
		var result RuntimeVal = FloatVal{
			Value: ToFloat(value).Value + 1,
		}
		switch ident := expr.Expr.(type) {
		case ast.Identifier:
			env.assignVar(ident.Name, result)
		}
		return result
	case lexer.MINUS_MINUS:
		var result RuntimeVal = FloatVal{
			Value: ToFloat(value).Value - 1,
		}
		switch ident := expr.Expr.(type) {
		case ast.Identifier:
			env.assignVar(ident.Name, result)
		}
		return result
	default:
		panic("Unknown Unary Operator")
	}
}

func evalCallExpr(expr ast.CallExpr, env Environment) RuntimeVal {
	var args []RuntimeVal
	for _, arg := range expr.Args {
		args = append(args, evaluateExpr(arg, env))
	}
	caller := evaluateExpr(expr.Caller, env)

	switch callerType := caller.(type) {
	case NativeFnVal:
		return callerType.Call(args, env)
	case FunctionVal:
		scope := Environment{
			parent:    &callerType.DeclarationEnv,
			variables: make(map[string]variable),
		}

		for i, param := range callerType.Params {
			if len(callerType.Params) > len(args) {
				panic("Arg num less then Param num")
			} else if len(callerType.Params) < len(args) {
				panic("Arg num more then Param num")
			}
			scope.declareVar(param.Name, args[i], false)
		}

		var result RuntimeVal = NullVal{}
		for _, stmt := range callerType.Body {
			result = EvaluateStmt(stmt, scope)
		}

		return result
	default:
		panic("Cannot call value that is not a function")
	}
}

func evalAssignExpr(expr ast.AssignExpr, env Environment) RuntimeVal {
	switch left := expr.Assigne.(type) {
	case ast.Identifier:
		right := evaluateExpr(expr.Expr, env)
		return env.assignVar(left.Name, right)
	default:
		panic("Invalid left hand side expr inside assignment expr")
	}
}
