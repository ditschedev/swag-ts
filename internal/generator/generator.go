package generator

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/ditschedev/swag-ts/internal/parser"
	"github.com/fatih/color"
	"github.com/getkin/kin-openapi/openapi3"
	"log"
	"os"
	"time"
)

type Type int

const (
	TypeScript Type = iota
)

type Generator interface {
	GenerateTypes(outputPath string) error
}

type generator struct {
	doc        *openapi3.T
	t          Type
	outputPath string

	start time.Time

	s *spinner.Spinner
}

func NewGenerator(specPath string, t Type) Generator {
	s := spinner.New(spinner.CharSets[40], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Suffix = " Loading OpenAPI spec"
	s.Start()

	color.Set(color.FgHiBlack)
	fmt.Printf("Loading OpenAPI spec from: %s\n", specPath)

	doc, err := parser.LoadAndValidateOpenAPISpec(specPath)
	if err != nil {
		log.Fatalf("Failed to load and validate OpenAPI spec: %s", err)
	}

	s.Stop()

	color.Set(color.FgGreen)
	fmt.Println("✓ OpenAPI spec loaded")
	color.Unset()

	return &generator{
		doc: doc,
		t:   t,
		s:   s,
	}
}

func (g *generator) GenerateTypes(outputPath string) error {

	g.outputPath = outputPath
	g.start = time.Now()

	switch g.t {
	case TypeScript:
		return g.generateTypescriptTypes()
	}

	return nil
}
