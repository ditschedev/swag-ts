package generator

import "testing"

func TestCreateGenerator(t *testing.T) {
	gen := NewGenerator("../../testdata/test.json", TypeScript)

	if gen == nil {
		t.Fatalf("Expected generator to be created, got nil")
	}
}
