package app

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/nessaee/concat/internal/config"
	"github.com/nessaee/concat/internal/core"
	"github.com/nessaee/concat/internal/infra"
	"github.com/nessaee/concat/internal/protocol"
)

// Run is the main application entry point
func Run(cfg *config.Config) error {
	// 1. Initialize Filter
	filter := core.NewFilter(cfg.Extensions, cfg.IgnorePatterns, cfg.ExcludeTests)

	// Determine Formatter
	var formatter protocol.Formatter
	if cfg.UseXML {
		formatter = &protocol.XMLFormatter{}
	} else {
		formatter = &protocol.MarkdownFormatter{}
	}

	// 2. Initialize Components
	concatenator := core.NewConcatenator(filter, cfg, formatter)

	// Determine Output Writer
	var outWriter io.Writer
	var clipboardBuffer *bytes.Buffer

	stat, _ := os.Stdout.Stat()
	isPipe := (stat.Mode() & os.ModeCharDevice) == 0

	if cfg.PrintToStdout || isPipe {
		outWriter = os.Stdout
	} else if cfg.Output != "" {
		f, err := os.Create(cfg.Output)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer f.Close()
		outWriter = f
	} else {
		clipboardBuffer = new(bytes.Buffer)
		outWriter = clipboardBuffer
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// 3. Generate Header
	header := fmt.Sprintf("---\nProject: %s\nGenerated: %s\n---\n\n", filepath.Base(cwd), time.Now().Format(time.RFC1123))
	fmt.Fprint(outWriter, header)

	// 4. Generate Tree (Optional)
	if cfg.IncludeTree {
		fmt.Fprintln(os.Stderr, "> Generating directory tree...")
		treeGen := core.NewTreeGenerator(filter)
		treeStr, err := treeGen.Generate(".")
		if err != nil {
			return fmt.Errorf("failed to generate tree: %w", err)
		}
		fmt.Fprint(outWriter, treeStr+"\n---\n\n")
	}

	// 5. Process Files
	fmt.Fprintln(os.Stderr, "> Searching for files to process...")
	count, size, err := concatenator.Process(".", outWriter)
	if err != nil {
		return fmt.Errorf("processing failed: %w", err)
	}

	// 6. Finalize (Clipboard logic)
	estTokens := size / 4

	if clipboardBuffer != nil {
		clipboard := infra.NewClipboard()
		err := clipboard.WriteAll(clipboardBuffer.String())
		if err != nil {
			return fmt.Errorf("failed to copy to clipboard: %w", err)
		}
		fmt.Printf("✓ Copied %d files (%d bytes, ~%d tokens) to clipboard.\n", count, size, estTokens)
	} else if cfg.Output != "" {
		fmt.Printf("✓ Wrote %d files (%d bytes, ~%d tokens) to '%s'.\n", count, size, estTokens, cfg.Output)
	} else {
		// Stdout logic: log to stderr
		fmt.Fprintf(os.Stderr, "✓ Output %d files (%d bytes, ~%d tokens) to stdout.\n", count, size, estTokens)
	}

	return nil
}
