package slices

import (
	"reflect"
	"testing"
)

// TestMapEmpty tests the Map function with an empty slice
func TestMapEmpty(t *testing.T) {
	t.Parallel()

	result := Map([]int{}, func(n int) int {
		return n * n
	})

	if len(result) != 0 {
		t.Errorf("Expected empty slice, got %v", result)
	}
}

// TestMapSquare tests the Map function with a slice of integers, squaring each element
func TestMapSquare(t *testing.T) {
	t.Parallel()

	src := []int{1, 2, 3, 4, 5}
	expected := []int{1, 4, 9, 16, 25}

	result := Map(src, func(n int) int {
		return n * n
	})

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestMapStringLength tests the Map function with a slice of strings, getting the length of each
func TestMapStringLength(t *testing.T) {
	t.Parallel()

	src := []string{"apple", "banana", "cherry"}
	expected := []int{5, 6, 6}

	result := Map(src, func(s string) int {
		return len(s)
	})

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
