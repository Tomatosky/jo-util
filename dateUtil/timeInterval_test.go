package dateUtil

import (
	"testing"
	"time"
)

func TestNewTimer(t *testing.T) {
	timer := NewTimer()
	if timer == nil {
		t.Error("NewTimer() returned nil")
	}
	if timer.start.IsZero() {
		t.Error("NewTimer() start time is zero")
	}
}

func TestInterval(t *testing.T) {
	timer := &TimeInterval{start: time.Now()}

	// Wait a short time
	time.Sleep(10 * time.Millisecond)

	elapsed := timer.Interval()
	if elapsed <= 0 {
		t.Errorf("Interval() returned non-positive value: %d", elapsed)
	}
	if elapsed > 100 {
		t.Errorf("Interval() returned unexpectedly large value: %d", elapsed)
	}
}

func TestIntervalRestart(t *testing.T) {
	timer := &TimeInterval{start: time.Now()}

	// Wait a short time
	time.Sleep(10 * time.Millisecond)

	firstElapsed := timer.IntervalRestart()
	if firstElapsed <= 0 {
		t.Errorf("IntervalRestart() first call returned non-positive value: %d", firstElapsed)
	}
	if firstElapsed > 100 {
		t.Errorf("IntervalRestart() first call returned unexpectedly large value: %d", firstElapsed)
	}

	// Check if timer was restarted
	time.Sleep(5 * time.Millisecond)
	secondElapsed := timer.Interval()
	if secondElapsed <= 0 {
		t.Errorf("Interval() after restart returned non-positive value: %d", secondElapsed)
	}
	if secondElapsed > 50 {
		t.Errorf("Interval() after restart returned unexpectedly large value: %d", secondElapsed)
	}
	if secondElapsed >= firstElapsed {
		t.Errorf("Interval() after restart (%d) should be less than first interval (%d)", secondElapsed, firstElapsed)
	}
}
