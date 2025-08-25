package cpfcnpj

import (
	"errors"
	"fmt"
	"strings"
)

// Predefined errors for validation
var (
	ErrAllSameDigits    = errors.New("document cannot have all same digits")
	ErrInvalidCharacter = errors.New("document contains invalid character")

	// CPF-specific errors
	ErrCPFInvalidLength   = errors.New("CPF must have exactly 11 digits")
	ErrCPFInvalidChecksum = errors.New("CPF checksum validation failed")

	// CNPJ-specific errors
	ErrCNPJInvalidLength       = errors.New("CNPJ must have exactly 14 characters")
	ErrCNPJInvalidChecksum     = errors.New("CNPJ checksum validation failed")
	ErrCNPJInvalidAlphanumeric = errors.New("CNPJ alphanumeric format invalid: " +
		"first 12 must be A-Z or 0-9, last 2 must be digits")

	// Security-related errors
	ErrInputTooLarge = errors.New("input string too large: maximum 1000 characters allowed")
)

// Security constants for DoS protection
const (
	// MaxInputSize defines the maximum allowed input string length to prevent DoS attacks.
	// This is generous enough for any valid formatted document (e.g., "12.ABC.345/01DE-35" = 18 chars)
	// while preventing extremely large inputs from consuming excessive resources.
	MaxInputSize = 1000
)

func isAlreadyClean(s string) bool {
	// DoS protection: reject oversized inputs immediately
	if len(s) > MaxInputSize {
		return false
	}

	// Quick length check for valid document lengths
	if len(s) != 11 && len(s) != 14 {
		return false
	}

	isDigitsOnly := len(s) == 11

	// Single pass to check if already clean - optimized character validation
	for i := 0; i < len(s); i++ {
		ch := s[i]

		// Optimize: combine digit check with early return for performance
		if ch >= '0' && ch <= '9' {
			continue // Valid for both CPF and CNPJ
		}

		// Optimize: combine letter check with CPF validation in single condition
		if ch >= 'A' && ch <= 'Z' && !isDigitsOnly {
			continue // Valid for CNPJ alphanumeric
		}

		// Any other character (including lowercase letters or invalid chars for CPF) is invalid
		return false
	}

	return true
}

func cleanString(s string) string {
	// DoS protection: validate input size before processing
	// This prevents potential memory exhaustion attacks from extremely large inputs
	if len(s) > MaxInputSize {
		// Return empty string for oversized inputs instead of panicking
		// This maintains backward compatibility while providing protection
		return ""
	}

	// Use strings.Map for efficient character transformation
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r // Digit
		} else if r >= 'A' && r <= 'Z' {
			return r // Uppercase letter
		} else if r >= 'a' && r <= 'z' {
			return r - 32 // Convert lowercase to uppercase
		}
		// Skip all other characters (formatting, punctuation, etc.)
		return -1
	}, s)
}

func getCharacterValue(char byte) (int, error) {
	if char >= '0' && char <= '9' {
		return int(char - '0'), nil // 0-9
	}
	if char >= 'A' && char <= 'Z' {
		return int(char - 48), nil // ASCII - 48: A=65, so 65-48=17
	}
	return 0, fmt.Errorf("invalid character '%c' (ASCII %d): %w", char, char, ErrInvalidCharacter)
}

func sumDigit(s string, table []int) (int, error) {
	// Process up to the minimum length to handle mismatched inputs gracefully
	length := len(s)
	if len(table) < length {
		length = len(table)
	}

	sum := 0
	for i := 0; i < length; i++ {
		charValue, err := getCharacterValue(s[i])
		if err != nil {
			return 0, fmt.Errorf("error processing character at position %d: %w", i, err)
		}
		sum += table[i] * charValue
	}
	return sum, nil
}

func calculateModule11Digits(base string, firstTable, secondTable []int) (firstDigit, secondDigit int, err error) {
	// Calculate first verification digit
	sum1, err := sumDigit(base, firstTable)
	if err != nil {
		return 0, 0, fmt.Errorf("error calculating first digit: %w", err)
	}
	remainder1 := sum1 % 11
	digit1 := 0
	if remainder1 >= 2 {
		digit1 = 11 - remainder1
	}

	// Calculate second verification digit
	sum2, err := sumDigit(base, secondTable[:len(base)]) // Only use weights for base length
	if err != nil {
		return 0, 0, fmt.Errorf("error calculating second digit: %w", err)
	}
	if len(secondTable) > len(base) {
		// Add the contribution of the first digit with its corresponding weight
		sum2 += digit1 * secondTable[len(base)]
	}

	remainder2 := sum2 % 11
	digit2 := 0
	if remainder2 >= 2 {
		digit2 = 11 - remainder2
	}

	return digit1, digit2, nil
}

func isSameCharacter(s string) bool {
	if len(s) <= 1 {
		return false
	}

	// Check if all characters are the same as the first character
	firstChar := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] != firstChar {
			return false
		}
	}
	return true
}

// Clean removes formatting and normalizes CPF/CNPJ documents.
func Clean(s string) string {
	if s == "" {
		return s
	}

	// Fast path: check if already clean
	if isAlreadyClean(s) {
		return s // Fast path: already clean, 0 allocations
	}

	// Perform full cleaning
	result := cleanString(s)

	// Apply document-type specific filtering based on detected length
	if len(result) == 11 {
		// CPF: should contain only digits, filter out any letters
		return filterToDigitsOnly(result)
	}

	// For CNPJ (14 chars) or invalid lengths, return as-is
	return result
}

func filterToDigitsOnly(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1 // Remove non-digit characters
	}, s)
}

func formatDocument(s, pattern string) string {
	var result strings.Builder
	result.Grow(len(pattern))

	pos := 0
	for i := 0; i < len(pattern); i++ {
		if pattern[i] == 'X' {
			if pos < len(s) {
				result.WriteByte(s[pos])
				pos++
			}
		} else {
			result.WriteByte(pattern[i])
		}
	}

	return result.String()
}
