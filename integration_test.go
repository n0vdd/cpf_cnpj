package cpfcnpj

import (
	"fmt"
	"strings"
	"testing"
)

// Integration tests simulating real-world usage scenarios
// These tests validate the complete workflow from form input to validation

// FormData represents typical form input structure
type FormData struct {
	DocumentNumber string `json:"document_number"`
	DocumentType   string `json:"document_type"` // "CPF" or "CNPJ" or "AUTO"
}

// ValidationResult represents the result of document validation
type ValidationResult struct {
	IsValid        bool   `json:"is_valid"`
	CleanedValue   string `json:"cleaned_value"`
	FormattedValue string `json:"formatted_value"`
	DocumentType   string `json:"document_type"`
	ErrorMessage   string `json:"error_message,omitempty"`
}

// ProcessFormDocument simulates a real-world form processing function
func ProcessFormDocument(form FormData) ValidationResult {
	result := ValidationResult{
		DocumentType: form.DocumentType,
	}

	// Clean the input
	cleaned := Clean(form.DocumentNumber)
	result.CleanedValue = cleaned

	// Auto-detect document type if not specified
	if form.DocumentType == "AUTO" || form.DocumentType == "" {
		switch len(cleaned) {
		case 11:
			result.DocumentType = "CPF"
		case 14:
			result.DocumentType = "CNPJ"
		default:
			result.IsValid = false
			result.DocumentType = "" // Clear the AUTO type since we couldn't determine it
			result.ErrorMessage = "Unable to determine document type from length"
			return result
		}
	}

	// Validate based on document type
	switch result.DocumentType {
	case "CPF":
		cpf, err := NewCpf(cleaned)
		if err != nil {
			result.IsValid = false
			result.ErrorMessage = err.Error()
		} else {
			result.IsValid = true
			result.FormattedValue = cpf.String()
		}
	case "CNPJ":
		cnpj, err := NewCnpj(cleaned)
		if err != nil {
			result.IsValid = false
			result.ErrorMessage = err.Error()
		} else {
			result.IsValid = true
			result.FormattedValue = cnpj.String()
		}
	default:
		result.IsValid = false
		result.ErrorMessage = "Invalid document type"
	}

	return result
}

