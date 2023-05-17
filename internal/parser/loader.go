package parser

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/getkin/kin-openapi/openapi3"
	"net/url"
	"os"
	"time"
)

func GenerateTypescriptTypes(specPath, outputPath string) {
	if outputPath == "" {
		outputPath = "./types/swagger.ts"
	}

	start := time.Now()
	color.Set(color.FgHiCyan)
	fmt.Println("                           __    ")
	fmt.Println("  ____    _____ ____ _____/ /____")
	fmt.Println(" (_-< |/|/ / _ `/ _ `/___/ __(_-<")
	fmt.Println("/___/__,__/\\_,_/\\_, /    \\__/___/")
	fmt.Println("               /___/             ")
	color.Unset()

	s := spinner.New(spinner.CharSets[40], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Suffix = " Loading OpenAPI spec"
	s.Start()

	color.Set(color.FgHiBlack)
	fmt.Printf("Loading OpenAPI spec from: %s\n", specPath)

	loader := openapi3.NewLoader()

	specUrl, err := url.Parse(specPath)

	var doc *openapi3.T
	if err == nil && specUrl != nil {
		doc, err = loader.LoadFromURI(specUrl)
		if err != nil {
			panic(err)
		}
	} else {
		doc, err = loader.LoadFromFile(specPath)
		if err != nil {
			panic(err)
		}
	}

	if err = doc.Validate(loader.Context); err != nil {
		panic(err)
	}

	s.Stop()

	color.Set(color.FgGreen)
	fmt.Println("✓ OpenAPI spec loaded")
	color.Unset()

	color.Set(color.FgHiBlack)
	s.Suffix = " Generating types page"
	s.Start()

	generateTypes(doc, outputPath)

	s.Stop()

	color.Set(color.FgGreen)
	fmt.Printf("✓ Generated typescript types\n")
	color.Unset()

	color.Set(color.FgHiBlack)
	fmt.Printf("✓ Finished generation in %s \n", time.Since(start).Round(time.Microsecond))
	color.Unset()
}
