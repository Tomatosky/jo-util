package mapUtil

import "encoding/json"

func ContainsKey[K comparable, V any](m map[K]V, key K) bool {
	if _, ok := m[key]; ok {
		return true
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
