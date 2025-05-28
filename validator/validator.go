// Package validator implements utility for input validation
package validator

import (
	"regexp"
	"slices"
)

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// Validator holds a map of validation errors keyed by field name
// This allows accumulating multiple validation errors before returning them
type Validator struct {
	Errors map[string]string
}

// New creates and returns a new Validator instance with an empty error map
func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// Valid returns true if there are no validation errors stored
// This is typically called after running all validations to check overall validity
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds a validation error message for the given field key
// If an error already exists for this key, it will not be overwritten
// This prevents duplicate error messages for the same field
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check evaluates a boolean condition and adds an error if the condition is false
// This is a convenience method for conditional validation
// Example: v.Check(len(name) > 0, "name", "Name is required")
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// PermittedValue checks if a value exists within a list of permitted values
// Uses Go generics to work with any comparable type (strings, ints, etc.)
// Returns true if the value is found in the permitted values slice
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// Matches checks if a string value matches the provided regular expression
// Returns true if the string matches the pattern, false otherwise
// Commonly used with the EmailRX variable for email validation
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Unique checks if all values in a slice are unique (no duplicates)
// Uses Go generics to work with slices of any comparable type
// Returns true if all values are unique, false if duplicates are found
func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
