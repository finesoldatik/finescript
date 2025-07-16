package ast

import (
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

type Error struct {
	Position *lexer.Position
}

func (e Error) stmt()  {}
func (e Error) expr()  {}
func (e Error) _type() {}
func (e Error) Pos() lexer.Position {
	return *e.Position
}
