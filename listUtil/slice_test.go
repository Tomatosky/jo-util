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
