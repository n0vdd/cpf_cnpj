package cpfcnpj

import (
	"testing"
)

// Test isAlreadyClean function
func TestIsAlreadyClean(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Fast path success cases - CPF
		{"CPF clean digits", "71656686759", true},
		{"CPF all zeros", "00000000000", true},
		{"CPF all nines", "99999999999", true},
		{"CPF mixed digits", "12345678901", true},
		{"CPF valid checksum", "11144477735", true},

		// Fast path success cases - CNPJ numeric
		{"CNPJ numeric clean", "22796729000159", true},
		{"CNPJ numeric all zeros", "00000000000000", true},
		{"CNPJ numeric all nines", "99999999999999", true},
		{"CNPJ numeric mixed", "12345678901234", true},

		// Fast path success cases - CNPJ alphanumeric
		{"CNPJ alphanumeric clean", "12ABC34501DE35", true},
		{"CNPJ all letters valid", "ABCDEFGHIJKLMN", true},
		{"CNPJ first char A", "A2796729000159", true},
		{"CNPJ last char Z", "2279672900015Z", true},
		{"CNPJ mixed boundary", "12ABC34501DE3Z", true},
		{"CNPJ all A's", "AAAAAAAAAAAAAA", true},
		{"CNPJ all Z's", "ZZZZZZZZZZZZZZ", true},
		{"CNPJ alternating", "1A2B3C4D5E6F7G", true},

		// Fast path failure cases - wrong length
		{"Empty string", "", false},
		{"Too short", "123", false},
		{"One digit", "1", false},
		{"Ten digits", "1234567890", false},
		{"CPF too short", "1234567890", false},
		{"CPF too long", "123456789012", false},
		{"Between CPF and CNPJ", "123456789012", false},
		{"CNPJ too short", "1234567890123", false},
		{"CNPJ too long", "123456789012345", false},
		{"Very long", "12345678901234567890", false},

		// Fast path failure cases - invalid characters for CPF (length 11)
		{"CPF with uppercase A", "7165668675A", false},
		{"CPF with uppercase Z", "7165668675Z", false},
		{"CPF with lowercase a", "7165668675a", false},
		{"CPF with lowercase z", "7165668675z", false},
		{"CPF with special char", "716566867.9", false},
		{"CPF with space", "71656686 59", false},
		{"CPF with dash", "716566867-9", false},
		{"CPF with mixed letters", "716A566B675", false},

		// Fast path failure cases - invalid characters for CNPJ (length 14)
		{"CNPJ with lowercase", "12abc34501de35", false},
		{"CNPJ mixed case", "12Abc34501De35", false},
		{"CNPJ with special chars", "12.ABC.345/0135", false},
		{"CNPJ with dash", "12-ABC-345-0135", false},
		{"CNPJ with space", "12ABC34501DE 5", false},
		{"CNPJ with dots", "12ABC34501DE.5", false},
		{"CNPJ with symbols", "12@BC34501DE35", false},
		{"CNPJ with invalid ASCII", "12\x80BC34501DE35", false},

		// Edge boundary cases
		{"ASCII boundary just below 0", "1234567890" + string(byte('0'-1)), false},         // 11 chars but invalid
		{"ASCII boundary just above 9", "1234567890" + string(byte('9'+1)), false},         // 11 chars but invalid
		{"ASCII boundary just below A", "12" + string(byte('A'-1)) + "BC34501DE35", false}, // 14 chars but invalid
		{"ASCII boundary just above Z", "12" + string(byte('Z'+1)) + "BC34501DE35", false}, // 14 chars but invalid
		{"ASCII boundary just below a", "12" + string(byte('a'-1)) + "BC34501DE35", false}, // 14 chars but invalid
		{"ASCII boundary just above z", "12" + string(byte('z'+1)) + "BC34501DE35", false}, // 14 chars but invalid

		// Control characters and whitespace
		{"CNPJ with tab", "12ABC34501\tDE35", false},
		{"CNPJ with newline", "12ABC34501\nDE35", false},
		{"CNPJ with carriage return", "12ABC34501\rDE35", false},
		{"CPF with null byte", "12345\x0067890", false},

		// Mixed valid and invalid patterns
		{"Length 11 mixed valid", "12345ABCDE6", false},  // CPF length with letters
		{"Length 14 all digits", "12345678901234", true}, // CNPJ length, all digits is valid
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAlreadyClean(tt.input)
			if result != tt.expected {
				t.Errorf("isAlreadyClean(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Test cleanString function
func TestCleanString(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedResult string
	}{
		// Basic cleaning cases - CPF
		{"CPF formatted", "716.566.867-59", "71656686759"},
		{"CPF with spaces", " 716 566 867 59 ", "71656686759"},
		{"CPF mixed formatting", "716-566.867/59", "71656686759"},
		{"CPF complex formatting", " 7!1@6#5$6%6^8&6*7(5)9 ", "71656686759"},

		// Basic cleaning cases - CNPJ numeric
		{"CNPJ numeric formatted", "22.796.729/0001-59", "22796729000159"},
		{"CNPJ with spaces", " 22 796 729 0001 59 ", "22796729000159"},
		{"CNPJ complex formatting", " 2!2@7#9$6%7^2&9*0(0)0{1}5[9] ", "22796729000159"},

		// Basic cleaning cases - CNPJ alphanumeric
		{"CNPJ alphanumeric formatted", "12.ABC.345/01DE-35", "12ABC34501DE35"},
		{"CNPJ mixed case", "12.abc.345/01de-35", "12ABC34501DE35"},
		{"CNPJ mixed case complex", " 1!2@.#a$b%c^.&3*4(5)/0{1}d[e]-3`5~ ", "12ABC34501DE35"},

		// Character type combinations
		{"Digits with letters", "123ABC456", "123ABC456"},
		{"Only letters", "ABCDEFG", "ABCDEFG"},
		{"Only lowercase letters", "abcdefg", "ABCDEFG"},
		{"Mixed case letters", "aBcDeF", "ABCDEF"},
		{"Only digits", "1234567", "1234567"},
		{"Alternating types", "a1B2c3D4e5", "A1B2C3D4E5"},

		// Edge cases
		{"Empty string", "", ""},
		{"Only symbols", "@#$%^&*()", ""},
		{"Only spaces", "   ", ""},
		{"Only whitespace", " \t\n\r ", ""},
		{"Mixed symbols and valid", "1@2#3A$B%4", "123AB4"},
		{"Single character cases", "1", "1"},
		{"Single letter upper", "A", "A"},
		{"Single letter lower", "a", "A"},
		{"Single symbol", "@", ""},

		// Capacity estimation cases
		{"Short input", "123AB", "123AB"},
		{"Medium input", "12345678901234567890", "12345678901234567890"},
		{"Long input triggers cap", "123456789012345678901234567890ABCDEFGHIJK", "123456789012345678901234567890ABCDEFGHIJK"},

		// All invalid characters
		{"Only formatting", "...-///()[]", ""},
		{"Long invalid", ".....................", ""},

		// ASCII boundary testing
		{"ASCII digits boundaries", "01239", "01239"},
		{"ASCII letters boundaries", "AZaz", "AZAZ"},
		{"ASCII edge chars", "/0:A[a{", "0AA"},

		// Real-world scenarios
		{"Brazilian CPF standard", "000.000.001-91", "00000000191"},
		{"Brazilian CNPJ standard", "11.222.333/0001-81", "11222333000181"},
		{"CNPJ alphanumeric standard", "12.ABC.345/01DE-35", "12ABC34501DE35"},
		{"User input with extra spaces", "  716.566.867-59  ", "71656686759"},
		{"Copy-paste with hidden chars", "716.566.867-59\n", "71656686759"},
		{"International formatting", "716 566 867 59", "71656686759"},

		// Case normalization verification
		{"All lowercase", "abcdefghijklmnopqrstuvwxyz", "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		{"All uppercase", "ABCDEFGHIJKLMNOPQRSTUVWXYZ", "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		{"Mixed case comprehensive", "aAbBcCdDeEfFgGhH", "AABBCCDDEEFFGGHH"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanString(tt.input)
			if result != tt.expectedResult {
				t.Errorf("cleanString(%q) = %q, want %q", tt.input, result, tt.expectedResult)
			}
		})
	}
}

// Test formatDocument function
func TestFormatDocument(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		pattern  string
		expected string
	}{
		// CPF formatting
		{"CPF clean digits", "71656686759", "XXX.XXX.XXX-XX", "716.566.867-59"},
		{"CPF all zeros", "00000000000", "XXX.XXX.XXX-XX", "000.000.000-00"},
		{"CPF all nines", "99999999999", "XXX.XXX.XXX-XX", "999.999.999-99"},

		// CNPJ numeric formatting
		{"CNPJ numeric clean", "22796729000159", "XX.XXX.XXX/XXXX-XX", "22.796.729/0001-59"},
		{"CNPJ all zeros", "00000000000000", "XX.XXX.XXX/XXXX-XX", "00.000.000/0000-00"},
		{"CNPJ all nines", "99999999999999", "XX.XXX.XXX/XXXX-XX", "99.999.999/9999-99"},

		// CNPJ alphanumeric formatting
		{"CNPJ alphanumeric", "12ABC34501DE35", "XX.XXX.XXX/XXXX-XX", "12.ABC.345/01DE-35"},
		{"CNPJ all letters", "ABCDEFGHIJKLMN", "XX.XXX.XXX/XXXX-XX", "AB.CDE.FGH/IJKL-MN"},
		{"CNPJ mixed", "A2B4C6D8E0F1G3", "XX.XXX.XXX/XXXX-XX", "A2.B4C.6D8/E0F1-G3"},

		// Edge cases
		{"Empty input", "", "XXX.XXX.XXX-XX", "..-"},
		{"Short input", "123", "XXX.XXX.XXX-XX", "123..-"},
		{"Pattern longer than input", "12345", "XX.XXX.XXX/XXXX-XX", "12.345./-"},
		{"No X in pattern", "123", "...-", "...-"},

		// Different patterns
		{"Custom pattern 1", "12345", "XX-XX-X", "12-34-5"},
		{"Custom pattern 2", "ABCDE", "X.X.X.X.X", "A.B.C.D.E"},
		{"Mixed separators", "123456", "XX/XX-XX", "12/34-56"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDocument(tt.input, tt.pattern)
			if got != tt.expected {
				t.Errorf("formatDocument(%q, %q) = %q, want %q", tt.input, tt.pattern, got, tt.expected)
			}
		})
	}
}
