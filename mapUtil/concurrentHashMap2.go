package mapUtil

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Tomatosky/jo-util/logger"
	"go.mongodb.org/mongo-driver/bson"
)

var _ IMap[string, int] = (*ConcurrentHashMap2[string, int])(nil)

type ConcurrentHashMap2[K comparable, V any] struct {
	m sync.Map
}

func NewConcurrentHashMap2[K comparable, V any](initMap ...map[K]V) *ConcurrentHashMap2[K, V] {
	cm := &ConcurrentHashMap2[K, V]{}

	if len(initMap) > 0 && initMap[0] != nil {
		for k, v := range initMap[0] {
			cm.m.Store(k, v)
		}
	}
	return cm
}

func (cm *ConcurrentHashMap2[K, V]) Get(key K) V {
	value, ok := cm.m.Load(key)
	if !ok {
		var zero V
		return zero
	}
	return value.(V)
}

func (cm *ConcurrentHashMap2[K, V]) Put(key K, value V) {
	cm.m.Store(key, value)
}

func (cm *ConcurrentHashMap2[K, V]) Remove(key K) {
	cm.m.Delete(key)
}

func (cm *ConcurrentHashMap2[K, V]) Size() int {
	size := 0
	cm.m.Range(func(_, _ interface{}) bool {
		size++
		return true
	})
	return size
}

func (cm *ConcurrentHashMap2[K, V]) ContainsKey(key K) bool {
	_, ok := cm.m.Load(key)
	return ok
}

func (cm *ConcurrentHashMap2[K, V]) Clear() {
	cm.m = sync.Map{}
}

func (cm *ConcurrentHashMap2[K, V]) Keys() []K {
	keys := make([]K, 0)
	cm.m.Range(func(key, _ interface{}) bool {
		keys = append(keys, key.(K))
		return true
	})
	return keys
}

func (cm *ConcurrentHashMap2[K, V]) Values() []V {
	values := make([]V, 0)
	cm.m.Range(func(_, value interface{}) bool {
		values = append(values, value.(V))
		return true
	})
	return values
}

func (cm *ConcurrentHashMap2[K, V]) PutIfAbsent(key K, value V) (existing V, loaded bool) {
	actual, loaded := cm.m.LoadOrStore(key, value)
	return actual.(V), loaded
}

func (cm *ConcurrentHashMap2[K, V]) GetOrDefault(key K, defaultValue V) V {
	if value, ok := cm.m.Load(key); ok {
		return value.(V)
	}
	return defaultValue
}

func (cm *ConcurrentHashMap2[K, V]) ToMap() map[K]V {
	result := make(map[K]V)
	cm.m.Range(func(key, value interface{}) bool {
		result[key.(K)] = value.(V)
		return true
	})
	return result
}

func (cm *ConcurrentHashMap2[K, V]) Range(f func(key K, value V) bool) {
	cm.m.Range(func(key, value interface{}) bool {
		return f(key.(K), value.(V))
	})
}

func (cm *ConcurrentHashMap2[K, V]) ToString() string {
	m := cm.ToMap()
	bytes, err := json.Marshal(m)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("%v", err))
		panic(err)
	}
	return string(bytes)
}

func (cm *ConcurrentHashMap2[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(cm.ToMap())
}

func (cm *ConcurrentHashMap2[K, V]) UnmarshalJSON(data []byte) error {
	var m map[K]V
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	cm.m = sync.Map{}
	for k, v := range m {
		cm.m.Store(k, v)
	}
	return nil
}

func (cm *ConcurrentHashMap2[K, V]) MarshalBSON() ([]byte, error) {
	return bson.Marshal(cm.ToMap())
}

func (cm *ConcurrentHashMap2[K, V]) UnmarshalBSON(data []byte) error {
	var m map[K]V
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}
	cm.m = sync.Map{}
	for k, v := range m {
		cm.m.Store(k, v)
	}
	return nil
}
