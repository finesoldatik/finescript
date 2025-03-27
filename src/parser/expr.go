package parser

import (
	"finescript/src/ast"
	"finescript/src/lexer"
	"fmt"
	"strconv"
)

func parseExpr(p *parser, bp binding_power) ast.Expr {
	startPos := p.currentToken().Pos
	tokenKind := p.currentTokenKind()
	nud_fn, exists := nud_lu[tokenKind]

	if !exists {
		panic(lexer.FormatError(p.initialSource, &lexer.Position{
			StartLine:   startPos.StartLine,
			StartColumn: startPos.StartColumn,
			EndLine:     p.currentToken().Pos.EndLine,
			EndColumn:   p.currentToken().Pos.EndColumn,
			Index:       startPos.Index,
		}, fmt.Sprintf("NUD Handler expected for token %s at %s\n", lexer.TokenKindString(tokenKind), startPos.String())))
	}

	left := nud_fn(p)

	for bp_lu[p.currentTokenKind()] > bp {
		tokenKind = p.currentTokenKind()
		led_fn, exists := led_lu[tokenKind]

		if !exists {
			panic(lexer.FormatError(p.initialSource, &lexer.Position{
				StartLine:   startPos.StartLine,
				StartColumn: startPos.StartColumn,
				EndLine:     p.currentToken().Pos.EndLine,
				EndColumn:   p.currentToken().Pos.EndColumn,
				Index:       startPos.Index,
			}, fmt.Sprintf("LED Handler expected for token %s at %s\n", lexer.TokenKindString(tokenKind), startPos.String())))
		}

		left = led_fn(p, left, bp)
	}

	return left
}

func parsePrimaryExpr(p *parser) ast.Expr {
	pos := p.currentToken().Pos
	switch p.currentTokenKind() {
	case lexer.INT:
		number, _ := strconv.ParseInt(p.advance().Value, 0, 64)
		return ast.IntLiteral{
			Value: number,
			Pos:   pos,
		}
	case lexer.FLOAT:
		number, _ := strconv.ParseFloat(p.advance().Value, 64)
		return ast.FloatLiteral{
			Value: number,
			Pos:   pos,
		}
	case lexer.STRING:
		return ast.StringLiteral{
			Value: p.advance().Value,
			Pos:   pos,
		}
	case lexer.IDENTIFIER:
		return ast.Identifier{
			Name: p.advance().Value,
			Pos:  pos,
		}
	case lexer.NULL:
		p.advance()
		return ast.NullLiteral{
			Pos: pos,
		}
	case lexer.TRUE:
		p.advance()
		return ast.BoolLiteral{
			Value: true,
			Pos:   pos,
		}
	case lexer.FALSE:
		p.advance()
		return ast.BoolLiteral{
			Value: false,
			Pos:   pos,
		}
	default:
		panic(fmt.Sprintf("Cannot create primary_expr from %s\n", lexer.TokenKindString(p.currentTokenKind())))
	}
}

func parseUnaryExpr(p *parser) ast.Expr {
	pos := p.currentToken().Pos
	operatorToken := p.advance()
	expr := parseExpr(p, unary)

	return ast.UnaryExpr{
		Op:    operatorToken,
		Value: expr,
		Pos: &lexer.Position{
			StartLine:   pos.StartLine,
			StartColumn: pos.StartColumn,
			EndLine:     expr.Position().EndLine,
			EndColumn:   expr.Position().EndColumn,
			Index:       pos.Index,
		},
	}
}

func parseLedUnaryExpr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	operatorToken := p.advance()

	if p.currentTokenKind() == lexer.SEMI_COLON {
		p.advance()
	}

	return ast.UnaryExpr{
		Op:    operatorToken,
		Value: left,
		Pos: &lexer.Position{
			StartLine:   left.Position().StartLine,
			StartColumn: left.Position().StartColumn,
			EndLine:     operatorToken.Pos.EndLine,
			EndColumn:   operatorToken.Pos.EndColumn,
			Index:       left.Position().Index,
		},
	}
}

func parseAssignExpr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	p.advance()
	rhs := parseExpr(p, bp)

	return ast.AssignExpr{
		Left:  left,
		Right: rhs,
		Pos: &lexer.Position{
			StartLine:   left.Position().StartLine,
			StartColumn: left.Position().StartColumn,
			EndLine:     rhs.Position().EndLine,
			EndColumn:   rhs.Position().EndColumn,
			Index:       left.Position().Index,
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
		Pos: &lexer.Position{
			StartLine:   left.Position().StartLine,
			StartColumn: left.Position().StartColumn,
			EndLine:     right.Position().EndLine,
			EndColumn:   right.Position().EndColumn,
			Index:       left.Position().Index,
		},
	}
}

