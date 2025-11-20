package plugins

import (
	"fmt"

	"github.com/CRSylar/rfcquery"
)

// FormURLEncodedParser implements application/x-www-form-urlencoded parsing
// This is RFC3986-compliant and preserves more metadata than URL stdlib
type FormURLEncodedParser struct {
	// PreserveInsertionOrder maintains the order of the keys as they appear
	PreserveInsertionOrder bool

	// AllowDuplicateKeys allows multiple values for the same key
	AllowDuplicateKeys bool
}

// Name returns the parser identifier
func (p *FormURLEncodedParser) Name() string {
	return "application/x-www-form-urlencoded"
}

// Parse implements the Parser interface
func (p *FormURLEncodedParser) Parse(scanner *rfcquery.Scanner) (any, error) {
	values := rfcquery.NewValues()

	var currKey rfcquery.TokenSlice
	var currValue rfcquery.TokenSlice

	for {
		tok, err := scanner.PeekToken()
		if err != nil {
			return nil, err
		}

		// Collect key
		currKey, err = scanner.CollectUntil(func(t rfcquery.Token) bool {
			return t.Type == rfcquery.TokenSubDelims && t.Value == "="
		})
		if err != nil {
			return nil, err
		}

		// Consume the '='
		eqTok, err := scanner.NextToken()
		if err != nil {
			return nil, err
		}

		if eqTok.Type == rfcquery.TokenEOF {
			// No more token, assume the value is empty
			if len(currKey) > 0 {
				keyStr := currKey.StringDecoded()
				val := rfcquery.Value{
					Value:       "",
					ValuePos:    eqTok.Start,
					KeyTokens:   currKey,
					ValueTokens: rfcquery.TokenSlice{},
				}
				values.Add(keyStr, val)
			}
			break
		}

		if eqTok.Value != "=" {
			return nil, fmt.Errorf("Expected '=', got %q", eqTok.Value)
		}

		currValue, err = scanner.CollectUntil(func(t rfcquery.Token) bool {
			return t.Type == rfcquery.TokenSubDelims && t.Value == "&"
		})
		if err != nil {
			return nil, err
		}

		keyStr := currKey.StringDecoded()
		valStr := currValue.StringDecoded()

		var valPos rfcquery.Position
		if len(currValue) > 0 {
			valPos = currValue[0].Start
		} else {
			valPos = rfcquery.Position{
				Offset: -1,
			}
		}

		value := rfcquery.Value{
			Value:       valStr,
			KeyPos:      currKey[0].Start,
			ValuePos:    valPos,
			KeyTokens:   currKey,
			ValueTokens: currValue,
		}

		values.Add(keyStr, value)

		tok, err = scanner.PeekToken()
		if err != nil {
			return nil, err
		}

		if tok.Type == rfcquery.TokenEOF {
			break
		}

		_, err = scanner.NextToken()
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}

// ParseFormURLEncoded - convenience function
func ParseFormURLEncoded(query string) (*rfcquery.Values, error) {
	scanner := rfcquery.NewScanner(query)
	if err := scanner.Valid(); err != nil {
		return nil, err
	}

	parser := &FormURLEncodedParser{
		PreserveInsertionOrder: true,
		AllowDuplicateKeys:     true,
	}

	result, err := parser.Parse(scanner)
	if err != nil {
		return nil, err
	}

	values, ok := result.(*rfcquery.Values)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return values, nil
}
