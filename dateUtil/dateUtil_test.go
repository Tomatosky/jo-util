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
			"Normal time",
			time.Date(2023, 5, 15, 14, 30, 15, 123, loc),
			time.Date(2023, 5, 15, 0, 0, 0, 0, loc),
		},
		{
			"Already at midnight",
			time.Date(2023, 5, 15, 0, 0, 0, 0, loc),
			time.Date(2023, 5, 15, 0, 0, 0, 0, loc),
		},
		{
			"Leap day",
			time.Date(2020, 2, 29, 23, 59, 59, 999, loc),
			time.Date(2020, 2, 29, 0, 0, 0, 0, loc),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BeginOfDay(tt.input)
			if !got.Equal(tt.expected) {
				t.Errorf("BeginOfDay() = %v, want %v", got, tt.expected)
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
			"Monday",
			time.Date(2023, 5, 15, 0, 0, 0, 0, loc), // 周一
			time.Date(2023, 5, 15, 0, 0, 0, 0, loc),
		},
		{
			"Wednesday",
			time.Date(2023, 5, 17, 12, 0, 0, 0, loc), // 周三
			time.Date(2023, 5, 15, 0, 0, 0, 0, loc),
		},
		{
			"Sunday",
			time.Date(2023, 5, 21, 23, 59, 59, 999, loc), // 周日
			time.Date(2023, 5, 15, 0, 0, 0, 0, loc),
		},
		{
			"Year boundary",
			time.Date(2023, 1, 1, 0, 0, 0, 0, loc), // 周日
			time.Date(2022, 12, 26, 0, 0, 0, 0, loc),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BeginOfWeek(tt.input)
			if !got.Equal(tt.expected) {
				t.Errorf("BeginOfWeek() = %v, want %v", got, tt.expected)
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
			"Normal date",
			time.Date(2023, 5, 15, 14, 30, 0, 0, loc),
			time.Date(2023, 5, 1, 0, 0, 0, 0, loc),
		},
		{
			"Already first day",
			time.Date(2023, 5, 1, 0, 0, 0, 0, loc),
			time.Date(2023, 5, 1, 0, 0, 0, 0, loc),
		},
		{
			"February leap year",
			time.Date(2020, 2, 29, 23, 59, 59, 999, loc),
			time.Date(2020, 2, 1, 0, 0, 0, 0, loc),
		},
		{
			"December to January",
			time.Date(2023, 12, 31, 0, 0, 0, 0, loc),
			time.Date(2023, 12, 1, 0, 0, 0, 0, loc),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BeginOfMonth(tt.input)
			if !got.Equal(tt.expected) {
				t.Errorf("BeginOfMonth() = %v, want %v", got, tt.expected)
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
			"Same day",
			time.Date(2023, 5, 15, 0, 0, 0, 0, loc),
			time.Date(2023, 5, 15, 23, 59, 59, 999, loc),
			0,
		},
		{
			"Consecutive days",
			time.Date(2023, 5, 15, 0, 0, 0, 0, loc),
			time.Date(2023, 5, 16, 0, 0, 0, 0, loc),
			1,
		},
		{
			"One week apart",
			time.Date(2023, 5, 15, 0, 0, 0, 0, loc),
			time.Date(2023, 5, 22, 0, 0, 0, 0, loc),
			7,
		},
		{
			"Month boundary",
			time.Date(2023, 4, 30, 23, 59, 59, 999, loc),
			time.Date(2023, 5, 1, 0, 0, 0, 0, loc),
			1,
		},
		{
			"Year boundary",
			time.Date(2022, 12, 31, 23, 59, 59, 999, loc),
			time.Date(2023, 1, 1, 0, 0, 0, 0, loc),
			1,
		},
		{
			"Leap year",
			time.Date(2020, 2, 28, 0, 0, 0, 0, loc),
			time.Date(2020, 3, 1, 0, 0, 0, 0, loc),
			2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BetweenDay(tt.t1, tt.t2)
			if got != tt.expected {
				t.Errorf("BetweenDay() = %v, want %v", got, tt.expected)
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
			"Normal time",
			time.Date(2023, 5, 15, 14, 30, 15, 123, loc),
			time.Date(2023, 5, 15, 23, 59, 59, 0, loc),
		},
		{
			"Already at end of day",
			time.Date(2023, 5, 15, 23, 59, 59, 0, loc),
			time.Date(2023, 5, 15, 23, 59, 59, 0, loc),
		},
		{
			"Leap day",
			time.Date(2020, 2, 29, 0, 0, 0, 0, loc),
			time.Date(2020, 2, 29, 23, 59, 59, 0, loc),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EndOfDay(tt.input)
			if !got.Equal(tt.expected) {
				t.Errorf("EndOfDay() = %v, want %v", got, tt.expected)
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
			"Monday",
			time.Date(2023, 5, 15, 0, 0, 0, 0, loc), // 周一
			time.Date(2023, 5, 21, 23, 59, 59, 0, loc),
		},
		{
			"Wednesday",
			time.Date(2023, 5, 17, 12, 0, 0, 0, loc), // 周三
			time.Date(2023, 5, 21, 23, 59, 59, 0, loc),
		},
		{
			"Sunday",
			time.Date(2023, 5, 21, 23, 59, 59, 999, loc), // 周日
			time.Date(2023, 5, 21, 23, 59, 59, 0, loc),
		},
		{
			"Year boundary",
			time.Date(2023, 1, 1, 0, 0, 0, 0, loc), // 周日
			time.Date(2023, 1, 1, 23, 59, 59, 0, loc),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EndOfWeek(tt.input)
			if !got.Equal(tt.expected) {
				t.Errorf("EndOfWeek() = %v, want %v", got, tt.expected)
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
			"Normal month",
			time.Date(2023, 5, 15, 14, 30, 0, 0, loc),
			time.Date(2023, 5, 31, 23, 59, 59, 0, loc),
		},
		{
			"Already last day",
			time.Date(2023, 5, 31, 23, 59, 59, 0, loc),
			time.Date(2023, 5, 31, 23, 59, 59, 0, loc),
		},
		{
			"February leap year",
			time.Date(2020, 2, 1, 0, 0, 0, 0, loc),
			time.Date(2020, 2, 29, 23, 59, 59, 0, loc),
		},
		{
			"February non-leap year",
			time.Date(2023, 2, 1, 0, 0, 0, 0, loc),
			time.Date(2023, 2, 28, 23, 59, 59, 0, loc),
		},
		{
			"December",
			time.Date(2023, 12, 1, 0, 0, 0, 0, loc),
			time.Date(2023, 12, 31, 23, 59, 59, 0, loc),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EndOfMonth(tt.input)
			if !got.Equal(tt.expected) {
				t.Errorf("EndOfMonth() = %v, want %v", got, tt.expected)
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
			"Same time",
			time.Date(2023, 5, 15, 14, 30, 0, 0, loc),
			time.Date(2023, 5, 15, 14, 30, 0, 0, loc),
			true,
		},
		{
			"Same day different time",
			time.Date(2023, 5, 15, 0, 0, 0, 0, loc),
			time.Date(2023, 5, 15, 23, 59, 59, 999, loc),
			true,
		},
		{
			"Different day same time",
			time.Date(2023, 5, 15, 14, 30, 0, 0, loc),
			time.Date(2023, 5, 16, 14, 30, 0, 0, loc),
			false,
		},
		{
			"Month boundary",
			time.Date(2023, 4, 30, 23, 59, 59, 999, loc),
			time.Date(2023, 5, 1, 0, 0, 0, 0, loc),
			false,
		},
		{
			"Year boundary",
			time.Date(2022, 12, 31, 23, 59, 59, 999, loc),
			time.Date(2023, 1, 1, 0, 0, 0, 0, loc),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSameDay(tt.t1, tt.t2)
			if got != tt.expected {
				t.Errorf("IsSameDay() = %v, want %v", got, tt.expected)
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
			"Same day",
			time.Date(2023, 5, 15, 0, 0, 0, 0, loc), // 周一
			time.Date(2023, 5, 15, 23, 59, 59, 999, loc),
			true,
		},
		{
			"Same week",
			time.Date(2023, 5, 15, 0, 0, 0, 0, loc),      // 周一
			time.Date(2023, 5, 21, 23, 59, 59, 999, loc), // 周日
			true,
		},
		{
			"Different week",
			time.Date(2023, 5, 15, 0, 0, 0, 0, loc), // 周一
			time.Date(2023, 5, 22, 0, 0, 0, 0, loc), // 下周一
			false,
		},
		{
			"Year boundary same week",
			time.Date(2022, 12, 31, 0, 0, 0, 0, loc), // 周六
			time.Date(2023, 1, 1, 0, 0, 0, 0, loc),   // 周日
			true,
		},
		{
			"Year boundary different week",
			time.Date(2022, 12, 25, 0, 0, 0, 0, loc), // 周日
			time.Date(2022, 12, 26, 0, 0, 0, 0, loc), // 周一
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSameWeek(tt.t1, tt.t2)
			if got != tt.expected {
				t.Errorf("IsSameWeek() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsSameMonth(t *testing.T) {
	tests := []struct {
		name     string
		t1       time.Time
		t2       time.Time
		expected bool
	}{
		{
			"Same day",
			time.Date(2023, 5, 15, 0, 0, 0, 0, loc),
			time.Date(2023, 5, 15, 23, 59, 59, 999, loc),
			true,
		},
		{
			"Same month different day",
			time.Date(2023, 5, 1, 0, 0, 0, 0, loc),
			time.Date(2023, 5, 31, 23, 59, 59, 999, loc),
			true,
		},
		{
			"Different month same day",
			time.Date(2023, 5, 15, 0, 0, 0, 0, loc),
			time.Date(2023, 6, 15, 0, 0, 0, 0, loc),
			false,
		},
		{
			"Year boundary same month",
			time.Date(2022, 12, 1, 0, 0, 0, 0, loc),
			time.Date(2022, 12, 31, 23, 59, 59, 999, loc),
			true,
		},
		{
			"Year boundary different month",
			time.Date(2022, 12, 31, 23, 59, 59, 999, loc),
			time.Date(2023, 1, 1, 0, 0, 0, 0, loc),
			false,
		},
		{
			"Leap year February",
			time.Date(2020, 2, 1, 0, 0, 0, 0, loc),
			time.Date(2020, 2, 29, 23, 59, 59, 999, loc),
			true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSameMonth(tt.t1, tt.t2)
			if got != tt.expected {
				t.Errorf("IsSameMonth() = %v, want %v", got, tt.expected)
			}
		})
	}
}
