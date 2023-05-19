package generator

import "testing"

func TestCreateGenerator(t *testing.T) {
	gen := NewGenerator("../../testdata/test.json", TypeScript)

	if gen == nil {
		t.Fatalf("Expected generator to be created, got nil")
	}
}

func TestParseRefName(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{input: "hello/world", expected: "world"},
		{input: "#/components/schemas/LoginResponse", expected: "LoginResponse"},
		{input: "nosegments", expected: "nosegments"},
		{input: "", expected: ""},
		{input: "with/empty/segment/", expected: ""},
	}

	for _, tc := range testCases {
		result := parseRefName(tc.input)
		if result != tc.expected {
			t.Errorf("parseRefName(%q) = %q; want %q", tc.input, result, tc.expected)
		}
	}
}
