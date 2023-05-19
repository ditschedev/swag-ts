package cmd

import (
	"fmt"
	"github.com/ditschedev/swag-ts/internal/config"
	"github.com/ditschedev/swag-ts/internal/generator"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	version  bool
	specPath string
	output   string
	genType  generator.Type = generator.TypeScript
)

var RootCmd = &cobra.Command{
	Use:   "swag-ts",
	Short: "A tiny cli tool that generates typescript types based on a provided OpenAPI Specifications",
	Long:  `Simply provide a OpenAPI Specification and swag-ts will generate typescript types for you. You can provide json or yaml definitions on your local filesystem or a remote url.`,
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			fmt.Printf("swag-ts - v%s (%s)\n", config.Version, config.Date)
			return
		}

		if specPath != "" {
			printBanner()

			g := generator.NewGenerator(specPath, genType)

			err := g.GenerateTypes(output)
			if err != nil {
				log.Fatalf("Failed to generate types: %s", err)
			}

			return
		}

		_ = cmd.Help()
	},
}

func init() {
	RootCmd.Flags().BoolVarP(&version, "version", "v", false, "shows the version of the cli")
	RootCmd.Flags().StringVarP(&specPath, "file", "f", "", "file path or url to the OpenAPI Specification")

	RootCmd.Flags().StringVarP(&output, "output", "o", "./types/swagger.ts", "output file for generated definitions")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Printf("Failed to execute command")
		os.Exit(1)
	}
}

func printBanner() {
	color.Set(color.FgHiCyan)
	fmt.Println("                           __    ")
	fmt.Println("  ____    _____ ____ _____/ /____")
	fmt.Println(" (_-< |/|/ / _ `/ _ `/___/ __(_-<")
	fmt.Println("/___/__,__/\\_,_/\\_, /    \\__/___/")
	fmt.Println("               /___/             ")
	color.Unset()
}
