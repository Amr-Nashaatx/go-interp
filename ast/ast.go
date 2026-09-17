package ast

import (
	"github.com/Amr-Nashaatx/go-interp/token"
)

type Node interface {
	Anchor() *token.Token
	String() string
}

type Expression interface {
	Node
	expressionNode() // marker method, unexported
}

type Statement interface {
	Node
	statementNode()
}
