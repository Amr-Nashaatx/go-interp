package ast

import (
	"testing"

	"github.com/Amr-Nashaatx/go-interp/token"
)

// Each node must satisfy its category. Catches a missing method at build time;
// TestAnchor catches one that exists but does not work.
var (
	_ Statement = (*LetStatement)(nil)
	_ Statement = (*ReturnStatement)(nil)
	_ Statement = (*ExpressionStatement)(nil)

	_ Expression = (*Identifier)(nil)
	_ Expression = (*NumberLiteral)(nil)
	_ Expression = (*PrefixExpression)(nil)
	_ Expression = (*InfixExpression)(nil)

	_ Node = (*Program)(nil)
)

func num(lexeme string, v float32) *NumberLiteral {
	return &NumberLiteral{Token: token.New(token.NUMBER, lexeme), Value: v}
}

func infix(t token.TokenType, op string, left, right Expression) *InfixExpression {
	return &InfixExpression{Token: token.New(t, op), Left: left, Operator: op, Right: right}
}

// The same tokens shaped two ways. String() has to tell them apart.
func TestTreeRendering(t *testing.T) {
	one, two, three := num("1", 1), num("2", 2), num("3", 3)

	tests := []struct {
		name  string
		value Expression
		want  string
	}{
		{
			name:  "1 + 2 * 3",
			value: infix(token.PLUS, "+", one, infix(token.MULTIP, "*", two, three)),
			want:  "let x = (1 + (2 * 3));",
		},
		{
			name:  "(1 + 2) * 3",
			value: infix(token.MULTIP, "*", infix(token.PLUS, "+", one, two), three),
			want:  "let x = ((1 + 2) * 3);",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := &LetStatement{
				Token: token.New(token.LET, "let"),
				Name:  "x",
				Value: tt.value,
			}
			if got := stmt.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

// Every node reports the token it starts at. A node that satisfies Node by
// embedding it compiles cleanly and then panics here.
func TestAnchor(t *testing.T) {
	tests := []struct {
		name string
		node Node
		want token.Token
	}{
		{
			name: "statement anchors at its keyword",
			node: &LetStatement{Token: token.New(token.LET, "let")},
			want: token.Token{Type: token.LET, Lexeme: "let"},
		},
		{
			name: "infix anchors at its operator",
			node: &InfixExpression{Token: token.New(token.PLUS, "+")},
			want: token.Token{Type: token.PLUS, Lexeme: "+"},
		},
		{
			name: "leaf anchors at itself",
			node: num("42", 42),
			want: token.Token{Type: token.NUMBER, Lexeme: "42"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.node.Anchor()
			if got == nil {
				t.Fatalf("Anchor() = nil, want %+v", tt.want)
			}
			if *got != tt.want {
				t.Errorf("Anchor() = %+v, want %+v", *got, tt.want)
			}
		})
	}
}
