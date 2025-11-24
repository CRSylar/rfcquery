package jsoninquery

import (
	"encoding/json"
	"fmt"

	"github.com/CRSylar/rfcquery"
	formurlencoded "github.com/CRSylar/rfcquery/plugins/form_urlencoded"
)

// JSONParser extracts and parses JSON values from query parameters
type JSONParser struct {
	// TargetParam specifies which parameter to extract JSON from
	// if empty, parses the entire query string as JSON
	TargetParam string

	// AllowMultiple allows multiple values for the target parameter
	// If false, returns error on duplicates
	AllowMultiple bool

	// StrictValidation requires valid RFC3986 before JSON parsing
	StrictValidation bool
}

// Name returns the parser identifier
func (p *JSONParser) Name() string {
	if p.TargetParam != "" {
		return fmt.Sprintf("json-in-query[%s]", p.TargetParam)
	}
	return "json-in-query"
}

// Parse implements the Parser interface
func (p *JSONParser) Parse(scanner *rfcquery.Scanner) (any, error) {
	if p.StrictValidation {
		if err := scanner.Valid(); err != nil {
			return nil, fmt.Errorf("validation failed: %w", err)
		}
		scanner.Reset()
	}

	if p.TargetParam == "" {
		return p.parseEntireQuery(scanner)
	}

	return p.parseTargetParam(scanner)
}

// parseEntireQuery treats the whole query string as JSON
func (p *JSONParser) parseEntireQuery(scanner *rfcquery.Scanner) (any, error) {
	tokens, err := scanner.CollectAll()
	if err != nil {
		return nil, fmt.Errorf("failed to collect query: %w", err)
	}

	jsonStr := tokens.StringDecoded()
	var result any
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("invalid JSON in query: %w", err)
	}

	return result, nil
}

// parseTargetParam extracts JSON from a specific parameter
func (p *JSONParser) parseTargetParam(scanner *rfcquery.Scanner) (map[string]any, error) {
	formParser := &formurlencoded.FormURLEncodedParser{}
	result, err := formParser.Parse(scanner)
	if err != nil {
		return nil, fmt.Errorf("failed to parse as form-urlencoded: %w", err)
	}

	values := result.(*rfcquery.Values)
	targetValues := values.Get(p.TargetParam)

	if len(targetValues) == 0 {
		return nil, fmt.Errorf("target parameter %q not found", p.TargetParam)
	}

	if len(targetValues) > 1 && !p.AllowMultiple {
		return nil, fmt.Errorf("multiple values found for parameter %q", p.TargetParam)
	}

	results := make(map[string]any)
	for i, val := range targetValues {
		var jsonData any
		if err := json.Unmarshal([]byte(val.Value), &jsonData); err != nil {
			return nil, fmt.Errorf("invalid JSON in parameter %q (value %d): %w", p.TargetParam, i, err)
		}

		key := p.TargetParam
		if len(targetValues) > 1 {
			key = fmt.Sprintf("%s[%d]", p.TargetParam, i)
		}
		results[key] = jsonData
	}

	if len(results) == 1 {
		for _, v := range results {
			return map[string]any{p.TargetParam: v}, nil
		}
	}
	return results, nil
}

func ParseJSONQuery(query string, targetParam string) (any, error) {
	scanner := rfcquery.NewScanner(query)
	if err := scanner.Valid(); err != nil {
		return nil, err
	}

	parser := &JSONParser{
		TargetParam:      targetParam,
		AllowMultiple:    false,
		StrictValidation: true,
	}

	result, err := parser.Parse(scanner)
	if err != nil {
		return nil, err
	}

	return result, nil
}
