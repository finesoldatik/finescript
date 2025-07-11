package parser

import (
	"finescript/src/ast"
	"finescript/src/lexer"
	"fmt"
)

func parseParams(p *parser, params []ast.Param) {
	for !p.currentToken().IsOneOfMany(lexer.CLOSE_PAREN, lexer.EOF) {
		params = append(params, ast.Param{
			Name: p.expect(lexer.IDENTIFIER).Value,
			Type: parse_type(p, defalt_bp),
		})

		if p.currentTokenKind() == lexer.COMMA {
			p.advance()
		} else {
			panic(fmt.Sprintf("Expected ',' between parameters in function declaration at %s", p.currentToken().Position.ToString()))
		}
	}
}
