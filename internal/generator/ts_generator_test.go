package generator

import (
	"os"
	"strings"
	"testing"
)

func TestTSGeneration(t *testing.T) {
	gen := NewGenerator("../../testdata/test.json", TypeScript)

	if gen == nil {
		t.Fatalf("Expected generator to be created, got nil")
	}

	err := gen.GenerateTypes("../../tmp/testtypes.ts")
	if err != nil {
		t.Fatalf("Expected typescript types to be generated, got error: %s", err)
	}

	data, err := os.ReadFile("../../tmp/testtypes.ts")
	if err != nil {
		t.Fatalf("Expected testtypes.ts to be read, got error: %s", err)
	}

	expected := `
export interface DynamicPayload {
  data: any;
  info?: any | null;
}

export interface LoginFormData {
  certificate: Blob;
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
}

export interface LoginResponseWrapper {
  data: LoginResponse;
  message?: string | null;
}`
	if !strings.Contains(string(data), expected) {
		t.Errorf("Expected testtypes.ts to be:\n%s\n\nGot:\n%s", expected, string(data))
	}
}

func TestNumberTypesGeneration(t *testing.T) {
	gen := NewGenerator("../../testdata/test-numbers.json", TypeScript)

	if gen == nil {
		t.Fatalf("Expected generator to be created, got nil")
	}

	err := gen.GenerateTypes("../../tmp/testtypes.ts")
	if err != nil {
		t.Fatalf("Expected typescript types to be generated, got error: %s", err)
	}

	data, err := os.ReadFile("../../tmp/testtypes.ts")
	if err != nil {
		t.Fatalf("Expected testtypes.ts to be read, got error: %s", err)
	}

	expected := `
export interface NumberTypes {
  anyNum?: number;
  doubleNum: number;
  floatNum?: number;
  int32Num?: number;
  int64Num?: number;
  intNum?: number;
}`
	if !strings.Contains(string(data), expected) {
		t.Errorf("Expected testtypes.ts to be:\n%s\n\nGot:\n%s", expected, string(data))
	}
}
