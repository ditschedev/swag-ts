package parser

import (
	"github.com/getkin/kin-openapi/openapi3"
	"net/url"
)

func LoadAndValidateOpenAPISpec(specPath string) (*openapi3.T, error) {
	l := openapi3.NewLoader()

	specUrl, err := url.Parse(specPath)

	var doc *openapi3.T
	if err == nil && specUrl != nil {
		doc, err = l.LoadFromURI(specUrl)
		if err != nil {
			return nil, err
		}
	} else {
		doc, err = l.LoadFromFile(specPath)
		if err != nil {
			return nil, err
		}
	}

	if err = doc.Validate(l.Context); err != nil {
		return nil, err
	}

	return doc, nil
}
