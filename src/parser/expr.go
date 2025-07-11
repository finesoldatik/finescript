package parser

import (
	"finescript/src/ast"
	"finescript/src/lexer"
	"fmt"
	"strconv"
)

func parseExpr(p *parser, bp binding_power) ast.Expr {
	tokenKind := p.currentTokenKind()
	nud_fn, exists := nud_lu[tokenKind]

	if !exists {
		panic(fmt.Sprintf("NUD Handler expected for token %s at %s\n", lexer.TokenKindString(tokenKind), p.currentToken().Position.ToString()))
	}

	left := nud_fn(p)

	for bp_lu[p.currentTokenKind()] > bp {
		tokenKind = p.currentTokenKind()
		led_fn, exists := led_lu[tokenKind]

		if !exists {
			panic(fmt.Sprintf("LED Handler expected for token %s at %s\n", lexer.TokenKindString(tokenKind), p.currentToken().Position.ToString()))
		}

		left = led_fn(p, left, bp)
	}

	if p.currentTokenKind() == lexer.SEMI_COLON {
		p.advance()
	}

	return left
}

func parsePrimaryExpr(p *parser) ast.Expr {
	switch p.currentTokenKind() {
	case lexer.INT:
		token := p.advance()
		number, _ := strconv.ParseInt(token.Value, 0, 64)
		return ast.IntLiteral{
			Value:    number,
			Position: token.Position,
		}
	case lexer.FLOAT:
		token := p.advance()
		number, _ := strconv.ParseFloat(token.Value, 64)
		return ast.FloatLiteral{
			Value:    number,
			Position: token.Position,
		}
	case lexer.STRING:
		token := p.advance()
		return ast.StringLiteral{
			Value:    token.Value,
			Position: token.Position,
		}
	case lexer.IDENTIFIER:
		token := p.advance()
		return ast.Identifier{
			Name:     token.Value,
			Position: token.Position,
		}
	case lexer.NULL:
		return ast.NullLiteral{
			Position: p.advance().Position,
		}
	case lexer.UNDEFINED:
		return ast.UndefinedLiteral{
			Position: p.advance().Position,
		}
	case lexer.TRUE:
		return ast.BoolLiteral{
			Value:    true,
			Position: p.advance().Position,
		}
	case lexer.FALSE:
		return ast.BoolLiteral{
			Value:    false,
			Position: p.advance().Position,
		}
	default:
		panic(fmt.Sprintf("Cannot create primary_expr from %s at %s", lexer.TokenKindString(p.currentTokenKind()), p.currentToken().Position.ToString()))
	}
}

func parseUnaryExpr(p *parser) ast.Expr {
	operatorToken := p.advance()
	expr := parseExpr(p, unary)

	return ast.UnaryExpr{
		Op:   operatorToken,
		Expr: expr,
		Position: lexer.Position{
			StartPos: operatorToken.Position.StartPos,
			EndPos:   expr.Pos().EndPos,
		},
	}
}

func parseLedUnaryExpr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	operatorToken := p.advance()

	return ast.UnaryExpr{
		Op:   operatorToken,
		Expr: left,
		Position: lexer.Position{
			StartPos: operatorToken.Position.StartPos,
			EndPos:   left.Pos().EndPos,
		},
	}
}

func parseAssignExpr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	p.advance()
	expr := parseExpr(p, bp)

	return ast.AssignExpr{
		Assigne: left,
		Expr:    expr,
		Position: lexer.Position{
			StartPos: left.Pos().StartPos,
			EndPos:   expr.Pos().EndPos,
		},
	}
}

func parseBinaryExpr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	operatorToken := p.advance()
	right := parseExpr(p, defalt_bp)

	return ast.BinaryExpr{
		Left:  left,
		Op:    operatorToken,
		Right: right,
		Position: lexer.Position{
			StartPos: left.Pos().StartPos,
			EndPos:   right.Pos().EndPos,
		},
	}
}

func parseGroupingExpr(p *parser) ast.Expr {
	p.advance()
	expr := parseExpr(p, defalt_bp)

	if p.currentTokenKind() == lexer.SEMI_COLON {
		p.advance()
	}

	p.expect(lexer.CLOSE_PAREN)
	return expr
}

func parseCallExpr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	startPos := p.advance().Position.StartPos
	arguments := make([]ast.Expr, 0)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_PAREN {
		arguments = append(arguments, parseExpr(p, assignment))

		if !p.currentToken().IsOneOfMany(lexer.EOF, lexer.CLOSE_PAREN) {
			if p.currentTokenKind() == lexer.COMMA {
				p.advance()
			} else {
				p.expect(lexer.SEMI_COLON)
			}
		}
	}

	if p.currentTokenKind() == lexer.COMMA {
		p.advance()
	} else if p.currentTokenKind() == lexer.SEMI_COLON {
		p.advance()
	}

	return ast.CallExpr{
		Caller: left,
		Args:   arguments,
		Position: lexer.Position{
			StartPos: startPos,
			EndPos:   p.expect(lexer.CLOSE_PAREN).Position.EndPos,
		},
	}
}
