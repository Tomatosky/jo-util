package eventUtil

import (
	"errors"
	"github.com/Tomatosky/jo-util/poolUtil"
	"github.com/panjf2000/ants/v2"
	"sync"
	"time"
)

// EventHandler 定义事件处理函数的类型
type EventHandler func(data interface{})

// EventManager 事件管理器结构
type EventManager struct {
	handlers   map[string][]EventHandler
	pool       *poolUtil.Pool
	lock       sync.RWMutex
	destroying bool
}

// NewEventManager 创建一个新的事件管理器
func NewEventManager() (*EventManager, error) {
	return &EventManager{
		handlers: make(map[string][]EventHandler),
		pool:     poolUtil.NewPool(ants.DefaultAntsPoolSize),
	}, nil
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

// Trigger 触发事件
func (em *EventManager) Trigger(eventName string, data interface{}) error {
	if em.isDestroying() {
		return errors.New("event manager is destroying, cannot trigger events")
	}

	if eventName == "" {
		return errors.New("event name cannot be empty")
	}

	em.lock.RLock()
	_, exists := em.handlers[eventName]
	em.lock.RUnlock()

	if !exists {
		return errors.New("event not found")
	}

	err := em.pool.Submit(func() {
		em.triggerSync(eventName, data)
	})
	if err != nil {
		return err
	}
	return nil
}

// triggerSync 同步触发事件
func (em *EventManager) triggerSync(eventName string, data interface{}) {
	em.lock.RLock()
	handlers, _ := em.handlers[eventName]
	em.lock.RUnlock()

	// 同步执行所有处理函数
	for _, handler := range handlers {
		handler(data)
	}
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

// ShutDown 销毁事件管理器
func (em *EventManager) ShutDown(timeout time.Duration) {
	defer em.pool.Release(timeout)

	em.lock.Lock()
	// 设置销毁标志，阻止新的事件注册和触发
	em.destroying = true
	em.lock.Unlock()
}

// isDestroying 检查是否正在销毁
func (em *EventManager) isDestroying() bool {
	em.lock.RLock()
	defer em.lock.RUnlock()
	return em.destroying
}
