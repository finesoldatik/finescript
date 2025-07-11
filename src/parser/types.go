package parser

import (
	"finescript/src/ast"
	"finescript/src/lexer"
	"fmt"
)

type type_nud_handler func(p *parser) ast.Type
type type_led_handler func(p *parser, left ast.Type, bp binding_power) ast.Type

type type_nud_lookup map[lexer.TokenKind]type_nud_handler
type type_led_lookup map[lexer.TokenKind]type_led_handler
type type_bp_lookup map[lexer.TokenKind]binding_power

var type_bp_lu = type_bp_lookup{}
var type_nud_lu = type_nud_lookup{}
var type_led_lu = type_led_lookup{}

func type_led(kind lexer.TokenKind, bp binding_power, led_fn type_led_handler) {
	type_bp_lu[kind] = bp
	type_led_lu[kind] = led_fn
}

func type_nud(kind lexer.TokenKind, bp binding_power, nud_fn type_nud_handler) {
	type_bp_lu[kind] = primary
	type_nud_lu[kind] = nud_fn
}

func createTypeTokenLookups() {
	type_nud(lexer.IDENTIFIER, primary, func(p *parser) ast.Type {
		token := p.advance()
		return ast.TypeAlias{
			Name:     token.Value,
			Position: token.Position,
		}
	})
	type_nud(lexer.TYPE, primary, parseTypeDecl)
	type_nud(lexer.STRUCT, primary, parseStructType)
	type_nud(lexer.NULL, primary, parsePrimaryType)
	type_nud(lexer.UNDEFINED, primary, parsePrimaryType)
	type_nud(lexer.FUN, primary, parsePrimaryType)
	type_nud(lexer.INT_TYPE, primary, parsePrimaryType)
	type_nud(lexer.FLOAT_TYPE, primary, parsePrimaryType)
	type_nud(lexer.STRING_TYPE, primary, parsePrimaryType)
	type_nud(lexer.BOOL_TYPE, primary, parsePrimaryType)
	type_nud(lexer.OBJECT_TYPE, primary, parsePrimaryType)
	type_nud(lexer.ARRAY_TYPE, primary, parsePrimaryType)
	type_nud(lexer.ANY_TYPE, primary, parsePrimaryType)
	type_nud(lexer.VOID_TYPE, primary, parsePrimaryType)
}

func parse_type(p *parser, bp binding_power) ast.Type {
	tokenKind := p.currentTokenKind()
	nud_fn, exists := type_nud_lu[tokenKind]

	if !exists {
		panic(fmt.Sprintf("type: NUD Handler expected for token %s\n", lexer.TokenKindString(tokenKind)))
	}

	left := nud_fn(p)

	for type_bp_lu[p.currentTokenKind()] > bp {
		tokenKind = p.currentTokenKind()
		led_fn, exists := type_led_lu[tokenKind]

		if !exists {
			panic(fmt.Sprintf("type: LED Handler expected for token %s\n", lexer.TokenKindString(tokenKind)))
		}

		left = led_fn(p, left, bp)
	}

	return left
}

func parseTypeDecl(p *parser) ast.Type {
	return ast.TypeAliasDecl{}
}

func parseStructType(p *parser) ast.Type {
	return ast.StructType{}
}

func parsePrimaryType(p *parser) ast.Type {
	switch p.currentTokenKind() {
	case lexer.NULL:
		return ast.NullKeyword{}
	case lexer.UNDEFINED:
		return ast.UndefinedKeyword{}
	case lexer.FUN:
		return ast.FunKeyword{}
	case lexer.INT_TYPE:
		return ast.IntKeyword{}
	case lexer.FLOAT_TYPE:
		return ast.FloatKeyword{}
	case lexer.STRING_TYPE:
		return ast.StringKeyword{}
	case lexer.BOOL_TYPE:
		return ast.BoolKeyword{}
	case lexer.OBJECT_TYPE:
		return ast.ObjectKeyword{}
	case lexer.ARRAY_TYPE:
		return ast.ArrayKeyword{}
	case lexer.ANY_TYPE:
		return ast.AnyKeyword{}
	case lexer.VOID_TYPE:
		return ast.VoidKeyword{}
	default:
		panic(fmt.Sprintf("Cannot create primary_expr from %s at %s", lexer.TokenKindString(p.currentTokenKind()), p.currentToken().Position.ToString()))
	}
}
