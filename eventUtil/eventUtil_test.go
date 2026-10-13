package eventUtil

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func newTestEventManager(t *testing.T) *EventManager {
	t.Helper()
	manager := NewEventManager(EventOpt{PoolSize: 2, QueueSize: 10})
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

	for _, opt := range []EventOpt{
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
	handler := EventHandler{Func: func(interface{}) {}}
	manager.Register("event", handler)
	manager.Register("event", handler)
	if got := len(manager.handlers["event"]); got != 2 {
		t.Errorf("registered handler count = %d, want 2", got)
	}

	for _, test := range []struct {
		name      string
		eventName string
		handler   EventHandler
	}{
		{name: "empty name", eventName: "", handler: handler},
		{name: "nil handler", eventName: "event", handler: EventHandler{}},
		{name: "negative priority", eventName: "event", handler: EventHandler{Func: handler.Func, Priority: -1}},
		{name: "negative priority for new event", eventName: "negative", handler: EventHandler{Func: handler.Func, Priority: -1}},
		{name: "priority above maximum", eventName: "event", handler: EventHandler{Func: handler.Func, Priority: 1000}},
		{name: "priority above maximum for new event", eventName: "above-maximum", handler: EventHandler{Func: handler.Func, Priority: 1000}},
	} {
		t.Run(test.name, func(t *testing.T) {
			manager.Register(test.eventName, test.handler)
			if len(manager.handlers) != 1 || len(manager.handlers["event"]) != 2 {
				t.Errorf("invalid registration changed handlers: %#v", manager.handlers)
			}
		})
	}
}

func TestTriggerSync(t *testing.T) {
	manager := newTestEventManager(t)
	var got []string
	manager.Register("event", EventHandler{Func: func(data interface{}) {
		got = append(got, fmt.Sprintf("first:%v", data))
	}})
	manager.Register("event", EventHandler{Func: func(data interface{}) {
		got = append(got, fmt.Sprintf("second:%v", data))
	}})

	manager.TriggerSync("event", "payload")
	if want := []string{"first:payload", "second:payload"}; !reflect.DeepEqual(got, want) {
		t.Errorf("TriggerSync handler calls = %#v, want %#v", got, want)
	}
	manager.TriggerSync("missing", nil)
	if len(got) != 2 {
		t.Errorf("TriggerSync(missing) unexpectedly invoked handlers: %#v", got)
	}
}

func TestTriggerAsync(t *testing.T) {
	for _, mode := range []string{"Trigger", "TriggerWithId"} {
		t.Run(mode, func(t *testing.T) {
			manager := newTestEventManager(t)
			received := make(chan interface{}, 3)
			manager.Register("event", EventHandler{Func: func(data interface{}) {
				received <- data
			}})
			trigger := manager.Trigger
			if mode == "TriggerWithId" {
				trigger = func(eventName string, data interface{}) {
					manager.TriggerWithId(1, eventName, data)
				}
			}

			trigger("event", "payload")
			select {
			case got := <-received:
				if got != "payload" {
					t.Errorf("handler data = %v, want payload", got)
				}
			case <-time.After(time.Second):
				t.Fatal("timed out waiting for asynchronous handler")
			}

			trigger("", nil)
			trigger("missing", nil)
			manager.ShutDown(time.Second)
			select {
			case got := <-received:
				t.Errorf("invalid event unexpectedly invoked handler with %v", got)
			default:
			}
		})
	}
}

func TestTriggerPriorityOrder(t *testing.T) {
	for _, mode := range []string{"sync", "async", "asyncWithID"} {
		t.Run(mode, func(t *testing.T) {
			manager := newTestEventManager(t)
			entries := []struct {
				name     string
				priority int
			}{
				{name: "default"},
				{name: "second-default"},
				{name: "high", priority: 1},
				{name: "equal", priority: 1},
				{name: "explicit-default", priority: 999},
				{name: "before-default", priority: 998},
			}
			received := make(chan string, len(entries))
			for _, entry := range entries {
				manager.Register("event", EventHandler{
					Func:     func(interface{}) { received <- entry.name },
					Priority: entry.priority,
				})
			}

			switch mode {
			case "sync":
				manager.TriggerSync("event", nil)
			case "async":
				manager.Trigger("event", nil)
			case "asyncWithID":
				manager.TriggerWithId(1, "event", nil)
			}

			var got []string
			timer := time.NewTimer(time.Second)
			defer timer.Stop()
			for range entries {
				select {
				case name := <-received:
					got = append(got, name)
				case <-timer.C:
					t.Fatalf("timed out waiting for handlers; received %v", got)
				}
			}
			want := []string{"high", "equal", "before-default", "default", "second-default", "explicit-default"}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("handler order = %v, want %v", got, want)
			}
		})
	}
}

