package strUtil

import (
	"testing"
)

// TestIsInt 测试 IsInt 函数
func TestIsInt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// 正常情况 - 有效整数
		{name: "正整数", input: "123", expected: true},
		{name: "零", input: "0", expected: true},
		{name: "多位数字", input: "123456789", expected: true},
		// 边界情况
		{name: "空字符串", input: "", expected: false},
		{name: "单个数字", input: "5", expected: true},
		// 异常情况 - 无效输入
		{name: "负整数", input: "-123", expected: true},
		{name: "正号前缀", input: "+123", expected: false},
		{name: "浮点数", input: "123.45", expected: false},
		{name: "科学计数法", input: "1e10", expected: false},
		{name: "包含字母", input: "123abc", expected: false},
		{name: "包含空格", input: "123 ", expected: false},
		{name: "前导零", input: "007", expected: false},
		{name: "包含特殊字符", input: "123-456", expected: false},
		{name: "仅包含小数点", input: ".", expected: false},
		{name: "十六进制数字", input: "0xFF", expected: false},
		{name: "中文字符", input: "一二三", expected: false},
		// 新增场景
		{name: "00_双个前导零", input: "00", expected: false},
		{name: "0123_有前导零的正数", input: "0123", expected: false},
		{name: "-0_负零", input: "-0", expected: false},
		{name: "1.2.3_多小数点", input: "1.2.3", expected: false},
		{name: "+_只有正号", input: "+", expected: false},
		{name: "-_只有负号", input: "-", expected: false},
		{name: " 123_带前导空格", input: " 123", expected: false},
		{name: "123 _带后导空格", input: "123 ", expected: false},
		{name: "1e10_科学计数法_确认", input: "1e10", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsInt(tt.input)
			if result != tt.expected {
				t.Errorf("IsInt(%q) = %v, 期望 %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsFloat 测试 IsFloat 函数
func TestIsFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// 正常情况 - 有效浮点数
		{name: "正浮点数", input: "123.45", expected: true},
		{name: "零浮点数", input: "0.0", expected: true},
		{name: "纯整数", input: "123", expected: true},
		{name: "负浮点数", input: "-123.45", expected: true},
		{name: "负零点几", input: "-0.5", expected: true},
		{name: "正号浮点数", input: "+123.45", expected: false},
		{name: "负整数", input: "-123", expected: true},
		{name: "正号整数", input: "+123", expected: false},
		// 边界情况
		{name: "空字符串", input: "", expected: false},
		{name: "仅小数点", input: ".", expected: false},
		{name: "仅负号和小数点", input: "-.", expected: false},
		{name: "前导小数点", input: ".123", expected: false},
		{name: "后缀小数点", input: "123.", expected: false},
		{name: "单个零", input: "0", expected: true},
		// 异常情况 - 无效输入
		{name: "包含空格", input: "123.45 ", expected: false},
		{name: "包含字母", input: "123.45abc", expected: false},
		{name: "多个小数点", input: "123.45.67", expected: false},
		{name: "科学计数法", input: "1.23e10", expected: false},
		{name: "逗号分隔符", input: "1,234.56", expected: false},
		{name: "十六进制", input: "0xFF.AA", expected: false},
		{name: "中文字符", input: "一二三", expected: false},
		{name: "特殊字符", input: "123!456", expected: false},
		{name: "多个符号", input: "+-123.45", expected: false},
		// 新增场景
		{name: "00_双个前导零", input: "00", expected: false},
		{name: "0123.45_有前导零的小数", input: "0123.45", expected: false},
		{name: "-0_负零", input: "-0", expected: false},
		{name: ".5_无整数部分", input: ".5", expected: false},
		{name: "5._无小数部分", input: "5.", expected: false},
		{name: "1.2.3_多小数点_确认", input: "1.2.3", expected: false},
		{name: "+_只有正号", input: "+", expected: false},
		{name: "-_只有负号", input: "-", expected: false},
		{name: " 1.23_带前导空格", input: " 1.23", expected: false},
		{name: "1.23 _带后导空格", input: "1.23 ", expected: false},
		{name: "1e10_科学计数法_确认", input: "1e10", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsFloat(tt.input)
			if result != tt.expected {
				t.Errorf("IsFloat(%q) = %v, 期望 %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsNumeric 测试 IsNumeric 函数
func TestIsNumeric(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// 正常情况 - 有效数字
		{name: "正浮点数", input: "123.45", expected: true},
		{name: "正整数", input: "123", expected: true},
		{name: "零", input: "0", expected: true},
		{name: "负浮点数", input: "-123.45", expected: true},
		{name: "负零点几", input: "-0.5", expected: true},
		{name: "负整数", input: "-123", expected: true},
		{name: "正号浮点数", input: "+123.45", expected: false},
		{name: "正号整数", input: "+123", expected: false},
		// 边界情况
		{name: "空字符串", input: "", expected: false},
		{name: "仅小数点", input: ".", expected: false},
		{name: "前导小数点", input: ".123", expected: false},
		{name: "后缀小数点", input: "123.", expected: false},
		{name: "仅负号", input: "-", expected: false},
		{name: "仅正号", input: "+", expected: false},
		// 异常情况 - 无效输入
		{name: "包含空格", input: "123.45 ", expected: false},
		{name: "包含字母", input: "123abc", expected: false},
		{name: "科学计数法", input: "1.23e10", expected: false},
		{name: "多个小数点", input: "123.45.67", expected: false},
		{name: "逗号分隔符", input: "1,234.56", expected: false},
		// 新增场景 - 与 IsFloat 相同
		{name: "00_双个前导零", input: "00", expected: false},
		{name: "0123.45_有前导零的小数", input: "0123.45", expected: false},
		{name: "-0_负零", input: "-0", expected: false},
		{name: ".5_无整数部分", input: ".5", expected: false},
		{name: "5._无小数部分", input: "5.", expected: false},
		{name: "1.2.3_多小数点_确认", input: "1.2.3", expected: false},
		{name: " 1.23_带前导空格", input: " 1.23", expected: false},
		{name: "1.23 _带后导空格", input: "1.23 ", expected: false},
		{name: "1e10_科学计数法_确认", input: "1e10", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNumeric(tt.input)
			if result != tt.expected {
				t.Errorf("IsNumeric(%q) = %v, 期望 %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsBlankChar 测试 IsBlankChar 函数
func TestIsBlankChar(t *testing.T) {
	tests := []struct {
		name     string
		input    rune
		expected bool
	}{
		// 正常情况 - 标准空白字符
		{name: "空格", input: ' ', expected: true},
		{name: "制表符", input: '\t', expected: true},
		{name: "换行符", input: '\n', expected: true},
		{name: "回车符", input: '\r', expected: true},
		{name: "垂直制表符", input: '\v', expected: true},
		{name: "换页符", input: '\f', expected: true},
		// 边界情况 - 特殊空白字符
		{name: "零宽不折行空格", input: '\ufeff', expected: true},
		{name: "从左到右嵌入", input: '\u202a', expected: true},
		{name: "空字符", input: '\u0000', expected: true},
		{name: "Hangul Filler", input: '\u3164', expected: true},
		{name: "Braille Pattern Blank", input: '\u2800', expected: true},
		{name: "MONGOLIAN VOWEL SEPARATOR", input: '\u180e', expected: true},
		// 新增场景
		{name: "不间断空格_NBSP", input: '\u00A0', expected: true}, // unicode.IsSpace 返回 true
		// 异常情况 - 非空白字符
		{name: "字母A", input: 'A', expected: false},
		{name: "字母a", input: 'a', expected: false},
		{name: "数字0", input: '0', expected: false},
		{name: "数字9", input: '9', expected: false},
		{name: "特殊字符下划线", input: '_', expected: false},
		{name: "特殊字符点号", input: '.', expected: false},
		{name: "中文字符", input: '你', expected: false},
		{name: "日文字符", input: 'あ', expected: false},
		{name: "表情符号", input: '😀', expected: false},
		{name: "ASCII感叹号", input: '!', expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsBlankChar(tt.input)
			if result != tt.expected {
				t.Errorf("IsBlankChar(%q) = %v, 期望 %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsBlank 测试 IsBlank 函数
func TestIsBlank(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// 正常情况 - 空白字符串
		{name: "空字符串", input: "", expected: true},
		{name: "单个空格", input: " ", expected: true},
		{name: "多个空格", input: "    ", expected: true},
		{name: "制表符", input: "\t", expected: true},
		{name: "换行符", input: "\n", expected: true},
		{name: "回车符", input: "\r", expected: true},
		{name: "混合空白字符", input: " \t\n\r", expected: true},
		// 边界情况 - 特殊空白字符
		{name: "零宽不折行空格", input: "\ufeff", expected: true},
		{name: "从左到右嵌入", input: "\u202a", expected: true},
		{name: "空字符", input: "\u0000", expected: true},
		{name: "Hangul Filler", input: "\u3164", expected: true},
		{name: "Braille Pattern Blank", input: "\u2800", expected: true},
		{name: "MONGOLIAN VOWEL SEPARATOR", input: "\u180e", expected: true},
		// 新增场景
		{name: "零宽字符_零宽不折行空格加从左到右嵌入", input: "\ufeff\u202a", expected: true},
		// 异常情况 - 非空白字符
		{name: "单个可见字符", input: "a", expected: false},
		{name: "英文字母", input: "abcd", expected: false},
		{name: "数字", input: "123", expected: false},
		{name: "中文字符", input: "你好", expected: false},
		{name: "特殊字符", input: "!@#$", expected: false},
		{name: "空格加可见字符", input: " a", expected: false},
		{name: "可见字符加空格", input: "a ", expected: false},
		{name: "混合空白和可见字符", input: " \t a \n", expected: false},
		{name: "零值数字", input: "0", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsBlank(tt.input)
			if result != tt.expected {
				t.Errorf("IsBlank(%q) = %v, 期望 %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsNumericAndIsFloatConsistency 测试 IsNumeric 和 IsFloat 的一致性
func TestIsNumericAndIsFloatConsistency(t *testing.T) {
	testCases := []string{
		"123",
		"123.45",
		"-123",
		"-123.45",
		"+123",
		"+123.45",
		"0",
		".123",
		"123.",
		"",
		"abc",
		"1.23e10",
	}

	for _, tc := range testCases {
		resultIsNumeric := IsNumeric(tc)
		resultIsFloat := IsFloat(tc)
		if resultIsNumeric != resultIsFloat {
			t.Errorf("IsNumeric(%q) = %v 与 IsFloat(%q) = %v 不一致",
				tc, resultIsNumeric, tc, resultIsFloat)
		}
	}
}

// TestIsNumericEdgeCases 测试 IsNumeric 的特殊边界情况
func TestIsNumericEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{name: "多个前导零", input: "0000", expected: false},
		{name: "多个前导零带小数点", input: "000.123", expected: false},
		{name: "小数点后多个零", input: "123.0000", expected: true},
		{name: "纯小数点后多零", input: ".0000", expected: false},
		{name: "负号和小数点", input: "-.", expected: false},
		{name: "数字中包含下划线", input: "1_234.56", expected: false},
		{name: "Infinity", input: "Infinity", expected: false},
		{name: "NaN", input: "NaN", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNumeric(tt.input)
			if result != tt.expected {
				t.Errorf("IsNumeric(%q) = %v, 期望 %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsIntLargeNumbers 测试 IsInt 的大数处理
func TestIsIntLargeNumbers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{name: "19位数字", input: "1234567890123456789", expected: true},
		{name: "20位数字", input: "12345678901234567890", expected: true},
		{name: "非常长的数字串", input: "1234567890123456789012345678901234567890", expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsInt(tt.input)
			if result != tt.expected {
				t.Errorf("IsInt(%q) = %v, 期望 %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsBlankUnicode 测试 IsBlank 对 Unicode 空白字符的处理
func TestIsBlankUnicode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{name: "不间断空格_NBSP", input: "\u00a0", expected: true},
		{name: "全角空格", input: "\u3000", expected: true},
		{name: "零宽空格", input: "\u200b", expected: false}, // 零宽空格不在 unicode.IsSpace 列表中
		{name: "零宽不折行空格已包含", input: "\ufeff", expected: true},
		{name: "行分隔符", input: "\u2028", expected: true},
		{name: "段落分隔符", input: "\u2029", expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsBlank(tt.input)
			if result != tt.expected {
				t.Errorf("IsBlank(%q) = %v, 期望 %v", tt.input, result, tt.expected)
			}
		})
	}
}
