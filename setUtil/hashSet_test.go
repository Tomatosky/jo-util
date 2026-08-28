package setUtil

import (
	"encoding/json"
	"reflect"
	"sort"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func assertIntSetElements(t *testing.T, got []int, want []int) {
	t.Helper()

	got = append([]int(nil), got...)
	want = append([]int(nil), want...)
	sort.Ints(got)
	sort.Ints(want)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("set elements = %v, want %v", got, want)
	}
}

func TestNewHashSet(t *testing.T) {
	// 测试空集合
	emptySet := NewHashSet[int]()
	if emptySet.Size() != 0 {
		t.Errorf("Expected empty set size 0, got %d", emptySet.Size())
	}

	// 测试带初始元素的集合
	set := NewHashSet(1, 2, 3, 2)
	if set.Size() != 3 {
		t.Errorf("Expected set size 3, got %d", set.Size())
	}
	assertIntSetElements(t, set.ToSlice(), []int{1, 2, 3})
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
	visited := make([]int, 0, set.Size())

	set.Range(func(n int) bool {
		visited = append(visited, n)
		return true
	})
	assertIntSetElements(t, visited, []int{1, 2, 3, 4, 5})

	count := 0
	set.Range(func(int) bool {
		count++
		return count < 3
	})
	if count != 3 {
		t.Errorf("Range should stop on the third callback, got %d callbacks", count)
	}
}

func TestToSlice(t *testing.T) {
	elements := []int{1, 2, 3, 4}
	set := NewHashSet(elements...)
	slice := set.ToSlice()

	assertIntSetElements(t, slice, elements)
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
		t.Fatalf("Failed to unmarshal set string: %v", err)
	}
	sort.Strings(slice)
	if want := []string{"apple", "banana", "cherry"}; !reflect.DeepEqual(slice, want) {
		t.Errorf("ToString elements = %v, want %v", slice, want)
	}

	if got := NewHashSet[int]().ToString(); got != "[]" {
		t.Errorf("empty set ToString() = %q, want []", got)
	}
}

func TestHashSetJSON(t *testing.T) {
	set := NewHashSet(1, 2, 2, 3)
	data, err := json.Marshal(set)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	var encoded []int
	if err := json.Unmarshal(data, &encoded); err != nil {
		t.Fatalf("marshaled data is not a JSON array: %v", err)
	}
	assertIntSetElements(t, encoded, []int{1, 2, 3})

	restored := NewHashSet(99)
	if err := json.Unmarshal([]byte(`[3,2,2,1]`), restored); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	assertIntSetElements(t, restored.ToSlice(), []int{1, 2, 3})

	before := restored.ToSlice()
	if err := json.Unmarshal([]byte(`{"invalid":true}`), restored); err == nil {
		t.Fatal("UnmarshalJSON should reject a non-array value")
	}
	assertIntSetElements(t, restored.ToSlice(), before)

	if err := json.Unmarshal([]byte(`[]`), restored); err != nil {
		t.Fatalf("UnmarshalJSON empty array failed: %v", err)
	}
	assertIntSetElements(t, restored.ToSlice(), []int{})
}

func TestHashSetBSONValue(t *testing.T) {
	original := NewHashSet(1, 2, 2, 3)
	typ, data, err := original.MarshalBSONValue()
	if err != nil {
		t.Fatalf("MarshalBSONValue failed: %v", err)
	}
	if bson.Type(typ) != bson.TypeArray {
		t.Fatalf("MarshalBSONValue type = %v, want array", bson.Type(typ))
	}

	restored := NewHashSet(99)
	if err := restored.UnmarshalBSONValue(typ, data); err != nil {
		t.Fatalf("UnmarshalBSONValue failed: %v", err)
	}
	assertIntSetElements(t, restored.ToSlice(), []int{1, 2, 3})

	before := restored.ToSlice()
	if err := restored.UnmarshalBSONValue(byte(bson.TypeArray), []byte{1, 2, 3}); err == nil {
		t.Fatal("UnmarshalBSONValue should reject malformed BSON")
	}
	assertIntSetElements(t, restored.ToSlice(), before)
}
