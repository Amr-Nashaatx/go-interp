package parser

import (
	"fmt"
	"strconv"

	"github.com/Amr-Nashaatx/go-interp/ast"
	"github.com/Amr-Nashaatx/go-interp/token"
)

type pratt struct{}

func isInfixOp(op string) bool {
	tok, ok := token.LookupSymbol(op)
	if !ok {
		return false
	}
	m := map[token.TokenType]bool{
		token.GREATER: true,
		token.LESS:    true,
		token.GTE:     true,
		token.LTE:     true,
		token.MINUS:   true,
		token.MULTIP:  true,
		token.PLUS:    true,
		token.DIV:     true,
		token.EQUAL:   true,
		token.NOT_EQ:  true,
	}

	if _, ok := m[tok]; ok {
		return true
	}
	return false
}

func (p *pratt) expression(pars *Parser) ast.Expression {
	return p.parseExpression(pars, 0)
}

func (p *pratt) parseExpression(pars *Parser, minBindPow int) ast.Expression {
	left := p.parsePrefix(pars)

	for {
		op := pars.peek()
		if !isInfixOp(op.Lexeme) {
			break
		}

		lbp, rbp := p.bindingPower(op.Lexeme)
		if lbp < minBindPow {
			break
		}
		pars.consume()
		right := p.parseExpression(pars, rbp)
		left = &ast.InfixExpression{Left: left, Right: right, Operator: op.Lexeme, Token: op}
	}

	return left
}

func (p *pratt) bindingPower(opType string) (int, int) {
	op, ok := token.LookupSymbol(opType)
	if !ok {
		panic(fmt.Sprintf("expected an operator, found: %v", op))
	}
	switch op {
	case token.EQUAL, token.NOT_EQ:
		return 1, 2
	case token.GREATER, token.LESS, token.LTE, token.GTE:
		return 3, 4
	case token.PLUS, token.MINUS:
		return 5, 6
	case token.MULTIP, token.DIV:
		return 7, 8
	default:
		panic(fmt.Sprintf("cannot parse operator %v", op))
	}
}
func (p *pratt) parsePrefix(pars *Parser) ast.Expression {
	if pars.expect(token.NUMBER) {
		tok := pars.consume()
		n, _ := strconv.ParseFloat(tok.Lexeme, 32)
		return &ast.NumberLiteral{Token: tok, Value: float32(n)}

	} else if pars.expect(token.IDENT) {
		tok := pars.consume()
		return &ast.Identifier{Token: tok, Name: tok.Lexeme}
	} else if pars.expect(token.LEFT_PAREN) {
		pars.consume()
		exp := p.expression(pars)
		if !pars.expect(token.RIGHT_PAREN) {
			panic("Invalid syntax expected a ')'")
		}
		pars.consume()
		return exp
	} else if pars.expect(token.NEGATE) || pars.expect(token.MINUS) {
		tok := pars.consume()
		return &ast.PrefixExpression{Token: tok, Operator: tok.Lexeme, Right: p.parseExpression(pars, 9)}
	} else {
		panic(fmt.Sprintf("can not parse %v", pars.consume().Lexeme))
	}

}
