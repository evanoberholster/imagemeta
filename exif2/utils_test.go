package exif2

import "testing"

func TestParseStrUint(t *testing.T) {
	tests := []struct {
		raw    string
		result int
	}{
		{"", 0},
		{"10000", 10000},
		{"15", 15},
		{"123", 123},
		{"2424", 2424},
		{"12345", 12345},
	}
	for _, test := range tests {
		result := parseStrUint([]byte(test.raw))
		if result != uint(test.result) {
			t.Errorf("error parseStrUint. Got %d wanted %d", result, test.result)
		}
	}
}

func TestTrim(t *testing.T) {
	tests := []struct {
		raw    string
		result string
	}{
		{"abcdefgh\000\000\000", "abcdefgh"},
		{"\n\n\n\n\000\000\000", ""},
	}
	for _, test := range tests {
		result := trimNULBuffer([]byte(test.raw))
		if string(result) != test.result {
			t.Errorf("error trim. got %s wanted %s", string(result), test.result)
		}
	}
}
