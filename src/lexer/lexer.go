package lexer

import (
	"finescript/src/helpers"
	"fmt"
	"regexp"
)

var escapeRegex = regexp.MustCompile(`\\(?:[nrt\\'"]|x[0-9a-fA-F]{2}|u[0-9a-fA-F]{4}|U[0-9a-fA-F]{8})`)

type regexHandler func(lex *lexer, regex *regexp.Regexp)
type regexPattern struct {
	regex   *regexp.Regexp
	handler regexHandler
}

type lexer struct {
	patterns []regexPattern
	Tokens   []Token
	source   string
	pos      int
	errors   []string
}

func Tokenize(source string) ([]Token, []string) {
	lex := createLexer(source)

	for !lex.at_eof() {
		matched := false
		remainder := lex.remainder()

		for _, pattern := range lex.patterns {
			loc := pattern.regex.FindStringIndex(remainder)
			if loc != nil && loc[0] == 0 {
				pattern.handler(lex, pattern.regex)
				matched = true
				break
			}
		}

		if !matched {
			lex.errors = append(lex.errors, fmt.Sprintf("unrecognized token near \"%v\" at %d", helpers.Ellipsis(remainder, 20), lex.pos))
			lex.advanceN(len(remainder))
		}
	}

	lex.push(Token{EOF, "eof", Position{
		StartPos: lex.pos,
		EndPos:   lex.pos,
	}})
	return lex.Tokens, lex.errors
}

func (lex *lexer) advanceN(n int) int {
	lex.pos += n
	return lex.pos
}

// func (lex *lexer) at() byte {
// 	return lex.source[lex.pos]
// }

// func (lex *lexer) advance() {
// 	lex.pos += 1
// }

func (lex *lexer) remainder() string {
	return lex.source[lex.pos:]
}

func (lex *lexer) push(token Token) {
	lex.Tokens = append(lex.Tokens, token)
}

func (lex *lexer) at_eof() bool {
	return lex.pos >= len(lex.source)
}

