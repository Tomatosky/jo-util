package eventUtil

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/Tomatosky/jo-util/logger"
	"github.com/Tomatosky/jo-util/poolUtil"
	"github.com/Tomatosky/jo-util/randomUtil"
)

var Inst *EventManager

// EventHandler 定义事件处理函数及其优先级。
type EventHandler struct {
	Func     func(data interface{})
	Priority int // 数字越小越先执行；未设置或为 0 时默认为 999。
}

// EventManager 事件管理器结构
type EventManager struct {
	handlers   map[string][]EventHandler
	pool       *poolUtil.IdPool
	lock       sync.RWMutex
	destroying bool
}

type EventOpt struct {
	PoolSize  int64
	QueueSize int
}

// NewEventManager 创建一个新的事件管理器
func NewEventManager(opt *EventOpt) *EventManager {
	if opt.PoolSize <= 0 || opt.QueueSize <= 0 {
		logger.Log.Error(fmt.Sprintf("pool size and queue size must be greater than 0"))
		panic(fmt.Sprintf("pool size and queue size must be greater than 0"))
	}

	manager := &EventManager{
		handlers: make(map[string][]EventHandler),
		pool: poolUtil.NewIdPool(&poolUtil.IdPoolOpt{
			PoolSize:  opt.PoolSize,
			QueueSize: opt.QueueSize,
			PoolName:  "event",
		}),
	}
	return manager
}

// Register 注册事件处理函数，相同优先级按注册顺序执行。
func (em *EventManager) Register(eventName string, handler EventHandler) {
	if eventName == "" {
		logger.Log.Error(fmt.Sprintf("event name must not be empty"))
		return
	}
	if handler.Func == nil {
		logger.Log.Error(fmt.Sprintf("event handler must not be nil"))
		return
	}
	if handler.Priority < 0 {
		logger.Log.Error(fmt.Sprintf("event handler priority must not be negative"))
		return
	}
	if handler.Priority > 999 {
		logger.Log.Error(fmt.Sprintf("event handler priority must be less than 999"))
		return
	}
	if handler.Priority == 0 {
		handler.Priority = 999
	}

	em.lock.Lock()
	defer em.lock.Unlock()

	if em.destroying {
		logger.Log.Warn(fmt.Sprintf("event handler is already destroying"))
		return
	}

	handlers := em.handlers[eventName]
	index := sort.Search(len(handlers), func(i int) bool {
		return handlers[i].Priority > handler.Priority
	})
	// 使用新切片插入，避免改变正在触发的处理函数列表。
	ordered := make([]EventHandler, len(handlers)+1)
	copy(ordered, handlers[:index])
	ordered[index] = handler
	copy(ordered[index+1:], handlers[index:])
	em.handlers[eventName] = ordered
}

func (em *EventManager) Trigger(eventName string, data interface{}) {
	em.TriggerWithId(int32(randomUtil.RandomInt(1, 100000)), eventName, data)
}

// TriggerWithId 触发事件
func (em *EventManager) TriggerWithId(id int32, eventName string, data interface{}) {
	if em.isDestroying() {
		logger.Log.Warn(fmt.Sprintf("event handler is already destroying"))
		return
	}

	if eventName == "" {
		logger.Log.Error(fmt.Sprintf("event handler must not be empty"))
		return
	}

	em.lock.RLock()
	_, exists := em.handlers[eventName]
	em.lock.RUnlock()

	if !exists {
		logger.Log.Error(fmt.Sprintf("event handler not found"))
		return
	}

	em.pool.SubmitWithId(id, func() {
		em.TriggerSync(eventName, data)
	})
}

// TriggerSync 同步触发事件
func (em *EventManager) TriggerSync(eventName string, data interface{}) {
	em.lock.RLock()
	handlers, _ := em.handlers[eventName]
	em.lock.RUnlock()

	// 同步执行所有处理函数
	for _, handler := range handlers {
		em.triggerHandle(handler.Func, data)
	}
}

// triggerHandle 触发事件
func (em *EventManager) triggerHandle(f func(data interface{}), data interface{}) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Log.Error(fmt.Sprintf("error: %v", err))
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
