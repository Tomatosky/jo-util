package dateUtil

import (
	"testing"
	"time"
)

func TestNewTimer(t *testing.T) {
	before := time.Now()
	timer := NewTimer()
	after := time.Now()
	if timer == nil {
		t.Fatal("NewTimer() returned nil")
	}
	if timer.start.Before(before) || timer.start.After(after) {
		t.Errorf("NewTimer() start = %v, want within [%v, %v]", timer.start, before, after)
	}
}

func TestInterval(t *testing.T) {
	start := time.Now().Add(-2 * time.Second)
	timer := &TimeInterval{start: start}
	before := time.Since(start).Milliseconds()
	got := timer.Interval()
	after := time.Since(start).Milliseconds()
	if got < before || got > after {
		t.Errorf("Interval() = %dms, want within [%dms, %dms]", got, before, after)
	}
}

func TestIntervalRestart(t *testing.T) {
	originalStart := time.Now().Add(-2 * time.Second)
	timer := &TimeInterval{start: originalStart}
	beforeCall := time.Now()
	wantLower := beforeCall.Sub(originalStart).Milliseconds()
	got := timer.IntervalRestart()
	afterCall := time.Now()
	wantUpper := afterCall.Sub(originalStart).Milliseconds()

	if got < wantLower || got > wantUpper {
		t.Errorf("IntervalRestart() = %dms, want within [%dms, %dms]", got, wantLower, wantUpper)
	}
	if timer.start.Before(beforeCall) || timer.start.After(afterCall) {
		t.Errorf("restart time = %v, want within [%v, %v]", timer.start, beforeCall, afterCall)
	}

	beforeInterval := time.Since(timer.start).Milliseconds()
	second := timer.Interval()
	afterInterval := time.Since(timer.start).Milliseconds()
	if second < beforeInterval || second > afterInterval {
		t.Errorf("Interval() after restart = %dms, want within [%dms, %dms]", second, beforeInterval, afterInterval)
	}
}
