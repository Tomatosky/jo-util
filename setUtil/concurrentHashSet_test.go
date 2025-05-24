package setUtil

import (
	"testing"
)

func TestNewConcurrentHashSet(t *testing.T) {
	// 测试空集合创建
	emptySet := NewConcurrentHashSet[int]()
	if emptySet.Size() != 0 {
		t.Errorf("Expected empty set size 0, got %d", emptySet.Size())
	}

	// 测试带初始元素的集合创建
	set := NewConcurrentHashSet(1, 2, 3)
	if set.Size() != 3 {
		t.Errorf("Expected set size 3, got %d", set.Size())
	}
}

func TestAddAndContains(t *testing.T) {
	set := NewConcurrentHashSet[string]()
	set.Add("apple")

	// 测试添加后包含
	if !set.Contains("apple") {
		t.Error("Expected set to contain 'apple'")
	}

	// 测试不存在的元素
	if set.Contains("banana") {
		t.Error("Set should not contain 'banana'")
	}
}

func TestAddAll3(t *testing.T) {
	set := NewConcurrentHashSet[int]()
	set.AddAll(1, 2, 3, 4, 5)

	// 测试批量添加后的数量
	if set.Size() != 5 {
		t.Errorf("Expected set size 5, got %d", set.Size())
	}

	// 测试所有元素都存在
	for i := 1; i <= 5; i++ {
		if !set.Contains(i) {
			t.Errorf("Expected set to contain %d", i)
		}
	}
}

func TestRemove3(t *testing.T) {
	set := NewConcurrentHashSet("a", "b", "c")
	set.Remove("b")

	// 测试移除后大小
	if set.Size() != 2 {
		t.Errorf("Expected set size 2 after removal, got %d", set.Size())
	}

	// 测试移除的元素不存在
	if set.Contains("b") {
		t.Error("Set should not contain 'b' after removal")
	}

	// 测试移除不存在的元素
	set.Remove("d") // 不应该panic
	if set.Size() != 2 {
		t.Error("Removing non-existent element should not change set size")
	}
}

func TestClear3(t *testing.T) {
	set := NewConcurrentHashSet(1.1, 2.2, 3.3)
	set.Clear()

	// 测试清空后大小
	if set.Size() != 0 {
		t.Errorf("Expected empty set after clear, got size %d", set.Size())
	}

	// 测试清空后不包含任何元素
	if set.Contains(1.1) {
		t.Error("Set should be empty after clear")
	}
}

func TestToSlice3(t *testing.T) {
	elements := []int{1, 2, 3, 4, 5}
	set := NewConcurrentHashSet(elements...)
	slice := set.ToSlice()

	// 测试切片长度
	if len(slice) != len(elements) {
		t.Errorf("Expected slice length %d, got %d", len(elements), len(slice))
	}

	// 测试所有元素都存在
	elementMap := make(map[int]bool)
	for _, e := range elements {
		elementMap[e] = true
	}

	for _, e := range slice {
		if !elementMap[e] {
			t.Errorf("Unexpected element %d in slice", e)
		}
	}
}

func TestIsEmpty3(t *testing.T) {
	// 测试空集合
	emptySet := NewConcurrentHashSet[string]()
	if !emptySet.IsEmpty() {
		t.Error("New set should be empty")
	}

	// 测试非空集合
	nonEmptySet := NewConcurrentHashSet("a")
	if nonEmptySet.IsEmpty() {
		t.Error("Set with elements should not be empty")
	}

	// 测试清空后的集合
	nonEmptySet.Clear()
	if !nonEmptySet.IsEmpty() {
		t.Error("Cleared set should be empty")
	}
}

func TestToString3(t *testing.T) {
	set := NewConcurrentHashSet(1, 2, 3)
	str := set.ToString()

	// 简单测试字符串格式
	if len(str) < 5 { // 至少包含 "[1,2,3]"
		t.Errorf("Unexpected string representation: %s", str)
	}

	// 测试空集合的字符串表示
	emptySet := NewConcurrentHashSet[int]()
	if emptySet.ToString() != "[]" {
		t.Errorf("Empty set should stringify to [], got %s", emptySet.ToString())
	}
}

func TestConcurrentOperations(t *testing.T) {
	set := NewConcurrentHashSet[int]()
	const numOperations = 1000
	done := make(chan bool, 2) // 缓冲通道避免goroutine泄漏

	// 并发添加偶数
	go func() {
		for i := 0; i < numOperations; i += 2 {
			set.Add(i)
		}
		done <- true
	}()

	// 并发添加奇数并移除偶数
	go func() {
		for i := 1; i < numOperations; i += 2 {
			set.Add(i)
			set.Remove(i - 1) // 尝试移除前一个偶数
		}
		done <- true
	}()

	// 等待两个goroutine完成
	<-done
	<-done

	// 验证结果
	for i := 0; i < numOperations; i++ {
		contains := set.Contains(i)
		// 偶数应该被移除（除非在奇数goroutine执行前添加goroutine已经添加了它）
		if i%2 == 0 {

		} else {
			// 奇数应该存在
			if !contains {
				t.Errorf("Odd number %d should be present", i)
			}
		}
	}

	// 更可靠的验证方式：检查至少所有奇数都存在
	for i := 1; i < numOperations; i += 2 {
		if !set.Contains(i) {
			t.Errorf("Odd number %d is missing", i)
		}
	}

	// 检查集合大小在合理范围内
	size := set.Size()
	if size < numOperations/2 || size > numOperations {
		t.Errorf("Unexpected set size %d, expected between %d and %d",
			size, numOperations/2, numOperations)
	}
}
