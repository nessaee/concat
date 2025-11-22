package core

import (
	"testing"
)

func TestFilter_IsTestFile(t *testing.T) {
	filter := &Filter{
		excludeTests: true,
	}

	tests := []struct {
		path     string
		expected bool
	}{
		{"main.go", false},
		{"main_test.go", true},
		{"utils.js", false},
		{"utils.test.js", true},
		{"utils.spec.js", true},
		{"component.ts", false},
		{"component.test.ts", true},
		{"component.spec.ts", true},
		{"script.py", false},
		{"test_script.py", true},
		{"random_file.txt", false},
	}

	for _, tt := range tests {
		if got := filter.IsTestFile(tt.path); got != tt.expected {
			t.Errorf("IsTestFile(%q) = %v; want %v", tt.path, got, tt.expected)
		}
	}
}