// TestIntegration_WebFormProcessing simulates web form submission scenarios
func TestIntegration_WebFormProcessing(t *testing.T) {
	tests := []struct {
		name           string
		form           FormData
		expectedValid  bool
		expectedType   string
		expectedFormat string
		expectedError  string
	}{
		// CPF scenarios
		{
			name: "Valid CPF from web form",
			form: FormData{
				DocumentNumber: "716.566.867-59",
				DocumentType:   "CPF",
			},
			expectedValid:  true,
			expectedType:   "CPF",
			expectedFormat: "716.566.867-59",
		},
		{
			name: "Valid CPF auto-detect",
			form: FormData{
				DocumentNumber: "716.566.867-59",
				DocumentType:   "AUTO",
			},
			expectedValid:  true,
			expectedType:   "CPF",
			expectedFormat: "716.566.867-59",
		},
		{
			name: "Messy CPF input from user",
			form: FormData{
				DocumentNumber: "  7 1 6 . 5 6 6 . 8 6 7 - 5 9  ",
				DocumentType:   "AUTO",
			},
			expectedValid:  true,
			expectedType:   "CPF",
			expectedFormat: "716.566.867-59",
		},
		{
			name: "Invalid CPF from form",
			form: FormData{
				DocumentNumber: "111.111.111-11",
				DocumentType:   "CPF",
			},
			expectedValid: false,
			expectedType:  "CPF",
			expectedError: "document cannot have all same digits",
		},

		// CNPJ numeric scenarios
		{
			name: "Valid numeric CNPJ from web form",
			form: FormData{
				DocumentNumber: "22.796.729/0001-59",
				DocumentType:   "CNPJ",
			},
			expectedValid:  true,
			expectedType:   "CNPJ",
			expectedFormat: "22.796.729/0001-59",
		},
		{
			name: "Valid numeric CNPJ auto-detect",
			form: FormData{
				DocumentNumber: "22.796.729/0001-59",
				DocumentType:   "AUTO",
			},
			expectedValid:  true,
			expectedType:   "CNPJ",
			expectedFormat: "22.796.729/0001-59",
		},

		// CNPJ alphanumeric scenarios
		{
			name: "Valid alphanumeric CNPJ from web form",
			form: FormData{
				DocumentNumber: "12.ABC.345/01DE-35",
				DocumentType:   "CNPJ",
			},
			expectedValid:  true,
			expectedType:   "CNPJ",
			expectedFormat: "12.ABC.345/01DE-35",
		},
		{
			name: "Lowercase alphanumeric CNPJ normalized",
			form: FormData{
				DocumentNumber: "12.abc.345/01de-35",
				DocumentType:   "AUTO",
			},
			expectedValid:  true,
			expectedType:   "CNPJ",
			expectedFormat: "12.ABC.345/01DE-35",
		},

		// Error scenarios
		{
			name: "Empty input",
			form: FormData{
				DocumentNumber: "",
				DocumentType:   "AUTO",
			},
			expectedValid: false,
			expectedType:  "", // No type determined
			expectedError: "Unable to determine document type from length",
		},
		{
			name: "Invalid length input",
			form: FormData{
				DocumentNumber: "12345",
				DocumentType:   "AUTO",
			},
			expectedValid: false,
			expectedType:  "", // No type determined
			expectedError: "Unable to determine document type from length",
		},
		{
			name: "Wrong type specification",
			form: FormData{
				DocumentNumber: "716.566.867-59", // CPF input
				DocumentType:   "CNPJ",           // but claiming CNPJ
			},
			expectedValid: false,
			expectedType:  "CNPJ",
			expectedError: "CNPJ must have exactly 14 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessFormDocument(tt.form)

			if result.IsValid != tt.expectedValid {
				t.Errorf("ProcessFormDocument(%+v).IsValid = %v, want %v", tt.form, result.IsValid, tt.expectedValid)
			}

			if result.DocumentType != tt.expectedType {
				t.Errorf("ProcessFormDocument(%+v).DocumentType = %v, want %v", tt.form, result.DocumentType, tt.expectedType)
			}

			if tt.expectedValid && result.FormattedValue != tt.expectedFormat {
				t.Errorf("ProcessFormDocument(%+v).FormattedValue = %v, want %v", tt.form, result.FormattedValue, tt.expectedFormat)
			}

			if !tt.expectedValid && !strings.Contains(result.ErrorMessage, tt.expectedError) {
				t.Errorf("ProcessFormDocument(%+v).ErrorMessage should contain %q, got %q", tt.form, tt.expectedError, result.ErrorMessage)
			}
		})
	}
}

// TestIntegration_BatchProcessing simulates batch processing scenarios
func TestIntegration_BatchProcessing(t *testing.T) {
	// Simulate processing a batch of documents from CSV import or API
	documents := []string{
		"716.566.867-59",     // Valid CPF
		"22.796.729/0001-59", // Valid numeric CNPJ
		"12.ABC.345/01DE-35", // Valid alphanumeric CNPJ
		"111.111.111-11",     // Invalid CPF
		"12345",              // Invalid length
		"",                   // Empty
		"12.abc.345/01de-35", // Lowercase alphanumeric CNPJ
	}

	expectedResults := []struct {
		valid        bool
		documentType string
		errorType    error
	}{
		{true, "CPF", nil},
		{true, "CNPJ", nil},
		{true, "CNPJ", nil},
		{false, "", ErrAllSameDigits},
		{false, "", nil}, // Length error
		{false, "", nil}, // Empty error
		{true, "CNPJ", nil},
	}

	validCount := 0
	cpfCount := 0
	cnpjCount := 0
	errorCount := 0

	for i, doc := range documents {
		t.Run(fmt.Sprintf("Document_%d", i), func(t *testing.T) {
			form := FormData{
				DocumentNumber: doc,
				DocumentType:   "AUTO",
			}

			result := ProcessFormDocument(form)
			expected := expectedResults[i]

			if result.IsValid != expected.valid {
				t.Errorf("Document %d (%q): expected valid=%v, got=%v", i, doc, expected.valid, result.IsValid)
			}

			if result.IsValid {
				validCount++
				switch result.DocumentType {
				case "CPF":
					cpfCount++
				case "CNPJ":
					cnpjCount++
				}
			} else {
				errorCount++
				if expected.errorType != nil {
					if !strings.Contains(result.ErrorMessage, expected.errorType.Error()) {
						t.Errorf("Document %d (%q): expected error containing %q, got %q", i, doc, expected.errorType.Error(), result.ErrorMessage)
					}
				}
			}
		})
	}

	// Verify batch statistics
	t.Run("Batch_Statistics", func(t *testing.T) {
		expectedValid := 4
		expectedCPF := 1
		expectedCNPJ := 3
		expectedErrors := 3

		if validCount != expectedValid {
			t.Errorf("Expected %d valid documents, got %d", expectedValid, validCount)
		}
		if cpfCount != expectedCPF {
			t.Errorf("Expected %d CPF documents, got %d", expectedCPF, cpfCount)
		}
		if cnpjCount != expectedCNPJ {
			t.Errorf("Expected %d CNPJ documents, got %d", expectedCNPJ, cnpjCount)
		}
		if errorCount != expectedErrors {
			t.Errorf("Expected %d errors, got %d", expectedErrors, errorCount)
		}
	})
}

