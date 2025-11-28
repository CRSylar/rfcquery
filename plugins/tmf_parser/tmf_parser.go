package tmfparser

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/CRSylar/rfcquery"
)

// TMFOperator represents a comparison operator in TMF syntax
type TMFOperator string

const (
	TMFOperatorEq  TMFOperator = "eq"  // = or %3D
	TMFOperatorGt  TMFOperator = "gt"  // %3E
	TMFOperatorLt  TMFOperator = "lt"  // %3C
	TMFOperatorGte TMFOperator = "gte" // %3E%3D
	TMFOperatorLte TMFOperator = "lte" // %3C%3D
	TMFOperatorNe  TMFOperator = "ne"  // %21%3D
)

// OperatorMap maps encoded operators to their string reprentation
var operatorMap = map[string]TMFOperator{
	"%3D":    TMFOperatorEq,
	"%3E":    TMFOperatorGt,
	".gt":    TMFOperatorGt,
	"%3C":    TMFOperatorLt,
	".lt":    TMFOperatorLt,
	"%3E%3D": TMFOperatorGte,
	".gte":   TMFOperatorGte,
	"%3C%3D": TMFOperatorLte,
	".lte":   TMFOperatorLte,
	"%21%3D": TMFOperatorNe,
}

type TMFExpression struct {
	Operator TMFOperator
	Value    string
	Token    rfcquery.TokenSlice
}

type TMFFilterGroup struct {
	Expressions []TMFExpression
	Tokens      rfcquery.TokenSlice
}

type TMFSortField struct {
	Field     string
	Direction string // "asc" or "desc"
	Tokens    rfcquery.TokenSlice
}

// TMFQuery represents a complete parsed TMF Query
type TMFQuery struct {
	Expressions map[string][]TMFExpression // Top-level filter expressions
	Sorting     []TMFSortField             // Parsed sort fields
	OtherParams map[string][]string        // Non-filter params ( limit, offset, etc...)
}

type TMFParser struct {
	// OperatorMap allows customizing recognized operators
	OperatorMap map[string]TMFOperator

	// StrictValidation enforces RFC3986 compliance
	StrictValidation bool

	// EnableGrouping enables support for parantheses grouping
	EnableGrouping bool
}

func NewTMFParser() *TMFParser {
	return &TMFParser{
		OperatorMap:      operatorMap,
		StrictValidation: true,
		EnableGrouping:   true,
	}
}

func (p *TMFParser) Name() string {
	return "tmf-query-parser"
}

func (p *TMFParser) Parse(scanner *rfcquery.Scanner) (any, error) {
	if p.StrictValidation {
		if err := scanner.Valid(); err != nil {
			return nil, fmt.Errorf("RFC3986 validation failed: %w", err)
		}
		scanner.Reset()
	}

	query := &TMFQuery{
		Expressions: make(map[string][]TMFExpression),
		Sorting:     make([]TMFSortField, 0),
		OtherParams: make(map[string][]string),
	}

	if err := p.parseQuery(scanner, query); err != nil {
		return nil, err
	}

	return query, nil
}

func (p *TMFParser) parseQuery(scanner *rfcquery.Scanner, query *TMFQuery) error {
outer:
	for {
		tk, err := scanner.PeekToken()
		if err != nil {
			return err
		}

		if tk.Type == rfcquery.TokenEOF {
			break outer
		}

		segment, err := p.parseSegment(scanner)
		if err != nil {
			return err
		}

		if segment.Key == "sort" {
			sortFields, err := p.parseSortValue(segment.ValueTokens)
			if err != nil {
				return fmt.Errorf("invalid sort syntax: %w", err)
			}
			query.Sorting = append(query.Sorting, sortFields...)
			continue
		}

		if p.isFilterSegment(segment.Key) {
			if segment.Expressions != nil {
				query.Expressions[segment.Key] = append(query.Expressions[segment.Key], *segment.Expressions)
			}
		} else {
			query.OtherParams[segment.Key] = append(query.OtherParams[segment.Key], segment.Value...)
		}

		tok, err := scanner.PeekToken()
		if err != nil {
			return err
		}
		if tok.Type == rfcquery.TokenSubDelims && (tok.Value == ";" || tok.Value == "&") {
			scanner.NextToken() // consume separator
		}
	}

	return nil
}

type segmentResult struct {
	Key         string
	Value       []string
	ValueTokens rfcquery.TokenSlice
	Expressions []TMFExpression
}

