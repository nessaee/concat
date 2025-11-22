package transform

import (
	"testing"
)

func TestTransformer_RemoveExcessWhitespace(t *testing.T) {
	transformer := NewTransformer(Options{})
	input := "func foo() {\n\n\n\n\treturn\n}"
	expected := "func foo() {\n\n\treturn\n}"
	got := transformer.removeExcessWhitespace(input)
	if got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

func TestTransformer_StripLicense(t *testing.T) {
	transformer := NewTransformer(Options{})
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "C-style block header",
			input: `/*
 * Copyright (c) 2023
 * License: MIT
 */
package main`,
			expected: "package main",
		},
		{
			name: "Go-style line header",
			input: `// Copyright 2023
// License: MIT

package main`,
			expected: "package main",
		},
		{
			name: "Hash-style header",
			input: `# Copyright 2023
# License: MIT

import os`,
			expected: "import os",
		},
		{
			name: "No header",
			input: `package main
func main() {}`,
			expected: "package main\nfunc main() {}",
		},
	}

	for _, tt := range tests {
		got := transformer.stripLicense(tt.input)
		if got != tt.expected {
			t.Errorf("%s: stripLicense() = %q, want %q", tt.name, got, tt.expected)
		}
	}
}
