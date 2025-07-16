package lexer

import (
	"fmt"
	// "slices"
)

type TokenKind int

const (
	// Специальные
	EOF TokenKind = iota
	NULL
	UNDEFINED

	// Литералы
	TRUE
	FALSE
	INT
	FLOAT
	STRING
	IDENTIFIER

	// Скобки
	OPEN_BRACKET
	CLOSE_BRACKET
	OPEN_CURLY
	CLOSE_CURLY
	OPEN_PAREN
	CLOSE_PAREN

	// Равенство
	EQUALS
	NOT_EQUALS
	NOT

	// Измерение
	LESS
	LESS_EQUALS
	GREATER
	GREATER_EQUALS

	// Логическое
	OR
	AND

	// Символы
	DOT
	DOT_DOT
	SEMI_COLON
	COLON
	QUESTION
	COMMA

	// Краткая запись
	PLUS_PLUS
	MINUS_MINUS

	// Присваивание
	ASSIGNMENT
	PLUS_EQUALS
	MINUS_EQUALS
	// STAR_EQUALS
	// SLASH_EQUALS
	// PERCENT_EQUALS

	// Математика
	PLUS
	MINUS
	SLASH
	STAR
	PERCENT

	// Ключевые слова
	LET
	VAR
	CONST

	TYPE
	STRUCT

	FUN
	IF
	ELSE

	YAY
	OOPS

	INT_TYPE
	FLOAT_TYPE
	STRING_TYPE
	BOOL_TYPE
	OBJECT_TYPE
	ARRAY_TYPE
	ANY_TYPE
	VOID_TYPE

	ERROR
)

var keywords = map[string]TokenKind{
	"true":      TRUE,
	"false":     FALSE,
	"null":      NULL,
	"undefined": UNDEFINED,

	"let":   LET,
	"var":   VAR,
	"const": CONST,

	"type":   TYPE,
	"struct": STRUCT,

	"fun":  FUN,
	"if":   IF,
	"else": ELSE,

	"yay":  YAY,
	"oops": OOPS,

	"int":    INT_TYPE,
	"float":  FLOAT_TYPE,
	"string": STRING_TYPE,
	"bool":   BOOL_TYPE,
	"object": OBJECT_TYPE,
	"array":  ARRAY_TYPE,
	"any":    ANY_TYPE,
	"void":   VOID_TYPE,
}

type Position struct {
	StartPos int
	EndPos   int
}

func (pos Position) String() string {
	return fmt.Sprintf("%d:%d", pos.StartPos, pos.EndPos)
}

type Token struct {
	Kind     TokenKind
	Value    string
	Position Position
}

// func (tk Token) IsOneOf(kinds ...TokenKind) bool {
// 	return slices.Contains(kinds, tk.Kind)
// }
// // Для проверки категорий токенов
// func (k TokenKind) IsLiteral() bool {
//     return k >= TRUE && k <= IDENTIFIER
// }

// func (k TokenKind) IsOperator() bool {
//     return k >= EQUALS && k <= PERCENT
// }

var tokenKindString = map[TokenKind]string{
	EOF:       "eof",
	NULL:      "null",
	UNDEFINED: "undefined",

	INT:        "int",
	FLOAT:      "float",
	STRING:     "string",
	TRUE:       "true",
	FALSE:      "false",
	IDENTIFIER: "identifier",

	OPEN_BRACKET:  "open_bracket",
	CLOSE_BRACKET: "close_bracket",
	OPEN_CURLY:    "open_curly",
	CLOSE_CURLY:   "close_curly",
	OPEN_PAREN:    "open_paren",
	CLOSE_PAREN:   "close_paren",

	EQUALS:     "equals",
	NOT_EQUALS: "not_equals",
	NOT:        "not",

	LESS:           "less",
	LESS_EQUALS:    "less_equals",
	GREATER:        "greater",
	GREATER_EQUALS: "greater_equals",

	OR:  "or",
	AND: "and",

	DOT:        "dot",
	DOT_DOT:    "dot_dot",
	SEMI_COLON: "semi_colon",
	COLON:      "colon",
	QUESTION:   "question",
	COMMA:      "comma",

	PLUS_PLUS:   "plus_plus",
	MINUS_MINUS: "minus_minus",

	ASSIGNMENT:   "assignment",
	PLUS_EQUALS:  "plus_equals",
	MINUS_EQUALS: "minus_equals",

	PLUS:    "plus",
	MINUS:   "minus",
	SLASH:   "slash",
	STAR:    "star",
	PERCENT: "percent",

	LET:   "let",
	VAR:   "var",
	CONST: "const",

	FUN:  "fun",
	IF:   "if",
	ELSE: "else",

	YAY:  "yay",
	OOPS: "oops",

	TYPE:   "type",
	STRUCT: "struct",

	INT_TYPE:    "int_type",
	FLOAT_TYPE:  "float_type",
	STRING_TYPE: "string_type",
	BOOL_TYPE:   "bool_type",
	OBJECT_TYPE: "object_type",
	ARRAY_TYPE:  "array_type",
	ANY_TYPE:    "any_type",
	VOID_TYPE:   "void_type",
}

func TokenKindString(kind TokenKind) string {
	if str, ok := tokenKindString[kind]; ok {
		return str
	}
	return fmt.Sprintf("unknown(%d)", kind)
}
