package mapUtil

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Tomatosky/jo-util/logger"
	"go.mongodb.org/mongo-driver/bson"
)

var _ IMap[string, int] = (*BiMap[string, int])(nil)

// BiMap 双向映射，支持根据key查找value和根据value查找key
// K和V都需要是可比较类型
type BiMap[K comparable, V comparable] struct {
	mu      sync.RWMutex
	forward map[K]V // key -> value
	inverse map[V]K // value -> key
}

// NewBiMap 创建一个新的BiMap
func NewBiMap[K comparable, V comparable](initMap ...map[K]V) *BiMap[K, V] {
	bm := &BiMap[K, V]{
		forward: make(map[K]V),
		inverse: make(map[V]K),
	}

	// 如果有传入初始化map
	if len(initMap) > 0 && initMap[0] != nil {
		bm.mu.Lock()
		defer bm.mu.Unlock()

		// 深拷贝原始map内容，同时建立反向映射
		for k, v := range initMap[0] {
			// 检查value是否已存在，如果存在则移除旧的key映射
			if oldKey, exists := bm.inverse[v]; exists {
				delete(bm.forward, oldKey)
			}
			bm.forward[k] = v
			bm.inverse[v] = k
		}
	}
	return bm
}

// Get 根据key获取value
func (bm *BiMap[K, V]) Get(key K) V {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	return bm.forward[key]
}

// GetKey 根据value获取key（BiMap特有方法）
func (bm *BiMap[K, V]) GetKey(value V) K {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	return bm.inverse[value]
}

// Put 设置key-value映射
// 如果key已存在，会覆盖旧的映射
// 如果value已存在，会移除旧的key映射
func (bm *BiMap[K, V]) Put(key K, value V) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	// 检查key是否已存在，如果存在则移除旧的value反向映射
	if oldValue, exists := bm.forward[key]; exists {
		delete(bm.inverse, oldValue)
	}

	// 检查value是否已存在，如果存在则移除旧的key映射
	if oldKey, exists := bm.inverse[value]; exists {
		delete(bm.forward, oldKey)
	}

	// 建立新的双向映射
	bm.forward[key] = value
	bm.inverse[value] = key
}

// Remove 移除指定key的映射
func (bm *BiMap[K, V]) Remove(key K) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if value, exists := bm.forward[key]; exists {
		delete(bm.forward, key)
		delete(bm.inverse, value)
	}
}

// RemoveValue 根据value移除映射（BiMap特有方法）
func (bm *BiMap[K, V]) RemoveValue(value V) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if key, exists := bm.inverse[value]; exists {
		delete(bm.inverse, value)
		delete(bm.forward, key)
	}
}

// Size 返回映射数量
func (bm *BiMap[K, V]) Size() int {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	return len(bm.forward)
}

// ContainsKey 检查key是否存在
func (bm *BiMap[K, V]) ContainsKey(key K) bool {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	_, ok := bm.forward[key]
	return ok
}

// ContainsValue 检查value是否存在（BiMap特有方法）
func (bm *BiMap[K, V]) ContainsValue(value V) bool {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	_, ok := bm.inverse[value]
	return ok
}

// Clear 清空所有映射
func (bm *BiMap[K, V]) Clear() {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	bm.forward = make(map[K]V)
	bm.inverse = make(map[V]K)
}

// Keys 返回所有key的切片
func (bm *BiMap[K, V]) Keys() []K {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	keys := make([]K, 0, len(bm.forward))
	for k := range bm.forward {
		keys = append(keys, k)
	}
	return keys
}

// Values 返回所有value的切片
func (bm *BiMap[K, V]) Values() []V {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	values := make([]V, 0, len(bm.forward))
	for _, v := range bm.forward {
		values = append(values, v)
	}
	return values
}

// PutIfAbsent 如果key不存在则设置，返回现有值和是否已存在
func (bm *BiMap[K, V]) PutIfAbsent(key K, value V) (existing V, loaded bool) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if existing, loaded = bm.forward[key]; !loaded {
		// 检查value是否已存在，如果存在则移除旧的key映射
		if oldKey, exists := bm.inverse[value]; exists {
			delete(bm.forward, oldKey)
		}
		bm.forward[key] = value
		bm.inverse[value] = key
		existing = value
	}
	return
}

// GetOrDefault 获取key对应的value，如果不存在返回默认值
func (bm *BiMap[K, V]) GetOrDefault(key K, defaultValue V) V {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	if value, ok := bm.forward[key]; ok {
		return value
	}
	return defaultValue
}

// ToMap 转换为普通map
func (bm *BiMap[K, V]) ToMap() map[K]V {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	result := make(map[K]V, len(bm.forward))
	for k, v := range bm.forward {
		result[k] = v
	}
	return result
}

// Range 遍历元素（返回false可提前终止）
func (bm *BiMap[K, V]) Range(f func(key K, value V) bool) {
	// 先复制数据再遍历
	bm.mu.RLock()
	tmp := make(map[K]V, len(bm.forward))
	for k, v := range bm.forward {
		tmp[k] = v
	}
	bm.mu.RUnlock()

	// 遍历复制后的数据，不需要锁
	for k, v := range tmp {
		if !f(k, v) {
			break
		}
	}
}

// ToString 转换为JSON字符串
func (bm *BiMap[K, V]) ToString() string {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	bytes, err := json.Marshal(bm.forward)
	if err != nil {
		logger.Log.Fatal(fmt.Sprintf("%v", err))
	}
	return string(bytes)
}

// MarshalJSON 实现 json.Marshaler 接口
func (bm *BiMap[K, V]) MarshalJSON() ([]byte, error) {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	return json.Marshal(bm.forward) // 只序列化正向映射
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (bm *BiMap[K, V]) UnmarshalJSON(data []byte) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	tempMap := make(map[K]V)
	if err := json.Unmarshal(data, &tempMap); err != nil {
		return err
	}

	// 重建双向映射
	bm.forward = make(map[K]V)
	bm.inverse = make(map[V]K)
	for k, v := range tempMap {
		// 检查value是否已存在，如果存在则移除旧的key映射
		if oldKey, exists := bm.inverse[v]; exists {
			delete(bm.forward, oldKey)
		}
		bm.forward[k] = v
		bm.inverse[v] = k
	}
	return nil
}

// MarshalBSON 实现 bson.Marshaler 接口
func (bm *BiMap[K, V]) MarshalBSON() ([]byte, error) {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	return bson.Marshal(bm.forward)
}

// UnmarshalBSON 实现 bson.Unmarshaler 接口
func (bm *BiMap[K, V]) UnmarshalBSON(data []byte) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	tempMap := make(map[K]V)
	if err := bson.Unmarshal(data, &tempMap); err != nil {
		return err
	}

	// 重建双向映射
	bm.forward = make(map[K]V)
	bm.inverse = make(map[V]K)
	for k, v := range tempMap {
		// 检查value是否已存在，如果存在则移除旧的key映射
		if oldKey, exists := bm.inverse[v]; exists {
			delete(bm.forward, oldKey)
		}
		bm.forward[k] = v
		bm.inverse[v] = k
	}
	return nil
}
