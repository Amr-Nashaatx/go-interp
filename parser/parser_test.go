package parser_test

import (
	"testing"

	"github.com/Amr-Nashaatx/go-interp/parser"
	"github.com/Amr-Nashaatx/go-interp/scanner"
)

// Both strategies must build the same tree from the same source.
var strategies = map[string]parser.Strategy{
	"descent": parser.Descent,
	"pratt":   parser.Pratt,
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

func TestParseProgram(t *testing.T) {
	tests := []struct {
		src  string
		want string
	}{
		{"let x = 1 + 2;", "let x = (1 + 2);\n"},
		{"return x;", "return x;\n"},
		{"x; { let y = x; { y; } }", "x;\n{ let y = x; { y; } }\n"},
		{"{ }", "{ }\n"},
		{"", ""},
	}

	for name, s := range strategies {
		for _, tt := range tests {
			t.Run(name+"/"+tt.src, func(t *testing.T) {
				got := parser.New(scanner.New(tt.src), s).ParseProgram().String()
				if got != tt.want {
					t.Errorf("got %q, want %q", got, tt.want)
				}
			})
		}
	}
}

// Each statement is anchored at the token that selected its rule.
func TestStatementAnchors(t *testing.T) {
	prog := parser.New(scanner.New("let a = 1; b + 2; { c; }"), parser.Pratt).ParseProgram()
	want := []string{"let", "b", "{"}

	for i, s := range prog.Statements {
		if got := s.Anchor().Lexeme; got != want[i] {
			t.Errorf("statement %d: anchor %q, want %q", i, got, want[i])
		}
	}
}

// Malformed programs must fail rather than return a tree.
func TestParseProgramRejects(t *testing.T) {
	for name, s := range strategies {
		for _, src := range []string{"let x 1;", "x", "{ x;", "x; }"} {
			t.Run(name+"/"+src, func(t *testing.T) {
				defer func() {
					if recover() == nil {
						t.Errorf("parsed %q without failing", src)
					}
				}()
				parser.New(scanner.New(src), s).ParseProgram()
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
