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

// 测试结构体
type testUser struct {
	ID   int
	Name string
	Age  int
}

// 测试指针结构体
type testProduct struct {
	Code  string
	Price float64
}

func TestToMap(t *testing.T) {
	// 测试用例1：正常情况，使用int类型作为key
	t.Run("Normal case with int key", func(t *testing.T) {
		users := []testUser{
			{ID: 1, Name: "Alice", Age: 25},
			{ID: 2, Name: "Bob", Age: 30},
		}

		result, err := ToMap[int, testUser]("ID", users)
		if err != nil {
			t.Error("Unexpected error:", err)
		}

		if len(result) != 2 {
			t.Error("Expected map length 2, got", len(result))
		}

		if result[1].Name != "Alice" || result[2].Name != "Bob" {
			t.Error("Map content incorrect")
		}
	})

	// 测试用例2：正常情况，使用string类型作为key
	t.Run("Normal case with string key", func(t *testing.T) {
		users := []testUser{
			{ID: 1, Name: "Alice", Age: 25},
			{ID: 2, Name: "Bob", Age: 30},
		}

		result, err := ToMap[string, testUser]("Name", users)
		if err != nil {
			t.Error("Unexpected error:", err)
		}

		if result["Alice"].ID != 1 || result["Bob"].Age != 30 {
			t.Error("Map content incorrect")
		}
	})

	// 测试用例3：测试指针结构体
	t.Run("Pointer struct case", func(t *testing.T) {
		products := []*testProduct{
			{Code: "P001", Price: 9.99},
			{Code: "P002", Price: 19.99},
		}

		result, err := ToMap[string, *testProduct]("Code", products)
		if err != nil {
			t.Error("Unexpected error:", err)
		}

		if result["P001"].Price != 9.99 || result["P002"].Price != 19.99 {
			t.Error("Map content incorrect")
		}
	})

	// 测试用例4：字段不存在的情况
	t.Run("Field not found case", func(t *testing.T) {
		users := []testUser{
			{ID: 1, Name: "Alice", Age: 25},
		}

		_, err := ToMap[int, testUser]("NonExistentField", users)
		if err == nil {
			t.Error("Expected error for non-existent field, got nil")
		}
	})

	// 测试用例5：字段类型不匹配的情况
	t.Run("Field type mismatch case", func(t *testing.T) {
		users := []testUser{
			{ID: 1, Name: "Alice", Age: 25},
		}

		_, err := ToMap[string, testUser]("ID", users)
		if err == nil {
			t.Error("Expected error for type mismatch, got nil")
		}
	})

	// 测试用例6：非结构体类型的情况
	t.Run("Non-struct type case", func(t *testing.T) {
		numbers := []int{1, 2, 3}

		_, err := ToMap[int, int]("Field", numbers)
		if err == nil {
			t.Error("Expected error for non-struct type, got nil")
		}
	})

	// 测试用例7：空切片的情况
	t.Run("Empty slice case", func(t *testing.T) {
		var users []testUser

		result, err := ToMap[int, testUser]("ID", users)
		if err != nil {
			t.Error("Unexpected error:", err)
		}

		if len(result) != 0 {
			t.Error("Expected empty map, got", len(result))
		}
	})

	// 测试用例8：重复key的情况（应该覆盖）
	t.Run("Duplicate key case", func(t *testing.T) {
		users := []testUser{
			{ID: 1, Name: "Alice", Age: 25},
			{ID: 1, Name: "Bob", Age: 30}, // 相同ID
		}

		result, err := ToMap[int, testUser]("ID", users)
		if err != nil {
			t.Error("Unexpected error:", err)
		}

		if len(result) != 1 {
			t.Error("Expected map length 1 due to duplicate keys, got", len(result))
		}

		// 应该保留最后一个值
		if result[1].Name != "Bob" {
			t.Error("Expected last value to be kept for duplicate key")
		}
	})
}

