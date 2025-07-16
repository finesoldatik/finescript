package ast

import "finescript/src/lexer"

/*
"string"
*/
type StringLiteralType struct {
	Type     string
	Position lexer.Position
}

func (t StringLiteralType) _type() {}
func (t StringLiteralType) Pos() lexer.Position {
	return t.Position
}

/*
42
*/
type IntLiteralType struct {
	Type     int64
	Position lexer.Position
}

func (t IntLiteralType) _type() {}
func (t IntLiteralType) Pos() lexer.Position {
	return t.Position
}

/*
42.0
*/
type FloatLiteralType struct {
	Type     float64
	Position lexer.Position
}

func (t FloatLiteralType) _type() {}
func (t FloatLiteralType) Pos() lexer.Position {
	return t.Position
}

/*
true
*/
type BoolLiteralType struct {
	Type     bool
	Position lexer.Position
}

func (t BoolLiteralType) _type() {}
func (t BoolLiteralType) Pos() lexer.Position {
	return t.Position
}

/*
int
*/
type IntKeyword struct {
	Position lexer.Position
}

func (t IntKeyword) _type() {}
func (t IntKeyword) Pos() lexer.Position {
	return t.Position
}

/*
float
*/
type FloatKeyword struct {
	Position lexer.Position
}

func (t FloatKeyword) _type() {}
func (t FloatKeyword) Pos() lexer.Position {
	return t.Position
}

/*
string
*/
type StringKeyword struct {
	Position lexer.Position
}

func (t StringKeyword) _type() {}
func (t StringKeyword) Pos() lexer.Position {
	return t.Position
}

/*
bool
*/
type BoolKeyword struct {
	Position lexer.Position
}

func (t BoolKeyword) _type() {}
func (t BoolKeyword) Pos() lexer.Position {
	return t.Position
}

/*
null
*/
type NullKeyword struct {
	Position lexer.Position
}

func (t NullKeyword) _type() {}
func (t NullKeyword) Pos() lexer.Position {
	return t.Position
}

/*
undefined
*/
type UndefinedKeyword struct {
	Position lexer.Position
}

func (t UndefinedKeyword) _type() {}
func (t UndefinedKeyword) Pos() lexer.Position {
	return t.Position
}

/*
object
*/
type ObjectKeyword struct {
	Position lexer.Position
}

func (t ObjectKeyword) _type() {}
func (t ObjectKeyword) Pos() lexer.Position {
	return t.Position
}

/*
array
*/
type ArrayKeyword struct {
	Position lexer.Position
}

func (t ArrayKeyword) _type() {}
func (t ArrayKeyword) Pos() lexer.Position {
	return t.Position
}

/*
any
*/
type AnyKeyword struct {
	Position lexer.Position
}

func (t AnyKeyword) _type() {}
func (t AnyKeyword) Pos() lexer.Position {
	return t.Position
}

/*
void
*/
type VoidKeyword struct {
	Position lexer.Position
}

func (t VoidKeyword) _type() {}
func (t VoidKeyword) Pos() lexer.Position {
	return t.Position
}

/*
fun
*/
type FunKeyword struct {
	Position lexer.Position
}

func (t FunKeyword) _type() {}
func (t FunKeyword) Pos() lexer.Position {
	return t.Position
}

/*
name
*/
type TypeAlias struct {
	Name     string
	Position lexer.Position
}

func (t TypeAlias) _type() {}
func (t TypeAlias) Pos() lexer.Position {
	return t.Position
}

/*
[]int
*/
type ArrayType struct {
	ElementType Type
	Position    lexer.Position
}

func (t ArrayType) _type() {}
func (t ArrayType) Pos() lexer.Position {
	return t.Position
}

/*
(int | string)
*/
type UnionType struct {
	Types    []Type
	Position lexer.Position
}

func (t UnionType) _type() {}
func (t UnionType) Pos() lexer.Position {
	return t.Position
}

/*
(int & 42)
*/
type IntersectionType struct {
	Types    []Type
	Position lexer.Position
}

func (t IntersectionType) _type() {}
func (t IntersectionType) Pos() lexer.Position {
	return t.Position
}

/*
fun (params) => type
*/
type FunType struct {
	Params   []Param
	ReturnType     Type
	Position lexer.Position
}

func (t FunType) _type() {}
func (t FunType) Pos() lexer.Position {
	return t.Position
}

/*
member
*/
type Member interface {
	member()
}

/*
name: type
*/
type PropertySignature struct {
	Name string
	Type Type
}

func (t PropertySignature) member() {}

/*
name(params): type
*/
type MethodSignature struct {
	Name   string
	Params []Param
	Type   Type
}

func (t MethodSignature) member() {}

/*
struct {
	name1: type,
	name2(params): type
}
*/
type Struct struct {
	Members  []Member
	Position lexer.Position
}

func (t Struct) _type() {}
func (t Struct) Pos() lexer.Position {
	return t.Position
}
