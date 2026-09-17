package parser

import (
	"github.com/Amr-Nashaatx/go-interp/ast"
	"github.com/Amr-Nashaatx/go-interp/token"
)

type Strategy interface {
	expression(p *Parser) ast.Expression
}

type TokenSource interface {
	NextToken() *token.Token
}

type Parser struct {
	src      TokenSource
	next     *token.Token
	strategy Strategy
}

func New(src TokenSource, s Strategy) *Parser {
	return &Parser{src: src, strategy: s}
}