func TestEvery(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		elements []int
		expected bool
	}{
		{
			name:     "空集合和空元素",
			input:    []int{},
			elements: []int{},
			expected: true,
		},
		{
			name:     "空集合但非空元素",
			input:    []int{},
			elements: []int{1, 2},
			expected: false,
		},
		{
			name:     "非空集合和空元素",
			input:    []int{1, 2, 3},
			elements: []int{},
			expected: true,
		},
		{
			name:     "包含所有元素",
			input:    []int{1, 2, 3, 4, 5},
			elements: []int{2, 4},
			expected: true,
		},
		{
			name:     "不包含所有元素",
			input:    []int{1, 2, 3, 4, 5},
			elements: []int{2, 6},
			expected: false,
		},
		{
			name:     "重复元素在集合中",
			input:    []int{1, 2, 2, 3, 3, 3},
			elements: []int{2, 3},
			expected: true,
		},
		{
			name:     "重复元素在查找列表中",
			input:    []int{1, 2, 3},
			elements: []int{2, 2, 3},
			expected: true,
		},
		{
			name:     "完全匹配",
			input:    []int{1, 2, 3},
			elements: []int{1, 2, 3},
			expected: true,
		},
		{
			name:     "部分匹配",
			input:    []int{1, 2, 3},
			elements: []int{1, 2, 4},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainAll(tt.input, tt.elements...)
			if result != tt.expected {
				t.Errorf("Every(%v, %v) = %v, want %v", tt.input, tt.elements, result, tt.expected)
			}
		})
	}
}

// 测试字符串类型的Every函数
func TestEveryString(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		elements []string
		expected bool
	}{
		{
			name:     "字符串匹配",
			input:    []string{"a", "b", "c", "d"},
			elements: []string{"b", "d"},
			expected: true,
		},
		{
			name:     "字符串不匹配",
			input:    []string{"a", "b", "c", "d"},
			elements: []string{"b", "e"},
			expected: false,
		},
		{
			name:     "空字符串",
			input:    []string{"", "b", "c"},
			elements: []string{""},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainAll(tt.input, tt.elements...)
			if result != tt.expected {
				t.Errorf("Every(%v, %v) = %v, want %v", tt.input, tt.elements, result, tt.expected)
			}
		})
	}
}

func TestContainOne(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		elements []int
		want     bool
	}{
		{
			name:     "空切片和空元素",
			input:    []int{},
			elements: []int{},
			want:     false,
		},
		{
			name:     "空切片但有元素",
			input:    []int{},
			elements: []int{1, 2, 3},
			want:     false,
		},
		{
			name:     "非空切片但空元素",
			input:    []int{1, 2, 3},
			elements: []int{},
			want:     false,
		},
		{
			name:     "包含单个匹配元素",
			input:    []int{1, 2, 3},
			elements: []int{2},
			want:     true,
		},
		{
			name:     "包含多个匹配元素中的一个",
			input:    []int{1, 2, 3},
			elements: []int{4, 2, 5},
			want:     true,
		},
		{
			name:     "不包含任何元素",
			input:    []int{1, 2, 3},
			elements: []int{4, 5, 6},
			want:     false,
		},
		{
			name:     "重复元素且匹配",
			input:    []int{1, 1, 2, 2, 3},
			elements: []int{2},
			want:     true,
		},
		{
			name:     "重复元素但不匹配",
			input:    []int{1, 1, 2, 2, 3},
			elements: []int{4},
			want:     false,
		},
		{
			name:     "大切片测试",
			input:    make([]int, 1000), // 1000个0
			elements: []int{0, 999},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainOne(tt.input, tt.elements...); got != tt.want {
				t.Errorf("ContainOne() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 测试字符串类型的ContainOne
func TestContainOneString(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		elements []string
		want     bool
	}{
		{
			name:     "字符串匹配",
			input:    []string{"apple", "banana", "orange"},
			elements: []string{"banana", "grape"},
			want:     true,
		},
		{
			name:     "字符串不匹配",
			input:    []string{"apple", "banana", "orange"},
			elements: []string{"pear", "grape"},
			want:     false,
		},
		{
			name:     "空字符串测试",
			input:    []string{"", "banana", "orange"},
			elements: []string{""},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainOne(tt.input, tt.elements...); got != tt.want {
				t.Errorf("ContainOne() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 测试自定义结构体类型的ContainOne
func TestContainOneStruct(t *testing.T) {
	type person struct {
		name string
		age  int
	}

	tests := []struct {
		name     string
		input    []person
		elements []person
		want     bool
	}{
		{
			name: "结构体匹配",
			input: []person{
				{"Alice", 20},
				{"Bob", 30},
			},
			elements: []person{
				{"Bob", 30},
				{"Charlie", 40},
			},
			want: true,
		},
		{
			name: "结构体不匹配",
			input: []person{
				{"Alice", 20},
				{"Bob", 30},
			},
			elements: []person{
				{"Alice", 30}, // 同名但不同年龄
				{"Bob", 20},   // 同年龄但不同名
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainOne(tt.input, tt.elements...); got != tt.want {
				t.Errorf("ContainOne() = %v, want %v", got, tt.want)
			}
		})
	}
}
