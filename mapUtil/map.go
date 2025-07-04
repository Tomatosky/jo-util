package mapUtil

import (
	"encoding/json"
	"golang.org/x/exp/constraints"
	"sort"
)

func ContainsKey[K comparable, V any](m map[K]V, key K) bool {
	if _, ok := m[key]; ok {
		return true
	}
	return false
}

func ContainValue[K comparable, V comparable](m map[K]V, value V) bool {
	for _, v := range m {
		if v == value {
			return true
		}
	}
	return false
}

func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func GetOrDefault[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	if value, ok := m[key]; ok {
		return value
	}
	return defaultValue
}

func PutIfAbsent[K comparable, V any](m map[K]V, key K, defaultValue V) {
	if _, ok := m[key]; !ok {
		m[key] = defaultValue
	}
}

func ToString[K comparable, V any](m map[K]V) string {
	marshal, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(marshal)
}

// SortByValue 根据map的值排序并返回键列表
// m: 要排序的map
// reverse: true为降序，false为升序
func SortByValue[K comparable, V constraints.Ordered](m map[K]V, reverse bool) []K {
	// 创建一个切片来保存map的键
	keys := make([]K, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	// 根据reverse参数决定排序方向
	if reverse {
		// 降序排序
		sort.Slice(keys, func(i, j int) bool {
			return m[keys[i]] > m[keys[j]]
		})
	} else {
		// 升序排序
		sort.Slice(keys, func(i, j int) bool {
			return m[keys[i]] < m[keys[j]]
		})
	}
	return keys
}
