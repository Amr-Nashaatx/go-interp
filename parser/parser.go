package parser

import (
	"fmt"

	"github.com/Amr-Nashaatx/go-interp/ast"
	"github.com/Amr-Nashaatx/go-interp/token"
)

type Strategy interface {
	expression(p *Parser) ast.Expression
}

var (
	Descent Strategy = &descent{}
	Pratt   Strategy = &pratt{}
)

type TokenSource interface {
	NextToken() *token.Token
}

type Parser struct {
	src      TokenSource
	next     *token.Token
	strategy Strategy
}

func New(src TokenSource, s Strategy) *Parser {
	return &Parser{src: src, strategy: s, next: src.NextToken()}
}

func (p *Parser) peek() *token.Token {
	return p.next
}

func (p *Parser) consume() *token.Token {
	prev := p.next
	p.next = p.src.NextToken()
	return prev
}

func (p *Parser) expect(t token.TokenType) bool {
	if p.peek().Type == t {
		return true
	}
	return false
}

func (p *Parser) ParseExpression() ast.Expression {
	exp := p.strategy.expression(p)
	if !p.expect(token.EOF) {
		panic(fmt.Sprintf("source did not end properly, found %v", p.peek().Lexeme))
	}
	return exp
}