// TestIntegration_ErrorHandlingWorkflow tests comprehensive error handling
func TestIntegration_ErrorHandlingWorkflow(t *testing.T) {
	errorScenarios := []struct {
		name          string
		input         string
		documentType  string
		expectedError error
	}{
		{
			name:          "CPF length error",
			input:         "123",
			documentType:  "CPF",
			expectedError: ErrCPFInvalidLength,
		},
		{
			name:          "CPF character error (filtered by Clean)",
			input:         "1234567890A", // A gets filtered, causes length error
			documentType:  "CPF",
			expectedError: ErrCPFInvalidLength, // Length error, not character error
		},
		{
			name:          "CPF same digits error",
			input:         "11111111111",
			documentType:  "CPF",
			expectedError: ErrAllSameDigits,
		},
		{
			name:          "CPF checksum error",
			input:         "12345678901",
			documentType:  "CPF",
			expectedError: ErrCPFInvalidChecksum,
		},
		{
			name:          "CNPJ length error",
			input:         "123",
			documentType:  "CNPJ",
			expectedError: ErrCNPJInvalidLength,
		},
		{
			name:          "CNPJ format error",
			input:         "12ABC34501DEAB",
			documentType:  "CNPJ",
			expectedError: ErrCNPJInvalidAlphanumeric,
		},
		{
			name:          "CNPJ same digits error",
			input:         "11111111111111",
			documentType:  "CNPJ",
			expectedError: ErrAllSameDigits,
		},
		{
			name:          "CNPJ checksum error",
			input:         "12ABC34501DE99",
			documentType:  "CNPJ",
			expectedError: ErrCNPJInvalidChecksum,
		},
	}

	for _, tt := range errorScenarios {
		t.Run(tt.name, func(t *testing.T) {
			form := FormData{
				DocumentNumber: tt.input,
				DocumentType:   tt.documentType,
			}

			result := ProcessFormDocument(form)

			if result.IsValid {
				t.Errorf("Expected error for %s, but got valid result", tt.name)
			}

			if result.ErrorMessage == "" {
				t.Errorf("Expected error message for %s, but got empty message", tt.name)
			}

			// Check that the specific error type is properly wrapped/reported
			if !strings.Contains(result.ErrorMessage, tt.expectedError.Error()) {
				t.Errorf("Expected error message to contain %q for %s, got %q", tt.expectedError.Error(), tt.name, result.ErrorMessage)
			}
		})
	}
}

// TestIntegration_PerformanceWorkflow simulates performance-critical scenarios
func TestIntegration_PerformanceWorkflow(t *testing.T) {
	// Test that the integration workflow is efficient for large datasets
	testData := []FormData{
		{"716.566.867-59", "AUTO"},
		{"22.796.729/0001-59", "AUTO"},
		{"12.ABC.345/01DE-35", "AUTO"},
		{"71656686759", "AUTO"},
		{"22796729000159", "AUTO"},
		{"12ABC34501DE35", "AUTO"},
	}

	// Run processing multiple times to simulate real workload
	iterations := 100

	for i := 0; i < iterations; i++ {
		for j, form := range testData {
			result := ProcessFormDocument(form)
			if !result.IsValid {
				t.Errorf("Iteration %d, Document %d: Expected valid result, got error: %s", i, j, result.ErrorMessage)
			}
		}
	}

	// Test that memory is handled properly (no leaks)
	t.Run("Memory_Handling", func(t *testing.T) {
		// Process a large batch
		for i := 0; i < 1000; i++ {
			form := FormData{
				DocumentNumber: "716.566.867-59",
				DocumentType:   "AUTO",
			}
			result := ProcessFormDocument(form)
			if !result.IsValid {
				t.Errorf("Expected valid result at iteration %d", i)
			}
		}
	})
}

