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

	return ast.ExprStmt{
		Expr: expr,
	}
}

func parseBlockStmt(p *parser) ast.Stmt {
	startPos := p.advance().Position.StartPos
	body := []ast.Stmt{}

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		body = append(body, parseStmt(p))
	}

	return ast.BlockStmt{
		Body: body,
		Position: lexer.Position{
			StartPos: startPos,
			EndPos:   p.expect(lexer.CLOSE_CURLY).Position.EndPos,
		},
	}
}

func parseVarDecl(p *parser) ast.Stmt {
	startToken := p.advance()
	isConstant := startToken.Kind == lexer.CONST
	identName := p.expectError(lexer.IDENTIFIER,
		fmt.Sprintf("Following %s expected variable name however instead recieved %s instead\n",
			lexer.TokenKindString(startToken.Kind), lexer.TokenKindString(p.currentTokenKind())))

	var assignmentValue ast.Expr

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
		p.error("Cannot define constant variable without providing default value.", nil)
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

func parseFunDecl(p *parser) ast.Stmt {
	startPos := p.advance().Position.StartPos
	name := p.expect(lexer.IDENTIFIER).Value
	params := make([]ast.Param, 0)

	p.expect(lexer.OPEN_PAREN)

	parseParams(p, params)

	p.expect(lexer.CLOSE_PAREN)
	var body []ast.Stmt

	var endPos int
	if p.currentTokenKind() == lexer.OPEN_CURLY {
		blockStmt := ast.ExpectStmt[ast.BlockStmt](parseBlockStmt(p))
		endPos = blockStmt.Pos().EndPos
		body = blockStmt.Body
	} else {
		stmt := parseStmt(p)
		endPos = stmt.Pos().EndPos
		body = []ast.Stmt{stmt}
	}

	return ast.FunDeclStmt{
		Params: params,
		Body:   body,
		Name:   name,
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

	consequentBlockStmt := ast.ExpectStmt[ast.BlockStmt](parseBlockStmt(p))
	var endPos int = consequentBlockStmt.Pos().EndPos

	var alternate []ast.Stmt
	if p.currentTokenKind() == lexer.ELSE {
		p.advance()

		if p.currentTokenKind() == lexer.IF {
			alternate = []ast.Stmt{parseIfStmt(p)}
		} else {
			alternateBlockStmt := ast.ExpectStmt[ast.BlockStmt](parseBlockStmt(p))
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
