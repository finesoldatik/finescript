package parser

import (
	"finescript/src/ast"
	"finescript/src/lexer"
	"fmt"
)

func parseParams(p *parser) []ast.Param {
	params := make([]ast.Param, 0)
	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_PAREN {
		name := p.expect(lexer.IDENTIFIER).Value
		p.expect(lexer.COLON)

		params = append(params, ast.Param{
			Name: name,
			Type: parse_type(p, defalt_bp),
		})

		if p.currentTokenKind() != lexer.CLOSE_PAREN {
			p.expectError(lexer.COMMA, fmt.Sprintf("Expected ',' between parameters in function declaration at %s", p.currentToken().Position.ToString()))
		}
	}

	return params
}
