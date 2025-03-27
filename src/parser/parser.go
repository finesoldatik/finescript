package parser

import (
	"finescript/src/ast"
	"finescript/src/lexer"
)

type parser struct {
	tokens        []lexer.Token
	pos           int
	initialSource string
}

func createParser(tokens []lexer.Token, initialSource string) *parser {
	createTokenLookups()

	p := &parser{
		tokens:        tokens,
		pos:           0,
		initialSource: initialSource,
	}

	return p
}

func Parse(tokens []lexer.Token, initialSource string) ast.Program {
	p := createParser(tokens, initialSource)
	body := make([]ast.Stmt, 0)

	for p.hasTokens() {
		body = append(body, parseStmt(p))
	}

	return ast.Program{
		Body: body,
	}
}
