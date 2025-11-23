package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/nessaee/concat/internal/infra"
	"github.com/nessaee/concat/internal/transform"
	"github.com/spf13/cobra"
)

var (
	flagCompact      bool
	flagStripHeaders bool
	flagCost         bool
	flagStdout       bool
)

// CountingWriter wraps an io.Writer to count bytes
type CountingWriter struct {
	w     io.Writer
	count int64
}

func (cw *CountingWriter) Write(p []byte) (int, error) {
	n, err := cw.w.Write(p)
	cw.count += int64(n)
	return n, err
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "opt",
		Short: "Stream optimizer for LLM context",
		Long: `opt (Optimizer) v0.1.5
Refines text streams for LLM consumption.
Handles cost estimation, whitespace compaction, and license stripping.`,
		Run: func(cmd *cobra.Command, args []string) {
			// 1. Setup Transformer
			transformer := transform.NewTransformer(transform.Options{
				Compact:      flagCompact,
				StripHeaders: flagStripHeaders,
			})

			// 2. Determine Output Strategy
			stat, _ := os.Stdout.Stat()
			isPipe := (stat.Mode() & os.ModeCharDevice) == 0

			if flagStdout || isPipe {
				// Streaming Mode (Stdout)
				writer := &CountingWriter{w: os.Stdout}
				
				if err := transformer.ProcessStream(os.Stdin, writer); err != nil {
					fmt.Fprintf(os.Stderr, "Error processing stream: %v\n", err)
					os.Exit(1)
				}

				if flagCost {
					tokens := writer.count / 4
					fmt.Fprintf(os.Stderr, "Estimated Tokens: ~%d\n", tokens)
				}
			} else {
				// Buffered Mode (Clipboard)
				// We must buffer to copy to clipboard, but we enforce a limit to prevent OOM.
				const ClipboardLimit = 10 * 1024 * 1024 // 10MB limit for clipboard
				
				var buf bytes.Buffer
				writer := &CountingWriter{w: &buf}

				if err := transformer.ProcessStream(os.Stdin, writer); err != nil {
					fmt.Fprintf(os.Stderr, "Error processing stream: %v\n", err)
					os.Exit(1)
				}

				if writer.count > ClipboardLimit {
					fmt.Fprintf(os.Stderr, "Error: Output size (%d bytes) exceeds clipboard limit (10MB).\nUse --stdout to pipe output or redirect to a file.\n", writer.count)
					os.Exit(1)
				}

				// Copy to Clipboard
				clipboard := infra.NewClipboard()
				if err := clipboard.WriteAll(buf.String()); err != nil {
					fmt.Fprintf(os.Stderr, "Error copying to clipboard: %v\nPrinting to stdout instead.\n", err)
					fmt.Print(buf.String())
				} else {
					estTokens := writer.count / 4
					fmt.Fprintf(os.Stderr, "✓ Copied to clipboard (%d bytes, ~%d tokens).\n", writer.count, estTokens)
				}
			}
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&flagCompact, "compact", "c", false, "Reduce whitespace to save tokens.")
	rootCmd.PersistentFlags().BoolVar(&flagStripHeaders, "strip-headers", false, "Strip copyright/license headers.")
	rootCmd.PersistentFlags().BoolVar(&flagCost, "cost", false, "Estimate tokens (output to stderr).")
	rootCmd.PersistentFlags().BoolVar(&flagCost, "dry-run", false, "Alias for --cost")
	rootCmd.PersistentFlags().BoolVarP(&flagStdout, "stdout", "s", false, "Print output to stdout instead of clipboard.")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
