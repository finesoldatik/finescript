package ast

import "finescript/src/lexer"

type Identifier struct {
	Name     string
	Position lexer.Position
}

func (e Identifier) expr() {}
func (e Identifier) Pos() lexer.Position {
	return e.Position
}

type IntLiteral struct {
	Value    int64
	Position lexer.Position
}

func (e IntLiteral) expr() {}
func (e IntLiteral) Pos() lexer.Position {
	return e.Position
}

type FloatLiteral struct {
	Value    float64
	Position lexer.Position
}

func (e FloatLiteral) expr() {}
func (e FloatLiteral) Pos() lexer.Position {
	return e.Position
}

type StringLiteral struct {
	Value    string
	Position lexer.Position
}

func (e StringLiteral) expr() {}
func (e StringLiteral) Pos() lexer.Position {
	return e.Position
}

type BoolLiteral struct {
	Value    bool
	Position lexer.Position
}

func (e BoolLiteral) expr() {}
func (e BoolLiteral) Pos() lexer.Position {
	return e.Position
}

type NullLiteral struct {
	Position lexer.Position
}

func (e NullLiteral) expr() {}
func (e NullLiteral) Pos() lexer.Position {
	return e.Position
}

type UndefinedLiteral struct {
	Position lexer.Position
}

func (e UndefinedLiteral) expr() {}
func (e UndefinedLiteral) Pos() lexer.Position {
	return e.Position
}

//////

type UnaryExpr struct {
	Op       lexer.Token
	Expr     Expr
	Position lexer.Position
}

func (e UnaryExpr) expr() {}
func (e UnaryExpr) Pos() lexer.Position {
	return e.Position
}

type BinaryExpr struct {
	Left     Expr
	Op       lexer.Token
	Right    Expr
	Position lexer.Position
}

func (e BinaryExpr) expr() {}
func (e BinaryExpr) Pos() lexer.Position {
	return e.Position
}

type AssignExpr struct {
	Assigne  Expr
	Expr     Expr
	Position lexer.Position
}

func (e AssignExpr) expr() {}
func (e AssignExpr) Pos() lexer.Position {
	return e.Position
}

type CallExpr struct {
	Caller   Expr
	Args     []Expr
	Position lexer.Position
}

func (e CallExpr) expr() {}
func (e CallExpr) Pos() lexer.Position {
	return e.Position
}

type ConditionalExpr struct {
	Condition   Expr
	Consequent  Expr
	Alternate   Expr
	Position    lexer.Position
}

func (e ConditionalExpr) expr() {}
func (e ConditionalExpr) Pos() lexer.Position {
	return e.Position
}
