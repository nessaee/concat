package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nessaee/concat/internal/config"
	"github.com/nessaee/concat/internal/protocol"
)

// Concatenator handles finding and reading files
type Concatenator struct {
	filter *Filter
	config *config.Config
}

// NewConcatenator creates a new Concatenator
func NewConcatenator(filter *Filter, cfg *config.Config) *Concatenator {
	return &Concatenator{filter: filter, config: cfg}
}

// Process walks the directory and returns the formatted content
func (c *Concatenator) Process(root string) (string, int, int64, error) {
	var sb strings.Builder
	var count int
	var totalSize int64

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Handle "."
		if path == root {
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		// We rely on ShouldProcess to handle both Ignores and Extensions
		if !c.filter.ShouldProcess(relPath, d.IsDir()) {
			if d.IsDir() {
				// If directory is ignored, skip it
				if c.filter.IsIgnored(relPath, true) {
					return filepath.SkipDir
				}
				// If directory is not ignored but passed (e.g. just a folder traversal), continue
				return nil
			}
			// File not ignored but extension doesn't match
			return nil
		}

		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read %s: %w", path, err)
			}

			// Binary Check: Skip files that look binary to save tokens and avoid garbage output
			if isBinary(content) {
				fmt.Fprintf(os.Stderr, "⚠ Skipping binary file: %s\n", relPath)
				return nil
			}

			contentStr := string(content)

			if c.config.UseXML {
				sb.WriteString(protocol.FormatHeaderXML(relPath) + "\n")
				sb.WriteString(contentStr)
				sb.WriteString("\n" + protocol.MarkerXMLEnd + "\n")
			} else {
				sb.WriteString(protocol.FormatHeaderMD(relPath) + "\n")
				sb.WriteString(contentStr)
				sb.WriteString("\n\n---\n\n")
			}

			count++
			totalSize += info.Size()
		}

		return nil
	})

	return sb.String(), count, totalSize, err
}
