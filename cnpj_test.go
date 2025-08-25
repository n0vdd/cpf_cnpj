package cpfcnpj

import (
	"errors"
	"strings"
	"testing"
)

// Test NewCnpj constructor with valid inputs
func TestNewCNPJ_Valid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Valid numeric CNPJ", "22796729000159"},
		{"Valid formatted numeric CNPJ", "22.796.729/0001-59"},
		{"Valid alphanumeric CNPJ", "12ABC34501DE35"},
		{"Valid formatted alphanumeric CNPJ", "12.ABC.345/01DE-35"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cnpj, err := NewCnpj(tt.input)
			if err != nil {
				t.Errorf("NewCnpj(%q) expected valid, got error: %v", tt.input, err)
				return
			}
			if len(string(cnpj)) != 14 {
				t.Errorf("NewCnpj(%q) returned CNPJ with wrong length: %d", tt.input, len(string(cnpj)))
			}
		})
	}
}

// Test NewCnpj constructor with invalid inputs
func TestNewCNPJ_Invalid(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expectCode string
	}{
		{"Empty string", "", "CNPJ must have exactly 14 characters"},
		{"Too short", "12345", "CNPJ must have exactly 14 characters"},
		{"Too long", "123456789012345", "CNPJ must have exactly 14 characters"},
		{"All zeros", "00000000000000", "CNPJ cannot have all same characters"},
		{"All ones", "11111111111111", "CNPJ cannot have all same characters"},
		{"All nines", "99999999999999", "CNPJ cannot have all same characters"},
		{"All twos", "22222222222222", "CNPJ cannot have all same characters"},
		{"All eights", "88888888888888", "CNPJ cannot have all same characters"},
		{"Invalid check digits", "12ABC34501DE99", "CNPJ check digits are invalid"},
		{"Invalid check digits not numeric", "12ABC34501DEAB", "CNPJ format is invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCnpj(tt.input)
			if err == nil {
				t.Errorf("NewCnpj(%q) expected error, got nil", tt.input)
				return
			}
			if !strings.Contains(err.Error(), tt.expectCode) {
				t.Errorf("NewCnpj(%q) expected error containing %q, got: %v", tt.input, tt.expectCode, err)
			}
		})
	}
}

