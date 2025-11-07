package dateUtil

import (
	"testing"
	"time"
)

func TestGetTime(t *testing.T) {
	tests := []struct {
		name      string
		timestamp int64
		expected  time.Time
	}{
		{"Zero timestamp", 0, time.Unix(0, 0).In(Loc)},
		{"Current timestamp", time.Now().Unix(), time.Unix(time.Now().Unix(), 0).In(Loc)},
		{"Future timestamp", 1893456000, time.Unix(1893456000, 0).In(Loc)}, // 2030-01-01 00:00:00
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTime(tt.timestamp)
			if !got.Equal(tt.expected) {
				t.Errorf("GetTime() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBeginOfDay(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "普通日期-中午时间",
			input:    time.Date(2023, 5, 15, 12, 30, 45, 0, Loc),
			expected: time.Date(2023, 5, 15, 0, 0, 0, 0, Loc),
		},
		{
			name:     "普通日期-午夜时间",
			input:    time.Date(2023, 5, 15, 0, 0, 0, 0, Loc),
			expected: time.Date(2023, 5, 15, 0, 0, 0, 0, Loc),
		},
		{
			name:     "普通日期-接近午夜时间",
			input:    time.Date(2023, 5, 15, 23, 59, 59, 999999, Loc),
			expected: time.Date(2023, 5, 15, 0, 0, 0, 0, Loc),
		},
		{
			name:     "闰年2月29日",
			input:    time.Date(2020, 2, 29, 15, 30, 0, 0, Loc),
			expected: time.Date(2020, 2, 29, 0, 0, 0, 0, Loc),
		},
		{
			name:     "跨年时间",
			input:    time.Date(2023, 12, 31, 23, 59, 59, 0, Loc),
			expected: time.Date(2023, 12, 31, 0, 0, 0, 0, Loc),
		},
		{
			name:     "夏令时转换时间",
			input:    time.Date(2023, 3, 12, 2, 30, 0, 0, Loc), // 夏令时开始
			expected: time.Date(2023, 3, 12, 0, 0, 0, 0, Loc),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BeginOfDay(tt.input)
			if !got.Equal(tt.expected) {
				t.Errorf("BeginOfDay(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestBeginOfWeek(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "周一当天",
			input:    time.Date(2023, 5, 15, 12, 30, 45, 0, Loc), // 2023-05-15 周一
			expected: time.Date(2023, 5, 15, 0, 0, 0, 0, Loc),    // 同一天零点
		},
		{
			name:     "周二",
			input:    time.Date(2023, 5, 16, 8, 0, 0, 0, Loc), // 2023-05-16 周二
			expected: time.Date(2023, 5, 15, 0, 0, 0, 0, Loc), // 本周一零点
		},
		{
			name:     "周日",
			input:    time.Date(2023, 5, 21, 23, 59, 59, 0, Loc), // 2023-05-21 周日
			expected: time.Date(2023, 5, 15, 0, 0, 0, 0, Loc),    // 本周一零点
		},
		{
			name:     "跨月情况-月初是周一",
			input:    time.Date(2023, 6, 4, 15, 0, 0, 0, Loc), // 2023-06-04 周日
			expected: time.Date(2023, 5, 29, 0, 0, 0, 0, Loc), // 上周一(5月29日)零点
		},
		{
			name:     "跨年情况-年初是周一",
			input:    time.Date(2023, 1, 1, 0, 0, 1, 0, Loc),   // 2023-01-01 周日
			expected: time.Date(2022, 12, 26, 0, 0, 0, 0, Loc), // 上周一(2022-12-26)零点
		},
		{
			name:     "闰年2月",
			input:    time.Date(2020, 2, 29, 12, 0, 0, 0, Loc), // 2020-02-29 周六
			expected: time.Date(2020, 2, 24, 0, 0, 0, 0, Loc),  // 本周一零点
		},
		{
			name:     "夏令时转换周",
			input:    time.Date(2023, 3, 12, 2, 30, 0, 0, Loc), // 2023-03-12 周日(夏令时开始)
			expected: time.Date(2023, 3, 6, 0, 0, 0, 0, Loc),   // 本周一零点
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BeginOfWeek(tt.input)
			if !got.Equal(tt.expected) {
				t.Errorf("BeginOfWeek(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestBeginOfMonth(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "当月第一天",
			input:    time.Date(2023, 5, 1, 0, 0, 0, 0, Loc), // 2023-05-01
			expected: time.Date(2023, 5, 1, 0, 0, 0, 0, Loc), // 同一天零点
		},
		{
			name:     "月中某天",
			input:    time.Date(2023, 5, 15, 12, 30, 45, 0, Loc), // 2023-05-15
			expected: time.Date(2023, 5, 1, 0, 0, 0, 0, Loc),     // 当月第一天零点
		},
		{
			name:     "当月最后一天",
			input:    time.Date(2023, 5, 31, 23, 59, 59, 0, Loc), // 2023-05-31
			expected: time.Date(2023, 5, 1, 0, 0, 0, 0, Loc),     // 当月第一天零点
		},
		{
			name:     "闰年2月",
			input:    time.Date(2020, 2, 29, 12, 0, 0, 0, Loc), // 2020-02-29
			expected: time.Date(2020, 2, 1, 0, 0, 0, 0, Loc),   // 当月第一天零点
		},
		{
			name:     "非闰年2月",
			input:    time.Date(2023, 2, 28, 23, 59, 59, 0, Loc), // 2023-02-28
			expected: time.Date(2023, 2, 1, 0, 0, 0, 0, Loc),     // 当月第一天零点
		},
		{
			name:     "跨年12月",
			input:    time.Date(2023, 12, 31, 0, 0, 1, 0, Loc), // 2023-12-31
			expected: time.Date(2023, 12, 1, 0, 0, 0, 0, Loc),  // 当月第一天零点
		},
		{
			name:     "跨年1月",
			input:    time.Date(2023, 1, 1, 0, 0, 0, 0, Loc), // 2023-01-01
			expected: time.Date(2023, 1, 1, 0, 0, 0, 0, Loc), // 同一天零点
		},
		{
			name:     "夏令时转换月",
			input:    time.Date(2023, 3, 15, 2, 30, 0, 0, Loc), // 2023-03-15
			expected: time.Date(2023, 3, 1, 0, 0, 0, 0, Loc),   // 当月第一天零点
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BeginOfMonth(tt.input)
			if !got.Equal(tt.expected) {
				t.Errorf("BeginOfMonth(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestEndOfDay(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "普通日期",
			input:    time.Date(2023, 1, 1, 12, 30, 15, 123456, Loc),
			expected: time.Date(2023, 1, 1, 23, 59, 59, 0, Loc),
		},
		{
			name:     "午夜零点",
			input:    time.Date(2023, 1, 1, 0, 0, 0, 0, Loc),
			expected: time.Date(2023, 1, 1, 23, 59, 59, 0, Loc),
		},
		{
			name:     "23:59:59",
			input:    time.Date(2023, 1, 1, 23, 59, 59, 999999, Loc),
			expected: time.Date(2023, 1, 1, 23, 59, 59, 0, Loc),
		},
		{
			name:     "闰年2月29日",
			input:    time.Date(2020, 2, 29, 15, 30, 0, 0, Loc),
			expected: time.Date(2020, 2, 29, 23, 59, 59, 0, Loc),
		},
		{
			name:     "非闰年2月28日",
			input:    time.Date(2023, 2, 28, 15, 30, 0, 0, Loc),
			expected: time.Date(2023, 2, 28, 23, 59, 59, 0, Loc),
		},
		{
			name:     "12月31日",
			input:    time.Date(2023, 12, 31, 15, 30, 0, 0, Loc),
			expected: time.Date(2023, 12, 31, 23, 59, 59, 0, Loc),
		},
		{
			name:     "1月1日",
			input:    time.Date(2023, 1, 1, 15, 30, 0, 0, Loc),
			expected: time.Date(2023, 1, 1, 23, 59, 59, 0, Loc),
		},
		{
			name:     "带纳秒的时间",
			input:    time.Date(2023, 1, 1, 12, 30, 15, 999999999, Loc),
			expected: time.Date(2023, 1, 1, 23, 59, 59, 0, Loc),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := EndOfDay(tt.input)
			if !actual.Equal(tt.expected) {
				t.Errorf("EndOfDay(%v) = %v, want %v", tt.input, actual, tt.expected)
			}
			// 验证时间部分确实是23:59:59
			if actual.Hour() != 23 || actual.Minute() != 59 || actual.Second() != 59 {
				t.Errorf("EndOfDay(%v) time part is not 23:59:59, got %02d:%02d:%02d",
					tt.input, actual.Hour(), actual.Minute(), actual.Second())
			}
			// 验证纳秒部分为0
			if actual.Nanosecond() != 0 {
				t.Errorf("EndOfDay(%v) nanosecond is not 0, got %d", tt.input, actual.Nanosecond())
			}
			// 验证日期部分不变
			if actual.Year() != tt.input.Year() || actual.Month() != tt.input.Month() || actual.Day() != tt.input.Day() {
				t.Errorf("EndOfDay(%v) date part changed, got %v", tt.input, actual)
			}
		})
	}
}

func TestEndOfWeek(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "周一",
			input:    time.Date(2023, 6, 5, 12, 0, 0, 0, Loc), // 2023-06-05 周一
			expected: time.Date(2023, 6, 11, 23, 59, 59, 0, Loc),
		},
		{
			name:     "周二",
			input:    time.Date(2023, 6, 6, 12, 0, 0, 0, Loc), // 2023-06-06 周二
			expected: time.Date(2023, 6, 11, 23, 59, 59, 0, Loc),
		},
		{
			name:     "周三",
			input:    time.Date(2023, 6, 7, 12, 0, 0, 0, Loc), // 2023-06-07 周三
			expected: time.Date(2023, 6, 11, 23, 59, 59, 0, Loc),
		},
		{
			name:     "周四",
			input:    time.Date(2023, 6, 8, 12, 0, 0, 0, Loc), // 2023-06-08 周四
			expected: time.Date(2023, 6, 11, 23, 59, 59, 0, Loc),
		},
		{
			name:     "周五",
			input:    time.Date(2023, 6, 9, 12, 0, 0, 0, Loc), // 2023-06-09 周五
			expected: time.Date(2023, 6, 11, 23, 59, 59, 0, Loc),
		},
		{
			name:     "周六",
			input:    time.Date(2023, 6, 10, 12, 0, 0, 0, Loc), // 2023-06-10 周六
			expected: time.Date(2023, 6, 11, 23, 59, 59, 0, Loc),
		},
		{
			name:     "周日",
			input:    time.Date(2023, 6, 11, 12, 0, 0, 0, Loc), // 2023-06-11 周日
			expected: time.Date(2023, 6, 11, 23, 59, 59, 0, Loc),
		},
		{
			name:     "跨月周日",
			input:    time.Date(2023, 5, 29, 12, 0, 0, 0, Loc),  // 2023-05-29 周一
			expected: time.Date(2023, 6, 4, 23, 59, 59, 0, Loc), // 2023-06-04 周日
		},
		{
			name:     "跨年周日",
			input:    time.Date(2022, 12, 26, 12, 0, 0, 0, Loc), // 2022-12-26 周一
			expected: time.Date(2023, 1, 1, 23, 59, 59, 0, Loc), // 2023-01-01 周日
		},
		{
			name:     "闰年2月",
			input:    time.Date(2020, 2, 24, 12, 0, 0, 0, Loc),  // 2020-02-24 周一
			expected: time.Date(2020, 3, 1, 23, 59, 59, 0, Loc), // 2020-03-01 周日
		},
		{
			name:     "周日23:59:59",
			input:    time.Date(2023, 6, 11, 23, 59, 59, 999, Loc),
			expected: time.Date(2023, 6, 11, 23, 59, 59, 0, Loc),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := EndOfWeek(tt.input)
			if !actual.Equal(tt.expected) {
				t.Errorf("EndOfWeek(%v) = %v, want %v", tt.input, actual, tt.expected)
			}
			// 验证确实是周日
			if actual.Weekday() != time.Sunday {
				t.Errorf("EndOfWeek(%v) is not Sunday, got %v", tt.input, actual.Weekday())
			}
			// 验证时间部分确实是23:59:59
			if actual.Hour() != 23 || actual.Minute() != 59 || actual.Second() != 59 {
				t.Errorf("EndOfWeek(%v) time part is not 23:59:59, got %02d:%02d:%02d",
					tt.input, actual.Hour(), actual.Minute(), actual.Second())
			}
			// 验证纳秒部分为0
			if actual.Nanosecond() != 0 {
				t.Errorf("EndOfWeek(%v) nanosecond is not 0, got %d", tt.input, actual.Nanosecond())
			}
		})
	}
}

func TestEndOfMonth(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "1月31天",
			input:    time.Date(2023, 1, 15, 12, 30, 15, 123456, Loc),
			expected: time.Date(2023, 1, 31, 23, 59, 59, 0, Loc),
		},
		{
			name:     "2月非闰年28天",
			input:    time.Date(2023, 2, 10, 0, 0, 0, 0, Loc),
			expected: time.Date(2023, 2, 28, 23, 59, 59, 0, Loc),
		},
		{
			name:     "2月闰年29天",
			input:    time.Date(2020, 2, 15, 12, 0, 0, 0, Loc),
			expected: time.Date(2020, 2, 29, 23, 59, 59, 0, Loc),
		},
		{
			name:     "4月30天",
			input:    time.Date(2023, 4, 1, 23, 59, 59, 999999, Loc),
			expected: time.Date(2023, 4, 30, 23, 59, 59, 0, Loc),
		},
		{
			name:     "12月31天",
			input:    time.Date(2023, 12, 31, 0, 0, 0, 0, Loc),
			expected: time.Date(2023, 12, 31, 23, 59, 59, 0, Loc),
		},
		{
			name:     "月末最后一天",
			input:    time.Date(2023, 1, 31, 23, 59, 59, 999999, Loc),
			expected: time.Date(2023, 1, 31, 23, 59, 59, 0, Loc),
		},
		{
			name:     "跨年12月",
			input:    time.Date(2022, 12, 1, 0, 0, 0, 0, Loc),
			expected: time.Date(2022, 12, 31, 23, 59, 59, 0, Loc),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := EndOfMonth(tt.input)
			if !actual.Equal(tt.expected) {
				t.Errorf("EndOfMonth(%v) = %v, want %v", tt.input, actual, tt.expected)
			}
			// 验证时间部分确实是23:59:59
			if actual.Hour() != 23 || actual.Minute() != 59 || actual.Second() != 59 {
				t.Errorf("EndOfMonth(%v) time part is not 23:59:59, got %02d:%02d:%02d",
					tt.input, actual.Hour(), actual.Minute(), actual.Second())
			}
			// 验证纳秒部分为0
			if actual.Nanosecond() != 0 {
				t.Errorf("EndOfMonth(%v) nanosecond is not 0, got %d", tt.input, actual.Nanosecond())
			}
			// 验证月份不变
			if actual.Month() != tt.input.Month() || actual.Year() != tt.input.Year() {
				t.Errorf("EndOfMonth(%v) month/year changed, got %v", tt.input, actual)
			}
		})
	}
}

