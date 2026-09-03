package eventUtil

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func newTestEventManager(t *testing.T) *EventManager {
	t.Helper()
	manager := NewEventManager(&EventOpt{PoolSize: 2, QueueSize: 10})
	t.Cleanup(func() {
		manager.ShutDown(time.Second)
	})
	return manager
}

func panicValue(f func()) (recovered any) {
	defer func() {
		recovered = recover()
	}()
	f()
	return nil
}

func TestNewEventManager(t *testing.T) {
	manager := newTestEventManager(t)
	if manager.handlers == nil || manager.pool == nil {
		t.Fatalf("NewEventManager() returned an incompletely initialized manager: %#v", manager)
	}
	if manager.destroying {
		t.Error("new manager is already destroying")
	}

	for _, opt := range []*EventOpt{
		{PoolSize: 0, QueueSize: 1},
		{PoolSize: 1, QueueSize: 0},
		{PoolSize: -1, QueueSize: -1},
	} {
		if got := panicValue(func() { NewEventManager(opt) }); got != "pool size and queue size must be greater than 0" {
			t.Errorf("NewEventManager(%+v) panic = %v, want validation message", opt, got)
		}
	}
}

func TestRegister(t *testing.T) {
	manager := newTestEventManager(t)
	handler := func(interface{}) {}
	if err := manager.Register("event", handler); err != nil {
		t.Fatalf("Register(valid): %v", err)
	}
	if err := manager.Register("event", handler); err != nil {
		t.Fatalf("Register(second handler): %v", err)
	}
	if got := len(manager.handlers["event"]); got != 2 {
		t.Errorf("registered handler count = %d, want 2", got)
	}

	if err := manager.Register("", handler); err == nil || err.Error() != "event name cannot be empty" {
		t.Errorf("Register(empty name) error = %v", err)
	}
	if err := manager.Register("event", nil); err == nil || err.Error() != "handler cannot be nil" {
		t.Errorf("Register(nil handler) error = %v", err)
	}
}

func TestTriggerSync(t *testing.T) {
	manager := newTestEventManager(t)
	var got []string
	if err := manager.Register("event", func(data interface{}) {
		got = append(got, fmt.Sprintf("first:%v", data))
	}); err != nil {
		t.Fatal(err)
	}
	if err := manager.Register("event", func(data interface{}) {
		got = append(got, fmt.Sprintf("second:%v", data))
	}); err != nil {
		t.Fatal(err)
	}

	manager.TriggerSync("event", "payload")
	if want := []string{"first:payload", "second:payload"}; !reflect.DeepEqual(got, want) {
		t.Errorf("TriggerSync handler calls = %#v, want %#v", got, want)
	}
	manager.TriggerSync("missing", nil)
	if len(got) != 2 {
		t.Errorf("TriggerSync(missing) unexpectedly invoked handlers: %#v", got)
	}
}

func TestTriggerWithID(t *testing.T) {
	manager := newTestEventManager(t)
	received := make(chan interface{}, 1)
	if err := manager.Register("event", func(data interface{}) {
		received <- data
	}); err != nil {
		t.Fatal(err)
	}

	if err := manager.TriggerWithId(1, "event", "payload"); err != nil {
		t.Fatalf("TriggerWithId(valid): %v", err)
	}
	select {
	case got := <-received:
		if got != "payload" {
			t.Errorf("handler data = %v, want payload", got)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for asynchronous handler")
	}

	if err := manager.TriggerWithId(1, "", nil); err == nil || err.Error() != "event name cannot be empty" {
		t.Errorf("TriggerWithId(empty name) error = %v", err)
	}
	if err := manager.TriggerWithId(1, "missing", nil); err == nil || err.Error() != "event not found" {
		t.Errorf("TriggerWithId(missing) error = %v", err)
	}
}

func TestHasEventAndClear(t *testing.T) {
	manager := newTestEventManager(t)
	if manager.HasEvent("event") {
		t.Error("HasEvent() = true before registration")
	}
	if err := manager.Register("event", func(interface{}) {}); err != nil {
		t.Fatal(err)
	}
	if !manager.HasEvent("event") {
		t.Error("HasEvent() = false after registration")
	}
	manager.Clear()
	if manager.HasEvent("event") || len(manager.handlers) != 0 {
		t.Errorf("Clear() left handlers: %#v", manager.handlers)
	}
}

func TestShutDownRejectsNewWork(t *testing.T) {
	manager := NewEventManager(&EventOpt{PoolSize: 1, QueueSize: 1})
	if err := manager.Register("event", func(interface{}) {}); err != nil {
		t.Fatal(err)
	}
	manager.ShutDown(time.Second)
	manager.ShutDown(time.Second)

	if !manager.isDestroying() {
		t.Error("ShutDown() did not set destroying")
	}
	if err := manager.Register("new", func(interface{}) {}); err == nil || err.Error() != "event manager is destroying, cannot register new handlers" {
		t.Errorf("Register(after shutdown) error = %v", err)
	}
	if err := manager.TriggerWithId(1, "event", nil); err == nil || err.Error() != "event manager is destroying, cannot trigger events" {
		t.Errorf("TriggerWithId(after shutdown) error = %v", err)
	}
}

func TestHandlerPanicDoesNotEscapeOrStopFollowingHandlers(t *testing.T) {
	manager := newTestEventManager(t)
	called := false
	if err := manager.Register("sync", func(interface{}) { panic("test panic") }); err != nil {
		t.Fatal(err)
	}
	if err := manager.Register("sync", func(interface{}) { called = true }); err != nil {
		t.Fatal(err)
	}
	manager.TriggerSync("sync", nil)
	if !called {
		t.Error("handler after panicking handler was not called")
	}

	done := make(chan struct{}, 1)
	if err := manager.Register("async", func(interface{}) { panic("test panic") }); err != nil {
		t.Fatal(err)
	}
	if err := manager.Register("async", func(interface{}) { done <- struct{}{} }); err != nil {
		t.Fatal(err)
	}
	if err := manager.TriggerWithId(1, "async", nil); err != nil {
		t.Fatal(err)
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("panicking asynchronous handler stopped the following handler")
	}
}
