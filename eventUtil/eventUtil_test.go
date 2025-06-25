package eventUtil

import (
	"errors"
	"sync"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	em := NewEventManager(32)

	tests := []struct {
		name        string
		eventName   string
		handler     EventHandler
		expectedErr error
	}{
		{"valid registration", "test", func(data interface{}) {}, nil},
		{"empty event name", "", func(data interface{}) {}, errors.New("event name cannot be empty")},
		{"nil handler", "test", nil, errors.New("handler cannot be nil")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := em.Register(tt.eventName, tt.handler)
			if (err != nil) != (tt.expectedErr != nil) {
				t.Errorf("Expected error: %v, got: %v", tt.expectedErr, err)
			}
			if err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error message: %q, got: %q", tt.expectedErr.Error(), err.Error())
			}
		})
	}
}

func TestRegisterDuringDestruction(t *testing.T) {
	em := NewEventManager(32)

	// Start destruction
	go em.ShutDown(1 * time.Second)

	// Wait a bit to ensure destruction flag is set
	time.Sleep(100 * time.Millisecond)

	err := em.Register("test", func(data interface{}) {})
	if err == nil || err.Error() != "event manager is destroying, cannot register new handlers" {
		t.Errorf("Expected error during destruction, got: %v", err)
	}
}

func TestTrigger(t *testing.T) {
	// 测试正常触发事件
	t.Run("NormalTrigger", func(t *testing.T) {
		em := NewEventManager(32)
		defer em.ShutDown(time.Second)

		var wg sync.WaitGroup
		wg.Add(1)

		// 注册事件处理函数
		err := em.Register("test", func(data interface{}) {
			if data != "test data" {
				t.Error("Expected 'test data', got", data)
			}
			wg.Done()
		})
		if err != nil {
			t.Error("Register failed:", err)
		}

		// 触发事件
		err = em.Trigger("test", "test data")
		if err != nil {
			t.Error("Trigger failed:", err)
		}

		wg.Wait() // 等待处理函数执行完成
	})

	// 测试触发不存在的事件
	t.Run("TriggerNonExistentEvent", func(t *testing.T) {
		em := NewEventManager(32)
		defer em.ShutDown(time.Second)

		err := em.Trigger("nonexistent", nil)
		if err == nil || err.Error() != "event not found" {
			t.Error("Expected 'event not found' error, got", err)
		}
	})

	// 测试空事件名
	t.Run("EmptyEventName", func(t *testing.T) {
		em := NewEventManager(32)
		defer em.ShutDown(time.Second)

		err := em.Trigger("", nil)
		if err == nil || err.Error() != "event name cannot be empty" {
			t.Error("Expected 'event name cannot be empty' error, got", err)
		}
	})

	// 测试在销毁状态下触发事件
	t.Run("TriggerWhileDestroying", func(t *testing.T) {
		em := NewEventManager(32)
		em.ShutDown(time.Second) // 立即开始销毁

		err := em.Trigger("test", nil)
		if err == nil || err.Error() != "event manager is destroying, cannot trigger events" {
			t.Error("Expected 'event manager is destroying' error, got", err)
		}
	})

	// 测试并发触发事件
	t.Run("ConcurrentTrigger", func(t *testing.T) {
		em := NewEventManager(32)
		defer em.ShutDown(time.Second)
		var counter int
		var mu sync.Mutex
		const numHandlers = 5
		const numTriggers = 10
		// 使用WaitGroup等待所有处理函数完成
		var handlersWg sync.WaitGroup
		handlersWg.Add(numHandlers * numTriggers)
		// 注册多个处理函数
		for i := 0; i < numHandlers; i++ {
			err := em.Register("concurrent", func(data interface{}) {
				defer handlersWg.Done()
				mu.Lock()
				counter++
				mu.Unlock()
			})
			if err != nil {
				t.Error("Register failed:", err)
			}
		}
		// 并发触发事件
		var triggersWg sync.WaitGroup
		for i := 0; i < numTriggers; i++ {
			triggersWg.Add(1)
			go func() {
				defer triggersWg.Done()
				err := em.Trigger("concurrent", nil)
				if err != nil {
					t.Error("Trigger failed:", err)
				}
			}()
		}
		triggersWg.Wait() // 等待所有触发完成
		// 等待所有处理函数完成
		handlersWg.Wait()
		// 验证所有处理函数都被调用了
		expected := numHandlers * numTriggers
		if counter != expected {
			t.Errorf("Expected counter to be %d, got %d", expected, counter)
		}
	})

	// 测试处理函数panic的情况
	t.Run("HandlerPanic", func(t *testing.T) {
		em := NewEventManager(32)
		defer em.ShutDown(time.Second)

		// 注册会panic的处理函数
		err := em.Register("panic", func(data interface{}) {
			panic("test panic")
		})
		if err != nil {
			t.Error("Register failed:", err)
		}

		// 触发事件，应该不会导致测试失败
		err = em.Trigger("panic", nil)
		if err != nil {
			t.Error("Trigger failed:", err)
		}

		// 等待一段时间确保处理函数执行完成
		time.Sleep(100 * time.Millisecond)
	})
}

