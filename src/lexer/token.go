package lexer

import (
	"fmt"
)

type TokenKind int

const (
	EOF TokenKind = iota
	NULL
	TRUE
	FALSE
	INT
	FLOAT
	STRING
	IDENTIFIER

	// Grouping & Braces
	OPEN_BRACKET
	CLOSE_BRACKET
	OPEN_CURLY
	CLOSE_CURLY
	OPEN_PAREN
	CLOSE_PAREN

	// Equivilance
	ASSIGNMENT
	EQUALS
	NOT_EQUALS
	NOT

	// Conditional
	LESS
	LESS_EQUALS
	GREATER
	GREATER_EQUALS

	// Logical
	OR
	AND

	// Symbols
	DOT
	DOT_DOT
	SEMI_COLON
	COLON
	QUESTION
	COMMA

	// Shorthand
	PLUS_PLUS
	MINUS_MINUS
	STAR_STAR

	PLUS_EQUALS
	MINUS_EQUALS
	STAR_EQUALS
	SLASH_EQUALS
	PERCENT_EQUALS

	// Maths
	PLUS
	MINUS
	SLASH
	STAR
	PERCENT

	// Reserved Keywords
	LET
	VAR
	CONST
	IMPORT
	FROM
	FUN
	LOOP
	BREAK
	CONTINUE
	IF
	ELSE
	EXPORT
)

var reserved_lu map[string]TokenKind = map[string]TokenKind{
	"true":     TRUE,
	"false":    FALSE,
	"null":     NULL,
	"let":      LET,
	"var":      VAR,
	"const":    CONST,
	"import":   IMPORT,
	"from":     FROM,
	"fun":      FUN,
	"loop":     LOOP,
	"break":    BREAK,
	"continue": CONTINUE,
	"if":       IF,
	"else":     ELSE,
	"export":   EXPORT,
}

type Position struct {
	StartLine   int
	EndLine     int
	StartColumn int
	EndColumn   int
	Index       int
}

func (p Position) String() string {
	if p.StartLine == p.EndLine {
		return fmt.Sprintf("line %d, cols %d-%d", p.StartLine, p.StartColumn, p.EndColumn)
	}
	return fmt.Sprintf("line %d, col %d - line %d, col %d", p.StartLine, p.StartColumn, p.EndLine, p.EndColumn)
}

type Token struct {
	Kind  TokenKind
	Value string
	Pos   *Position
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
	if token.Kind == IDENTIFIER || token.Kind == INT || token.Kind == FLOAT || token.Kind == STRING {
		fmt.Printf("%s(%s)\n", TokenKindString(token.Kind), token.Value)
	} else {
		fmt.Printf("%s()\n", TokenKindString(token.Kind))
	}
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
		return "VAR"
	case CONST:
		return "const"
	case IMPORT:
		return "import"
	case FROM:
		return "from"
	case FUN:
		return "fun"
	case LOOP:
		return "loop"
	case BREAK:
		return "break"
	case CONTINUE:
		return "continue"
	case IF:
		return "if"
	case ELSE:
		return "else"
	case EXPORT:
		return "export"
	default:
		return fmt.Sprintf("unknown(%d)", kind)
	}
}

func newToken(kind TokenKind, value string, pos *Position) Token {
	return Token{
		Kind:  kind,
		Value: value,
		Pos:   pos,
	}
}
