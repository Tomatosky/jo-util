package maputil

import (
	"encoding/json"
	"sync"
)

type ConcurrentHashMap[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]V
}

func NewConcurrentHashMap[K comparable, V any](initMap ...map[K]V) *ConcurrentHashMap[K, V] {
	cm := &ConcurrentHashMap[K, V]{
		m: make(map[K]V),
	}

	// 如果有传入初始化map
	if len(initMap) > 0 && initMap[0] != nil {
		// 使用写锁确保线程安全
		cm.mu.Lock()
		defer cm.mu.Unlock()

		// 深拷贝原始map内容
		for k, v := range initMap[0] {
			cm.m[k] = v
		}
	}
	return cm
}

func (cm *ConcurrentHashMap[K, V]) Get(key K) (V, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	value, ok := cm.m[key]
	return value, ok
}

func (cm *ConcurrentHashMap[K, V]) Put(key K, value V) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.m[key] = value
}

func (cm *ConcurrentHashMap[K, V]) Remove(key K) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.m, key)
}

func (cm *ConcurrentHashMap[K, V]) Size() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return len(cm.m)
}

func (cm *ConcurrentHashMap[K, V]) ContainsKey(key K) bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	_, ok := cm.m[key]
	return ok
}

func (cm *ConcurrentHashMap[K, V]) Clear() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.m = make(map[K]V)
}

func (cm *ConcurrentHashMap[K, V]) Keys() []K {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	keys := make([]K, 0, len(cm.m))
	for k := range cm.m {
		keys = append(keys, k)
	}
	return keys
}

func (cm *ConcurrentHashMap[K, V]) Values() []V {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	values := make([]V, 0, len(cm.m))
	for _, v := range cm.m {
		values = append(values, v)
	}
	return values
}

func (cm *ConcurrentHashMap[K, V]) PutIfAbsent(key K, value V) (existing V, loaded bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if existing, loaded = cm.m[key]; !loaded {
		cm.m[key] = value
	}
	return
}

func (cm *ConcurrentHashMap[K, V]) GetOrDefault(key K, defaultValue V) V {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	if value, ok := cm.m[key]; ok {
		return value
	}
	return defaultValue
}

func (cm *ConcurrentHashMap[K, V]) ToMap() map[K]V {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.m
}

func (cm *ConcurrentHashMap[K, V]) ToString() string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	bytes, err := json.Marshal(cm.m)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
