package eventUtil

import (
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestEventManager(t *testing.T) {
	// 初始化测试用的logger
	logger, _ := zap.NewDevelopment()

	// 测试NewEventManager
	t.Run("TestNewEventManager", func(t *testing.T) {
		// 测试正常情况
		opt := &EventOpt{
			PoolSize:  2,
			QueueSize: 10,
			Logger:    logger,
		}
		manager := NewEventManager(opt)
		if manager == nil {
			t.Error("NewEventManager should not return nil")
		}

		// 测试异常情况 - 参数不合法
		invalidOpt := &EventOpt{
			PoolSize:  0,
			QueueSize: 0,
		}
		defer func() {
			if r := recover(); r == nil {
				t.Error("NewEventManager should panic with invalid options")
			}
		}()
		_ = NewEventManager(invalidOpt)
	})

	// 测试Register
	t.Run("TestRegister", func(t *testing.T) {
		manager := NewEventManager(&EventOpt{PoolSize: 2, QueueSize: 10, Logger: logger})

		// 测试正常注册
		err := manager.Register("testEvent", func(data interface{}) {})
		if err != nil {
			t.Errorf("Register should not return error, got: %v", err)
		}

		// 测试重复注册
		err = manager.Register("testEvent", func(data interface{}) {})
		if err != nil {
			t.Errorf("Register should allow multiple handlers, got: %v", err)
		}

		// 测试空事件名
		err = manager.Register("", func(data interface{}) {})
		if err == nil {
			t.Error("Register should return error for empty event name")
		}

		// 测试nil handler
		err = manager.Register("testEvent", nil)
		if err == nil {
			t.Error("Register should return error for nil handler")
		}

		// 测试销毁状态下注册
		manager.ShutDown(time.Second)
		err = manager.Register("testEvent", func(data interface{}) {})
		if err == nil {
			t.Error("Register should return error when manager is destroying")
		}
	})

	// 测试Trigger和TriggerSync
	t.Run("TestTrigger", func(t *testing.T) {
		manager := NewEventManager(&EventOpt{PoolSize: 2, QueueSize: 10, Logger: logger})
		var wg sync.WaitGroup
		var triggered bool

		// 注册测试handler
		err := manager.Register("testEvent", func(data interface{}) {
			defer wg.Done()
			triggered = true
			if data != "testData" {
				t.Errorf("Expected data to be 'testData', got %v", data)
			}
		})
		if err != nil {
			t.Fatalf("Register failed: %v", err)
		}

		// 测试Trigger
		wg.Add(1)
		err = manager.Trigger("testEvent", "testData")
		if err != nil {
			t.Errorf("Trigger failed: %v", err)
		}
		wg.Wait()
		if !triggered {
			t.Error("Handler was not triggered")
		}

		// 测试TriggerSync
		triggered = false
		wg.Add(1)
		manager.TriggerSync("testEvent", "testData")
		if !triggered {
			t.Error("Handler was not triggered in TriggerSync")
		}

		// 测试不存在的event
		err = manager.Trigger("nonexistent", nil)
		if err == nil {
			t.Error("Trigger should return error for nonexistent event")
		}

		// 测试空事件名
		err = manager.Trigger("", nil)
		if err == nil {
			t.Error("Trigger should return error for empty event name")
		}

		// 测试销毁状态下触发
		manager.ShutDown(time.Second)
		err = manager.Trigger("testEvent", nil)
		if err == nil {
			t.Error("Trigger should return error when manager is destroying")
		}
	})

	// 测试HasEvent
	t.Run("TestHasEvent", func(t *testing.T) {
		manager := NewEventManager(&EventOpt{PoolSize: 2, QueueSize: 10, Logger: logger})

		// 注册测试事件
		_ = manager.Register("testEvent", func(data interface{}) {})

		// 测试存在的事件
		if !manager.HasEvent("testEvent") {
			t.Error("HasEvent should return true for registered event")
		}

		// 测试不存在的事件
		if manager.HasEvent("nonexistent") {
			t.Error("HasEvent should return false for unregistered event")
		}
	})

	// 测试Clear
	t.Run("TestClear", func(t *testing.T) {
		manager := NewEventManager(&EventOpt{PoolSize: 2, QueueSize: 10, Logger: logger})

		// 注册测试事件
		_ = manager.Register("testEvent", func(data interface{}) {})
		_ = manager.Register("testEvent2", func(data interface{}) {})

		// 清除所有事件
		manager.Clear()

		// 验证事件是否被清除
		if manager.HasEvent("testEvent") || manager.HasEvent("testEvent2") {
			t.Error("Clear should remove all events")
		}
	})

	// 测试ShutDown
	t.Run("TestShutDown", func(t *testing.T) {
		manager := NewEventManager(&EventOpt{PoolSize: 2, QueueSize: 10, Logger: logger})

		// 注册测试事件
		_ = manager.Register("testEvent", func(data interface{}) {})

		// 关闭管理器
		manager.ShutDown(time.Second)

		// 验证销毁标志
		if !manager.isDestroying() {
			t.Error("ShutDown should set destroying flag")
		}

		// 验证是否阻止新的事件注册
		err := manager.Register("newEvent", func(data interface{}) {})
		if err == nil {
			t.Error("Register should fail after ShutDown")
		}

		// 验证是否阻止事件触发
		err = manager.Trigger("testEvent", nil)
		if err == nil {
			t.Error("Trigger should fail after ShutDown")
		}
	})

	// 测试panic恢复
	t.Run("TestPanicRecovery", func(t *testing.T) {
		manager := NewEventManager(&EventOpt{PoolSize: 2, QueueSize: 10, Logger: logger})

		// 注册会panic的handler
		_ = manager.Register("panicEvent", func(data interface{}) {
			panic("test panic")
		})

		// 测试TriggerSync是否会捕获panic
		manager.TriggerSync("panicEvent", nil)

		// 测试Trigger是否会捕获panic
		_ = manager.Trigger("panicEvent", nil)
		time.Sleep(100 * time.Millisecond) // 等待goroutine执行
	})

	// 测试无logger情况下的panic恢复
	t.Run("TestPanicRecoveryWithoutLogger", func(t *testing.T) {
		manager := NewEventManager(&EventOpt{PoolSize: 2, QueueSize: 10}) // 不传入logger

		// 注册会panic的handler
		_ = manager.Register("panicEvent", func(data interface{}) {
			panic("test panic")
		})

		// 测试TriggerSync是否会捕获panic
		manager.TriggerSync("panicEvent", nil)
	})
}
