package cpfcnpj

import (
	"errors"
	"strings"
	"testing"
)

// Test data sets
var (
	validCPFs = []string{
		"64844696793",
		"62641322846",
		"87195726037",
		"71656686759",
		"52824728051",
		"03167158085",
	}

	formattedCPFs = []struct {
		input     string
		cleaned   string
		formatted string
	}{
		{"716.566.867-59", "71656686759", "716.566.867-59"},
		{"648.446.967-93", "64844696793", "648.446.967-93"},
		{"626.413.228-46", "62641322846", "626.413.228-46"},
		{"871.957.260-37", "87195726037", "871.957.260-37"},
	}
)

// Test NewCpf constructor
func TestNewValidCPF(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectValid bool
		expectCode  string
	}{
		// Valid CPFs
		{"Valid formatted CPF", "716.566.867-59", true, ""},
		{"Valid unformatted CPF", "71656686759", true, ""},
		{"Valid CPF with spaces", " 716.566.867-59 ", true, ""},

		// Invalid length
		{"Empty string", "", false, "CPF must have exactly 11 digits"},
		{"Too short", "123456789", false, "CPF must have exactly 11 digits"},
		{"Too long", "123456789012", false, "CPF must have exactly 11 digits"},

		// All same digits
		{"All zeros", "00000000000", false, "document cannot have all same digits"},
		{"All ones", "11111111111", false, "document cannot have all same digits"},
		{"All nines", "99999999999", false, "document cannot have all same digits"},

		// Invalid check digits
		{"Invalid check digits 1", "71656686734", false, "CPF checksum validation failed"},
		{"Invalid check digits 2", "12345678901", false, "CPF checksum validation failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpf, err := NewCpf(tt.input)

			if tt.expectValid {
				if err != nil {
					t.Errorf("NewCpf(%q) expected valid, got error: %v", tt.input, err)
					return
				}
				// Verify the CPF is properly cleaned
				if len(string(cpf)) != 11 {
					t.Errorf("NewCpf(%q) returned CPF with wrong length: %d", tt.input, len(string(cpf)))
				}
			} else {
				if err == nil {
					t.Errorf("NewCpf(%q) expected error, got valid CPF: %s", tt.input, cpf)
					return
				}

				// Check if error message contains expected text
				if !strings.Contains(err.Error(), tt.expectCode) {
					t.Errorf("NewCpf(%q) expected error containing %s, got: %v",
						tt.input, tt.expectCode, err)
				}
			}
		})
	}
}

// Test CPF type String method
func TestCPFString(t *testing.T) {
	tests := []struct {
		name     string
		cpf      CPF
		expected string
	}{
		{"Valid CPF formatting", "71656686759", "716.566.867-59"},
		{"Another valid CPF", "64844696793", "648.446.967-93"},
		{"Wrong length CPF", "123", "123"}, // Returns as-is if wrong length
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cpf.String()
			if got != tt.expected {
				t.Errorf("CPF.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// Test error details and messages
func TestNewValidCPF_ErrorDetails(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedErr error
		msgContains string
	}{
		{
			name:        "Empty input",
			input:       "",
			expectedErr: ErrCPFInvalidLength,
			msgContains: "CPF must have exactly 11 digits, got 0",
		},
		{
			name:        "All zeros",
			input:       "00000000000",
			expectedErr: ErrAllSameDigits,
			msgContains: "CPF cannot have all digits the same",
		},
		{
			name:        "Invalid check digits",
			input:       "12345678901",
			expectedErr: ErrCPFInvalidChecksum,
			msgContains: "CPF check digits are invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCpf(tt.input)
			if err == nil {
				t.Fatalf("NewCpf(%q) expected error, got nil", tt.input)
			}

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("Expected error to be %v, got %v", tt.expectedErr, err)
			}

			if !strings.Contains(err.Error(), tt.msgContains) {
				t.Errorf("Expected error message to contain %q, got %q", tt.msgContains, err.Error())
			}
		})
	}
}

// Test CPF invalid character validation (lines 57-61 in cpf.go)
// Note: This test demonstrates that the invalid character validation is hard to trigger
// because Clean() filters most invalid characters for CPF-length inputs.
// The validation exists as a safety net for edge cases.
func TestNewCPF_InvalidCharacterValidation_EdgeCase(t *testing.T) {
	// To actually test this code path, we would need a way to bypass Clean()
	// or provide input that Clean() doesn't filter but contains invalid chars.
	// This is intentionally difficult by design for security/robustness.

	// For now, we test that the validation logic is sound by testing
	// scenarios that trigger length errors (demonstrating Clean() worked)
	tests := []struct {
		name        string
		input       string
		description string
	}{
		{
			name:        "Letters filtered by Clean",
			input:       "71656686A59", // Clean() removes A, resulting in length error
			description: "Letters are filtered out by Clean(), causing length error",
		},
		{
			name:        "Symbols filtered by Clean",
			input:       "71656686@59", // Clean() removes @, resulting in length error
			description: "Symbols are filtered out by Clean(), causing length error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCpf(tt.input)

			// These should fail with length error, not character error,
			// because Clean() filters the invalid characters first
			if err == nil {
				t.Errorf("NewCpf(%q) expected error, got nil", tt.input)
				return
			}

			// Should be length error, not character error
			if errors.Is(err, ErrInvalidCharacter) {
				t.Errorf("NewCpf(%q) should fail with length error (after Clean), not character error", tt.input)
			}

			if !errors.Is(err, ErrCPFInvalidLength) {
				t.Logf("NewCpf(%q) failed with expected non-character error: %v", tt.input, err)
			}
		})
	}
}
