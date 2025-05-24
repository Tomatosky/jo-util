package listUtil

import (
	"testing"
)

func TestNewArrayList(t *testing.T) {
	list := NewArrayList[int]()
	if list.Size() != 0 {
		t.Errorf("Expected size 0, got %d", list.Size())
	}
}

func TestAddAndGet(t *testing.T) {
	list := NewArrayList[string]()
	list.Add("a")
	list.Add("b")
	list.Add("c")

	if list.Size() != 3 {
		t.Errorf("Expected size 3, got %d", list.Size())
	}

	if val := list.Get(0); val != "a" {
		t.Errorf("Expected 'a' at index 0, got %s", val)
	}

	if val := list.Get(1); val != "b" {
		t.Errorf("Expected 'b' at index 1, got %s", val)
	}

	if val := list.Get(2); val != "c" {
		t.Errorf("Expected 'c' at index 2, got %s", val)
	}
}

func TestInsert(t *testing.T) {
	list := NewArrayList[int]()
	list.Add(1)
	list.Add(3)
	list.Insert(1, 2)

	if list.Size() != 3 {
		t.Errorf("Expected size 3, got %d", list.Size())
	}

	if val := list.Get(1); val != 2 {
		t.Errorf("Expected 2 at index 1, got %d", val)
	}
}

func TestNegativeIndex(t *testing.T) {
	list := NewArrayList[string]()
	list.Add("a")
	list.Add("b")
	list.Add("c")

	if val := list.Get(-1); val != "c" {
		t.Errorf("Expected 'c' at index -1, got %s", val)
	}

	if val := list.Get(-2); val != "b" {
		t.Errorf("Expected 'b' at index -2, got %s", val)
	}
}

func TestRemove(t *testing.T) {
	list := NewArrayList[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	removed := list.Remove(1)
	if removed != 2 {
		t.Errorf("Expected removed value 2, got %d", removed)
	}

	if list.Size() != 2 {
		t.Errorf("Expected size 2 after removal, got %d", list.Size())
	}

	if val := list.Get(1); val != 3 {
		t.Errorf("Expected 3 at index 1 after removal, got %d", val)
	}
}

func TestRange(t *testing.T) {
	list := NewArrayList[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	sum := 0
	list.Range(func(i int, v int) bool {
		sum += v
		return true
	})

	if sum != 6 {
		t.Errorf("Expected sum 6, got %d", sum)
	}

	// Test early termination
	sum = 0
	list.Range(func(i int, v int) bool {
		sum += v
		return i < 1 // stop after second element
	})

	if sum != 3 {
		t.Errorf("Expected sum 3 with early termination, got %d", sum)
	}
}

func TestContains(t *testing.T) {
	list := NewArrayList[string]()
	list.Add("apple")
	list.Add("banana")
	list.Add("cherry")

	if !list.Contains("banana") {
		t.Error("Expected to find 'banana'")
	}

	if list.Contains("orange") {
		t.Error("Did not expect to find 'orange'")
	}
}

func TestRemoveObject(t *testing.T) {
	list := NewArrayList[int]()
	list.Add(1)
	list.Add(2)
	list.Add(2)
	list.Add(3)
	list.Add(2)

	count := list.RemoveObject(2)
	if count != 3 {
		t.Errorf("Expected to remove 3 elements, removed %d", count)
	}

	if list.Size() != 2 {
		t.Errorf("Expected size 2 after removal, got %d", list.Size())
	}

	if list.Contains(2) {
		t.Error("Did not expect to find any 2s after removal")
	}
}

func TestToSlice(t *testing.T) {
	list := NewArrayList[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	slice := list.ToSlice()
	if len(slice) != 3 {
		t.Errorf("Expected slice length 3, got %d", len(slice))
	}

	if slice[0] != 1 || slice[1] != 2 || slice[2] != 3 {
		t.Error("Slice contents do not match expected values")
	}
}

func TestToString2(t *testing.T) {
	list := NewArrayList[string]()
	list.Add("a")
	list.Add("b")
	list.Add("c")

	str := list.ToString()
	expected := `["a","b","c"]`
	if str != expected {
		t.Errorf("Expected %s, got %s", expected, str)
	}
}

func TestPanicCases(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic did not occur")
		}
	}()

	list := NewArrayList[int]()
	list.Get(0) // should panic
}

func TestPanicCasesNegative(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic did not occur")
		}
	}()

	list := NewArrayList[int]()
	list.Add(1)
	list.Get(-2) // should panic (only -1 is valid for size 1)
}

func TestInsertPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic did not occur")
		}
	}()

	list := NewArrayList[int]()
	list.Insert(1, 1) // should panic (index out of range)
}

func TestRemovePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic did not occur")
		}
	}()

	list := NewArrayList[int]()
	list.Remove(0) // should panic
}
