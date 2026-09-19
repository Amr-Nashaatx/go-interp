package scanner

import (
	"strings"
	"testing"

	"github.com/Amr-Nashaatx/go-interp/token"
)

// scanAll drains the scanner and returns every token including the final EOF.
//
// The iteration cap matters: without it, a scanner that fails to advance hangs
// `go test` until the 10-minute timeout and tells you nothing. With it, the
// same bug fails in milliseconds and names the input.
func scanAll(t *testing.T, src string) []token.Token {
	t.Helper()

	s := New(src)
	var got []token.Token

	const maxTokens = 1000
	for range maxTokens {
		tok := s.NextToken()
		if tok == nil {
			t.Fatalf("NextToken returned a nil token after %d tokens", len(got))
		}

		got = append(got, *tok)
		if tok.Type == token.EOF {
			return got
		}
	}

	t.Fatalf("scanner produced %d tokens without reaching EOF — not advancing?", maxTokens)
	return nil
}

// assertTokens compares a token stream against what was expected, reporting
// every mismatch rather than stopping at the first.
func assertTokens(t *testing.T, src string, want []token.Token) {
	t.Helper()

	got := scanAll(t, src)

	if len(got) != len(want) {
		t.Errorf("got %d tokens, want %d", len(got), len(want))
		t.Logf("  got:  %s", formatTokens(got))
		t.Logf("  want: %s", formatTokens(want))
		return
	}

	for i := range want {
		if got[i].Type != want[i].Type || got[i].Lexeme != want[i].Lexeme {
			t.Errorf("token %d: got {%s %q}, want {%s %q}",
				i, got[i].Type, got[i].Lexeme, want[i].Type, want[i].Lexeme)
		}
	}
}

func formatTokens(toks []token.Token) string {
	var b strings.Builder
	for i, tk := range toks {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(string(tk.Type))
	}
	return b.String()
}

// eof is the token every scan must end with.
var eof = token.Token{Type: token.EOF, Lexeme: ""}

// --- trivia only: these must all produce nothing but EOF ------------------

func TestScansNothingButTrivia(t *testing.T) {
	cases := map[string]string{
		"empty":                   "",
		"spaces":                  "   ",
		"tabs and newlines":       "\t\n\n\t ",
		"carriage returns":        "\r\n\r\n",
		"comment only":            "// just a comment\n",
		"comment without newline": "// no trailing newline",
		"comment then blank line": "// one\n\n",
		"trivia then comment":     "  \n  // trailing",
	}

	for name, src := range cases {
		t.Run(name, func(t *testing.T) {
			assertTokens(t, src, []token.Token{eof})
		})
	}
}

// --- trivia must not swallow the tokens around it -------------------------

func TestTriviaDoesNotConsumeTokens(t *testing.T) {
	cases := []struct {
		name string
		src  string
		want []token.Token
	}{
		{
			name: "leading whitespace",
			src:  "   let",
			want: []token.Token{{Type: token.LET, Lexeme: "let"}, eof},
		},
		{
			name: "trailing whitespace",
			src:  "let   ",
			want: []token.Token{{Type: token.LET, Lexeme: "let"}, eof},
		},
		{
			// A comment must end at its newline. If it runs to end of input
			// instead, everything after it silently disappears.
			name: "comment in the middle",
			src:  "a // hidden\nb",
			want: []token.Token{
				{Type: token.IDENT, Lexeme: "a"},
				{Type: token.IDENT, Lexeme: "b"},
				eof,
			},
		},
		{
			name: "comment between two statements",
			src:  "x;\n// note\ny;",
			want: []token.Token{
				{Type: token.IDENT, Lexeme: "x"},
				{Type: token.SEMICOLON, Lexeme: ";"},
				{Type: token.IDENT, Lexeme: "y"},
				{Type: token.SEMICOLON, Lexeme: ";"},
				eof,
			},
		},
		{
			// A lone slash is division, not the start of a comment.
			name: "single slash is division",
			src:  "a / b",
			want: []token.Token{
				{Type: token.IDENT, Lexeme: "a"},
				{Type: token.DIV, Lexeme: "/"},
				{Type: token.IDENT, Lexeme: "b"},
				eof,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assertTokens(t, c.src, c.want)
		})
	}
}

// --- maximal munch on symbols --------------------------------------------

