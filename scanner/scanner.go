package scanner

import (
	"fmt"

	"github.com/Amr-Nashaatx/go-interp/token"
)

type Scanner struct {
	src   string
	curr  int
	start int
	line  int
}

func (s *Scanner) peek() byte {
	if s.curr+1 >= len(s.src) {
		return 0
	}
	peeked := s.src[s.curr+1]
	return peeked
}

func (s *Scanner) peekNext() byte {
	if s.curr+2 >= len(s.src) {
		return 0
	}
	peeked := s.src[s.curr+2]
	return peeked
}

func (s *Scanner) advance() {
	s.curr += 1
}

func (s *Scanner) lexeme() string {
	t := s.src[s.start : s.curr+1]
	return t
}

func (s *Scanner) skipWsAndComments() {
	for {
		switch s.peek() {
		case ' ', '\t', '\r':
			s.advance()
		case '\n':
			s.line++
			s.advance()
		case '/':
			if s.peekNext() != '/' {
				s.start = s.curr + 1
				return
			}
			for s.peek() != '\n' && !(s.curr+1 > len(s.src)) {
				s.advance()
			}
		default:
			s.start = s.curr + 1
			return
		}
	}
}

func (s *Scanner) scanIdentifier() string {
	for {
		peeked := s.peek()
		if token.IsDigit(peeked) || token.IsLetter(peeked) || peeked == '_' {
			s.advance()
		} else {
			break
		}
	}

	return s.lexeme()
}

func (s *Scanner) scanNumber() string {
	for token.IsDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && token.IsDigit(s.peekNext()) {
		s.advance()
		for token.IsDigit(s.peek()) {
			s.advance()
		}
	}

	return s.lexeme()
}

func New(content string) *Scanner {
	scanner := Scanner{src: content, curr: -1, start: 0, line: 1}
	return &scanner
}

func (s *Scanner) NextToken() (*token.Token, error) {
	// Eat(ignore) all white spaces and comments, and leave the start pointer s.start just before the next token
	s.skipWsAndComments()

	// Did we consume all source string?
	if s.curr >= len(s.src)-1 {
		return token.EOFToken, nil
	}

	// Check first 2 character symbols (maximum-munch)
	symbol := string(s.peek()) + string(s.peekNext())
	t, ok := token.LookupSymbol(symbol)
	if ok {
		s.advance()
		s.advance()
		return token.New(t, symbol), nil
	}

	// If lookup of 2 character symbol failed, we try 1 character symbol
	symbol = string(s.peek())
	t, ok = token.LookupSymbol(symbol)
	if ok {
		s.advance()
		return token.New(t, symbol), nil
	}

	var scanned string
	// Lastly, is this byte a letter? if so we invoke indentifier scanner.
	// in which case the result could be a keyword so we check for this condition.
	if token.IsLetter(s.peek()) || s.peek() == '_' {
		scanned = s.scanIdentifier()
		return token.New(token.LookupIdent(scanned), scanned), nil
	}

	// if the peeked character is a digit we scan for a number
	if token.IsDigit(s.peek()) {
		scanned = s.scanNumber()
		return token.New(token.NUMBER, scanned), nil
	}

	return nil, fmt.Errorf("could not match token")
}
