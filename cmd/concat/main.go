package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/nessaee/concat/internal/app"
	"github.com/nessaee/concat/internal/config"
)

var (
	cfg config.Config
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "concat",
		Short: "Concatenates project files for LLM context",
		Long: `Project Concatenator v1.0.2
Concatenates project files and copies the result to the clipboard or a file.
Designed for easily grabbing project context for LLMs.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(cfg.Extensions) == 0 {
                // Fail if no extensions provided, matching original script behavior
				fmt.Println("Error: You must specify at least one file type to include with -p.")
                cmd.Usage()
				os.Exit(1)
			}

			if err := app.Run(&cfg); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	// Flags
	rootCmd.PersistentFlags().StringSliceVarP(&cfg.Extensions, "pattern", "p", []string{}, "Include files with this extension (e.g., 'py', 'js'). Can be used multiple times.")
	rootCmd.PersistentFlags().StringSliceVarP(&cfg.IgnorePatterns, "ignore", "i", []string{}, "Ignore files or directories matching this pattern. Can be used multiple times.")
	rootCmd.PersistentFlags().StringVarP(&cfg.Output, "output", "o", "", "Output to a file instead of the clipboard.")
	rootCmd.PersistentFlags().BoolVarP(&cfg.IncludeTree, "tree", "t", false, "Include a directory tree structure at the top of the output.")
	rootCmd.PersistentFlags().BoolVarP(&cfg.UseXML, "xml", "x", false, "Format output in XML (<file path='...'>) instead of Markdown.")
    
    // Version flag is automatic with Cobra if we set Version field, but let's leave it for now.

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