// Test CNPJ type String method
func TestCNPJString(t *testing.T) {
	tests := []struct {
		name     string
		cnpj     CNPJ
		expected string
	}{
		{"Valid numeric CNPJ", "22796729000159", "22.796.729/0001-59"},
		{"Valid alphanumeric CNPJ", "12ABC34501DE35", "12.ABC.345/01DE-35"},
		{"Wrong length CNPJ", "123", "123"}, // Returns as-is if wrong length
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cnpj.String()
			if got != tt.expected {
				t.Errorf("CNPJ.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// Test error details and wrapping
func TestNewCNPJ_ErrorDetails(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedErr error
		msgContains string
	}{
		{
			name:        "Empty input",
			input:       "",
			expectedErr: ErrCNPJInvalidLength,
			msgContains: "CNPJ must have exactly 14 characters, got 0",
		},
		{
			name:        "All zeros",
			input:       "00000000000000",
			expectedErr: ErrAllSameDigits,
			msgContains: "CNPJ cannot have all same characters",
		},
		{
			name:        "Invalid check digits",
			input:       "12ABC34501DE99",
			expectedErr: ErrCNPJInvalidChecksum,
			msgContains: "CNPJ check digits are invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCnpj(tt.input)
			if err == nil {
				t.Fatalf("NewCnpj(%q) expected error, got nil", tt.input)
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

// Test alphanumeric character value conversion
func TestGetCharacterValue(t *testing.T) {
	tests := []struct {
		name    string
		char    byte
		want    int
		wantErr bool
	}{
		{"Digit 0", '0', 0, false},
		{"Digit 9", '9', 9, false},
		{"Letter A", 'A', 17, false}, // 65 - 48 = 17
		{"Letter Z", 'Z', 42, false}, // 90 - 48 = 42
		{"Invalid character", '@', 0, true},
		{"Lowercase letter", 'a', 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getCharacterValue(tt.char)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCharacterValue(%v) error = %v, wantErr %v", tt.char, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getCharacterValue(%v) = %v, want %v", tt.char, got, tt.want)
			}
		})
	}
}

// Test isValidCNPJFormat function comprehensively - targeting edge cases for 77.8% coverage
func TestIsValidCNPJFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Valid formats
		{"Valid numeric CNPJ", "22796729000159", true},
		{"Valid alphanumeric CNPJ", "12ABC34501DE35", true},
		{"Valid mixed alphanumeric", "1A2B3C4D5E0F12", true},
		{"Valid all letters first 12", "ABCDEFGHIJKL12", true},
		{"Valid all numbers", "12345678901234", true},

		// Invalid length (line 126-128)
		{"Too short", "123", false},
		{"Too long", "123456789012345", false},
		{"Empty string", "", false},
		{"Length 13", "1234567890123", false},
		{"Length 15", "123456789012345", false},

		// Invalid characters in first 12 positions (lines 131-135)
		{"Invalid char @ in position 0", "@2345678901234", false},
		{"Invalid char @ in position 1", "1@345678901234", false},
		{"Invalid char @ in position 5", "12345@78901234", false},
		{"Invalid char @ in position 11", "12345678901@34", false},
		{"Lowercase letter in first 12", "12abc678901234", false},
		{"Symbol # in first 12", "12#45678901234", false},
		{"Symbol $ in first 12", "12$45678901234", false},
		{"Symbol % in first 12", "12%45678901234", false},
		{"Symbol & in first 12", "12&45678901234", false},
		{"Symbol * in first 12", "12*45678901234", false},
		{"Space in first 12", "12 45678901234", false},
		{"Tab in first 12", "12\t45678901234", false},
		{"Newline in first 12", "12\n45678901234", false},

		// Invalid characters in check digit positions (lines 138-142)
		{"Letter A in check digit pos 12", "12345678901AA4", false},
		{"Letter A in check digit pos 13", "123456789012A4", false},
		{"Letter B in check digit pos 13", "123456789012B4", false},
		{"Letter Z in both check positions", "123456789012ZZ", false},
		{"Symbol @ in check digit pos 12", "12345678901@34", false},
		{"Symbol @ in check digit pos 13", "123456789012@4", false},
		{"Lowercase in check digit pos", "12345678901a34", false},
		{"Space in check digit position", "123456789012 4", false},
		{"Tab in check digit position", "123456789012\t4", false},

		// Boundary cases for check digits - ASCII boundary testing
		{"ASCII 47 (/) in check digit", "123456789012/4", false}, // Just before '0'
		{"ASCII 58 (:) in check digit", "123456789012:4", false}, // Just after '9'
		{"Valid ASCII 48 (0) in check", "12345678901203", true},  // '0' - 14 chars total
		{"Valid ASCII 57 (9) in check", "12345678901293", true},  // '9' - 14 chars total

		// Mixed valid/invalid scenarios
		{"Valid alphanumeric + invalid check", "12ABC34501DE@@", false},
		{"Invalid first part + valid check", "12@BC34501DE35", false},
		{"All valid except one char", "12ABC345@1DE35", false},

		// Unicode and special character edge cases
		{"Unicode digit lookalike", "1234567890１２34", false},  // Full-width digits in check positions
		{"High Unicode in first 12", "1234567890１234", false}, // Full-width digit in first 12
		{"Null character", "12345678901\x00234", false},
		{"Extended ASCII", "123456789012ä4", false}, // Extended ASCII in check digit

		// Real-world formatting remnants that might slip through Clean()
		{"Dot remnant", "12345678901.34", false},
		{"Dash remnant", "12345678901-34", false},
		{"Slash remnant", "12345678901/34", false},

		// Edge cases with boundary alphanumeric characters
		{"Boundary letter-number", "ZZZZZZZZZZZ999", true},  // All Z's (valid) + digits
		{"Boundary number-letter", "000000000000AA", false}, // Digits + letters in check positions (invalid)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidCNPJFormat(tt.input)
			if got != tt.expected {
				t.Errorf("isValidCNPJFormat(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

// Test specific getCharacterValue error cases that affect isValidCNPJFormat
func TestIsValidCNPJFormat_GetCharacterValueErrors(t *testing.T) {
	// These tests specifically target the error handling in getCharacterValue
	// which is called from isValidCNPJFormat (line 132)
	errorTests := []struct {
		name  string
		input string
		pos   int // Position where the invalid character appears
	}{
		{"Invalid char position 0", "@23456789012AB", 0},
		{"Invalid char position 1", "1@3456789012AB", 1},
		{"Invalid char position 2", "12@456789012AB", 2},
		{"Invalid char position 3", "123@56789012AB", 3},
		{"Invalid char position 4", "1234@6789012AB", 4},
		{"Invalid char position 5", "12345@789012AB", 5},
		{"Invalid char position 6", "123456@89012AB", 6},
		{"Invalid char position 7", "1234567@9012AB", 7},
		{"Invalid char position 8", "12345678@012AB", 8},
		{"Invalid char position 9", "123456789@12AB", 9},
		{"Invalid char position 10", "1234567890@2AB", 10},
		{"Invalid char position 11", "12345678901@AB", 11},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			// All of these should return false because getCharacterValue will error
			result := isValidCNPJFormat(tt.input)
			if result != false {
				t.Errorf("isValidCNPJFormat(%q) = %v, want false (due to invalid char at pos %d)",
					tt.input, result, tt.pos)
			}
		})
	}
}

// Test comprehensive boundary conditions for CNPJ format validation
func TestIsValidCNPJFormat_BoundaryConditions(t *testing.T) {
	tests := []struct {
		name        string
		description string
		input       string
		expected    bool
	}{
		{
			name:        "All A's first 12, digits last 2",
			description: "Test maximum letter values with valid check digits",
			input:       "AAAAAAAAAAAA99",
			expected:    true,
		},
		{
			name:        "All Z's first 12, digits last 2",
			description: "Test maximum letter boundary with valid check digits",
			input:       "ZZZZZZZZZZZZ99",
			expected:    true,
		},
		{
			name:        "All 0's first 12, 00 last 2",
			description: "Test minimum digit values throughout",
			input:       "00000000000000",
			expected:    true,
		},
		{
			name:        "All 9's first 12, 99 last 2",
			description: "Test maximum digit values throughout",
			input:       "99999999999999",
			expected:    true,
		},
		{
			name:        "ASCII 47 in first 12",
			description: "Test character just before '0' in ASCII",
			input:       "123456789012//",
			expected:    false,
		},
		{
			name:        "ASCII 58 in first 12",
			description: "Test character just after '9' in ASCII",
			input:       "123456789012::",
			expected:    false,
		},
		{
			name:        "ASCII 64 in first 12",
			description: "Test character just before 'A' in ASCII",
			input:       "123456789012@@",
			expected:    false,
		},
		{
			name:        "ASCII 91 in first 12",
			description: "Test character just after 'Z' in ASCII",
			input:       "123456789012[[",
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidCNPJFormat(tt.input)
			if result != tt.expected {
				t.Errorf("%s: isValidCNPJFormat(%q) = %v, want %v",
					tt.description, tt.input, result, tt.expected)
			}
		})
	}
}

// Test CNPJ Raw method with numeric CNPJs
func TestCNPJRaw_Numeric(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Valid numeric CNPJ formatted", "22.796.729/0001-59", "22796729000159"},
		{"Valid numeric CNPJ unformatted", "22796729000159", "22796729000159"},
		{"Valid numeric CNPJ with spaces", " 22.796.729/0001-59 ", "22796729000159"},
		{"Another valid numeric CNPJ", "11222333000181", "11222333000181"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cnpj, err := NewCnpj(tt.input)
			if err != nil {
				t.Errorf("NewCnpj(%q) unexpected error: %v", tt.input, err)
				return
			}

			got := cnpj.Raw()
			if got != tt.expected {
				t.Errorf("CNPJ.Raw() = %q, want %q", got, tt.expected)
			}

			// Verify Raw returns exactly 14 characters
			if len(got) != CNPJLength {
				t.Errorf("CNPJ.Raw() returned wrong length: got %d, want %d", len(got), CNPJLength)
			}

			// Verify all characters are digits for numeric CNPJ
			for i, char := range got {
				if char < '0' || char > '9' {
					t.Errorf("CNPJ.Raw()[%d] = %q, expected digit", i, char)
				}
			}
		})
	}
}

// Test CNPJ Raw method with alphanumeric CNPJs
func TestCNPJRaw_Alphanumeric(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Valid alphanumeric CNPJ formatted", "12.ABC.345/01DE-35", "12ABC34501DE35"},
		{"Valid alphanumeric CNPJ unformatted", "12ABC34501DE35", "12ABC34501DE35"},
		{"Valid alphanumeric with spaces", " 12.ABC.345/01DE-35 ", "12ABC34501DE35"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cnpj, err := NewCnpj(tt.input)
			if err != nil {
				t.Errorf("NewCnpj(%q) unexpected error: %v", tt.input, err)
				return
			}

			got := cnpj.Raw()
			if got != tt.expected {
				t.Errorf("CNPJ.Raw() = %q, want %q", got, tt.expected)
			}

			// Verify Raw returns exactly 14 characters
			if len(got) != CNPJLength {
				t.Errorf("CNPJ.Raw() returned wrong length: got %d, want %d", len(got), CNPJLength)
			}

			// Verify first 12 characters are alphanumeric, last 2 are digits
			for i, char := range got[:12] {
				if !((char >= '0' && char <= '9') || (char >= 'A' && char <= 'Z')) {
					t.Errorf("CNPJ.Raw()[%d] = %q, expected alphanumeric character", i, char)
				}
			}

			// Verify last 2 characters are digits (check digits)
			for i, char := range got[12:] {
				if char < '0' || char > '9' {
					t.Errorf("CNPJ.Raw()[%d+12] = %q, expected digit", i, char)
				}
			}
		})
	}
}

