package generator

import (
	"fmt"
	"github.com/ditschedev/swag-ts/internal/config"
	"github.com/ditschedev/swag-ts/internal/templates"
	"github.com/fatih/color"
	"github.com/getkin/kin-openapi/openapi3"
	"html/template"
	"log"
	"os"
	"sort"
	"time"
)

func (g *generator) generateTypescriptTypes() error {
	color.Set(color.FgHiBlack)
	g.s.Suffix = " Generating types page"
	g.s.Start()

	err := generateTypes(g.doc, g.outputPath)
	g.s.Stop()

	if err != nil {
		return fmt.Errorf("failed to generate types: %s", err)
	}

	color.Set(color.FgGreen)
	fmt.Printf("✓ Generated typescript types\n")
	color.Unset()

	color.Set(color.FgHiBlack)
	fmt.Printf("✓ Finished generation in %s \n", time.Since(g.start).Round(time.Microsecond))
	color.Unset()

	return nil
}

type Field struct {
	Name     string
	Type     string
	Nullable bool
	Optional bool
}

type Model struct {
	Name   string
	Type   string
	IsEnum bool
	Fields []Field
	Values []interface{}
}

type TypesPlaceholder struct {
	Models  []Model
	Now     string
	Version string
}

var swaggerTypeToTypescriptType = map[string]string{
	"string":  "string",
	"integer": "number",
	"boolean": "boolean",
	"array":   "any[]",
	"object":  "any",
}

func generateTypes(doc *openapi3.T, outputPath string) error {
	now := time.Now()

	var models []Model

	for name, ref := range doc.Components.Schemas {
		schema := ref.Value

		if schema == nil || (schema.Type != "object" && schema.Enum == nil) {

			log.Printf("Schema %s is not an object or enum, skipping...\n", name)
			continue
		}

		model := Model{
			Name:   name,
			Type:   schema.Type,
			IsEnum: false,
			Fields: make([]Field, 0),
		}

		if schema.Enum != nil {
			model.IsEnum = true
			model.Values = schema.Enum
			models = append(models, model)
			continue
		}

		var requiredFields []string
		if schema.Required != nil {
			requiredFields = schema.Required
		}

		for fieldName, field := range ref.Value.Properties {
			var f Field

			switch field.Value.Type {
			case "array":
				f = Field{
					Name:     fieldName,
					Nullable: field.Value.Nullable,
					Optional: true,
				}

				if field.Value.Items.Ref != "" {
					f.Type = parseRefName(field.Value.Items.Ref) + "[]"
				} else {
					f.Type = swaggerTypeToTypescriptType[field.Value.Items.Value.Type] + "[]"
				}
			case "object":
				f = Field{
					Name:     fieldName,
					Type:     parseRefName(field.Ref),
					Nullable: field.Value.Nullable,
					Optional: true,
				}
			default:
				f = Field{
					Name:     fieldName,
					Type:     swaggerTypeToTypescriptType[field.Value.Type],
					Nullable: field.Value.Nullable,
					Optional: true,
				}

				if field.Ref != "" {
					f.Type = parseRefName(field.Ref)
				}

			}

			for _, requiredField := range requiredFields {
				if requiredField == fieldName {
					f.Optional = false
					break
				}
			}

			model.Fields = append(model.Fields, f)

		}

		// sort fields by name
		sort.Slice(model.Fields, func(i, j int) bool {
			return model.Fields[i].Name < model.Fields[j].Name
		})

		models = append(models, model)
	}

	sort.Slice(models, func(i, j int) bool {
		return models[i].Name < models[j].Name
	})

	typesPlaceholder := TypesPlaceholder{
		Models:  models,
		Now:     now.Format(time.RFC3339),
		Version: config.Version,
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	funcMap := template.FuncMap{
		"raw": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	tmpl, err := template.New("types").Delims("<%", "%>").Funcs(funcMap).ParseFS(templates.FS, "types.ts.tmpl")

	if err != nil {
		return err
	}

	err = tmpl.ExecuteTemplate(f, "types", typesPlaceholder)
	if err != nil {
		return err
	}

	return nil
}
