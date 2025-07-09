package strUtil

import (
	"encoding/json"
	"strconv"
	"unicode"
)

// ToInt 将字符串转换为 int
func ToInt(s string) int {
	atoi, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return atoi
}

// ToInt32 将字符串转换为 int32
func ToInt32(s string) int32 {
	i32, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		panic(err)
	}
	return int32(i32)
}

// ToInt64 将字符串转换为 int64
func ToInt64(s string) int64 {
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return i64
}

// ToFloat32 将字符串转换为 float32
func ToFloat32(s string) float32 {
	f32, err := strconv.ParseFloat(s, 32)
	if err != nil {
		panic(err)
	}
	return float32(f32)
}

// ToFloat64 将字符串转换为 float64
func ToFloat64(s string) float64 {
	f64, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f64
}

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

func ToString(value any) string {
	if value == nil {
		return ""
	}

	switch val := value.(type) {
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case int:
		return strconv.FormatInt(int64(val), 10)
	case int8:
		return strconv.FormatInt(int64(val), 10)
	case int16:
		return strconv.FormatInt(int64(val), 10)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case int64:
		return strconv.FormatInt(val, 10)
	case uint:
		return strconv.FormatUint(uint64(val), 10)
	case uint8:
		return strconv.FormatUint(uint64(val), 10)
	case uint16:
		return strconv.FormatUint(uint64(val), 10)
	case uint32:
		return strconv.FormatUint(uint64(val), 10)
	case uint64:
		return strconv.FormatUint(val, 10)
	case string:
		return val
	case []byte:
		return string(val)
	default:
		b, err := json.Marshal(val)
		if err != nil {
			return ""
		}
		return string(b)
	}
}
