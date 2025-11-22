package main

import (
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

func main() {
	rootCmd := &cobra.Command{
		Use:   "opt",
		Short: "Stream optimizer for LLM context",
		Long: `opt (Optimizer) v0.1.3
Refines text streams for LLM consumption.
Handles cost estimation, whitespace compaction, and license stripping.`,
		Run: func(cmd *cobra.Command, args []string) {
			// 1. Read Stream
			input, err := io.ReadAll(os.Stdin)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
				os.Exit(1)
			}
			content := string(input)

			// 2. Handle Analysis (Cost)
			if flagCost {
				tokens := len(content) / 4
				fmt.Fprintf(os.Stderr, "Tokens: ~%d\n", tokens)
			}

			// 3. Apply Transformations
			transformer := transform.NewTransformer(transform.Options{
				Compact:      flagCompact,
				StripHeaders: flagStripHeaders,
			})
			result := transformer.Process(content)

			// 4. Output Strategy
			// Check if stdout is a pipe
			stat, _ := os.Stdout.Stat()
			isPipe := (stat.Mode() & os.ModeCharDevice) == 0

			if flagStdout || isPipe {
				fmt.Print(result)
				// If strictly piping, we might not want logs to stderr if it's being consumed programmatically,
				// but Concat does print to stderr. Let's stick to Concat behavior: simple logging.
				// Actually Concat ONLY logs "Output X files" to stderr if !isPipe.
				// If it IS a pipe, it stays silent (except for errors).
				// Let's follow that.
			} else {
				// TTY -> Clipboard
				clipboard := infra.NewClipboard()
				if err := clipboard.WriteAll(result); err != nil {
					fmt.Fprintf(os.Stderr, "Error copying to clipboard: %v\nPrinting to stdout instead.\n", err)
					fmt.Print(result)
				} else {
					estTokens := len(result) / 4
					fmt.Fprintf(os.Stderr, "✓ Copied to clipboard (%d bytes, ~%d tokens).\n", len(result), estTokens)
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
