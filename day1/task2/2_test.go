package main

import (
	"testing"
	"strings"
)

func TestMostFamousWords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		K        int
		expected string
	}{
		{
			name: "Basic Test",
			input: "apple banana apple orange banana apple",
			K: 2,
			expected: "apple banana",
		},
		{
			name: "K Greater than Unique Words",
			input: "red blue green red blue",
			K: 5,
			expected: "blue red green",
		},
		{
			name: "Empty Input",
			input: "",
			K: 3,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MostFamousWords(tt.K, strings.Fields(tt.input))
			if got != tt.expected {
				t.Errorf("MostFamousWords(%d, %q) = %q; want %q", tt.K, tt.input, got, tt.expected)
			}
		})
	}
}
