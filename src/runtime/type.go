package runtime

import (
	"finescript/src/ast"
	"finescript/src/lexer"
)

func resolveType(typ ast.Type, env Environment) ast.Type {
	switch t := typ.(type) {

	case ast.TypeAlias:
		val := env.lookupVar(t.Name).Value
		typeAlias, ok := val.(TypeAliasVal)
		if !ok {
			panic("Expected type alias, got something else")
		}
		return resolveType(typeAlias.Type, env)

	case ast.ArrayType:
		return ast.ArrayType{
			ElementType: resolveType(t.ElementType, env),
		}

	case ast.UnionType:
		resolved := make([]ast.Type, 0, len(t.Types))
		for _, inner := range t.Types {
			resolved = append(resolved, resolveType(inner, env))
		}
		return ast.UnionType{Types: resolved}

	case ast.IntersectionType:
		resolved := make([]ast.Type, 0, len(t.Types))
		for _, inner := range t.Types {
			resolved = append(resolved, resolveType(inner, env))
		}
		return ast.IntersectionType{Types: resolved}

	case ast.FunType:
		params := make([]ast.Param, 0, len(t.Params))
		for _, p := range t.Params {
			params = append(params, ast.Param{
				Name: p.Name,
				Type: resolveType(p.Type, env),
			})
		}
		return ast.FunType{
			Params:     params,
			ReturnType: resolveType(t.ReturnType, env),
		}

	case ast.Struct:
		members := make([]ast.Member, 0, len(t.Members))
		for _, m := range t.Members {
			switch member := m.(type) {
			case ast.PropertySignature:
				members = append(members, ast.PropertySignature{
					Name: member.Name,
					Type: resolveType(member.Type, env),
				})
			case ast.MethodSignature:
				params := make([]ast.Param, 0, len(member.Params))
				for _, p := range member.Params {
					params = append(params, ast.Param{
						Name: p.Name,
						Type: resolveType(p.Type, env),
					})
				}
				members = append(members, ast.MethodSignature{
					Name:   member.Name,
					Params: params,
					Type:   resolveType(member.Type, env),
				})
			default:
				panic("Unknown struct member type")
			}
		}
		return ast.Struct{
			Members:  members,
			Position: t.Position,
		}

	// Примитивные и literal-типы возвращаем как есть
	case ast.IntKeyword, ast.StringKeyword, ast.FloatKeyword,
		ast.BoolKeyword, ast.NullKeyword, ast.UndefinedKeyword,
		ast.VoidKeyword, ast.AnyKeyword, ast.ObjectKeyword, ast.ArrayKeyword,
		ast.FunKeyword,
		ast.StringLiteralType, ast.IntLiteralType, ast.FloatLiteralType, ast.BoolLiteralType:
		return t

	default:
		panic("resolveType: Unknown type variant")
	}
}

func InferType(val RuntimeVal) ast.Type {
	switch val.(type) {
	case IntVal:
		return ast.IntKeyword{
			Position: lexer.Position{},
		}
	case FloatVal:
		return ast.FloatKeyword{
			Position: lexer.Position{},
		}
	case StringVal:
		return ast.StringKeyword{
			Position: lexer.Position{},
		}
	case BoolVal:
		return ast.BoolKeyword{
			Position: lexer.Position{},
		}
	case NullVal:
		return ast.NullKeyword{
			Position: lexer.Position{},
		}
	case FunctionVal, NativeFnVal:
		return ast.FunKeyword{
			Position: lexer.Position{},
		}
	default:
		panic("InferType: unsupported RuntimeVal")
	}
}
