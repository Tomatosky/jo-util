package dateUtil

import "time"

// TimeInterval 计时器结构体
type TimeInterval struct {
	start time.Time
}

// NewTimer 创建一个新的计时器
func NewTimer() *TimeInterval {
	return &TimeInterval{
		start: time.Now(),
	}
}

// Interval 返回从开始到现在的毫秒数
func (t *TimeInterval) Interval() int64 {
	return time.Since(t.start).Milliseconds()
}

// IntervalRestart 返回从开始到现在的毫秒数，并重置开始时间
func (t *TimeInterval) IntervalRestart() int64 {
	elapsed := time.Since(t.start)
	t.start = time.Now()
	return elapsed.Milliseconds()
}
