package mapUtil

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestContainValue(t *testing.T) {
	tests := []struct {
		name  string
		m     map[string]int
		value int
		want  bool
	}{
		{name: "value exists", m: map[string]int{"a": 1, "b": 2}, value: 2, want: true},
		{name: "value does not exist", m: map[string]int{"a": 1}, value: 2, want: false},
		{name: "zero value exists", m: map[string]int{"zero": 0}, value: 0, want: true},
		{name: "empty map", m: map[string]int{}, value: 0, want: false},
		{name: "nil map", m: nil, value: 0, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainValue(tt.m, tt.value); got != tt.want {
				t.Fatalf("ContainValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		{name: "nil map", m: nil, want: []string{}},
		{name: "single element", m: map[string]int{"a": 1}, want: []string{"a"}},
		{name: "multiple elements", m: map[string]int{"a": 1, "b": 2, "c": 3}, want: []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Keys(tt.m)
			gotSet := make(map[string]int, len(got))
			for _, key := range got {
				gotSet[key]++
			}
			wantSet := make(map[string]int, len(tt.want))
			for _, key := range tt.want {
				wantSet[key]++
			}
			if !reflect.DeepEqual(gotSet, wantSet) {
				t.Errorf("Keys() = %v, want keys %v", got, tt.want)
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
			name: "nil map",
			m:    nil,
			want: []int{},
		},
		{
			name: "single element",
			m:    map[string]int{"a": 1},
			want: []int{1},
		},
		{
			name: "multiple elements",
			m:    map[string]int{"a": 1, "b": 1, "c": 2},
			want: []int{1, 1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Values(tt.m)
			gotCounts := make(map[int]int, len(got))
			for _, value := range got {
				gotCounts[value]++
			}
			wantCounts := make(map[int]int, len(tt.want))
			for _, value := range tt.want {
				wantCounts[value]++
			}
			if !reflect.DeepEqual(gotCounts, wantCounts) {
				t.Errorf("Values() = %v, want values %v", got, tt.want)
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
		{name: "nil map", m: nil, key: "a", defaultValue: 1, want: 1},
		{name: "existing zero value", m: map[string]int{"a": 0}, key: "a", defaultValue: 1, want: 0},
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
		{name: "nil map", m: nil, want: `null`},
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

			if !reflect.DeepEqual(gotMap, wantMap) {
				t.Errorf("ToString() parsed to %#v, want %#v", gotMap, wantMap)
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
		{name: "equal values", input: map[string]int{"a": 1, "b": 1, "c": 1}, reverse: false, expected: []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SortByValue(tt.input, tt.reverse)

			if len(got) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(got))
				return
			}

			seen := make(map[string]bool, len(got))
			for i, key := range got {
				if _, exists := tt.input[key]; !exists || seen[key] {
					t.Fatalf("result contains invalid or duplicate key %q: %v", key, got)
				}
				seen[key] = true
				if i > 0 {
					previous, current := tt.input[got[i-1]], tt.input[key]
					if (!tt.reverse && previous > current) || (tt.reverse && previous < current) {
						t.Errorf("values are not sorted at index %d: %d then %d", i, previous, current)
					}
				}
			}
			if tt.name != "equal values" && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("SortByValue() = %v, want %v", got, tt.expected)
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
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("SortByValue() = %v, want %v", got, expected)
		}
	})

	t.Run("string values", func(t *testing.T) {
		input := map[int]string{1: "z", 2: "a", 3: "m"}
		expected := []int{2, 3, 1}
		got := SortByValue(input, false)
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("SortByValue() = %v, want %v", got, expected)
		}
	})
}

func TestMapImplementationsMarshalJSONErrors(t *testing.T) {
	value := make(chan int)
	ordered := NewOrderedMap[string, chan int]()
	ordered.Put("key", value)
	hash := NewConcurrentHashMap[string, chan int]()
	hash.Put("key", value)
	hash2 := NewConcurrentHashMap2[string, chan int]()
	hash2.Put("key", value)
	skipList := NewConcurrentSkipListMap[string, chan int]()
	skipList.Put("key", value)
	tree := NewTreeMap[string, chan int](func(a, b string) bool { return a < b })
	tree.Put("key", value)
	biMap := NewBiMap[string, chan int]()
	biMap.Put("key", value)

	tests := []struct {
		name      string
		marshaler interface{ MarshalJSON() ([]byte, error) }
	}{
		{name: "ordered map", marshaler: ordered},
		{name: "concurrent hash map", marshaler: hash},
		{name: "sync map", marshaler: hash2},
		{name: "skip list map", marshaler: skipList},
		{name: "tree map", marshaler: tree},
		{name: "bi-map", marshaler: biMap},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.marshaler.MarshalJSON(); err == nil {
				t.Error("MarshalJSON with an unsupported channel value expected an error")
			}
		})
	}
}
