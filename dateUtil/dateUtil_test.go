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
		{"Zero timestamp", 0, time.Unix(0, 0).In(loc)},
		{"Current timestamp", time.Now().Unix(), time.Unix(time.Now().Unix(), 0).In(loc)},
		{"Future timestamp", 1893456000, time.Unix(1893456000, 0).In(loc)}, // 2030-01-01 00:00:00
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
			input:    time.Date(2023, 5, 15, 12, 30, 45, 0, loc),
			expected: time.Date(2023, 5, 15, 0, 0, 0, 0, loc),
		},
		{
			name:     "普通日期-午夜时间",
			input:    time.Date(2023, 5, 15, 0, 0, 0, 0, loc),
			expected: time.Date(2023, 5, 15, 0, 0, 0, 0, loc),
		},
		{
			name:     "普通日期-接近午夜时间",
			input:    time.Date(2023, 5, 15, 23, 59, 59, 999999, loc),
			expected: time.Date(2023, 5, 15, 0, 0, 0, 0, loc),
		},
		{
			name:     "闰年2月29日",
			input:    time.Date(2020, 2, 29, 15, 30, 0, 0, loc),
			expected: time.Date(2020, 2, 29, 0, 0, 0, 0, loc),
		},
		{
			name:     "跨年时间",
			input:    time.Date(2023, 12, 31, 23, 59, 59, 0, loc),
			expected: time.Date(2023, 12, 31, 0, 0, 0, 0, loc),
		},
		{
			name:     "夏令时转换时间",
			input:    time.Date(2023, 3, 12, 2, 30, 0, 0, loc), // 夏令时开始
			expected: time.Date(2023, 3, 12, 0, 0, 0, 0, loc),
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
			input:    time.Date(2023, 5, 15, 12, 30, 45, 0, loc), // 2023-05-15 周一
			expected: time.Date(2023, 5, 15, 0, 0, 0, 0, loc),    // 同一天零点
		},
		{
			name:     "周二",
			input:    time.Date(2023, 5, 16, 8, 0, 0, 0, loc), // 2023-05-16 周二
			expected: time.Date(2023, 5, 15, 0, 0, 0, 0, loc), // 本周一零点
		},
		{
			name:     "周日",
			input:    time.Date(2023, 5, 21, 23, 59, 59, 0, loc), // 2023-05-21 周日
			expected: time.Date(2023, 5, 15, 0, 0, 0, 0, loc),    // 本周一零点
		},
		{
			name:     "跨月情况-月初是周一",
			input:    time.Date(2023, 6, 4, 15, 0, 0, 0, loc), // 2023-06-04 周日
			expected: time.Date(2023, 5, 29, 0, 0, 0, 0, loc), // 上周一(5月29日)零点
		},
		{
			name:     "跨年情况-年初是周一",
			input:    time.Date(2023, 1, 1, 0, 0, 1, 0, loc),   // 2023-01-01 周日
			expected: time.Date(2022, 12, 26, 0, 0, 0, 0, loc), // 上周一(2022-12-26)零点
		},
		{
			name:     "闰年2月",
			input:    time.Date(2020, 2, 29, 12, 0, 0, 0, loc), // 2020-02-29 周六
			expected: time.Date(2020, 2, 24, 0, 0, 0, 0, loc),  // 本周一零点
		},
		{
			name:     "夏令时转换周",
			input:    time.Date(2023, 3, 12, 2, 30, 0, 0, loc), // 2023-03-12 周日(夏令时开始)
			expected: time.Date(2023, 3, 6, 0, 0, 0, 0, loc),   // 本周一零点
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
			input:    time.Date(2023, 5, 1, 0, 0, 0, 0, loc), // 2023-05-01
			expected: time.Date(2023, 5, 1, 0, 0, 0, 0, loc), // 同一天零点
		},
		{
			name:     "月中某天",
			input:    time.Date(2023, 5, 15, 12, 30, 45, 0, loc), // 2023-05-15
			expected: time.Date(2023, 5, 1, 0, 0, 0, 0, loc),     // 当月第一天零点
		},
		{
			name:     "当月最后一天",
			input:    time.Date(2023, 5, 31, 23, 59, 59, 0, loc), // 2023-05-31
			expected: time.Date(2023, 5, 1, 0, 0, 0, 0, loc),     // 当月第一天零点
		},
		{
			name:     "闰年2月",
			input:    time.Date(2020, 2, 29, 12, 0, 0, 0, loc), // 2020-02-29
			expected: time.Date(2020, 2, 1, 0, 0, 0, 0, loc),   // 当月第一天零点
		},
		{
			name:     "非闰年2月",
			input:    time.Date(2023, 2, 28, 23, 59, 59, 0, loc), // 2023-02-28
			expected: time.Date(2023, 2, 1, 0, 0, 0, 0, loc),     // 当月第一天零点
		},
		{
			name:     "跨年12月",
			input:    time.Date(2023, 12, 31, 0, 0, 1, 0, loc), // 2023-12-31
			expected: time.Date(2023, 12, 1, 0, 0, 0, 0, loc),  // 当月第一天零点
		},
		{
			name:     "跨年1月",
			input:    time.Date(2023, 1, 1, 0, 0, 0, 0, loc), // 2023-01-01
			expected: time.Date(2023, 1, 1, 0, 0, 0, 0, loc), // 同一天零点
		},
		{
			name:     "夏令时转换月",
			input:    time.Date(2023, 3, 15, 2, 30, 0, 0, loc), // 2023-03-15
			expected: time.Date(2023, 3, 1, 0, 0, 0, 0, loc),   // 当月第一天零点
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

func TestBetweenDay(t *testing.T) {
	tests := []struct {
		name     string
		t1       time.Time
		t2       time.Time
		expected int
	}{
		{
			name:     "同一天",
			t1:       time.Date(2023, 1, 1, 12, 0, 0, 0, loc),
			t2:       time.Date(2023, 1, 1, 23, 59, 59, 0, loc),
			expected: 0,
		},
		{
			name:     "相差1天-同月",
			t1:       time.Date(2023, 1, 1, 0, 0, 0, 0, loc),
			t2:       time.Date(2023, 1, 2, 0, 0, 0, 0, loc),
			expected: 1,
		},
		{
			name:     "相差1天-跨月",
			t1:       time.Date(2023, 1, 31, 0, 0, 0, 0, loc),
			t2:       time.Date(2023, 2, 1, 0, 0, 0, 0, loc),
			expected: 1,
		},
		{
			name:     "相差1天-跨年",
			t1:       time.Date(2022, 12, 31, 0, 0, 0, 0, loc),
			t2:       time.Date(2023, 1, 1, 0, 0, 0, 0, loc),
			expected: 1,
		},
		{
			name:     "相差7天",
			t1:       time.Date(2023, 1, 1, 0, 0, 0, 0, loc),
			t2:       time.Date(2023, 1, 8, 0, 0, 0, 0, loc),
			expected: 7,
		},
		{
			name:     "相差30天-非闰年2月",
			t1:       time.Date(2023, 1, 1, 0, 0, 0, 0, loc),
			t2:       time.Date(2023, 1, 31, 0, 0, 0, 0, loc),
			expected: 30,
		},
		{
			name:     "相差31天",
			t1:       time.Date(2023, 1, 1, 0, 0, 0, 0, loc),
			t2:       time.Date(2023, 2, 1, 0, 0, 0, 0, loc),
			expected: 31,
		},
		{
			name:     "相差365天-非闰年",
			t1:       time.Date(2023, 1, 1, 0, 0, 0, 0, loc),
			t2:       time.Date(2024, 1, 1, 0, 0, 0, 0, loc),
			expected: 365,
		},
		{
			name:     "相差366天-闰年",
			t1:       time.Date(2020, 1, 1, 0, 0, 0, 0, loc),
			t2:       time.Date(2021, 1, 1, 0, 0, 0, 0, loc),
			expected: 366,
		},
		{
			name:     "时间顺序相反",
			t1:       time.Date(2023, 1, 8, 0, 0, 0, 0, loc),
			t2:       time.Date(2023, 1, 1, 0, 0, 0, 0, loc),
			expected: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := BetweenDay(tt.t1, tt.t2)
			if actual != tt.expected {
				t.Errorf("BetweenDay(%v, %v) = %d, want %d", tt.t1, tt.t2, actual, tt.expected)
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
			input:    time.Date(2023, 1, 1, 12, 30, 15, 123456, loc),
			expected: time.Date(2023, 1, 1, 23, 59, 59, 0, loc),
		},
		{
			name:     "午夜零点",
			input:    time.Date(2023, 1, 1, 0, 0, 0, 0, loc),
			expected: time.Date(2023, 1, 1, 23, 59, 59, 0, loc),
		},
		{
			name:     "23:59:59",
			input:    time.Date(2023, 1, 1, 23, 59, 59, 999999, loc),
			expected: time.Date(2023, 1, 1, 23, 59, 59, 0, loc),
		},
		{
			name:     "闰年2月29日",
			input:    time.Date(2020, 2, 29, 15, 30, 0, 0, loc),
			expected: time.Date(2020, 2, 29, 23, 59, 59, 0, loc),
		},
		{
			name:     "非闰年2月28日",
			input:    time.Date(2023, 2, 28, 15, 30, 0, 0, loc),
			expected: time.Date(2023, 2, 28, 23, 59, 59, 0, loc),
		},
		{
			name:     "12月31日",
			input:    time.Date(2023, 12, 31, 15, 30, 0, 0, loc),
			expected: time.Date(2023, 12, 31, 23, 59, 59, 0, loc),
		},
		{
			name:     "1月1日",
			input:    time.Date(2023, 1, 1, 15, 30, 0, 0, loc),
			expected: time.Date(2023, 1, 1, 23, 59, 59, 0, loc),
		},
		{
			name:     "带纳秒的时间",
			input:    time.Date(2023, 1, 1, 12, 30, 15, 999999999, loc),
			expected: time.Date(2023, 1, 1, 23, 59, 59, 0, loc),
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
			input:    time.Date(2023, 6, 5, 12, 0, 0, 0, loc), // 2023-06-05 周一
			expected: time.Date(2023, 6, 11, 23, 59, 59, 0, loc),
		},
		{
			name:     "周二",
			input:    time.Date(2023, 6, 6, 12, 0, 0, 0, loc), // 2023-06-06 周二
			expected: time.Date(2023, 6, 11, 23, 59, 59, 0, loc),
		},
		{
			name:     "周三",
			input:    time.Date(2023, 6, 7, 12, 0, 0, 0, loc), // 2023-06-07 周三
			expected: time.Date(2023, 6, 11, 23, 59, 59, 0, loc),
		},
		{
			name:     "周四",
			input:    time.Date(2023, 6, 8, 12, 0, 0, 0, loc), // 2023-06-08 周四
			expected: time.Date(2023, 6, 11, 23, 59, 59, 0, loc),
		},
		{
			name:     "周五",
			input:    time.Date(2023, 6, 9, 12, 0, 0, 0, loc), // 2023-06-09 周五
			expected: time.Date(2023, 6, 11, 23, 59, 59, 0, loc),
		},
		{
			name:     "周六",
			input:    time.Date(2023, 6, 10, 12, 0, 0, 0, loc), // 2023-06-10 周六
			expected: time.Date(2023, 6, 11, 23, 59, 59, 0, loc),
		},
		{
			name:     "周日",
			input:    time.Date(2023, 6, 11, 12, 0, 0, 0, loc), // 2023-06-11 周日
			expected: time.Date(2023, 6, 11, 23, 59, 59, 0, loc),
		},
		{
			name:     "跨月周日",
			input:    time.Date(2023, 5, 29, 12, 0, 0, 0, loc),  // 2023-05-29 周一
			expected: time.Date(2023, 6, 4, 23, 59, 59, 0, loc), // 2023-06-04 周日
		},
		{
			name:     "跨年周日",
			input:    time.Date(2022, 12, 26, 12, 0, 0, 0, loc), // 2022-12-26 周一
			expected: time.Date(2023, 1, 1, 23, 59, 59, 0, loc), // 2023-01-01 周日
		},
		{
			name:     "闰年2月",
			input:    time.Date(2020, 2, 24, 12, 0, 0, 0, loc),  // 2020-02-24 周一
			expected: time.Date(2020, 3, 1, 23, 59, 59, 0, loc), // 2020-03-01 周日
		},
		{
			name:     "周日23:59:59",
			input:    time.Date(2023, 6, 11, 23, 59, 59, 999, loc),
			expected: time.Date(2023, 6, 11, 23, 59, 59, 0, loc),
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
			input:    time.Date(2023, 1, 15, 12, 30, 15, 123456, loc),
			expected: time.Date(2023, 1, 31, 23, 59, 59, 0, loc),
		},
		{
			name:     "2月非闰年28天",
			input:    time.Date(2023, 2, 10, 0, 0, 0, 0, loc),
			expected: time.Date(2023, 2, 28, 23, 59, 59, 0, loc),
		},
		{
			name:     "2月闰年29天",
			input:    time.Date(2020, 2, 15, 12, 0, 0, 0, loc),
			expected: time.Date(2020, 2, 29, 23, 59, 59, 0, loc),
		},
		{
			name:     "4月30天",
			input:    time.Date(2023, 4, 1, 23, 59, 59, 999999, loc),
			expected: time.Date(2023, 4, 30, 23, 59, 59, 0, loc),
		},
		{
			name:     "12月31天",
			input:    time.Date(2023, 12, 31, 0, 0, 0, 0, loc),
			expected: time.Date(2023, 12, 31, 23, 59, 59, 0, loc),
		},
		{
			name:     "月末最后一天",
			input:    time.Date(2023, 1, 31, 23, 59, 59, 999999, loc),
			expected: time.Date(2023, 1, 31, 23, 59, 59, 0, loc),
		},
		{
			name:     "跨年12月",
			input:    time.Date(2022, 12, 1, 0, 0, 0, 0, loc),
			expected: time.Date(2022, 12, 31, 23, 59, 59, 0, loc),
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
			t1:       time.Date(2023, 5, 15, 12, 30, 45, 0, loc),
			t2:       time.Date(2023, 5, 15, 12, 30, 45, 0, loc),
			expected: true,
		},
		{
			name:     "同一天不同时间",
			t1:       time.Date(2023, 5, 15, 0, 0, 0, 0, loc),
			t2:       time.Date(2023, 5, 15, 23, 59, 59, 0, loc),
			expected: true,
		},
		{
			name:     "不同天",
			t1:       time.Date(2023, 5, 15, 23, 59, 59, 0, loc),
			t2:       time.Date(2023, 5, 16, 0, 0, 0, 0, loc),
			expected: false,
		},
		{
			name:     "不同月",
			t1:       time.Date(2023, 5, 31, 23, 59, 59, 0, loc),
			t2:       time.Date(2023, 6, 1, 0, 0, 0, 0, loc),
			expected: false,
		},
		{
			name:     "不同年",
			t1:       time.Date(2022, 12, 31, 23, 59, 59, 0, loc),
			t2:       time.Date(2023, 1, 1, 0, 0, 0, 0, loc),
			expected: false,
		},
		{
			name:     "闰年2月29日",
			t1:       time.Date(2020, 2, 29, 0, 0, 0, 0, loc),
			t2:       time.Date(2020, 2, 29, 23, 59, 59, 0, loc),
			expected: true,
		},
		{
			name:     "边界时间-午夜前",
			t1:       time.Date(2023, 5, 15, 23, 59, 59, 999999999, loc),
			t2:       time.Date(2023, 5, 15, 0, 0, 0, 0, loc),
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
			t1:       time.Date(2023, 5, 15, 12, 30, 45, 0, loc), // 周一
			t2:       time.Date(2023, 5, 15, 12, 30, 45, 0, loc),
			expected: true,
		},
		{
			name:     "同一周不同天",
			t1:       time.Date(2023, 5, 15, 0, 0, 0, 0, loc),    // 周一
			t2:       time.Date(2023, 5, 21, 23, 59, 59, 0, loc), // 周日
			expected: true,
		},
		{
			name:     "跨周边界-周日和周一",
			t1:       time.Date(2023, 5, 21, 23, 59, 59, 0, loc), // 周日
			t2:       time.Date(2023, 5, 22, 0, 0, 0, 0, loc),    // 周一
			expected: false,
		},
		{
			name:     "不同周",
			t1:       time.Date(2023, 5, 14, 23, 59, 59, 0, loc), // 第19周周日
			t2:       time.Date(2023, 5, 15, 0, 0, 0, 0, loc),    // 第20周周一
			expected: false,
		},
		{
			name:     "跨年周-第52周和第1周",
			t1:       time.Date(2022, 12, 31, 23, 59, 59, 0, loc), // 2022年第52周周六
			t2:       time.Date(2023, 1, 1, 0, 0, 0, 0, loc),      // 2023年第52周周日
			expected: true,                                        // 注意：ISO周计算中，2023-01-01仍属于2022年第52周
		},
		{
			name:     "跨年周-不同年",
			t1:       time.Date(2021, 12, 31, 23, 59, 59, 0, loc), // 2021年第52周周五
			t2:       time.Date(2022, 1, 1, 0, 0, 0, 0, loc),      // 2022年第52周周六
			expected: true,                                        // 2021-12-31属于2021年第52周，2022-01-01属于2021年第52周
		},
		{
			name:     "闰年周",
			t1:       time.Date(2020, 2, 24, 0, 0, 0, 0, loc),   // 2020年第9周周一
			t2:       time.Date(2020, 3, 1, 23, 59, 59, 0, loc), // 2020年第9周周日
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