func TestIsSameDay(t *testing.T) {
	tests := []struct {
		name     string
		t1       time.Time
		t2       time.Time
		expected bool
	}{
		{
			name:     "完全相同的时间",
			t1:       time.Date(2023, 5, 15, 12, 30, 45, 0, Loc),
			t2:       time.Date(2023, 5, 15, 12, 30, 45, 0, Loc),
			expected: true,
		},
		{
			name:     "同一天不同时间",
			t1:       time.Date(2023, 5, 15, 0, 0, 0, 0, Loc),
			t2:       time.Date(2023, 5, 15, 23, 59, 59, 0, Loc),
			expected: true,
		},
		{
			name:     "不同天",
			t1:       time.Date(2023, 5, 15, 23, 59, 59, 0, Loc),
			t2:       time.Date(2023, 5, 16, 0, 0, 0, 0, Loc),
			expected: false,
		},
		{
			name:     "不同月",
			t1:       time.Date(2023, 5, 31, 23, 59, 59, 0, Loc),
			t2:       time.Date(2023, 6, 1, 0, 0, 0, 0, Loc),
			expected: false,
		},
		{
			name:     "不同年",
			t1:       time.Date(2022, 12, 31, 23, 59, 59, 0, Loc),
			t2:       time.Date(2023, 1, 1, 0, 0, 0, 0, Loc),
			expected: false,
		},
		{
			name:     "闰年2月29日",
			t1:       time.Date(2020, 2, 29, 0, 0, 0, 0, Loc),
			t2:       time.Date(2020, 2, 29, 23, 59, 59, 0, Loc),
			expected: true,
		},
		{
			name:     "边界时间-午夜前",
			t1:       time.Date(2023, 5, 15, 23, 59, 59, 999999999, Loc),
			t2:       time.Date(2023, 5, 15, 0, 0, 0, 0, Loc),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSameDay(tt.t1, tt.t2)
			if result != tt.expected {
				t.Errorf("IsSameDay(%v, %v) = %v, want %v", tt.t1, tt.t2, result, tt.expected)
			}
		})
	}
}

