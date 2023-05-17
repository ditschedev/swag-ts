package parser

import (
	"github.com/ditschedev/swag-ts/internal/templates"
	"github.com/getkin/kin-openapi/openapi3"
	"html/template"
	"log"
	"os"
	"strings"
	"time"
)

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
	Models []Model
	Now    string
}

var swaggerTypeToTypescriptType = map[string]string{
	"string":  "string",
	"integer": "number",
	"boolean": "boolean",
	"array":   "any[]",
	"object":  "any",
}

func generateTypes(doc *openapi3.T, outputPath string) {
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

				break
			case "object":
				f = Field{
					Name:     fieldName,
					Type:     parseRefName(field.Ref),
					Nullable: field.Value.Nullable,
					Optional: true,
				}
				break
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

			if requiredFields != nil {
				for _, requiredField := range requiredFields {
					if requiredField == fieldName {
						f.Optional = false
						break
					}
				}
			}

			model.Fields = append(model.Fields, f)

		}

		models = append(models, model)
	}

	typesPlaceholder := TypesPlaceholder{
		Models: models,
		Now:    now.Format(time.RFC3339),
	}

	f, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}

	funcMap := template.FuncMap{
		"raw": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	tmpl, err := template.New("types").Delims("<%", "%>").Funcs(funcMap).ParseFS(templates.FS, "types.ts.tmpl")

	if err != nil {
		panic(err)
	}

	err = tmpl.ExecuteTemplate(f, "types", typesPlaceholder)
	if err != nil {
		panic(err)
	}
}

func parseRefName(input string) string {
	segments := strings.Split(input, "/")
	return segments[len(segments)-1]
}
