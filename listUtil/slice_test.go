package listUtil

import (
	"reflect"
	"testing"
)

func TestContain(t *testing.T) {
	tests := []struct {
		name   string
		slice  []int
		target int
		want   bool
	}{
		{"Contains", []int{1, 2, 3}, 2, true},
		{"NotContains", []int{1, 2, 3}, 4, false},
		{"EmptySlice", []int{}, 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contain(tt.slice, tt.target); got != tt.want {
				t.Errorf("Contain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	tests := []struct {
		name  string
		slice []int
		want  []int
	}{
		{"NoDuplicates", []int{1, 2, 3}, []int{1, 2, 3}},
		{"WithDuplicates", []int{1, 2, 2, 3, 3, 3}, []int{1, 2, 3}},
		{"EmptySlice", []int{}, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Unique(tt.slice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		name  string
		slice []int
		want  string
	}{
		{"NormalCase", []int{1, 2, 3}, "[1,2,3]"},
		{"EmptySlice", []int{}, "[]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToString(tt.slice); got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		name  string
		slice []int
		want  []int
	}{
		{"EvenLength", []int{1, 2, 3, 4}, []int{4, 3, 2, 1}},
		{"OddLength", []int{1, 2, 3}, []int{3, 2, 1}},
		{"EmptySlice", []int{}, []int{}},
		{"SingleElement", []int{1}, []int{1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 复制切片以避免修改原始测试数据
			slice := make([]int, len(tt.slice))
			copy(slice, tt.slice)

			Reverse(slice)
			if !reflect.DeepEqual(slice, tt.want) {
				t.Errorf("Reverse() = %v, want %v", slice, tt.want)
			}
		})
	}
}

func TestShuffle(t *testing.T) {
	original := []int{1, 2, 3, 4, 5}
	shuffled := Shuffle(original)

	// 检查长度是否相同
	if len(shuffled) != len(original) {
		t.Errorf("Shuffle() changed slice length, got %v, want %v", len(shuffled), len(original))
	}

	// 检查元素是否相同(顺序可能不同)
	originalMap := make(map[int]int)
	shuffledMap := make(map[int]int)

	for _, v := range original {
		originalMap[v]++
	}

	for _, v := range shuffled {
		shuffledMap[v]++
	}

	if !reflect.DeepEqual(originalMap, shuffledMap) {
		t.Errorf("Shuffle() changed elements, got %v, want same elements as %v", shuffled, original)
	}
}

func TestAddIfAbsent(t *testing.T) {
	tests := []struct {
		name  string
		slice []int
		item  int
		want  []int
	}{
		{"ItemPresent", []int{1, 2, 3}, 2, []int{1, 2, 3}},
		{"ItemNotPresent", []int{1, 2}, 3, []int{1, 2, 3}},
		{"EmptySlice", []int{}, 1, []int{1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slice := make([]int, len(tt.slice))
			copy(slice, tt.slice)

			AddIfAbsent(&slice, tt.item)
			if !reflect.DeepEqual(slice, tt.want) {
				t.Errorf("AddIfAbsent() = %v, want %v", slice, tt.want)
			}
		})
	}
}

func TestRemove2(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		target   int
		all      bool
		expected []int
	}{
		{
			name:     "Remove first occurrence",
			input:    []int{1, 2, 3, 2, 4},
			target:   2,
			all:      false,
			expected: []int{1, 3, 2, 4},
		},
		{
			name:     "Remove all occurrences",
			input:    []int{1, 2, 3, 2, 4},
			target:   2,
			all:      true,
			expected: []int{1, 3, 4},
		},
		{
			name:     "Target not found",
			input:    []int{1, 2, 3},
			target:   4,
			all:      false,
			expected: []int{1, 2, 3},
		},
		{
			name:     "Empty slice",
			input:    []int{},
			target:   1,
			all:      false,
			expected: []int{},
		},
		{
			name:     "Remove all from empty slice",
			input:    []int{},
			target:   1,
			all:      true,
			expected: []int{},
		},
		{
			name:     "Remove first from single element",
			input:    []int{5},
			target:   5,
			all:      false,
			expected: []int{},
		},
		{
			name:     "Remove all from single element",
			input:    []int{5},
			target:   5,
			all:      true,
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Remove(tt.input, tt.target, tt.all)
			if !equal(result, tt.expected) {
				t.Errorf("Remove(%v, %d, %t) = %v, want %v", tt.input, tt.target, tt.all, result, tt.expected)
			}
		})
	}
}

// equal 辅助函数用于比较两个切片是否相等
func equal[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestGetByIndex(t *testing.T) {
	tests := []struct {
		name      string
		slice     []int
		index     int
		want      int
		wantError bool
	}{
		{"empty slice", []int{}, 0, 0, true},
		{"first element", []int{1, 2, 3}, 0, 1, false},
		{"last element", []int{1, 2, 3}, 2, 3, false},
		{"negative index -1", []int{1, 2, 3}, -1, 3, false},
		{"negative index -2", []int{1, 2, 3}, -2, 2, false},
		{"index out of range", []int{1, 2, 3}, 3, 0, true},
		{"negative index out of range", []int{1, 2, 3}, -4, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetByIndex(tt.slice, tt.index)
			if (err != nil) != tt.wantError {
				t.Errorf("GetByIndex() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && got != tt.want {
				t.Errorf("GetByIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsertByIndex(t *testing.T) {
	tests := []struct {
		name      string
		slice     []int
		index     int
		value     int
		want      []int
		wantError bool
	}{
		{"insert at beginning", []int{2, 3}, 0, 1, []int{1, 2, 3}, false},
		{"insert at end", []int{1, 2}, 2, 3, []int{1, 2, 3}, false},
		{"insert in middle", []int{1, 3}, 1, 2, []int{1, 2, 3}, false},
		{"negative index -1", []int{1, 3}, -1, 2, []int{1, 3, 2}, false},
		{"negative index -2", []int{1, 3}, -2, 2, []int{1, 2, 3}, false},
		{"index out of range", []int{1, 2}, 3, 3, nil, true},
		{"negative index out of range", []int{1, 2}, -4, 0, nil, true},
		{"empty slice insert at 0", []int{}, 0, 1, []int{1}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InsertByIndex(tt.slice, tt.index, tt.value)
			if (err != nil) != tt.wantError {
				t.Errorf("InsertByIndex() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError {
				if len(got) != len(tt.want) {
					t.Errorf("InsertByIndex() length = %v, want %v", len(got), len(tt.want))
					return
				}
				for i := range got {
					if got[i] != tt.want[i] {
						t.Errorf("InsertByIndex() = %v, want %v", got, tt.want)
						return
					}
				}
			}
		})
	}
}
