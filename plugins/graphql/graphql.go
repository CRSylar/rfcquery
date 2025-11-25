package graphql

import (
	"encoding/json"
	"fmt"

	"github.com/CRSylar/rfcquery"
	formurlencoded "github.com/CRSylar/rfcquery/plugins/form_urlencoded"
)

// GraphQLQuery represents a GraphQL request from URL parameters
type GraphQLQuery struct {
	// The query document
	Query string

	// Optional Operation name to support multiple operations in document
	OperationName string

	// Optional variables JSON object
	Variables map[string]any

	// Raw tokens for metadata
	QueryTokens     rfcquery.TokenSlice
	VariablesTokens rfcquery.TokenSlice
	OperationTokens rfcquery.TokenSlice
}

// GraphQLParser extracts and validates GraphQL query from URL params
type GraphQLParser struct {

	// TargetParam is the paramter name for the GraphQL query ( default: "query")
	TargetParam string

	// ParseVariables enables parsing the "variables" parameter as JSON
	ParseVariables bool

	// ParseOperationName enable parsing the <operation_name> parameter
	ParseOperationName bool

	// Strict validate using the RFC
	StrictValidation bool
}

// NewGraphQLParser crates a parser with default setting
func NewGraphQLParser() *GraphQLParser {
	return &GraphQLParser{
		TargetParam:        "query",
		ParseVariables:     true,
		ParseOperationName: true,
		StrictValidation:   true,
	}
}

// Name returns the parser indentifier
func (p *GraphQLParser) Name() string {
	return "graphql-over-http"
}

func (p *GraphQLParser) Parse(scanner *rfcquery.Scanner) (any, error) {
	if p.StrictValidation {
		if err := scanner.Valid(); err != nil {
			return nil, fmt.Errorf("RFC3986 validation failed: %w", err)
		}
		scanner.Reset()
	}

	formParser := &formurlencoded.FormURLEncodedParser{}
	result, err := formParser.Parse(scanner)
	if err != nil {
		return nil, fmt.Errorf("failed to parse query paramters: %w", err)
	}

	values := result.(*rfcquery.Values)
	graphql := &GraphQLQuery{}

	queryVals := values.Get(p.TargetParam)
	if len(queryVals) == 0 {
		return nil, fmt.Errorf("GrapQL query parameter %q not found", p.TargetParam)
	}
	if len(queryVals) > 1 {
		return nil, fmt.Errorf("Multiple values found for GrapQL query paramter %q", p.TargetParam)
	}

	graphql.Query = queryVals[0].Value
	graphql.QueryTokens = queryVals[0].ValueTokens

	if p.ParseVariables {
		if err := p.parseVariables(values, graphql); err != nil {
			return nil, err
		}
	}

	if p.ParseOperationName {
		if err := p.parseOperationName(values, graphql); err != nil {
			return nil, err
		}
	}

	return graphql, nil
}

func (p *GraphQLParser) parseVariables(values *rfcquery.Values, query *GraphQLQuery) error {
	varVals := values.Get("variables")
	if len(varVals) == 0 {
		return nil
	}

	if len(varVals) > 1 {
		return fmt.Errorf("Multiple values found for variables paramter")
	}

	query.VariablesTokens = varVals[0].ValueTokens

	if err := json.Unmarshal([]byte(varVals[0].Value), &query.Variables); err != nil {
		return fmt.Errorf("invalid JSON in variables parameter: %w", err)
	}

	return nil
}

func (p *GraphQLParser) parseOperationName(values *rfcquery.Values, query *GraphQLQuery) error {
	opVals := values.Get("operationName")
	if len(opVals) == 0 {
		return nil
	}

	if len(opVals) > 1 {
		return fmt.Errorf("multiple values found for operationName parameter")
	}

	query.OperationName = opVals[0].Value
	query.OperationTokens = opVals[0].ValueTokens

	return nil
}

func ParseGraphQLQuery(query string) (*GraphQLQuery, error) {
	scanner := rfcquery.NewScanner(query)
	if err := scanner.Valid(); err != nil {
		return nil, err
	}

	parser := NewGraphQLParser()
	result, err := parser.Parse(scanner)
	if err != nil {
		return nil, err
	}

	parsed, ok := result.(*GraphQLQuery)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return parsed, nil
}
