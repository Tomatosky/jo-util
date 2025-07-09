package strUtil

import (
	"testing"
	"unicode"
)

func TestToInt(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{"valid int", "123", 123, false},
		{"invalid int", "abc", 0, true},
		{"empty string", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tt.wantErr {
					t.Errorf("ToInt() panicked unexpectedly: %v", r)
				}
			}()

			got := ToInt(tt.input)
			if got != tt.want && !tt.wantErr {
				t.Errorf("ToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt32(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int32
		wantErr bool
	}{
		{"valid int32", "12345", 12345, false},
		{"out of range", "2147483648", 0, true}, // 超过int32最大值
		{"invalid int32", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tt.wantErr {
					t.Errorf("ToInt32() panicked unexpectedly: %v", r)
				}
			}()

			got := ToInt32(tt.input)
			if got != tt.want && !tt.wantErr {
				t.Errorf("ToInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt64(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		{"valid int64", "9223372036854775807", 9223372036854775807, false}, // int64最大值
		{"invalid int64", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tt.wantErr {
					t.Errorf("ToInt64() panicked unexpectedly: %v", r)
				}
			}()

			got := ToInt64(tt.input)
			if got != tt.want && !tt.wantErr {
				t.Errorf("ToInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToFloat32(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    float32
		wantErr bool
	}{
		{"valid float32", "3.14", 3.14, false},
		{"invalid float32", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tt.wantErr {
					t.Errorf("ToFloat32() panicked unexpectedly: %v", r)
				}
			}()

			got := ToFloat32(tt.input)
			if got != tt.want && !tt.wantErr {
				t.Errorf("ToFloat32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToFloat64(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    float64
		wantErr bool
	}{
		{"valid float64", "3.141592653589793", 3.141592653589793, false},
		{"invalid float64", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tt.wantErr {
					t.Errorf("ToFloat64() panicked unexpectedly: %v", r)
				}
			}()

			got := ToFloat64(tt.input)
			if got != tt.want && !tt.wantErr {
				t.Errorf("ToFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{"nil", nil, ""},
		{"float32", float32(3.14), "3.14"},
		{"float64", 3.141592653589793, "3.141592653589793"},
		{"int", 42, "42"},
		{"int8", int8(8), "8"},
		{"int16", int16(16), "16"},
		{"int32", int32(32), "32"},
		{"int64", int64(64), "64"},
		{"uint", uint(42), "42"},
		{"uint8", uint8(8), "8"},
		{"uint16", uint16(16), "16"},
		{"uint32", uint32(32), "32"},
		{"uint64", uint64(64), "64"},
		{"string", "hello", "hello"},
		{"[]byte", []byte("world"), "world"},
		{"struct", struct{ Name string }{Name: "test"}, `{"Name":"test"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToString(tt.input)
			if got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsBlankChar(t *testing.T) {
	tests := []struct {
		name string
		c    rune
		want bool
	}{
		{name: "regular space", c: ' ', want: true},
		{name: "tab", c: '\t', want: true},
		{name: "newline", c: '\n', want: true},
		{name: "zero width no-break space", c: '\ufeff', want: true},
		{name: "left-to-right embedding", c: '\u202a', want: true},
		{name: "null character", c: '\u0000', want: true},
		{name: "Hangul Filler", c: '\u3164', want: true},
		{name: "Braille Pattern Blank", c: '\u2800', want: true},
		{name: "Mongolian Vowel Separator", c: '\u180e', want: true},
		{name: "non-blank character", c: 'a', want: false},
		{name: "digit", c: '1', want: false},
		{name: "symbol", c: '@', want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsBlankChar(tt.c)
			if got != tt.want {
				t.Errorf("IsBlankChar(%q) = %v, want %v", tt.c, got, tt.want)
			}
		})
	}
}

func TestIsBlankChar_UnicodeSpace(t *testing.T) {
	// Test all unicode space characters
	for r := rune(0); r <= unicode.MaxRune; r++ {
		if unicode.IsSpace(r) && !IsBlankChar(r) {
			t.Errorf("IsBlankChar(%q) = false, want true for unicode space character", r)
		}
	}
}

func TestIsBlank(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"empty string", "", true},
		{"only spaces", "   ", true},
		{"only tabs", "\t\t\t", true},
		{"only newlines", "\n\n\n", true},
		{"mixed whitespace", " \t\n\r", true},
		{"non-whitespace chars", "abc", false},
		{"whitespace with chars", "  abc  ", false},
		{"leading whitespace", "  abc", false},
		{"trailing whitespace", "abc  ", false},
		{"unicode whitespace", "\u2000\u2001", true}, // en quad, em quad
		{"unicode non-whitespace", "你好", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsBlank(tt.input)
			if result != tt.expected {
				t.Errorf("IsBlank(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