// Test CNPJ Raw vs String consistency
func TestCNPJRawStringConsistency(t *testing.T) {
	testCNPJs := []struct {
		name     string
		input    string
		cleanRaw string
	}{
		{"Numeric CNPJ", "22.796.729/0001-59", "22796729000159"},
		{"Alphanumeric CNPJ", "12.ABC.345/01DE-35", "12ABC34501DE35"},
		{"Mixed case cleaned", "12.abc.345/01de-35", "12ABC34501DE35"}, // Clean() normalizes case
	}

	for _, tt := range testCNPJs {
		t.Run(tt.name, func(t *testing.T) {
			cnpj, err := NewCnpj(tt.input)
			if err != nil {
				t.Errorf("NewCnpj(%q) unexpected error: %v", tt.input, err)
				return
			}

			raw := cnpj.Raw()
			formatted := cnpj.String()

			// Raw should be the cleaned characters
			if raw != tt.cleanRaw {
				t.Errorf("CNPJ.Raw() = %q, expected %q", raw, tt.cleanRaw)
			}

			// String should be formatted version
			expectedFormatted := formatDocument(tt.cleanRaw, "XX.XXX.XXX/XXXX-XX")
			if formatted != expectedFormatted {
				t.Errorf("CNPJ.String() = %q, expected %q", formatted, expectedFormatted)
			}

			// Cleaning formatted should give raw
			cleaned := Clean(formatted)
			if cleaned != raw {
				t.Errorf("Clean(CNPJ.String()) = %q, expected CNPJ.Raw() = %q", cleaned, raw)
			}
		})
	}
}

