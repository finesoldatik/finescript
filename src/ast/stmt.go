package ast

import "finescript/src/lexer"

//
// Декларации
//

// Декларация переменной
type VarDeclStmt struct {
	IsConstant bool   // Константное выражение
	Name       string // Имя переменной
	Value      Expr   // Значение (опционально)
	Pos        *lexer.Position
}

func (s VarDeclStmt) stmt() {}
func (s VarDeclStmt) Position() *lexer.Position {
	return s.Pos
}

// Декларация функции
type FunctionDeclStmt struct {
	Name   string   // Имя функции
	Params []string // Список параметров
	Body   []Stmt   // Тело функции
	Pos    *lexer.Position
}

func (s FunctionDeclStmt) stmt() {}
func (s FunctionDeclStmt) Position() *lexer.Position {
	return s.Pos
}

//
// Утверждения
//

type Program struct {
	Body []Stmt // Список утверждений
}

func (s Program) stmt() {}
func (s Program) Position() *lexer.Position {
	return &lexer.Position{}
}

// Утверждение выражения
type ExprStmt struct {
	Expr Expr // Выражение
}

func (s ExprStmt) stmt() {}
func (s ExprStmt) Position() *lexer.Position {
	return &lexer.Position{}
}

// Блок утверждений
type BlockStmt struct {
	Body []Stmt // Список утверждений
	Pos  *lexer.Position
}

func (s BlockStmt) stmt() {}
func (s BlockStmt) Position() *lexer.Position {
	return s.Pos
}

// Утверждение If
type IfStmt struct {
	Condition  Expr   // Условие
	Consequent []Stmt // Если условие верно
	Alternate  []Stmt // Иначе (опционально)
	Pos        *lexer.Position
}

func (s IfStmt) stmt() {}
func (s IfStmt) Position() *lexer.Position {
	return s.Pos
}

// Цикл For
type ForStmt struct {
	Init      Stmt // Инициализация (например, объявление переменной)
	Condition Expr // Условие продолжения
	Update    Stmt // Обновление (например, инкремент)
	Body      Stmt // Тело цикла
	Pos       *lexer.Position
}

func (s ForStmt) stmt() {}
func (s ForStmt) Position() *lexer.Position {
	return s.Pos
}

// Цикл While
type WhileStmt struct {
	Condition Expr // Условие
	Body      Stmt // Тело цикла
	Pos       *lexer.Position
}

func (s WhileStmt) stmt() {}
func (s WhileStmt) Position() *lexer.Position {
	return s.Pos
}

// Цикл Loop
type LoopStmt struct {
	Body []Stmt // Тело цикла
	Pos  *lexer.Position
}

func (s LoopStmt) stmt() {}
func (s LoopStmt) Position() *lexer.Position {
	return s.Pos
}

// Прерывание цикла Break
type BreakStmt struct {
	Pos *lexer.Position
}

func (s BreakStmt) stmt() {}
func (s BreakStmt) Position() *lexer.Position {
	return s.Pos
}

// Прерывание цикла Continue
type ContinueStmt struct {
	Pos *lexer.Position
}

func (s ContinueStmt) stmt() {}
func (s ContinueStmt) Position() *lexer.Position {
	return s.Pos
}

// Утверждение возвращаемого значения
type ReturnStmt struct {
	Value Expr // Возвращаемое значение (опционально)
	Pos   *lexer.Position
}

func (s ReturnStmt) stmt() {}
func (s ReturnStmt) Position() *lexer.Position {
	return s.Pos
}

// Утверждение импорта
type ImportStmt struct {
	Path string // Путь к модулю
	Pos  *lexer.Position
}

func (s ImportStmt) stmt() {}
func (s ImportStmt) Position() *lexer.Position {
	return s.Pos
}
