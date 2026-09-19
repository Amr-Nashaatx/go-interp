package token

// TokenType identifies what kind of thing a token is. The parser branches on
// this and nothing else.
//
// Each constant's value is its own name, not the source text it matches. The
// spelling of a token lives in exactly one place — the symbols and keywords
// tables below — so the two cannot drift apart. It also means a token prints as
// "PLUS +" rather than "+ +", which says more.
//
// These values are deliberately what `stringer` would generate for an integer
// enum. If TokenType ever becomes an int, every log line and test expectation
// keeps reading exactly the same.
type TokenType string

const (
	// Punctuation.
	LEFT_PAREN   TokenType = "LEFT_PAREN"
	RIGHT_PAREN  TokenType = "RIGHT_PAREN"
	L_CURL_BRACE TokenType = "L_CURL_BRACE"
	R_CURL_BRACE TokenType = "R_CURL_BRACE"
	COMMA        TokenType = "COMMA"
	SEMICOLON    TokenType = "SEMICOLON"

	// One-character operators.
	PLUS    TokenType = "PLUS"
	MINUS   TokenType = "MINUS"
	MULTIP  TokenType = "MULTIP"
	DIV     TokenType = "DIV"
	ASSIGN  TokenType = "ASSIGN"
	NEGATE  TokenType = "NEGATE"
	GREATER TokenType = "GREATER"
	LESS    TokenType = "LESS"

	// Two-character operators.
	EQUAL  TokenType = "EQUAL"
	NOT_EQ TokenType = "NOT_EQ"
	GTE    TokenType = "GTE"
	LTE    TokenType = "LTE"

	// Keywords.
	LET      TokenType = "LET"
	FUNC     TokenType = "FUNC"
	TRUE     TokenType = "TRUE"
	FALSE    TokenType = "FALSE"
	IFCOND   TokenType = "IFCOND"
	ELSECOND TokenType = "ELSECOND"
	RETURN   TokenType = "RETURN"

	// Literals and control.
	IDENT  TokenType = "IDENT"
	NUMBER TokenType = "NUMBER"
	EOF    TokenType = "EOF"

	// ILLEGAL carries a character the scanner could not match. It is a token
	// rather than an error so that scanning continues past bad input and the
	// consumer has one channel to read from, not two. Its lexeme is the
	// offending text.
	ILLEGAL TokenType = "ILLEGAL"
)

// symbols holds every punctuation mark and operator, of any length. This is the
// only place their spellings are written down.
var symbols = map[string]TokenType{
	"(": LEFT_PAREN,
	")": RIGHT_PAREN,
	"{": L_CURL_BRACE,
	"}": R_CURL_BRACE,
	",": COMMA,
	";": SEMICOLON,

	"+": PLUS,
	"-": MINUS,
	"*": MULTIP,
	"/": DIV,
	"=": ASSIGN,
	"!": NEGATE,
	">": GREATER,
	"<": LESS,

	"==": EQUAL,
	"!=": NOT_EQ,
	">=": GTE,
	"<=": LTE,
}

// keywords holds the reserved words. A word is scanned as an identifier first
// and looked up here afterwards.
var keywords = map[string]TokenType{
	"let":    LET,
	"fn":     FUNC,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IFCOND,
	"else":   ELSECOND,
	"return": RETURN,
}

// LookupSymbol reports the type of a punctuation mark or operator, and whether
// text is one at all.
func LookupSymbol(text string) (TokenType, bool) {
	t, ok := symbols[text]
	return t, ok
}

// LookupIdent returns the keyword type for text, or IDENT when text is an
// ordinary identifier.
func LookupIdent(text string) TokenType {
	if t, ok := keywords[text]; ok {
		return t
	}
	return IDENT
}

// Token is one lexical unit: what kind of thing it is, and the exact source
// text it came from.
type Token struct {
	Type   TokenType
	Lexeme string
	// Literal any
	// Line    int
}

// New builds a token. Lexeme is always the text the scanner actually consumed,
// so the recorded source text cannot drift from the type.
func New(t TokenType, lexeme string) *Token {
	return &Token{Type: t, Lexeme: lexeme}
}

// EOFToken marks the end of the source. It is a token rather than an error so
// that consumers can loop until they see it without a second "am I done"
// channel. Its lexeme is empty because it consumes no characters.
var EOFToken = &Token{Type: EOF, Lexeme: ""}
