package ast

import (
	"strings"

	"github.com/Amr-Nashaatx/go-interp/token"
)

type Program struct {
	Statements []Statement
}

func (p *Program) Add(s Statement) {
	p.Statements = append(p.Statements, s)
}

func (p *Program) String() string {
	var builder strings.Builder

	for _, s := range p.Statements {
		builder.WriteString(s.String() + "\n")
	}

	return builder.String()
}

func (p *Program) Anchor() *token.Token {
	if len(p.Statements) == 0 {
		return nil
	}
	return p.Statements[0].Anchor()
}
