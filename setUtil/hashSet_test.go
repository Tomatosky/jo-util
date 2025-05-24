package setUtil

import (
	"encoding/json"
	"testing"
)

func TestNewHashSet(t *testing.T) {
	// 测试空集合
	emptySet := NewHashSet[int]()
	if emptySet.Size() != 0 {
		t.Errorf("Expected empty set size 0, got %d", emptySet.Size())
	}

	// 测试带初始元素的集合
	set := NewHashSet(1, 2, 3)
	if set.Size() != 3 {
		t.Errorf("Expected set size 3, got %d", set.Size())
	}
}

func TestAdd(t *testing.T) {
	set := NewHashSet[string]()
	set.Add("a")
	if !set.Contains("a") {
		t.Error("Expected set to contain 'a'")
	}
	if set.Size() != 1 {
		t.Errorf("Expected set size 1, got %d", set.Size())
	}

	// 测试重复添加
	set.Add("a")
	if set.Size() != 1 {
		t.Errorf("Expected set size still 1 after duplicate add, got %d", set.Size())
	}
}

func TestAddAll(t *testing.T) {
	set := NewHashSet[int]()
	set.AddAll(1, 2, 3, 2) // 包含重复元素

	if set.Size() != 3 {
		t.Errorf("Expected set size 3, got %d", set.Size())
	}
	for _, v := range []int{1, 2, 3} {
		if !set.Contains(v) {
			t.Errorf("Expected set to contain %d", v)
		}
	}
}

func TestRemove(t *testing.T) {
	set := NewHashSet("a", "b", "c")
	set.Remove("b")

	if set.Contains("b") {
		t.Error("Expected set not to contain 'b' after removal")
	}
	if set.Size() != 2 {
		t.Errorf("Expected set size 2, got %d", set.Size())
	}

	// 测试移除不存在的元素
	set.Remove("d") // 不应该panic
	if set.Size() != 2 {
		t.Errorf("Expected set size still 2 after removing non-existent element, got %d", set.Size())
	}
}

func TestContains(t *testing.T) {
	set := NewHashSet(1.1, 2.2, 3.3)

	if !set.Contains(2.2) {
		t.Error("Expected set to contain 2.2")
	}
	if set.Contains(4.4) {
		t.Error("Expected set not to contain 4.4")
	}
}

func TestSize(t *testing.T) {
	set := NewHashSet[rune]()
	if set.Size() != 0 {
		t.Errorf("Expected initial size 0, got %d", set.Size())
	}

	set.Add('a')
	set.Add('b')
	if set.Size() != 2 {
		t.Errorf("Expected size 2, got %d", set.Size())
	}

	set.Remove('a')
	if set.Size() != 1 {
		t.Errorf("Expected size 1 after removal, got %d", set.Size())
	}
}

func TestClear(t *testing.T) {
	set := NewHashSet("x", "y", "z")
	set.Clear()

	if !set.IsEmpty() {
		t.Error("Expected set to be empty after clear")
	}
	if set.Size() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", set.Size())
	}
}

func TestRange(t *testing.T) {
	set := NewHashSet(1, 2, 3, 4, 5)
	count := 0

	set.Range(func(n int) bool {
		count++
		return count < 3 // 只处理前两个元素
	})

	if count != 3 { // 因为我们在处理第三个元素时返回了false
		t.Errorf("Expected to process 3 elements, processed %d", count)
	}
}

func TestToSlice(t *testing.T) {
	elements := []int{1, 2, 3, 4}
	set := NewHashSet(elements...)
	slice := set.ToSlice()

	if len(slice) != len(elements) {
		t.Errorf("Expected slice length %d, got %d", len(elements), len(slice))
	}

	// 检查所有元素都存在
	for _, v := range elements {
		found := false
		for _, s := range slice {
			if v == s {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected element %d in slice", v)
		}
	}
}

func TestIsEmpty(t *testing.T) {
	emptySet := NewHashSet[int]()
	if !emptySet.IsEmpty() {
		t.Error("Expected new set to be empty")
	}

	nonEmptySet := NewHashSet(1)
	if nonEmptySet.IsEmpty() {
		t.Error("Expected set with elements to not be empty")
	}
}

func TestToString(t *testing.T) {
	set := NewHashSet("apple", "banana", "cherry")
	str := set.ToString()

	// 验证JSON格式
	var slice []string
	err := json.Unmarshal([]byte(str), &slice)
	if err != nil {
		t.Errorf("Failed to unmarshal set string: %v", err)
	}

	if len(slice) != 3 {
		t.Errorf("Expected 3 elements in JSON, got %d", len(slice))
	}

	// 检查所有元素都存在
	for _, v := range []string{"apple", "banana", "cherry"} {
		found := false
		for _, s := range slice {
			if v == s {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected element %s in JSON", v)
		}
	}
}
