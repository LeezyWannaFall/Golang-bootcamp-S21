package main

import (
	"testing"
	"strings"
	"os"
	"io"
)

func TestMostFamousWords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		K        int
		expected string
	}{
		{
			name:	 "Basic Test",
			input:    "apple banana apple orange banana apple\n",
			K:        2,
			expected: "apple banana",
		},
		{
			name:     "K Greater Than Unique Words",
			input:    "red blue green red blue\n",
			K:        5,
			expected: "blue red green",
		},
		{
			name:     "Empty Input",
			input:    "",
			K:        3,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Redirect stdin
			oldStdin := os.Stdin
			r, w, _ := os.Pipe()
			w.WriteString(tt.input)
			w.Close()
			os.Stdin = r

			// Capture stdout
			oldStdout := os.Stdout
			rOut, wOut, _ := os.Pipe()
			os.Stdout = wOut

			var result []WordCount
			MostFamousWords(tt.K, result)

			wOut.Close()
			var outputBuilder strings.Builder
			io.Copy(&outputBuilder, rOut)
			os.Stdout = oldStdout
			os.Stdin = oldStdin

			output := strings.TrimSpace(outputBuilder.String())
			if output != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, output)
			}
		})
	}
}