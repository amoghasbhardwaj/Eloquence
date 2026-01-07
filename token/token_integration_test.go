// ==============================================================================================
// FILE: token/token_integration_test.go
// ==============================================================================================
// PURPOSE: Tests the integration of the keyword map with the lookup function across various
//          categories of keywords to ensure no category is missing.
// ==============================================================================================

package token

import "testing"

func TestIntegrationKeywordCategories(t *testing.T) {
	// We categorize tests to ensure broad coverage of the language features.
	categories := map[string][]struct {
		input string
		want  TokenType
	}{
		"Math": {
			{"adds", ADDS},
			{"times", TIMES},
			{"modulo", MODULO},
		},
		"Logic": {
			{"and", AND},
			{"or", OR},
			{"not", NOT},
		},
		"Control Flow": {
			{"if", IF},
			{"else", ELSE},
			{"while", WHILE},
			{"in", IN}, // Range loop support
		},
		"Pointers": {
			{"pointing to", POINTING_TO},
			{"pointing from", POINTING_FROM},
		},
		"Structures": {
			{"define", DEFINE},
			{"as", AS},
			{"struct", STRUCT},
		},
		"Modules": {
			{"include", INCLUDE},
		},
	}

	for category, tests := range categories {
		t.Run(category, func(t *testing.T) {
			for _, tt := range tests {
				got := LookupIdent(tt.input)
				if got != tt.want {
					t.Errorf("FAIL [%s]: LookupIdent(%q) = %q, want %q", category, tt.input, got, tt.want)
				}
			}
		})
	}
}
