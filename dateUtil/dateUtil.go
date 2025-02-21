package dateUtil

import (
	"math"
	"time"
)

var loc *time.Location

func getTime(timestamp int64) time.Time {
	if loc == nil {
		var err error
		loc, err = time.LoadLocation("Asia/Shanghai")
		if err != nil {
			panic(err) // 初始化时区失败，程序无法继续
		}
	}
	return time.Unix(timestamp, 0).In(loc)
}

// BeginOfDay 获取某天的起始时间戳(秒时间戳)
func BeginOfDay(timestamp int64) int64 {
	t := getTime(timestamp)
	year, month, day := t.Date()
	begin := time.Date(year, month, day, 0, 0, 0, 0, loc)
	return begin.Unix()
}

// BeginOfWeek 获取本周一零点时间戳(秒时间戳)
func BeginOfWeek(timestamp int64) int64 {
	t := getTime(timestamp)

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
	t := getTime(timestamp)
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
	t := getTime(timestamp)
	// 计算下个月第一天的零点时间戳，减1秒即为当月最后时间
	nextMonthFirstDay := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, loc)
	return nextMonthFirstDay.Unix() - 1
}

// Format 格式化时间戳
func Format(timestamp int64) string {
	return getTime(timestamp).Format("2006-01-02 15:04:05")
}

// IsSameDay 判断两个时间戳是否在同一天
func IsSameDay(t1, t2 int64) bool {
	time1 := getTime(t1)
	time2 := getTime(t2)

	y1, m1, d1 := time1.Date()
	y2, m2, d2 := time2.Date()

	return y1 == y2 && m1 == m2 && d1 == d2
}

// Parse 解析时间字符串
func Parse(str string) int64 {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", str, loc)
	if err != nil {
		panic(err)
	}
	return t.Unix()
}
