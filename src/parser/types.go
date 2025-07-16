package parser

import (
	"finescript/src/ast"
	"finescript/src/lexer"
	"fmt"
)

type typeNUDHandler func(p *parser) ast.Type
type typeLEDHandler func(p *parser, left ast.Type, bp bindingPower) ast.Type

type typeNUDLookup map[lexer.TokenKind]typeNUDHandler
type typeLEDLookup map[lexer.TokenKind]typeLEDHandler
type typeBPLookup map[lexer.TokenKind]bindingPower

var typeNUDLU = typeNUDLookup{}
var typeLEDLU = typeLEDLookup{}
var typeBPLU = typeBPLookup{}

func typeNUD(kind lexer.TokenKind, bp bindingPower, nudFn typeNUDHandler) {
	typeBPLU[kind] = bp
	typeNUDLU[kind] = nudFn
}

func typeLED(kind lexer.TokenKind, bp bindingPower, ledFn typeLEDHandler) {
	typeBPLU[kind] = bp
	typeLEDLU[kind] = ledFn
}

func createTypeTokenLookups() {
	typeNUD(lexer.IDENTIFIER, primary, func(p *parser) ast.Type {
		token := p.advance()
		return ast.TypeAlias{
			Name:     token.Value,
			Position: token.Position,
		}
	})
	typeNUD(lexer.STRUCT, primary, parseStruct)
	typeNUD(lexer.NULL, primary, parsePrimaryType)
	typeNUD(lexer.UNDEFINED, primary, parsePrimaryType)
	typeNUD(lexer.FUN, primary, parsePrimaryType)
	typeNUD(lexer.INT_TYPE, primary, parsePrimaryType)
	typeNUD(lexer.FLOAT_TYPE, primary, parsePrimaryType)
	typeNUD(lexer.STRING_TYPE, primary, parsePrimaryType)
	typeNUD(lexer.BOOL_TYPE, primary, parsePrimaryType)
	typeNUD(lexer.OBJECT_TYPE, primary, parsePrimaryType)
	typeNUD(lexer.ARRAY_TYPE, primary, parsePrimaryType)
	typeNUD(lexer.ANY_TYPE, primary, parsePrimaryType)
	typeNUD(lexer.VOID_TYPE, primary, parsePrimaryType)
}

func parseType(p *parser, bp bindingPower) ast.Type {
	token := p.currentToken()
	nudFn, exists := typeNUDLU[token.Kind]

	if !exists {
		p.errors = append(p.errors, fmt.Sprintf("TYPE_NUD Handler expected for token %s at %s\n", lexer.TokenKindString(token.Kind), token.Position.String()))
		return ast.Error{
			Position: &token.Position,
		}
	}

	left := nudFn(p)

	for typeBPLU[p.currentTokenKind()] > bp {
		token = p.currentToken()
		ledFn, exists := typeLEDLU[token.Kind]

		if !exists {
			p.errors = append(p.errors, fmt.Sprintf("TYPE_LED Handler expected for token %s at %s\n", lexer.TokenKindString(token.Kind), token.Position.String()))
			return ast.Error{
				Position: &token.Position,
			}
		}

		left = ledFn(p, left, bp)
	}

	return left
}

func parseStruct(p *parser) ast.Type {
	startPos := p.advance().Position.StartPos
	members := make([]ast.Member, 0)
	properties := make([]ast.PropertySignature, 0)
	methods := make([]ast.MethodSignature, 0)

	expected := p.expect(lexer.OPEN_CURLY)
	if expected.Kind == lexer.ERROR {
		return ast.Error{
			Position: &expected.Position,
		}
	}
	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		expectedName := p.expect(lexer.IDENTIFIER)
		name := expectedName.Value
		if expectedName.Kind == lexer.ERROR {
			return ast.Error{
				Position: &expectedName.Position,
			}
		}
		if p.currentTokenKind() == lexer.COLON {
			p.advance()
			properties = append(properties, ast.PropertySignature{
				Name: name,
				Type: parseType(p, defaultBP),
			})
		} else {
			expectedOpenParen := p.expect(lexer.OPEN_PAREN)
			if expectedOpenParen.Kind == lexer.ERROR {
				return ast.Error{
					Position: &expectedOpenParen.Position,
				}
			}
			params, err := parseParams(p)
			if err.Position != nil {
				return ast.Error{
					Position: err.Position,
				}
			}
			expectedCloseParen := p.expect(lexer.CLOSE_PAREN)
			if expectedCloseParen.Kind == lexer.ERROR {
				return ast.Error{
					Position: &expectedCloseParen.Position,
				}
			}

			var methodType ast.Type = ast.VoidKeyword{}
			if p.currentTokenKind() == lexer.COLON {
				p.advance()
				methodType = parseType(p, defaultBP)
			}

			methods = append(methods, ast.MethodSignature{
				Name:   name,
				Params: params,
				Type:   methodType,
			})
		}

		if p.currentTokenKind() != lexer.CLOSE_CURLY {
			expected := p.expectError(lexer.COMMA, fmt.Sprintf("Expected ',' between properties in structure declaration at %s", p.currentToken().Position.String()))
			return ast.Error{
				Position: &expected.Position,
			}
		}
	}

	for _, p := range properties {
		members = append(members, p)
	}
	for _, m := range methods {
		members = append(members, m)
	}

	expectedCloseCurly := p.expect(lexer.CLOSE_CURLY)
	if expectedCloseCurly.Kind == lexer.ERROR {
		return ast.Error{
			Position: &expectedCloseCurly.Position,
		}
	}

	return ast.Struct{
		Members: members,
		Position: lexer.Position{
			StartPos: startPos,
			EndPos:   expectedCloseCurly.Position.EndPos,
		},
	}
}

func parsePrimaryType(p *parser) ast.Type {
	token := p.currentToken()
	switch token.Kind {
	case lexer.NULL:
		p.advance()
		return ast.NullKeyword{}
	case lexer.UNDEFINED:
		p.advance()
		return ast.UndefinedKeyword{}
	case lexer.FUN:
		p.advance()
		return ast.FunKeyword{}
	case lexer.INT_TYPE:
		p.advance()
		return ast.IntKeyword{}
	case lexer.FLOAT_TYPE:
		p.advance()
		return ast.FloatKeyword{}
	case lexer.STRING_TYPE:
		p.advance()
		return ast.StringKeyword{}
	case lexer.BOOL_TYPE:
		p.advance()
		return ast.BoolKeyword{}
	case lexer.OBJECT_TYPE:
		p.advance()
		return ast.ObjectKeyword{}
	case lexer.ARRAY_TYPE:
		p.advance()
		return ast.ArrayKeyword{}
	case lexer.ANY_TYPE:
		p.advance()
		return ast.AnyKeyword{}
	case lexer.VOID_TYPE:
		p.advance()
		return ast.VoidKeyword{}
	default:
		p.errors = append(p.errors, fmt.Sprintf("Cannot create primary_expr from %s at %s", lexer.TokenKindString(token.Kind), token.Position.String()))
		return ast.Error{
			Position: &token.Position,
		}
	}
}
