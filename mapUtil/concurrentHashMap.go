package mapUtil

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
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

func (cm *ConcurrentHashMap[K, V]) Get(key K) V {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.m[key]
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

// Range 遍历元素（返回false可提前终止）
func (cm *ConcurrentHashMap[K, V]) Range(f func(key K, value V) bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	for k, v := range cm.m {
		if !f(k, v) {
			break
		}
	}
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

// MarshalJSON 实现 json.Marshaler 接口
func (cm *ConcurrentHashMap[K, V]) MarshalJSON() ([]byte, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return json.Marshal(cm.m) // 只序列化内部的 map
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (cm *ConcurrentHashMap[K, V]) UnmarshalJSON(data []byte) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.m = make(map[K]V)
	return json.Unmarshal(data, &cm.m)
}

// MarshalBSON 实现 bson.Marshaler 接口
func (cm *ConcurrentHashMap[K, V]) MarshalBSON() ([]byte, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return bson.Marshal(cm.m)
}

// UnmarshalBSON 实现 bson.Unmarshaler 接口
func (cm *ConcurrentHashMap[K, V]) UnmarshalBSON(data []byte) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.m = make(map[K]V)
	return bson.Unmarshal(data, &cm.m)
}