func TestSymbolsTakeTheLongestMatch(t *testing.T) {
	cases := []struct {
		name string
		src  string
		want []token.Token
	}{
		{
			// Two operators in a row must NOT merge into one lexeme.
			name: "adjacent operators stay separate",
			src:  "a=-b",
			want: []token.Token{
				{Type: token.IDENT, Lexeme: "a"},
				{Type: token.ASSIGN, Lexeme: "="},
				{Type: token.MINUS, Lexeme: "-"},
				{Type: token.IDENT, Lexeme: "b"},
				eof,
			},
		},
		{
			name: "two-character operators win over one",
			src:  "a == b != c >= d <= e",
			want: []token.Token{
				{Type: token.IDENT, Lexeme: "a"},
				{Type: token.EQUAL, Lexeme: "=="},
				{Type: token.IDENT, Lexeme: "b"},
				{Type: token.NOT_EQ, Lexeme: "!="},
				{Type: token.IDENT, Lexeme: "c"},
				{Type: token.GTE, Lexeme: ">="},
				{Type: token.IDENT, Lexeme: "d"},
				{Type: token.LTE, Lexeme: "<="},
				{Type: token.IDENT, Lexeme: "e"},
				eof,
			},
		},
		{
			name: "one-character operators when no second char matches",
			src:  "= ! < > + - * /",
			want: []token.Token{
				{Type: token.ASSIGN, Lexeme: "="},
				{Type: token.NEGATE, Lexeme: "!"},
				{Type: token.LESS, Lexeme: "<"},
				{Type: token.GREATER, Lexeme: ">"},
				{Type: token.PLUS, Lexeme: "+"},
				{Type: token.MINUS, Lexeme: "-"},
				{Type: token.MULTIP, Lexeme: "*"},
				{Type: token.DIV, Lexeme: "/"},
				eof,
			},
		},
		{
			// At end of input peekNext() returns the sentinel, so the
			// two-character candidate cannot match and the one-character
			// lookup must catch it.
			name: "one-character operator at end of input",
			src:  "a=",
			want: []token.Token{
				{Type: token.IDENT, Lexeme: "a"},
				{Type: token.ASSIGN, Lexeme: "="},
				eof,
			},
		},
		{
			name: "two-character operator at end of input",
			src:  "a>=",
			want: []token.Token{
				{Type: token.IDENT, Lexeme: "a"},
				{Type: token.GTE, Lexeme: ">="},
				eof,
			},
		},
		{
			name: "punctuation with no whitespace",
			src:  "({,;})",
			want: []token.Token{
				{Type: token.LEFT_PAREN, Lexeme: "("},
				{Type: token.L_CURL_BRACE, Lexeme: "{"},
				{Type: token.COMMA, Lexeme: ","},
				{Type: token.SEMICOLON, Lexeme: ";"},
				{Type: token.R_CURL_BRACE, Lexeme: "}"},
				{Type: token.RIGHT_PAREN, Lexeme: ")"},
				eof,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assertTokens(t, c.src, c.want)
		})
	}
}

// --- maximal munch on words: keywords are classified AFTER the munch ------

func TestKeywordsAreClassifiedAfterMunching(t *testing.T) {
	cases := []struct {
		name string
		src  string
		want []token.Token
	}{
		{
			name: "bare keywords",
			src:  "let fn true false if else return",
			want: []token.Token{
				{Type: token.LET, Lexeme: "let"},
				{Type: token.FUNC, Lexeme: "fn"},
				{Type: token.TRUE, Lexeme: "true"},
				{Type: token.FALSE, Lexeme: "false"},
				{Type: token.IFCOND, Lexeme: "if"},
				{Type: token.ELSECOND, Lexeme: "else"},
				{Type: token.RETURN, Lexeme: "return"},
				eof,
			},
		},
		{
			// Every one of these begins with a keyword. If the scanner ever
			// matches keywords during the munch rather than after it, these
			// split into a keyword plus a fragment.
			name: "identifiers that start with keywords",
			src:  "lettuce fnord truest ifer elsewhere returned",
			want: []token.Token{
				{Type: token.IDENT, Lexeme: "lettuce"},
				{Type: token.IDENT, Lexeme: "fnord"},
				{Type: token.IDENT, Lexeme: "truest"},
				{Type: token.IDENT, Lexeme: "ifer"},
				{Type: token.IDENT, Lexeme: "elsewhere"},
				{Type: token.IDENT, Lexeme: "returned"},
				eof,
			},
		},
		{
			name: "keyword immediately followed by punctuation",
			src:  "if(x)",
			want: []token.Token{
				{Type: token.IFCOND, Lexeme: "if"},
				{Type: token.LEFT_PAREN, Lexeme: "("},
				{Type: token.IDENT, Lexeme: "x"},
				{Type: token.RIGHT_PAREN, Lexeme: ")"},
				eof,
			},
		},
		{
			// Settled language rule: an identifier may START with a letter or
			// an underscore, and may CONTINUE with letters, digits or
			// underscores. Digits are the asymmetry — they can continue an
			// identifier but must not begin one, or "123abc" would be
			// ambiguous with a number.
			name: "identifiers with digits and underscores",
			src:  "a1 _x x_1 snake_case _ __private",
			want: []token.Token{
				{Type: token.IDENT, Lexeme: "a1"},
				{Type: token.IDENT, Lexeme: "_x"},
				{Type: token.IDENT, Lexeme: "x_1"},
				{Type: token.IDENT, Lexeme: "snake_case"},
				{Type: token.IDENT, Lexeme: "_"},
				{Type: token.IDENT, Lexeme: "__private"},
				eof,
			},
		},
		{
			// The other half of that rule: a digit must not start an
			// identifier. "1abc" is a NUMBER followed by an IDENT, never one
			// token.
			name: "digits cannot start an identifier",
			src:  "1abc",
			want: []token.Token{
				{Type: token.NUMBER, Lexeme: "1"},
				{Type: token.IDENT, Lexeme: "abc"},
				eof,
			},
		},
		{
			name: "single-character identifier at end of input",
			src:  "x",
			want: []token.Token{{Type: token.IDENT, Lexeme: "x"}, eof},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assertTokens(t, c.src, c.want)
		})
	}
}

