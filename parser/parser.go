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

func (p *Parser) check(t token.TokenType) bool {
	if p.peek().Type == t {
		return true
	}
	return false
}

func (p *Parser) require(t token.TokenType) *token.Token {
	tok := p.peek()
	if tok.Type != t {
		panic(fmt.Sprintf("Invalid token: expected %v, found %v", t, tok))
	}
	return p.consume()
}

func (p *Parser) ParseExpression() ast.Expression {
	exp := p.strategy.expression(p)
	if !p.check(token.EOF) {
		panic(fmt.Sprintf("source did not end properly, found %v", p.peek().Lexeme))
	}
	return exp
}

func (p *Parser) ParseProgram() *ast.Program {
	prog := &ast.Program{}
	for !p.check(token.EOF) {
		prog.Statements = append(prog.Statements, p.statement())
	}
	return prog
}

func (p *Parser) statement() ast.Statement {
	tok := p.peek()
	switch tok.Type {
	case token.LET:
		return p.letStatement()
	case token.RETURN:
		return p.returnStatement()
	case token.L_CURL_BRACE:
		return p.block()
	default:
		return p.expressionStatement()
	}
}

func (p *Parser) letStatement() ast.Statement {
	tok := p.consume() // the LET keyword
	ident := p.require(token.IDENT)
	p.require(token.ASSIGN)
	value := p.strategy.expression(p)
	p.require(token.SEMICOLON)
	return &ast.LetStatement{Token: tok, Name: ident.Lexeme, Value: value}
}

func (p *Parser) returnStatement() ast.Statement {
	tok := p.consume() // RETURN keyword
	value := p.strategy.expression(p)
	p.require(token.SEMICOLON)
	return &ast.ReturnStatement{Token: tok, Value: value}
}

func (p *Parser) expressionStatement() ast.Statement {
	tok := p.peek()
	expr := p.strategy.expression(p)
	p.require(token.SEMICOLON)
	return &ast.ExpressionStatement{Token: tok, Value: expr}
}

func (p *Parser) block() ast.Statement {
	tok := p.consume() // {
	blk := ast.Block{Token: tok}

	for !p.check(token.EOF) && !p.check(token.R_CURL_BRACE) {
		blk.Add(p.statement())
	}

	p.require(token.R_CURL_BRACE)
	return &blk
}
