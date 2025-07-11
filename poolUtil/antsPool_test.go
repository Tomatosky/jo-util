package poolUtil

import (
	"sync"
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{"positive size", 10},
		{"zero size", 0},
		{"negative size", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := NewPool(tt.size)
			if pool == nil {
				t.Error("NewPool returned nil")
			}
			if pool.pool == nil {
				t.Error("ants pool not initialized")
			}
		})
	}
}

func TestPool_Submit(t *testing.T) {
	pool := NewPool(1)
	defer pool.Release(time.Second)

	// 测试正常提交任务
	var executed bool
	err := pool.Submit(func() {
		executed = true
	})
	if err != nil {
		t.Errorf("Submit failed: %v", err)
	}
	time.Sleep(100 * time.Millisecond) // 等待任务执行
	if !executed {
		t.Error("Task was not executed")
	}

	// 测试提交nil任务
	err = pool.Submit(nil)
	if err == nil {
		t.Error("Expected error when submitting nil task")
	}

	// 测试池已关闭的情况
	pool.Release(time.Second)
	err = pool.Submit(func() {})
	if err == nil {
		t.Error("Expected error when submitting to closed pool")
	}
}

func TestPool_ScheduleAtFixedRate(t *testing.T) {
	pool := NewPool(2)
	defer pool.Release(time.Second)

	var counter int
	var mu sync.Mutex

	// 测试固定频率调度
	stop := pool.ScheduleAtFixedRate(0, 100*time.Millisecond, func() {
		mu.Lock()
		counter++
		mu.Unlock()
	})

	time.Sleep(350 * time.Millisecond) // 应该执行大约3-4次
	stop()

	mu.Lock()
	if counter < 3 || counter > 4 {
		t.Errorf("Expected 3-4 executions, got %d", counter)
	}
	mu.Unlock()

	// 测试初始延迟
	counter = 0
	stop = pool.ScheduleAtFixedRate(200*time.Millisecond, 100*time.Millisecond, func() {
		mu.Lock()
		counter++
		mu.Unlock()
	})

	time.Sleep(150 * time.Millisecond) // 初始延迟未到，应该为0
	mu.Lock()
	if counter != 0 {
		t.Errorf("Expected 0 executions before initial delay, got %d", counter)
	}
	mu.Unlock()

	time.Sleep(200 * time.Millisecond) // 总共350ms，应该执行1-2次
	stop()

	mu.Lock()
	if counter < 1 || counter > 2 {
		t.Errorf("Expected 1-2 executions after initial delay, got %d", counter)
	}
	mu.Unlock()
}

func TestPool_ScheduleWithFixedDelay(t *testing.T) {
	pool := NewPool(2)
	defer pool.Release(time.Second)

	var counter int
	var mu sync.Mutex

	// 测试固定延迟调度
	stop := pool.ScheduleWithFixedDelay(0, 100*time.Millisecond, func() {
		mu.Lock()
		counter++
		mu.Unlock()
		time.Sleep(50 * time.Millisecond) // 模拟任务执行时间
	})

	time.Sleep(450 * time.Millisecond) // 应该执行大约3次 (每次任务+延迟约150ms)
	stop()

	mu.Lock()
	if counter < 2 || counter > 3 {
		t.Errorf("Expected 2-3 executions, got %d", counter)
	}
	mu.Unlock()

	// 测试初始延迟
	counter = 0
	stop = pool.ScheduleWithFixedDelay(200*time.Millisecond, 100*time.Millisecond, func() {
		mu.Lock()
		counter++
		mu.Unlock()
		time.Sleep(50 * time.Millisecond)
	})

	time.Sleep(150 * time.Millisecond) // 初始延迟未到，应该为0
	mu.Lock()
	if counter != 0 {
		t.Errorf("Expected 0 executions before initial delay, got %d", counter)
	}
	mu.Unlock()

	time.Sleep(300 * time.Millisecond) // 总共450ms，应该执行1-2次
	stop()

	mu.Lock()
	if counter < 1 || counter > 2 {
		t.Errorf("Expected 1-2 executions after initial delay, got %d", counter)
	}
	mu.Unlock()
}

func TestPool_Release(t *testing.T) {
	// 测试正常释放
	pool := NewPool(2)
	var wg sync.WaitGroup
	wg.Add(1)
	_ = pool.Submit(func() {
		defer wg.Done()
		time.Sleep(200 * time.Millisecond)
	})

	// 超时时间足够
	isTimeout := pool.Release(300 * time.Millisecond)
	if isTimeout {
		t.Error("Expected no timeout when waiting for tasks")
	}

	// 测试超时情况
	pool = NewPool(1)
	wg.Add(1)
	_ = pool.Submit(func() {
		defer wg.Done()
		time.Sleep(500 * time.Millisecond)
	})

	// 超时时间不足
	isTimeout = pool.Release(100 * time.Millisecond)
	if !isTimeout {
		t.Error("Expected timeout when not waiting long enough")
	}
	wg.Wait() // 确保测试不会泄漏goroutine

	// 测试多次释放
	pool = NewPool(1)
	pool.Release(time.Second)
	isTimeout = pool.Release(time.Second) // 第二次释放应该无害
	if isTimeout {
		t.Error("Second release should not timeout")
	}
}

func TestConcurrentUsage(t *testing.T) {
	pool := NewPool(10)
	defer pool.Release(time.Second)

	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 并发提交任务
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := pool.Submit(func() {
				mu.Lock()
				counter++
				mu.Unlock()
			})
			if err != nil {
				t.Errorf("Submit failed: %v", err)
			}
		}()
	}

	wg.Wait()
	time.Sleep(100 * time.Millisecond) // 等待所有任务完成

	mu.Lock()
	if counter != 100 {
		t.Errorf("Expected 100 increments, got %d", counter)
	}
	mu.Unlock()
}
