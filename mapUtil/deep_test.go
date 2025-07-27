package mapUtil

import (
	"reflect"
	"testing"
)

func TestDeepGet(t *testing.T) {
	testMap := map[string]any{
		"top": map[string]any{
			"sub":    "value",
			"number": 42,
			"nested": map[string]any{
				"key": "nested_value",
			},
			"array": []map[string]any{
				{"id": 1, "name": "first"},
				{"id": 2, "name": "second"},
			},
		},
		"direct": "direct_value",
	}

	tests := []struct {
		name     string
		path     string
		expected any
	}{
		{
			name:     "direct key",
			path:     "direct",
			expected: "direct_value",
		},
		{
			name:     "nested key",
			path:     "top.sub",
			expected: "value",
		},
		{
			name:     "deep nested key",
			path:     "top.nested.key",
			expected: "nested_value",
		},
		{
			name:     "number value",
			path:     "top.number",
			expected: 42,
		},
		{
			name:     "non-existent key",
			path:     "top.not.exists",
			expected: nil,
		},
		{
			name:     "empty path",
			path:     "",
			expected: testMap,
		},
		{
			name: "wildcard on array",
			path: "top.array.*",
			expected: []map[string]any{
				{"id": 1, "name": "first"},
				{"id": 2, "name": "second"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DeepGet(testMap, tt.path)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("DeepGet() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDeepSet(t *testing.T) {
	tests := []struct {
		name     string
		initial  map[string]any
		path     string
		value    any
		expected map[string]any
	}{
		{
			name:    "set direct key",
			initial: map[string]any{},
			path:    "key",
			value:   "value",
			expected: map[string]any{
				"key": "value",
			},
		},
		{
			name:    "set nested key",
			initial: map[string]any{},
			path:    "parent.child",
			value:   "child_value",
			expected: map[string]any{
				"parent": map[string]any{
					"child": "child_value",
				},
			},
		},
		{
			name:    "set deep nested key",
			initial: map[string]any{},
			path:    "level1.level2.level3",
			value:   "deep_value",
			expected: map[string]any{
				"level1": map[string]any{
					"level2": map[string]any{
						"level3": "deep_value",
					},
				},
			},
		},
		{
			name: "update existing key",
			initial: map[string]any{
				"existing": "old_value",
			},
			path:  "existing",
			value: "new_value",
			expected: map[string]any{
				"existing": "new_value",
			},
		},
		{
			name: "update nested existing key",
			initial: map[string]any{
				"parent": map[string]any{
					"child": "old_value",
				},
			},
			path:  "parent.child",
			value: "new_value",
			expected: map[string]any{
				"parent": map[string]any{
					"child": "new_value",
				},
			},
		},
		{
			name:    "set array index",
			initial: map[string]any{},
			path:    "array[0]",
			value:   "first",
			expected: map[string]any{
				"array": []string{"first"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DeepSet(&tt.initial, tt.path, tt.value)
			if !reflect.DeepEqual(tt.initial, tt.expected) {
				t.Errorf("DeepSet() result = %v, want %v", tt.initial, tt.expected)
			}
		})
	}
}
