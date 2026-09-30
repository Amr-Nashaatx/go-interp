package parser

import (
	"fmt"
	"strconv"

	"github.com/Amr-Nashaatx/go-interp/ast"
	"github.com/Amr-Nashaatx/go-interp/token"
)

type descent struct{}

func (d *descent) primary(pars *Parser) ast.Expression {
	if pars.check(token.NUMBER) {
		t := pars.consume()
		n, _ := strconv.ParseFloat(t.Lexeme, 32)
		return &ast.NumberLiteral{Token: t, Value: float32(n)}
	} else if pars.check(token.IDENT) {
		t := pars.consume()
		return &ast.Identifier{Token: t, Name: t.Lexeme}
	} else if pars.check(token.LEFT_PAREN) {
		pars.consume()
		exp := d.expression(pars)
		if !pars.check(token.RIGHT_PAREN) {
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
	if pars.check(token.NEGATE) || pars.check(token.MINUS) {
		op := pars.consume()
		right = &ast.PrefixExpression{Token: op, Operator: op.Lexeme, Right: d.unary(pars)}
		return right
	}
	right = d.primary(pars)
	return right
}

func (d *descent) factor(pars *Parser) ast.Expression {
	left := d.unary(pars)

	for pars.check(token.MULTIP) || pars.check(token.DIV) {
		op := pars.consume()
		right := d.unary(pars)
		left = &ast.InfixExpression{Left: left, Right: right, Operator: op.Lexeme, Token: op}
	}
	return left
}

func (d *descent) term(pars *Parser) ast.Expression {
	left := d.factor(pars)

	for pars.check(token.PLUS) || pars.check(token.MINUS) {
		op := pars.consume()
		right := d.factor(pars)
		left = &ast.InfixExpression{Left: left, Right: right, Token: op, Operator: op.Lexeme}
	}

	return left
}

func (d *descent) comparison(pars *Parser) ast.Expression {
	left := d.term(pars)

	for pars.check(token.LESS) || pars.check(token.GREATER) || pars.check(token.GTE) || pars.check(token.LTE) {
		op := pars.consume()
		right := d.term(pars)
		left = &ast.InfixExpression{Token: op, Right: right, Left: left, Operator: op.Lexeme}
	}

	return left
}
func (d *descent) equality(pars *Parser) ast.Expression {
	left := d.comparison(pars)

	for pars.check(token.EQUAL) || pars.check(token.NOT_EQ) {
		op := pars.consume()
		right := d.comparison(pars)
		left = &ast.InfixExpression{Token: op, Left: left, Right: right, Operator: op.Lexeme}
	}

	return left
}

func (d *descent) expression(pars *Parser) ast.Expression {
	return d.equality(pars)
}
