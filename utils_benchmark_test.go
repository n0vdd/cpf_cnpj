package cpfcnpj

import (
	"testing"
)

// Realistic test data for benchmarking core functions
var (
	// CPF test data
	cpfClean     = "71656686759"
	cpfFormatted = "716.566.867-59"
	cpfDirty     = " 7!1@6#5$6%6^8&6*7(-5)9 "
	cpfInvalid   = "11111111111"

	// CNPJ numeric test data
	cnpjNumericClean     = "22796729000159"
	cnpjNumericFormatted = "22.796.729/0001-59"
	cnpjNumericDirty     = " 2!2@7#9$6%7^2&9*0(0)0{1}5[9] "
	cnpjNumericInvalid   = "00000000000000"

	// CNPJ alphanumeric test data
	cnpjAlphaClean     = "12ABC34501DE35"
	cnpjAlphaFormatted = "12.ABC.345/01DE-35"
	cnpjAlphaLowercase = "12.abc.345/01de-35"
	cnpjAlphaInvalid   = "12ABC34501DE99"
)

// BenchmarkClean tests the Clean function with realistic input patterns
func BenchmarkClean(b *testing.B) {
	b.Run("CPF_Clean", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			Clean(cpfClean)
		}
	})

	b.Run("CPF_Formatted", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			Clean(cpfFormatted)
		}
	})

	b.Run("CPF_Dirty", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			Clean(cpfDirty)
		}
	})

	b.Run("CNPJ_Numeric_Clean", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			Clean(cnpjNumericClean)
		}
	})

	b.Run("CNPJ_Numeric_Formatted", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			Clean(cnpjNumericFormatted)
		}
	})

	b.Run("CNPJ_Alpha_Clean", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			Clean(cnpjAlphaClean)
		}
	})

	b.Run("CNPJ_Alpha_Formatted", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			Clean(cnpjAlphaFormatted)
		}
	})

	b.Run("CNPJ_Alpha_Lowercase", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			Clean(cnpjAlphaLowercase)
		}
	})
}

// BenchmarkNewCpf tests CPF validation performance
func BenchmarkNewCpf(b *testing.B) {
	b.Run("Valid_Clean", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = NewCpf(cpfClean)
		}
	})

	b.Run("Valid_Formatted", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = NewCpf(cpfFormatted)
		}
	})

	b.Run("Valid_Dirty", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = NewCpf(cpfDirty)
		}
	})

	b.Run("Invalid_SameDigits", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = NewCpf(cpfInvalid)
		}
	})

	b.Run("Invalid_WrongLength", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = NewCpf("123")
		}
	})
}

// BenchmarkNewCnpj tests CNPJ validation performance for both numeric and alphanumeric formats
func BenchmarkNewCnpj(b *testing.B) {
	b.Run("Numeric_Valid_Clean", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = NewCnpj(cnpjNumericClean)
		}
	})

	b.Run("Numeric_Valid_Formatted", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = NewCnpj(cnpjNumericFormatted)
		}
	})

	b.Run("Numeric_Valid_Dirty", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = NewCnpj(cnpjNumericDirty)
		}
	})

	b.Run("Alpha_Valid_Clean", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = NewCnpj(cnpjAlphaClean)
		}
	})

	b.Run("Alpha_Valid_Formatted", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = NewCnpj(cnpjAlphaFormatted)
		}
	})

	b.Run("Alpha_Valid_Lowercase", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = NewCnpj(cnpjAlphaLowercase)
		}
	})

	b.Run("Numeric_Invalid_SameDigits", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = NewCnpj(cnpjNumericInvalid)
		}
	})

	b.Run("Alpha_Invalid_Checksum", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = NewCnpj(cnpjAlphaInvalid)
		}
	})

	b.Run("Invalid_WrongLength", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = NewCnpj("123")
		}
	})
}

