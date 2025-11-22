package e2e

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var (
	concatBin string
	optBin    string
)

func TestMain(m *testing.M) {
	// 1. Build Binaries
	tmpBin, err := os.MkdirTemp("", "concat_e2e_bin")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpBin)

	concatBin = filepath.Join(tmpBin, "concat")
	optBin = filepath.Join(tmpBin, "opt")

	if err := buildBin("../cmd/concat", concatBin); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to build concat: %v\n", err)
		os.Exit(1)
	}
	if err := buildBin("../cmd/opt", optBin); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to build opt: %v\n", err)
		os.Exit(1)
	}

	// 2. Run Tests
	os.Exit(m.Run())
}

func buildBin(src, dst string) error {
	cmd := exec.Command("go", "build", "-o", dst, src)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, out)
	}
	return nil
}

func TestConcatOptWorkflow(t *testing.T) {
	// Setup Fixture
	fixtureDir := t.TempDir()
	createFile(t, fixtureDir, "main.go", "package main\n\n\nfunc main() {}")
	createFile(t, fixtureDir, "ignored.log", "should be ignored")
	createFile(t, fixtureDir, ".gitignore", "*.log")

	t.Run("Concat_Stdout_Flag", func(t *testing.T) {
		cmd := exec.Command(concatBin, "-p", "go", "--stdout")
		cmd.Dir = fixtureDir
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Run failed: %v\nOutput: %s", err, out)
		}
		output := string(out)

		if !strings.Contains(output, "package main") {
			t.Error("Output missing main.go content")
		}
		if strings.Contains(output, "ignored.log") {
			t.Error("Output contains ignored file")
		}
	})

	t.Run("Concat_Pipe_Opt", func(t *testing.T) {
		// This tests the critical "Auto-Pipe" logic
		// We construct a pipeline: concat -p go | opt -c
		
		cCmd := exec.Command(concatBin, "-p", "go") // Note: NO --stdout flag!
		cCmd.Dir = fixtureDir
		
oCmd := exec.Command(optBin, "-c") // Compact flag
		
		// Pipe
		reader, writer, _ := os.Pipe()
		cCmd.Stdout = writer
		oCmd.Stdin = reader
		
		var outBuf bytes.Buffer
		var errBuf bytes.Buffer
		oCmd.Stdout = &outBuf
		oCmd.Stderr = &errBuf // Capture opt logs (cost, etc) if any

		// Start Concat
		if err := cCmd.Start(); err != nil {
			t.Fatalf("Failed to start concat: %v", err)
		}
		
		// Start Opt
		if err := oCmd.Start(); err != nil {
			t.Fatalf("Failed to start opt: %v", err)
		}
		
		// Wait for Concat to finish writing to pipe
		if err := cCmd.Wait(); err != nil {
			t.Fatalf("Concat failed: %v", err)
		}
		writer.Close() // Close pipe so Opt sees EOF
		
		// Wait for Opt
		if err := oCmd.Wait(); err != nil {
			t.Fatalf("Opt failed: %v\nStderr: %s", err, errBuf.String())
		}
		
		output := outBuf.String()
		
		// Assertions
		if !strings.Contains(output, "package main") {
			t.Error("Pipe output missing content")
		}
		
		// Verify Compaction (3 newlines -> 2)
		// Original: package main\n\n\nfunc main()
		// Compacted: package main\n\nfunc main()
		if strings.Contains(output, "package main\n\n\nfunc") {
			t.Error("Output was not compacted (found 3 newlines)")
		}
		if !strings.Contains(output, "package main\n\nfunc") {
			t.Error("Output expected compacted structure")
		}
	})
}

func createFile(t *testing.T, dir, name, content string) {
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
