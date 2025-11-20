package rfcquery

import (
	"testing"
)

func TestScannerNextToken(t *testing.T) {
	type wantTokens struct {
		tokenType TokenType
		value     string
		decoded   string
	}
	tests := []struct {
		name       string
		input      string
		wantTokens []wantTokens
		wantErr    bool
	}{
		{
			name:  "simple key-value",
			input: "key=value",
			wantTokens: []wantTokens{
				{tokenType: TokenUnreserved, value: "k", decoded: ""},
				{tokenType: TokenUnreserved, value: "e", decoded: ""},
				{tokenType: TokenUnreserved, value: "y", decoded: ""},
				{tokenType: TokenSubDelims, value: "=", decoded: ""},
				{tokenType: TokenUnreserved, value: "v", decoded: ""},
				{tokenType: TokenUnreserved, value: "a", decoded: ""},
				{tokenType: TokenUnreserved, value: "l", decoded: ""},
				{tokenType: TokenUnreserved, value: "u", decoded: ""},
				{tokenType: TokenUnreserved, value: "e", decoded: ""},
			},
		},
		{
			name:  "percent-encoded space",
			input: "hello%20world",
			wantTokens: []wantTokens{
				{tokenType: TokenUnreserved, value: "h", decoded: ""},
				{tokenType: TokenUnreserved, value: "e", decoded: ""},
				{tokenType: TokenUnreserved, value: "l", decoded: ""},
				{tokenType: TokenUnreserved, value: "l", decoded: ""},
				{tokenType: TokenUnreserved, value: "o", decoded: ""},
				{tokenType: TokenPercentEncoded, value: "%20", decoded: " "},
				{tokenType: TokenUnreserved, value: "w", decoded: ""},
				{tokenType: TokenUnreserved, value: "o", decoded: ""},
				{tokenType: TokenUnreserved, value: "r", decoded: ""},
				{tokenType: TokenUnreserved, value: "l", decoded: ""},
				{tokenType: TokenUnreserved, value: "d", decoded: ""},
			},
		},
		{
			name:       "special characters - '>' not urlEncoded",
			input:      "filter:age>25",
			wantTokens: []wantTokens{},
			wantErr:    true,
		},
		{
			name:  "path-style query",
			input: "path/to/file?search",
			wantTokens: []wantTokens{
				{tokenType: TokenUnreserved, value: "p", decoded: ""},
				{tokenType: TokenUnreserved, value: "a", decoded: ""},
				{tokenType: TokenUnreserved, value: "t", decoded: ""},
				{tokenType: TokenUnreserved, value: "h", decoded: ""},
				{tokenType: TokenPathChar, value: "/", decoded: ""},
				{tokenType: TokenUnreserved, value: "t", decoded: ""},
				{tokenType: TokenUnreserved, value: "o", decoded: ""},
				{tokenType: TokenPathChar, value: "/", decoded: ""},
				{tokenType: TokenUnreserved, value: "f", decoded: ""},
				{tokenType: TokenUnreserved, value: "i", decoded: ""},
				{tokenType: TokenUnreserved, value: "l", decoded: ""},
				{tokenType: TokenUnreserved, value: "e", decoded: ""},
				{tokenType: TokenPathChar, value: "?", decoded: ""},
				{tokenType: TokenUnreserved, value: "s", decoded: ""},
				{tokenType: TokenUnreserved, value: "e", decoded: ""},
				{tokenType: TokenUnreserved, value: "a", decoded: ""},
				{tokenType: TokenUnreserved, value: "r", decoded: ""},
				{tokenType: TokenUnreserved, value: "c", decoded: ""},
				{tokenType: TokenUnreserved, value: "h", decoded: ""},
				{tokenType: TokenEOF, value: "", decoded: ""},
			},
		},
		{
			name:  "sub-delimiters",
			input: "a=!$&'()*+,;=",
			wantTokens: []wantTokens{
				{tokenType: TokenUnreserved, value: "a", decoded: ""},
				{tokenType: TokenSubDelims, value: "=", decoded: ""},
				{tokenType: TokenSubDelims, value: "!", decoded: ""},
				{tokenType: TokenSubDelims, value: "$", decoded: ""},
				{tokenType: TokenSubDelims, value: "&", decoded: ""},
				{tokenType: TokenSubDelims, value: "'", decoded: ""},
				{tokenType: TokenSubDelims, value: "(", decoded: ""},
				{tokenType: TokenSubDelims, value: ")", decoded: ""},
				{tokenType: TokenSubDelims, value: "*", decoded: ""},
				{tokenType: TokenSubDelims, value: "+", decoded: ""},
				{tokenType: TokenSubDelims, value: ",", decoded: ""},
				{tokenType: TokenSubDelims, value: ";", decoded: ""},
				{tokenType: TokenSubDelims, value: "=", decoded: ""},
				{tokenType: TokenEOF, value: "", decoded: ""},
			},
		},
		{
			name:  "percent encoded unicode",
			input: "emoji=%F0%9F%91%8D",
			wantTokens: []wantTokens{
				{tokenType: TokenUnreserved, value: "e", decoded: ""},
				{tokenType: TokenUnreserved, value: "m", decoded: ""},
				{tokenType: TokenUnreserved, value: "o", decoded: ""},
				{tokenType: TokenUnreserved, value: "j", decoded: ""},
				{tokenType: TokenUnreserved, value: "i", decoded: ""},
				{tokenType: TokenSubDelims, value: "=", decoded: ""},
				{tokenType: TokenPercentEncoded, value: "%F0", decoded: "\xF0"},
				{tokenType: TokenPercentEncoded, value: "%9F", decoded: "\x9F"},
				{tokenType: TokenPercentEncoded, value: "%91", decoded: "\x91"},
				{tokenType: TokenPercentEncoded, value: "%8D", decoded: "\x8D"},
				{tokenType: TokenEOF, value: "", decoded: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewScanner(tt.input)

			for i, want := range tt.wantTokens {
				tok, err := scanner.NextToken()

				if (err != nil) != tt.wantErr {
					t.Errorf("NextToken() at token %d error = %v, wantErr %v", i, err, tt.wantErr)
					return
				}

				if tt.wantErr {
					return
				}

				if tok.Type != want.tokenType {
					t.Errorf("token %d: Type = %v, want %v", i, tok.Type, want.tokenType)
				}

				if tok.Value != want.value {
					t.Errorf("token %d: Value = %q, want %q", i, tok.Value, want.value)
				}

				if tok.Decoded != want.decoded {
					t.Errorf("token %d: Decoded = %q, want %q", i, tok.Decoded, want.decoded)
				}
			}
		})
	}
}

func TestScannerPeekToken(t *testing.T) {
	scanner := NewScanner("key=value")

	tok1, err := scanner.PeekToken()
	if err != nil {
		t.Fatalf("PeekToken() error = %v", err)
	}

	tok2, err := scanner.PeekToken()
	if err != nil {
		t.Fatalf("PeekToken() error = %v", err)
	}

	if tok1 != tok2 {
		t.Errorf("PeekToken() should return same token, got %v and %v", tok1, tok2)
	}

	tok3, err := scanner.NextToken()
	if err != nil {
		t.Fatalf("NextToken() error = %v", err)
	}

	if tok3 != tok1 {
		t.Errorf("NextToken() = %v, want %v", tok3, tok1)
	}

	tok4, _ := scanner.PeekToken()
	if tok4.Type == tok3.Type && tok4.Value == tok3.Value {
		t.Errorf("PeekToken() after NextToken() should return different token")
	}
}

func TestScannerRewind(t *testing.T) {
	scanner := NewScanner("hello")

	tok1, _ := scanner.NextToken()
	tok2, _ := scanner.NextToken()

	scanner.Rewind(2)

	tok1again, _ := scanner.NextToken()
	tok2again, _ := scanner.NextToken()

	if tok1 != tok1again {
		t.Errorf("Rewind failed: fist token mismatch")
	}
	if tok2 != tok2again {
		t.Errorf("Rewind failed: second token mismatch")
	}
}

func TestScannerReset(t *testing.T) {
	scanner := NewScanner("test")

	for {
		tok, _ := scanner.NextToken()
		if tok.Type == TokenEOF {
			break
		}
	}

	scanner.Reset()

	tok, err := scanner.PeekToken()
	if err != nil {
		t.Fatalf("PeekToken() after Reset() error = %v", err)
	}

	if tok.Type != TokenUnreserved || tok.Value != "t" {
		t.Errorf("Reset() failed: first token = %v, want Unreserved 't'", tok)
	}
}

func TestScannerErrorPosition(t *testing.T) {
	scanner := NewScanner("test%GG")

	// Should error at pos 4
	_, err := scanner.NextToken() // 't'
	_, err = scanner.NextToken()  // 'e'
	_, err = scanner.NextToken()  // 's'
	_, err = scanner.NextToken()  // 't'
	_, err = scanner.NextToken()  // '%GG - error here'

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	rfcErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("expected *Error type, got %T", err)
	}

	if rfcErr.Pos.Offset != 4 {
		t.Errorf("expected position 4, got %d", rfcErr.Pos.Offset)
	}
}