func parseMemberExpr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	firstToken := p.advance()
	isComputed := firstToken.Kind == lexer.OPEN_BRACKET

	if isComputed {
		rhs := parseExpr(p, bp)
		end := p.expect(lexer.CLOSE_BRACKET, &lexer.Position{
			StartLine:   firstToken.Pos.StartLine,
			StartColumn: firstToken.Pos.StartColumn,
			EndLine:     p.currentToken().Pos.EndLine,
			EndColumn:   p.currentToken().Pos.EndColumn,
			Index:       firstToken.Pos.Index,
		})
		return ast.MemberExpr{
			Object: left,
			Field:  rhs,
			Pos: &lexer.Position{
				StartLine:   left.Position().StartLine,
				StartColumn: left.Position().StartColumn,
				EndLine:     end.Pos.EndLine,
				EndColumn:   end.Pos.EndColumn,
				Index:       left.Position().Index,
			},
		}
	}

	endPos := &lexer.Position{
		StartLine:   left.Position().StartLine,
		StartColumn: left.Position().StartColumn,
		EndLine:     p.currentToken().Pos.EndLine,
		EndColumn:   p.currentToken().Pos.EndColumn,
		Index:       left.Position().Index,
	}
	ident := p.expect(lexer.IDENTIFIER, endPos)
	return ast.MemberExpr{
		Object: left,
		Field: ast.Identifier{
			Name: ident.Value,
		},
		Pos: endPos,
	}
}

func parseArrayLiteralExpr(p *parser) ast.Expr {
	startPos := p.currentToken().Pos
	p.advance()
	arrayContents := make([]ast.Expr, 0)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_BRACKET {
		elemPos := p.currentToken().Pos
		arrayContents = append(arrayContents, parseExpr(p, logical))

		if !p.currentToken().IsOneOfMany(lexer.EOF, lexer.CLOSE_BRACKET) {
			p.expect(lexer.COMMA, &lexer.Position{
				StartLine:   startPos.StartLine,
				StartColumn: startPos.StartColumn,
				EndLine:     elemPos.EndLine,
				EndColumn:   elemPos.EndColumn,
				Index:       startPos.Index,
			})
		}
	}

	if p.currentTokenKind() == lexer.COMMA {
		p.advance()
	}

	endPos := &lexer.Position{
		StartLine:   startPos.StartLine,
		StartColumn: startPos.StartColumn,
		EndLine:     p.currentToken().Pos.EndLine,
		EndColumn:   p.currentToken().Pos.EndColumn,
		Index:       startPos.Index,
	}
	p.expect(lexer.CLOSE_BRACKET, endPos)
	return ast.ArrayLiteral{
		Elements: arrayContents,
		Pos:      endPos,
	}
}

func parseGroupingExpr(p *parser) ast.Expr {
	startPos := p.advance().Pos
	expr := parseExpr(p, defalt_bp)

	if p.currentTokenKind() == lexer.SEMI_COLON {
		p.advance()
	}

	endPos := &lexer.Position{
		StartLine:   startPos.StartLine,
		StartColumn: startPos.StartColumn,
		EndLine:     p.currentToken().Pos.EndLine,
		EndColumn:   p.currentToken().Pos.EndColumn,
		Index:       startPos.Index,
	}

	p.expect(lexer.CLOSE_PAREN, endPos)
	return expr
}

func parseCallExpr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	startPos := p.currentToken().Pos
	p.advance()
	arguments := make([]ast.Expr, 0)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_PAREN {
		arguments = append(arguments, parseExpr(p, assignment))

		if !p.currentToken().IsOneOfMany(lexer.EOF, lexer.CLOSE_PAREN) {
			if p.currentTokenKind() == lexer.COMMA {
				p.advance()
			} else {
				p.expect(lexer.SEMI_COLON, &lexer.Position{
					StartLine:   startPos.StartLine,
					StartColumn: startPos.StartColumn,
					EndLine:     p.currentToken().Pos.EndLine,
					EndColumn:   p.currentToken().Pos.EndColumn,
					Index:       startPos.Index,
				})
			}
		}
	}

	if p.currentTokenKind() == lexer.COMMA {
		p.advance()
	} else if p.currentTokenKind() == lexer.SEMI_COLON {
		p.advance()
	}

	endPos := &lexer.Position{
		StartLine:   startPos.StartLine,
		StartColumn: startPos.StartColumn,
		EndLine:     p.currentToken().Pos.EndLine,
		EndColumn:   p.currentToken().Pos.EndColumn,
		Index:       startPos.Index,
	}
	p.expect(lexer.CLOSE_PAREN, endPos)
	return ast.CallExpr{
		Caller: left,
		Args:   arguments,
		Pos:    endPos,
	}
}