// Test CNPJ Raw method edge cases
func TestCNPJRaw_EdgeCases(t *testing.T) {
	// Note: We can't test with actually invalid CNPJs since NewCnpj validates them
	// These tests focus on valid CNPJs that exercise different character ranges

	validTestCNPJs := []string{
		"11444777000161", // Valid numeric CNPJ
	}

	for _, testCNPJ := range validTestCNPJs {
		t.Run("Valid_"+testCNPJ, func(t *testing.T) {
			cnpj, err := NewCnpj(testCNPJ)
			if err != nil {
				t.Errorf("NewCnpj(%q) unexpected error: %v", testCNPJ, err)
				return
			}

			raw := cnpj.Raw()

			// Raw should equal the clean input for valid CNPJs
			cleaned := Clean(testCNPJ)
			if raw != cleaned {
				t.Errorf("CNPJ.Raw() = %q, expected Clean(%q) = %q", raw, testCNPJ, cleaned)
			}

			// Raw should be exactly 14 characters
			if len(raw) != CNPJLength {
				t.Errorf("CNPJ.Raw() length = %d, expected %d", len(raw), CNPJLength)
			}
		})
	}
}

// Benchmark CNPJ Raw method for performance verification
func BenchmarkCNPJRaw(b *testing.B) {
	cnpj, err := NewCnpj("22796729000159")
	if err != nil {
		b.Fatal("Failed to create CNPJ for benchmark:", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cnpj.Raw()
	}
}
