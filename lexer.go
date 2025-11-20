package rfcquery

import (
	"unicode"

	"github.com/CRSylar/rfcquery/internal/percent"
)

// Lexer performs RFC3986 lexical analysis on query strings
type Lexer struct {
	input string
	pos   int
}

// NewLexer creates a new Lexer for the given query string
func NewLexer(input string) *Lexer {
	return &Lexer{input: input, pos: 0}
}

// isUnreserved returns true if the byte is an unreserved character per RFC3986
func isUnreserved(c byte) bool {
	return (c >= 'A' && c <= 'Z') ||
		(c >= 'a' && c <= 'z') ||
		(c >= '0' && c <= '9') ||
		(c == '-' || c == '.' || c == '_' || c == '~')
}

// isSubDelim returns true if the byte is a sub-delims character
func isSubDelim(c byte) bool {
	switch c {
	case '!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=':
		return true
	default:
		return false
	}
}

// isPcharOther returns true for ":" and "@"
func isPcharOther(c byte) bool {
	return c == ':' || c == '@'
}

// isPathChar returns true for '/' and '?'
func isPathChar(c byte) bool {
	return c == '/' || c == '?'
}

// Valid performs strict RFC3986 validation of the query string
// Returns nil if valid, or an error with position information
func (l *Lexer) Valid() error {
	i := 0
	for i < len(l.input) {
		c := l.input[i]

		if c == '%' {
			if i+2 >= len(l.input) {
				return newError(i, "incomplete percent-encoded sequence")
			}

			hex1, hex2 := l.input[i+1], l.input[i+2]
			if !isHexDigit(hex1) || !isHexDigit(hex2) {
				return newError(i, "invalid percent-encoded sequence %%%c%c", hex1, hex2)
			}

			i += 3
			continue
		}

		if isUnreserved(c) || isSubDelim(c) ||
			isPcharOther(c) || isPathChar(c) {
			i++
			continue
		}

		if c > unicode.MaxASCII {
			return newError(i, "non-ASCII character %q not allowed ( must be percent-encoded)", c)
		}

		return newError(i, "invalid character %q in query string", c)
	}

	return nil
}

// isHexDigit returns true if the byte is a valid hexadecimal digit
func isHexDigit(c byte) bool {
	return (c >= '0' && c <= '9') ||
		(c >= 'A' && c <= 'F') ||
		(c >= 'a' && c <= 'f')
}

// Decode returns the fully decoded query string
// This performs both validation and percent-encoding
func (l *Lexer) Decode() (string, error) {
	if err := l.Valid(); err != nil {
		return "", err
	}

	return percent.Decode(l.input)
}
