package mapUtil

import (
	"reflect"
	"testing"
)

func TestDeepGet(t *testing.T) {
	testMap := map[string]any{
		"direct":  "direct_value",
		"top.sub": "literal_dotted_key",
		"top": map[string]any{
			"sub":     "value",
			"number":  42,
			"nested":  map[string]any{"value": "nested_value"},
			"strings": map[string]string{"key": "string_map_value"},
			"anyMap":  map[any]any{"key": "any_map_value"},
			"array": []map[string]any{
				{"id": 1, "name": "first"},
				{"id": 2, "name": "second"},
			},
		},
	}

	tests := []struct {
		name string
		path string
		want any
		ok   bool
	}{
		{name: "direct key", path: "direct", want: "direct_value", ok: true},
		{name: "literal dotted key takes precedence", path: "top.sub", want: "literal_dotted_key", ok: true},
		{name: "nested any map", path: "top.nested.value", want: "nested_value", ok: true},
		{name: "nested numeric value", path: "top.number", want: 42, ok: true},
		{name: "nested string map", path: "top.strings.key", want: "string_map_value", ok: true},
		{name: "nested interface-key map", path: "top.anyMap.key", want: "any_map_value", ok: true},
		{name: "slice index", path: "top.array.1.name", want: "second", ok: true},
		{name: "slice wildcard", path: "top.array.*", want: testMap["top"].(map[string]any)["array"], ok: true},
		{name: "slice wildcard projection", path: "top.array.*.name", want: []any{"first", "second"}, ok: true},
		{name: "empty path", path: "", want: testMap, ok: true},
		{name: "missing key", path: "top.missing", want: nil, ok: false},
		{name: "invalid slice index", path: "top.array.nope", want: nil, ok: false},
		{name: "slice index out of range", path: "top.array.2", want: nil, ok: false},
		{name: "negative slice index", path: "top.array.-1", want: nil, ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := GetByPath(testMap, tt.path)
			if ok != tt.ok || !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("GetByPath(%q) = (%#v, %v), want (%#v, %v)", tt.path, got, ok, tt.want, tt.ok)
			}
			if deepGot := DeepGet(testMap, tt.path); !reflect.DeepEqual(deepGot, tt.want) {
				t.Errorf("DeepGet(%q) = %#v, want %#v", tt.path, deepGot, tt.want)
			}
		})
	}

	got, ok := GetByPathKeys(testMap, nil)
	if !ok || !reflect.DeepEqual(got, testMap) {
		t.Fatalf("GetByPathKeys with no keys = (%#v, %v), want original map", got, ok)
	}
}

func TestSetByPath(t *testing.T) {
	tests := []struct {
		name    string
		initial map[string]any
		path    string
		value   any
		want    map[string]any
	}{
		{
			name: "set direct key", initial: map[string]any{}, path: "key", value: "value",
			want: map[string]any{"key": "value"},
		},
		{
			name: "initialize nil map", initial: nil, path: "key", value: "value",
			want: map[string]any{"key": "value"},
		},
		{
			name: "set nested key", initial: map[string]any{}, path: "parent.child", value: "child_value",
			want: map[string]any{"parent": map[string]any{"child": "child_value"}},
		},
		{
			name:    "update nested key",
			initial: map[string]any{"parent": map[string]any{"child": "old_value"}},
			path:    "parent.child", value: "new_value",
			want: map[string]any{"parent": map[string]any{"child": "new_value"}},
		},
		{
			name: "set sparse array", initial: map[string]any{}, path: "array[2]", value: "third",
			want: map[string]any{"array": []string{"", "", "third"}},
		},
		{
			name: "set nested value in sparse array", initial: map[string]any{}, path: "array[2].name", value: "third",
			want: map[string]any{"array": []map[string]any{nil, nil, {"name": "third"}}},
		},
		{
			name:    "expand existing array",
			initial: map[string]any{"array": []string{"first"}}, path: "array[2]", value: "third",
			want: map[string]any{"array": []string{"first", "", "third"}},
		},
		{
			name:    "set nested array map",
			initial: map[string]any{"array": []map[string]any{{"name": "first"}, {"name": "second"}}},
			path:    "array.1.name", value: "updated",
			want: map[string]any{"array": []map[string]any{{"name": "first"}, {"name": "updated"}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetByPath(&tt.initial, tt.path, tt.value); err != nil {
				t.Fatalf("SetByPath(%q) unexpected error: %v", tt.path, err)
			}
			if !reflect.DeepEqual(tt.initial, tt.want) {
				t.Errorf("SetByPath(%q) result = %#v, want %#v", tt.path, tt.initial, tt.want)
			}
		})
	}
}

func TestDeepSet(t *testing.T) {
	m := map[string]any{}
	DeepSet(&m, "top.value", 42)
	want := map[string]any{"top": map[string]any{"value": 42}}
	if !reflect.DeepEqual(m, want) {
		t.Errorf("DeepSet result = %#v, want %#v", m, want)
	}
}

func TestSetByPathErrors(t *testing.T) {
	tests := []struct {
		name    string
		initial map[string]any
		path    string
	}{
		{name: "cannot descend into scalar", initial: map[string]any{"parent": "value"}, path: "parent.child"},
		{name: "cannot index scalar", initial: map[string]any{"array": "value"}, path: "array[0]"},
		{name: "cannot use named slice index", initial: map[string]any{"array": []string{"first"}}, path: "array.name"},
		{name: "negative slice index", initial: map[string]any{"array": []string{"first"}}, path: "array[-1]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := cloneStringAnyMap(tt.initial)
			if err := SetByPath(&tt.initial, tt.path, "updated"); err == nil {
				t.Fatalf("SetByPath(%q) expected an error", tt.path)
			}
			if !reflect.DeepEqual(tt.initial, before) {
				t.Errorf("failed SetByPath(%q) mutated map: got %#v, want %#v", tt.path, tt.initial, before)
			}
		})
	}

	if err := SetByKeys(nil, []string{"key"}, "value"); err == nil {
		t.Error("SetByKeys(nil, ...) expected an error")
	}
}

func TestMakeByKeysDoesNotMutateInput(t *testing.T) {
	keys := []string{"top", "nested", "value"}
	wantKeys := append([]string(nil), keys...)
	wantMap := map[string]any{"top": map[string]any{"nested": map[string]any{"value": 42}}}

	if got := MakeByKeys(keys, 42); !reflect.DeepEqual(got, wantMap) {
		t.Errorf("MakeByKeys() = %#v, want %#v", got, wantMap)
	}
	if !reflect.DeepEqual(keys, wantKeys) {
		t.Errorf("MakeByKeys mutated keys: got %v, want %v", keys, wantKeys)
	}
	if got := MakeByKeys(nil, 42); len(got) != 0 {
		t.Errorf("MakeByKeys(nil, ...) = %#v, want empty map", got)
	}
}

func cloneStringAnyMap(input map[string]any) map[string]any {
	result := make(map[string]any, len(input))
	for key, value := range input {
		switch typed := value.(type) {
		case []string:
			result[key] = append([]string(nil), typed...)
		default:
			result[key] = value
		}
	}
	return result
}
