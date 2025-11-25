package graphql_test

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"

	"github.com/CRSylar/rfcquery"
	"github.com/CRSylar/rfcquery/plugins/graphql"
)

func TestGraphQLParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *graphql.GraphQLQuery
		wantErr bool
	}{
		{
			name:  "simple query",
			input: fmt.Sprintf(`query=%s`, url.QueryEscape("{user{name}}")),
			want: &graphql.GraphQLQuery{
				Query: `{user{name}}`,
			},
		},
		{
			name:  "query with variables",
			input: fmt.Sprintf(`query=%s&variables=%s`, url.QueryEscape("{user(id:$id){name}}"), url.QueryEscape(`{"id":"123"}`)),
			want: &graphql.GraphQLQuery{
				Query: `{user(id:$id){name}}`,
				Variables: map[string]any{
					"id": "123",
				},
			},
		},
		{
			name:  "query with operation name",
			input: fmt.Sprintf(`query=%s&operationName=GetUser`, url.QueryEscape("queryGetUser{user{name}}")),
			want: &graphql.GraphQLQuery{
				Query:         `queryGetUser{user{name}}`,
				OperationName: "GetUser",
			},
		},
		{
			name:  "full GraphQL-over-HTTP",
			input: fmt.Sprintf(`query=%s%s&variables=%s&operationName=GetUser`, "query%20", url.QueryEscape("GetUser($id:ID!){user(id:$id){name}}"), url.QueryEscape(`{"id":"123"}`)),
			want: &graphql.GraphQLQuery{
				Query:         `query GetUser($id:ID!){user(id:$id){name}}`,
				Variables:     map[string]any{"id": "123"},
				OperationName: "GetUser",
			},
		},
		{
			name:  "percent-encoded query",
			input: `query=%7B%20user%20%7B%20name%20%7D%20%7D`, // { user { name } }
			want: &graphql.GraphQLQuery{
				Query: `{ user { name } }`,
			},
		},
		{
			name:  "complex query with special chars",
			input: fmt.Sprintf(`query=%s`, url.QueryEscape(`{search(filter:{field:"value"}){results}}`)),
			want: &graphql.GraphQLQuery{
				Query: `{search(filter:{field:"value"}){results}}`,
			},
		},
		{
			name:    "missing query parameter",
			input:   fmt.Sprintf(`variables=%s`, url.QueryEscape(`{"id":"123"}`)),
			wantErr: true,
		},
		{
			name:    "multiple query parameters",
			input:   `query=foo&query=bar`,
			wantErr: true,
		},
		{
			name:    "invalid variables JSON",
			input:   fmt.Sprintf(`query=%s&variables=%s`, url.QueryEscape("{ foo }"), url.QueryEscape("{broken json}")),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := graphql.NewGraphQLParser()
			scanner := rfcquery.NewScanner(tt.input)

			result, err := parser.Parse(scanner)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, WantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			gpql, ok := result.(*graphql.GraphQLQuery)
			if !ok {
				t.Fatalf("expected *GraphQLQuery, got %T", result)
			}

			if gpql.Query != tt.want.Query {
				t.Errorf("Query = %q, want %q", gpql.Query, tt.want.Query)
			}

			if tt.want.OperationName != "" && gpql.OperationName != tt.want.OperationName {
				t.Errorf("OperationName = %q, want %q", gpql.OperationName, tt.want.OperationName)
			}

			if tt.want.Variables != nil {
				if !reflect.DeepEqual(gpql.Variables, tt.want.Variables) {
					t.Errorf("Variables = %v, want %v", gpql.Variables, tt.want.Variables)
				}
			}
		})
	}
}

func TestGraphQLParser_CustomTargetParam(t *testing.T) {
	input := fmt.Sprintf(`gql=%s`, url.QueryEscape(`{user{name}}`))

	parser := graphql.NewGraphQLParser()
	parser.TargetParam = "gql"

	scanner := rfcquery.NewScanner(input)
	result, err := parser.Parse(scanner)
	if err != nil {
		t.Fatalf("Parse() failed: %v", err)
	}

	gpql := result.(*graphql.GraphQLQuery)
	if gpql.Query != `{user{name}}` {
		t.Errorf("unexpected query: %q", gpql.Query)
	}
}

func TestGraphQLParser_DisableFeatures(t *testing.T) {

	input := fmt.Sprintf(`query=%s&variables=%s&operationName=GetUser`, url.QueryEscape(`{ user { name } }`), url.QueryEscape(`{"id": "123"}`))

	parser := graphql.NewGraphQLParser()
	parser.ParseVariables = false
	parser.ParseOperationName = false

	scanner := rfcquery.NewScanner(input)
	result, err := parser.Parse(scanner)
	if err != nil {
		t.Fatalf("Parse() failed: %v", err)
	}

	graphql := result.(*graphql.GraphQLQuery)

	if graphql.Query == "" {
		t.Error("Query should be parsed")
	}

	if graphql.Variables != nil {
		t.Error("Variables should not be parsed when disabled")
	}

	if graphql.OperationName != "" {
		t.Error("OperationName should not be parsed when disabled")
	}
}

func TestGraphQLParser_RFC3986Advantage(t *testing.T) {
	// GraphQL queries often contain special characters that break stdlib
	tests := []string{
		fmt.Sprintf(`query=%s`, url.QueryEscape(`{ user(id: "123") { field @include(if: $show) } }`)),   // @ character
		fmt.Sprintf(`query=%s`, url.QueryEscape(`{ search(filter: { path: "/api/v1" }) { results } }`)), // / character
		fmt.Sprintf(`query=%s`, url.QueryEscape(`{ user(email: "test@example.com") { id } }`)),          // @ character
		fmt.Sprintf(`query=%s`, url.QueryEscape(`query GetUser($id: ID!) { user(id: $id) { name } }`)),  // ! character
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			// Should parse successfully
			graphql, err := graphql.ParseGraphQLQuery(input)
			if err != nil {
				t.Errorf("rfcquery should parse %q but got: %v", input, err)
			}

			if graphql.Query == "" {
				t.Error("Query should not be empty")
			}
		})
	}
}
