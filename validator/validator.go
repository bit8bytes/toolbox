// Package validator provides input validation utilities with error collection
package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// Validator collects validation errors by field name
type Validator struct {
	Errors map[string]string
}

// New creates a Validator instance ready to collect validation errors
func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// Valid returns true when no validation errors exist
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError records a validation error for a field
// Existing errors for the same field are preserved (no overwrite)
func (v *Validator) AddError(key, message string) {
	if v.Errors == nil {
		v.Errors = make(map[string]string)
	}

	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check adds an error when validation fails
// Example: v.Check(len(name) > 0, "name", "Name is required")
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// PermittedValue returns true if value is in the permitted list
// Works with any comparable type (strings, ints, etc.)
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// Matches returns true if value matches the regex pattern
// Commonly used with EmailRX for email validation
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Unique returns true if all slice values are distinct
// Works with any comparable type
func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}

// NotBlank returns true if value contains non-whitespace characters
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars returns true if value has n characters or fewer
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// MinChars returns true if value has n characters or more
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}
