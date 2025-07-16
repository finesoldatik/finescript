package parser

import (
	"finescript/src/ast"
	"finescript/src/lexer"
	"fmt"
)

type parser struct {
	tokens        []lexer.Token
	pos           int
	initialSource string
	errors        []string
}

func newParser(tokens []lexer.Token, initialSource string) *parser {
	createTokenLookups()
	createTypeTokenLookups()

	p := &parser{
		tokens:        tokens,
		pos:           0,
		initialSource: initialSource,
		errors:        make([]string, 0),
	}

	return p
}

func Parse(tokens []lexer.Token, initialSource string) (ast.Program, []string) {
	p := newParser(tokens, initialSource)
	body := make([]ast.Stmt, 0)

	for p.hasTokens() {
		body = append(body, parseStmt(p))
	}

	return ast.Program{
			Body: body,
			Position: lexer.Position{
				StartPos: 0,
				EndPos:   body[len(body)-1].Pos().EndPos,
			},
		},
		p.errors
}

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

func (p *parser) expectError(expectedKind lexer.TokenKind, err any) lexer.Token {
	token := p.currentToken()
	if token.Kind != expectedKind {
		if err == nil {
			p.errors = append(p.errors, fmt.Sprintf("Syntax error: expected %s but got %s (\"%s\") at %s:\n%s",
				lexer.TokenKindString(expectedKind),
				lexer.TokenKindString(token.Kind),
				token.Value,
				token.Position.String(),
				p.initialSource[token.Position.StartPos:token.Position.EndPos],
			))
			return lexer.Token{
				Kind:     lexer.ERROR,
				Value:    token.Position.String(),
				Position: token.Position,
			}
		} else {
			p.errors = append(p.errors, p.error(err, &token.Position))
			return lexer.Token{
				Kind:     lexer.ERROR,
				Value:    token.Position.String(),
				Position: token.Position,
			}
		}
	}
	return p.advance()
}

func (p *parser) expect(expectedKind lexer.TokenKind) lexer.Token {
	return p.expectError(expectedKind, nil)
}

func (p *parser) error(err any, pos *lexer.Position) string {
	if pos == nil {
		tokenPos := p.currentToken().Position
		return fmt.Sprintf("Parser Error at %s:\n%s\n%s", tokenPos.String(), p.initialSource[tokenPos.StartPos:tokenPos.EndPos], err)
	} else {
		return fmt.Sprintf("Parser Error at %s:\n%s\n%s", pos.String(), p.initialSource[pos.StartPos:pos.EndPos], err)
	}
}
