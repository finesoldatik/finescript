package parser

import (
	"finescript/src/ast"
	"finescript/src/lexer"
)

type binding_power int

const (
	defalt_bp binding_power = iota
	comma
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

func nud(kind lexer.TokenKind, bp binding_power, nud_fn nud_handler) {
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
	led(lexer.STAR_STAR, multiplicative, parseBinaryExpr)

	// Literals & Symbols
	nud(lexer.INT, primary, parsePrimaryExpr)
	nud(lexer.FLOAT, primary, parsePrimaryExpr)
	nud(lexer.STRING, primary, parsePrimaryExpr)
	nud(lexer.IDENTIFIER, primary, parsePrimaryExpr)
	nud(lexer.TRUE, primary, parsePrimaryExpr)
	nud(lexer.FALSE, primary, parsePrimaryExpr)

	// Unary/Prefix
	nud(lexer.MINUS, unary, parseUnaryExpr)
	nud(lexer.NOT, unary, parseUnaryExpr)
	nud(lexer.PLUS_PLUS, unary, parseUnaryExpr)
	nud(lexer.MINUS_MINUS, unary, parseUnaryExpr)
	led(lexer.PLUS_PLUS, unary, parseLedUnaryExpr)
	led(lexer.MINUS_MINUS, unary, parseLedUnaryExpr)
	nud(lexer.OPEN_BRACKET, primary, parseArrayLiteralExpr)

	// Member / Computed // Call
	led(lexer.DOT, member, parseMemberExpr)
	led(lexer.OPEN_BRACKET, member, parseMemberExpr)
	led(lexer.OPEN_PAREN, call, parseCallExpr)

	// Grouping Expr
	nud(lexer.OPEN_PAREN, defalt_bp, parseGroupingExpr)

	// Stmt
	stmt(lexer.OPEN_CURLY, parseBlockStmt)
	stmt(lexer.LET, parseVarDeclStmt)
	stmt(lexer.VAR, parseVarDeclStmt)
	stmt(lexer.CONST, parseVarDeclStmt)
	stmt(lexer.FUN, parseFunDeclaration)
	stmt(lexer.IF, parseIfStmt)
	stmt(lexer.LOOP, parseLoopStmt)
	stmt(lexer.BREAK, parseLoopControl)
	stmt(lexer.CONTINUE, parseLoopControl)
}
