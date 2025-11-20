package rfcquery

import (
	"bytes"
	"strings"
)

type TokenType int

type TokenSlice []Token

// Token represents a lexical token in a query string
type Token struct {
	Type    TokenType
	Value   string // Raw value (percent-encoded if present)
	Decoded string // Decoded value for percent-encoded tokens
	Start   Position
	End     Position
}

const (
	TokenInvalid        TokenType = iota
	TokenPercentEncoded           // %HH sequence
	TokenUnreserved               // ALPHA / DIGIT/ - / . / _  / ~
	TokenSubDelims                // ! / $ / & / ' / ( / ) / * / + / , / ; / =
	TokenPcharOther               // : / @
	TokenPathChar                 // '/' / ?
	TokenEOF
)

// String returns a readable representation of the token type
// this implement the Stringer interface
func (tt TokenType) String() string {
	switch tt {
	case TokenEOF:
		return "EOF"
	case TokenPercentEncoded:
		return "PercentEncoded"
	case TokenUnreserved:
		return "Unreserved"
	case TokenSubDelims:
		return "SubDelims"
	case TokenPcharOther:
		return "PCharOther"
	case TokenPathChar:
		return "PathChar"
	default:
		return "Invalid"
	}
}

// String reconstructs the original query string
// implemtation of Stringer interface
func (ts TokenSlice) String() string {
	var sb strings.Builder
	for _, tok := range ts {
		sb.WriteString(tok.Value)
	}
	return sb.String()
}

// StringDecoded reconstructs the fully decoded query string
func (ts TokenSlice) StringDecoded() string {
	var sb strings.Builder
	for _, tok := range ts {
		if tok.Type == TokenPercentEncoded {
			sb.WriteString(tok.Decoded)
		} else {
			sb.WriteString(tok.Value)
		}
	}

	return sb.String()
}

// Bytes returns the raw byte representation
func (ts TokenSlice) Bytes() []byte {
	var buf bytes.Buffer
	for _, tok := range ts {
		buf.WriteString(tok.Value)
	}
	return buf.Bytes()
}
