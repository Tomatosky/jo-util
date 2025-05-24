package listUtil

import (
	"encoding/json"
	"sync"
	"testing"
)

func TestNewCopyOnWriteArrayList(t *testing.T) {
	list := NewCopyOnWriteArrayList[int]()
	if list.Size() != 0 {
		t.Errorf("Expected empty list, got size %d", list.Size())
	}
}

func TestAddAndGet3(t *testing.T) {
	list := NewCopyOnWriteArrayList[string]()
	list.Add("first")
	list.Add("second")

	if size := list.Size(); size != 2 {
		t.Errorf("Expected size 2, got %d", size)
	}

	if val := list.Get(0); val != "first" {
		t.Errorf("Expected 'first', got '%s'", val)
	}

	if val := list.Get(1); val != "second" {
		t.Errorf("Expected 'second', got '%s'", val)
	}
}

func TestInsert3(t *testing.T) {
	list := NewCopyOnWriteArrayList[int]()
	list.Add(1)
	list.Add(3)
	list.Insert(1, 2)

	if size := list.Size(); size != 3 {
		t.Errorf("Expected size 3, got %d", size)
	}

	if val := list.Get(1); val != 2 {
		t.Errorf("Expected 2, got %d", val)
	}
}

func TestInsertPanic3(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for out of range index")
		}
	}()

	list := NewCopyOnWriteArrayList[int]()
	list.Insert(1, 1) // Should panic
}

func TestRemove3(t *testing.T) {
	list := NewCopyOnWriteArrayList[string]()
	list.Add("a")
	list.Add("b")
	list.Add("c")

	removed := list.Remove(1)
	if removed != "b" {
		t.Errorf("Expected 'b' removed, got '%s'", removed)
	}

	if size := list.Size(); size != 2 {
		t.Errorf("Expected size 2, got %d", size)
	}

	if val := list.Get(1); val != "c" {
		t.Errorf("Expected 'c', got '%s'", val)
	}
}

func TestRemovePanic3(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for out of range index")
		}
	}()

	list := NewCopyOnWriteArrayList[int]()
	list.Remove(0) // Should panic
}

func TestContains3(t *testing.T) {
	list := NewCopyOnWriteArrayList[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	if !list.Contains(2) {
		t.Error("Expected to contain 2")
	}

	if list.Contains(4) {
		t.Error("Expected not to contain 4")
	}
}

func TestRemoveObject3(t *testing.T) {
	list := NewCopyOnWriteArrayList[int]()
	list.Add(1)
	list.Add(2)
	list.Add(2)
	list.Add(3)
	list.Add(2)

	count := list.RemoveObject(2)
	if count != 3 {
		t.Errorf("Expected 3 removals, got %d", count)
	}

	if list.Size() != 2 {
		t.Errorf("Expected size 2, got %d", list.Size())
	}

	if list.Contains(2) {
		t.Error("Expected no 2s remaining")
	}
}

func TestRange3(t *testing.T) {
	list := NewCopyOnWriteArrayList[int]()
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
	count := 0
	list.Range(func(i int, v int) bool {
		count++
		return count < 2
	})

	if count != 2 {
		t.Errorf("Expected to process 2 elements, got %d", count)
	}
}

func TestToSlice3(t *testing.T) {
	list := NewCopyOnWriteArrayList[string]()
	list.Add("a")
	list.Add("b")
	list.Add("c")

	slice := list.ToSlice()
	if len(slice) != 3 {
		t.Errorf("Expected slice length 3, got %d", len(slice))
	}

	// Modify the returned slice should not affect the original list
	slice[0] = "x"
	if list.Get(0) != "a" {
		t.Error("Modifying returned slice affected original list")
	}
}

func TestToString3(t *testing.T) {
	list := NewCopyOnWriteArrayList[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	str := list.ToString()
	expected := "[1,2,3]"
	if str != expected {
		t.Errorf("Expected '%s', got '%s'", expected, str)
	}

	// Test JSON unmarshaling
	var data []int
	err := json.Unmarshal([]byte(str), &data)
	if err != nil {
		t.Errorf("Failed to unmarshal JSON: %v", err)
	}

	if len(data) != 3 || data[0] != 1 || data[1] != 2 || data[2] != 3 {
		t.Error("Unmarshaled data doesn't match original")
	}
}

func TestConcurrentAccess(t *testing.T) {
	list := NewCopyOnWriteArrayList[int]()
	var wg sync.WaitGroup

	// Concurrent writers
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			list.Add(val)
		}(i)
	}

	// Concurrent readers
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = list.Size()
			_ = list.Contains(42)
		}()
	}

	wg.Wait()

	if list.Size() != 100 {
		t.Errorf("Expected size 100, got %d", list.Size())
	}
}

func TestNegativeIndex3(t *testing.T) {
	list := NewCopyOnWriteArrayList[string]()
	list.Add("a")
	list.Add("b")
	list.Add("c")

	if val := list.Get(-1); val != "c" {
		t.Errorf("Expected 'c', got '%s'", val)
	}

	if val := list.Get(-2); val != "b" {
		t.Errorf("Expected 'b', got '%s'", val)
	}
}

func TestNegativeIndexPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for out of range negative index")
		}
	}()

	list := NewCopyOnWriteArrayList[int]()
	list.Add(1)
	list.Get(-2) // Should panic
}
