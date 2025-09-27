package tool

import (
	"testing"
)

func TestStrdup(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic string",
			input:    "hello, world",
			expected: "hello, world",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "string with special characters",
			input:    "a/b\\c\x00d",
			expected: "a/b\\c\x00d",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Strdup(tc.input)
			if result != tc.expected {
				t.Errorf("Strdup(%q) = %q; want %q", tc.input, result, tc.expected)
			}
			// In Go, we can also check if the address is the same, as the underlying
			// data should not be copied for this operation, but this is an
			// implementation detail. The functional correctness is what matters.
		})
	}
}

func TestMemdup0(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "basic byte slice",
			input:    []byte("hello, world"),
			expected: "hello, world",
		},
		{
			name:     "empty byte slice",
			input:    []byte{},
			expected: "",
		},
		{
			name:     "byte slice with null terminator",
			input:    []byte("hello\x00world"),
			expected: "hello\x00world",
		},
		{
			name:     "nil byte slice",
			input:    nil,
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Memdup0(tc.input)
			if result != tc.expected {
				t.Errorf("Memdup0(%v) = %q; want %q", tc.input, result, tc.expected)
			}
		})
	}
}