package strUtil

import (
	"unicode"
)

// IsBlankChar 判断是否为空白字符
// 空白符包括空格、制表符、全角空格和不间断空格等
func IsBlankChar(c rune) bool {
	return unicode.IsSpace(c) ||
		c == '\ufeff' || // 零宽不折行空格
		c == '\u202a' || // 从左到右嵌入
		c == '\u0000' || // 空字符
		c == '\u3164' || // Hangul Filler
		c == '\u2800' || // Braille Pattern Blank
		c == '\u180e' // MONGOLIAN VOWEL SEPARATOR
}

func IsBlank(s string) bool {
	if s == "" {
		return true
	}
	for _, c := range s {
		if !IsBlankChar(c) {
			return false
		}
	}
	return true
}
