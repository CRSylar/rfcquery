package jsoninquery_test

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/CRSylar/rfcquery"
	jsoninquery "github.com/CRSylar/rfcquery/plugins/json_in_query"
)

func TestJSONParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		target  string
		want    any
		wantErr bool
	}{
		{
			name:   "simple object in parameter",
			input:  `filter=%7B%22name%22:%22John%22,%22age%22:30%7D`,
			target: "filter",
			want: map[string]any{
				"filter": map[string]any{
					"name": "John",
					"age":  float64(30),
				},
			},
		},
		{
			name:   "array in paramter",
			input:  `data=%5B%22a%22,%22b%22,%22c%22%5D`,
			target: "data",
			want: map[string]any{
				"data": []any{"a", "b", "c"},
			},
		},
		{
			name:   "null value",
			input:  `value=null`,
			target: "value",
			want: map[string]any{
				"value": nil,
			},
		},
		{
			name:   "boolean and number",
			input:  `flags=%7B%22active%22:true,%22count%22:42%7D`,
			target: "flags",
			want: map[string]any{
				"flags": map[string]any{
					"active": true,
					"count":  float64(42),
				},
			},
		},
		{
			name:   "nested objects",
			input:  `config=%7B%22user%22:%7B%22name%22:%22Alice%22%7D,%22enabled%22:true%7D`,
			target: "config",
			want: map[string]any{
				"config": map[string]any{
					"user": map[string]any{
						"name": "Alice",
					},
					"enabled": true,
				},
			},
		},
		{
			name:   "percent-encoded JSON",
			input:  `filter=%7B%22name%22%3A%22John%22%7D`,
			target: "filter",
			want: map[string]any{
				"filter": map[string]any{
					"name": "John",
				},
			},
		},
		{
			name:   "multiple parameters with one JSON",
			input:  `type=user&filter=%7B%22active%22:true%7D&limit=10`,
			target: "filter",
			want: map[string]any{
				"filter": map[string]any{
					"active": true,
				},
			},
		},
		{
			name:    "parameter not found",
			input:   `foo=bar`,
			target:  "missing",
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			input:   `filter=%7Bbroken json%7D`,
			target:  "filter",
			wantErr: true,
		},
		{
			name:    "empty paramter value",
			input:   `filter=`,
			target:  "filter",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			scanner := rfcquery.NewScanner(tt.input)
			parser := &jsoninquery.JSONParser{
				TargetParam:      tt.target,
				StrictValidation: true,
			}

			result, err := parser.Parse(scanner)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("Parse() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestJSONParse_MultipleValues(t *testing.T) {
	input := `data=%5B%22first%22%5D&data=%5B%22second%22%2C%22third%22%5D`

	parser := &jsoninquery.JSONParser{
		TargetParam:   "data",
		AllowMultiple: true,
	}

	scanner := rfcquery.NewScanner(input)
	result, err := parser.Parse(scanner)
	if err != nil {
		t.Fatalf("Parse() with AllowMultiple=true failed: %v", err)
	}

	resultMap := result.(map[string]any)
	if len(resultMap) != 2 {
		t.Errorf("Expected error when AllowMultiple=false but got none")
	}
}

func TestJSONParser_EntireQuery(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    any
		wantErr bool
	}{
		{
			name:  "query in JSON object",
			input: `%7B%22name%22%3A%22John%22,%22age%22%3A30%7D`,
			want: map[string]any{
				"name": "John",
				"age":  float64(30),
			},
		},
		{
			name:  "query is JSON array",
			input: `%5B%22apple%22,%22banana%22,%22cherry%22%5D`,
			want:  []any{"apple", "banana", "cherry"},
		},
		{
			name:  "query is nested JSON",
			input: `%7B%22users%22%3A%5B%7B%22id%22%3A1%7D,%7B%22id%22%3A2%7D%5D%7D`,
			want: map[string]any{
				"users": []any{
					map[string]any{"id": float64(1)},
					map[string]any{"id": float64(2)},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &jsoninquery.JSONParser{
				TargetParam:      "",
				StrictValidation: true,
			}

			scanner := rfcquery.NewScanner(tt.input)
			result, err := parser.Parse(scanner)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("Parse() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestJSONParser_StdlibComparison(t *testing.T) {
	tests := []struct {
		key   string
		input string
	}{
		{
			key:   "filter",
			input: `filter=%7B%22range%22%3A%7B%22age%22:%7B%22gt%22:25%7D%7D%7D`,
		}, {
			key:   "config",
			input: `config=%7B%22server%22%3A%22localhost%3A8080%22%2C%22path%22%3A%22/api/v1%22%7D`,
		}, {
			key:   "data",
			input: `data=%7B%22url%22%3A%22https%3A//example.com?search=test%22%7D`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {

			parser := &jsoninquery.JSONParser{
				TargetParam: tt.key,
			}

			scanner := rfcquery.NewScanner(tt.input)

			_, err := parser.Parse(scanner)
			if err != nil {
				t.Errorf("rfcquery should parse %q but got: %v", tt, err)
			}

			// stdlib should fail these due to special chars
			_, stdlibErr := url.ParseQuery(tt.input)
			if stdlibErr == nil {
				t.Logf("stdlib parsed %q but would mangle special chars", tt)
			}
		})
	}
}

func BenchmarkJSONParser(b *testing.B) {
	input := `filter=%7B%22name%22:%22John Doe%22,%22age%22:30,%22tags%22:%5B%22go%22,%22json%22%5D%7D&limit=10`

	b.ResetTimer()
	for b.Loop() {

		scanner := rfcquery.NewScanner(input)
		parser := &jsoninquery.JSONParser{TargetParam: "filter"}
		_, err := parser.Parse(scanner)
		if err != nil {
			b.Fatal(err)
		}
	}
}
