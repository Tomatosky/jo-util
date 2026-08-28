package setUtil

import (
	"encoding/json"
	"sync"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestNewConcurrentHashSet(t *testing.T) {
	// 测试空集合创建
	emptySet := NewConcurrentHashSet[int]()
	if emptySet.Size() != 0 {
		t.Errorf("Expected empty set size 0, got %d", emptySet.Size())
	}

	// 测试带初始元素的集合创建
	set := NewConcurrentHashSet(1, 2, 3, 2)
	if set.Size() != 3 {
		t.Errorf("Expected set size 3, got %d", set.Size())
	}
	assertIntSetElements(t, set.ToSlice(), []int{1, 2, 3})
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
	set.AddAll(1, 2, 3, 4, 5, 3)

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

	assertIntSetElements(t, slice, elements)
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

	var elements []int
	if err := json.Unmarshal([]byte(str), &elements); err != nil {
		t.Fatalf("ToString returned invalid JSON: %v", err)
	}
	assertIntSetElements(t, elements, []int{1, 2, 3})

	// 测试空集合的字符串表示
	emptySet := NewConcurrentHashSet[int]()
	if emptySet.ToString() != "[]" {
		t.Errorf("Empty set should stringify to [], got %s", emptySet.ToString())
	}
}

func TestConcurrentHashSetRange(t *testing.T) {
	set := NewConcurrentHashSet(1, 2, 3, 4, 5)
	for name, rangeFn := range map[string]func(func(int) bool){
		"Range":     set.Range,
		"CopyRange": set.CopyRange,
	} {
		t.Run(name+" visits every element", func(t *testing.T) {
			visited := make([]int, 0, set.Size())
			rangeFn(func(value int) bool {
				visited = append(visited, value)
				return true
			})
			assertIntSetElements(t, visited, []int{1, 2, 3, 4, 5})
		})

		t.Run(name+" stops immediately", func(t *testing.T) {
			count := 0
			rangeFn(func(int) bool {
				count++
				return false
			})
			if count != 1 {
				t.Fatalf("%s callbacks = %d, want 1", name, count)
			}
		})
	}
}

func TestConcurrentHashSetJSON(t *testing.T) {
	set := NewConcurrentHashSet(1, 2, 2, 3)
	data, err := json.Marshal(set)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	var encoded []int
	if err := json.Unmarshal(data, &encoded); err != nil {
		t.Fatalf("marshaled data is not a JSON array: %v", err)
	}
	assertIntSetElements(t, encoded, []int{1, 2, 3})

	restored := NewConcurrentHashSet(99)
	originalMap := restored.m
	if err := json.Unmarshal([]byte(`[3,2,2,1]`), restored); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	assertIntSetElements(t, restored.ToSlice(), []int{1, 2, 3})
	if restored.m != originalMap {
		t.Error("UnmarshalJSON replaced the synchronized map instance")
	}

	before := restored.ToSlice()
	if err := json.Unmarshal([]byte(`{"invalid":true}`), restored); err == nil {
		t.Fatal("UnmarshalJSON should reject a non-array value")
	}
	assertIntSetElements(t, restored.ToSlice(), before)

	if err := json.Unmarshal([]byte(`[]`), restored); err != nil {
		t.Fatalf("UnmarshalJSON empty array failed: %v", err)
	}
	assertIntSetElements(t, restored.ToSlice(), []int{})

	var zeroValue ConcurrentHashSet[int]
	if err := json.Unmarshal([]byte(`[4,5,5]`), &zeroValue); err != nil {
		t.Fatalf("UnmarshalJSON into zero value failed: %v", err)
	}
	assertIntSetElements(t, zeroValue.ToSlice(), []int{4, 5})
}

func TestConcurrentHashSetBSONValue(t *testing.T) {
	original := NewConcurrentHashSet(1, 2, 2, 3)
	typ, data, err := original.MarshalBSONValue()
	if err != nil {
		t.Fatalf("MarshalBSONValue failed: %v", err)
	}
	if bson.Type(typ) != bson.TypeArray {
		t.Fatalf("MarshalBSONValue type = %v, want array", bson.Type(typ))
	}

	restored := NewConcurrentHashSet[int]()
	originalMap := restored.m
	if err := restored.UnmarshalBSONValue(typ, data); err != nil {
		t.Fatalf("UnmarshalBSONValue failed: %v", err)
	}
	assertIntSetElements(t, restored.ToSlice(), []int{1, 2, 3})
	if restored.m != originalMap {
		t.Error("UnmarshalBSONValue replaced the synchronized map instance")
	}

	before := restored.ToSlice()
	if err := restored.UnmarshalBSONValue(byte(bson.TypeArray), []byte{1, 2, 3}); err == nil {
		t.Fatal("UnmarshalBSONValue should reject malformed BSON")
	}
	assertIntSetElements(t, restored.ToSlice(), before)

	var zeroValue ConcurrentHashSet[int]
	if err := zeroValue.UnmarshalBSONValue(typ, data); err != nil {
		t.Fatalf("UnmarshalBSONValue into zero value failed: %v", err)
	}
	assertIntSetElements(t, zeroValue.ToSlice(), []int{1, 2, 3})
}

func TestConcurrentOperations(t *testing.T) {
	set := NewConcurrentHashSet[int]()
	const numOperations = 1000
	for i := 0; i < numOperations; i += 2 {
		set.Add(i)
	}

	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		for i := 1; i < numOperations; i += 2 {
			set.Add(i)
		}
	}()
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < numOperations; i += 2 {
			set.Remove(i)
		}
	}()
	close(start)
	wg.Wait()

	want := make([]int, 0, numOperations/2)
	for i := 1; i < numOperations; i += 2 {
		want = append(want, i)
	}
	assertIntSetElements(t, set.ToSlice(), want)
}
