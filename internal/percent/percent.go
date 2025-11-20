// Package percent provides full RFC3986 percent-encoding/decoding
package percent

import (
	"fmt"
)

// Decode decodes percent-encoded sequences in a string
// Returns an error for malformed sequences (% not followed by [2] hex digits)
func Decode(s string) (string, error) {
	var result []byte
	i := 0

	for i < len(s) {
		c := s[i]
		if c == '%' {
			if i+2 >= len(s) {
				return "", fmt.Errorf("incomplete percent-encoded sequence at position %d", i)
			}

			hex1, hex2 := s[i+1], s[i+2]
			v1, ok1 := hexToByte(hex1)
			v2, ok2 := hexToByte(hex2)

			if !ok1 || !ok2 {
				return "", fmt.Errorf("invalid percent-encoded sequence %%%c%c at position %d", hex1, hex2, i)
			}

			result = append(result, v1<<4|v2)
			i += 3
		} else {
			result = append(result, c)
			i++
		}
	}

	return string(result), nil
}

// hexToByte converts a hex character to its byte value
func hexToByte(c byte) (byte, bool) {
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	default:
		return 0, false
	}
}
