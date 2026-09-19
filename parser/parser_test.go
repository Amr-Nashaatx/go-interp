package parser_test

import (
	"testing"

	"github.com/Amr-Nashaatx/go-interp/parser"
	"github.com/Amr-Nashaatx/go-interp/scanner"
)

// Both strategies must build the same tree from the same source.
var strategies = map[string]parser.Strategy{
	"descent": parser.Descent,
}

func TestParseExpression(t *testing.T) {
	tests := []struct {
		src  string
		want string
	}{
		{"1 + 2 * 3", "(1 + (2 * 3))"},
		{"1 * 2 + 3", "((1 * 2) + 3)"},
		{"a - b - c", "((a - b) - c)"},
		{"(1 + 2) * 3", "((1 + 2) * 3)"},
		{"-a * b", "((-a) * b)"},
		{"a + b == c < d", "((a + b) == (c < d))"},
	}

	for name, s := range strategies {
		for _, tt := range tests {
			t.Run(name+"/"+tt.src, func(t *testing.T) {
				got := parser.New(scanner.New(tt.src), s).ParseExpression().String()
				if got != tt.want {
					t.Errorf("got %q, want %q", got, tt.want)
				}
			})
		}
	}
}

// Malformed input must fail rather than return a tree.
func TestParseExpressionRejects(t *testing.T) {
	for name, s := range strategies {
		for _, src := range []string{"1 +", "(1 + 2", "* 3", "1 2"} {
			t.Run(name+"/"+src, func(t *testing.T) {
				defer func() {
					if recover() == nil {
						t.Errorf("parsed %q without failing", src)
					}
				}()
				parser.New(scanner.New(src), s).ParseExpression()
			})
		}
	}
}
