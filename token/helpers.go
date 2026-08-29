package token

// Character classes used by the scanner to decide which strategy to run.
//
// These are deliberately positive tests — "is this byte inside the set I care
// about" — so an unknown byte, including the scanner's end-of-input sentinel,
// is rejected automatically.

// IsLetter reports whether b can appear in an identifier as a letter.
// ASCII only: identifiers in this language are deliberately not Unicode.
func IsLetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

// IsDigit reports whether b is an ASCII decimal digit.
func IsDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
