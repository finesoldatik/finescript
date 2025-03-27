package parser

import (
	"finescript/src/ast"
	"finescript/src/lexer"
	"fmt"
)

func parseStmt(p *parser) ast.Stmt {
	stmt_fn, exists := stmt_lu[p.currentTokenKind()]

	if exists {
		return stmt_fn(p)
	}

	return parseExprStmt(p)
}

func parseExprStmt(p *parser) ast.ExprStmt {
	expr := parseExpr(p, defalt_bp)

	if p.currentTokenKind() == lexer.SEMI_COLON {
		p.advance()
	}

	return ast.ExprStmt{
		Expr: expr,
	}
}

func parseBlockStmt(p *parser) ast.Stmt {
	startPos := p.currentToken().Pos
	p.advance()
	body := []ast.Stmt{}

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		body = append(body, parseStmt(p))
	}

	endPos := &lexer.Position{
		StartLine:   startPos.StartLine,
		StartColumn: startPos.StartColumn,
		EndLine:     p.currentToken().Pos.EndLine,
		EndColumn:   p.currentToken().Pos.EndColumn,
		Index:       startPos.Index,
	}
	p.expect(lexer.CLOSE_CURLY, endPos)
	return ast.BlockStmt{
		Body: body,
		Pos:  endPos,
	}
}

func parseVarDeclStmt(p *parser) ast.Stmt {
	startToken := p.advance()
	isConstant := startToken.Kind == lexer.CONST
	identName := p.expectError(lexer.IDENTIFIER, &lexer.Position{
		StartLine:   startToken.Pos.StartLine,
		StartColumn: startToken.Pos.StartColumn,
		EndLine:     p.currentToken().Pos.EndLine,
		EndColumn:   p.currentToken().Pos.EndColumn,
		Index:       startToken.Pos.Index,
	},
		fmt.Sprintf("Following %s expected variable name however instead recieved %s instead\n",
			lexer.TokenKindString(startToken.Kind), lexer.TokenKindString(p.currentTokenKind())))

	var assignmentValue ast.Expr

	if p.currentTokenKind() == lexer.ASSIGNMENT {
		p.advance()
		assignmentValue = parseExpr(p, assignment)
	}

	var end *lexer.Position
	if p.currentTokenKind() == lexer.SEMI_COLON {
		end = p.advance().Pos
	} else {
		end = identName.Pos
		if assignmentValue != nil {
			end = assignmentValue.Position()
		}
	}

	if isConstant && assignmentValue == nil {
		panic("Cannot define constant variable without providing default value.")
	}

	return ast.VarDeclStmt{
		IsConstant: isConstant,
		Name:       identName.Value,
		Value:      assignmentValue,
		Pos: &lexer.Position{
			StartLine:   startToken.Pos.StartLine,
			StartColumn: startToken.Pos.StartColumn,
			EndLine:     end.EndLine,
			EndColumn:   end.EndColumn,
			Index:       startToken.Pos.Index,
		},
	}
}

