package eventUtil

import (
	"errors"
	"sync"
	"testing"
	"time"
)

func TestNewEventManager(t *testing.T) {
	em, err := NewEventManager()
	if err != nil {
		t.Fatalf("Failed to create EventManager: %v", err)
	}

	if em.handlers == nil {
		t.Error("handlers map not initialized")
	}

	if em.pool == nil {
		t.Error("worker pool not initialized")
	}
}

func TestRegister(t *testing.T) {
	em, _ := NewEventManager()

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
	em, _ := NewEventManager()

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
	em, _ := NewEventManager()

	var wg sync.WaitGroup
	wg.Add(1)

	handlerCalled := false
	err := em.Register("test", func(data interface{}) {
		handlerCalled = true
		if data != "test data" {
			t.Errorf("Expected data: 'test data', got: %v", data)
		}
		wg.Done()
	})
	if err != nil {
		t.Fatalf("Failed to register handler: %v", err)
	}

	err = em.Trigger("test", "test data")
	if err != nil {
		t.Fatalf("Failed to trigger event: %v", err)
	}

	wg.Wait()
	if !handlerCalled {
		t.Error("Handler was not called")
	}
}

func TestTriggerInvalidEvent(t *testing.T) {
	em, _ := NewEventManager()

	err := em.Trigger("nonexistent", nil)
	if err == nil || err.Error() != "event not found" {
		t.Errorf("Expected 'event not found' error, got: %v", err)
	}
}

func TestTriggerDuringDestruction(t *testing.T) {
	em, _ := NewEventManager()

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
	em, _ := NewEventManager()

	if em.HasEvent("test") {
		t.Error("Event should not exist before registration")
	}

	em.Register("test", func(data interface{}) {})
	if !em.HasEvent("test") {
		t.Error("Event should exist after registration")
	}
}

func TestClear(t *testing.T) {
	em, _ := NewEventManager()

	em.Register("test1", func(data interface{}) {})
	em.Register("test2", func(data interface{}) {})

	em.Clear()

	if em.HasEvent("test1") || em.HasEvent("test2") {
		t.Error("Events should be cleared")
	}
}

func TestOnDestroy(t *testing.T) {
	em, _ := NewEventManager()

	// Register a handler
	em.Register("test", func(data interface{}) {})

	// Start destruction
	em.ShutDown(100 * time.Millisecond)

	// Verify destruction flag
	if !em.isDestroying() {
		t.Error("Destroying flag should be set")
	}

	// Verify pool is released by trying to trigger an event
	err := em.Trigger("test", nil)
	if err == nil || err.Error() != "event manager is destroying, cannot trigger events" {
		t.Errorf("Expected error after destruction, got: %v", err)
	}
}

func TestConcurrentAccess(t *testing.T) {
	em, _ := NewEventManager()
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
