package validator

import (
	"regexp"
	"testing"
)

func TestNew(t *testing.T) {
	v := New()
	if v == nil {
		t.Error("New() returned nil")
	}
	if v.Errors == nil {
		t.Error("Errors map is nil")
	}
	if len(v.Errors) != 0 {
		t.Error("New validator should have no errors")
	}
}

func TestValid(t *testing.T) {
	v := New()
	if !v.Valid() {
		t.Error("New validator should be valid")
	}

	v.AddError("test", "error")
	if v.Valid() {
		t.Error("Validator with errors should not be valid")
	}
}

func TestAddError(t *testing.T) {
	v := New()

	v.AddError("field1", "error1")
	if len(v.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(v.Errors))
	}
	if v.Errors["field1"] != "error1" {
		t.Errorf("Expected 'error1', got '%s'", v.Errors["field1"])
	}

	// Test no overwrite of existing error
	v.AddError("field1", "error2")
	if v.Errors["field1"] != "error1" {
		t.Error("Existing error should not be overwritten")
	}

	// Test nil map handling
	v2 := &Validator{}
	v2.AddError("test", "message")
	if v2.Errors == nil {
		t.Error("Errors map should be initialized")
	}
	if v2.Errors["test"] != "message" {
		t.Error("Error should be added to nil map")
	}
}

func TestCheck(t *testing.T) {
	v := New()

	v.Check(true, "field1", "error1")
	if len(v.Errors) != 0 {
		t.Error("No error should be added when condition is true")
	}

	v.Check(false, "field1", "error1")
	if len(v.Errors) != 1 {
		t.Error("Error should be added when condition is false")
	}
	if v.Errors["field1"] != "error1" {
		t.Error("Correct error message should be added")
	}
}

func TestPermittedValue(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func() bool
		expected bool
	}{
		{
			name:     "String value in permitted list",
			testFunc: func() bool { return PermittedValue("apple", "apple", "banana", "cherry") },
			expected: true,
		},
		{
			name:     "String value not in permitted list",
			testFunc: func() bool { return PermittedValue("grape", "apple", "banana", "cherry") },
			expected: false,
		},
		{
			name:     "Integer value in permitted list",
			testFunc: func() bool { return PermittedValue(2, 1, 2, 3) },
			expected: true,
		},
		{
			name:     "Integer value not in permitted list",
			testFunc: func() bool { return PermittedValue(4, 1, 2, 3) },
			expected: false,
		},
		{
			name:     "Empty permitted list",
			testFunc: func() bool { return PermittedValue("any", []string{}...) },
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.testFunc()
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestMatches(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		pattern  *regexp.Regexp
		expected bool
	}{
		{
			name:     "Valid email matches pattern",
			value:    "test@example.com",
			pattern:  EmailRX,
			expected: true,
		},
		{
			name:     "Invalid email does not match pattern",
			value:    "invalid-email",
			pattern:  EmailRX,
			expected: false,
		},
		{
			name:     "Empty string does not match pattern",
			value:    "",
			pattern:  EmailRX,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Matches(test.value, test.pattern)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestEmailRX(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{
			name:     "Valid email with subdomain",
			email:    "user+tag@example.co.uk",
			expected: true,
		},
		{
			name:     "Valid email with numbers",
			email:    "123@example.org",
			expected: true,
		},
		{
			name:     "Valid standard email",
			email:    "test@example.com",
			expected: true,
		},
		{
			name:     "Valid email with dots",
			email:    "user.name@example.com",
			expected: true,
		},
		{
			name:     "Invalid email format",
			email:    "invalid-email",
			expected: false,
		},
		{
			name:     "Missing local part",
			email:    "@example.com",
			expected: false,
		},
		{
			name:     "Missing domain",
			email:    "test@",
			expected: false,
		},
		{
			name:     "Invalid domain format",
			email:    "test@.com",
			expected: false,
		},
		{
			name:     "Missing TLD",
			email:    "test@example",
			expected: false,
		},
		{
			name:     "Missing @ symbol",
			email:    "test.example.com",
			expected: false,
		},
		{
			name:     "Single character TLD",
			email:    "test@example.c",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Matches(test.email, EmailRX)
			if result != test.expected {
				t.Errorf("Email %s: expected %v, got %v", test.email, test.expected, result)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func() bool
		expected bool
	}{
		{
			name:     "Unique strings return true",
			testFunc: func() bool { return Unique([]string{"a", "b", "c"}) },
			expected: true,
		},
		{
			name:     "Duplicate strings return false",
			testFunc: func() bool { return Unique([]string{"a", "b", "a"}) },
			expected: false,
		},
		{
			name:     "Unique integers return true",
			testFunc: func() bool { return Unique([]int{1, 2, 3}) },
			expected: true,
		},
		{
			name:     "Duplicate integers return false",
			testFunc: func() bool { return Unique([]int{1, 2, 1}) },
			expected: false,
		},
		{
			name:     "Empty slice is unique",
			testFunc: func() bool { return Unique([]string{}) },
			expected: true,
		},
		{
			name:     "Single element is unique",
			testFunc: func() bool { return Unique([]string{"a"}) },
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.testFunc()
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestNotBlank(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{
			name:     "Non-empty string is not blank",
			value:    "hello",
			expected: true,
		},
		{
			name:     "Empty string is blank",
			value:    "",
			expected: false,
		},
		{
			name:     "Whitespace-only string is blank",
			value:    "   ",
			expected: false,
		},
		{
			name:     "Tab and newline characters are blank",
			value:    "\t\n ",
			expected: false,
		},
		{
			name:     "String with content and whitespace is not blank",
			value:    "  hello  ",
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := NotBlank(test.value)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestMaxChars(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		maxChars int
		expected bool
	}{
		{
			name:     "String with exact max length passes",
			value:    "hello",
			maxChars: 5,
			expected: true,
		},
		{
			name:     "String under max length passes",
			value:    "hi",
			maxChars: 5,
			expected: true,
		},
		{
			name:     "String over max length fails",
			value:    "hello world",
			maxChars: 5,
			expected: false,
		},
		{
			name:     "Empty string passes any max length",
			value:    "",
			maxChars: 5,
			expected: true,
		},
		{
			name:     "Unicode string with exact max length passes",
			value:    "héllo",
			maxChars: 5,
			expected: true,
		},
		{
			name:     "Unicode string over max length fails",
			value:    "héllo world",
			maxChars: 5,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := MaxChars(test.value, test.maxChars)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestMinChars(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		minChars int
		expected bool
	}{
		{
			name:     "String with exact min length passes",
			value:    "hello",
			minChars: 5,
			expected: true,
		},
		{
			name:     "String over min length passes",
			value:    "hello world",
			minChars: 5,
			expected: true,
		},
		{
			name:     "String under min length fails",
			value:    "hi",
			minChars: 5,
			expected: false,
		},
		{
			name:     "Empty string fails non-zero min length",
			value:    "",
			minChars: 1,
			expected: false,
		},
		{
			name:     "Unicode string with exact min length passes",
			value:    "héllo",
			minChars: 5,
			expected: true,
		},
		{
			name:     "Unicode string under min length fails",
			value:    "hé",
			minChars: 5,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := MinChars(test.value, test.minChars)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}
