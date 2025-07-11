package parser

import (
	"finescript/src/ast"
	"finescript/src/lexer"
)

type binding_power int

const (
	defalt_bp binding_power = iota
	assignment
	logical
	relational
	additive
	multiplicative
	unary
	call
	member
	primary
)

type stmt_handler func(p *parser) ast.Stmt
type nud_handler func(p *parser) ast.Expr
type led_handler func(p *parser, left ast.Expr, bp binding_power) ast.Expr

type stmt_lookup map[lexer.TokenKind]stmt_handler
type nud_lookup map[lexer.TokenKind]nud_handler
type led_lookup map[lexer.TokenKind]led_handler
type bp_lookup map[lexer.TokenKind]binding_power

var bp_lu = bp_lookup{}
var nud_lu = nud_lookup{}
var led_lu = led_lookup{}
var stmt_lu = stmt_lookup{}

func led(kind lexer.TokenKind, bp binding_power, led_fn led_handler) {
	bp_lu[kind] = bp
	led_lu[kind] = led_fn
}

func nud(kind lexer.TokenKind, nud_fn nud_handler) {
	bp_lu[kind] = primary
	nud_lu[kind] = nud_fn
}

func stmt(kind lexer.TokenKind, stmt_fn stmt_handler) {
	bp_lu[kind] = defalt_bp
	stmt_lu[kind] = stmt_fn
}

func createTokenLookups() {
	// Assignment
	led(lexer.ASSIGNMENT, assignment, parseAssignExpr)
	led(lexer.PLUS_EQUALS, assignment, parseAssignExpr)
	led(lexer.MINUS_EQUALS, assignment, parseAssignExpr)

	// Logical
	led(lexer.AND, logical, parseBinaryExpr)
	led(lexer.OR, logical, parseBinaryExpr)

	// Relational
	led(lexer.LESS, relational, parseBinaryExpr)
	led(lexer.LESS_EQUALS, relational, parseBinaryExpr)
	led(lexer.GREATER, relational, parseBinaryExpr)
	led(lexer.GREATER_EQUALS, relational, parseBinaryExpr)
	led(lexer.EQUALS, relational, parseBinaryExpr)
	led(lexer.NOT_EQUALS, relational, parseBinaryExpr)

	// Additive & Multiplicitave
	led(lexer.PLUS, additive, parseBinaryExpr)
	led(lexer.MINUS, additive, parseBinaryExpr)
	led(lexer.SLASH, multiplicative, parseBinaryExpr)
	led(lexer.STAR, multiplicative, parseBinaryExpr)
	led(lexer.PERCENT, multiplicative, parseBinaryExpr)

	// Literals & Symbols
	nud(lexer.INT, parsePrimaryExpr)
	nud(lexer.FLOAT, parsePrimaryExpr)
	nud(lexer.STRING, parsePrimaryExpr)
	nud(lexer.IDENTIFIER, parsePrimaryExpr)
	nud(lexer.TRUE, parsePrimaryExpr)
	nud(lexer.FALSE, parsePrimaryExpr)

	// Unary/Prefix
	nud(lexer.MINUS, parseUnaryExpr)
	nud(lexer.NOT, parseUnaryExpr)
	nud(lexer.PLUS_PLUS, parseUnaryExpr)
	nud(lexer.MINUS_MINUS, parseUnaryExpr)
	led(lexer.PLUS_PLUS, unary, parseLedUnaryExpr)
	led(lexer.MINUS_MINUS, unary, parseLedUnaryExpr)
	// nud(lexer.OPEN_BRACKET, parseArrayLiteralExpr)

	// Member / Computed // Call
	// led(lexer.DOT, member, parseMemberExpr)
	// led(lexer.OPEN_BRACKET, member, parseMemberExpr)
	led(lexer.OPEN_PAREN, call, parseCallExpr)

	// Grouping Expr
	nud(lexer.OPEN_PAREN, parseGroupingExpr)

	// Stmt
	stmt(lexer.OPEN_CURLY, parseBlockStmt)
	stmt(lexer.LET, parseVarDeclStmt)
	stmt(lexer.VAR, parseVarDeclStmt)
	stmt(lexer.CONST, parseVarDeclStmt)
	stmt(lexer.FUN, parseFunDeclaration)
	stmt(lexer.IF, parseIfStmt)
}
