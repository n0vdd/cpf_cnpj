// Package cpfcnpj provides validation and formatting for Brazilian taxpayer identification documents (CPF and CNPJ).
//
// Basic usage:
//
//	cpf, err := cpfcnpj.NewCpf("716.566.867-59")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(cpf.String()) // "716.566.867-59"
//
// The package supports both numeric and alphanumeric CNPJ formats.
package cpfcnpj

import (
	"fmt"
	"strconv"
)

// Constants for CPF validation
const (
	CPFLength = 11
)

// CPF validation tables for Module 11 algorithm
var (
	cpfFirstDigitTable  = []int{10, 9, 8, 7, 6, 5, 4, 3, 2}
	cpfSecondDigitTable = []int{11, 10, 9, 8, 7, 6, 5, 4, 3, 2}
)

// CPF type definition
type CPF string

// NewCpf creates and validates a CPF from a string.
// Returns error if CPF is invalid.
func NewCpf(s string) (CPF, error) {
	// Clean input: keep only digits
	cleaned := Clean(s)

	// Validate length
	if len(cleaned) != CPFLength {
		return "", fmt.Errorf("CPF must have exactly %d digits, got %d: %w", CPFLength, len(cleaned),
			ErrCPFInvalidLength)
	}

	// Reject invalid patterns (all same digits)
	if isSameCharacter(cleaned) {
		return "", fmt.Errorf("CPF cannot have all digits the same: %w", ErrAllSameDigits)
	}

	// Validate check digits using Module 11 algorithm
	firstPart := cleaned[:9]
	d1, d2, err := calculateModule11Digits(firstPart, cpfFirstDigitTable, cpfSecondDigitTable)
	if err != nil {
		return "", fmt.Errorf("error calculating CPF check digits: %w", err)
	}

	expectedCPF := firstPart + strconv.Itoa(d1) + strconv.Itoa(d2)
	if expectedCPF != cleaned {
		return "", fmt.Errorf("CPF check digits are invalid: %w", ErrCPFInvalidChecksum)
	}

	return CPF(cleaned), nil
}

// String returns the CPF formatted as XXX.XXX.XXX-XX.
func (c *CPF) String() string {
	str := string(*c)

	// Safety check: only format if exactly 11 digits
	if len(str) != CPFLength {
		return str
	}

	// Use shared formatting function
	return formatDocument(str, "XXX.XXX.XXX-XX")
}
