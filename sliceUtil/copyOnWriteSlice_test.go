package sliceUtil

import (
	"encoding/json"
	"reflect"
	"sort"
	"sync"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// TestNewCopyOnWriteSlice 测试创建新实例
func TestNewCopyOnWriteSlice(t *testing.T) {
	slice := NewCopyOnWriteSlice[int]()
	if slice == nil {
		t.Fatal("NewCopyOnWriteSlice() returned nil")
	}
	if slice.Size() != 0 {
		t.Errorf("Expected size 0, got %d", slice.Size())
	}
}

// TestAdd 测试添加单个元素
func TestAdd(t *testing.T) {
	slice := NewCopyOnWriteSlice[int]()

	slice.Add(1)
	if slice.Size() != 1 {
		t.Errorf("Expected size 1, got %d", slice.Size())
	}
	if slice.Get(0) != 1 {
		t.Errorf("Expected element 1, got %d", slice.Get(0))
	}

	slice.Add(2)
	slice.Add(3)
	if slice.Size() != 3 {
		t.Errorf("Expected size 3, got %d", slice.Size())
	}
	if slice.Get(1) != 2 {
		t.Errorf("Expected element 2 at index 1, got %d", slice.Get(1))
	}
	if slice.Get(2) != 3 {
		t.Errorf("Expected element 3 at index 2, got %d", slice.Get(2))
	}
}

// TestAddAll 测试批量添加元素
func TestAddAll(t *testing.T) {
	slice := NewCopyOnWriteSlice[int]()

	slice.AddAll(1, 2, 3, 4, 5)
	if slice.Size() != 5 {
		t.Errorf("Expected size 5, got %d", slice.Size())
	}

	for i := 0; i < 5; i++ {
		if slice.Get(i) != i+1 {
			t.Errorf("Expected element %d at index %d, got %d", i+1, i, slice.Get(i))
		}
	}

	// 测试空AddAll
	slice.AddAll()
	if slice.Size() != 5 {
		t.Errorf("Expected size 5 after empty AddAll, got %d", slice.Size())
	}

	// 继续添加
	slice.AddAll(6, 7)
	if got, want := slice.ToSlice(), []int{1, 2, 3, 4, 5, 6, 7}; !reflect.DeepEqual(got, want) {
		t.Errorf("AddAll result = %v, want %v", got, want)
	}
}

// TestInsert 测试插入元素
func TestInsert(t *testing.T) {
	slice := NewCopyOnWriteSlice[int]()
	slice.AddAll(1, 3, 5)

	// 在中间插入
	slice.Insert(1, 2)
	if slice.Size() != 4 {
		t.Errorf("Expected size 4, got %d", slice.Size())
	}
	if slice.Get(1) != 2 {
		t.Errorf("Expected element 2 at index 1, got %d", slice.Get(1))
	}

	// 在开头插入
	slice.Insert(0, 0)
	if slice.Size() != 5 {
		t.Errorf("Expected size 5, got %d", slice.Size())
	}
	if slice.Get(0) != 0 {
		t.Errorf("Expected element 0 at index 0, got %d", slice.Get(0))
	}

	// 在末尾插入
	slice.Insert(slice.Size(), 6)
	if got, want := slice.ToSlice(), []int{0, 1, 2, 3, 5, 6}; !reflect.DeepEqual(got, want) {
		t.Errorf("Insert result = %v, want %v", got, want)
	}
}

// TestInsertOutOfRange 测试插入越界
func TestInsertOutOfRange(t *testing.T) {
	slice := NewCopyOnWriteSlice[int]()
	slice.AddAll(1, 2, 3)

	defer func() {
		if r := recover(); r != "index out of range" {
			t.Errorf("Insert(-1) panic = %v, want %q", r, "index out of range")
		}
	}()
	slice.Insert(-1, 0)
}

// TestInsertOutOfRangeHigh 测试插入索引过大
func TestInsertOutOfRangeHigh(t *testing.T) {
	slice := NewCopyOnWriteSlice[int]()
	slice.AddAll(1, 2, 3)

	defer func() {
		if r := recover(); r != "index out of range" {
			t.Errorf("Insert(4) panic = %v, want %q", r, "index out of range")
		}
	}()
	slice.Insert(4, 0)
}

// TestGet 测试获取元素
func TestGet(t *testing.T) {
	slice := NewCopyOnWriteSlice[int]()
	slice.AddAll(10, 20, 30, 40, 50)

	// 正常索引
	if slice.Get(0) != 10 {
		t.Errorf("Expected 10, got %d", slice.Get(0))
	}
	if slice.Get(4) != 50 {
		t.Errorf("Expected 50, got %d", slice.Get(4))
	}

	// 负数索引
	if slice.Get(-1) != 50 {
		t.Errorf("Expected 50 for index -1, got %d", slice.Get(-1))
	}
	if slice.Get(-5) != 10 {
		t.Errorf("Expected 10 for index -5, got %d", slice.Get(-5))
	}
}

// TestGetOutOfRange 测试获取越界
func TestGetOutOfRange(t *testing.T) {
	slice := NewCopyOnWriteSlice[int]()
	slice.AddAll(1, 2, 3)

	defer func() {
		if r := recover(); r != "index out of range" {
			t.Errorf("Get(3) panic = %v, want %q", r, "index out of range")
		}
	}()
	slice.Get(3)
}

// TestGetNegativeOutOfRange 测试负索引越界
func TestGetNegativeOutOfRange(t *testing.T) {
	slice := NewCopyOnWriteSlice[int]()
	slice.AddAll(1, 2, 3)

	defer func() {
		if r := recover(); r != "index out of range" {
			t.Errorf("Get(-10) panic = %v, want %q", r, "index out of range")
		}
	}()
	slice.Get(-10)
}

// TestRemove 测试移除元素
func TestRemove(t *testing.T) {
	slice := NewCopyOnWriteSlice[int]()
	slice.AddAll(1, 2, 3, 4, 5)

	// 移除中间元素
	removed := slice.Remove(2)
	if removed != 3 {
		t.Errorf("Expected removed element 3, got %d", removed)
	}
	if slice.Size() != 4 {
		t.Errorf("Expected size 4, got %d", slice.Size())
	}
	if slice.Get(2) != 4 {
		t.Errorf("Expected element 4 at index 2, got %d", slice.Get(2))
	}

	// 移除第一个元素
	removed = slice.Remove(0)
	if removed != 1 {
		t.Errorf("Expected removed element 1, got %d", removed)
	}
	if slice.Size() != 3 {
		t.Errorf("Expected size 3, got %d", slice.Size())
	}

	// 移除最后一个元素
	removed = slice.Remove(slice.Size() - 1)
	if removed != 5 {
		t.Errorf("Expected removed element 5, got %d", removed)
	}
	if slice.Size() != 2 {
		t.Errorf("Expected size 2, got %d", slice.Size())
	}
	if got, want := slice.ToSlice(), []int{2, 4}; !reflect.DeepEqual(got, want) {
		t.Errorf("Remove result = %v, want %v", got, want)
	}
}

// TestRemoveOutOfRange 测试移除越界
func TestRemoveOutOfRange(t *testing.T) {
	for _, tt := range []struct {
		name  string
		index int
	}{
		{name: "negative", index: -1},
		{name: "equal to size", index: 3},
		{name: "greater than size", index: 5},
	} {
		t.Run(tt.name, func(t *testing.T) {
			slice := NewCopyOnWriteSlice[int]()
			slice.AddAll(1, 2, 3)
			defer func() {
				if r := recover(); r != "index out of range" {
					t.Errorf("Remove(%d) panic = %v, want %q", tt.index, r, "index out of range")
				}
			}()
			slice.Remove(tt.index)
		})
	}
}

// TestSize 测试大小
func TestSize(t *testing.T) {
	slice := NewCopyOnWriteSlice[int]()

	if slice.Size() != 0 {
		t.Errorf("Expected size 0, got %d", slice.Size())
	}

	slice.Add(1)
	if slice.Size() != 1 {
		t.Errorf("Expected size 1, got %d", slice.Size())
	}

	slice.AddAll(2, 3, 4)
	if slice.Size() != 4 {
		t.Errorf("Expected size 4, got %d", slice.Size())
	}

	slice.Remove(0)
	if slice.Size() != 3 {
		t.Errorf("Expected size 3, got %d", slice.Size())
	}
}

// TestRange 测试遍历
func TestRange(t *testing.T) {
	slice := NewCopyOnWriteSlice[int]()
	slice.AddAll(1, 2, 3, 4, 5)

	// 完整遍历
	count := 0
	slice.Range(func(i int, v int) bool {
		count++
		if v != i+1 {
			t.Errorf("Expected value %d at index %d, got %d", i+1, i, v)
		}
		return true
	})
	if count != 5 {
		t.Errorf("Expected to iterate 5 times, got %d", count)
	}

	// 提前终止
	count = 0
	slice.Range(func(i int, v int) bool {
		count++
		return i < 2 // 只遍历前3个元素
	})
	if count != 3 {
		t.Errorf("Expected to iterate 3 times, got %d", count)
	}
}

// TestContains 测试包含检查
func TestContains(t *testing.T) {
	slice := NewCopyOnWriteSlice[int]()
	slice.AddAll(1, 2, 3, 4, 5)

	if !slice.Contains(3) {
		t.Error("Expected slice to contain 3")
	}

	if slice.Contains(10) {
		t.Error("Expected slice not to contain 10")
	}

	// 空切片
	emptySlice := NewCopyOnWriteSlice[int]()
	if emptySlice.Contains(1) {
		t.Error("Expected empty slice not to contain any element")
	}
}

// TestRemoveObject 测试删除所有匹配元素
func TestRemoveObject(t *testing.T) {
	slice := NewCopyOnWriteSlice[int]()
	slice.AddAll(1, 2, 3, 2, 4, 2, 5)

	// 删除所有的2
	count := slice.RemoveObject(2)
	if count != 3 {
		t.Errorf("Expected to remove 3 elements, got %d", count)
	}
	if slice.Size() != 4 {
		t.Errorf("Expected size 4, got %d", slice.Size())
	}
	if slice.Contains(2) {
		t.Error("Slice should not contain 2 after removal")
	}
	if got, want := slice.ToSlice(), []int{1, 3, 4, 5}; !reflect.DeepEqual(got, want) {
		t.Errorf("RemoveObject result = %v, want %v", got, want)
	}

	// 删除不存在的元素
	count = slice.RemoveObject(100)
	if count != 0 {
		t.Errorf("Expected to remove 0 elements, got %d", count)
	}
	if slice.Size() != 4 {
		t.Errorf("Expected size 4, got %d", slice.Size())
	}
}

// TestToSlice 测试转换为切片
func TestToSlice(t *testing.T) {
	slice := NewCopyOnWriteSlice[int]()
	slice.AddAll(1, 2, 3, 4, 5)

	result := slice.ToSlice()
	if len(result) != 5 {
		t.Errorf("Expected slice length 5, got %d", len(result))
	}

	for i := 0; i < 5; i++ {
		if result[i] != i+1 {
			t.Errorf("Expected element %d at index %d, got %d", i+1, i, result[i])
		}
	}

	// 确保返回的是副本，修改不影响原始数据
	result[0] = 100
	if slice.Get(0) == 100 {
		t.Error("ToSlice should return a copy, not the original slice")
	}
}

// TestCOWToString 测试转换为字符串
func TestCOWToString(t *testing.T) {
	slice := NewCopyOnWriteSlice[int]()
	slice.AddAll(1, 2, 3)

	str := slice.ToString()
	expected := "[1,2,3]"
	if str != expected {
		t.Errorf("Expected string %s, got %s", expected, str)
	}

	// 空切片
	emptySlice := NewCopyOnWriteSlice[int]()
	str = emptySlice.ToString()
	expected = "[]"
	if str != expected {
		t.Errorf("Expected string %s, got %s", expected, str)
	}
}

// TestMarshalJSON 测试JSON序列化
func TestMarshalJSON(t *testing.T) {
	slice := NewCopyOnWriteSlice[int]()
	slice.AddAll(1, 2, 3, 4, 5)

	data, err := json.Marshal(slice)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	expected := "[1,2,3,4,5]"
	if string(data) != expected {
		t.Errorf("Expected JSON %s, got %s", expected, string(data))
	}
}

// TestUnmarshalJSON 测试JSON反序列化
func TestUnmarshalJSON(t *testing.T) {
	jsonData := []byte("[1,2,3,4,5]")

	slice := NewCopyOnWriteSlice[int]()
	err := json.Unmarshal(jsonData, slice)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	if slice.Size() != 5 {
		t.Errorf("Expected size 5, got %d", slice.Size())
	}

	for i := 0; i < 5; i++ {
		if slice.Get(i) != i+1 {
			t.Errorf("Expected element %d at index %d, got %d", i+1, i, slice.Get(i))
		}
	}
}

func TestUnmarshalJSONRejectsInvalidDataWithoutChangingSlice(t *testing.T) {
	slice := NewCopyOnWriteSlice[int]()
	slice.AddAll(1, 2, 3)
	if err := json.Unmarshal([]byte(`{"invalid":true}`), slice); err == nil {
		t.Fatal("UnmarshalJSON should reject a non-array value")
	}
	if got, want := slice.ToSlice(), []int{1, 2, 3}; !reflect.DeepEqual(got, want) {
		t.Errorf("slice after invalid JSON = %v, want %v", got, want)
	}
}

func TestUnmarshalJSONIsAtomicForConcurrentReaders(t *testing.T) {
	slice := NewCopyOnWriteSlice[int]()
	slice.AddAll(1, 2, 3)
	stop := make(chan struct{})
	errors := make(chan []int, 1)
	var readers sync.WaitGroup
	for i := 0; i < 8; i++ {
		readers.Add(1)
		go func() {
			defer readers.Done()
			for {
				select {
				case <-stop:
					return
				default:
					snapshot := slice.ToSlice()
					if !reflect.DeepEqual(snapshot, []int{1, 2, 3}) && !reflect.DeepEqual(snapshot, []int{4, 5, 6}) {
						select {
						case errors <- snapshot:
						default:
						}
						return
					}
				}
			}
		}()
	}

	for i := 0; i < 100; i++ {
		data := []byte(`[1,2,3]`)
		if i%2 == 1 {
			data = []byte(`[4,5,6]`)
		}
		if err := json.Unmarshal(data, slice); err != nil {
			t.Fatalf("UnmarshalJSON failed: %v", err)
		}
	}
	close(stop)
	readers.Wait()
	close(errors)
	for invalid := range errors {
		t.Errorf("reader observed partial replacement: %v", invalid)
	}
	if got, want := slice.ToSlice(), []int{4, 5, 6}; !reflect.DeepEqual(got, want) {
		t.Errorf("final slice = %v, want %v", got, want)
	}
}

// TestJSONRoundTrip 测试JSON序列化往返
func TestJSONRoundTrip(t *testing.T) {
	original := NewCopyOnWriteSlice[string]()
	original.AddAll("hello", "world", "test")

	// 序列化
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	restored := NewCopyOnWriteSlice[string]()
	err = json.Unmarshal(data, restored)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if restored.Size() != original.Size() {
		t.Errorf("Expected size %d, got %d", original.Size(), restored.Size())
	}

	for i := 0; i < original.Size(); i++ {
		if restored.Get(i) != original.Get(i) {
			t.Errorf("Element mismatch at index %d: expected %s, got %s",
				i, original.Get(i), restored.Get(i))
		}
	}
}

// TestMarshalBSONValue 测试BSON序列化
func TestMarshalBSONValue(t *testing.T) {
	slice := NewCopyOnWriteSlice[int]()
	slice.AddAll(1, 2, 3, 4, 5)

	bsonType, data, err := slice.MarshalBSONValue()
	if err != nil {
		t.Fatalf("MarshalBSONValue failed: %v", err)
	}

	if bson.Type(bsonType) != bson.TypeArray {
		t.Errorf("MarshalBSONValue type = %v, want array", bson.Type(bsonType))
	}
	var decoded []int
	if err := bson.UnmarshalValue(bson.Type(bsonType), data, &decoded); err != nil {
		t.Fatalf("cannot decode marshaled BSON: %v", err)
	}
	if want := []int{1, 2, 3, 4, 5}; !reflect.DeepEqual(decoded, want) {
		t.Errorf("marshaled BSON elements = %v, want %v", decoded, want)
	}
}

// TestUnmarshalBSONValue 测试BSON反序列化
func TestUnmarshalBSONValue(t *testing.T) {
	// 创建原始切片
	original := NewCopyOnWriteSlice[int]()
	original.AddAll(10, 20, 30)

	// 序列化
	bsonType, data, err := original.MarshalBSONValue()
	if err != nil {
		t.Fatalf("MarshalBSONValue failed: %v", err)
	}

	// 反序列化
	restored := NewCopyOnWriteSlice[int]()
	err = restored.UnmarshalBSONValue(bsonType, data)
	if err != nil {
		t.Fatalf("UnmarshalBSONValue failed: %v", err)
	}

	// 验证
	if restored.Size() != original.Size() {
		t.Errorf("Expected size %d, got %d", original.Size(), restored.Size())
	}

	for i := 0; i < original.Size(); i++ {
		if restored.Get(i) != original.Get(i) {
			t.Errorf("Element mismatch at index %d: expected %d, got %d",
				i, original.Get(i), restored.Get(i))
		}
	}
}

func TestUnmarshalBSONValueRejectsInvalidDataWithoutChangingSlice(t *testing.T) {
	slice := NewCopyOnWriteSlice[int]()
	slice.AddAll(1, 2, 3)
	if err := slice.UnmarshalBSONValue(byte(bson.TypeArray), []byte{1, 2, 3}); err == nil {
		t.Fatal("UnmarshalBSONValue should reject malformed BSON")
	}
	if got, want := slice.ToSlice(), []int{1, 2, 3}; !reflect.DeepEqual(got, want) {
		t.Errorf("slice after invalid BSON = %v, want %v", got, want)
	}
}

// TestBSONRoundTrip 测试BSON序列化往返
func TestBSONRoundTrip(t *testing.T) {
	original := NewCopyOnWriteSlice[string]()
	original.AddAll("foo", "bar", "baz")

	// 序列化
	bsonType, data, err := original.MarshalBSONValue()
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	restored := NewCopyOnWriteSlice[string]()
	err = restored.UnmarshalBSONValue(bsonType, data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if restored.Size() != original.Size() {
		t.Errorf("Expected size %d, got %d", original.Size(), restored.Size())
	}

	for i := 0; i < original.Size(); i++ {
		if restored.Get(i) != original.Get(i) {
			t.Errorf("Element mismatch at index %d: expected %s, got %s",
				i, original.Get(i), restored.Get(i))
		}
	}
}

// TestConcurrentAdd 测试并发添加
func TestConcurrentAdd(t *testing.T) {
	slice := NewCopyOnWriteSlice[int]()
	numGoroutines := 100
	elementsPerGoroutine := 10

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(start int) {
			defer wg.Done()
			for j := 0; j < elementsPerGoroutine; j++ {
				slice.Add(start + j)
			}
		}(i * elementsPerGoroutine)
	}

	wg.Wait()

	expectedSize := numGoroutines * elementsPerGoroutine
	if slice.Size() != expectedSize {
		t.Errorf("Expected size %d, got %d", expectedSize, slice.Size())
	}
	got := slice.ToSlice()
	sort.Ints(got)
	want := make([]int, expectedSize)
	for i := range want {
		want[i] = i
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("concurrently added elements = %v, want 0..%d", got, expectedSize-1)
	}
}

// TestConcurrentReadWrite 测试并发读写
func TestConcurrentReadWrite(t *testing.T) {
	slice := NewCopyOnWriteSlice[int]()
	slice.AddAll(1, 2, 3, 4, 5)

	numReaders := 50
	numWriters := 50
	iterations := 100

	var wg sync.WaitGroup
	wg.Add(numReaders + numWriters)
	errors := make(chan string, numReaders*iterations*3)

	// 读取协程
	for i := 0; i < numReaders; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				if got := slice.Get(0); got != 1 {
					errors <- "first element changed during append-only writes"
				}
				if !slice.Contains(3) {
					errors <- "stable element disappeared during append-only writes"
				}
				snapshot := slice.ToSlice()
				if len(snapshot) < 5 || !reflect.DeepEqual(snapshot[:5], []int{1, 2, 3, 4, 5}) {
					errors <- "snapshot lost its stable prefix"
				}
			}
		}()
	}

	// 写入协程
	for i := 0; i < numWriters; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				slice.Add(id*1000 + j)
			}
		}(i)
	}

	wg.Wait()
	close(errors)
	for message := range errors {
		t.Error(message)
	}

	// 验证最终大小
	expectedSize := 5 + numWriters*iterations
	if got := slice.Size(); got != expectedSize {
		t.Errorf("final size = %d, want %d", got, expectedSize)
	}
	got := slice.ToSlice()[5:]
	sort.Ints(got)
	want := make([]int, 0, numWriters*iterations)
	for writer := 0; writer < numWriters; writer++ {
		for iteration := 0; iteration < iterations; iteration++ {
			want = append(want, writer*1000+iteration)
		}
	}
	sort.Ints(want)
	if !reflect.DeepEqual(got, want) {
		t.Error("concurrent writes did not preserve the exact submitted elements")
	}
}