// --- numbers: the only place needing two characters of lookahead ---------

func TestNumbers(t *testing.T) {
	cases := []struct {
		name string
		src  string
		want []token.Token
	}{
		{
			name: "integers",
			src:  "0 7 42 1000",
			want: []token.Token{
				{Type: token.NUMBER, Lexeme: "0"},
				{Type: token.NUMBER, Lexeme: "7"},
				{Type: token.NUMBER, Lexeme: "42"},
				{Type: token.NUMBER, Lexeme: "1000"},
				eof,
			},
		},
		{
			// The dot continues the number only because a digit follows it.
			name: "decimals",
			src:  "1.5 12.5 0.25",
			want: []token.Token{
				{Type: token.NUMBER, Lexeme: "1.5"},
				{Type: token.NUMBER, Lexeme: "12.5"},
				{Type: token.NUMBER, Lexeme: "0.25"},
				eof,
			},
		},
		{
			name: "number at end of input",
			src:  "42",
			want: []token.Token{{Type: token.NUMBER, Lexeme: "42"}, eof},
		},
		{
			name: "decimal at end of input",
			src:  "1.5",
			want: []token.Token{{Type: token.NUMBER, Lexeme: "1.5"}, eof},
		},
		{
			name: "number followed by punctuation",
			src:  "42;",
			want: []token.Token{
				{Type: token.NUMBER, Lexeme: "42"},
				{Type: token.SEMICOLON, Lexeme: ";"},
				eof,
			},
		},
		{
			name: "arithmetic without whitespace",
			src:  "1+2*3",
			want: []token.Token{
				{Type: token.NUMBER, Lexeme: "1"},
				{Type: token.PLUS, Lexeme: "+"},
				{Type: token.NUMBER, Lexeme: "2"},
				{Type: token.MULTIP, Lexeme: "*"},
				{Type: token.NUMBER, Lexeme: "3"},
				eof,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assertTokens(t, c.src, c.want)
		})
	}
}

// --- EOF is sticky --------------------------------------------------------

// Calling NextToken past the end must keep returning EOF rather than panicking
// or producing junk. A parser will read one token beyond what it needs.
func TestEOFIsRepeatable(t *testing.T) {
	s := New("x")

	for i := range 5 {
		tok := s.NextToken()
		if tok == nil {
			t.Fatalf("call %d: nil token", i)
		}
		if i == 0 {
			if tok.Type != token.IDENT {
				t.Fatalf("call 0: got %s, want IDENT", tok.Type)
			}
			continue
		}
		if tok.Type != token.EOF {
			t.Errorf("call %d: got %s, want EOF", i, tok.Type)
		}
	}
}

// --- unknown characters ---------------------------------------------------

// An unrecognised character comes back as an ILLEGAL token carrying the
// character, and scanning continues. Two bad characters therefore produce two
// tokens rather than stopping at the first.
func TestUnknownCharacterIsReported(t *testing.T) {
	cases := []struct {
		src  string
		want []token.Token
	}{
		{"@", []token.Token{
			{Type: token.ILLEGAL, Lexeme: "@"},
			*token.EOFToken,
		}},
		{"a @ b", []token.Token{
			{Type: token.IDENT, Lexeme: "a"},
			{Type: token.ILLEGAL, Lexeme: "@"},
			{Type: token.IDENT, Lexeme: "b"},
			*token.EOFToken,
		}},
		{"@#", []token.Token{
			{Type: token.ILLEGAL, Lexeme: "@"},
			{Type: token.ILLEGAL, Lexeme: "#"},
			*token.EOFToken,
		}},
	}

	for _, c := range cases {
		t.Run(c.src, func(t *testing.T) {
			assertTokens(t, c.src, c.want)
		})
	}
}
