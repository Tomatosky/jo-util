package eventUtil

import (
	"sync"
	"testing"
	"time"
)

//func TestPool(t *testing.T) {
//	logger.InitLog(nil)
//
//	pool := NewPool(100)
//	logger.Log.Info("")
//	for i := 0; i < 3; i++ {
//		pool.Submit(func() {
//			time.Sleep(2 * time.Second)
//			logger.Log.Info(fmt.Sprintf("%d", i))
//		})
//	}
//
//	pool.Release(10 * time.Second)
//	logger.Log.Info("")
//}

func TestNewPool(t *testing.T) {
	size := 10
	newPool := NewPool(size)
	if newPool == nil {
		t.Fatal("NewPool returned nil")
	}
	if newPool.pool == nil {
		t.Error("Pool's internal ants newPool is nil")
	}
	if newPool.pool.Cap() != size {
		t.Errorf("Expected newPool capacity %d, got %d", size, newPool.pool.Cap())
	}
}

func TestSubmit(t *testing.T) {
	pool := NewPool(1)
	defer pool.Release(time.Second)

	var executed bool
	err := pool.Submit(func() {
		executed = true
	})
	if err != nil {
		t.Errorf("Submit returned error: %v", err)
	}

	// Wait a bit for the task to execute
	time.Sleep(100 * time.Millisecond)
	if !executed {
		t.Error("Submitted task was not executed")
	}
}

func TestSubmitWithFullPool(t *testing.T) {
	pool := NewPool(1)
	defer pool.Release(time.Second)
	// First task will occupy the pool
	firstTaskDone := make(chan struct{})
	err := pool.Submit(func() {
		defer close(firstTaskDone)
		time.Sleep(500 * time.Millisecond)
	})
	if err != nil {
		t.Errorf("First submit failed: %v", err)
	}
	// Second submit should wait, not fail
	secondTaskDone := make(chan struct{})
	start := time.Now()
	err = pool.Submit(func() {
		close(secondTaskDone)
	})
	if err != nil {
		t.Errorf("Second submit failed: %v", err)
	}
	// Wait for first task to complete
	<-firstTaskDone
	// Check if second task was executed
	select {
	case <-secondTaskDone:
		duration := time.Since(start)
		if duration < 500*time.Millisecond {
			t.Error("Second task executed before first task completed")
		}
	case <-time.After(600 * time.Millisecond):
		t.Error("Second task was not executed after first task completed")
	}
}

func TestRelease(t *testing.T) {
	pool := NewPool(2)

	var wg sync.WaitGroup
	wg.Add(1)

	start := time.Now()
	err := pool.Submit(func() {
		defer wg.Done()
		time.Sleep(200 * time.Millisecond)
	})
	if err != nil {
		t.Errorf("Submit failed: %v", err)
	}

	// Release with timeout longer than task duration
	pool.Release(300 * time.Millisecond)
	duration := time.Since(start)

	if duration < 200*time.Millisecond {
		t.Error("Release returned before task completed")
	}
	if duration > 250*time.Millisecond {
		t.Error("Release took longer than expected")
	}
}

func TestReleaseTimeout(t *testing.T) {
	pool := NewPool(2)

	var wg sync.WaitGroup
	wg.Add(1)

	start := time.Now()
	err := pool.Submit(func() {
		defer wg.Done()
		time.Sleep(500 * time.Millisecond)
	})
	if err != nil {
		t.Errorf("Submit failed: %v", err)
	}

	// Release with timeout shorter than task duration
	pool.Release(100 * time.Millisecond)
	duration := time.Since(start)

	if duration > 150*time.Millisecond {
		t.Error("Release didn't timeout as expected")
	}
}
