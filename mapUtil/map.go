package mapUtil

import "encoding/json"

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

func ToString[K comparable, V any](m map[K]V) string {
	marshal, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(marshal)
}
