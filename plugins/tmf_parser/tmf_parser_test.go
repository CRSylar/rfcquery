package tmfparser_test

import (
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
		// {
		// 	name:  "simple equality with =",
		// 	input: "name=John",
		// 	check: func(t *testing.T, q *tmfparser.TMFQuery) {
		// 		if len(q.Expressions) != 1 {
		// 			t.Fatalf("expected 1 expression, got %d", len(q.Expressions))
		// 		}
		// 		expr := q.Expressions[0]
		// 		if expr.Field != "name" || expr.Operator != tmfparser.TMFOperatorEq || expr.Value != "John" {
		// 			t.Errorf("unexpected expression: %+v", expr)
		// 		}
		// 	},
		// },
		{
			name:  "greater than operator",
			input: "dateTime%3E2013-04-20",
			check: func(t *testing.T, q *tmfparser.TMFQuery) {
				if len(q.Expressions) != 1 {
					t.Fatalf("expected 1 expression, got %d", len(q.Expressions))
				}
				expr := q.Expressions[0]
				if expr.Field != "dateTime" || expr.Operator != tmfparser.TMFOperatorGt || expr.Value != "2013-04-20" {
					t.Errorf("unexpected expression: %+v", expr)
				}
			},
		},
		// {
		// 	name:  "less than operator",
		// 	input: "age%3C65",
		// 	check: func(t *testing.T, q *tmfparser.TMFQuery) {
		// 		if len(q.Expressions) != 1 {
		// 			t.Fatalf("expected 1 expression, got %d", len(q.Expressions))
		// 		}
		// 		expr := q.Expressions[0]
		// 		if expr.Field != "age" || expr.Operator != tmfparser.TMFOperatorLt || expr.Value != "65" {
		// 			t.Errorf("unexpected expression: %+v", expr)
		// 		}
		// 	},
		// },
		// {
		// 	name:  "greater than or equal operator",
		// 	input: "dateTime%3E%3D2013-04-20",
		// 	check: func(t *testing.T, q *tmfparser.TMFQuery) {
		// 		if len(q.Expressions) != 1 {
		// 			t.Fatalf("expected 1 expression, got %d", len(q.Expressions))
		// 		}
		// 		expr := q.Expressions[0]
		// 		if expr.Field != "dateTime" || expr.Operator != tmfparser.TMFOperatorGte || expr.Value != "2013-04-20" {
		// 			t.Errorf("unexpected expression: %+v", expr)
		// 		}
		// 	},
		// },
		// {
		// 	name:  "less than or equal operator",
		// 	input: "age%3C%3D65",
		// 	check: func(t *testing.T, q *tmfparser.TMFQuery) {
		// 		if len(q.Expressions) != 1 {
		// 			t.Fatalf("expected 1 expression, got %d", len(q.Expressions))
		// 		}
		// 		expr := q.Expressions[0]
		// 		if expr.Field != "age" || expr.Operator != tmfparser.TMFOperatorLte || expr.Value != "65" {
		// 			t.Errorf("unexpected expression: %+v", expr)
		// 		}
		// 	},
		// },
		// {
		// 	name:  "not equal operator",
		// 	input: "status%21%3Ddeleted",
		// 	check: func(t *testing.T, q *tmfparser.TMFQuery) {
		// 		if len(q.Expressions) != 1 {
		// 			t.Fatalf("expected 1 expression, got %d", len(q.Expressions))
		// 		}
		// 		expr := q.Expressions[0]
		// 		if expr.Field != "status" || expr.Operator != tmfparser.TMFOperatorNe || expr.Value != "deleted" {
		// 			t.Errorf("unexpected expression: %+v", expr)
		// 		}
		// 	},
		// },
		// {
		// 	name:  "multiple expressions with ;",
		// 	input: "name=John;age%3E25;status=active",
		// 	check: func(t *testing.T, q *tmfparser.TMFQuery) {
		// 		if len(q.Expressions) != 3 {
		// 			t.Fatalf("expected 3 expression, got %d", len(q.Expressions))
		// 		}
		//
		// 		if q.Expressions[0].Field != "name" || q.Expressions[0].Value != "John" {
		// 			t.Errorf("first expression if wrong: %+v", q.Expressions[0])
		// 		}
		//
		// 		if q.Expressions[1].Field != "age" || q.Expressions[1].Operator != tmfparser.TMFOperatorGt || q.Expressions[1].Value != "25" {
		// 			t.Errorf("second expression if wrong: %+v", q.Expressions[0])
		// 		}
		// 		if q.Expressions[2].Field != "status" || q.Expressions[2].Value != "active" {
		// 			t.Errorf("third expression if wrong: %+v", q.Expressions[0])
		// 		}
		// 	},
		// },
		// {
		// 	name:  "multiple expressions with &",
		// 	input: "name=John&age%3E25&status=active",
		// 	check: func(t *testing.T, q *tmfparser.TMFQuery) {
		// 		if len(q.Expressions) != 3 {
		// 			t.Fatalf("expected 3 expression, got %d", len(q.Expressions))
		// 		}
		//
		// 		if q.Expressions[0].Field != "name" || q.Expressions[0].Value != "John" {
		// 			t.Errorf("first expression if wrong: %+v", q.Expressions[0])
		// 		}
		//
		// 		if q.Expressions[1].Field != "age" || q.Expressions[1].Operator != tmfparser.TMFOperatorGt || q.Expressions[1].Value != "25" {
		// 			t.Errorf("second expression if wrong: %+v", q.Expressions[0])
		// 		}
		// 		if q.Expressions[2].Field != "status" || q.Expressions[2].Value != "active" {
		// 			t.Errorf("third expression if wrong: %+v", q.Expressions[0])
		// 		}
		// 	},
		// },
		// {
		// 	name:  "mixed separators",
		// 	input: "name=John;age%3E25&status=active",
		// 	check: func(t *testing.T, q *tmfparser.TMFQuery) {
		// 		if len(q.Expressions) != 3 {
		// 			t.Fatalf("expected 3 expression, got %d", len(q.Expressions))
		// 		}
		//
		// 		if q.Expressions[0].Field != "name" || q.Expressions[0].Value != "John" {
		// 			t.Errorf("first expression if wrong: %+v", q.Expressions[0])
		// 		}
		//
		// 		if q.Expressions[1].Field != "age" || q.Expressions[1].Operator != tmfparser.TMFOperatorGt || q.Expressions[1].Value != "25" {
		// 			t.Errorf("second expression if wrong: %+v", q.Expressions[0])
		// 		}
		// 		if q.Expressions[2].Field != "status" || q.Expressions[2].Value != "active" {
		// 			t.Errorf("third expression if wrong: %+v", q.Expressions[0])
		// 		}
		// 	},
		// },
		// {
		// 	name:  "value encoded operator in value (no operator meaning)",
		// 	input: "description=value%3Etest",
		// 	check: func(t *testing.T, q *tmfparser.TMFQuery) {
		// 		if len(q.Expressions) != 1 {
		// 			t.Fatalf("expected 1 expression, got %d", len(q.Expressions))
		// 		}
		//
		// 		expr := q.Expressions[0]
		// 		if expr.Operator != tmfparser.TMFOperatorEq || expr.Value != "value>test" {
		// 			t.Errorf("expected equality with value 'value>test', got: %+v", expr)
		// 		}
		// 	},
		// },
		// {
		// 	name:  "multiple operators on same field",
		// 	input: "dateTime%E2013-04-20;dateTime%3C2017-04-20",
		// 	check: func(t *testing.T, q *tmfparser.TMFQuery) {
		// 		if len(q.Expressions) != 2 {
		// 			t.Fatalf("expected 2 expressions got: %d", len(q.Expressions))
		// 		}
		//
		// 		if q.Expressions[0].Field != "dateTime" || q.Expressions[0].Operator != tmfparser.TMFOperatorGt || q.Expressions[0].Value != "2013-04-20" {
		// 			t.Errorf("first dateTime expression error: %+v", q.Expressions[0])
		// 		}
		//
		// 		if q.Expressions[1].Field != "dateTime" || q.Expressions[1].Operator != tmfparser.TMFOperatorLt || q.Expressions[0].Value != "2017-04-20" {
		// 			t.Errorf("second dateTime expression error: %+v", q.Expressions[1])
		// 		}
		// 	},
		// },
		// {
		// 	name:  "comma-separated values (list)",
		// 	input: "status=active,suspended,pending",
		// 	check: func(t *testing.T, q *tmfparser.TMFQuery) {
		// 		if len(q.Expressions) != 1 {
		// 			t.Fatalf("expected 1 expression got: %d", len(q.Expressions))
		// 		}
		// 		expr := q.Expressions[0]
		// 		if expr.Value != "active,suspended,pending" {
		// 			t.Errorf("expected value list, got %s", expr.Value)
		// 		}
		// 	},
		// },
		// {
		// 	name:  "repeated keys treated as comma-separated values",
		// 	input: "status=active&status=suspended&status=pending",
		// 	check: func(t *testing.T, q *tmfparser.TMFQuery) {
		// 		if len(q.Expressions) != 1 {
		// 			t.Fatalf("expected 1 expression got: %d", len(q.Expressions))
		// 		}
		//
		// 		expr := q.Expressions[0]
		// 		if expr.Value != "active,suspended,pending" {
		// 			t.Errorf("expected values list, got %s", expr.Value)
		// 		}
		// 	},
		// },
		// {
		// 	name:  "sort parameter ascending",
		// 	input: "sort=+name",
		// 	check: func(t *testing.T, q *tmfparser.TMFQuery) {
		// 		if len(q.Sorting) != 1 {
		// 			t.Fatalf("expected 1 sort field, got %d", len(q.Sorting))
		// 		}
		// 		sort := q.Sorting[0]
		// 		if sort.Field != "name" || sort.Direction != "asc" {
		// 			t.Errorf("unexpected sort: %+v", sort)
		// 		}
		// 	},
		// },
		// {
		// 	name:  "sort parameter descending",
		// 	input: "sort=-created",
		// 	check: func(t *testing.T, q *tmfparser.TMFQuery) {
		// 		if len(q.Sorting) != 1 {
		// 			t.Fatalf("expected 1 sort field, got %d", len(q.Sorting))
		// 		}
		// 		sort := q.Sorting[0]
		// 		if sort.Field != "created" || sort.Direction != "desc" {
		// 			t.Errorf("unexpected sort: %+v", sort)
		// 		}
		// 	},
		// },
		// {
		// 	name:  "sort multiple fields mixed",
		// 	input: "sort=+name,-created,age",
		// 	check: func(t *testing.T, q *tmfparser.TMFQuery) {
		// 		if len(q.Sorting) != 3 {
		// 			t.Fatalf("expected 3 sort field, got %d", len(q.Sorting))
		// 		}
		// 		if q.Sorting[0].Field != "name" || q.Sorting[0].Direction != "asc" {
		// 			t.Errorf("first sort field error: %+v", q.Sorting[0])
		// 		}
		// 		if q.Sorting[1].Field != "created" || q.Sorting[1].Direction != "desc" {
		// 			t.Errorf("second sort field error: %+v", q.Sorting[1])
		// 		}
		// 		if q.Sorting[2].Field != "age" || q.Sorting[2].Direction != "asc" {
		// 			t.Errorf("third sort field error: %+v", q.Sorting[2])
		// 		}
		// 	},
		// },
		// {
		// 	name:  "other paramters (limit, offset)",
		// 	input: "limit=10&offset=20",
		// 	check: func(t *testing.T, q *tmfparser.TMFQuery) {
		// 		if q.OtherParams == nil {
		// 			t.Fatalf("OtherParams should not be nil")
		// 		}
		// 		if limitVals, ok := q.OtherParams["limit"]; !ok || len(limitVals) == 0 || limitVals[0] != "10" {
		// 			t.Errorf("limit parameter invalid, want 10, have: %v", limitVals)
		// 		}
		// 		if offsetVals, ok := q.OtherParams["offset"]; !ok || len(offsetVals) == 0 || offsetVals[0] != "20" {
		// 			t.Errorf("offset parameter invalid, want 10, have: %v", offsetVals)
		// 		}
		// 	},
		// },
		// {
		// 	name:  "complex example from TMF spec, with mixed params",
		// 	input: "name=John;age%3E25;status=active,suspended&sort=-created,+name&limit=10",
		// 	check: func(t *testing.T, q *tmfparser.TMFQuery) {
		// 		if len(q.Expressions) != 3 {
		// 			t.Fatalf("expected 3 expressions, got %d", len(q.Expressions))
		// 		}
		// 		if len(q.Sorting) != 2 {
		// 			t.Fatalf("expected 2 sort fields, got %d", len(q.Sorting))
		// 		}
		// 		if limitVals, ok := q.OtherParams["limit"]; !ok || limitVals[0] != "10" {
		// 			t.Errorf("limit parameter missing or wrong: %+v", limitVals)
		// 		}
		// 	},
		// },
		// {
		// 	name:    "invalid operator encoding",
		// 	input:   "dateTime>2013-04-20",
		// 	wantErr: true,
		// },
		// {
		// 	name:  "empty value",
		// 	input: "name=",
		// 	check: func(t *testing.T, q *tmfparser.TMFQuery) {
		// 		if len(q.Expressions) != 1 || q.Expressions[0].Value != "" {
		// 			t.Errorf("empty value not handled correctly")
		// 		}
		// 	},
		// },
		// {
		// 	name:  "key only",
		// 	input: "key",
		// 	check: func(t *testing.T, q *tmfparser.TMFQuery) {
		// 		if len(q.Expressions) != 1 || q.Expressions[0].Field != "key" {
		// 			t.Errorf("key-only expression not handled correctly")
		// 		}
		// 	},
		// },
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