func TestTriggerInvalidEvent(t *testing.T) {
	em := NewEventManager(32)

	err := em.Trigger("nonexistent", nil)
	if err == nil || err.Error() != "event not found" {
		t.Errorf("Expected 'event not found' error, got: %v", err)
	}
}

func TestTriggerDuringDestruction(t *testing.T) {
	em := NewEventManager(32)

	// Register a handler first
	em.Register("test", func(data interface{}) {})

	// Start destruction
	go em.ShutDown(1 * time.Second)

	// Wait a bit to ensure destruction flag is set
	time.Sleep(100 * time.Millisecond)

	err := em.Trigger("test", nil)
	if err == nil || err.Error() != "event manager is destroying, cannot trigger events" {
		t.Errorf("Expected error during destruction, got: %v", err)
	}
}

func TestHasEvent(t *testing.T) {
	em := NewEventManager(32)

	if em.HasEvent("test") {
		t.Error("Event should not exist before registration")
	}

	em.Register("test", func(data interface{}) {})
	if !em.HasEvent("test") {
		t.Error("Event should exist after registration")
	}
}

func TestClear(t *testing.T) {
	em := NewEventManager(32)

	em.Register("test1", func(data interface{}) {})
	em.Register("test2", func(data interface{}) {})

	em.Clear()

	if em.HasEvent("test1") || em.HasEvent("test2") {
		t.Error("Events should be cleared")
	}
}

func TestEventManager_ShutDown(t *testing.T) {
	// 测试用例1: 正常关闭
	t.Run("正常关闭", func(t *testing.T) {
		em := NewEventManager(10)
		em.ShutDown(time.Second)

		if !em.isDestroying() {
			t.Error("期望销毁标志为true，但实际为false")
		}
	})

	// 测试用例2: 关闭后不能注册新事件
	t.Run("关闭后不能注册新事件", func(t *testing.T) {
		em := NewEventManager(10)
		em.ShutDown(time.Second)

		err := em.Register("test", func(data interface{}) {})
		if err == nil || err.Error() != "event manager is destroying, cannot register new handlers" {
			t.Errorf("期望错误消息为'event manager is destroying, cannot register new handlers'，但实际为'%v'", err)
		}
	})

	// 测试用例3: 关闭后不能触发事件
	t.Run("关闭后不能触发事件", func(t *testing.T) {
		em := NewEventManager(10)
		em.ShutDown(time.Second)

		err := em.Trigger("test", nil)
		if err == nil || err.Error() != "event manager is destroying, cannot trigger events" {
			t.Errorf("期望错误消息为'event manager is destroying, cannot trigger events'，但实际为'%v'", err)
		}
	})

	// 测试用例4: 关闭后检查销毁标志
	t.Run("关闭后检查销毁标志", func(t *testing.T) {
		em := NewEventManager(10)
		em.ShutDown(time.Second)

		if !em.isDestroying() {
			t.Error("期望销毁标志为true，但实际为false")
		}
	})

	// 测试用例5: 关闭后池是否释放
	t.Run("关闭后池是否释放", func(t *testing.T) {
		em := NewEventManager(10)
		// 注册一个事件用于测试
		_ = em.Register("test", func(data interface{}) {})

		// 触发事件确保池在工作
		_ = em.Trigger("test", nil)

		// 关闭管理器
		em.ShutDown(time.Second)

		// 尝试再次提交任务应该失败
		err := em.pool.Submit(func() {})
		if err == nil {
			t.Error("期望池已关闭不能提交任务，但实际可以提交")
		}
	})

	// 测试用例6: 超时关闭
	t.Run("超时关闭", func(t *testing.T) {
		em := NewEventManager(10)
		// 注册一个长时间运行的事件
		_ = em.Register("long", func(data interface{}) {
			time.Sleep(2 * time.Second)
		})

		// 触发长时间运行的事件
		_ = em.Trigger("long", nil)

		// 尝试快速关闭(超时时间很短)
		start := time.Now()
		em.ShutDown(100 * time.Millisecond)
		duration := time.Since(start)

		if duration > 200*time.Millisecond {
			t.Error("期望在超时时间内关闭，但实际关闭时间过长")
		}
	})
}

func TestConcurrentAccess(t *testing.T) {
	em := NewEventManager(32)
	var wg sync.WaitGroup

	// Concurrent registrations
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			em.Register("test", func(data interface{}) {})
		}(i)
	}

	// Concurrent triggers
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = em.Trigger("test", nil)
		}()
	}

	wg.Wait()

	// Verify all handlers were registered
	em.lock.RLock()
	defer em.lock.RUnlock()
	if len(em.handlers["test"]) != 100 {
		t.Errorf("Expected 100 handlers, got %d", len(em.handlers["test"]))
	}
}
