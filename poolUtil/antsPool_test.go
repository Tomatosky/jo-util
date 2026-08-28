package poolUtil

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func waitForSignal(t *testing.T, ch <-chan struct{}, timeout time.Duration, description string) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(timeout):
		t.Fatalf("timed out waiting for %s", description)
	}
}

func TestNewAntsPool(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		wantCap int
	}{
		{name: "fixed size", size: 10, wantCap: 10},
		{name: "zero means unlimited", size: 0, wantCap: -1},
		{name: "negative means unlimited", size: -1, wantCap: -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := NewAntsPool(tt.size)
			if pool == nil || pool.pool == nil {
				t.Fatal("NewAntsPool returned an uninitialized pool")
			}
			if got := pool.pool.Cap(); got != tt.wantCap {
				t.Errorf("pool capacity = %d, want %d", got, tt.wantCap)
			}
			if timedOut := pool.Shutdown(time.Second); timedOut {
				t.Error("unused pool shutdown timed out")
			}
		})
	}
}

func TestPoolSubmit(t *testing.T) {
	pool := NewAntsPool(1)
	executed := make(chan struct{}, 2)
	pool.Submit(func() { executed <- struct{}{} })
	pool.SubmitWithId("ignored", func() { executed <- struct{}{} })
	waitForSignal(t, executed, time.Second, "submitted task")
	waitForSignal(t, executed, time.Second, "task submitted with id")

	if timedOut := pool.Shutdown(time.Second); timedOut {
		t.Fatal("shutdown timed out after completed tasks")
	}
}

func TestPoolSubmitNilPanics(t *testing.T) {
	pool := NewAntsPool(1)
	defer pool.Shutdown(time.Second)

	defer func() {
		if got := recover(); got != "task cannot be nil" {
			t.Fatalf("Submit(nil) panic = %v, want %q", got, "task cannot be nil")
		}
	}()
	pool.Submit(nil)
}

func TestPoolSubmitAfterShutdownDoesNotLeakWaitGroup(t *testing.T) {
	pool := NewAntsPool(1)
	if timedOut := pool.Shutdown(time.Second); timedOut {
		t.Fatal("initial shutdown timed out")
	}

	var executed atomic.Bool
	pool.Submit(func() { executed.Store(true) })
	if timedOut := pool.Shutdown(200 * time.Millisecond); timedOut {
		t.Fatal("submission rejected by a closed ants pool leaked the wait group")
	}
	if executed.Load() {
		t.Error("task submitted after shutdown was executed")
	}
}

func TestPoolScheduleAtFixedRate(t *testing.T) {
	pool := NewAntsPool(2)
	executions := make(chan struct{})
	var count atomic.Int32
	stop := pool.ScheduleAtFixedRate(0, 40*time.Millisecond, func() {
		count.Add(1)
		executions <- struct{}{}
	})

	for i := 1; i <= 3; i++ {
		waitForSignal(t, executions, time.Second, "fixed-rate execution")
	}
	stop()
	stop()
	select {
	case <-executions:
		t.Error("fixed-rate schedule executed again after stop")
	case <-time.After(80 * time.Millisecond):
	}
	if got := count.Load(); got != 3 {
		t.Errorf("fixed-rate execution count after stop = %d, want 3", got)
	}
	if timedOut := pool.Shutdown(time.Second); timedOut {
		t.Fatal("fixed-rate pool shutdown timed out")
	}
}

func TestPoolScheduleWithFixedDelay(t *testing.T) {
	pool := NewAntsPool(2)
	executions := make(chan struct{})
	var count atomic.Int32
	stop := pool.ScheduleWithFixedDelay(0, 40*time.Millisecond, func() {
		count.Add(1)
		executions <- struct{}{}
	})

	for i := 1; i <= 3; i++ {
		waitForSignal(t, executions, time.Second, "fixed-delay execution")
	}
	stop()
	stop()
	select {
	case <-executions:
		t.Error("fixed-delay schedule executed again after stop")
	case <-time.After(80 * time.Millisecond):
	}
	if got := count.Load(); got != 3 {
		t.Errorf("fixed-delay execution count after stop = %d, want 3", got)
	}
	if timedOut := pool.Shutdown(time.Second); timedOut {
		t.Fatal("fixed-delay pool shutdown timed out")
	}
}

func TestPoolFixedDelayDoesNotRunWhenStopAndDelayAreReady(t *testing.T) {
	pool := NewAntsPool(4)
	done := make(chan struct{})
	close(done)
	var once sync.Once
	var count atomic.Int32

	for i := 0; i < 1000; i++ {
		pool.scheduleNextWithDelay(0, func() { count.Add(1) }, done, &once)
	}
	if timedOut := pool.Shutdown(time.Second); timedOut {
		t.Fatal("fixed-delay pool shutdown timed out")
	}
	if got := count.Load(); got != 0 {
		t.Errorf("fixed-delay schedule executed %d tasks after stop", got)
	}
}

