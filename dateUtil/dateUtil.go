package dateUtil

import (
	"fmt"
	"math"
	"strings"
	"time"
)

var loc *time.Location
var timeFormat map[string]string

func init() {
	loc, _ = time.LoadLocation("Asia/Shanghai")
	timeFormat = map[string]string{
		"yyyy-mm-dd hh:mm:ss": "2006-01-02 15:04:05",
		"yyyy-mm-dd hh:mm":    "2006-01-02 15:04",
		"yyyy-mm-dd hh":       "2006-01-02 15",
		"yyyy-mm-dd":          "2006-01-02",
		"yyyy-mm":             "2006-01",
		"mm-dd":               "01-02",
		"dd-mm-yy hh:mm:ss":   "02-01-06 15:04:05",
		"yyyy/mm/dd hh:mm:ss": "2006/01/02 15:04:05",
		"yyyy/mm/dd hh:mm":    "2006/01/02 15:04",
		"yyyy/mm/dd hh":       "2006/01/02 15",
		"yyyy/mm/dd":          "2006/01/02",
		"yyyy/mm":             "2006/01",
		"mm/dd":               "01/02",
		"dd/mm/yy hh:mm:ss":   "02/01/06 15:04:05",
		"yyyymmdd":            "20060102",
		"mmddyy":              "010206",
		"yyyy":                "2006",
		"yy":                  "06",
		"mm":                  "01",
		"hh:mm:ss":            "15:04:05",
		"hh:mm":               "15:04",
		"mm:ss":               "04:05",
	}
}

func GetTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0).In(loc)
}

// BeginOfDay 获取某天的起始时间戳(秒时间戳)
func BeginOfDay(timestamp int64) int64 {
	t := GetTime(timestamp)
	year, month, day := t.Date()
	begin := time.Date(year, month, day, 0, 0, 0, 0, loc)
	return begin.Unix()
}

// BeginOfWeek 获取本周一零点时间戳(秒时间戳)
func BeginOfWeek(timestamp int64) int64 {
	t := GetTime(timestamp)

	// 计算与周一的日期差
	weekday := t.Weekday()
	offsetDays := int(weekday - time.Monday)
	if offsetDays < 0 {
		offsetDays += 7 // 处理周日的情况
	}

	// 获取周一日期并构造零点时间
	monday := t.AddDate(0, 0, -offsetDays)
	year, month, day := monday.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, loc).Unix()
}

// BeginOfMonth 获取当月第一天零点时间戳(秒时间戳)
func BeginOfMonth(timestamp int64) int64 {
	t := GetTime(timestamp)
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, loc).Unix()
}

// BetweenDay 计算两个时间戳之间的天数差
func BetweenDay(t1, t2 int64) int {
	// 计算两时间的零点时间差
	midnight1 := BeginOfDay(t1)
	midnight2 := BeginOfDay(t2)
	diffSeconds := midnight1 - midnight2

	// 计算绝对值天数差
	days := int(math.Abs(float64(diffSeconds)) / 86400)
	return days
}

// EndOfDay 获取某天的结束时间戳(秒时间戳)
func EndOfDay(timestamp int64) int64 {
	begin := BeginOfDay(timestamp)
	return begin + 86400 - 1
}

// EndOfWeek 获取本周日23:59:59时间戳(秒时间戳)
func EndOfWeek(timestamp int64) int64 {
	begin := BeginOfWeek(timestamp)
	return begin + 7*86400 - 1
}

// EndOfMonth 获取当月最后一天23:59:59时间戳(秒时间戳)
func EndOfMonth(timestamp int64) int64 {
	t := GetTime(timestamp)
	// 计算下个月第一天的零点时间戳，减1秒即为当月最后时间
	nextMonthFirstDay := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, loc)
	return nextMonthFirstDay.Unix() - 1
}

// Format 格式化时间戳
func Format(timestamp int64) string {
	return GetTime(timestamp).Format("2006-01-02 15:04:05")
}

// IsSameDay 判断两个时间戳是否在同一天
func IsSameDay(t1, t2 int64) bool {
	time1 := GetTime(t1)
	time2 := GetTime(t2)

	y1, m1, d1 := time1.Date()
	y2, m2, d2 := time2.Date()

	return y1 == y2 && m1 == m2 && d1 == d2
}

func FormatToStr(t time.Time, format string, timezone ...string) string {
	tf, ok := timeFormat[strings.ToLower(format)]
	if !ok {
		return ""
	}

	if timezone != nil && timezone[0] != "" {
		loc, err := time.LoadLocation(timezone[0])
		if err != nil {
			return ""
		}
		return t.In(loc).Format(tf)
	}
	return t.Format(tf)
}

func ParseToTime(str, format string, timezone ...string) (time.Time, error) {
	tf, ok := timeFormat[strings.ToLower(format)]
	if !ok {
		return time.Time{}, fmt.Errorf("format %s not support", format)
	}

	if timezone != nil && timezone[0] != "" {
		loc, err := time.LoadLocation(timezone[0])
		if err != nil {
			return time.Time{}, err
		}

		return time.ParseInLocation(tf, str, loc)
	}

	return time.Parse(tf, str)
}