func TestRegisterDuringTrigger(t *testing.T) {
	manager := newTestEventManager(t)
	var got []int
	registered := false
	for _, priority := range []int{10, 20, 30, 40} {
		manager.Register("event", EventHandler{
			Func: func(interface{}) {
				got = append(got, priority)
				if priority == 10 && !registered {
					registered = true
					manager.Register("event", EventHandler{
						Func:     func(interface{}) { got = append(got, 1) },
						Priority: 1,
					})
				}
			},
			Priority: priority,
		})
	}

	manager.TriggerSync("event", nil)
	if want := []int{10, 20, 30, 40}; !reflect.DeepEqual(got, want) {
		t.Errorf("current trigger order = %v, want %v", got, want)
	}
	got = nil
	manager.TriggerSync("event", nil)
	if want := []int{1, 10, 20, 30, 40}; !reflect.DeepEqual(got, want) {
		t.Errorf("next trigger order = %v, want %v", got, want)
	}
}

func TestHasEventAndClear(t *testing.T) {
	manager := newTestEventManager(t)
	if manager.HasEvent("event") {
		t.Error("HasEvent() = true before registration")
	}
	manager.Register("event", EventHandler{Func: func(interface{}) {}})
	if !manager.HasEvent("event") {
		t.Error("HasEvent() = false after registration")
	}
	manager.Clear()
	if manager.HasEvent("event") || len(manager.handlers) != 0 {
		t.Errorf("Clear() left handlers: %#v", manager.handlers)
	}
}

func TestShutDownRejectsNewWork(t *testing.T) {
	manager := newTestEventManager(t)
	received := make(chan struct{}, 2)
	handler := EventHandler{Func: func(interface{}) { received <- struct{}{} }}
	manager.Register("event", handler)
	manager.ShutDown(time.Second)
	manager.ShutDown(time.Second)

	if !manager.isDestroying() {
		t.Error("ShutDown() did not set destroying")
	}
	manager.Register("new", handler)
	manager.Register("event", handler)
	if manager.HasEvent("new") || len(manager.handlers["event"]) != 1 {
		t.Errorf("Register(after shutdown) changed handlers: %#v", manager.handlers)
	}
	manager.Trigger("event", nil)
	manager.TriggerWithId(1, "event", nil)
	select {
	case <-received:
		t.Error("trigger after shutdown unexpectedly invoked handler")
	default:
	}
}

func TestHandlerPanicDoesNotEscapeOrStopFollowingHandlers(t *testing.T) {
	manager := newTestEventManager(t)
	called := false
	manager.Register("sync", EventHandler{Func: func(interface{}) { panic("test panic") }})
	manager.Register("sync", EventHandler{Func: func(interface{}) { called = true }})
	manager.TriggerSync("sync", nil)
	if !called {
		t.Error("handler after panicking handler was not called")
	}

	done := make(chan struct{}, 1)
	manager.Register("async", EventHandler{Func: func(interface{}) { panic("test panic") }})
	manager.Register("async", EventHandler{Func: func(interface{}) { done <- struct{}{} }})
	manager.TriggerWithId(1, "async", nil)
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("panicking asynchronous handler stopped the following handler")
	}
}
