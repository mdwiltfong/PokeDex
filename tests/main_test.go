package main

import (
	"testing"
)

func TestSanitizeInput(t *testing.T) {
	input := "  HELLO  "
	expected := "hello"
	actual := SanitizeInput(input)
	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}
