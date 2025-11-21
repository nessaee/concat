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
	filter := core.NewFilter(cfg.Extensions, cfg.IgnorePatterns)

	// 2. Initialize Components
	concatenator := core.NewConcatenator(filter)
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
		fmt.Println("> Generating directory tree...")
		treeGen := core.NewTreeGenerator(filter)
		treeStr, err := treeGen.Generate(".")
		if err != nil {
			return fmt.Errorf("failed to generate tree: %w", err)
		}
		outputBuilder += treeStr + "\n---\n\n"
	}

	// 5. Process Files
	fmt.Println("> Searching for files to process...")
	content, count, size, err := concatenator.Process(".")
	if err != nil {
		return fmt.Errorf("processing failed: %w", err)
	}
    outputBuilder += content

	// 6. Output
	if cfg.Output != "" {
		err := os.WriteFile(cfg.Output, []byte(outputBuilder), 0644)
		if err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Printf("✓ Wrote %d files (%d bytes) to '%s'.\n", count, size, cfg.Output)
	} else {
		err := clipboard.WriteAll(outputBuilder)
		if err != nil {
			return fmt.Errorf("failed to copy to clipboard: %w", err)
		}
		fmt.Printf("✓ Copied %d files (%d bytes) to clipboard.\n", count, size)
	}

	return nil
}