# rfcquery

[![Go Reference](https://pkg.go.dev/badge/github.com/CRSylar/rfcquery.svg)](https://pkg.go.dev/github.com/CRSylar/rfcquery)
[![Go Report Card](https://goreportcard.com/badge/github.com/CRSylar/rfcquery)](https://goreportcard.com/report/github.com/CRSylar/rfcquery)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A strict RFC3986-compliant query string parser for Go with a pluggable architecture. Parse URL queries with precision, extensibility, and performance.

## The Problem

Go's standard library `url.ParseQuery()` implements `application/x-www-form-urlencoded`, not the full RFC3986 query specification:

- ❌ Rejects valid RFC3986 characters (`:`, `/`, `?`, `@`) unless encoded
- ❌ Treats `+` as space (HTML form-specific, not RFC3986)
- ❌ Assumes key-value pairs only (RFC3986 allows any structure)
- ❌ No token-level access for custom parsing logic

## The Solution

`rfcquery` separates **lexical validation** from **semantic parsing**:

1. **Lexer Layer**: Strict RFC3986 validation with position-aware errors
2. **Token Stream**: Fine-grained access to query characters with lookahead
3. **Plugin System**: Pluggable parsers for different query formats

## Quick Start

```bash
go get github.com/CRSylar/rfcquery
```

## Examples

```go
package main

import (
    "fmt"
    "log"
    "github.com/CRSylar/rfcquery"
)

func main() {
    query := "filter[name]=John%20Doe&sort=created@asc"
    
    // Validate RFC3986 compliance
    scanner := rfcquery.NewScanner(query)
    if err := scanner.Valid(); err != nil {
        log.Fatal(err) // rfcquery: invalid character ' ' at position 7
    }
    
    // Parse as form-urlencoded (RFC3986-compliant)
    values, err := rfcquery.ParseFormURLEncoded(query)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Filter: %s\n", values.Get("filter[name]").Value)
    // Output: Filter: John Doe
}
```

## Features

### RFC3986 Strict Validation
```go
query := "user:pass@host/path?search"
scanner := rfcquery.NewScanner(url.QueryEscape(query))
err := scanner.Valid() // nil - valid RFC3986
```
validates percent-encoding, character classes, and provides precise error positions.

### Token Stream API
```go
scanner := rfcquery.NewScanner("name=John%20Doe")

// Token-by-token access
for {
    tok, err := scanner.NextToken()
    if err != nil {
        log.Fatal(err)
    }
    if tok.Type == rfcquery.TokenEOF {
        break
    }
    fmt.Printf("%s: %q (decoded: %q)\n", tok.Type, tok.Value, tok.Decoded)
}
```
### Bulk Collection for perfomance:
```go
// Collect until condition
tokens, err := scanner.CollectWhile(func(tok rfcquery.Token) bool {
    return tok.Type != rfcquery.TokenSubDelims || tok.Value != "&"
})

// Reconstruct strings
original := tokens.String()           // "name=John%20Doe"
decoded := tokens.StringDecoded()     // "name=John Doe"
```


### Plugin Architecture:

Built-in parsers with a common interface:

1. Form URL-Encoded ( `application/x-www-form-urlencoded`)
    ```go
    parser := &rfcquery.FormURLEncodedParser{
        PreserveInsertionOrder: true,
        AllowDuplicateKeys:     true,
    }

    scanner := rfcquery.NewScanner("tags=go&tags=library")
    values, err := parser.Parse(scanner)

    // Access all values for a key
    tags := values.Get("tags") // ["go", "library"]
    ```
    Advantages over `net/url`:
    - Preserves insertion order
    - Handles RFC3986 special characters ( :, @, /, ? )
    - Token-level metadata ( position, raw values)

2. JSON-in-query
    extract JSON from query parameter values:
    ```go
    query := `filter={"name":"John","age":30}&sort=created` // <-- NOTE: `{ / " / , / }` characters must be encoded to be valid in RFC3986, here is kept in plain text just for you to visually understand what is going on

    result, err := rfcquery.ParseJSONQuery(query, "filter")
    if err != nil {
        log.Fatal(err)
    }

    // Access the parsed JSON
    filterData := result.(map[string]interface{})
    fmt.Println(filterData["name"]) // "John"
    ```
    Features:
    - Parses percent-encoded JSON (`%7B%22name%22%3A%22John%22%7D`)
    - Handles special characters without mangling
    - Supports arrays, objects, primitives
    - Optional: can parse entire query string as JSON ( without the 'key' )

3. GraphQL-over-HTTP
    Extract GraphQL queries from URL parameters (per [GraphQL-over-HTTP spec](https://graphql.github.io/graphql-over-http/)):

    ```go
    query := `query=query GetUser($id: ID!) { user(id: $id) { name } }&variables={"id": "123"}`
    
    graphql, err := rfcquery.ParseGraphQLQuery(query)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(graphql.Query)         // The GraphQL query document
    fmt.Println(graphql.OperationName) // "GetUser"
    fmt.Println(graphql.Variables)     // map[id:123]
    ```
    Features:
    - Handles special GraphQL characters (@, !, $) without percent-encoding
    - Parses optional variables JSON parameter
    - Supports operationName for multiple operations
    - Works with percent-encoded queries
    
4. TMF API Guidelines (TMF630)
    Parse complex filter expressions following TMF630 guidelines:

    ```go
    query := "dateTime%3E%3D2013-04-20;status=active,suspended&sort=-created,+name&limit=10"

    tmf, err := rfcquery.ParseTMFQuery(query)
    if err != nil {
        log.Fatal(err)
    }

    // Access filter expressions
    for _, expr := range tmf.Expressions {
        fmt.Printf("%s %s %s\n", expr.Field, expr.Operator, expr.Value)
        // Output: dateTime gte 2013-04-20
        //         status eq active,suspended
    }

    // Access sorting
    for _, sort := range tmf.Sorting {
        fmt.Printf("Sort by %s (%s)\n", sort.Field, sort.Direction)
        // Output: Sort by created (desc)
        //         Sort by name (asc)
    }

    // Access pagination params
    limit := tmf.OtherParams["limit"][0] // "10"
    ```
    - Encoded Operators ( %3E for `>`, %3C for `<`, etc..)
    - Multiple separators ( `=` and `;` treated identically)
    - List values (comma-separated)
    - Sorting with +/- prefixes
    - MultipleOperators on same field (date%3E2017-04-01;date%3C2017-05-01)
    - Values containing encoded opeartors ( no false positives)
    - RFC3986-compliant ( special chars like `@`, `:`, `/` work correctly)

5. Custom Parser
    To implement a custom parser implement the `Parser` interface
    ```go
    type MyCustomParser struct{}

    func (p *MyCustomParser) Parse(scanner *rfcquery.Scanner) (interface{}, error) {
        // Use scanner.CollectWhile(), scanner.NextToken(), etc.
        // Return your custom data structure
    }

    func (p *MyCustomParser) Name() string {
        return "my-custom-parser"
    }
    ```



### Roadmap
 - [X] GraphQL query parser plugin
 - [X] TMF query parser plugin
 - [ ] Query builder API (fluent interface)
 - [ ] Streaming parser for very large queries
 - [ ] JSON Schema validation for JSON-in-query
 - [ ] Performance optimizations with pooled scanner
 - [ ] encoder package for strict rfc encoding

## Contributing

We welcome contributions! Please see [CONTRIBUTING](https://github.com/CRSylar/rfcquery/CONTRIBUTING.md) for guidelines.


## License
MIT License - see [LICENSE](https://github.com/CRSylar/rfcquery/LICENSE.md) file for details.

## Why choose rfcquery?
✅ Correctness: Strict RFC3986 compliance
✅ Flexibility: Plugin system for any query format
✅ Performance: Bulk operations and minimal allocations
✅ Developer Experience: Clear errors with positions
