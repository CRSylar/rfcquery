package percent

import "testing"

func TestDecode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"simple", "hello", "hello", false},
		{"space encoded", "hello%20world", "hello world", false},
		{"mixed case hex", "test%2A%2b%2f", "test*+/", false},
		{"malformed incomplete", "test%", "", true},
		{"malformed invalid hex", "test%GG", "", true},
		{"invalid hex digit", "test%1G", "", true},
		{"empty string", "", "", false},
		{"single percent", "%", "", true},
		{"double percent", "%%", "", true},
		{"percent at end", "test%", "", true},
		{"percent with one digit", "test%1", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}