func TestPoolScheduleCanStopDuringInitialDelay(t *testing.T) {
	tests := []struct {
		name     string
		schedule func(*AntsPool, func()) func()
	}{
		{
			name: "fixed rate",
			schedule: func(pool *AntsPool, task func()) func() {
				return pool.ScheduleAtFixedRate(80*time.Millisecond, 10*time.Millisecond, task)
			},
		},
		{
			name: "fixed delay",
			schedule: func(pool *AntsPool, task func()) func() {
				return pool.ScheduleWithFixedDelay(80*time.Millisecond, 10*time.Millisecond, task)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := NewAntsPool(1)
			executed := make(chan struct{}, 1)
			stop := tt.schedule(pool, func() { executed <- struct{}{} })
			stop()
			stop()

			select {
			case <-executed:
				t.Error("schedule executed after being stopped during initial delay")
			case <-time.After(150 * time.Millisecond):
			}
			if timedOut := pool.Shutdown(time.Second); timedOut {
				t.Fatal("cancelled schedule pool shutdown timed out")
			}
		})
	}
}

func TestPoolSchedulesHonorInitialDelay(t *testing.T) {
	const initialDelay = 80 * time.Millisecond
	tests := []struct {
		name     string
		schedule func(*AntsPool, func()) func()
	}{
		{
			name: "fixed rate",
			schedule: func(pool *AntsPool, task func()) func() {
				return pool.ScheduleAtFixedRate(initialDelay, time.Second, task)
			},
		},
		{
			name: "fixed delay",
			schedule: func(pool *AntsPool, task func()) func() {
				return pool.ScheduleWithFixedDelay(initialDelay, time.Second, task)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := NewAntsPool(1)
			executed := make(chan struct{}, 1)
			stop := tt.schedule(pool, func() { executed <- struct{}{} })
			t.Cleanup(func() {
				stop()
				pool.Shutdown(time.Second)
			})

			select {
			case <-executed:
				t.Fatal("schedule executed before its initial delay")
			case <-time.After(initialDelay / 4):
			}
			waitForSignal(t, executed, time.Second, "execution after initial delay")
			stop()
			if timedOut := pool.Shutdown(time.Second); timedOut {
				t.Fatal("scheduled pool shutdown timed out")
			}
		})
	}
}

func TestPoolFixedDelayStartsAfterPreviousTaskCompletes(t *testing.T) {
	const delay = 60 * time.Millisecond
	pool := NewAntsPool(2)
	starts := make(chan int32, 2)
	release := make(chan struct{})
	var releaseOnce sync.Once
	var count atomic.Int32
	stop := pool.ScheduleWithFixedDelay(0, delay, func() {
		call := count.Add(1)
		starts <- call
		if call == 1 {
			<-release
		}
	})
	t.Cleanup(func() {
		releaseOnce.Do(func() { close(release) })
		stop()
		pool.Shutdown(time.Second)
	})

	select {
	case call := <-starts:
		if call != 1 {
			t.Fatalf("first call number = %d, want 1", call)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for first fixed-delay task")
	}
	select {
	case call := <-starts:
		t.Fatalf("call %d started while the first task was still running", call)
	case <-time.After(delay + 20*time.Millisecond):
	}

	releaseOnce.Do(func() { close(release) })
	select {
	case call := <-starts:
		t.Fatalf("call %d started before the post-completion delay elapsed", call)
	case <-time.After(delay / 3):
	}
	select {
	case call := <-starts:
		if call != 2 {
			t.Fatalf("second call number = %d, want 2", call)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for second fixed-delay task")
	}
	stop()
	if timedOut := pool.Shutdown(time.Second); timedOut {
		t.Fatal("fixed-delay pool shutdown timed out")
	}
}

func TestPoolShutdown(t *testing.T) {
	t.Run("waits for running task", func(t *testing.T) {
		pool := NewAntsPool(1)
		started := make(chan struct{})
		release := make(chan struct{})
		pool.Submit(func() {
			close(started)
			<-release
		})
		waitForSignal(t, started, time.Second, "task start")

		result := make(chan bool, 1)
		go func() { result <- pool.Shutdown(time.Second) }()
		select {
		case <-result:
			t.Fatal("Shutdown returned before the running task completed")
		case <-time.After(20 * time.Millisecond):
		}
		close(release)
		select {
		case timedOut := <-result:
			if timedOut {
				t.Error("Shutdown timed out after task was released")
			}
		case <-time.After(time.Second):
			t.Fatal("Shutdown did not return after task completion")
		}
	})

	t.Run("reports timeout", func(t *testing.T) {
		pool := NewAntsPool(1)
		started := make(chan struct{})
		release := make(chan struct{})
		finished := make(chan struct{})
		pool.Submit(func() {
			close(started)
			<-release
			close(finished)
		})
		waitForSignal(t, started, time.Second, "task start")

		if timedOut := pool.Shutdown(20 * time.Millisecond); !timedOut {
			t.Error("Shutdown should report timeout while task is blocked")
		}
		close(release)
		waitForSignal(t, finished, time.Second, "timed-out task completion")
		pool.wg.Wait()
	})
}

func TestConcurrentUsage(t *testing.T) {
	pool := NewAntsPool(10)
	var counter atomic.Int32
	var submitters sync.WaitGroup

	for i := 0; i < 100; i++ {
		submitters.Add(1)
		go func() {
			defer submitters.Done()
			pool.Submit(func() { counter.Add(1) })
		}()
	}

	submitters.Wait()
	if timedOut := pool.Shutdown(time.Second); timedOut {
		t.Fatal("shutdown timed out after concurrent submissions")
	}
	if got := counter.Load(); got != 100 {
		t.Errorf("executed task count = %d, want 100", got)
	}
}