// TestConcurrentRemove 测试并发删除
func TestConcurrentRemove(t *testing.T) {
	slice := NewCopyOnWriteSlice[int]()

	// 添加足够多的元素
	for i := 0; i < 1000; i++ {
		slice.Add(i)
	}

	numGoroutines := 10
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// 并发删除
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				slice.Remove(0)
			}
		}()
	}

	wg.Wait()

	expectedSize := 1000 - numGoroutines*50
	if slice.Size() != expectedSize {
		t.Errorf("Expected size %d, got %d", expectedSize, slice.Size())
	}
	got := slice.ToSlice()
	want := make([]int, expectedSize)
	for i := range want {
		want[i] = i + numGoroutines*50
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("remaining elements = %v, want %v", got, want)
	}
}

// TestConcurrentRange 测试并发遍历
func TestConcurrentRange(t *testing.T) {
	slice := NewCopyOnWriteSlice[int]()
	slice.AddAll(1, 2, 3, 4, 5)

	numGoroutines := 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2)
	snapshots := make([][]int, numGoroutines)

	// 遍历协程
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			slice.Range(func(i int, v int) bool {
				snapshots[id] = append(snapshots[id], v)
				return true
			})
		}(i)
	}

	// 同时进行修改
	for i := 0; i < numGoroutines; i++ {
		go func(val int) {
			defer wg.Done()
			slice.Add(val)
		}(i + 100)
	}

	wg.Wait()

	// 验证没有panic，并且大小正确
	expectedSize := 5 + numGoroutines
	if slice.Size() != expectedSize {
		t.Errorf("Expected size %d, got %d", expectedSize, slice.Size())
	}
	for id, snapshot := range snapshots {
		if len(snapshot) < 5 || !reflect.DeepEqual(snapshot[:5], []int{1, 2, 3, 4, 5}) {
			t.Errorf("snapshot %d has invalid stable prefix: %v", id, snapshot)
			continue
		}
		seen := make(map[int]struct{}, len(snapshot)-5)
		for _, value := range snapshot[5:] {
			if value < 100 || value >= 100+numGoroutines {
				t.Errorf("snapshot %d contains unexpected value %d", id, value)
			}
			if _, duplicate := seen[value]; duplicate {
				t.Errorf("snapshot %d contains duplicate value %d", id, value)
			}
			seen[value] = struct{}{}
		}
	}
}