func (p *TMFParser) parseSegment(scanner *rfcquery.Scanner) (*segmentResult, error) {
	keyTokens, err := scanner.CollectUntil(func(t rfcquery.Token) bool {
		separator := t.Type == rfcquery.TokenSubDelims && (t.Value == "=" || t.Value == ";" || t.Value == "&")

		// 																		checking operators    =										>									<											<=										>=											!=
		operator := t.Type == rfcquery.TokenPercentEncoded && (t.Value == "%3D" || t.Value == "%3E" || t.Value == "%3C" || t.Value == "%3E%3D" || t.Value == "%3C%3D" || t.Value == "%21%3D")
		return separator || operator
	})
	if err != nil {
		return nil, err
	}

	sepToken, err := scanner.PeekToken()
	if err != nil {
		return nil, err
	}

	result := &segmentResult{
		Key: keyTokens.StringDecoded(),
	}

	if sepToken.Type == rfcquery.TokenSubDelims && (sepToken.Value == ";" || sepToken.Value == "&") {
		// TODO: perform some simulation, maybe here is more appropriate return an error since i've the key with no value ( or just skip it )
		result.Value = []string{}
		result.ValueTokens = rfcquery.TokenSlice{}

		result.Expressions = append(result.Expressions, TMFExpression{
			Operator: TMFOperatorEq,
			Value:    "",
			Token:    keyTokens,
		})
		return result, nil
	}

	if sepToken.Value == "=" {
		_, err := scanner.NextToken()
		if err != nil {
			return nil, err
		}
	}

	valueTokens, err := scanner.CollectUntil(func(t rfcquery.Token) bool {
		return t.Type == rfcquery.TokenSubDelims && (t.Value == ";" || t.Value == "&")
	})
	if err != nil {
		return nil, err
	}

	result.Value = strings.Split(valueTokens.StringDecoded(), ",")
	result.ValueTokens = valueTokens

	if p.isFilterSegment(keyTokens.StringDecoded()) {
		result.Expressions = append(result.Expressions, p.parseFilterValue(keyTokens.StringDecoded(), valueTokens)...)
	}

	return result, nil
}

func (p *TMFParser) isFilterSegment(key string) bool {
	return key != "" && key != "sort" && key != "limit" && key != "offset"
}

func (p *TMFParser) parseFilterValue(field string, tokens rfcquery.TokenSlice) []TMFExpression {
	slog.Info("parseFilterValue - received", "field", field, "tokens", tokens)
	currPos := 0

	results := []TMFExpression{}

	for currPos < len(tokens) {
		operatorFound := false

		for opStr, opType := range p.OperatorMap {
			opLen := len(opStr)

			slog.Info("currPos", "cp", currPos, "op", opStr, "tokens", tokens)
			if currPos+opLen <= len(tokens) {
				match := true
				for i := range opLen {
					slog.Info("ranging", "v", tokens[currPos+i].Value)
					if tokens[currPos+i].Value != string(opStr[i]) {
						match = false
						break
					}
				}

				slog.Info("match", "m", match)
				if match {
					valueTokens := tokens[currPos+opLen:]

					results = append(results, TMFExpression{
						Operator: opType,
						Value:    valueTokens.StringDecoded(),
						Token:    tokens,
					})

					operatorFound = true
					break
				}
			}
		}

		if operatorFound {
			break
		}
		currPos++
	}

	return results
}

func (p *TMFParser) parseSortValue(tokens rfcquery.TokenSlice) ([]TMFSortField, error) {
	var fields []TMFSortField

	splitTokens := p.splitTokens(tokens, ",")
	for _, fieldTokens := range splitTokens {
		if len(fieldTokens) == 0 {
			continue
		}

		firstTok := fieldTokens[0]
		direction := "asc"

		switch firstTok.Value {
		case "+":
			direction = "asc"
			fieldTokens = fieldTokens[1:]
		case "-":
			direction = "desc"
			fieldTokens = fieldTokens[1:]
		}

		if len(fieldTokens) == 0 {
			return nil, fmt.Errorf("empty sort field")
		}

		fieldName := fieldTokens.StringDecoded()
		fields = append(fields, TMFSortField{
			Field:     fieldName,
			Direction: direction,
			Tokens:    fieldTokens,
		})
	}
	return fields, nil
}

func (p *TMFParser) splitTokens(tokens rfcquery.TokenSlice, sep string) []rfcquery.TokenSlice {
	var result []rfcquery.TokenSlice
	currStart := 0

	for i, tok := range tokens {
		if tok.Type == rfcquery.TokenSubDelims && tok.Value == sep {
			if i > currStart {
				result = append(result, tokens[currStart:i])
			}
			currStart = i + 1
		}
	}

	if currStart < len(tokens) {
		result = append(result, tokens[currStart:])
	}

	return result
}

func ParseTMFQuery(query string) (*TMFQuery, error) {
	scanner := rfcquery.NewScanner(query)
	if err := scanner.Valid(); err != nil {
		return nil, err
	}

	parser := NewTMFParser()
	result, err := parser.Parse(scanner)
	if err != nil {
		return nil, err
	}

	tmfQuery, ok := result.(*TMFQuery)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return tmfQuery, nil
}
