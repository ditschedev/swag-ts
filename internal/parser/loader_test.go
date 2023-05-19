package parser

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

const addr = "localhost:7965"

func TestLoadLocalOpenAPISpecJson(t *testing.T) {
	const specPath = "../../testdata/test.json"

	doc, err := LoadAndValidateOpenAPISpec(specPath)
	if err != nil {
		t.Errorf("Error loading spec: %v", err)
		return
	}

	if doc.Info.Title != "Swagger Test" {
		t.Errorf("Expected title to be 'Swagger Test', got '%s'", doc.Info.Title)
	}
}

func createTestServer(t *testing.T, handler http.Handler) *httptest.Server {
	ts := httptest.NewUnstartedServer(handler)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatalf("Error loading spec: %v", err)
	}
	ts.Listener.Close()
	ts.Listener = l
	return ts
}

func TestLoadRemoteOpenAPISpecJson(t *testing.T) {
	fs := http.FileServer(http.Dir("../../testdata"))
	ts := createTestServer(t, fs)
	ts.Start()
	defer ts.Close()

	doc, err := LoadAndValidateOpenAPISpec("http://" + addr + "/test.json")
	if err != nil {
		t.Errorf("Error loading spec: %v", err)
		return
	}

	if doc.Info.Title != "Swagger Test" {
		t.Errorf("Expected title to be 'Swagger Test', got '%s'", doc.Info.Title)
	}
}
