package strUtil

import (
	"regexp"
	"unicode"
)

var intReg = regexp.MustCompile(`^\d+$`)
var floatReg = regexp.MustCompile(`^[-+]?\d*\.?\d+$`)
var numReg = regexp.MustCompile(`^[-+]?\d*\.?\d+$`)

// IsInt check the string is an integer number
func IsInt(s string) bool { return intReg.MatchString(s) }

// IsFloat check the string is a float number
func IsFloat(s string) bool { return floatReg.MatchString(s) }

// IsNumeric returns true if the given string is a numeric, otherwise false.
func IsNumeric(s string) bool { return numReg.MatchString(s) }

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
