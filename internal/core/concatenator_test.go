package core

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nessaee/concat/internal/config"
)

func TestConcatenator_Process(t *testing.T) {
	// Setup temporary directory
	tmpDir, err := os.MkdirTemp("", "concat_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	files := map[string]string{
		"main.go": `package main

func main() {


}
`,
		"binary.bin": string([]byte{0x00, 0x01, 0x02}), // Null byte = binary
		"ignored.log": "log content",
	}

	for name, content := range files {
		path := filepath.Join(tmpDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Setup Filter and Config
	cfg := &config.Config{
		Extensions:     []string{"go", "bin", "log"},
		IgnorePatterns: []string{"*.log"},
	}
	filter := NewFilter(cfg.Extensions, cfg.IgnorePatterns, false)
	concatenator := NewConcatenator(filter, cfg)

	// Run Process
	output, count, _, err := concatenator.Process(tmpDir)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	// Assertions
	if count != 1 {
		t.Errorf("Expected 1 processed file (main.go), got %d", count)
	}

	// Check if binary was skipped (implied by count=1 if logs are ignored)
	if strings.Contains(output, "binary.bin") {
		t.Error("Output contains binary file that should have been skipped")
	}

	// Check if log was ignored
	if strings.Contains(output, "ignored.log") {
		t.Error("Output contains ignored log file")
	}

	// Check Content Integrity (No Compact Check anymore)
	if !strings.Contains(output, "func main() {") {
		t.Error("Output missing main function")
	}
}

func TestIsBinary(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected bool
	}{
		{"Text", []byte("Hello World"), false},
		{"Binary", []byte{0x00, 0xFF}, true},
		{"LargeText", bytes.Repeat([]byte("A"), 10000), false},
		{"LargeBinary", append(bytes.Repeat([]byte("A"), 100), 0x00), true}, // Null byte within limit
	}

	for _, tt := range tests {
		if got := isBinary(tt.content); got != tt.expected {
			t.Errorf("%s: expected %v, got %v", tt.name, tt.expected, got)
		}
	}
}
