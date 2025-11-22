package main

import (
	"fmt"
	"io"
	"os"

	"github.com/nessaee/concat/internal/transform"
	"github.com/spf13/cobra"
)

var (
	flagCompact      bool
	flagStripHeaders bool
	flagCost         bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "opt",
		Short: "Stream optimizer for LLM context",
		Long: `opt (Optimizer) v0.1.0
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
			// If --cost is set, we print analysis to Stderr.
			if flagCost {
				tokens := len(content) / 4 // Approximation
				cost := float64(tokens) * 0.000000075 // Gemini Flash Approx
				fmt.Fprintf(os.Stderr, "Tokens: ~%d | Est Cost: ~$%.5f (Gemini Flash)\n", tokens, cost)
			}

			// 3. Apply Transformations
			transformer := transform.NewTransformer(transform.Options{
				Compact:      flagCompact,
				StripHeaders: flagStripHeaders,
			})
			result := transformer.Process(content)

			// 4. Output
			// Always output result to Stdout unless the user redirects it.
			fmt.Print(result)
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&flagCompact, "compact", "c", false, "Reduce whitespace to save tokens.")
	rootCmd.PersistentFlags().BoolVar(&flagStripHeaders, "strip-headers", false, "Strip copyright/license headers.")
	rootCmd.PersistentFlags().BoolVar(&flagCost, "cost", false, "Estimate tokens and cost (output to stderr).")
	rootCmd.PersistentFlags().BoolVar(&flagCost, "dry-run", false, "Alias for --cost")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
