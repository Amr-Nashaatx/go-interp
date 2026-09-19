package scanner

import (
	"testing"

	"github.com/Amr-Nashaatx/go-interp/token"
)

func TestScanner(t *testing.T) {
	input := `let lettuce = 12.5;
	if (lettuce >= 5) {
		// check it
		return lettuce != 0;
	}
	`
	s := New(input)

	want := []token.Token{
		{Type: token.LET, Lexeme: "let"},
		{Type: token.IDENT, Lexeme: "lettuce"},
		{Type: token.ASSIGN, Lexeme: "="},
		{Type: token.NUMBER, Lexeme: "12.5"},
		{Type: token.SEMICOLON, Lexeme: ";"},
		{Type: token.IFCOND, Lexeme: "if"},
		{Type: token.LEFT_PAREN, Lexeme: "("},
		{Type: token.IDENT, Lexeme: "lettuce"},
		{Type: token.GTE, Lexeme: ">="},
		{Type: token.NUMBER, Lexeme: "5"},
		{Type: token.RIGHT_PAREN, Lexeme: ")"},
		{Type: token.L_CURL_BRACE, Lexeme: "{"},
		{Type: token.RETURN, Lexeme: "return"},
		{Type: token.IDENT, Lexeme: "lettuce"},
		{Type: token.NOT_EQ, Lexeme: "!="},
		{Type: token.NUMBER, Lexeme: "0"},
		{Type: token.SEMICOLON, Lexeme: ";"},
		{Type: token.R_CURL_BRACE, Lexeme: "}"},
		*token.EOFToken,
	}

	got := make([]token.Token, 20)
	for range len(want) {
		tok := s.NextToken()
		got = append(got, *token.New(tok.Type, tok.Lexeme))
	}

	t.Run("Scanner", func(t *testing.T) {
		assertTokens(t, input, want)
	})
}
