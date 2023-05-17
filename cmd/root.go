package cmd

import (
	"fmt"
	"github.com/ditschedev/swag-ts/internal/config"
	"github.com/ditschedev/swag-ts/internal/parser"
	"github.com/spf13/cobra"
	"os"
)

var (
	version  bool
	specPath string
	output   string
)

var rootCmd = &cobra.Command{
	Use:   "swag-ts",
	Short: "A tiny cli tool that generates typescript types based on a provided OpenAPI Specifications",
	Long:  `Simply provide a OpenAPI Specification and swag-ts will generate typescript types for you. You can provide json or yaml definitions on your local filesystem or a remote url.`,
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			fmt.Printf("swag-ts - v%s (%s)\n", config.Version, config.Date)
			return
		}

		if specPath != "" {
			parser.GenerateTypescriptTypes(specPath, output)
		}

		_ = cmd.Help()
	},
}

func Execute() {
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "shows the version of the cli")
	rootCmd.Flags().StringVarP(&specPath, "file", "f", "", "file path or url to the OpenAPI Specification")

	rootCmd.Flags().StringVarP(&output, "output", "o", "./types/swagger.ts", "output file for generated definitions")

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Failed to execute command")
		os.Exit(1)
	}
}
