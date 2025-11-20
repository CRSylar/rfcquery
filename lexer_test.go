package rfcquery

import (
	"strings"
	"testing"
)

func TextLexerValid(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{"empty", "", false, ""},
		{"simple", "key=value", false, ""},
		{"space encoded", "name=John%20Doe", false, ""},
		{"special chars", "filters=name%3Dtest%26active%3Dtrue", false, ""},
		{"unreserved", "test-ABC_123.~xyz", false, ""},
		{"sub-delims", "a=!$&'()*+,;=", false, ""},
		{"pchar other", "user:pass@host", false, ""},
		{"path chars", "path/to/file?search", false, ""},
		{"mixed", "q=search&filter[age]=25", false, ""},

		// Invalid cases, expect fails
		{"space literal", "hello world", true, "invalid character"},
		{"quote", `test"value"`, true, "invalid character"},
		{"backslash", "test\\value", true, "invalid character"},
		{"incomplete percent encode", "test%", true, "incomplete percent-encoded"},
		{"invalid hex upper", "test%GG", true, "invalid percent-encoded"},
		{"invalid hex lower", "test%gg", true, "invalid percent-encoded"},
		{"percent at end", "test%2", true, "incomplete percent-encoded"},
		{"non-ascii", "text\xc3\xa9", true, "non-ascii character"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			err := l.Valid()

			if (err != nil) != tt.wantErr {
				t.Errorf("Valid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("Valid() error = %v, want error containing %q", err, tt.errMsg)
			}
		})
	}
}

func TestLexerDecode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"no encoding", "hello", "hello", false},
		{"space", "hello%20world", "hello world", false},
		{"mixed", "name=John%20Doe&age=30", "name=John Doe&age=30", false},
		{"unicode via encoding", "emoji=%F0%9F%91%8D", "emoji=👍", false},
		{"invalid", "test%GG", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			got, err := l.Decode()

			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("Decode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPositionInError(t *testing.T) {
	l := NewLexer("test%GGvalue")
	err := l.Valid()

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
