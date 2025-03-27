package ast

import "finescript/src/lexer"

//
// Литералы
//

// Идентификатор
type Identifier struct {
	Name string // identifier
	Pos  *lexer.Position
}

func (e Identifier) expr() {}
func (e Identifier) Position() *lexer.Position {
	return e.Pos
}

// Литерал целого числа
type IntLiteral struct {
	Value int64 // 0
	Pos   *lexer.Position
}

func (e IntLiteral) expr() {}
func (e IntLiteral) Position() *lexer.Position {
	return e.Pos
}

// Литерал числа с плавающей точкой
type FloatLiteral struct {
	Value float64 // 0.0
	Pos   *lexer.Position
}

func (e FloatLiteral) expr() {}
func (e FloatLiteral) Position() *lexer.Position {
	return e.Pos
}

// Литерал строки
type StringLiteral struct {
	Value string // "", ''
	Pos   *lexer.Position
}

func (e StringLiteral) expr() {}
func (e StringLiteral) Position() *lexer.Position {
	return e.Pos
}

// Литерал булевого значения
type BoolLiteral struct {
	Value bool // true, false
	Pos   *lexer.Position
}

func (e BoolLiteral) expr() {}
func (e BoolLiteral) Position() *lexer.Position {
	return e.Pos
}

// Литерал списка
type ArrayLiteral struct {
	Elements []Expr // []
	Pos      *lexer.Position
}

func (e ArrayLiteral) expr() {}
func (e ArrayLiteral) Position() *lexer.Position {
	return e.Pos
}

// Литерал объекта
type ObjectLiteral struct {
	Fields map[string]Expr // Object{}
	Pos    *lexer.Position
}

func (e ObjectLiteral) expr() {}
func (e ObjectLiteral) Position() *lexer.Position {
	return e.Pos
}

type NullLiteral struct {
	Pos *lexer.Position
}

func (e NullLiteral) expr() {}
func (e NullLiteral) Position() *lexer.Position {
	return e.Pos
}

//
// Выражения
//

// Унарное выражение
type UnaryExpr struct {
	Op    lexer.Token // "!", "-", "++", "--"
	Value Expr        // Значение
	Pos   *lexer.Position
}

func (e UnaryExpr) expr() {}
func (e UnaryExpr) Position() *lexer.Position {
	return e.Pos
}

// Бинарное выражение
type BinaryExpr struct {
	Left  Expr        // Левое выражение
	Op    lexer.Token // "+", "-"
	Right Expr        // Правое выражение
	Pos   *lexer.Position
}

func (e BinaryExpr) expr() {}
func (e BinaryExpr) Position() *lexer.Position {
	return e.Pos
}

// Выражение вызова
type CallExpr struct {
	Caller Expr   // Идентификатор функции
	Args   []Expr // Список аргументов
	Pos    *lexer.Position
}

func (e CallExpr) expr() {}
func (e CallExpr) Position() *lexer.Position {
	return e.Pos
}

// Тернарное выражение
type TernaryExpr struct {
	Condition Expr // Условие
	Then      Expr // Если условие верно
	Else      Expr // Иначе (опционально)
	Pos       *lexer.Position
}

func (e TernaryExpr) expr() {}
func (e TernaryExpr) Position() *lexer.Position {
	return e.Pos
}

// Выражение члена
type MemberExpr struct {
	Object Expr // Объект
	Field  Expr // Поле объекта
	Pos    *lexer.Position
}

func (e MemberExpr) expr() {}
func (e MemberExpr) Position() *lexer.Position {
	return e.Pos
}

// Выражение присваивания
type AssignExpr struct {
	Left  Expr // Присвоить
	Right Expr // Присваиваемое
	Pos   *lexer.Position
}

func (e AssignExpr) expr() {}
func (e AssignExpr) Position() *lexer.Position {
	return e.Pos
}

// Try-Catch выражение
type TryCatchExpr struct {
	Try   []Stmt // Попытаться
	Catch []Stmt // Иначе
	Pos   *lexer.Position
}

func (e TryCatchExpr) expr() {}
func (e TryCatchExpr) Position() *lexer.Position {
	return e.Pos
}
