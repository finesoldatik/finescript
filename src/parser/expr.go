package parser

import (
	"finescript/src/ast"
	"finescript/src/lexer"
	"fmt"
	"strconv"
)

func parseExpr(p *parser, bp bindingPower) ast.Expr {
	token := p.currentToken()
	nudFn, exists := nudLU[token.Kind]

	if !exists {
		p.errors = append(p.errors, fmt.Sprintf("NUD Handler expected for token %s at %s\n", lexer.TokenKindString(token.Kind), token.Position.String()))
		return ast.Error{
			Position: &token.Position,
		}
	}

	left := nudFn(p)

	for bpLU[p.currentTokenKind()] > bp {
		token = p.currentToken()
		ledFn, exists := ledLU[token.Kind]

		if !exists {
			p.errors = append(p.errors, fmt.Sprintf("LED Handler expected for token %s at %s\n", lexer.TokenKindString(token.Kind), token.Position.String()))
			return ast.Error{
				Position: &token.Position,
			}
		}

		left = ledFn(p, left, bp)
	}

	if p.currentTokenKind() == lexer.SEMI_COLON {
		p.advance()
	}

	return left
}

func parsePrimaryExpr(p *parser) ast.Expr {
	token := p.currentToken()
	switch token.Kind {
	case lexer.INT:
		token := p.advance()
		number, err := strconv.ParseInt(token.Value, 0, 64)
		if err != nil {
			panic(fmt.Sprintf("invalid integer literal: %s", token.Value))
		}
		return ast.IntLiteral{
			Value:    number,
			Position: token.Position,
		}
	case lexer.FLOAT:
		token := p.advance()
		number, err := strconv.ParseFloat(token.Value, 64)
		if err != nil {
			panic(fmt.Sprintf("invalid float literal: %s", token.Value))
		}
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
		p.errors = append(p.errors, fmt.Sprintf("Cannot create primary_expr from %s at %s", lexer.TokenKindString(token.Kind), token.Position.String()))
		return ast.Error{
			Position: &token.Position,
		}
	}
}

func parseUnaryExpr(p *parser) ast.Expr {
	operatorToken := p.advance()
	expr := parseExpr(p, unary)
	if err, ok := expr.(ast.Error); ok {
		return err
	}

	return ast.UnaryExpr{
		Op:   operatorToken,
		Expr: expr,
		Position: lexer.Position{
			StartPos: operatorToken.Position.StartPos,
			EndPos:   expr.Pos().EndPos,
		},
	}
}

func parseLedUnaryExpr(p *parser, left ast.Expr, bp bindingPower) ast.Expr {
	if err, ok := left.(ast.Error); ok {
		return err
	}
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

func parseAssignExpr(p *parser, left ast.Expr, bp bindingPower) ast.Expr {
	if err, ok := left.(ast.Error); ok {
		return err
	}
	op := p.advance()
	expr := parseExpr(p, bp)
	if err, ok := expr.(ast.Error); ok {
		return err
	}

	return ast.AssignExpr{
		Assigne: left,
		Op:      op,
		Expr:    expr,
		Position: lexer.Position{
			StartPos: left.Pos().StartPos,
			EndPos:   expr.Pos().EndPos,
		},
	}
}

func parseBinaryExpr(p *parser, left ast.Expr, bp bindingPower) ast.Expr {
	if err, ok := left.(ast.Error); ok {
		return err
	}
	operatorToken := p.advance()
	right := parseExpr(p, defaultBP)
	if err, ok := right.(ast.Error); ok {
		return err
	}

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
	expr := parseExpr(p, defaultBP)
	if err, ok := expr.(ast.Error); ok {
		return err
	}

	if p.currentTokenKind() == lexer.SEMI_COLON {
		p.advance()
	}

	expected := p.expect(lexer.CLOSE_PAREN)
	if expected.Kind == lexer.ERROR {
		return ast.Error{
			Position: &expected.Position,
		}
	}
	return expr
}

func parseCallExpr(p *parser, left ast.Expr, bp bindingPower) ast.Expr {
	if err, ok := left.(ast.Error); ok {
		return err
	}
	startPos := p.advance().Position.StartPos
	arguments := make([]ast.Expr, 0)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_PAREN {
		expr := parseExpr(p, assignment)
		if err, ok := expr.(ast.Error); ok {
			return err
		}
		arguments = append(arguments, expr)

		if p.currentTokenKind() != lexer.CLOSE_PAREN {
			expected := p.expectError(lexer.COMMA, fmt.Sprintf("Expected ',' between parameters in function declaration at %s", p.currentToken().Position.String()))
			if expected.Kind == lexer.ERROR {
				return ast.Error{
					Position: &expected.Position,
				}
			}
		}
	}

	expected := p.expect(lexer.CLOSE_PAREN)
	if expected.Kind == lexer.ERROR {
		return ast.Error{
			Position: &expected.Position,
		}
	}
	return ast.CallExpr{
		Caller: left,
		Args:   arguments,
		Position: lexer.Position{
			StartPos: startPos,
			EndPos:   expected.Position.EndPos,
		},
	}
}

func parseConditionalExpr(p *parser, left ast.Expr, bp bindingPower) ast.Expr {
	if err, ok := left.(ast.Error); ok {
		return err
	}
	p.advance()

	consequent := parseExpr(p, defaultBP)
	if err, ok := consequent.(ast.Error); ok {
		return err
	}

	
	expected := p.expect(lexer.COLON)
	if expected.Kind == lexer.ERROR {
		return ast.Error{
			Position: &expected.Position,
		}
	}
	alternate := parseExpr(p, defaultBP)
	if err, ok := alternate.(ast.Error); ok {
		return err
	}

	return ast.ConditionalExpr{
		Condition:  left,
		Consequent: consequent,
		Alternate:  alternate,
		Position: lexer.Position{
			StartPos: left.Pos().StartPos,
			EndPos:   alternate.Pos().EndPos,
		},
	}
}