// TestDifferentTypes 测试不同类型
func TestDifferentTypes(t *testing.T) {
	// 测试string类型
	stringSlice := NewCopyOnWriteSlice[string]()
	stringSlice.AddAll("hello", "world")
	if stringSlice.Size() != 2 {
		t.Errorf("String slice: expected size 2, got %d", stringSlice.Size())
	}
	if !stringSlice.Contains("hello") {
		t.Error("String slice should contain 'hello'")
	}

	// 测试float64类型
	floatSlice := NewCopyOnWriteSlice[float64]()
	floatSlice.AddAll(1.1, 2.2, 3.3)
	if floatSlice.Size() != 3 {
		t.Errorf("Float slice: expected size 3, got %d", floatSlice.Size())
	}
	if floatSlice.Get(1) != 2.2 {
		t.Errorf("Float slice: expected 2.2, got %f", floatSlice.Get(1))
	}

	// 测试bool类型
	boolSlice := NewCopyOnWriteSlice[bool]()
	boolSlice.AddAll(true, false, true)
	if boolSlice.Size() != 3 {
		t.Errorf("Bool slice: expected size 3, got %d", boolSlice.Size())
	}
	if !boolSlice.Contains(true) {
		t.Error("Bool slice should contain true")
	}
}

// TestEmptySliceOperations 测试空切片操作
func TestEmptySliceOperations(t *testing.T) {
	slice := NewCopyOnWriteSlice[int]()

	// 测试空切片的各种操作
	if slice.Size() != 0 {
		t.Errorf("Empty slice: expected size 0, got %d", slice.Size())
	}

	if slice.Contains(1) {
		t.Error("Empty slice should not contain any element")
	}

	result := slice.ToSlice()
	if len(result) != 0 {
		t.Errorf("Empty slice ToSlice: expected length 0, got %d", len(result))
	}

	str := slice.ToString()
	if str != "[]" {
		t.Errorf("Empty slice ToString: expected '[]', got '%s'", str)
	}

	count := slice.RemoveObject(1)
	if count != 0 {
		t.Errorf("Empty slice RemoveObject: expected 0, got %d", count)
	}

	rangeCount := 0
	slice.Range(func(i int, v int) bool {
		rangeCount++
		return true
	})
	if rangeCount != 0 {
		t.Errorf("Empty slice Range: expected 0 iterations, got %d", rangeCount)
	}
}

// BenchmarkAdd 添加性能测试
func BenchmarkAdd(b *testing.B) {
	slice := NewCopyOnWriteSlice[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slice.Add(i)
	}
}

// BenchmarkGet 读取性能测试
func BenchmarkGet(b *testing.B) {
	slice := NewCopyOnWriteSlice[int]()
	for i := 0; i < 1000; i++ {
		slice.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = slice.Get(i % 1000)
	}
}

// BenchmarkConcurrentRead 并发读取性能测试
func BenchmarkConcurrentRead(b *testing.B) {
	slice := NewCopyOnWriteSlice[int]()
	for i := 0; i < 1000; i++ {
		slice.Add(i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_ = slice.Get(i % 1000)
			i++
		}
	})
}
