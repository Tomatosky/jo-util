package mapUtil

import (
	"encoding/json"
	"testing"
)

func TestKeys(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]int
		want []string
	}{
		{
			name: "empty map",
			m:    map[string]int{},
			want: []string{},
		},
		{
			name: "single element",
			m:    map[string]int{"a": 1},
			want: []string{"a"},
		},
		{
			name: "multiple elements",
			m:    map[string]int{"a": 1, "b": 2, "c": 3},
			want: []string{"a", "b", "c"},
		},
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
		{
			name:         "key exists",
			m:            map[string]int{"a": 1, "b": 2},
			key:          "a",
			defaultValue: 0,
			want:         1,
		},
		{
			name:         "key does not exist",
			m:            map[string]int{"a": 1, "b": 2},
			key:          "c",
			defaultValue: 3,
			want:         3,
		},
		{
			name:         "empty map",
			m:            map[string]int{},
			key:          "a",
			defaultValue: 1,
			want:         1,
		},
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
		{
			name: "empty map",
			m:    map[string]int{},
			want: "{}",
		},
		{
			name: "single element",
			m:    map[string]int{"a": 1},
			want: `{"a":1}`,
		},
		{
			name: "multiple elements",
			m:    map[string]int{"a": 1, "b": 2},
			want: `{"a":1,"b":2}`,
		},
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
