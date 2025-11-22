package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/nessaee/concat/internal/config"
	"github.com/nessaee/concat/internal/protocol"
)

// Concatenator handles finding and reading files
type Concatenator struct {
	filter    *Filter
	config    *config.Config
	formatter protocol.Formatter
}

// NewConcatenator creates a new Concatenator
func NewConcatenator(filter *Filter, cfg *config.Config, formatter protocol.Formatter) *Concatenator {
	return &Concatenator{
		filter:    filter,
		config:    cfg,
		formatter: formatter,
	}
}

// Process walks the directory and returns the formatted content
func (c *Concatenator) Process(root string, w io.Writer) (int, int64, error) {
	var count int

	// Wrap the writer
	cw := &CountingWriter{Writer: w}

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
			// Open file instead of ReadFile
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open %s: %w", path, err)
			}
			defer file.Close()

			// Binary Check: Read small buffer first
			// 8192 bytes (8KB) is a safe bet for detection without reading huge files
			header := make([]byte, 8192)
			n, err := file.Read(header)
			if err != nil && err != io.EOF {
				return fmt.Errorf("failed to read header of %s: %w", path, err)
			}

			if isBinary(header[:n]) {
				fmt.Fprintf(os.Stderr, "⚠ Skipping binary file: %s\n", relPath)
				return nil
			}

			// Reset file pointer to start
			if _, err := file.Seek(0, 0); err != nil {
				return fmt.Errorf("failed to seek %s: %w", path, err)
			}

			c.formatter.WriteHeader(cw, relPath)

			// Copy content to writer
			if _, err := io.Copy(cw, file); err != nil {
				return fmt.Errorf("failed to copy content of %s: %w", path, err)
			}

			c.formatter.WriteFooter(cw)

			count++
		}

		return nil
	})

	return count, cw.Count, err
}

// CountingWriter wraps an io.Writer and counts bytes written
type CountingWriter struct {
	Writer io.Writer
	Count  int64
}

func (cw *CountingWriter) Write(p []byte) (int, error) {
	n, err := cw.Writer.Write(p)
	cw.Count += int64(n)
	return n, err
}
