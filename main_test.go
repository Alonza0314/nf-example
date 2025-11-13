package main

import "testing"

func TestSum(t *testing.T) {
	a, b := 2, 3
	result := a + b
	if result != 5 {
		t.Errorf("Expected 5, got %d", result)
	}
}
