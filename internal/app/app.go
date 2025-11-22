package app

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/nessaee/concat/internal/config"
	"github.com/nessaee/concat/internal/core"
	"github.com/nessaee/concat/internal/infra"
)

// Run is the main application entry point
func Run(cfg *config.Config) error {
	// 1. Initialize Filter
	filter := core.NewFilter(cfg.Extensions, cfg.IgnorePatterns, cfg.ExcludeTests)

	// 2. Initialize Components
	concatenator := core.NewConcatenator(filter, cfg)
	clipboard := infra.NewClipboard()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// 3. Generate Header
	var outputBuilder string
	header := fmt.Sprintf("---\nProject: %s\nGenerated: %s\n---\n\n", filepath.Base(cwd), time.Now().Format(time.RFC1123))
	outputBuilder += header

	// 4. Generate Tree (Optional)
	if cfg.IncludeTree {
		fmt.Fprintln(os.Stderr, "> Generating directory tree...")
		treeGen := core.NewTreeGenerator(filter)
		treeStr, err := treeGen.Generate(".")
		if err != nil {
			return fmt.Errorf("failed to generate tree: %w", err)
		}
		outputBuilder += treeStr + "\n---\n\n"
	}

	// 5. Process Files
	fmt.Fprintln(os.Stderr, "> Searching for files to process...")
	content, count, size, err := concatenator.Process(".")
	if err != nil {
		return fmt.Errorf("processing failed: %w", err)
	}
	outputBuilder += content

	// 6. Output
	// Check if stdout is a pipe
	stat, _ := os.Stdout.Stat()
	isPipe := (stat.Mode() & os.ModeCharDevice) == 0

	if cfg.PrintToStdout || isPipe {
		fmt.Print(outputBuilder)
		// Only print to stderr if NOT a pipe (or if user forced stdout)
		// Actually, if we are piping, we definitely want the log to go to Stderr so it doesn't mix with output.
		// But if the user simply ran `concat -s`, they might want to see the log.
		// Let's print the log to Stderr always, as it's safe (won't corrupt the pipe).
		estTokens := size / 4
		fmt.Fprintf(os.Stderr, "✓ Output %d files (%d bytes, ~%d tokens) to stdout.\n", count, size, estTokens)
	} else if cfg.Output != "" {
		err := os.WriteFile(cfg.Output, []byte(outputBuilder), 0644)
		if err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		estTokens := size / 4
		fmt.Printf("✓ Wrote %d files (%d bytes, ~%d tokens) to '%s'.\n", count, size, estTokens, cfg.Output)
	} else {
		err := clipboard.WriteAll(outputBuilder)
		if err != nil {
			return fmt.Errorf("failed to copy to clipboard: %w", err)
		}
		estTokens := size / 4
		fmt.Printf("✓ Copied %d files (%d bytes, ~%d tokens) to clipboard.\n", count, size, estTokens)
	}

	return nil
}
