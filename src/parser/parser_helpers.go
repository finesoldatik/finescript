package parser

import (
	"finescript/src/lexer"
	"fmt"

	"github.com/sanity-io/litter"
)

func (p *parser) currentToken() lexer.Token {
	return p.tokens[p.pos]
}

func (p *parser) advance() lexer.Token {
	tk := p.currentToken()
	p.pos++
	return tk
}

func (p *parser) hasTokens() bool {
	return p.pos < len(p.tokens) && p.currentTokenKind() != lexer.EOF
}

// func (p *parser) nextToken() lexer.Token {
// 	return p.tokens[p.pos+1]
// }

// func (p *parser) previousToken() lexer.Token {
// 	return p.tokens[p.pos-1]
// }

func (p *parser) currentTokenKind() lexer.TokenKind {
	return p.currentToken().Kind
}

func (p *parser) expectError(expectedKind lexer.TokenKind, pos *lexer.Position, err any) lexer.Token {
	token := p.currentToken()
	if token.Kind != expectedKind {
		if err == nil {
			err = fmt.Sprintf("Syntax error: expected %s but got %s ('%s')",
				lexer.TokenKindString(expectedKind),
				lexer.TokenKindString(token.Kind),
				token.Value)
			litter.Dump(pos)
			panic(lexer.FormatError(p.initialSource, pos, fmt.Sprintf("%s", err)))
		}
	}
	return p.advance()
}

func (p *parser) expect(expectedKind lexer.TokenKind, pos *lexer.Position) lexer.Token {
	return p.expectError(expectedKind, pos, nil)
}
