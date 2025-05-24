package dateUtil

import (
	"testing"
	"time"
)

func TestGetTime(t *testing.T) {
	// 1640995200 是 UTC 时间 2022-01-01 00:00:00
	// 在 Asia/Shanghai (UTC+8) 应该是 2022-01-01 08:00:00
	timestamp := int64(1640995200)
	expected := time.Date(2022, 1, 1, 8, 0, 0, 0, loc)
	result := GetTime(timestamp)

	if !result.Equal(expected) {
		t.Errorf("GetTime() = %v, want %v", result, expected)
	}
}

func TestBeginOfDay(t *testing.T) {
	tests := []struct {
		name      string
		timestamp int64
		want      int64
	}{
		// 1641024000 是 UTC 时间 2022-01-01 08:00:00 (Asia/Shanghai 16:00:00)
		// 当天开始时间应该是 Asia/Shanghai 00:00:00 (UTC 前一天 16:00:00)
		{"noon", 1641024000, 1640966400},       // 2022-01-01 16:00:00 CST -> 2022-01-01 00:00:00 CST
		{"midnight", 1640966400, 1640966400},   // 2022-01-01 00:00:00 CST -> same
		{"end of day", 1641052799, 1640966400}, // 2022-01-01 23:59:59 CST -> 2022-01-01 00:00:00 CST
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BeginOfDay(tt.timestamp); got != tt.want {
				t.Errorf("BeginOfDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBeginOfWeek(t *testing.T) {
	tests := []struct {
		name      string
		timestamp int64
		want      int64
	}{
		// 2022-01-03 (Monday) 00:00:00 CST = 2022-01-02 16:00:00 UTC (1641168000)
		{"Monday", 1641168000, 1641139200},    // 2022-01-03 (Monday) -> same
		{"Wednesday", 1641340800, 1641139200}, // 2022-01-05 -> 2022-01-03
		{"Sunday", 1641686400, 1641139200},    // 2022-01-09 (Sunday) -> 2022-01-03
		{"next week", 1641772800, 1641744000}, // 2022-01-10 (next Monday) -> same
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BeginOfWeek(tt.timestamp); got != tt.want {
				t.Errorf("BeginOfWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBeginOfMonth(t *testing.T) {
	tests := []struct {
		name      string
		timestamp int64
		want      int64
	}{
		{"first day", 1640966400, 1640966400},    // 2022-01-01 00:00:00 CST -> same
		{"middle month", 1643644800, 1643644800}, // 2022-02-01 00:00:00 CST -> same
		{"end month", 1646063999, 1643644800},    // 2022-02-28 23:59:59 CST -> 2022-02-01 00:00:00 CST
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BeginOfMonth(tt.timestamp); got != tt.want {
				t.Errorf("BeginOfMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBetweenDay(t *testing.T) {
	tests := []struct {
		name string
		t1   int64
		t2   int64
		want int
	}{
		{"same day", 1640966400, 1640995200, 0},  // 2022-01-01 00:00:00 CST and 2022-01-01 08:00:00 CST
		{"one day", 1640966400, 1641052800, 1},   // 2022-01-01 and 2022-01-02
		{"four days", 1640966400, 1641312000, 4}, // 2022-01-01 and 2022-01-05
		{"negative", 1641052800, 1640966400, 1},  // should return absolute value
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BetweenDay(tt.t1, tt.t2); got != tt.want {
				t.Errorf("BetweenDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndOfDay(t *testing.T) {
	tests := []struct {
		name      string
		timestamp int64
		want      int64
	}{
		{"noon", 1640995200, 1641052799},     // 2022-01-01 08:00:00 UTC -> 2022-01-01 23:59:59 CST
		{"midnight", 1640966400, 1641052799}, // 2022-01-01 00:00:00 CST -> 2022-01-01 23:59:59 CST
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EndOfDay(tt.timestamp); got != tt.want {
				t.Errorf("EndOfDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndOfWeek(t *testing.T) {
	tests := []struct {
		name      string
		timestamp int64
		want      int64
	}{
		{"Monday", 1641168000, 1641743999},    // 2022-01-03 (Monday) -> 2022-01-09 23:59:59 CST
		{"Wednesday", 1641340800, 1641743999}, // 2022-01-05 -> 2022-01-09 23:59:59 CST
		{"Sunday", 1641686400, 1641743999},    // 2022-01-09 (Sunday) -> same day 23:59:59 CST
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EndOfWeek(tt.timestamp); got != tt.want {
				t.Errorf("EndOfWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndOfMonth(t *testing.T) {
	tests := []struct {
		name      string
		timestamp int64
		want      int64
	}{
		{"January", 1640966400, 1643644799},       // 2022-01-01 -> 2022-01-31 23:59:59 CST
		{"February", 1643644800, 1646063999},      // 2022-02-01 -> 2022-02-28 23:59:59 CST
		{"February leap", 1582992000, 1585670399}, // 2020-02-01 -> 2020-02-29 23:59:59 CST
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EndOfMonth(tt.timestamp); got != tt.want {
				t.Errorf("EndOfMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	timestamp := int64(1640995200) // UTC 2022-01-01 00:00:00 = CST 2022-01-01 08:00:00
	want := "2022-01-01 08:00:00"
	if got := Format(timestamp); got != want {
		t.Errorf("Format() = %v, want %v", got, want)
	}
}

func TestIsSameDay(t *testing.T) {
	tests := []struct {
		name string
		t1   int64
		t2   int64
		want bool
	}{
		{"same time", 1640995200, 1640995200, true},
		{"same day", 1640995200, 1641009600, true},       // 2022-01-01 08:00:00 and 2022-01-01 12:00:00 CST
		{"different day", 1640995200, 1641081600, false}, // 2022-01-01 and 2022-01-02
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSameDay(tt.t1, tt.t2); got != tt.want {
				t.Errorf("IsSameDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatToStr(t *testing.T) {
	testTime := time.Date(2022, 1, 1, 12, 30, 45, 0, loc)

	tests := []struct {
		name   string
		t      time.Time
		format string
		want   string
	}{
		{"default format", testTime, "yyyy-mm-dd hh:mm:ss", "2022-01-01 12:30:45"},
		{"date only", testTime, "yyyy-mm-dd", "2022-01-01"},
		{"time only", testTime, "hh:mm:ss", "12:30:45"},
		{"different format", testTime, "dd/mm/yy hh:mm:ss", "01/01/22 12:30:45"},
		{"invalid format", testTime, "invalid", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatToStr(tt.t, tt.format); got != tt.want {
				t.Errorf("FormatToStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