// BenchmarkStringMethods tests the performance of String() methods for formatting
func BenchmarkStringMethods(b *testing.B) {
	// Create valid instances for benchmarking String() methods
	validCPF, _ := NewCpf(cpfClean)
	validCNPJNumeric, _ := NewCnpj(cnpjNumericClean)
	validCNPJAlpha, _ := NewCnpj(cnpjAlphaClean)

	b.Run("CPF_String", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = validCPF.String()
		}
	})

	b.Run("CNPJ_Numeric_String", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = validCNPJNumeric.String()
		}
	})

	b.Run("CNPJ_Alpha_String", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = validCNPJAlpha.String()
		}
	})

	// Test String() with invalid length (edge case performance)
	invalidCPF := CPF("123")     // Wrong length
	invalidCNPJ := CNPJ("12345") // Wrong length

	b.Run("CPF_String_InvalidLength", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = invalidCPF.String()
		}
	})

	b.Run("CNPJ_String_InvalidLength", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = invalidCNPJ.String()
		}
	})
}

// BenchmarkErrorPaths tests performance of error handling paths
func BenchmarkErrorPaths(b *testing.B) {
	// Test error path performance for various validation failures

	b.Run("CPF_ErrorPath_InvalidLength", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = NewCpf("123") // Too short
		}
	})

	b.Run("CPF_ErrorPath_InvalidCharacters", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = NewCpf("1234567890A") // Contains letter
		}
	})

	b.Run("CPF_ErrorPath_SameDigits", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = NewCpf("11111111111") // All same digits
		}
	})

	b.Run("CPF_ErrorPath_InvalidChecksum", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = NewCpf("12345678901") // Invalid checksum
		}
	})

	b.Run("CNPJ_ErrorPath_InvalidLength", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = NewCnpj("123") // Too short
		}
	})

	b.Run("CNPJ_ErrorPath_InvalidFormat", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = NewCnpj("12ABC34501DEAB") // Letters in check digits
		}
	})

	b.Run("CNPJ_ErrorPath_SameDigits", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = NewCnpj("11111111111111") // All same digits
		}
	})

	b.Run("CNPJ_ErrorPath_InvalidChecksum", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = NewCnpj("12ABC34501DE99") // Invalid checksum
		}
	})
}

// BenchmarkUtilityFunctions tests performance of internal utility functions
func BenchmarkUtilityFunctions(b *testing.B) {
	b.Run("FilterToDigitsOnly", func(b *testing.B) {
		b.ReportAllocs()
		input := "12ABC34501DE35" // Mixed alphanumeric
		for i := 0; i < b.N; i++ {
			_ = filterToDigitsOnly(input)
		}
	})

	b.Run("FilterToDigitsOnly_EmptyResult", func(b *testing.B) {
		b.ReportAllocs()
		input := "ABCDEFGHIJKLMN" // Only letters
		for i := 0; i < b.N; i++ {
			_ = filterToDigitsOnly(input)
		}
	})

	b.Run("FilterToDigitsOnly_AllDigits", func(b *testing.B) {
		b.ReportAllocs()
		input := "1234567890123" // Only digits
		for i := 0; i < b.N; i++ {
			_ = filterToDigitsOnly(input)
		}
	})

	b.Run("IsSameCharacter_True", func(b *testing.B) {
		b.ReportAllocs()
		input := "11111111111111" // All same
		for i := 0; i < b.N; i++ {
			_ = isSameCharacter(input)
		}
	})

	b.Run("IsSameCharacter_False", func(b *testing.B) {
		b.ReportAllocs()
		input := "12345678901234" // Different chars
		for i := 0; i < b.N; i++ {
			_ = isSameCharacter(input)
		}
	})

	b.Run("GetCharacterValue_Digit", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = getCharacterValue('5')
		}
	})

	b.Run("GetCharacterValue_Letter", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = getCharacterValue('A')
		}
	})

	b.Run("GetCharacterValue_Invalid", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = getCharacterValue('@')
		}
	})

	b.Run("IsValidCNPJFormat_Valid", func(b *testing.B) {
		b.ReportAllocs()
		input := "12ABC34501DE35" // Valid format
		for i := 0; i < b.N; i++ {
			_ = isValidCNPJFormat(input)
		}
	})

	b.Run("IsValidCNPJFormat_Invalid", func(b *testing.B) {
		b.ReportAllocs()
		input := "12ABC34501DEAB" // Invalid format
		for i := 0; i < b.N; i++ {
			_ = isValidCNPJFormat(input)
		}
	})
}

