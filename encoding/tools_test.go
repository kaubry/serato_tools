package encoding

import (
	"testing"
)

func TestInt32ToByteArray(t *testing.T) {
	arrayLength := 4
	number := uint32(1234567890)
	expected := []byte{73, 150, 2, 210}

	result := Int32ToByteArray(arrayLength, number)

	if len(result) != arrayLength {
		t.Errorf("Expected byte array length: %d, but got: %d", arrayLength, len(result))
	}

	for i, b := range result {
		if b != expected[i] {
			t.Errorf("Expected byte at index %d: %d, but got: %d", i, expected[i], b)
		}
	}
}

func TestEncodeUTF16WithBOM(t *testing.T) {
	s := "Hello, 世界!"
	addBOM := true

	expected := []byte{
		0xFE, 0xFF, //BIG-Endian BOM
		0x00, 0x48, // 'H'
		0x00, 0x65, // 'e'
		0x00, 0x6C, // 'l'
		0x00, 0x6C, // 'l'
		0x00, 0x6F, // 'o'
		0x00, 0x2C, // ','
		0x00, 0x20, // ' '
		0x4E, 0x16, // '世'
		0x75, 0x4C, // '界'
		0x00, 0x21, // '!'
	}

	result := EncodeUTF16(s, addBOM)

	compareResultAndExpected(t, result, expected)

}

func TestEncodeUTF16WithoutBOM(t *testing.T) {
	s := "Hello, 世界!"
	addBOM := false

	expected := []byte{
		0x00, 0x48, // 'H'
		0x00, 0x65, // 'e'
		0x00, 0x6C, // 'l'
		0x00, 0x6C, // 'l'
		0x00, 0x6F, // 'o'
		0x00, 0x2C, // ','
		0x00, 0x20, // ' '
		0x4E, 0x16, // '世'
		0x75, 0x4C, // '界'
		0x00, 0x21, // '!'
	}

	result := EncodeUTF16(s, addBOM)

	compareResultAndExpected(t, result, expected)

}

func compareResultAndExpected(t *testing.T, result []byte, expected []byte) {
	if len(result) != len(expected) {
		t.Errorf("Expected byte array length: %d, but got: %d", len(expected), len(result))
	}

	for i, b := range result {
		if b != expected[i] {
			t.Errorf("Expected byte at index %d: %d, but got: %d", i, expected[i], b)
		}
	}
}

func TestUTF16Bom(t *testing.T) {
	tests := []struct {
		name     string
		bom      []byte
		expected ByteOrder
	}{
		{
			name:     "Invalid BOM",
			bom:      []byte{},
			expected: InvalidBOM,
		},
		{
			name:     "Big-Endian BOM",
			bom:      []byte{0xFE, 0xFF, 0x00},
			expected: BigEndian,
		},
		{
			name:     "Little-Endian BOM",
			bom:      []byte{0xFF, 0xFE, 0x00},
			expected: LittleEndian,
		},
		{
			name:     "No BOM",
			bom:      []byte{0x00, 0x00, 0x00},
			expected: NoBOM,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UTF16Bom(tt.bom)
			if result != tt.expected {
				t.Errorf("Expected ByteOrder: %d, but got: %d", tt.expected, result)
			}
		})
	}
}

func TestDecodeUTF16(t *testing.T) {
	tests := []struct {
		name          string
		input         []byte
		expected      string
		expectedError string
	}{
		{
			name:     "Valid UTF-16 (Big-Endian)",
			input:    []byte{0xFE, 0xFF, 0x00, 0x48, 0x00, 0x65, 0x00, 0x6C, 0x00, 0x6C, 0x00, 0x6F},
			expected: "\ufeffHello",
		},
		{
			name: "Valid UTF-16 (Big-Endian)",
			input: []byte{0xFE, 0xFF, //BIG-Endian BOM
				0x00, 0x48, // 'H'
				0x00, 0x65, // 'e'
				0x00, 0x6C, // 'l'
				0x00, 0x6C, // 'l'
				0x00, 0x6F, // 'o'
				0x00, 0x2C, // ','
				0x00, 0x20, // ' '
				0x4E, 0x16, // '世'
				0x75, 0x4C, // '界'
				0x00, 0x21, // '!'
			},
			expected: "\ufeffHello, 世界!",
		},
		{
			name:          "Invalid length",
			input:         []byte{0xFE, 0xFF, 0x48, 0x65, 0x6C}, // Incomplete UTF-16 encoding
			expectedError: "must have even length byte slice",
		},
		{
			name:     "No BOM",
			input:    []byte{0x00, 0x48, 0x00, 0x65, 0x00, 0x6C, 0x00, 0x6C, 0x00, 0x6F}, // No BOM
			expected: "Hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := DecodeUTF16(tt.input)
			if err != nil {
				if err.Error() != tt.expectedError {
					t.Errorf("Test %q: Expected error: %s, but got: %s", tt.name, tt.expectedError, err.Error())
				}
			} else {
				if result != tt.expected {
					t.Errorf("Test %q: Expected result: %q, but got: %q", tt.name, tt.expected, result)
				}
			}
		})
	}
}
