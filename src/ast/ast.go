package ast

import (
	"finescript/src/helpers"
	"finescript/src/lexer"
)

// Утверждение
type Stmt interface {
	stmt()
	Position() *lexer.Position
}

// Выражение
type Expr interface {
	expr()
	Position() *lexer.Position
}

func ExpectExpr[T Expr](expr Expr) T {
	return helpers.ExpectType[T](expr)
}

func ExpectStmt[T Stmt](stmt Stmt) T {
	return helpers.ExpectType[T](stmt)
}
