package ast

import (
	"strings"

	"github.com/Amr-Nashaatx/go-interp/token"
)

// Statements satisfy both Node and Statement interfaces

type LetStatement struct {
	Token *token.Token
	Name  string
	Value Expression
}

func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) Anchor() *token.Token {
	return ls.Token
}
func (ls *LetStatement) String() string {
	var builder strings.Builder

	builder.WriteString("let ")
	builder.WriteString(ls.Name)
	builder.WriteString(" = ")
	builder.WriteString(ls.Value.String())
	builder.WriteString(";")

	return builder.String()
}

type ReturnStatement struct {
	Token *token.Token
	Value Expression
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) Anchor() *token.Token {
	return rs.Token
}
func (rs *ReturnStatement) String() string {
	var builder strings.Builder

	builder.WriteString("return ")
	builder.WriteString(rs.Value.String())
	builder.WriteString(";")
	return builder.String()
}

type ExpressionStatement struct {
	Token *token.Token
	Value Expression
}

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) Anchor() *token.Token {
	return es.Token
}
func (es *ExpressionStatement) String() string {
	var builder strings.Builder

	builder.WriteString(es.Value.String())
	builder.WriteString(";")

	return builder.String()
}

type Block struct {
	Token      *token.Token
	Statements []Statement
}

func (blk *Block) statementNode() {}

func (blk *Block) Add(stmt Statement) {
	blk.Statements = append(blk.Statements, stmt)
}

func (blk *Block) Anchor() *token.Token {
	return blk.Token
}

func (blk *Block) String() string {
	var builder strings.Builder

	builder.Write([]byte("{"))
	for _, stmt := range blk.Statements {
		builder.Write([]byte(" "))
		builder.Write([]byte(stmt.String()))
	}

	builder.Write([]byte(" }"))
	return builder.String()
}