func parseFunDeclaration(p *parser) ast.Stmt {
	startPos := p.currentToken().Pos
	p.advance()
	functionName := p.expect(lexer.IDENTIFIER, &lexer.Position{
		StartLine:   startPos.StartLine,
		StartColumn: startPos.StartColumn,
		EndLine:     p.currentToken().Pos.EndLine,
		EndColumn:   p.currentToken().Pos.EndColumn,
		Index:       startPos.Index,
	}).Value
	functionParams := make([]string, 0)

	p.expect(lexer.OPEN_PAREN, &lexer.Position{
		StartLine:   startPos.StartLine,
		StartColumn: startPos.StartColumn,
		EndLine:     p.currentToken().Pos.EndLine,
		EndColumn:   p.currentToken().Pos.EndColumn,
		Index:       startPos.Index,
	})
	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_PAREN {
		paramName := p.expect(lexer.IDENTIFIER, &lexer.Position{
			StartLine:   startPos.StartLine,
			StartColumn: startPos.StartColumn,
			EndLine:     p.currentToken().Pos.EndLine,
			EndColumn:   p.currentToken().Pos.EndColumn,
			Index:       startPos.Index,
		}).Value

		functionParams = append(functionParams, paramName)

		if !p.currentToken().IsOneOfMany(lexer.CLOSE_PAREN, lexer.EOF) {
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

	p.expect(lexer.CLOSE_PAREN, &lexer.Position{
		StartLine:   startPos.StartLine,
		StartColumn: startPos.StartColumn,
		EndLine:     p.currentToken().Pos.EndLine,
		EndColumn:   p.currentToken().Pos.EndColumn,
		Index:       startPos.Index,
	})
	var functionBody []ast.Stmt

	var end *lexer.Position
	if p.currentTokenKind() == lexer.OPEN_CURLY {
		blockStmt := ast.ExpectStmt[ast.BlockStmt](parseBlockStmt(p))
		end = blockStmt.Pos
		functionBody = blockStmt.Body
	} else {
		stmt := parseStmt(p)
		end = stmt.Position()
		functionBody = []ast.Stmt{stmt}
	}

	return ast.FunctionDeclStmt{
		Params: functionParams,
		Body:   functionBody,
		Name:   functionName,
		Pos: &lexer.Position{
			StartLine:   startPos.StartLine,
			StartColumn: startPos.StartColumn,
			EndLine:     end.EndLine,
			EndColumn:   end.EndColumn,
			Index:       startPos.Index,
		},
	}
}

func parseIfStmt(p *parser) ast.Stmt {
	startPos := p.currentToken().Pos
	p.advance()
	condition := parseExpr(p, assignment)

	if p.currentTokenKind() == lexer.SEMI_COLON {
		p.advance()
	}

	consequentBlockStmt := ast.ExpectStmt[ast.BlockStmt](parseBlockStmt(p))
	var end = consequentBlockStmt.Pos

	var alternate []ast.Stmt
	if p.currentTokenKind() == lexer.ELSE {
		p.advance()

		if p.currentTokenKind() == lexer.IF {
			alternate = []ast.Stmt{parseIfStmt(p)}
		} else {
			alternateBlockStmt := ast.ExpectStmt[ast.BlockStmt](parseBlockStmt(p))
			alternate = alternateBlockStmt.Body
			end = alternateBlockStmt.Pos
		}
	}

	return ast.IfStmt{
		Condition:  condition,
		Consequent: consequentBlockStmt.Body,
		Alternate:  alternate,
		Pos: &lexer.Position{
			StartLine:   startPos.StartLine,
			StartColumn: startPos.StartColumn,
			EndLine:     end.EndLine,
			EndColumn:   end.EndColumn,
			Index:       startPos.Index,
		},
	}
}

func parseLoopStmt(p *parser) ast.Stmt {
	startPos := p.currentToken().Pos
	p.advance()
	blockStmt := ast.ExpectStmt[ast.BlockStmt](parseBlockStmt(p))
	return ast.LoopStmt{
		Body: blockStmt.Body,
		Pos: &lexer.Position{
			StartLine:   startPos.StartLine,
			StartColumn: startPos.StartColumn,
			EndLine:     blockStmt.Pos.EndLine,
			EndColumn:   blockStmt.Pos.EndColumn,
			Index:       startPos.Index,
		},
	}
}

func parseLoopControl(p *parser) ast.Stmt {
	pos := p.currentToken().Pos
	switch p.currentTokenKind() {
	case lexer.BREAK:
		p.advance()
		return ast.BreakStmt{
			Pos: pos,
		}
	case lexer.CONTINUE:
		p.advance()
		return ast.ContinueStmt{
			Pos: pos,
		}
	default:
		panic("Unknown Loop Control Operator")
	}
}
