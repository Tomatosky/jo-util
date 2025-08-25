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

func DayOfWeek(t time.Time) int {
	weekday := t.Weekday()
	if weekday == time.Sunday {
		return 7
	}
	return int(weekday)
}

func OffsetDay(t time.Time, offset int) time.Time {
	return t.AddDate(0, 0, offset)
}

// BeginOfDay 获取某天的起始时间
func BeginOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	begin := time.Date(year, month, day, 0, 0, 0, 0, loc)
	return begin
}

// BeginOfWeek 获取本周一零点时间
func BeginOfWeek(t time.Time) time.Time {
	// 计算与周一的日期差
	weekday := t.Weekday()
	offsetDays := int(weekday - time.Monday)
	if offsetDays < 0 {
		offsetDays += 7 // 处理周日的情况
	}

	// 获取周一日期并构造零点时间
	monday := t.AddDate(0, 0, -offsetDays)
	year, month, day := monday.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, loc)
}

// BeginOfMonth 获取当月第一天零点时间
func BeginOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, loc)
}

func BeginOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

// EndOfDay 获取某天的结束时间(23:59:59)
func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
}

// EndOfWeek 获取本周日23:59:59时间
func EndOfWeek(t time.Time) time.Time {
	// 计算到本周日的天数差(周日是0，周一是1...周六是6)
	daysUntilSunday := (7 - int(t.Weekday())) % 7
	endOfWeek := t.AddDate(0, 0, daysUntilSunday)
	return EndOfDay(endOfWeek)
}

// EndOfMonth 获取当月最后一天23:59:59时间
func EndOfMonth(t time.Time) time.Time {
	// 下个月的第0天就是本月的最后一天
	firstOfNextMonth := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
	endOfMonth := firstOfNextMonth.Add(-time.Second)
	return endOfMonth
}

func EndOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 12, 31, 23, 59, 59, 0, t.Location())
}

func BetweenMinute(t1, t2 time.Time, isReset bool) int {
	if isReset {
		t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), 0, 0, t1.Location())
		t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), t2.Hour(), t2.Minute(), 0, 0, t2.Location())
	}
	diffSeconds := t1.Unix() - t2.Unix()
	minutes := int(math.Abs(float64(diffSeconds)) / 60)
	return minutes
}

func BetweenHour(t1, t2 time.Time, isReset bool) int {
	if isReset {
		t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), 0, 0, 0, t1.Location())
		t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), t2.Hour(), 0, 0, 0, t2.Location())
	}
	diffSeconds := t1.Unix() - t2.Unix()
	hours := int(math.Abs(float64(diffSeconds)) / 3600)
	return hours
}

// BetweenDay 计算两个时间之间的天数差
func BetweenDay(t1, t2 time.Time, isReset bool) int {
	if isReset {
		t1 = BeginOfDay(t1)
		t2 = BeginOfDay(t2)
	}
	diffSeconds := t1.Unix() - t2.Unix()

	// 计算绝对值天数差
	days := int(math.Abs(float64(diffSeconds)) / 86400)
	return days
}

// BetweenWeek 计算两个时间之间的周数差
func BetweenWeek(t1, t2 time.Time, isReset bool) int {
	if isReset {
		t1 = BeginOfWeek(t1)
		t2 = BeginOfWeek(t2)
	}
	diffSeconds := t1.Unix() - t2.Unix()
	// 计算绝对值周数差
	weeks := int(math.Abs(float64(diffSeconds)) / (86400 * 7))
	return weeks
}

// BetweenMonth 计算两个时间之间的月数差
func BetweenMonth(t1, t2 time.Time, isReset bool) int {
	if isReset {
		t1 = BeginOfMonth(t1)
		t2 = BeginOfMonth(t2)
	}

	// 计算月数差
	years := t1.Year() - t2.Year()
	months := int(t1.Month()) - int(t2.Month())
	totalMonths := years*12 + months
	return int(math.Abs(float64(totalMonths)))
}

func BetweenYear(t1, t2 time.Time, isReset bool) int {
	if isReset {
		t1 = BeginOfYear(t1)
		t2 = BeginOfYear(t2)
	}
	years := t1.Year() - t2.Year()
	return int(math.Abs(float64(years)))
}

// IsSameDay 判断两个时间是否在同一天
func IsSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func IsSameWeek(t1, t2 time.Time) bool {
	// 获取年份和周数
	y1, w1 := t1.ISOWeek()
	y2, w2 := t2.ISOWeek()
	return y1 == y2 && w1 == w2
}

func IsSameMonth(t1, t2 time.Time) bool {
	y1, m1, _ := t1.Date()
	y2, m2, _ := t2.Date()
	return y1 == y2 && m1 == m2
}

func IsSameYear(t1, t2 time.Time) bool {
	y1, _, _ := t1.Date()
	y2, _, _ := t2.Date()
	return y1 == y2
}

func ParseToTime(str, format string, timezone ...string) (time.Time, error) {
	tf, ok := timeFormat[strings.ToLower(format)]
	if !ok {
		return time.Time{}, fmt.Errorf("format %s not support", format)
	}

	location := loc
	if timezone != nil && timezone[0] != "" {
		var err error
		location, err = time.LoadLocation(timezone[0])
		if err != nil {
			return time.Time{}, err
		}
	}

	return time.ParseInLocation(tf, str, location)
}