func createLexer(source string) *lexer {
	return &lexer{
		pos:    0,
		errors: make([]string, 0),
		source: source,
		Tokens: make([]Token, 0),
		patterns: []regexPattern{
			{regexp.MustCompile(`\/\/.*|\/\*[\s\S]*?\*\/`), skipHandler},
			{regexp.MustCompile(`\s+`), skipHandler},
			{regexp.MustCompile(`^\s*$`), skipHandler},
			{regexp.MustCompile(`"(?:\\.|[^"\\])*"|'(?:\\.|[^'\\])*'`), stringHandler},
			{regexp.MustCompile(`[0-9]+\.[0-9]+`), numberHandler(FLOAT)},
			{regexp.MustCompile(`[0-9]+`), numberHandler(INT)},
			{regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`), identifierHandler},
			{regexp.MustCompile(`\[`), defaultHandler(OPEN_BRACKET, "[")},
			{regexp.MustCompile(`\]`), defaultHandler(CLOSE_BRACKET, "]")},
			{regexp.MustCompile(`\{`), defaultHandler(OPEN_CURLY, "{")},
			{regexp.MustCompile(`\}`), defaultHandler(CLOSE_CURLY, "}")},
			{regexp.MustCompile(`\(`), defaultHandler(OPEN_PAREN, "(")},
			{regexp.MustCompile(`\)`), defaultHandler(CLOSE_PAREN, ")")},
			{regexp.MustCompile(`==`), defaultHandler(EQUALS, "==")},
			{regexp.MustCompile(`!=`), defaultHandler(NOT_EQUALS, "!=")},
			{regexp.MustCompile(`=`), defaultHandler(ASSIGNMENT, "=")},
			{regexp.MustCompile(`!`), defaultHandler(NOT, "!")},
			{regexp.MustCompile(`<=`), defaultHandler(LESS_EQUALS, "<=")},
			{regexp.MustCompile(`<`), defaultHandler(LESS, "<")},
			{regexp.MustCompile(`>=`), defaultHandler(GREATER_EQUALS, ">=")},
			{regexp.MustCompile(`>`), defaultHandler(GREATER, ">")},
			{regexp.MustCompile(`\|\|`), defaultHandler(OR, "||")},
			{regexp.MustCompile(`&&`), defaultHandler(AND, "&&")},
			{regexp.MustCompile(`\.\.`), defaultHandler(DOT_DOT, "..")},
			{regexp.MustCompile(`\.`), defaultHandler(DOT, ".")},
			{regexp.MustCompile(`;`), defaultHandler(SEMI_COLON, ";")},
			{regexp.MustCompile(`:`), defaultHandler(COLON, ":")},
			{regexp.MustCompile(`\?`), defaultHandler(QUESTION, "?")},
			{regexp.MustCompile(`,`), defaultHandler(COMMA, ",")},
			{regexp.MustCompile(`\+\+`), defaultHandler(PLUS_PLUS, "++")},
			{regexp.MustCompile(`--`), defaultHandler(MINUS_MINUS, "--")},
			{regexp.MustCompile(`\+=`), defaultHandler(PLUS_EQUALS, "+=")},
			{regexp.MustCompile(`-=`), defaultHandler(MINUS_EQUALS, "-=")},
			// {regexp.MustCompile(`/=`), defaultHandler(SLASH_EQUALS, "/=")},
			// {regexp.MustCompile(`\*=`), defaultHandler(STAR_EQUALS, "*=")},
			// {regexp.MustCompile(`%=`), defaultHandler(PERCENT_EQUALS, "%=")},
			{regexp.MustCompile(`\+`), defaultHandler(PLUS, "+")},
			{regexp.MustCompile(`-`), defaultHandler(MINUS, "-")},
			{regexp.MustCompile(`/`), defaultHandler(SLASH, "/")},
			{regexp.MustCompile(`\*`), defaultHandler(STAR, "*")},
			{regexp.MustCompile(`%`), defaultHandler(PERCENT, "%")},
		},
	}
}

func defaultHandler(kind TokenKind, value string) regexHandler {
	return func(lex *lexer, _ *regexp.Regexp) {
		lex.push(Token{
			Kind:  kind,
			Value: value,
			Position: Position{
				StartPos: lex.pos,
				EndPos:   lex.advanceN(len(value)),
			},
		})
	}
}

func stringHandler(lex *lexer, regex *regexp.Regexp) {
	startPos := lex.pos
	match := regex.FindStringIndex(lex.remainder())
	if match == nil {
		lex.errors = append(lex.errors, fmt.Sprintf("unterminated string near \"%v\" at %d", lex.remainder(), lex.pos))
	}

	stringWithQuotes := lex.remainder()[match[0]:match[1]]

	stringLiteral := stringWithQuotes[1 : len(stringWithQuotes)-1]
	stringLiteral = escapeRegex.ReplaceAllStringFunc(stringLiteral, helpers.RemoveEscapeSigns)

	lex.push(Token{
		Kind:  STRING,
		Value: stringLiteral,
		Position: Position{
			StartPos: startPos,
			EndPos:   lex.advanceN(len(stringWithQuotes)),
		},
	})
}

func numberHandler(kind TokenKind) func(*lexer, *regexp.Regexp) {
	return func(lex *lexer, regex *regexp.Regexp) {
		startPos := lex.pos
		match := regex.FindString(lex.remainder())
		lex.push(Token{
			Kind:  kind,
			Value: match,
			Position: Position{
				StartPos: startPos,
				EndPos:   lex.advanceN(len(match)),
			},
		})
	}
}

func identifierHandler(lex *lexer, regex *regexp.Regexp) {
	startPos := lex.pos
	match := regex.FindString(lex.remainder())

	tokenKind := IDENTIFIER
	if kind, found := keywords[match]; found {
		tokenKind = kind
	}

	lex.push(Token{
		Kind:  tokenKind,
		Value: match,
		Position: Position{
			StartPos: startPos,
			EndPos:   lex.advanceN(len(match)),
		},
	})
}

func skipHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	if match == "" {
		return
	}
	lex.advanceN(len(match))
}
