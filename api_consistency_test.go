package cpfcnpj

import (
	"errors"
	"testing"
)

// Test API consistency between actual function names and documented function names
// This ensures the documented APIs in CLAUDE.md work as expected

// TestNewValidCPF_APIConsistency tests the documented NewValidCPF function
// The documentation references NewValidCPF but the actual function is NewCpf
func TestNewValidCPF_APIConsistency(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectValid bool
		expectError error
	}{
		{
			name:        "Valid CPF documented API",
			input:       "716.566.867-59",
			expectValid: true,
			expectError: nil,
		},
		{
			name:        "Invalid CPF documented API",
			input:       "11111111111",
			expectValid: false,
			expectError: ErrAllSameDigits,
		},
		{
			name:        "Invalid length CPF documented API",
			input:       "123",
			expectValid: false,
			expectError: ErrCPFInvalidLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the actual API (NewCpf)
			cpf, err := NewCpf(tt.input)

			if tt.expectValid {
				if err != nil {
					t.Errorf("NewCpf(%q) expected valid, got error: %v", tt.input, err)
				}
				if cpf == "" {
					t.Errorf("NewCpf(%q) returned empty CPF", tt.input)
				}
				// Test String method
				formatted := cpf.String()
				if formatted == "" {
					t.Errorf("CPF.String() returned empty string")
				}
			} else {
				if err == nil {
					t.Errorf("NewCpf(%q) expected error, got nil", tt.input)
				}
				if tt.expectError != nil && !errors.Is(err, tt.expectError) {
					t.Errorf("NewCpf(%q) expected error %v, got %v", tt.input, tt.expectError, err)
				}
			}
		})
	}
}

// TestNewValidCNPJ_APIConsistency tests the documented NewValidCNPJ function
// The documentation references NewValidCNPJ but the actual function is NewCnpj
func TestNewValidCNPJ_APIConsistency(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectValid bool
		expectError error
	}{
		{
			name:        "Valid numeric CNPJ documented API",
			input:       "22.796.729/0001-59",
			expectValid: true,
			expectError: nil,
		},
		{
			name:        "Valid alphanumeric CNPJ documented API",
			input:       "12.ABC.345/01DE-35",
			expectValid: true,
			expectError: nil,
		},
		{
			name:        "Invalid CNPJ documented API",
			input:       "00000000000000",
			expectValid: false,
			expectError: ErrAllSameDigits,
		},
		{
			name:        "Invalid length CNPJ documented API",
			input:       "123",
			expectValid: false,
			expectError: ErrCNPJInvalidLength,
		},
		{
			name:        "Invalid format CNPJ documented API",
			input:       "12ABC34501DEAB",
			expectValid: false,
			expectError: ErrCNPJInvalidAlphanumeric,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the actual API (NewCnpj)
			cnpj, err := NewCnpj(tt.input)

			if tt.expectValid {
				if err != nil {
					t.Errorf("NewCnpj(%q) expected valid, got error: %v", tt.input, err)
				}
				if cnpj == "" {
					t.Errorf("NewCnpj(%q) returned empty CNPJ", tt.input)
				}
				// Test String method
				formatted := cnpj.String()
				if formatted == "" {
					t.Errorf("CNPJ.String() returned empty string")
				}
			} else {
				if err == nil {
					t.Errorf("NewCnpj(%q) expected error, got nil", tt.input)
				}
				if tt.expectError != nil && !errors.Is(err, tt.expectError) {
					t.Errorf("NewCnpj(%q) expected error %v, got %v", tt.input, tt.expectError, err)
				}
			}
		})
	}
}

// Test documented usage patterns from CLAUDE.md
func TestDocumentedUsagePatterns(t *testing.T) {
	// Test CPF usage as documented
	t.Run("CPF documented usage", func(t *testing.T) {
		cpf, err := NewCpf("716.566.867-59")
		if err != nil {
			t.Errorf("Documented CPF example should work: %v", err)
		}

		formatted := cpf.String()
		expected := "716.566.867-59"
		if formatted != expected {
			t.Errorf("CPF.String() = %q, want %q", formatted, expected)
		}
	})

	// Test CNPJ numeric usage as documented
	t.Run("CNPJ numeric documented usage", func(t *testing.T) {
		cnpj, err := NewCnpj("22.796.729/0001-59")
		if err != nil {
			t.Errorf("Documented numeric CNPJ example should work: %v", err)
		}

		formatted := cnpj.String()
		expected := "22.796.729/0001-59"
		if formatted != expected {
			t.Errorf("CNPJ.String() = %q, want %q", formatted, expected)
		}
	})

	// Test CNPJ alphanumeric usage as documented
	t.Run("CNPJ alphanumeric documented usage", func(t *testing.T) {
		cnpj, err := NewCnpj("12.ABC.345/01DE-35")
		if err != nil {
			t.Errorf("Documented alphanumeric CNPJ example should work: %v", err)
		}

		formatted := cnpj.String()
		expected := "12.ABC.345/01DE-35"
		if formatted != expected {
			t.Errorf("CNPJ.String() = %q, want %q", formatted, expected)
		}
	})
}

