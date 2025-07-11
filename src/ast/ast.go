package ast

import (
	"finescript/src/helpers"
	"finescript/src/lexer"
)

type Stmt interface {
	stmt()
	Pos() lexer.Position
}

type Expr interface {
	expr()
	Pos() lexer.Position
}

type Type interface {
	_type()
	Pos() lexer.Position
}

func ExpectExpr[T Expr](expr Expr) T {
	return helpers.ExpectType[T](expr)
}

func ExpectStmt[T Stmt](stmt Stmt) T {
	return helpers.ExpectType[T](stmt)
}
