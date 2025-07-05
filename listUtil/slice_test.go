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

func TestIntersectionInt(t *testing.T) {
	tests := []struct {
		name     string
		a        []int
		b        []int
		expected []int
	}{
		{
			name:     "两个空slice",
			a:        []int{},
			b:        []int{},
			expected: []int{},
		},
		{
			name:     "第一个slice为空",
			a:        []int{},
			b:        []int{1, 2, 3},
			expected: []int{},
		},
		{
			name:     "第二个slice为空",
			a:        []int{1, 2, 3},
			b:        []int{},
			expected: []int{},
		},
		{
			name:     "有共同元素",
			a:        []int{1, 2, 3, 4},
			b:        []int{3, 4, 5, 6},
			expected: []int{3, 4},
		},
		{
			name:     "无共同元素",
			a:        []int{1, 2, 3},
			b:        []int{4, 5, 6},
			expected: []int{},
		},
		{
			name:     "完全相同的slice",
			a:        []int{1, 2, 3},
			b:        []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		{
			name:     "重复元素",
			a:        []int{1, 2, 2, 3},
			b:        []int{2, 2, 3, 4},
			expected: []int{2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Intersection(tt.a, tt.b)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Intersection() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUnion(t *testing.T) {
	tests := []struct {
		name     string
		a        []int
		b        []int
		expected []int
	}{
		{
			name:     "两个空slice",
			a:        []int{},
			b:        []int{},
			expected: []int{},
		},
		{
			name:     "第一个slice为空",
			a:        []int{},
			b:        []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		{
			name:     "第二个slice为空",
			a:        []int{1, 2, 3},
			b:        []int{},
			expected: []int{1, 2, 3},
		},
		{
			name:     "两个slice完全相同",
			a:        []int{1, 2, 3},
			b:        []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		{
			name:     "两个slice完全不同",
			a:        []int{1, 2, 3},
			b:        []int{4, 5, 6},
			expected: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:     "部分重复元素",
			a:        []int{1, 2, 3, 4},
			b:        []int{3, 4, 5, 6},
			expected: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:     "包含重复元素的slice",
			a:        []int{1, 2, 2, 3},
			b:        []int{3, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Union(tt.a, tt.b)
			if len(result) != len(tt.expected) {
				t.Errorf("期望长度 %d, 实际长度 %d", len(tt.expected), len(result))
				return
			}
			// 检查结果中是否包含所有期望的元素
			expectedSet := make(map[int]bool)
			for _, item := range tt.expected {
				expectedSet[item] = true
			}
			for _, item := range result {
				if !expectedSet[item] {
					t.Errorf("结果中包含不期望的元素 %d", item)
				}
			}
			// 检查所有期望的元素是否都在结果中
			resultSet := make(map[int]bool)
			for _, item := range result {
				resultSet[item] = true
			}
			for _, item := range tt.expected {
				if !resultSet[item] {
					t.Errorf("结果中缺少期望的元素 %d", item)
				}
			}
		})
	}
}

// 测试结构体
type Person struct {
	Name string
	Age  int
}

func TestFieldExtractor(t *testing.T) {
	t.Run("提取Person的Name字段", func(t *testing.T) {
		people := []Person{
			{Name: "Alice", Age: 25},
			{Name: "Bob", Age: 30},
		}
		got := FieldExtractor(people, func(p Person) string { return p.Name })
		want := []string{"Alice", "Bob"}
		if len(got) != len(want) {
			t.Errorf("长度不匹配，期望 %d，得到 %d", len(want), len(got))
		}
		for i := range got {
			if got[i] != want[i] {
				t.Errorf("索引 %d 不匹配，期望 %v，得到 %v", i, want[i], got[i])
			}
		}
	})

	t.Run("提取Person的Age字段", func(t *testing.T) {
		people := []Person{
			{Name: "Alice", Age: 25},
			{Name: "Bob", Age: 30},
		}
		got := FieldExtractor(people, func(p Person) int { return p.Age })
		want := []int{25, 30}
		if len(got) != len(want) {
			t.Errorf("长度不匹配，期望 %d，得到 %d", len(want), len(got))
		}
		for i := range got {
			if got[i] != want[i] {
				t.Errorf("索引 %d 不匹配，期望 %v，得到 %v", i, want[i], got[i])
			}
		}
	})

	t.Run("空切片", func(t *testing.T) {
		var people []Person
		got := FieldExtractor(people, func(p Person) string { return p.Name })
		if len(got) != 0 {
			t.Errorf("期望空切片，得到长度 %d", len(got))
		}
	})
}

func TestToMap(t *testing.T) {
	// 测试用例1：正常情况，多个Person结构体
	t.Run("Normal case with multiple persons", func(t *testing.T) {
		people := []Person{
			{Name: "Alice", Age: 25},
			{Name: "Bob", Age: 30},
			{Name: "Charlie", Age: 35},
		}
		// 使用Name作为key
		result := ToMap(people, func(p Person) string { return p.Name })
		if len(result) != 3 {
			t.Errorf("Expected map length 3, got %d", len(result))
		}
		if result["Alice"].Age != 25 {
			t.Errorf("Expected Alice's age 25, got %d", result["Alice"].Age)
		}
		if result["Bob"].Name != "Bob" {
			t.Errorf("Expected Bob's name Bob, got %s", result["Bob"].Name)
		}
	})
	// 测试用例2：空切片
	t.Run("Empty slice", func(t *testing.T) {
		var empty []Person
		result := ToMap(empty, func(p Person) string { return p.Name })
		if len(result) != 0 {
			t.Errorf("Expected empty map, got %d", len(result))
		}
	})
	// 测试用例3：重复Name的情况（后出现的会覆盖前面的）
	t.Run("Duplicate names", func(t *testing.T) {
		people := []Person{
			{Name: "Alice", Age: 25},
			{Name: "Alice", Age: 30}, // 同名，年龄不同
			{Name: "Bob", Age: 35},
		}
		result := ToMap(people, func(p Person) string { return p.Name })
		if len(result) != 2 {
			t.Errorf("Expected 2 unique names, got %d", len(result))
		}
		if result["Alice"].Age != 30 {
			t.Errorf("Expected last Alice's age 30, got %d", result["Alice"].Age)
		}
	})
	// 测试用例4：使用Age作为key（测试不同类型的key）
	t.Run("Using age as key", func(t *testing.T) {
		people := []Person{
			{Name: "Alice", Age: 25},
			{Name: "Bob", Age: 30},
			{Name: "Charlie", Age: 25}, // 相同年龄
		}
		result := ToMap(people, func(p Person) int { return p.Age })
		if len(result) != 2 {
			t.Errorf("Expected 2 unique ages, got %d", len(result))
		}
		if result[25].Name != "Charlie" {
			t.Errorf("Expected last person with age 25 to be Charlie, got %s", result[25].Name)
		}
	})
}
