package main

import (
	"os"
	"testing"
)

func TestParseBody(t *testing.T) {

	exampleFull, err := os.ReadFile("examples/full-feedback.json")
	if err != nil {
		t.Errorf("Cannot read example file")
	}
	testCases := []struct {
		name        string
		input       []byte
		expected    string
		expectError bool
	}{
		{
			name:  "valid json",
			input: exampleFull,
			expected: `
Q: What's great (if anything)?
A: I like speed.

Q: What do you find product useful for?
A: To transfer personal files.

Q: What's missing or what's not great?
A: Ability to do multiple file transfer

`,
		},
		{
			name:        "invalid json",
			input:       []byte(`{"Questions": [`),
			expectError: true,
		},
		{
			name:        "json, incorrect template",
			input:       []byte(`{"Test": []}`),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parseBody(tc.input)
			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("expected: %q \n, but got: %q", tc.expected, result)
			}
		})
	}
}
