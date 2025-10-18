package main

import (
	"testing"
	"strings"
	"os"
	"io"
)

func TestPrintresult(t *testing.T) {
	tests := []struct {
		name     string
		K        int
		result   []WordCount
		expected string
	}{
		{
			name:   "Test Case 1",
			K:      3,
			result: []WordCount{{"apple", 4}, {"banana", 2}, {"orange", 2}, {"grape", 1}},
			expected: "apple banana orange",
		},
		{
			name:   "Test Case 2",
			K:      0,
			result: []WordCount{},
			expected: "",
		},
		{
			name:   "Test Case 3 - K greater than length",
			K:      5,
			result: []WordCount{{"red", 2}, {"blue", 2}, {"green", 1}},
			expected: "red blue green",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output strings.Builder
			originalStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			PrintResult(tt.K, tt.result)

			w.Close()
			os.Stdout = originalStdout
			out, _ := io.ReadAll(r)
			output.Write(out)

			if output.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, output.String())
			}
		})
	}
}