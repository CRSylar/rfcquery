package rfcquery

import (
	"unicode"

	"github.com/CRSylar/rfcquery/internal/percent"
)

// Scanner provides a token-by-token access to the query string
type Scanner struct {
	input string
	pos   int
	// Allow lookahead without consuming
	nextToken *Token
	nextErr   error
}

// NewScanner creates a new scanner for the query string
func NewScanner(input string) *Scanner {
	return &Scanner{
		input: input,
		pos:   0,
	}
}

// Valid performs full validation without tokenizing
func (s *Scanner) Valid() error {
	l := NewLexer(s.input)
	return l.Valid()
}

// Nextoken returns the next token and advances the scanner
func (s *Scanner) NextToken() (Token, error) {
	if s.nextToken != nil {
		tok := *s.nextToken
		s.nextToken = nil
		s.nextErr = nil
		return tok, nil
	}

	return s.scanToken()
}

// PeekToken returns the next token without advancing the scanner
func (s *Scanner) PeekToken() (Token, error) {
	if s.nextToken == nil && s.nextErr == nil {
		tok, err := s.scanToken()
		s.nextToken = &tok
		s.nextErr = err
	}

	if s.nextErr != nil {
		return Token{}, s.nextErr
	}

	return *s.nextToken, nil
}

func (s *Scanner) Pos() int {
	return s.pos
}

func (s *Scanner) scanToken() (Token, error) {
	startPos := s.pos

	if s.pos >= len(s.input) {
		return Token{
			Type:  TokenEOF,
			Value: "",
			Start: Position{Offset: startPos},
			End:   Position{Offset: startPos},
		}, nil
	}

	c := s.input[s.pos]

	if c == '%' {
		if s.pos+2 >= len(s.input) {
			return Token{}, newError(s.pos, "incomplete percent-encoded sequence")
		}

		hex1, hex2 := s.input[s.pos+1], s.input[s.pos+2]
		if !isHexDigit(hex1) || !isHexDigit(hex2) {
			return Token{}, newError(s.pos, "invalid percent-encoded sequence %%%c%c", hex1, hex2)
		}

		// Decode the sequence
		encoded := s.input[s.pos : s.pos+3]
		decoded, err := percent.Decode(encoded)
		if err != nil {
			return Token{}, newError(s.pos, "invalid percent-encoded sequence")
		}

		tok := Token{
			Type:    TokenPercentEncoded,
			Value:   encoded,
			Decoded: decoded,
			Start:   Position{Offset: s.pos},
			End:     Position{Offset: s.pos + 3},
		}
		s.pos += 3
		return tok, nil
	}

	var tokenType TokenType
	if isUnreserved(c) {
		tokenType = TokenUnreserved
	} else if isSubDelim(c) {
		tokenType = TokenSubDelims
	} else if isPcharOther(c) {
		tokenType = TokenPcharOther
	} else if isPathChar(c) {
		tokenType = TokenPathChar
	} else {
		if c > unicode.MaxASCII {
			return Token{}, newError(s.pos, "non-ASCII character %q must be percent-encoded", c)
		}
		return Token{}, newError(s.pos, "invalid character %q in query string", c)
	}

	tok := Token{
		Type:  tokenType,
		Value: s.input[s.pos : s.pos+1],
		Start: Position{Offset: s.pos},
		End:   Position{Offset: s.pos + 1},
	}
	s.pos++

	return tok, nil
}

func (s *Scanner) Rewind(n int) {
	s.pos -= n
	if s.pos < 0 {
		s.pos = 0
	}
	s.nextToken = nil
	s.nextErr = nil
}

func (s *Scanner) Reset() {
	s.pos = 0
	s.nextToken = nil
	s.nextErr = nil
}

// CollectAll reads all remaining tokens into a slice
func (s *Scanner) CollectAll() (TokenSlice, error) {
	var ts TokenSlice

	for {
		tok, err := s.NextToken()
		if err != nil {
			return nil, err
		}

		if tok.Type == TokenEOF {
			break
		}

		ts = append(ts, tok)
	}

	return ts, nil
}

// CollectWhile collects tokens while the predicate returns true
// The Token that fails the predicate is left unconsumed
func (s *Scanner) CollectWhile(predicate func(Token) bool) (TokenSlice, error) {
	var ts TokenSlice

	for {
		tok, err := s.PeekToken()
		if err != nil {
			return nil, err
		}

		if tok.Type == TokenEOF || !predicate(tok) {
			break
		}
		// consume the token ( is the same we already peeked)
		s.NextToken()
		ts = append(ts, tok)
	}

	return ts, nil
}

// CollectUntil collects tokens until the preciate returns true
// The Token that fails the predicate is left unconsumed
func (s *Scanner) CollectUntil(predicate func(Token) bool) (TokenSlice, error) {
	var ts TokenSlice

	for {
		tok, err := s.PeekToken()
		if err != nil {
			return nil, err
		}

		if tok.Type == TokenEOF || predicate(tok) {
			break
		}
		// consume the token ( is the same we already peeked)
		s.NextToken()
		ts = append(ts, tok)
	}

	return ts, nil
}

// CollectN collects exactly n tokens
// Returns error if fewer than n tokens are available
func (s *Scanner) CollectN(n int) (TokenSlice, error) {
	var ts TokenSlice

	for i := range n {
		tok, err := s.NextToken()
		if err != nil {
			return nil, err
		}

		if tok.Type == TokenEOF {
			return nil, newError(s.pos, "unexpected EOF, expected %d more tokens", n-i)
		}

		ts = append(ts, tok)
	}
	return ts, nil
}

// SkipWhile skips tokens while the predicate returns true
// Returns the number of skipped tokens
func (s *Scanner) SkipWhile(predicate func(Token) bool) (int, error) {
	count := 0

	for {
		tok, err := s.PeekToken()
		if err != nil {
			return 0, err
		}

		if tok.Type == TokenEOF || !predicate(tok) {
			break
		}

		s.NextToken()
		count++
	}

	return count, nil
}

// PeekN returns the next n tokens without consuming them
func (s *Scanner) PeekN(n int) (TokenSlice, error) {
	savedPos := s.pos

	tokens, err := s.CollectN(n)

	s.pos = savedPos
	s.nextToken = nil
	s.nextErr = nil

	return tokens, err
}
