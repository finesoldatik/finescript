package parser

import (
	"finescript/src/ast"
	"finescript/src/lexer"
)

type bindingPower int

const (
	defaultBP bindingPower = iota
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

type stmtHandler func(p *parser) ast.Stmt
type NUDHandler func(p *parser) ast.Expr
type LEDHandler func(p *parser, left ast.Expr, bp bindingPower) ast.Expr

type stmtLookup map[lexer.TokenKind]stmtHandler
type NUDLookup map[lexer.TokenKind]NUDHandler
type LEDLookup map[lexer.TokenKind]LEDHandler
type BPLookup map[lexer.TokenKind]bindingPower

var stmtLU = stmtLookup{}
var nudLU = NUDLookup{}
var ledLU = LEDLookup{}
var bpLU = BPLookup{}

func stmt(kind lexer.TokenKind, stmt_fn stmtHandler) {
	bpLU[kind] = defaultBP
	stmtLU[kind] = stmt_fn
}

func NUD(kind lexer.TokenKind, nudFn NUDHandler) {
	bpLU[kind] = primary
	nudLU[kind] = nudFn
}

func LED(kind lexer.TokenKind, bp bindingPower, ledFn LEDHandler) {
	bpLU[kind] = bp
	ledLU[kind] = ledFn
}

func createTokenLookups() {
	// Assignment
	LED(lexer.ASSIGNMENT, assignment, parseAssignExpr)
	LED(lexer.PLUS_EQUALS, assignment, parseAssignExpr)
	LED(lexer.MINUS_EQUALS, assignment, parseAssignExpr)

	// Logical
	LED(lexer.AND, logical, parseBinaryExpr)
	LED(lexer.OR, logical, parseBinaryExpr)

	// Relational
	LED(lexer.LESS, relational, parseBinaryExpr)
	LED(lexer.LESS_EQUALS, relational, parseBinaryExpr)
	LED(lexer.GREATER, relational, parseBinaryExpr)
	LED(lexer.GREATER_EQUALS, relational, parseBinaryExpr)
	LED(lexer.EQUALS, relational, parseBinaryExpr)
	LED(lexer.NOT_EQUALS, relational, parseBinaryExpr)

	// Additive & Multiplicitave
	LED(lexer.PLUS, additive, parseBinaryExpr)
	LED(lexer.MINUS, additive, parseBinaryExpr)
	LED(lexer.SLASH, multiplicative, parseBinaryExpr)
	LED(lexer.STAR, multiplicative, parseBinaryExpr)
	LED(lexer.PERCENT, multiplicative, parseBinaryExpr)

	// Literals & Symbols
	NUD(lexer.INT, parsePrimaryExpr)
	NUD(lexer.FLOAT, parsePrimaryExpr)
	NUD(lexer.STRING, parsePrimaryExpr)
	NUD(lexer.IDENTIFIER, parsePrimaryExpr)
	NUD(lexer.TRUE, parsePrimaryExpr)
	NUD(lexer.FALSE, parsePrimaryExpr)

	// Unary/Prefix
	NUD(lexer.MINUS, parseUnaryExpr)
	NUD(lexer.NOT, parseUnaryExpr)
	NUD(lexer.PLUS_PLUS, parseUnaryExpr)
	NUD(lexer.MINUS_MINUS, parseUnaryExpr)
	LED(lexer.PLUS_PLUS, unary, parseLedUnaryExpr)
	LED(lexer.MINUS_MINUS, unary, parseLedUnaryExpr)
	// NUD(lexer.OPEN_BRACKET, parseArrayLiteralExpr)

	// Member / Computed // Call
	// LED(lexer.DOT, member, parseMemberExpr)
	// LED(lexer.OPEN_BRACKET, member, parseMemberExpr)
	LED(lexer.OPEN_PAREN, call, parseCallExpr)
	LED(lexer.OPEN_PAREN, logical, parseConditionalExpr)

	// Grouping Expr
	NUD(lexer.OPEN_PAREN, parseGroupingExpr)

	// Stmt
	stmt(lexer.OPEN_CURLY, parseBlockStmt)
	stmt(lexer.LET, parseVarDecl)
	stmt(lexer.VAR, parseVarDecl)
	stmt(lexer.CONST, parseVarDecl)
	stmt(lexer.FUN, parseFunDecl)
	stmt(lexer.IF, parseIfStmt)

	// Types
	stmt(lexer.TYPE, parseTypeDecl)
}
