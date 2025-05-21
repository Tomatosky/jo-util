package eventUtil

import (
	"errors"
	"sync"
)

// EventHandler 定义事件处理函数的类型
type EventHandler func(data interface{})

// EventManager 事件管理器结构
type EventManager struct {
	handlers map[string][]EventHandler
	lock     sync.RWMutex
}

// NewEventManager 创建一个新的事件管理器
func NewEventManager() *EventManager {
	return &EventManager{
		handlers: make(map[string][]EventHandler),
	}
}

// Register 注册事件处理函数
func (em *EventManager) Register(eventName string, handler EventHandler) {
	if eventName == "" {
		panic("event name cannot be empty")
	}
	if handler == nil {
		panic("handler cannot be nil")
	}

	em.lock.Lock()
	defer em.lock.Unlock()

	em.handlers[eventName] = append(em.handlers[eventName], handler)
}

// Unregister 取消注册事件处理函数
func (em *EventManager) Unregister(eventName string, handler EventHandler) error {
	if eventName == "" {
		return errors.New("event name cannot be empty")
	}
	if handler == nil {
		return errors.New("handler cannot be nil")
	}

	em.lock.Lock()
	defer em.lock.Unlock()

	handlers, exists := em.handlers[eventName]
	if !exists {
		return errors.New("event not found")
	}

	for i, h := range handlers {
		if &h == &handler { // 比较函数指针
			em.handlers[eventName] = append(handlers[:i], handlers[i+1:]...)
			return nil
		}
	}

	return errors.New("handler not found for the event")
}

// Trigger 触发事件
func (em *EventManager) Trigger(eventName string, data interface{}) {
	go em.TriggerSync(eventName, data)
}

// TriggerSync 同步触发事件
func (em *EventManager) TriggerSync(eventName string, data interface{}) error {
	if eventName == "" {
		return errors.New("event name cannot be empty")
	}

	em.lock.RLock()
	defer em.lock.RUnlock()

	handlers, exists := em.handlers[eventName]
	if !exists {
		return errors.New("event not found")
	}

	// 同步执行所有处理函数
	for _, handler := range handlers {
		handler(data)
	}

	return nil
}

// HasEvent 检查事件是否存在
func (em *EventManager) HasEvent(eventName string) bool {
	em.lock.RLock()
	defer em.lock.RUnlock()

	_, exists := em.handlers[eventName]
	return exists
}

// Clear 清除所有事件处理函数
func (em *EventManager) Clear() {
	em.lock.Lock()
	defer em.lock.Unlock()

	em.handlers = make(map[string][]EventHandler)
}
