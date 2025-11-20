# rfcquery

A strict RFC3986-compliant query string parser for Go with a pluggable architecture.

## Problem

Go's standard library `url.ParseQuery()` implements `application/x-www-form-urlencoded` parsing, not the full RFC3986 query specification. This means:

- It rejects valid RFC3986 characters like `:` or `@` unless encoded
- It treats `+` as a space (HTML form specific)
- It assumes key-value pairs, but RFC3986 allows any structure

This library separates lexical validation (RFC3986 compliance) from semantic parsing.

## Quick Start

```go
import "github.com/yourusername/rfcquery"

// Just validate the query string is RFC3986 compliant
query := "filter=name:test&sort=created@asc"
lexer := rfcquery.NewLexer(query)
if err := lexer.Valid(); err != nil {
    // err includes position: "invalid character ' ' at position 7"
}

// Decode percent-encoded sequences
decoded, err := lexer.Decode()
// decoded: "filter=name:test&sort=created@asc"


