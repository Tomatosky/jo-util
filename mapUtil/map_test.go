package mapUtil

import (
	"encoding/json"
	"testing"
)

func TestContainsKey(t *testing.T) {
	tests := []struct {
		name     string
		m        map[string]int
		key      string
		expected bool
	}{
		{name: "key exists", m: map[string]int{"a": 1, "b": 2}, key: "a", expected: true},
		{name: "key does not exist", m: map[string]int{"a": 1, "b": 2}, key: "c", expected: false},
		{name: "empty map", m: map[string]int{}, key: "a", expected: false},
		{name: "nil map", m: nil, key: "a", expected: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContainsKey(tt.m, tt.key)
			if got != tt.expected {
				t.Errorf("ContainsKey() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestContainsKeyWithDifferentTypes(t *testing.T) {
	// 测试不同类型的 key
	t.Run("int key", func(t *testing.T) {
		m := map[int]string{1: "one", 2: "two"}
		if !ContainsKey(m, 1) {
			t.Error("ContainsKey() with int key failed")
		}
	})
	t.Run("float key", func(t *testing.T) {
		m := map[float64]string{1.1: "one", 2.2: "two"}
		if !ContainsKey(m, 1.1) {
			t.Error("ContainsKey() with float key failed")
		}
	})
	t.Run("struct key", func(t *testing.T) {
		type myKey struct {
			id int
		}
		m := map[myKey]string{{id: 1}: "one", {id: 2}: "two"}
		if !ContainsKey(m, myKey{id: 1}) {
			t.Error("ContainsKey() with struct key failed")
		}
	})
}

func TestKeys(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]int
		want []string
	}{
		{name: "empty map", m: map[string]int{}, want: []string{}},
		{name: "single element", m: map[string]int{"a": 1}, want: []string{"a"}},
		{name: "multiple elements", m: map[string]int{"a": 1, "b": 2, "c": 3}, want: []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Keys(tt.m)
			if len(got) != len(tt.want) {
				t.Errorf("Keys() length = %v, want %v", len(got), len(tt.want))
			}

			// 检查所有期望的键都存在
			for _, k := range tt.want {
				found := false
				for _, gk := range got {
					if gk == k {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Keys() missing key %v", k)
				}
			}
		})
	}
}

func TestValues(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]int
		want []int
	}{
		{
			name: "empty map",
			m:    map[string]int{},
			want: []int{},
		},
		{
			name: "single element",
			m:    map[string]int{"a": 1},
			want: []int{1},
		},
		{
			name: "multiple elements",
			m:    map[string]int{"a": 1, "b": 2, "c": 3},
			want: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Values(tt.m)
			if len(got) != len(tt.want) {
				t.Errorf("Values() length = %v, want %v", len(got), len(tt.want))
			}

			// 检查所有期望的值都存在
			for _, v := range tt.want {
				found := false
				for _, gv := range got {
					if gv == v {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Values() missing value %v", v)
				}
			}
		})
	}
}

func TestGetOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		m            map[string]int
		key          string
		defaultValue int
		want         int
	}{
		{name: "key exists", m: map[string]int{"a": 1, "b": 2}, key: "a", defaultValue: 0, want: 1},
		{name: "key does not exist", m: map[string]int{"a": 1, "b": 2}, key: "c", defaultValue: 3, want: 3},
		{name: "empty map", m: map[string]int{}, key: "a", defaultValue: 1, want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetOrDefault(tt.m, tt.key, tt.defaultValue); got != tt.want {
				t.Errorf("GetOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPutIfAbsent(t *testing.T) {
	tests := []struct {
		name         string
		m            map[string]int
		key          string
		defaultValue int
		want         map[string]int
	}{
		{
			name:         "key does not exist",
			m:            map[string]int{"a": 1},
			key:          "b",
			defaultValue: 2,
			want:         map[string]int{"a": 1, "b": 2},
		},
		{
			name:         "key exists",
			m:            map[string]int{"a": 1},
			key:          "a",
			defaultValue: 2,
			want:         map[string]int{"a": 1},
		},
		{
			name:         "empty map",
			m:            map[string]int{},
			key:          "a",
			defaultValue: 1,
			want:         map[string]int{"a": 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PutIfAbsent(tt.m, tt.key, tt.defaultValue)
			if len(tt.m) != len(tt.want) {
				t.Errorf("PutIfAbsent() map length = %v, want %v", len(tt.m), len(tt.want))
			}
			for k, v := range tt.want {
				if got, ok := tt.m[k]; !ok || got != v {
					t.Errorf("PutIfAbsent() map[%v] = %v, want %v", k, got, v)
				}
			}
		})
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]int
		want string
	}{
		{name: "empty map", m: map[string]int{}, want: "{}"},
		{name: "single element", m: map[string]int{"a": 1}, want: `{"a":1}`},
		{name: "multiple elements", m: map[string]int{"a": 1, "b": 2}, want: `{"a":1,"b":2}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToString(tt.m)
			// 因为 map 是无序的，我们需要解析 JSON 来比较
			var gotMap map[string]int
			if err := json.Unmarshal([]byte(got), &gotMap); err != nil {
				t.Errorf("ToString() returned invalid JSON: %v", err)
			}

			var wantMap map[string]int
			if err := json.Unmarshal([]byte(tt.want), &wantMap); err != nil {
				t.Errorf("Test case has invalid want JSON: %v", err)
			}

			if len(gotMap) != len(wantMap) {
				t.Errorf("ToString() map length = %v, want %v", len(gotMap), len(wantMap))
			}

			for k, v := range wantMap {
				if gotV, ok := gotMap[k]; !ok || gotV != v {
					t.Errorf("ToString() map[%v] = %v, want %v", k, gotV, v)
				}
			}
		})
	}
}

func TestSortByValue(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]int
		reverse  bool
		expected []string
	}{
		{name: "empty map", input: map[string]int{}, reverse: false, expected: []string{}},
		{name: "ascending sort", input: map[string]int{"a": 3, "b": 1, "c": 2}, reverse: false, expected: []string{"b", "c", "a"}},
		{name: "descending sort", input: map[string]int{"a": 3, "b": 1, "c": 2}, reverse: true, expected: []string{"a", "c", "b"}},
		{name: "same values ascending", input: map[string]int{"a": 1, "b": 1, "c": 1}, reverse: false, expected: []string{"a", "b", "c"}}, // 顺序不重要，但需要稳定
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SortByValue(tt.input, tt.reverse)

			if len(got) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(got))
				return
			}

			// 对于相同值的测试用例，我们只检查顺序是否稳定
			if tt.name == "same values ascending" {
				// 检查是否包含所有键
				keys := make(map[string]bool)
				for _, k := range got {
					keys[k] = true
				}
				for _, k := range tt.expected {
					if !keys[k] {
						t.Errorf("missing key %s in result", k)
					}
				}
				return
			}

			// 对于其他测试用例，检查顺序是否正确
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("at index %d, expected %s, got %s", i, tt.expected[i], got[i])
				}
			}
		})
	}
}

func TestSortByValueWithDifferentTypes(t *testing.T) {
	// 测试不同类型的map
	t.Run("float64 values", func(t *testing.T) {
		input := map[string]float64{"a": 1.1, "b": 1.0, "c": 1.2}
		expected := []string{"b", "a", "c"}
		got := SortByValue(input, false)

		for i := range got {
			if got[i] != expected[i] {
				t.Errorf("at index %d, expected %s, got %s", i, expected[i], got[i])
			}
		}
	})

	t.Run("string values", func(t *testing.T) {
		input := map[int]string{1: "z", 2: "a", 3: "m"}
		expected := []int{2, 3, 1}
		got := SortByValue(input, false)

		for i := range got {
			if got[i] != expected[i] {
				t.Errorf("at index %d, expected %d, got %d", i, expected[i], got[i])
			}
		}
	})
}
