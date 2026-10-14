package poolUtil

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type blockingJSONID struct {
	started chan<- struct{}
	release <-chan struct{}
}

func (id blockingJSONID) MarshalJSON() ([]byte, error) {
	close(id.started)
	<-id.release
	return []byte("0"), nil
}

func newTestIDPool(poolSize int64, queueSize int) *IdPool {
	return NewIdPool(&IdPoolOpt{
		PoolSize:  poolSize,
		QueueSize: queueSize,
		PoolName:  "test-id-pool",
	})
}

func TestIdPool(t *testing.T) {
	t.Run("rejects invalid sizes", func(t *testing.T) {
		tests := []IdPoolOpt{
			{PoolSize: 0, QueueSize: 1},
			{PoolSize: 1, QueueSize: 0},
			{PoolSize: -1, QueueSize: 1},
			{PoolSize: 1, QueueSize: -1},
		}
		for _, opt := range tests {
			func() {
				defer func() {
					if got := recover(); got != "pool size and queue size must be greater than 0" {
						t.Errorf("NewIdPool(%+v) panic = %v", opt, got)
					}
				}()
				NewIdPool(&opt)
			}()
		}
	})

	t.Run("executes tasks and cleans bookkeeping", func(t *testing.T) {
		pool := newTestIDPool(2, 4)
		if got := pool.poolName; got != "test-id-pool" {
			t.Errorf("pool name = %q, want %q", got, "test-id-pool")
		}
		var counter atomic.Int32
		done := make(chan struct{}, 2)
		for i := 0; i < 2; i++ {
			pool.SubmitWithId("1", func() {
				counter.Add(1)
				done <- struct{}{}
			})
		}
		waitForSignal(t, done, time.Second, "first id-pool task")
		waitForSignal(t, done, time.Second, "second id-pool task")

		if timedOut := pool.Shutdown(time.Second); timedOut {
			t.Fatal("shutdown timed out")
		}
		if got := counter.Load(); got != 2 {
			t.Errorf("executed task count = %d, want 2", got)
		}
		if got := pool.GetTaskCount(1); got != 0 {
			t.Errorf("task count after completion = %d, want 0", got)
		}
	})

	t.Run("negative id selects a valid worker", func(t *testing.T) {
		pool := newTestIDPool(2, 2)
		executed := make(chan struct{}, 1)
		pool.SubmitWithId(-1, func() { executed <- struct{}{} })
		waitForSignal(t, executed, time.Second, "task with negative id")
		if timedOut := pool.Shutdown(time.Second); timedOut {
			t.Fatal("shutdown timed out")
		}
		if got := pool.GetTaskCount(-1); got != 0 {
			t.Errorf("task count for negative id = %d, want 0", got)
		}
	})

	t.Run("Submit executes task without explicit id", func(t *testing.T) {
		pool := newTestIDPool(2, 2)
		executed := make(chan struct{}, 1)
		pool.Submit(func() { executed <- struct{}{} })
		waitForSignal(t, executed, time.Second, "task submitted without id")
		if timedOut := pool.Shutdown(time.Second); timedOut {
			t.Fatal("shutdown timed out")
		}
		if got := pool.GetTaskCount(999999); got != 0 {
			t.Errorf("unknown id task count = %d, want 0", got)
		}
	})

	t.Run("concurrent submissions all complete", func(t *testing.T) {
		pool := newTestIDPool(8, 512)
		const numTasks = 2000
		var counter atomic.Int32
		var submitters sync.WaitGroup
		for i := 0; i < numTasks; i++ {
			submitters.Add(1)
			go func(id int) {
				defer submitters.Done()
				pool.SubmitWithId(id%8, func() { counter.Add(1) })
			}(i)
		}
		submitters.Wait()

		if timedOut := pool.Shutdown(3 * time.Second); timedOut {
			t.Fatal("shutdown timed out after concurrent submissions")
		}
		if got := counter.Load(); got != numTasks {
			t.Errorf("executed task count = %d, want %d", got, numTasks)
		}
		for id := 0; id < 8; id++ {
			if got := pool.GetTaskCount(id); got != 0 {
				t.Errorf("task count for id %d = %d, want 0", id, got)
			}
		}
		if got := len(pool.idTaskCounts); got != 0 {
			t.Fatalf("completed task counters retained %d ids", got)
		}
	})

	t.Run("full queue rolls back rejected task", func(t *testing.T) {
		pool := newTestIDPool(1, 1)
		started := make(chan struct{})
		release := make(chan struct{})
		var executed atomic.Int32
		pool.SubmitWithId(0, func() {
			close(started)
			<-release
			executed.Add(1)
		})
		waitForSignal(t, started, time.Second, "blocking id-pool task")
		pool.SubmitWithId(0, func() { executed.Add(1) })
		pool.SubmitWithId(0, func() { executed.Add(100) })

		if got := pool.GetTaskCount(0); got != 2 {
			t.Errorf("accepted task count = %d, want 2", got)
		}
		if got := pool.MaxQueue(); got != 1 {
			t.Errorf("maximum queued tasks = %d, want 1", got)
		}

		close(release)
		if timedOut := pool.Shutdown(time.Second); timedOut {
			t.Fatal("shutdown timed out while draining accepted tasks")
		}
		if got := executed.Load(); got != 2 {
			t.Errorf("executed value = %d, want 2; rejected task must not run", got)
		}
		if got := pool.GetTaskCount(0); got != 0 {
			t.Errorf("task count after drain = %d, want 0", got)
		}
	})

	t.Run("shutdown drains pending tasks", func(t *testing.T) {
		pool := newTestIDPool(1, 10)
		const numTasks = 5
		var counter atomic.Int32
		for i := 0; i < numTasks; i++ {
			pool.SubmitWithId(0, func() { counter.Add(1) })
		}

		if timedOut := pool.Shutdown(time.Second); timedOut {
			t.Fatal("shutdown timed out while draining pending tasks")
		}
		if got := counter.Load(); got != numTasks {
			t.Errorf("drained task count = %d, want %d", got, numTasks)
		}
	})

	t.Run("shutdown reports timeout without leaking worker", func(t *testing.T) {
		pool := newTestIDPool(1, 1)
		started := make(chan struct{})
		release := make(chan struct{})
		pool.SubmitWithId(0, func() {
			close(started)
			<-release
		})
		waitForSignal(t, started, time.Second, "blocking id-pool task")

		if timedOut := pool.Shutdown(20 * time.Millisecond); !timedOut {
			t.Error("Shutdown should report timeout while task is blocked")
		}
		close(release)
		if timedOut := pool.Shutdown(time.Second); timedOut {
			t.Fatal("Shutdown timed out after timed-out task completed")
		}
		if got := pool.GetTaskCount(0); got != 0 {
			t.Errorf("task count after timed-out task completes = %d, want 0", got)
		}
	})

	t.Run("task counts are exact while blocked", func(t *testing.T) {
		pool := newTestIDPool(4, 64)
		const tasksPerID = 50
		release := make(chan struct{})
		for id := 0; id < 4; id++ {
			for i := 0; i < tasksPerID; i++ {
				pool.SubmitWithId(id, func() { <-release })
			}
		}

		for id := 0; id < 4; id++ {
			if got := pool.GetTaskCount(id); got != tasksPerID {
				t.Errorf("blocked task count for id %d = %d, want %d", id, got, tasksPerID)
			}
		}
		close(release)
		if timedOut := pool.Shutdown(2 * time.Second); timedOut {
			t.Fatal("shutdown timed out after releasing blocked tasks")
		}
		for id := 0; id < 4; id++ {
			if got := pool.GetTaskCount(id); got != 0 {
				t.Errorf("completed task count for id %d = %d, want 0", id, got)
			}
		}
	})

	t.Run("task panic still cleans bookkeeping", func(t *testing.T) {
		pool := newTestIDPool(1, 1)
		pool.SubmitWithId(7, func() { panic("test panic") })
		if timedOut := pool.Shutdown(time.Second); timedOut {
			t.Fatal("shutdown timed out after panicking task")
		}
		if got := pool.GetTaskCount(7); got != 0 {
			t.Errorf("task count after panic = %d, want 0", got)
		}
	})

	t.Run("submission after shutdown is ignored", func(t *testing.T) {
		pool := newTestIDPool(1, 1)
		if timedOut := pool.Shutdown(time.Second); timedOut {
			t.Fatal("initial shutdown timed out")
		}
		var executed atomic.Bool
		pool.SubmitWithId(0, func() { executed.Store(true) })
		if executed.Load() {
			t.Error("task submitted after shutdown was executed")
		}
		if got := pool.GetTaskCount(0); got != 0 {
			t.Errorf("task count after rejected submission = %d, want 0", got)
		}
	})

	t.Run("shutdown is idempotent", func(t *testing.T) {
		pool := newTestIDPool(2, 2)
		if timedOut := pool.Shutdown(time.Second); timedOut {
			t.Fatal("initial shutdown timed out")
		}
		if timedOut := pool.Shutdown(time.Second); timedOut {
			t.Fatal("repeated shutdown timed out")
		}
	})

	t.Run("shutdown rejects a submission paused before enqueue", func(t *testing.T) {
		pool := newTestIDPool(1, 1)
		started := make(chan struct{})
		release := make(chan struct{})
		submitted := make(chan struct{})
		var executed atomic.Bool

		go func() {
			defer close(submitted)
			pool.SubmitWithId(blockingJSONID{started: started, release: release}, func() {
				executed.Store(true)
			})
		}()

		waitForSignal(t, started, time.Second, "blocked id conversion")
		if timedOut := pool.Shutdown(time.Second); timedOut {
			t.Fatal("shutdown timed out while submission was paused before enqueue")
		}
		close(release)
		waitForSignal(t, submitted, time.Second, "paused submission return")

		if executed.Load() {
			t.Error("submission that lost the shutdown race was executed")
		}
		if got := pool.GetTaskCount(0); got != 0 {
			t.Errorf("task count after shutdown race = %d, want 0", got)
		}
	})

}
