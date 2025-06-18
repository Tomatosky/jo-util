package dateUtil

import (
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	// 固定时间戳用于测试(2022-01-02 17:04:05.123 CST)
	fixedTime := int64(1641114245123)

	tests := []struct {
		name      string
		format    FormatType
		opt       FormatOpt
		want      string
		expectErr bool
	}{
		// 测试自定义格式
		{
			name:   "Custom format",
			format: Format_Custom,
			opt: FormatOpt{
				CustomTime: fixedTime,
				CustomType: "2006/01/02 03:04:05 PM",
			},
			want: "2022/01/02 05:04:05 PM",
		},
		{
			name:      "Custom format empty",
			format:    Format_Custom,
			opt:       FormatOpt{},
			expectErr: true,
		},

		// 测试2位年份格式
		{
			name:   "YY_MM_DD_HHmmss_SSS",
			format: Format_YY_MM_DD_HHmmss_SSS,
			opt:    FormatOpt{CustomTime: fixedTime},
			want:   "22-01-02 17:04:05.123",
		},
		{
			name:   "YY_MM_DD_HHmmss",
			format: Format_YY_MM_DD_HHmmss,
			opt:    FormatOpt{CustomTime: fixedTime},
			want:   "22-01-02 17:04:05",
		},
		{
			name:   "YY_MM_DD_HHmm",
			format: Format_YY_MM_DD_HHmm,
			opt:    FormatOpt{CustomTime: fixedTime},
			want:   "22-01-02 17:04",
		},
		{
			name:   "YY_MM_DD",
			format: Format_YY_MM_DD,
			opt:    FormatOpt{CustomTime: fixedTime},
			want:   "22-01-02",
		},
		{
			name:   "YY_MM",
			format: Format_YY_MM,
			opt:    FormatOpt{CustomTime: fixedTime},
			want:   "22-01",
		},
		{
			name:   "MM_DD",
			format: Format_MM_DD,
			opt:    FormatOpt{CustomTime: fixedTime},
			want:   "01-02",
		},

		// 测试2位年份紧凑格式
		{
			name:   "YYMMDDHHmmss",
			format: Format_YYMMDDHHmmss,
			opt:    FormatOpt{CustomTime: fixedTime},
			want:   "220102170405",
		},
		{
			name:   "YYMMDD",
			format: Format_YYMMDD,
			opt:    FormatOpt{CustomTime: fixedTime},
			want:   "220102",
		},
		{
			name:   "YYMM",
			format: Format_YYMM,
			opt:    FormatOpt{CustomTime: fixedTime},
			want:   "2201",
		},
		{
			name:   "MMDD",
			format: Format_MMDD,
			opt:    FormatOpt{CustomTime: fixedTime},
			want:   "0102",
		},
		{
			name:   "YY",
			format: Format_YY,
			opt:    FormatOpt{CustomTime: fixedTime},
			want:   "22",
		},
		{
			name:   "MM",
			format: Format_MM,
			opt:    FormatOpt{CustomTime: fixedTime},
			want:   "01",
		},
		{
			name:   "DD",
			format: Format_DD,
			opt:    FormatOpt{CustomTime: fixedTime},
			want:   "02",
		},

		// 测试4位年份格式
		{
			name:   "YYYY_MM_DD_HHmmss_SSS",
			format: Format_YYYY_MM_DD_HHmmss_SSS,
			opt:    FormatOpt{CustomTime: fixedTime},
			want:   "2022-01-02 17:04:05.123",
		},
		{
			name:   "YYYY_MM_DD_HHmmss",
			format: Format_YYYY_MM_DD_HHmmss,
			opt:    FormatOpt{CustomTime: fixedTime},
			want:   "2022-01-02 17:04:05",
		},
		{
			name:   "YYYY_MM_DD_HHmm",
			format: Format_YYYY_MM_DD_HHmm,
			opt:    FormatOpt{CustomTime: fixedTime},
			want:   "2022-01-02 17:04",
		},
		{
			name:   "YYYY_MM_DD",
			format: Format_YYYY_MM_DD,
			opt:    FormatOpt{CustomTime: fixedTime},
			want:   "2022-01-02",
		},
		{
			name:   "YYYY_MM",
			format: Format_YYYY_MM,
			opt:    FormatOpt{CustomTime: fixedTime},
			want:   "2022-01",
		},

		// 测试4位年份紧凑格式
		{
			name:   "YYYYMMDDHHmmss",
			format: Format_YYYYMMDDHHmmss,
			opt:    FormatOpt{CustomTime: fixedTime},
			want:   "20220102170405",
		},
		{
			name:   "YYYYMMDD",
			format: Format_YYYYMMDD,
			opt:    FormatOpt{CustomTime: fixedTime},
			want:   "20220102",
		},
		{
			name:   "YYYYMM",
			format: Format_YYYYMM,
			opt:    FormatOpt{CustomTime: fixedTime},
			want:   "202201",
		},
		{
			name:   "YYYY",
			format: Format_YYYY,
			opt:    FormatOpt{CustomTime: fixedTime},
			want:   "2022",
		},

		// 测试时区
		{
			name:   "Timezone UTC",
			format: Format_YYYY_MM_DD_HHmmss,
			opt: FormatOpt{
				CustomTime: fixedTime,
				TimeZone:   "UTC",
			},
			want: "2022-01-02 09:04:05",
		},
		{
			name:   "Timezone America/New_York",
			format: Format_YYYY_MM_DD_HHmmss,
			opt: FormatOpt{
				CustomTime: fixedTime,
				TimeZone:   "America/New_York",
			},
			want: "2022-01-02 04:04:05",
		},

		// 测试当前时间
		{
			name:   "Current time",
			format: Format_YYYY_MM_DD,
			opt:    FormatOpt{},
			want:   time.Now().Format("2006-01-02"),
		},

		// 测试无效格式类型
		{
			name:      "Invalid format type",
			format:    FormatType(999),
			opt:       FormatOpt{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.expectErr {
						t.Errorf("Format() panicked unexpectedly: %v", r)
					}
				} else if tt.expectErr {
					t.Error("Format() should have panicked but didn't")
				}
			}()

			got := Format(tt.format, tt.opt)
			if !tt.expectErr && got != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormat_InvalidTimezone(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for invalid timezone, but got none")
		}
	}()

	Format(Format_YYYY_MM_DD, FormatOpt{
		TimeZone: "Invalid/Timezone",
	})
}
