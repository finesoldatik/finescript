package ast

import "finescript/src/lexer"

type StringLiteralType struct {
	Type     string
	Position lexer.Position
}

func (t StringLiteralType) _type() {}
func (t StringLiteralType) Pos() lexer.Position {
	return t.Position
}

type IntLiteralType struct {
	Type     int64
	Position lexer.Position
}

func (t IntLiteralType) _type() {}
func (t IntLiteralType) Pos() lexer.Position {
	return t.Position
}

type FloatLiteralType struct {
	Type     float64
	Position lexer.Position
}

func (t FloatLiteralType) _type() {}
func (t FloatLiteralType) Pos() lexer.Position {
	return t.Position
}

type BoolLiteralType struct {
	Type     bool
	Position lexer.Position
}

func (t BoolLiteralType) _type() {}
func (t BoolLiteralType) Pos() lexer.Position {
	return t.Position
}

type IntKeyword struct {
	Position lexer.Position
}

func (t IntKeyword) _type() {}
func (t IntKeyword) Pos() lexer.Position {
	return t.Position
}

type FloatKeyword struct {
	Position lexer.Position
}

func (t FloatKeyword) _type() {}
func (t FloatKeyword) Pos() lexer.Position {
	return t.Position
}

type StringKeyword struct {
	Position lexer.Position
}

func (t StringKeyword) _type() {}
func (t StringKeyword) Pos() lexer.Position {
	return t.Position
}

type BoolKeyword struct {
	Position lexer.Position
}

func (t BoolKeyword) _type() {}
func (t BoolKeyword) Pos() lexer.Position {
	return t.Position
}

type NullKeyword struct {
	Position lexer.Position
}

func (t NullKeyword) _type() {}
func (t NullKeyword) Pos() lexer.Position {
	return t.Position
}

type UndefinedKeyword struct {
	Position lexer.Position
}

func (t UndefinedKeyword) _type() {}
func (t UndefinedKeyword) Pos() lexer.Position {
	return t.Position
}

type ObjectKeyword struct {
	Position lexer.Position
}

func (t ObjectKeyword) _type() {}
func (t ObjectKeyword) Pos() lexer.Position {
	return t.Position
}

type ArrayKeyword struct {
	Position lexer.Position
}

func (t ArrayKeyword) _type() {}
func (t ArrayKeyword) Pos() lexer.Position {
	return t.Position
}

type AnyKeyword struct {
	Position lexer.Position
}

func (t AnyKeyword) _type() {}
func (t AnyKeyword) Pos() lexer.Position {
	return t.Position
}

type VoidKeyword struct {
	Position lexer.Position
}

func (t VoidKeyword) _type() {}
func (t VoidKeyword) Pos() lexer.Position {
	return t.Position
}

type FunKeyword struct {
	Position lexer.Position
}

func (t FunKeyword) _type() {}
func (t FunKeyword) Pos() lexer.Position {
	return t.Position
}

type TypeAlias struct {
	Name     string
	Position lexer.Position
}

func (t TypeAlias) _type() {}
func (t TypeAlias) Pos() lexer.Position {
	return t.Position
}

type TypeAliasDecl struct {
	Name     string
	Type     Type
	Position lexer.Position
}

func (t TypeAliasDecl) _type() {}
func (t TypeAliasDecl) Pos() lexer.Position {
	return t.Position
}

type ArrayType struct {
	ElementType Type
	Position    lexer.Position
}

func (t ArrayType) _type() {}
func (t ArrayType) Pos() lexer.Position {
	return t.Position
}

type UnionType struct {
	Types    []Type
	Position lexer.Position
}

func (t UnionType) _type() {}
func (t UnionType) Pos() lexer.Position {
	return t.Position
}

type IntersectionType struct {
	Types    []Type
	Position lexer.Position
}

func (t IntersectionType) _type() {}
func (t IntersectionType) Pos() lexer.Position {
	return t.Position
}

type FunType struct {
	Params   map[string]Type
	Type     Type
	Position lexer.Position
}

func (t FunType) _type() {}
func (t FunType) Pos() lexer.Position {
	return t.Position
}

type Member interface {
	member()
}

type PropertySignature struct {
	Name string
	Type Type
}

func (t PropertySignature) member() {}

type MethodSignature struct {
	Name   string
	Params map[string]Type
	Type   Type
}

func (t MethodSignature) member() {}

type StructType struct {
	Members  []Member
	Type     Type
	Position lexer.Position
}

func (t StructType) _type() {}
func (t StructType) Pos() lexer.Position {
	return t.Position
}
