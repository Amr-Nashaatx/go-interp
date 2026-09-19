package parser

import (
	"fmt"
	"strconv"

	"github.com/Amr-Nashaatx/go-interp/ast"
	"github.com/Amr-Nashaatx/go-interp/token"
)

type descent struct{}

func (d *descent) primary(pars *Parser) ast.Expression {
	if pars.expect(token.NUMBER) {
		t := pars.consume()
		n, _ := strconv.ParseFloat(t.Lexeme, 32)
		return &ast.NumberLiteral{Token: t, Value: float32(n)}
	} else if pars.expect(token.IDENT) {
		t := pars.consume()
		return &ast.Identifier{Token: t, Name: t.Lexeme}
	} else if pars.expect(token.LEFT_PAREN) {
		pars.consume()
		exp := d.expression(pars)
		if !pars.expect(token.RIGHT_PAREN) {
			panic("Invalid syntax expected a ')'")
		}
		pars.consume()
		return exp
	} else {
		panic(fmt.Sprintf("can not parse %v", pars.consume().Lexeme))
	}
}

func (d *descent) unary(pars *Parser) ast.Expression {
	var right ast.Expression
	if pars.expect(token.NEGATE) || pars.expect(token.MINUS) {
		op := pars.consume()
		right = &ast.PrefixExpression{Token: op, Operator: op.Lexeme, Right: d.unary(pars)}
		return right
	}
	right = d.primary(pars)
	return right
}

func (d *descent) factor(pars *Parser) ast.Expression {
	left := d.unary(pars)

	for pars.expect(token.MULTIP) || pars.expect(token.DIV) {
		op := pars.consume()
		right := d.unary(pars)
		left = &ast.InfixExpression{Left: left, Right: right, Operator: op.Lexeme, Token: op}
	}
	return left
}

func (d *descent) term(pars *Parser) ast.Expression {
	left := d.factor(pars)

	for pars.expect(token.PLUS) || pars.expect(token.MINUS) {
		op := pars.consume()
		right := d.factor(pars)
		left = &ast.InfixExpression{Left: left, Right: right, Token: op, Operator: op.Lexeme}
	}

	return left
}

func (d *descent) comparison(pars *Parser) ast.Expression {
	left := d.term(pars)

	for pars.expect(token.LESS) || pars.expect(token.GREATER) || pars.expect(token.GTE) || pars.expect(token.LTE) {
		op := pars.consume()
		right := d.term(pars)
		left = &ast.InfixExpression{Token: op, Right: right, Left: left, Operator: op.Lexeme}
	}

	return left
}
func (d *descent) equality(pars *Parser) ast.Expression {
	left := d.comparison(pars)

	for pars.expect(token.EQUAL) || pars.expect(token.NOT_EQ) {
		op := pars.consume()
		right := d.comparison(pars)
		left = &ast.InfixExpression{Token: op, Left: left, Right: right, Operator: op.Lexeme}
	}

	return left
}

func (d *descent) expression(pars *Parser) ast.Expression {
	return d.equality(pars)
}