// BenchmarkModule11Algorithm tests the core Module 11 calculation performance
func BenchmarkModule11Algorithm(b *testing.B) {
	cpfBase := "716566867"
	cnpjBase := "12ABC34501DE"

	b.Run("CPF_Module11", func(b *testing.B) {
		b.ReportAllocs()
		firstTable := []int{10, 9, 8, 7, 6, 5, 4, 3, 2}
		secondTable := []int{11, 10, 9, 8, 7, 6, 5, 4, 3, 2}
		for i := 0; i < b.N; i++ {
			_, _, _ = calculateModule11Digits(cpfBase, firstTable, secondTable)
		}
	})

	b.Run("CNPJ_Module11", func(b *testing.B) {
		b.ReportAllocs()
		firstTable := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
		secondTable := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
		for i := 0; i < b.N; i++ {
			_, _, _ = calculateModule11Digits(cnpjBase, firstTable, secondTable)
		}
	})

	b.Run("SumDigit_CPF", func(b *testing.B) {
		b.ReportAllocs()
		table := []int{10, 9, 8, 7, 6, 5, 4, 3, 2}
		for i := 0; i < b.N; i++ {
			_, _ = sumDigit(cpfBase, table)
		}
	})

	b.Run("SumDigit_CNPJ_Alphanumeric", func(b *testing.B) {
		b.ReportAllocs()
		table := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
		for i := 0; i < b.N; i++ {
			_, _ = sumDigit(cnpjBase, table)
		}
	})
}

// BenchmarkMemoryAllocations tests memory allocation patterns
func BenchmarkMemoryAllocations(b *testing.B) {
	// Test Clean function allocations under different scenarios
	b.Run("Clean_FastPath_CPF", func(b *testing.B) {
		b.ReportAllocs()
		input := cpfClean // Already clean, should use fast path
		for i := 0; i < b.N; i++ {
			_ = Clean(input)
		}
	})

	b.Run("Clean_FastPath_CNPJ", func(b *testing.B) {
		b.ReportAllocs()
		input := cnpjAlphaClean // Already clean, should use fast path
		for i := 0; i < b.N; i++ {
			_ = Clean(input)
		}
	})

	b.Run("Clean_SlowPath_CPF", func(b *testing.B) {
		b.ReportAllocs()
		input := cpfFormatted // Needs cleaning
		for i := 0; i < b.N; i++ {
			_ = Clean(input)
		}
	})

	b.Run("Clean_SlowPath_CNPJ", func(b *testing.B) {
		b.ReportAllocs()
		input := cnpjAlphaFormatted // Needs cleaning
		for i := 0; i < b.N; i++ {
			_ = Clean(input)
		}
	})

	// Test String formatting allocations
	validCPF, _ := NewCpf(cpfClean)
	validCNPJ, _ := NewCnpj(cnpjAlphaClean)

	b.Run("String_CPF_Allocations", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = validCPF.String()
		}
	})

	b.Run("String_CNPJ_Allocations", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = validCNPJ.String()
		}
	})
}

// BenchmarkRealWorldWorkloads simulates realistic usage patterns
func BenchmarkRealWorldWorkloads(b *testing.B) {
	// Simulate typical web form validation workflow
	b.Run("WebForm_Workflow_CPF", func(b *testing.B) {
		b.ReportAllocs()
		input := "716.566.867-59" // Typical user input
		for i := 0; i < b.N; i++ {
			cleaned := Clean(input)
			cpf, err := NewCpf(cleaned)
			if err == nil {
				_ = cpf.String()
			}
		}
	})

	b.Run("WebForm_Workflow_CNPJ", func(b *testing.B) {
		b.ReportAllocs()
		input := "12.ABC.345/01DE-35" // Typical user input
		for i := 0; i < b.N; i++ {
			cleaned := Clean(input)
			cnpj, err := NewCnpj(cleaned)
			if err == nil {
				_ = cnpj.String()
			}
		}
	})

	// Simulate batch processing workflow
	inputs := []string{
		"716.566.867-59",
		"22.796.729/0001-59",
		"12.ABC.345/01DE-35",
		"invalid-input",
	}

	b.Run("Batch_Processing", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, input := range inputs {
				cleaned := Clean(input)
				switch len(cleaned) {
				case 11:
					_, _ = NewCpf(cleaned)
				case 14:
					_, _ = NewCnpj(cleaned)
				}
			}
		}
	})
}
