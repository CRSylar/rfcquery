package formurlencoded_test

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/CRSylar/rfcquery"
	formurlencoded "github.com/CRSylar/rfcquery/plugins/form_urlencoded"
)

func TestFormURLEncodedParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    map[string][]string
		wantErr bool
	}{
		{
			name:  "simple key-value",
			input: "key=value",
			want: map[string][]string{
				"key": {"value"},
			},
		},
		{
			name:  "multiple pairs",
			input: "a=1&b=2&c=3",
			want: map[string][]string{
				"a": {"1"},
				"b": {"2"},
				"c": {"3"},
			},
		},
		{
			name:  "duplicate keys",
			input: "tag=go&tag=library&tag=rfc3986",
			want: map[string][]string{
				"tag": {"go", "library", "rfc3986"},
			},
		},
		{
			name:  "percent-encoded values",
			input: "name=John%20Doe&city=New%20York",
			want: map[string][]string{
				"name": {"John Doe"},
				"city": {"New York"},
			},
		},
		{
			name:  "special characters in keys",
			input: "filter:name=test&sort=created@asc",
			want: map[string][]string{
				"filter:name": {"test"},
				"sort":        {"created@asc"},
			},
		},
		{
			name:  "empty value",
			input: "key=",
			want: map[string][]string{
				"key": {""},
			},
		},
		{
			name:  "complex nesting simulation",
			input: "filter%5Bage%5D=25&filter%5Bname%5D=John",
			want: map[string][]string{
				"filter[age]":  {"25"},
				"filter[name]": {"John"},
			},
		},
		{
			name:  "sub-delimiters in values",
			input: "query=a!b$c&data=x*y+z",
			want: map[string][]string{
				"query": {"a!b$c"},
				"data":  {"x*y+z"},
			},
		},
		{
			name:  "path-style query",
			input: "path/to/file?search=value",
			want: map[string][]string{
				"path/to/file?search": {"value"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &formurlencoded.FormURLEncodedParser{}
			scanner := rfcquery.NewScanner(tt.input)

			result, err := parser.Parse(scanner)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			values := result.(*rfcquery.Values)

			got := make(map[string][]string)
			for _, key := range values.AllKeys() {
				vals := values.Get(key)
				for _, v := range vals {
					got[key] = append(got[key], v.Value)
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStdlibCompatibility(t *testing.T) {
	// test cases that should behave identically to stdlib
	testCases := []string{
		"key=value",
		"a=1&b=2&c=3",
		"tag=go&tag=library",
		"name=John%20Doe",
		"key=",
		"key",
		"%41=%42",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			stdLibVals, err := url.ParseQuery(tc)
			if err != nil {
				t.Fatalf("stdlib ParseQuery failed: %v", err)
			}

			rfcValues, err := formurlencoded.ParseFormURLEncoded(tc)
			if err != nil {
				t.Fatalf("rfcQuery ParseFormURLEncoded failed: %v", err)
			}

			for key := range stdLibVals {
				stdLibValues := stdLibVals[key]
				rfcValues := rfcValues.Get(key)

				if len(stdLibValues) != len(rfcValues) {
					t.Errorf("Key %s: stdlib has %d values, rfcquery has %d", key, len(stdLibValues), len(rfcValues))
					continue
				}

				for i, sstdLibVal := range stdLibValues {
					if rfcValues[i].Value != sstdLibVal {
						t.Errorf("key %s[%d]: stdlib=%q, rfcquery=%q", key, i, sstdLibVal, rfcValues[i].Value)
					}
				}
			}
		})
	}
}
