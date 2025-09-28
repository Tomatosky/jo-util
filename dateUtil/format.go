package dateUtil

import (
	"runtime/debug"
	"time"
)

// FormatType 定义时间格式类型
type FormatType int

// 时间格式化样式枚举
const (
	Format_Custom FormatType = iota // 自定义格式
	// 2位年份格式
	Format_YY_MM_DD_HHmmss_SSS // 06-01-02 15:04:05.000
	Format_YY_MM_DD_HHmmss     // 06-01-02 15:04:05
	Format_YY_MM_DD_HHmm       // 06-01-02 15:04
	Format_YY_MM_DD            // 06-01-02
	Format_YY_MM               // 06-01
	Format_MM_DD               // 01-02
	// 2位年份紧凑格式
	Format_YYMMDDHHmmss // 060102150405
	Format_YYMMDD       // 060102
	Format_YYMM         // 0601
	Format_MMDD         // 0102
	Format_YY           // 06
	Format_MM           // 01
	Format_DD           // 02
	// 4位年份格式
	Format_YYYY_MM_DD_HHmmss_SSS // 2006-01-02 15:04:05.000
	Format_YYYY_MM_DD_HHmmss     // 2006-01-02 15:04:05
	Format_YYYY_MM_DD_HHmm       // 2006-01-02 15:04
	Format_YYYY_MM_DD            // 2006-01-02
	Format_YYYY_MM               // 2006-01
	// 4位年份紧凑格式
	Format_YYYYMMDDHHmmss // 20060102150405
	Format_YYYYMMDD       // 20060102
	Format_YYYYMM         // 200601
	Format_YYYY           // 2006
)

// FormatOpt 格式化选项
type FormatOpt struct {
	CustomTime int64  // 自定义时间(毫秒时间戳)，0表示使用当前时间
	CustomType string // 自定义格式化的样式(使用Go时间格式)
	TimeZone   string // 时区，默认"Asia/Shanghai"
}

// 默认时区
const defaultTimeZone = "Asia/Shanghai"

// Format 时间格式化方法
func Format(formatType FormatType, opt ...FormatOpt) string {
	var t time.Time
	option := FormatOpt{}
	if len(opt) > 0 {
		option = opt[0]
	}

	// 设置时间
	if option.CustomTime > 0 {
		t = time.Unix(0, option.CustomTime*int64(time.Millisecond))
	} else {
		t = time.Now()
	}
	// 设置时区
	timeZone := option.TimeZone
	if timeZone == "" {
		timeZone = defaultTimeZone
	}
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		debug.PrintStack()
		panic(err)
	}
	t = t.In(loc)

	// 根据格式类型返回格式化后的字符串
	switch formatType {
	case Format_Custom:
		if option.CustomType == "" {
			debug.PrintStack()
			panic("custom type cannot be empty")
		}
		return t.Format(option.CustomType)
	// 2位年份格式
	case Format_YY_MM_DD_HHmmss_SSS:
		return t.Format("06-01-02 15:04:05.000")
	case Format_YY_MM_DD_HHmmss:
		return t.Format("06-01-02 15:04:05")
	case Format_YY_MM_DD_HHmm:
		return t.Format("06-01-02 15:04")
	case Format_YY_MM_DD:
		return t.Format("06-01-02")
	case Format_YY_MM:
		return t.Format("06-01")
	case Format_MM_DD:
		return t.Format("01-02")
	// 2位年份紧凑格式
	case Format_YYMMDDHHmmss:
		return t.Format("060102150405")
	case Format_YYMMDD:
		return t.Format("060102")
	case Format_YYMM:
		return t.Format("0601")
	case Format_MMDD:
		return t.Format("0102")
	case Format_YY:
		return t.Format("06")
	case Format_MM:
		return t.Format("01")
	case Format_DD:
		return t.Format("02")
	// 4位年份格式
	case Format_YYYY_MM_DD_HHmmss_SSS:
		return t.Format("2006-01-02 15:04:05.000")
	case Format_YYYY_MM_DD_HHmmss:
		return t.Format("2006-01-02 15:04:05")
	case Format_YYYY_MM_DD_HHmm:
		return t.Format("2006-01-02 15:04")
	case Format_YYYY_MM_DD:
		return t.Format("2006-01-02")
	case Format_YYYY_MM:
		return t.Format("2006-01")
	// 4位年份紧凑格式
	case Format_YYYYMMDDHHmmss:
		return t.Format("20060102150405")
	case Format_YYYYMMDD:
		return t.Format("20060102")
	case Format_YYYYMM:
		return t.Format("200601")
	case Format_YYYY:
		return t.Format("2006")
	default:
		debug.PrintStack()
		panic("invalid format type")
	}
}
