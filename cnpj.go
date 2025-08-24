package cpfcnpj

import (
	"fmt"
	"strconv"
	"strings"
)

// Constants for CNPJ validation
const (
	CNPJLength = 14
)

// CNPJ validation tables for Module 11 algorithm
var (
	// Keep lowercase for internal use
	cnpjFirstDigitTable  = []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	cnpjSecondDigitTable = []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
)

// CNPJ represents a Brazilian tax identification number.
// Supports both numeric and alphanumeric formats.
type CNPJ string

// NewCnpj creates and validates a CNPJ from a string.
// Supports both numeric and alphanumeric formats.
// Returns error if CNPJ is invalid.
func NewCnpj(s string) (CNPJ, error) {
	// Clean input: keep alphanumeric chars, normalize case
	cleaned := Clean(s)

	// Validate length
	if len(cleaned) != CNPJLength {
		return "", fmt.Errorf("CNPJ must have exactly %d characters, got %d: %w", CNPJLength, len(cleaned),
			ErrCNPJInvalidLength)
	}

	// Validate character format
	if !isValidCNPJFormat(cleaned) {
		return "", fmt.Errorf("CNPJ format is invalid: first 12 characters must be A-Z or 0-9, "+
			"last 2 must be 0-9: %w", ErrCNPJInvalidAlphanumeric)
	}

	// Reject invalid patterns (all same characters)
	if isSameCharacter(cleaned) {
		return "", fmt.Errorf("CNPJ cannot have all same characters: %w", ErrAllSameDigits)
	}

	// Validate check digits using Module 11 algorithm
	firstPart := cleaned[:12]
	d1, d2, err := calculateModule11Digits(firstPart, cnpjFirstDigitTable, cnpjSecondDigitTable)
	if err != nil {
		return "", fmt.Errorf("error calculating CNPJ check digits: %w", err)
	}

	expectedCNPJ := firstPart + strconv.Itoa(d1) + strconv.Itoa(d2)
	if expectedCNPJ != cleaned {
		return "", fmt.Errorf("CNPJ check digits are invalid: %w", ErrCNPJInvalidChecksum)
	}

	return CNPJ(cleaned), nil
}

// String returns the CNPJ formatted as XX.XXX.XXX/XXXX-XX.
func (c *CNPJ) String() string {
	str := string(*c)

	// Safety check: only format if exactly 14 characters
	if len(str) != CNPJLength {
		return str
	}

	// Use shared formatting function
	return formatDocument(str, "XX.XXX.XXX/XXXX-XX")
}

// isValidCNPJFormat validates the character format of CNPJ
func isValidCNPJFormat(cnpj string) bool {
	if len(cnpj) != CNPJLength {
		return false
	}

	// First 12 characters must be alphanumeric (A-Z, 0-9)
	firstTwelve := cnpj[:12]
	if invalidIndex := strings.IndexFunc(firstTwelve, func(r rune) bool {
		return !((r >= '0' && r <= '9') || (r >= 'A' && r <= 'Z'))
	}); invalidIndex != -1 {
		return false
	}

	// Last 2 characters must be numeric (check digits)
	lastTwo := cnpj[12:]
	if invalidIndex := strings.IndexFunc(lastTwo, func(r rune) bool {
		return !(r >= '0' && r <= '9')
	}); invalidIndex != -1 {
		return false
	}

	return true
}