func TestIsSameWeek(t *testing.T) {
	tests := []struct {
		name     string
		t1       time.Time
		t2       time.Time
		expected bool
	}{
		{
			name:     "完全相同的时间",
			t1:       time.Date(2023, 5, 15, 12, 30, 45, 0, Loc), // 周一
			t2:       time.Date(2023, 5, 15, 12, 30, 45, 0, Loc),
			expected: true,
		},
		{
			name:     "同一周不同天",
			t1:       time.Date(2023, 5, 15, 0, 0, 0, 0, Loc),    // 周一
			t2:       time.Date(2023, 5, 21, 23, 59, 59, 0, Loc), // 周日
			expected: true,
		},
		{
			name:     "跨周边界-周日和周一",
			t1:       time.Date(2023, 5, 21, 23, 59, 59, 0, Loc), // 周日
			t2:       time.Date(2023, 5, 22, 0, 0, 0, 0, Loc),    // 周一
			expected: false,
		},
		{
			name:     "不同周",
			t1:       time.Date(2023, 5, 14, 23, 59, 59, 0, Loc), // 第19周周日
			t2:       time.Date(2023, 5, 15, 0, 0, 0, 0, Loc),    // 第20周周一
			expected: false,
		},
		{
			name:     "跨年周-第52周和第1周",
			t1:       time.Date(2022, 12, 31, 23, 59, 59, 0, Loc), // 2022年第52周周六
			t2:       time.Date(2023, 1, 1, 0, 0, 0, 0, Loc),      // 2023年第52周周日
			expected: true,                                        // 注意：ISO周计算中，2023-01-01仍属于2022年第52周
		},
		{
			name:     "跨年周-不同年",
			t1:       time.Date(2021, 12, 31, 23, 59, 59, 0, Loc), // 2021年第52周周五
			t2:       time.Date(2022, 1, 1, 0, 0, 0, 0, Loc),      // 2022年第52周周六
			expected: true,                                        // 2021-12-31属于2021年第52周，2022-01-01属于2021年第52周
		},
		{
			name:     "闰年周",
			t1:       time.Date(2020, 2, 24, 0, 0, 0, 0, Loc),   // 2020年第9周周一
			t2:       time.Date(2020, 3, 1, 23, 59, 59, 0, Loc), // 2020年第9周周日
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSameWeek(tt.t1, tt.t2)
			if result != tt.expected {
				t.Errorf("IsSameWeek(%v, %v) = %v, want %v", tt.t1, tt.t2, result, tt.expected)
			}
		})
	}
}

