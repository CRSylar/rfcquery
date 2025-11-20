package rfcquery

// Value represent a parsed query value with metadata
type Value struct {
	// the decoded value
	Value string

	// Whether this key was seen multiple times
	HasMultiple bool

	// Positions information for precise error report
	KeyPos   Position
	ValuePos Position

	// Original Tokens for inspection
	KeyTokens   TokenSlice
	ValueTokens TokenSlice
}

// Values is a collection of parsed query params
// It preserves insertion order and allow duplicate keys
type Values struct {
	// map for fast lookup
	values map[string][]Value

	// slice for preserving order
	orderedKeys []string
}

func NewValues() *Values {
	return &Values{
		values:      make(map[string][]Value),
		orderedKeys: make([]string, 0),
	}
}

// Add a key-value pair
func (v *Values) Add(key string, value Value) {
	if _, exists := v.values[key]; !exists {
		v.orderedKeys = append(v.orderedKeys, key)
	}
	v.values[key] = append(v.values[key], value)
}

// First returns the first value for a key
func (v *Values) First(key string) (Value, bool) {
	if vals, ok := v.values[key]; ok && len(vals) > 0 {
		return vals[0], true
	}
	return Value{}, false
}

// Get return all values for a key
func (v *Values) Get(key string) []Value {
	return v.values[key]
}

// AllKeys returns all keys in insertion order
func (v *Values) AllKeys() []string {
	return v.orderedKeys
}

// Len return the total number of key-value pairs
func (v *Values) Len() int {
	count := 0
	for _, vals := range v.values {
		count += len(vals)
	}
	return count
}

// Parser is the interface that all query parsers must implement
type Parser interface {
	Parse(scanner *Scanner) (any, error)

	Name() string
}