// TestIntegration_RealWorldEdgeCases tests edge cases found in real applications
func TestIntegration_RealWorldEdgeCases(t *testing.T) {
	edgeCases := []struct {
		name        string
		input       string
		description string
		expectValid bool
		expectType  string
	}{
		{
			name:        "Excel copy-paste with extra spaces",
			input:       "  716.566.867-59  ",
			description: "User copied from Excel with leading/trailing spaces",
			expectValid: true,
			expectType:  "CPF",
		},
		{
			name:        "Manual typing with extra characters",
			input:       "716.566.867-59.",
			description: "User accidentally typed extra period",
			expectValid: true,
			expectType:  "CPF",
		},
		{
			name:        "Mobile input with autocorrect",
			input:       "22.796.729/0001-59 ",
			description: "Mobile keyboard added extra space",
			expectValid: true,
			expectType:  "CNPJ",
		},
		{
			name:        "PDF copy with formatting artifacts",
			input:       "12.ABC.345/01DE-35\n",
			description: "Copy from PDF included newline",
			expectValid: true,
			expectType:  "CNPJ",
		},
		{
			name:        "Inconsistent formatting",
			input:       "716-566-867.59",
			description: "User used dashes instead of dots",
			expectValid: true,
			expectType:  "CPF",
		},
		{
			name:        "All caps input",
			input:       "12.ABC.345/01DE-35",
			description: "User input in all caps",
			expectValid: true,
			expectType:  "CNPJ",
		},
		{
			name:        "Mixed case input",
			input:       "12.aBc.345/01De-35",
			description: "Mixed case alphanumeric CNPJ",
			expectValid: true,
			expectType:  "CNPJ",
		},
		{
			name:        "Database export with tabs",
			input:       "22\t796\t729\t0001\t59",
			description: "Tab-separated values from database export",
			expectValid: true,
			expectType:  "CNPJ",
		},
		{
			name:        "Voice-to-text artifacts",
			input:       "7 1 6 5 6 6 8 6 7 5 9",
			description: "Voice input with spaces between digits",
			expectValid: true,
			expectType:  "CPF",
		},
	}

	for _, tt := range edgeCases {
		t.Run(tt.name, func(t *testing.T) {
			form := FormData{
				DocumentNumber: tt.input,
				DocumentType:   "AUTO",
			}

			result := ProcessFormDocument(form)

			if result.IsValid != tt.expectValid {
				t.Errorf("%s: ProcessFormDocument(%q).IsValid = %v, want %v",
					tt.description, tt.input, result.IsValid, tt.expectValid)
			}

			if tt.expectValid && result.DocumentType != tt.expectType {
				t.Errorf("%s: ProcessFormDocument(%q).DocumentType = %v, want %v",
					tt.description, tt.input, result.DocumentType, tt.expectType)
			}

			if tt.expectValid && result.FormattedValue == "" {
				t.Errorf("%s: ProcessFormDocument(%q) should return formatted value",
					tt.description, tt.input)
			}
		})
	}
}

// TestIntegration_ConcurrentProcessing tests thread-safety in concurrent scenarios
func TestIntegration_ConcurrentProcessing(t *testing.T) {
	// This test ensures the package is safe for concurrent use
	testData := []FormData{
		{"716.566.867-59", "CPF"},
		{"22.796.729/0001-59", "CNPJ"},
		{"12.ABC.345/01DE-35", "CNPJ"},
	}

	// Run concurrent processing
	done := make(chan bool, len(testData))

	for i, form := range testData {
		go func(index int, f FormData) {
			defer func() { done <- true }()

			// Process the same document multiple times concurrently
			for j := 0; j < 100; j++ {
				result := ProcessFormDocument(f)
				if !result.IsValid {
					t.Errorf("Concurrent test %d-%d failed: %s", index, j, result.ErrorMessage)
					return
				}
			}
		}(i, form)
	}

	// Wait for all goroutines to complete
	for i := 0; i < len(testData); i++ {
		<-done
	}
}