// Test Clean() function as documented in usage examples
func TestClean_DocumentedUsage(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "CPF cleaning documented example",
			input:    "716.566.867-59",
			expected: "71656686759",
		},
		{
			name:     "CNPJ numeric cleaning documented example",
			input:    "22.796.729/0001-59",
			expected: "22796729000159",
		},
		{
			name:     "CNPJ alphanumeric cleaning documented example",
			input:    "12.ABC.345/01DE-35",
			expected: "12ABC34501DE35",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Clean(tt.input)
			if result != tt.expected {
				t.Errorf("Clean(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// Test ProcessDocumentInput pattern from documentation
func TestProcessDocumentInput_DocumentedPattern(t *testing.T) {
	// This implements the recommended pattern from CLAUDE.md
	ProcessDocumentInput := func(input string) (interface{}, error) {
		// Clean and auto-detect document type
		cleaned := Clean(input)

		// Validate based on length
		switch len(cleaned) {
		case 11:
			return NewCpf(cleaned) // CPF
		case 14:
			return NewCnpj(cleaned) // CNPJ
		default:
			return nil, errors.New("invalid document length")
		}
	}

	tests := []struct {
		name        string
		input       string
		expectType  string
		expectError bool
	}{
		{
			name:        "CPF input processed",
			input:       "716.566.867-59",
			expectType:  "CPF",
			expectError: false,
		},
		{
			name:        "CNPJ numeric input processed",
			input:       "22.796.729/0001-59",
			expectType:  "CNPJ",
			expectError: false,
		},
		{
			name:        "CNPJ alphanumeric input processed",
			input:       "12.ABC.345/01DE-35",
			expectType:  "CNPJ",
			expectError: false,
		},
		{
			name:        "Invalid input processed",
			input:       "12345",
			expectType:  "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ProcessDocumentInput(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("ProcessDocumentInput(%q) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("ProcessDocumentInput(%q) expected success, got error: %v", tt.input, err)
				return
			}

			switch tt.expectType {
			case "CPF":
				if _, ok := result.(CPF); !ok {
					t.Errorf("ProcessDocumentInput(%q) expected CPF type, got %T", tt.input, result)
				}
			case "CNPJ":
				if _, ok := result.(CNPJ); !ok {
					t.Errorf("ProcessDocumentInput(%q) expected CNPJ type, got %T", tt.input, result)
				}
			}
		})
	}
}

// Test all documented error types are available and work correctly
func TestDocumentedErrorTypes(t *testing.T) {
	errorTests := []struct {
		name        string
		input       string
		constructor func(string) (interface{}, error)
		expectedErr error
	}{
		// CPF errors
		{
			name:        "CPF invalid length",
			input:       "123",
			constructor: func(s string) (interface{}, error) { return NewCpf(s) },
			expectedErr: ErrCPFInvalidLength,
		},
		{
			name:        "CPF all same digits",
			input:       "11111111111",
			constructor: func(s string) (interface{}, error) { return NewCpf(s) },
			expectedErr: ErrAllSameDigits,
		},
		{
			name:        "CPF invalid checksum",
			input:       "12345678901",
			constructor: func(s string) (interface{}, error) { return NewCpf(s) },
			expectedErr: ErrCPFInvalidChecksum,
		},
		// CNPJ errors
		{
			name:        "CNPJ invalid length",
			input:       "123",
			constructor: func(s string) (interface{}, error) { return NewCnpj(s) },
			expectedErr: ErrCNPJInvalidLength,
		},
		{
			name:        "CNPJ all same digits",
			input:       "11111111111111",
			constructor: func(s string) (interface{}, error) { return NewCnpj(s) },
			expectedErr: ErrAllSameDigits,
		},
		{
			name:        "CNPJ invalid checksum",
			input:       "12ABC34501DE99",
			constructor: func(s string) (interface{}, error) { return NewCnpj(s) },
			expectedErr: ErrCNPJInvalidChecksum,
		},
		{
			name:        "CNPJ invalid alphanumeric",
			input:       "12ABC34501DEAB",
			constructor: func(s string) (interface{}, error) { return NewCnpj(s) },
			expectedErr: ErrCNPJInvalidAlphanumeric,
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.constructor(tt.input)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", tt.name)
				return
			}

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("Expected error %v for %s, got %v", tt.expectedErr, tt.name, err)
			}
		})
	}
}

// Test that all documented examples work exactly as shown
func TestDocumentationExamples(t *testing.T) {
	// Test exact examples from CLAUDE.md
	examples := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "CPF validation example",
			testFunc: func(t *testing.T) {
				cpf, err := NewCpf("716.566.867-59")
				if err != nil {
					t.Errorf("CPF example failed: %v", err)
					return
				}
				result := cpf.String()
				expected := "716.566.867-59"
				if result != expected {
					t.Errorf("CPF.String() = %q, want %q", result, expected)
				}
			},
		},
		{
			name: "CNPJ numeric validation example",
			testFunc: func(t *testing.T) {
				cnpj, err := NewCnpj("22.796.729/0001-59")
				if err != nil {
					t.Errorf("CNPJ numeric example failed: %v", err)
					return
				}
				result := cnpj.String()
				expected := "22.796.729/0001-59"
				if result != expected {
					t.Errorf("CNPJ.String() = %q, want %q", result, expected)
				}
			},
		},
		{
			name: "CNPJ alphanumeric validation example",
			testFunc: func(t *testing.T) {
				cnpj, err := NewCnpj("12.ABC.345/01DE-35")
				if err != nil {
					t.Errorf("CNPJ alphanumeric example failed: %v", err)
					return
				}
				result := cnpj.String()
				expected := "12.ABC.345/01DE-35"
				if result != expected {
					t.Errorf("CNPJ.String() = %q, want %q", result, expected)
				}
			},
		},
		{
			name: "Clean function examples",
			testFunc: func(t *testing.T) {
				tests := []struct {
					input    string
					expected string
				}{
					{"716.566.867-59", "71656686759"},
					{"12.ABC.345/01DE-35", "12ABC34501DE35"},
				}

				for _, tt := range tests {
					result := Clean(tt.input)
					if result != tt.expected {
						t.Errorf("Clean(%q) = %q, want %q", tt.input, result, tt.expected)
					}
				}
			},
		},
	}

	for _, example := range examples {
		t.Run(example.name, example.testFunc)
	}
}
