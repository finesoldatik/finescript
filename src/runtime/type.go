package runtime

import (
	"finescript/src/ast"
)

func resolveType(typ ast.Type, env Environment) ast.Type {
	switch t := typ.(type) {

	case ast.TypeAlias:
		val := env.lookupVar(t.Name).Value
		typeAlias, ok := val.(TypeAliasVal)
		if !ok {
			// p.errors = append(p.errors, fmt.Sprintf("Expected type alias, got something else at %s", typeAlias.Type.Pos().String()))
			// return ast.Error{
			// 	Position: typeAlias.Type.Pos(),
			// }
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

func inferType(val RuntimeVal) []ast.Type {
	types := make([]ast.Type, 0)
	switch r := val.(type) {
	case IntVal:
		types = append(types, ast.IntKeyword{}, ast.IntLiteralType{Type: r.Value})
	case FloatVal:
		types = append(types, ast.FloatKeyword{}, ast.FloatLiteralType{Type: r.Value})
	case StringVal:
		types = append(types, ast.StringKeyword{}, ast.StringLiteralType{Type: r.Value})
	case BoolVal:
		types = append(types, ast.BoolKeyword{}, ast.BoolLiteralType{Type: r.Value})
	case NullVal:
		types = append(types, ast.NullKeyword{})
	case UndefinedVal:
		types = append(types, ast.UndefinedKeyword{})
	case FunctionVal:
		types = append(types, ast.FunKeyword{}, ast.FunType{
			Params:     r.Params,
			ReturnType: r.ReturnType,
		})
	case NativeFnVal:
		types = append(types, ast.FunKeyword{})
	case TypeAliasVal:
		types = append(types, ast.TypeAlias{Name: r.Name})
	}

	return types
}