func TestParseToTime(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name     string
		str      string
		format   string
		timezone []string
		want     time.Time
		wantErr  bool
	}{
		// 测试常见日期时间格式
		{
			name:    "标准日期时间格式",
			str:     "2023-05-15 14:30:00",
			format:  "yyyy-mm-dd hh:mm:ss",
			want:    time.Date(2023, 5, 15, 14, 30, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "简略日期格式",
			str:     "2023-05-15",
			format:  "yyyy-mm-dd",
			want:    time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "年月格式",
			str:     "2023-05",
			format:  "yyyy-mm",
			want:    time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "时间格式",
			str:     "14:30:00",
			format:  "hh:mm:ss",
			want:    time.Date(0, 1, 1, 14, 30, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "斜杠日期格式",
			str:     "2023/05/15 14:30:00",
			format:  "yyyy/mm/dd hh:mm:ss",
			want:    time.Date(2023, 5, 15, 14, 30, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "紧凑日期格式",
			str:     "20230515",
			format:  "yyyymmdd",
			want:    time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},

		// 测试时区
		{
			name:     "带时区转换",
			str:      "2023-05-15 14:30:00",
			format:   "yyyy-mm-dd hh:mm:ss",
			timezone: []string{"Asia/Shanghai"},
			want:     time.Date(2023, 5, 15, 14, 30, 0, 0, time.FixedZone("CST", 8*60*60)),
			wantErr:  false,
		},
		{
			name:     "无效时区",
			str:      "2023-05-15 14:30:00",
			format:   "yyyy-mm-dd hh:mm:ss",
			timezone: []string{"Invalid/Zone"},
			wantErr:  true,
		},

		// 测试错误情况
		{
			name:    "不支持的格式",
			str:     "2023-05-15",
			format:  "invalid-format",
			wantErr: true,
		},
		{
			name:    "日期与格式不匹配",
			str:     "2023-05-15",
			format:  "yyyy-mm-dd hh:mm:ss",
			wantErr: true,
		},
		{
			name:    "无效日期",
			str:     "2023-02-30", // 2月没有30号
			format:  "yyyy-mm-dd",
			wantErr: true,
		},
		{
			name:    "空字符串",
			str:     "",
			format:  "yyyy-mm-dd",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseToTime(tt.str, tt.format, tt.timezone...)

			// 检查错误情况是否符合预期
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseToTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 如果没有错误，检查结果是否正确
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("ParseToTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBetweenMinute(t *testing.T) {
	t1 := time.Date(2023, 1, 1, 12, 30, 15, 0, Loc)
	t2 := time.Date(2023, 1, 1, 12, 31, 45, 0, Loc)

	// 测试不重置秒数
	got := BetweenMinute(t1, t2, false)
	if got != 1 {
		t.Errorf("BetweenMinute(false) = %d, want 1", got)
	}

	// 测试重置秒数
	got = BetweenMinute(t1, t2, true)
	if got != 1 {
		t.Errorf("BetweenMinute(true) = %d, want 1", got)
	}

	// 测试更大的时间差
	t3 := time.Date(2023, 1, 1, 12, 45, 0, 0, Loc)
	got = BetweenMinute(t1, t3, true)
	if got != 15 {
		t.Errorf("BetweenMinute(false) = %d, want 15", got)
	}
}

func TestBetweenHour(t *testing.T) {
	t1 := time.Date(2023, 1, 1, 12, 30, 15, 0, Loc)
	t2 := time.Date(2023, 1, 1, 14, 31, 45, 0, Loc)

	// 测试不重置分钟和秒数
	got := BetweenHour(t1, t2, false)
	if got != 2 {
		t.Errorf("BetweenHour(false) = %d, want 2", got)
	}

	// 测试重置分钟和秒数
	got = BetweenHour(t1, t2, true)
	if got != 2 {
		t.Errorf("BetweenHour(true) = %d, want 2", got)
	}

	// 测试跨天的情况
	t3 := time.Date(2023, 1, 2, 1, 0, 0, 0, Loc)
	got = BetweenHour(t1, t3, true)
	if got != 13 {
		t.Errorf("BetweenHour(false) = %d, want 13", got)
	}
}

func TestBetweenDay(t *testing.T) {
	t1 := time.Date(2023, 1, 1, 12, 30, 15, 0, Loc)
	t2 := time.Date(2023, 1, 3, 14, 31, 45, 0, Loc)

	// 测试不重置时间
	got := BetweenDay(t1, t2, false)
	if got != 2 {
		t.Errorf("BetweenDay(false) = %d, want 2", got)
	}

	// 测试重置时间
	got = BetweenDay(t1, t2, true)
	if got != 2 {
		t.Errorf("BetweenDay(true) = %d, want 2", got)
	}

	// 测试跨月的情况
	t3 := time.Date(2023, 2, 1, 0, 0, 0, 0, Loc)
	got = BetweenDay(t1, t3, true)
	if got != 31 {
		t.Errorf("BetweenDay(true) = %d, want 31", got)
	}
}

func TestBetweenWeek(t *testing.T) {
	t1 := time.Date(2023, 1, 1, 12, 0, 0, 0, Loc)  // 周日
	t2 := time.Date(2023, 1, 10, 12, 0, 0, 0, Loc) // 下周二

	// 测试不重置
	got := BetweenWeek(t1, t2, false)
	if got != 1 {
		t.Errorf("BetweenWeek(false) = %d, want 1", got)
	}

	// 测试多周差
	t3 := time.Date(2023, 1, 22, 0, 0, 0, 0, Loc)
	got = BetweenWeek(t1, t3, true)
	if got != 3 {
		t.Errorf("BetweenWeek(true) = %d, want 3", got)
	}
}

func TestBetweenMonth(t *testing.T) {
	t1 := time.Date(2023, 1, 15, 12, 0, 0, 0, Loc)
	t2 := time.Date(2023, 3, 20, 12, 0, 0, 0, Loc)

	// 测试不重置
	got := BetweenMonth(t1, t2, false)
	if got != 2 {
		t.Errorf("BetweenMonth(false) = %d, want 2", got)
	}

	// 测试重置
	got = BetweenMonth(t1, t2, true)
	if got != 2 {
		t.Errorf("BetweenMonth(true) = %d, want 2", got)
	}

	// 测试跨年
	t3 := time.Date(2024, 1, 1, 0, 0, 0, 0, Loc)
	got = BetweenMonth(t1, t3, true)
	if got != 12 {
		t.Errorf("BetweenMonth(true) = %d, want 12", got)
	}
}

func TestBetweenYear(t *testing.T) {
	t1 := time.Date(2023, 1, 15, 12, 0, 0, 0, Loc)
	t2 := time.Date(2025, 3, 20, 12, 0, 0, 0, Loc)

	// 测试不重置
	got := BetweenYear(t1, t2, false)
	if got != 2 {
		t.Errorf("BetweenYear(false) = %d, want 2", got)
	}

	// 测试重置
	got = BetweenYear(t1, t2, true)
	if got != 2 {
		t.Errorf("BetweenYear(true) = %d, want 2", got)
	}

	// 测试同一年
	t3 := time.Date(2023, 12, 31, 23, 59, 59, 0, Loc)
	got = BetweenYear(t1, t3, true)
	if got != 0 {
		t.Errorf("BetweenYear(true) = %d, want 0", got)
	}
}
