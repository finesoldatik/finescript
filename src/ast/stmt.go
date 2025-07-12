package ast

import "finescript/src/lexer"

type Program struct {
	Body     []Stmt
	Position lexer.Position
}

func (s Program) stmt() {}
func (s Program) Pos() lexer.Position {
	return s.Position
}

type BlockStmt struct {
	Body     []Stmt
	Position lexer.Position
}

func (s BlockStmt) stmt() {}
func (s BlockStmt) Pos() lexer.Position {
	return s.Position
}

type ExprStmt struct {
	Expr     Expr
	Position lexer.Position
}

func (s ExprStmt) stmt() {}
func (s ExprStmt) Pos() lexer.Position {
	return s.Position
}

type VarDeclStmt struct {
	IsConstant bool
	Name       string
	Value      Expr // Опционально или нет зависит от IsConstant
	Position   lexer.Position
}

func (s VarDeclStmt) stmt() {}
func (s VarDeclStmt) Pos() lexer.Position {
	return s.Position
}

type Param struct {
	Name string
	Type Type
}

type FunDeclStmt struct {
	Name       string
	Params     []Param
	Body       []Stmt
	ReturnType Type
	Position   lexer.Position
}

func (s FunDeclStmt) stmt() {}
func (s FunDeclStmt) Pos() lexer.Position {
	return s.Position
}

type IfStmt struct {
	Condition  Expr
	Consequent []Stmt
	Alternate  []Stmt // Опционально
	Position   lexer.Position
}

func (s IfStmt) stmt() {}
func (s IfStmt) Pos() lexer.Position {
	return s.Position
}

type TypeAliasDecl struct {
	Name     string
	Type     Type
	Position lexer.Position
}

func (t TypeAliasDecl) stmt() {}
func (t TypeAliasDecl) Pos() lexer.Position {
	return t.Position
}
