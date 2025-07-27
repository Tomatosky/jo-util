package strUtil

import (
	"testing"
	"unicode"
)

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
