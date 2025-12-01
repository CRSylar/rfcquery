package tmfparser_test

import (
	"log/slog"
	"testing"

	"github.com/CRSylar/rfcquery"
	tmfparser "github.com/CRSylar/rfcquery/plugins/tmf_parser"
)

func TestTMFParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		check   func(t *testing.T, q *tmfparser.TMFQuery)
		wantErr bool
	}{
		{
			name:  "simple equality with =",
			input: "name=John",
			check: func(t *testing.T, q *tmfparser.TMFQuery) {
				if len(q.Expressions) != 1 {
					t.Fatalf("expected 1 expression, got %d - expr %+v", len(q.Expressions), q.Expressions)
				}
				expr := q.Expressions
				if len(expr["name"]) != 1 || expr["name"][0].Operator != tmfparser.TMFOperatorEq || expr["name"][0].Value != "John" {
					t.Errorf("unexpected expression: %+v", expr)
				}
			},
		},
		{
			name:  "greater than operator",
			input: "dateTime%3E2013-04-20",
			check: func(t *testing.T, q *tmfparser.TMFQuery) {
				if len(q.Expressions) != 1 {
					t.Fatalf("expected 1 expression, got %d", len(q.Expressions))
				}
				exprs := q.Expressions
				if len(exprs["dateTime"]) != 1 {
					t.Fatalf("expected 1 value for the 'dateTime' expression, got %+v", exprs["dateTime"])
				}

				if exprs["dateTime"][0].Operator != tmfparser.TMFOperatorGt {
					t.Fatalf("expected > operator, got %s", exprs["dateTime"][0].Operator)
				}

				if exprs["dateTime"][0].Value != "2013-04-20" {
					t.Fatalf("unexpected value, want 2013-04-20, got %v", exprs["dateTime"][0].Value)
				}
			},
		},
		{
			name:  "less than operator",
			input: "age%3C65",
			check: func(t *testing.T, q *tmfparser.TMFQuery) {
				if len(q.Expressions) != 1 {
					t.Fatalf("expected 1 expression, got %d", len(q.Expressions))
				}
				expr := q.Expressions
				if len(expr["age"]) != 1 || expr["age"][0].Operator != tmfparser.TMFOperatorLt || expr["age"][0].Value != "65" {
					t.Errorf("unexpected expression: %+v", expr)
				}
			},
		},
		{
			name:  "greater than or equal operator",
			input: "dateTime%3E%3D2013-04-20",
			check: func(t *testing.T, q *tmfparser.TMFQuery) {
				if len(q.Expressions) != 1 {
					t.Fatalf("expected 1 expression, got %d", len(q.Expressions))
				}
				expr := q.Expressions
				if len(expr["dateTime"]) != 1 || expr["dateTime"][0].Operator != tmfparser.TMFOperatorGte || expr["dateTime"][0].Value != "2013-04-20" {
					t.Errorf("unexpected expression: %+v", expr)
				}
			},
		},
		{
			name:  "less than or equal operator",
			input: "age%3C%3D65",
			check: func(t *testing.T, q *tmfparser.TMFQuery) {
				if len(q.Expressions) != 1 {
					t.Fatalf("expected 1 expression, got %d", len(q.Expressions))
				}
				expr := q.Expressions
				if len(expr["age"]) != 1 || expr["age"][0].Operator != tmfparser.TMFOperatorLte || expr["age"][0].Value != "65" {
					t.Errorf("unexpected expression: %+v", expr)
				}
			},
		},
		{
			name:  "not equal operator",
			input: "status%21%3Ddeleted",
			check: func(t *testing.T, q *tmfparser.TMFQuery) {
				if len(q.Expressions) != 1 {
					t.Fatalf("expected 1 expression, got %d", len(q.Expressions))
				}
				expr := q.Expressions
				if len(expr["status"]) != 1 || expr["status"][0].Operator != tmfparser.TMFOperatorNe || expr["status"][0].Value != "deleted" {
					t.Errorf("unexpected expression: %+v", expr)
				}
			},
		},
		{
			name:  "multiple expressions with ;",
			input: "name=John;age%3E25;status=active",
			check: func(t *testing.T, q *tmfparser.TMFQuery) {
				if len(q.Expressions) != 3 {
					t.Fatalf("expected 3 expression, got %d", len(q.Expressions))
				}

				if q.Expressions["name"][0].Value != "John" {
					t.Errorf("first expression if wrong: %+v", q.Expressions["name"])
				}

				if q.Expressions["age"][0].Operator != tmfparser.TMFOperatorGt || q.Expressions["age"][0].Value != "25" {
					t.Errorf("second expression if wrong: %+v", q.Expressions["age"])
				}
				if q.Expressions["status"][0].Value != "active" {
					t.Errorf("third expression if wrong: %+v", q.Expressions["status"])
				}
			},
		},
		{
			name:  "multiple expressions with &",
			input: "name=John&age%3E25&status=active",
			check: func(t *testing.T, q *tmfparser.TMFQuery) {
				if len(q.Expressions) != 3 {
					t.Fatalf("expected 3 expression, got %d", len(q.Expressions))
				}

				if q.Expressions["name"][0].Value != "John" {
					t.Errorf("first expression if wrong: %+v", q.Expressions["name"])
				}

				if q.Expressions["age"][0].Operator != tmfparser.TMFOperatorGt || q.Expressions["age"][0].Value != "25" {
					t.Errorf("second expression if wrong: %+v", q.Expressions["age"])
				}
				if q.Expressions["status"][0].Value != "active" {
					t.Errorf("third expression if wrong: %+v", q.Expressions["status"])
				}
			},
		},
		{
			name:  "mixed separators",
			input: "name=John;age%3E25&status=active",
			check: func(t *testing.T, q *tmfparser.TMFQuery) {
				if len(q.Expressions) != 3 {
					t.Fatalf("expected 3 expression, got %d", len(q.Expressions))
				}

				if q.Expressions["name"][0].Value != "John" {
					t.Errorf("first expression if wrong: %+v", q.Expressions["name"])
				}

				if q.Expressions["age"][0].Operator != tmfparser.TMFOperatorGt || q.Expressions["age"][0].Value != "25" {
					t.Errorf("second expression if wrong: %+v", q.Expressions["age"])
				}
				if q.Expressions["status"][0].Value != "active" {
					t.Errorf("third expression if wrong: %+v", q.Expressions["status"])
				}
			},
		},
		{
			name:  "value encoded operator in value (no operator meaning)",
			input: "description=value%3Etest",
			check: func(t *testing.T, q *tmfparser.TMFQuery) {
				if len(q.Expressions) != 1 {
					t.Fatalf("expected 1 expression, got %d", len(q.Expressions))
				}

				expr := q.Expressions["description"]
				if expr[0].Operator != tmfparser.TMFOperatorEq || expr[0].Value != "value>test" {
					t.Errorf("expected equality with value 'value>test', got: %+v", expr)
				}
			},
		},
		{
			name:  "multiple operators on same field",
			input: "dateTime%3E2013-04-20;dateTime%3C2017-04-20",
			check: func(t *testing.T, q *tmfparser.TMFQuery) {
				if len(q.Expressions) != 1 {
					t.Fatalf("expected 1 expressions got: %d", len(q.Expressions))
				}

				if len(q.Expressions["dateTime"]) != 2 {
					t.Fatalf("expected 2 expressions for the 'dateTime' key got: %d", len(q.Expressions))
				}

				slog.Info("got", "expr", q.Expressions["dateTime"])
				if q.Expressions["dateTime"][0].Operator != tmfparser.TMFOperatorGt || q.Expressions["dateTime"][0].Value != "2013-04-20" {
					t.Errorf("first dateTime expression error: %+v", q.Expressions["dateTime"][0])
				}

				if q.Expressions["dateTime"][1].Operator != tmfparser.TMFOperatorLt || q.Expressions["dateTime"][1].Value != "2017-04-20" {
					t.Errorf("second dateTime expression error: %+v", q.Expressions["dateTime"][1])
				}
			},
		},
		{
			name:  "comma-separated values (list)",
			input: "status=active,suspended,pending",
			check: func(t *testing.T, q *tmfparser.TMFQuery) {
				if len(q.Expressions) != 1 {
					t.Fatalf("expected 1 expression got: %d", len(q.Expressions))
				}
				expr := q.Expressions["status"]
				if len(expr) != 3 {
					t.Errorf("expected 3 values in list, got %d", len(expr))
				}
				if expr[0].Value != "active" || expr[1].Value != "suspended" || expr[2].Value != "pending" {
					got := []string{expr[0].Value, expr[1].Value, expr[2].Value}
					t.Errorf("Expressions array is not right, expected ['active', 'suspended', 'pending'] - got: %+v", got)
				}
			},
		},
		{
			name:  "repeated keys treated as comma-separated values",
			input: "status=active&status=suspended&status=pending",
			check: func(t *testing.T, q *tmfparser.TMFQuery) {
				if len(q.Expressions) != 1 {
					t.Fatalf("expected 1 expression got: %d", len(q.Expressions))
				}

				expr := q.Expressions["status"]
				if len(expr) != 3 {
					t.Fatalf("expected 3 values for expression 'status' got: %d", len(q.Expressions["status"]))
				}

				if expr[0].Value != "active" || expr[1].Value != "suspended" || expr[2].Value != "pending" {
					t.Errorf("Expressions array is not right, expected ['active', 'suspended', 'pending'] - got: %+v", []string{expr[0].Value, expr[1].Value, expr[2].Value})
				}
				if expr[0].Operator != expr[1].Operator || expr[1].Operator != expr[2].Operator {
					t.Errorf("Operators mismatch, expected operators to be equal to each other: %+v", expr)
				}
				if expr[0].Operator != tmfparser.TMFOperatorEq {
					t.Errorf("expected operator to be '=', got %+v", expr[0].Operator)
				}
			},
		},
		{
			name:  "sort parameter ascending",
			input: "sort=+name",
			check: func(t *testing.T, q *tmfparser.TMFQuery) {
				if len(q.Sorting) != 1 {
					t.Fatalf("expected 1 sort field, got %d", len(q.Sorting))
				}
				sort := q.Sorting[0]
				if sort.Field != "name" || sort.Direction != "asc" {
					t.Errorf("unexpected sort: %+v", sort)
				}
			},
		},
		{
			name:  "sort parameter descending",
			input: "sort=-created",
			check: func(t *testing.T, q *tmfparser.TMFQuery) {
				if len(q.Sorting) != 1 {
					t.Fatalf("expected 1 sort field, got %d", len(q.Sorting))
				}
				sort := q.Sorting[0]
				if sort.Field != "created" || sort.Direction != "desc" {
					t.Errorf("unexpected sort: %+v", sort)
				}
			},
		},
		{
			name:  "sort multiple fields mixed",
			input: "sort=+name,-created,age",
			check: func(t *testing.T, q *tmfparser.TMFQuery) {
				if len(q.Sorting) != 3 {
					t.Fatalf("expected 3 sort field, got %d", len(q.Sorting))
				}
				if q.Sorting[0].Field != "name" || q.Sorting[0].Direction != "asc" {
					t.Errorf("first sort field error: %+v", q.Sorting[0])
				}
				if q.Sorting[1].Field != "created" || q.Sorting[1].Direction != "desc" {
					t.Errorf("second sort field error: %+v", q.Sorting[1])
				}
				if q.Sorting[2].Field != "age" || q.Sorting[2].Direction != "asc" {
					t.Errorf("third sort field error: %+v", q.Sorting[2])
				}
			},
		},
		{
			name:  "other paramters (limit, offset)",
			input: "limit=10&offset=20",
			check: func(t *testing.T, q *tmfparser.TMFQuery) {
				if q.OtherParams == nil {
					t.Fatalf("OtherParams should not be nil")
				}
				if limitVals, ok := q.OtherParams["limit"]; !ok || len(limitVals) == 0 || limitVals[0] != "10" {
					t.Errorf("limit parameter invalid, want 10, have: %v", limitVals)
				}
				if offsetVals, ok := q.OtherParams["offset"]; !ok || len(offsetVals) == 0 || offsetVals[0] != "20" {
					t.Errorf("offset parameter invalid, want 10, have: %v", offsetVals)
				}
			},
		},
		{
			name:  "complex example from TMF spec, with mixed params",
			input: "name=John;age%3E25;status=active,suspended&sort=-created,+name&limit=10",
			check: func(t *testing.T, q *tmfparser.TMFQuery) {
				if len(q.Expressions) != 3 {
					t.Fatalf("expected 3 expressions, got %d", len(q.Expressions))
				}
				if len(q.Sorting) != 2 {
					t.Fatalf("expected 2 sort fields, got %d", len(q.Sorting))
				}
				if limitVals, ok := q.OtherParams["limit"]; !ok || limitVals[0] != "10" {
					t.Errorf("limit parameter missing or wrong: %+v", limitVals)
				}
			},
		},
		{
			name:    "invalid operator encoding",
			input:   "dateTime>2013-04-20",
			wantErr: true,
		},
		{
			name:  "empty value",
			input: "name=",
			check: func(t *testing.T, q *tmfparser.TMFQuery) {
				if len(q.Expressions) != 1 || len(q.Expressions["name"]) != 0 {
					t.Errorf("empty value not handled correctly - got: %+v", q.Expressions)
				}
			},
		},
		{
			name:  "key only",
			input: "key",
			check: func(t *testing.T, q *tmfparser.TMFQuery) {
				if len(q.Expressions) != 1 {
					t.Errorf("key-only expression not handled correctly, got: %+v", q.Expressions)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := tmfparser.NewTMFParser()
			scanner := rfcquery.NewScanner(tt.input)

			result, err := parser.Parse(scanner)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			tmfQuery, ok := result.(*tmfparser.TMFQuery)
			if !ok {
				t.Fatalf("expected *TMFQuery, got %T", result)
			}

			tt.check(t, tmfQuery)
		})
	}
}
