package lexer

import "fmt"

type TokenKind int

const (
	EOF TokenKind = iota
	NULL
	UNDEFINED
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
	ASSIGNMENT
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

	PLUS_EQUALS
	MINUS_EQUALS
	STAR_EQUALS
	SLASH_EQUALS
	PERCENT_EQUALS

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
	"array": ARRAY_TYPE,
	"any": ANY_TYPE,
	"void": VOID_TYPE,
}

type Position struct {
	StartPos int
	EndPos   int
}

func (pos Position) ToString() string {
	return fmt.Sprintf("%d:%d", pos.StartPos, pos.EndPos)
}

type Token struct {
	Kind     TokenKind
	Value    string
	Position Position
}

func (tk Token) IsOneOfMany(expectedTokens ...TokenKind) bool {
	for _, expected := range expectedTokens {
		if expected == tk.Kind {
			return true
		}
	}

	return false
}

func (token Token) Debug() {
	fmt.Printf("%s()\n", TokenKindString(token.Kind))
}

func TokenKindString(kind TokenKind) string {
	switch kind {
	case EOF:
		return "eof"
	case NULL:
		return "null"
	case INT:
		return "int"
	case FLOAT:
		return "float"
	case STRING:
		return "string"
	case TRUE:
		return "true"
	case FALSE:
		return "false"
	case IDENTIFIER:
		return "identifier"
	case OPEN_BRACKET:
		return "open_bracket"
	case CLOSE_BRACKET:
		return "close_bracket"
	case OPEN_CURLY:
		return "open_curly"
	case CLOSE_CURLY:
		return "close_curly"
	case OPEN_PAREN:
		return "open_paren"
	case CLOSE_PAREN:
		return "close_paren"
	case ASSIGNMENT:
		return "assignment"
	case EQUALS:
		return "equals"
	case NOT_EQUALS:
		return "not_equals"
	case NOT:
		return "not"
	case LESS:
		return "less"
	case LESS_EQUALS:
		return "less_equals"
	case GREATER:
		return "greater"
	case GREATER_EQUALS:
		return "greater_equals"
	case OR:
		return "or"
	case AND:
		return "and"
	case DOT:
		return "dot"
	case DOT_DOT:
		return "dot_dot"
	case SEMI_COLON:
		return "semi_colon"
	case COLON:
		return "colon"
	case QUESTION:
		return "question"
	case COMMA:
		return "comma"
	case PLUS_PLUS:
		return "plus_plus"
	case MINUS_MINUS:
		return "minus_minus"
	case PLUS_EQUALS:
		return "plus_equals"
	case MINUS_EQUALS:
		return "minus_equals"
	case PLUS:
		return "plus"
	case MINUS:
		return "minus"
	case SLASH:
		return "slash"
	case STAR:
		return "star"
	case PERCENT:
		return "percent"
	case LET:
		return "let"
	case VAR:
		return "var"
	case CONST:
		return "const"
	case FUN:
		return "fun"
	case IF:
		return "if"
	case ELSE:
		return "else"
	case YAY:
		return "yay"
	case OOPS:
		return "oops"
	default:
		return fmt.Sprintf("unknown(%d)", kind)
	}
}
