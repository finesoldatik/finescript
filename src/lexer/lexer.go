package lexer

import (
	"finescript/src/helpers"
	"fmt"
	"regexp"
	"strings"
)

var escapeRegex = regexp.MustCompile(`\\(?:[nrt\\'"]|x[0-9a-fA-F]{2}|u[0-9a-fA-F]{4}|U[0-9a-fA-F]{8})`)

type regexPattern struct {
	regex   *regexp.Regexp
	handler regexHandler
}

type lexer struct {
	patterns      []regexPattern
	Tokens        []Token
	source        string
	initialSource string
	pos           int
	line          int
	lineStart     int
}

func Tokenize(source string) []Token {
	lex := createLexer(source)

	for !lex.at_eof() {
		matched := false

		for _, pattern := range lex.patterns {
			loc := pattern.regex.FindStringIndex(lex.remainder())
			if loc != nil && loc[0] == 0 {
				pattern.handler(lex, pattern.regex)
				matched = true
				break
			}
		}

		if !matched {
			lex.error("unrecognized token near \"%v\"", helpers.Ellipsis(lex.remainder(), 20))
		}
	}

	lex.push(newToken(EOF, "EOF", lex.currentPosition(0)))
	return lex.Tokens
}

func (lex *lexer) currentPosition(length int) *Position {
	startCol := lex.pos - lex.lineStart + 1
	return &Position{
		StartLine:   lex.line,
		StartColumn: startCol,
		EndLine:     lex.line,
		EndColumn:   startCol + length - 1,
		Index:       lex.pos,
	}
}

func (lex *lexer) advanceN(n int) {
	lex.pos += n
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

func (lex *lexer) error(format string, args ...interface{}) {
	pos := lex.currentPosition(0)
	msg := fmt.Sprintf(format, args...)
	panic(FormatError(lex.initialSource, pos, msg))
}

func createLexer(source string) *lexer {
	return &lexer{
		pos:           0,
		line:          1,
		lineStart:     0,
		source:        source,
		initialSource: source,
		Tokens:        make([]Token, 0),
		patterns: []regexPattern{
			{regexp.MustCompile(`\/\/.*|\/\*[\s\S]*?\*\/`), commentHandler},
			{regexp.MustCompile(`\s+`), skipHandler},
			{regexp.MustCompile(`^\s*$`), skipHandler},
			{regexp.MustCompile(`"(?:\\.|[^"\\])*"|'(?:\\.|[^'\\])*'`), stringHandler},
			{regexp.MustCompile(`[0-9]+\.[0-9]+`), floatHandler},
			{regexp.MustCompile(`[0-9]+`), intHandler},
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
			{regexp.MustCompile(`\*\*`), defaultHandler(STAR_STAR, "**")},
			{regexp.MustCompile(`\+=`), defaultHandler(PLUS_EQUALS, "+=")},
			{regexp.MustCompile(`-=`), defaultHandler(MINUS_EQUALS, "-=")},
			{regexp.MustCompile(`/=`), defaultHandler(SLASH_EQUALS, "/=")},
			{regexp.MustCompile(`\*=`), defaultHandler(STAR_EQUALS, "*=")},
			{regexp.MustCompile(`%=`), defaultHandler(PERCENT_EQUALS, "%=")},
			{regexp.MustCompile(`\+`), defaultHandler(PLUS, "+")},
			{regexp.MustCompile(`-`), defaultHandler(MINUS, "-")},
			{regexp.MustCompile(`/`), defaultHandler(SLASH, "/")},
			{regexp.MustCompile(`\*`), defaultHandler(STAR, "*")},
			{regexp.MustCompile(`%`), defaultHandler(PERCENT, "%")},
		},
	}
}

type regexHandler func(lex *lexer, regex *regexp.Regexp)

func defaultHandler(kind TokenKind, value string) regexHandler {
	return func(lex *lexer, _ *regexp.Regexp) {
		pos := lex.currentPosition(len(value))
		lex.advanceN(len(value))
		lex.push(newToken(kind, value, pos))
	}
}

func stringHandler(lex *lexer, regex *regexp.Regexp) {
	startPos := lex.currentPosition(0)

	match := regex.FindStringIndex(lex.remainder())
	if match == nil {
		lex.error("unterminated string near '%v'", lex.remainder())
	}

	stringWithQuotes := lex.remainder()[match[0]:match[1]]

	lineCount := strings.Count(stringWithQuotes, "\n")
	endLine := lex.line + lineCount
	lastNewline := strings.LastIndex(stringWithQuotes, "\n")
	var endColumn = len(stringWithQuotes) - lastNewline - 1 // -1 для учета кавычки

	if lineCount == 0 {
		endColumn = startPos.StartColumn + len(stringWithQuotes) - 1
	} else {
		endColumn = len(stringWithQuotes[lastNewline+1:])
	}

	pos := &Position{
		StartLine:   lex.line,
		StartColumn: startPos.StartColumn,
		EndLine:     endLine,
		EndColumn:   endColumn,
		Index:       startPos.Index,
	}

	lex.handleNewlines(stringWithQuotes)

	stringLiteral := stringWithQuotes[1 : len(stringWithQuotes)-1]
	stringLiteral = escapeRegex.ReplaceAllStringFunc(stringLiteral, helpers.RemoveEscapeSigns)

	lex.push(newToken(STRING, stringLiteral, pos))
	lex.advanceN(len(stringWithQuotes))
}

func floatHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.push(newToken(FLOAT, match, lex.currentPosition(len(match))))
	lex.advanceN(len(match))
}

func intHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.push(newToken(INT, match, lex.currentPosition(len(match))))
	lex.advanceN(len(match))
}

func identifierHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())

	if kind, found := reserved_lu[match]; found {
		lex.push(newToken(kind, match, lex.currentPosition(len(match))))
	} else {
		lex.push(newToken(IDENTIFIER, match, lex.currentPosition(len(match))))
	}

	lex.advanceN(len(match))
}

func skipHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	if match == "" {
		return
	}
	lex.handleNewlines(match)
	lex.advanceN(len(match))
}

func commentHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	if match == "" {
		return
	}
	lex.handleNewlines(match)
	lex.advanceN(len(match)) // Пропускаем комментарий полностью
}

func (lex *lexer) handleNewlines(s string) {
	for i, c := range s {
		if c == '\n' {
			lex.line++
			lex.lineStart = lex.pos + i + 1
		}
	}
}

func FormatError(source string, pos *Position, message string) string {
	lines := strings.Split(source, "\n")

	errorLines := lines[pos.StartLine-1 : pos.EndLine]

	var builder strings.Builder

	builder.WriteString("\n=== Syntax Error ===\n")
	builder.WriteString(fmt.Sprintf("Message: %s\n", message))
	builder.WriteString(fmt.Sprintf("Location: %s\n\n", pos.String()))

	for i := pos.StartLine; i <= pos.EndLine; i++ {
		// Выделяем строку с ошибкой
		if i == pos.EndLine {
			builder.WriteString(fmt.Sprintf("> %4d | %s\n", i, errorLines[i-pos.StartLine]))

			if pos.StartColumn > 0 {
				underline := strings.Repeat(" ", pos.StartColumn+6) + "^"
				if pos.EndColumn > pos.StartColumn {
					underline += strings.Repeat("~", pos.EndColumn-pos.StartColumn)
				}
				builder.WriteString(fmt.Sprintf("       | %s\n", underline))
			}
		} else {
			builder.WriteString(fmt.Sprintf("  %4d | %s\n", i, errorLines[i-pos.StartLine]))
		}
	}

	return builder.String()
}
