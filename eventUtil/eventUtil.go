package eventUtil

import (
	"errors"
	"fmt"
	"github.com/Tomatosky/jo-util/poolUtil"
	"github.com/Tomatosky/jo-util/randomUtil"
	"go.uber.org/zap"
	"sync"
	"time"
)

// EventHandler 定义事件处理函数的类型
type EventHandler func(data interface{})

// EventManager 事件管理器结构
type EventManager struct {
	handlers   map[string][]EventHandler
	pool       *poolUtil.IdPool
	lock       sync.RWMutex
	destroying bool
	logger     *zap.Logger
}

type EventOpt struct {
	PoolSize  int32
	QueueSize int
	Logger    *zap.Logger
}

// NewEventManager 创建一个新的事件管理器
func NewEventManager(opt *EventOpt) *EventManager {
	if opt.PoolSize <= 0 || opt.QueueSize <= 0 {
		panic("pool size and queue size must be greater than 0")
	}

	manager := &EventManager{
		handlers: make(map[string][]EventHandler),
		logger:   opt.Logger,
		pool: poolUtil.NewIdPool(&poolUtil.IdPoolOpt{
			PoolSize:  opt.PoolSize,
			QueueSize: opt.QueueSize,
			Logger:    opt.Logger,
			PoolName:  "event",
		}),
	}
	return manager
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

func (em *EventManager) Trigger(eventName string, data interface{}) error {
	return em.TriggerWithId(int32(randomUtil.RandomInt(1, 100000)), eventName, data)
}

// TriggerWithId 触发事件
func (em *EventManager) TriggerWithId(id int32, eventName string, data interface{}) error {
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

	em.pool.Submit(func() {
		em.TriggerSync(eventName, data)
	})
	return nil
}

// TriggerSync 同步触发事件
func (em *EventManager) TriggerSync(eventName string, data interface{}) {
	em.lock.RLock()
	handlers, _ := em.handlers[eventName]
	em.lock.RUnlock()

	// 同步执行所有处理函数
	for _, handler := range handlers {
		em.triggerHandle(handler, data)
	}
}

// triggerHandle 触发事件
func (em *EventManager) triggerHandle(f EventHandler, data interface{}) {
	defer func() {
		err := recover()
		if err != nil {
			if em.logger != nil {
				em.logger.Error(fmt.Sprintf("error: %v", err))
			} else {
				fmt.Printf("panic: %v", err)
			}
		}
	}()

	f(data)
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
	defer em.pool.Shutdown(timeout)

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
