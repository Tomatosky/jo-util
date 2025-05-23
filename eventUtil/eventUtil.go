package eventUtil

import (
	"errors"
	"sync"
)

// EventHandler 定义事件处理函数的类型
type EventHandler func(data interface{})

// EventManager 事件管理器结构
type EventManager struct {
	handlers   map[string][]EventHandler
	lock       sync.RWMutex
	destroying bool
	wg         sync.WaitGroup // 用于等待所有事件处理完成
}

// NewEventManager 创建一个新的事件管理器
func NewEventManager() *EventManager {
	return &EventManager{
		handlers: make(map[string][]EventHandler),
	}
}

// Register 注册事件处理函数
func (em *EventManager) Register(eventName string, handler EventHandler) error {
	if eventName == "" {
		return errors.New("event name cannot be empty")
	}
	if handler == nil {
		return errors.New("handler cannot be nil")
	}

	em.lock.Lock()
	defer em.lock.Unlock()

	if em.destroying {
		return errors.New("event manager is destroying, cannot register new handlers")
	}

	em.handlers[eventName] = append(em.handlers[eventName], handler)
	return nil
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

	if em.destroying {
		return errors.New("event manager is destroying, cannot unregister handlers")
	}

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
func (em *EventManager) Trigger(eventName string, data interface{}) error {
	if em.isDestroying() {
		return errors.New("event manager is destroying, cannot trigger events")
	}

	em.wg.Add(1)
	go func() {
		defer em.wg.Done()
		em.triggerSync(eventName, data)
	}()
	return nil
}

// triggerSync 同步触发事件
func (em *EventManager) triggerSync(eventName string, data interface{}) error {
	if eventName == "" {
		return errors.New("event name cannot be empty")
	}

	if em.isDestroying() {
		return errors.New("event manager is destroying, cannot trigger events")
	}

	em.lock.RLock()
	handlers, exists := em.handlers[eventName]
	em.lock.RUnlock()

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

// OnDestroy 销毁事件管理器
func (em *EventManager) OnDestroy() {
	em.lock.Lock()
	// 设置销毁标志，阻止新的事件注册和触发
	em.destroying = true
	em.lock.Unlock()

	// 等待所有正在处理的事件完成
	em.wg.Wait()

	// 清除所有事件处理函数
	em.Clear()
}

// isDestroying 检查是否正在销毁
func (em *EventManager) isDestroying() bool {
	em.lock.RLock()
	defer em.lock.RUnlock()
	return em.destroying
}
