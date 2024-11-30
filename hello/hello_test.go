package hello_test

import (
	"testing"

	"github.com/bit8bytes/toolbox/hello"
)

func TestPurpose(t *testing.T) {
	expected := "We exist to build great things"

	actual := hello.Purpose

	if actual != expected {
		t.Fatalf("Expected '%s', got '%s'", expected, actual)
	}
}
