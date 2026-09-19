package ast

import (
	"github.com/Amr-Nashaatx/go-interp/token"
	"strconv"
	"strings"
)

type Identifier struct {
	Token *token.Token
	Name  string
}

func (id *Identifier) expressionNode() {}

func (id *Identifier) Anchor() *token.Token {
	return id.Token
}
func (id *Identifier) String() string {
	return id.Name
}

// ----------------------------
type NumberLiteral struct {
	Token *token.Token
	Value float32
}

func (nl *NumberLiteral) expressionNode() {}

func (nl *NumberLiteral) Anchor() *token.Token {
	return nl.Token
}
func (nl *NumberLiteral) String() string {
	return strconv.FormatFloat(float64(nl.Value), 'f', -1, 32)
}

// ----------------------------
type PrefixExpression struct {
	Token    *token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

func (pf *PrefixExpression) Anchor() *token.Token {
	return pf.Token
}
func (pf *PrefixExpression) String() string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(pf.Operator)
	builder.WriteString(pf.Right.String())
	builder.WriteString(")")

	return builder.String()
}

// ----------------------------
type InfixExpression struct {
	Token    *token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}

func (ie *InfixExpression) Anchor() *token.Token {
	return ie.Token
}
func (ie *InfixExpression) String() string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(ie.Left.String())
	builder.WriteString(" " + ie.Operator + " ")
	builder.WriteString(ie.Right.String())
	builder.WriteString(")")

	return builder.String()
}
