package parser

import (
	"finescript/src/ast"
	"finescript/src/lexer"
	"fmt"
)

func parseStmt(p *parser) ast.Stmt {
	if handler, exists := stmtLU[p.currentTokenKind()]; exists {
		return handler(p)
	}

	return parseExprStmt(p)
}

func parseExprStmt(p *parser) ast.ExprStmt {
	expr := parseExpr(p, defaultBP)

	return ast.ExprStmt{
		Expr:     expr,
		Position: expr.Pos(),
	}
}

func parseBlockStmt(p *parser) ast.Stmt {
	startPos := p.advance().Position.StartPos
	body := []ast.Stmt{}

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		body = append(body, parseStmt(p))
	}

	expected := p.expect(lexer.CLOSE_CURLY)
	if expected.Kind == lexer.ERROR {
		return ast.Error{
			Position: &expected.Position,
		}
	}
	return ast.BlockStmt{
		Body: body,
		Position: lexer.Position{
			StartPos: startPos,
			EndPos:   expected.Position.EndPos,
		},
	}
}

func parseVarDecl(p *parser) ast.Stmt {
	startToken := p.advance()
	isConstant := startToken.Kind == lexer.CONST
	identName := p.expectError(lexer.IDENTIFIER,
		fmt.Sprintf("Following %s expected variable name however instead recieved %s instead\n",
			lexer.TokenKindString(startToken.Kind), lexer.TokenKindString(p.currentTokenKind())))
	if identName.Kind == lexer.ERROR {
		return ast.Error{
			Position: &identName.Position,
		}
	}

	var assignmentValue ast.Expr = nil

	if p.currentTokenKind() == lexer.ASSIGNMENT {
		p.advance()
		assignmentValue = parseExpr(p, assignment)
	}

	var endPos int
	if p.currentTokenKind() == lexer.SEMI_COLON {
		endPos = p.advance().Position.EndPos
	} else {
		endPos = identName.Position.EndPos
		if assignmentValue != nil {
			endPos = assignmentValue.Pos().EndPos
		}
	}

	if isConstant && assignmentValue == nil {
		p.errors = append(p.errors, "Cannot define constant variable without providing default value.")
	}

	if assignmentValue != nil {
		assignmentValue = ast.UndefinedLiteral{}
	}

	return ast.VarDeclStmt{
		IsConstant: isConstant,
		Name:       identName.Value,
		Value:      assignmentValue,
		Position: lexer.Position{
			StartPos: startToken.Position.StartPos,
			EndPos:   endPos,
		},
	}
}

func parseParams(p *parser) ([]ast.Param, ast.Error) {
	params := make([]ast.Param, 0)
	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_PAREN {
		expectedName := p.expect(lexer.IDENTIFIER)
		name := expectedName.Value

		if expectedName.Kind == lexer.ERROR {
			return params, ast.Error{
				Position: &expectedName.Position,
			}
		}
		expected := p.expect(lexer.COLON)
		if expected.Kind == lexer.ERROR {
			return params, ast.Error{
				Position: &expected.Position,
			}
		}

		params = append(params, ast.Param{
			Name: name,
			Type: parseType(p, defaultBP),
		})

		if p.currentTokenKind() != lexer.CLOSE_PAREN {
			expected := p.expectError(lexer.COMMA, fmt.Sprintf("Expected ',' between parameters in function declaration at %s", p.currentToken().Position.String()))
			if expected.Kind == lexer.ERROR {
				return params, ast.Error{
					Position: &expected.Position,
				}
			}
		}
	}

	return params, ast.Error{}
}

func parseFunDecl(p *parser) ast.Stmt {
	startPos := p.advance().Position.StartPos
	expectedName := p.expect(lexer.IDENTIFIER)
	name := expectedName.Value
	if expectedName.Kind == lexer.ERROR {
		return ast.Error{
			Position: &expectedName.Position,
		}
	}

	expectedOpenParen := p.expect(lexer.OPEN_PAREN)
	if expectedOpenParen.Kind == lexer.ERROR {
		return ast.Error{
			Position: &expectedOpenParen.Position,
		}
	}
	params, err := parseParams(p)
	if err.Position != nil {
		return ast.Error{
			Position: err.Position,
		}
	}
	expectedCloseParen := p.expect(lexer.CLOSE_PAREN)
	if expectedCloseParen.Kind == lexer.ERROR {
		return ast.Error{
			Position: &expectedCloseParen.Position,
		}
	}

	var returnType ast.Type = ast.VoidKeyword{}
	if p.currentTokenKind() == lexer.COLON {
		p.advance()
		returnType = parseType(p, defaultBP)
	}

	var body []ast.Stmt

	var endPos int
	if p.currentTokenKind() == lexer.OPEN_CURLY {
		blockStmt := parseBlockStmt(p).(ast.BlockStmt)
		endPos = blockStmt.Pos().EndPos
		body = blockStmt.Body
	} else {
		stmt := parseStmt(p)
		endPos = stmt.Pos().EndPos
		body = []ast.Stmt{stmt}
	}

	return ast.FunDeclStmt{
		Name:       name,
		Params:     params,
		Body:       body,
		ReturnType: returnType,
		Position: lexer.Position{
			StartPos: startPos,
			EndPos:   endPos,
		},
	}
}

func parseIfStmt(p *parser) ast.Stmt {
	startPos := p.advance().Position.StartPos
	condition := parseExpr(p, assignment)

	if p.currentTokenKind() == lexer.SEMI_COLON {
		p.advance()
	}

	consequentBlockStmt := parseBlockStmt(p).(ast.BlockStmt)
	var endPos int = consequentBlockStmt.Pos().EndPos

	var alternate []ast.Stmt
	if p.currentTokenKind() == lexer.ELSE {
		p.advance()

		if p.currentTokenKind() == lexer.IF {
			alternate = []ast.Stmt{parseIfStmt(p)}
		} else {
			alternateBlockStmt := parseBlockStmt(p).(ast.BlockStmt)
			alternate = alternateBlockStmt.Body
			endPos = alternateBlockStmt.Pos().EndPos
		}
	}

	return ast.IfStmt{
		Condition:  condition,
		Consequent: consequentBlockStmt.Body,
		Alternate:  alternate,
		Position: lexer.Position{
			StartPos: startPos,
			EndPos:   endPos,
		},
	}
}

func parseTypeDecl(p *parser) ast.Stmt {
	startPos := p.advance().Position.StartPos
	aliasExpected := p.expect(lexer.IDENTIFIER)
	alias := aliasExpected.Value
	if aliasExpected.Kind == lexer.ERROR {
		return ast.Error{
			Position: &aliasExpected.Position,
		}
	}

	expected := p.expect(lexer.ASSIGNMENT)
	if expected.Kind == lexer.ERROR {
		return ast.Error{
			Position: &expected.Position,
		}
	}

	aliasType := parseType(p, defaultBP)

	return ast.TypeAliasDecl{
		Name: alias,
		Type: aliasType,
		Position: lexer.Position{
			StartPos: startPos,
			EndPos:   aliasType.Pos().EndPos,
		},
	}
}
