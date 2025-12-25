package mapUtil

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"testing"

	"github.com/Tomatosky/jo-util/dateUtil"
	"go.mongodb.org/mongo-driver/bson"
)

// ============ 基础功能测试 ============

func TestNewConcurrentSkipListMap(t *testing.T) {
	// 测试空初始化
	emptyMap := NewConcurrentSkipListMap[string, int]()
	if emptyMap.Size() != 0 {
		t.Errorf("Expected empty map size 0, got %d", emptyMap.Size())
	}

	// 测试带初始值的初始化
	initData := map[string]int{"a": 1, "b": 2, "c": 3}
	initMap := NewConcurrentSkipListMap(initData)
	if initMap.Size() != 3 {
		t.Errorf("Expected initialized map size 3, got %d", initMap.Size())
	}
	if val := initMap.Get("a"); val != 1 {
		t.Errorf("Expected value 1 for key 'a', got %d", val)
	}
	if val := initMap.Get("b"); val != 2 {
		t.Errorf("Expected value 2 for key 'b', got %d", val)
	}
}

func TestSkipListMapGetPutRemove(t *testing.T) {
	csm := NewConcurrentSkipListMap[string, string]()

	// 测试Put和Get
	csm.Put("key1", "value1")
	if val := csm.Get("key1"); val != "value1" {
		t.Errorf("Expected 'value1', got '%s'", val)
	}

	// 测试更新已存在的key
	csm.Put("key1", "value2")
	if val := csm.Get("key1"); val != "value2" {
		t.Errorf("Expected 'value2' after update, got '%s'", val)
	}

	// 测试获取不存在的key
	if val := csm.Get("nonexistent"); val != "" {
		t.Errorf("Expected empty value for nonexistent key, got '%s'", val)
	}

	// 测试Remove
	csm.Remove("key1")
	if val := csm.Get("key1"); val != "" {
		t.Errorf("Expected empty value after remove, got '%s'", val)
	}

	// 测试删除不存在的key
	csm.Remove("nonexistent") // 不应panic
	if csm.Size() != 0 {
		t.Errorf("Expected size 0 after removing all, got %d", csm.Size())
	}
}

func TestSkipListMapSizeAndContainsKey(t *testing.T) {
	csm := NewConcurrentSkipListMap[int, bool]()

	// 初始大小
	if csm.Size() != 0 {
		t.Errorf("Expected size 0, got %d", csm.Size())
	}

	// 添加元素后测试
	csm.Put(1, true)
	if csm.Size() != 1 {
		t.Errorf("Expected size 1, got %d", csm.Size())
	}
	if !csm.ContainsKey(1) {
		t.Error("Expected to contain key 1")
	}
	if csm.ContainsKey(2) {
		t.Error("Should not contain key 2")
	}

	// 添加更多元素
	csm.Put(2, false)
	csm.Put(3, true)
	if csm.Size() != 3 {
		t.Errorf("Expected size 3, got %d", csm.Size())
	}

	// 删除元素
	csm.Remove(2)
	if csm.Size() != 2 {
		t.Errorf("Expected size 2 after removal, got %d", csm.Size())
	}
	if csm.ContainsKey(2) {
		t.Error("Should not contain key 2 after removal")
	}
}

func TestSkipListMapClear(t *testing.T) {
	csm := NewConcurrentSkipListMap[string, float64]()
	csm.Put("pi", 3.14)
	csm.Put("e", 2.71)
	csm.Put("phi", 1.618)

	if csm.Size() != 3 {
		t.Errorf("Expected size 3 before clear, got %d", csm.Size())
	}

	csm.Clear()
	if csm.Size() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", csm.Size())
	}

	// 验证清空后可以继续使用
	csm.Put("new", 1.0)
	if csm.Size() != 1 {
		t.Errorf("Expected size 1 after adding to cleared map, got %d", csm.Size())
	}
}

// ============ 有序性测试 ============

func TestSkipListMapOrder(t *testing.T) {
	csm := NewConcurrentSkipListMap[int, string]()

	// 乱序插入
	keys := []int{5, 2, 8, 1, 9, 3, 7, 4, 6}
	for _, k := range keys {
		csm.Put(k, fmt.Sprintf("value%d", k))
	}

	// 验证Keys()返回有序结果
	resultKeys := csm.Keys()
	expectedKeys := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	if len(resultKeys) != len(expectedKeys) {
		t.Errorf("Expected %d keys, got %d", len(expectedKeys), len(resultKeys))
	}

	for i := range expectedKeys {
		if resultKeys[i] != expectedKeys[i] {
			t.Errorf("At index %d: expected key %d, got %d", i, expectedKeys[i], resultKeys[i])
		}
	}

	// 验证Values()按键的顺序返回
	resultValues := csm.Values()
	for i, k := range expectedKeys {
		expectedValue := fmt.Sprintf("value%d", k)
		if resultValues[i] != expectedValue {
			t.Errorf("At index %d: expected value %s, got %s", i, expectedValue, resultValues[i])
		}
	}
}

func TestSkipListMapOrderWithStrings(t *testing.T) {
	csm := NewConcurrentSkipListMap[string, int]()

	// 插入字符串键
	words := []string{"dog", "apple", "cat", "banana", "elephant"}
	for i, w := range words {
		csm.Put(w, i)
	}

	// 验证有序性
	keys := csm.Keys()
	expected := []string{"apple", "banana", "cat", "dog", "elephant"}

	for i := range expected {
		if keys[i] != expected[i] {
			t.Errorf("At index %d: expected key %s, got %s", i, expected[i], keys[i])
		}
	}
}

func TestSkipListMapRangeOrder(t *testing.T) {
	csm := NewConcurrentSkipListMap[int, string]()

	// 插入数据
	for i := 10; i >= 1; i-- {
		csm.Put(i, fmt.Sprintf("v%d", i))
	}

	// 使用Range验证顺序
	lastKey := 0
	csm.Range(func(key int, value string) bool {
		if key <= lastKey {
			t.Errorf("Keys not in ascending order: %d after %d", key, lastKey)
		}
		lastKey = key
		expectedValue := fmt.Sprintf("v%d", key)
		if value != expectedValue {
			t.Errorf("Expected value %s for key %d, got %s", expectedValue, key, value)
		}
		return true
	})
}

// ============ 高级功能测试 ============

func TestSkipListMapPutIfAbsent(t *testing.T) {
	csm := NewConcurrentSkipListMap[string, int]()

	// 第一次放入
	existing, loaded := csm.PutIfAbsent("a", 1)
	if loaded {
		t.Error("Expected key 'a' not to exist")
	}
	if existing != 1 {
		t.Errorf("Expected returned value 1, got %d", existing)
	}

	// 第二次尝试放入相同key
	existing, loaded = csm.PutIfAbsent("a", 2)
	if !loaded {
		t.Error("Expected key 'a' to exist")
	}
	if existing != 1 {
		t.Errorf("Expected existing value 1, got %d", existing)
	}
	if val := csm.Get("a"); val != 1 {
		t.Errorf("Expected value to remain 1, got %d", val)
	}

	// 验证大小没有变化
	if csm.Size() != 1 {
		t.Errorf("Expected size 1, got %d", csm.Size())
	}
}

func TestSkipListMapGetOrDefault(t *testing.T) {
	csm := NewConcurrentSkipListMap[string, string]()

	// 测试不存在的key
	defaultVal := csm.GetOrDefault("missing", "default")
	if defaultVal != "default" {
		t.Errorf("Expected 'default', got '%s'", defaultVal)
	}

	// 测试存在的key
	csm.Put("exists", "value")
	val := csm.GetOrDefault("exists", "default")
	if val != "value" {
		t.Errorf("Expected 'value', got '%s'", val)
	}

	// 测试零值
	csm.Put("zero", "")
	val = csm.GetOrDefault("zero", "default")
	if val != "" {
		t.Errorf("Expected empty string, got '%s'", val)
	}
}

func TestSkipListMapToMap(t *testing.T) {
	csm := NewConcurrentSkipListMap[int, string]()
	csm.Put(1, "one")
	csm.Put(2, "two")
	csm.Put(3, "three")

	m := csm.ToMap()
	if len(m) != 3 {
		t.Errorf("Expected map size 3, got %d", len(m))
	}
	if m[1] != "one" {
		t.Errorf("Expected 'one' for key 1, got '%s'", m[1])
	}
	if m[2] != "two" {
		t.Errorf("Expected 'two' for key 2, got '%s'", m[2])
	}
	if m[3] != "three" {
		t.Errorf("Expected 'three' for key 3, got '%s'", m[3])
	}
}

func TestSkipListMapRange(t *testing.T) {
	// 准备测试数据
	testData := map[int]string{
		1: "one",
		2: "two",
		3: "three",
		4: "four",
		5: "five",
	}
	csm := NewConcurrentSkipListMap(testData)

	// 测试完整遍历
	visited := make(map[int]string)
	csm.Range(func(key int, value string) bool {
		visited[key] = value
		return true
	})
	if len(visited) != len(testData) {
		t.Errorf("Expected %d items, got %d", len(testData), len(visited))
	}
	for k, v := range testData {
		if visited[k] != v {
			t.Errorf("Expected value %s for key %d, got %s", v, k, visited[k])
		}
	}

	// 测试提前终止
	count := 0
	csm.Range(func(key int, value string) bool {
		count++
		return count < 3 // 只遍历前3个元素
	})
	if count != 3 {
		t.Errorf("Expected to stop after 3 items, but processed %d", count)
	}

	// 验证遍历顺序
	var keys []int
	csm.Range(func(key int, value string) bool {
		keys = append(keys, key)
		return true
	})
	for i := 1; i < len(keys); i++ {
		if keys[i] <= keys[i-1] {
			t.Errorf("Keys not in ascending order: %d after %d", keys[i], keys[i-1])
		}
	}
}

func TestSkipListMapToString(t *testing.T) {
	csm := NewConcurrentSkipListMap[string, int]()
	csm.Put("a", 1)
	csm.Put("b", 2)
	csm.Put("c", 3)

	str := csm.ToString()
	var m map[string]int
	err := json.Unmarshal([]byte(str), &m)
	if err != nil {
		t.Errorf("Failed to unmarshal string: %v", err)
	}
	if len(m) != 3 {
		t.Errorf("Expected unmarshaled map size 3, got %d", len(m))
	}
	if m["a"] != 1 || m["b"] != 2 || m["c"] != 3 {
		t.Error("Unmarshaled values don't match")
	}
}

// ============ FirstKey/LastKey 测试 ============

func TestSkipListMapFirstKey(t *testing.T) {
	csm := NewConcurrentSkipListMap[int, string]()

	// 空map测试
	_, ok := csm.FirstKey()
	if ok {
		t.Error("Expected FirstKey to return false for empty map")
	}

	// 插入数据
	keys := []int{5, 2, 8, 1, 9}
	for _, k := range keys {
		csm.Put(k, fmt.Sprintf("v%d", k))
	}

	// 获取第一个键
	firstKey, ok := csm.FirstKey()
	if !ok {
		t.Error("Expected FirstKey to return true")
	}
	if firstKey != 1 {
		t.Errorf("Expected first key 1, got %d", firstKey)
	}

	// 删除第一个键后测试
	csm.Remove(1)
	firstKey, ok = csm.FirstKey()
	if !ok {
		t.Error("Expected FirstKey to return true after removing first")
	}
	if firstKey != 2 {
		t.Errorf("Expected first key 2 after removal, got %d", firstKey)
	}
}

func TestSkipListMapLastKey(t *testing.T) {
	csm := NewConcurrentSkipListMap[int, string]()

	// 空map测试
	_, ok := csm.LastKey()
	if ok {
		t.Error("Expected LastKey to return false for empty map")
	}

	// 插入数据
	keys := []int{5, 2, 8, 1, 9}
	for _, k := range keys {
		csm.Put(k, fmt.Sprintf("v%d", k))
	}

	// 获取最后一个键
	lastKey, ok := csm.LastKey()
	if !ok {
		t.Error("Expected LastKey to return true")
	}
	if lastKey != 9 {
		t.Errorf("Expected last key 9, got %d", lastKey)
	}

	// 删除最后一个键后测试
	csm.Remove(9)
	lastKey, ok = csm.LastKey()
	if !ok {
		t.Error("Expected LastKey to return true after removing last")
	}
	if lastKey != 8 {
		t.Errorf("Expected last key 8 after removal, got %d", lastKey)
	}
}

func TestSkipListMapFirstEntry(t *testing.T) {
	csm := NewConcurrentSkipListMap[string, int]()

	// 空map测试
	_, _, ok := csm.FirstEntry()
	if ok {
		t.Error("Expected FirstEntry to return false for empty map")
	}

	// 插入数据
	csm.Put("dog", 1)
	csm.Put("apple", 2)
	csm.Put("cat", 3)

	// 获取第一个条目
	key, value, ok := csm.FirstEntry()
	if !ok {
		t.Error("Expected FirstEntry to return true")
	}
	if key != "apple" {
		t.Errorf("Expected first key 'apple', got '%s'", key)
	}
	if value != 2 {
		t.Errorf("Expected first value 2, got %d", value)
	}
}

func TestSkipListMapLastEntry(t *testing.T) {
	csm := NewConcurrentSkipListMap[string, int]()

	// 空map测试
	_, _, ok := csm.LastEntry()
	if ok {
		t.Error("Expected LastEntry to return false for empty map")
	}

	// 插入数据
	csm.Put("dog", 1)
	csm.Put("apple", 2)
	csm.Put("cat", 3)

	// 获取最后一个条目
	key, value, ok := csm.LastEntry()
	if !ok {
		t.Error("Expected LastEntry to return true")
	}
	if key != "dog" {
		t.Errorf("Expected last key 'dog', got '%s'", key)
	}
	if value != 1 {
		t.Errorf("Expected last value 1, got %d", value)
	}
}

// ============ 并发安全测试 ============

func TestSkipListMapConcurrentAccess(t *testing.T) {
	csm := NewConcurrentSkipListMap[int, int]()
	const numRoutines = 100
	const numIterations = 1000

	var wg sync.WaitGroup
	wg.Add(numRoutines * 3)

	// 并发写入
	for i := 0; i < numRoutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				key := id*numIterations + j
				csm.Put(key, key)
			}
		}(i)
	}

	// 并发读取
	for i := 0; i < numRoutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				key := id*numIterations + j
				_ = csm.Get(key)
			}
		}(i)
	}

	// 并发删除
	for i := 0; i < numRoutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numIterations/2; j++ {
				key := id*numIterations + j
				csm.Remove(key)
			}
		}(i)
	}

	wg.Wait()

	// 验证最终大小
	expectedSize := numRoutines * numIterations / 2
	if csm.Size() != expectedSize {
		t.Logf("Warning: Expected size around %d, got %d (due to concurrent operations)", expectedSize, csm.Size())
	}
}

func TestSkipListMapConcurrentReadWrite(t *testing.T) {
	csm := NewConcurrentSkipListMap[int, string]()
	const numRoutines = 50
	const numOperations = 500

	var wg sync.WaitGroup
	wg.Add(numRoutines * 4)

	// 并发Put
	for i := 0; i < numRoutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				csm.Put(j, fmt.Sprintf("value-%d-%d", id, j))
			}
		}(i)
	}

	// 并发Get
	for i := 0; i < numRoutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				_ = csm.Get(j)
			}
		}()
	}

	// 并发ContainsKey
	for i := 0; i < numRoutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				_ = csm.ContainsKey(j)
			}
		}()
	}

	// 并发Range
	for i := 0; i < numRoutines; i++ {
		go func() {
			defer wg.Done()
			csm.Range(func(key int, value string) bool {
				return true
			})
		}()
	}

	wg.Wait()

	// 验证基本属性
	if csm.Size() > numOperations {
		t.Errorf("Size should not exceed %d, got %d", numOperations, csm.Size())
	}
}

func TestSkipListMapConcurrentPutIfAbsent(t *testing.T) {
	csm := NewConcurrentSkipListMap[int, int]()
	const numRoutines = 100
	const targetKey = 42

	var wg sync.WaitGroup
	wg.Add(numRoutines)

	// 多个goroutine同时尝试插入同一个key
	successCount := 0
	var mu sync.Mutex

	for i := 0; i < numRoutines; i++ {
		go func(id int) {
			defer wg.Done()
			_, loaded := csm.PutIfAbsent(targetKey, id)
			if !loaded {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	// 只有一个goroutine应该成功插入
	if successCount != 1 {
		t.Errorf("Expected exactly 1 successful PutIfAbsent, got %d", successCount)
	}

	// 验证key存在
	if !csm.ContainsKey(targetKey) {
		t.Error("Expected key to exist after concurrent PutIfAbsent")
	}

	if csm.Size() != 1 {
		t.Errorf("Expected size 1, got %d", csm.Size())
	}
}

// ============ 序列化测试 ============

func TestSkipListMapJSONMarshalUnmarshal(t *testing.T) {
	csm := NewConcurrentSkipListMap[string, int]()
	csm.Put("x", 10)
	csm.Put("y", 20)
	csm.Put("z", 30)

	// 测试序列化
	data, err := json.Marshal(csm)
	if err != nil {
		t.Fatal("JSON序列化失败:", err)
	}

	// 测试反序列化
	var newCsm ConcurrentSkipListMap[string, int]
	if err := json.Unmarshal(data, &newCsm); err != nil {
		t.Fatal("JSON反序列化失败:", err)
	}

	if newCsm.Size() != 3 {
		t.Errorf("反序列化后的map大小不正确, expected 3, got %d", newCsm.Size())
	}
	if newCsm.Get("x") != 10 {
		t.Error("反序列化后的值不正确")
	}
	if newCsm.Get("y") != 20 {
		t.Error("反序列化后的值不正确")
	}

	// 验证顺序保持
	keys := newCsm.Keys()
	expectedKeys := []string{"x", "y", "z"}
	for i, k := range expectedKeys {
		if keys[i] != k {
			t.Errorf("Key order not preserved after unmarshal, expected %s at index %d, got %s", k, i, keys[i])
		}
	}
}

func TestSkipListMapBSONMarshalUnmarshal(t *testing.T) {
	csm := NewConcurrentSkipListMap[string, float64]()
	csm.Put("pi", 3.14159)
	csm.Put("e", 2.71828)
	csm.Put("phi", 1.618)

	// 测试BSON序列化
	data, err := bson.Marshal(csm)
	if err != nil {
		t.Fatal("BSON序列化失败:", err)
	}

	// 测试BSON反序列化
	var newCsm ConcurrentSkipListMap[string, float64]
	if err := bson.Unmarshal(data, &newCsm); err != nil {
		t.Fatal("BSON反序列化失败:", err)
	}

	if newCsm.Size() != 3 {
		t.Errorf("反序列化后的map大小不正确, expected 3, got %d", newCsm.Size())
	}

	if newCsm.Get("pi") != 3.14159 {
		t.Error("反序列化后的值不正确")
	}
}

func TestSkipListMapJSONWithIntKeys(t *testing.T) {
	csm := NewConcurrentSkipListMap[int, string]()
	for i := 1; i <= 5; i++ {
		csm.Put(i, fmt.Sprintf("value%d", i))
	}

	// 序列化
	data, err := json.Marshal(csm)
	if err != nil {
		t.Fatal("JSON序列化失败:", err)
	}

	// 反序列化
	var newCsm ConcurrentSkipListMap[int, string]
	if err := json.Unmarshal(data, &newCsm); err != nil {
		t.Fatal("JSON反序列化失败:", err)
	}

	// 验证数据
	if newCsm.Size() != 5 {
		t.Errorf("Expected size 5, got %d", newCsm.Size())
	}

	for i := 1; i <= 5; i++ {
		expected := fmt.Sprintf("value%d", i)
		if newCsm.Get(i) != expected {
			t.Errorf("Expected %s for key %d, got %s", expected, i, newCsm.Get(i))
		}
	}
}

// ============ 边界情况测试 ============

func TestSkipListMapEmptyOperations(t *testing.T) {
	csm := NewConcurrentSkipListMap[string, string]()

	// 空map上的操作
	if csm.Get("any") != "" {
		t.Error("Get on empty map should return zero value")
	}

	if csm.Size() != 0 {
		t.Error("Empty map size should be 0")
	}

	if csm.ContainsKey("any") {
		t.Error("Empty map should not contain any key")
	}

	keys := csm.Keys()
	if len(keys) != 0 {
		t.Error("Empty map should return empty keys slice")
	}

	values := csm.Values()
	if len(values) != 0 {
		t.Error("Empty map should return empty values slice")
	}

	// Range on empty map
	count := 0
	csm.Range(func(key string, value string) bool {
		count++
		return true
	})
	if count != 0 {
		t.Error("Range on empty map should not iterate")
	}

	// FirstKey/LastKey on empty map
	_, ok := csm.FirstKey()
	if ok {
		t.Error("FirstKey should return false on empty map")
	}

	_, ok = csm.LastKey()
	if ok {
		t.Error("LastKey should return false on empty map")
	}
}

func TestSkipListMapSingleElement(t *testing.T) {
	csm := NewConcurrentSkipListMap[int, string]()
	csm.Put(42, "answer")

	// 验证单元素操作
	if csm.Size() != 1 {
		t.Errorf("Expected size 1, got %d", csm.Size())
	}

	firstKey, ok := csm.FirstKey()
	if !ok || firstKey != 42 {
		t.Error("FirstKey failed on single-element map")
	}

	lastKey, ok := csm.LastKey()
	if !ok || lastKey != 42 {
		t.Error("LastKey failed on single-element map")
	}

	// FirstKey和LastKey应该相同
	if firstKey != lastKey {
		t.Error("FirstKey and LastKey should be equal for single-element map")
	}

	// 删除唯一元素
	csm.Remove(42)
	if csm.Size() != 0 {
		t.Error("Size should be 0 after removing only element")
	}
}

func TestSkipListMapLargeDataSet(t *testing.T) {
	csm := NewConcurrentSkipListMap[int, int]()
	const size = 10000

	// 插入大量数据
	for i := 0; i < size; i++ {
		csm.Put(i, i*2)
	}

	if csm.Size() != size {
		t.Errorf("Expected size %d, got %d", size, csm.Size())
	}

	// 验证数据正确性
	for i := 0; i < size; i++ {
		if csm.Get(i) != i*2 {
			t.Errorf("Expected value %d for key %d, got %d", i*2, i, csm.Get(i))
		}
	}

	// 验证有序性
	keys := csm.Keys()
	for i := 1; i < len(keys); i++ {
		if keys[i] <= keys[i-1] {
			t.Error("Keys not in ascending order")
			break
		}
	}

	// 删除一半元素
	for i := 0; i < size; i += 2 {
		csm.Remove(i)
	}

	if csm.Size() != size/2 {
		t.Errorf("Expected size %d after removing half, got %d", size/2, csm.Size())
	}
}

func TestSkipListMapDuplicatePuts(t *testing.T) {
	csm := NewConcurrentSkipListMap[string, int]()

	// 多次Put同一个key
	for i := 0; i < 100; i++ {
		csm.Put("key", i)
	}

	// Size应该仍然是1
	if csm.Size() != 1 {
		t.Errorf("Expected size 1 after duplicate puts, got %d", csm.Size())
	}

	// 值应该是最后一次Put的值
	if csm.Get("key") != 99 {
		t.Errorf("Expected value 99, got %d", csm.Get("key"))
	}
}

func TestSkipListMapZeroValues(t *testing.T) {
	csm := NewConcurrentSkipListMap[int, int]()

	// 插入零值
	csm.Put(1, 0)
	csm.Put(2, 0)
	csm.Put(3, 0)

	if csm.Size() != 3 {
		t.Errorf("Expected size 3, got %d", csm.Size())
	}

	// 零值应该能正确存储和获取
	if csm.Get(1) != 0 {
		t.Error("Zero value not stored correctly")
	}

	if !csm.ContainsKey(1) {
		t.Error("Should contain key with zero value")
	}
}

func TestSkipListMapUpdateValues(t *testing.T) {
	csm := NewConcurrentSkipListMap[string, int]()

	// 初始插入
	keys := []string{"a", "b", "c", "d", "e"}
	for i, k := range keys {
		csm.Put(k, i)
	}

	// 更新所有值
	for i, k := range keys {
		csm.Put(k, i*10)
	}

	// 验证更新
	for i, k := range keys {
		expected := i * 10
		if csm.Get(k) != expected {
			t.Errorf("Expected value %d for key %s, got %d", expected, k, csm.Get(k))
		}
	}

	// Size不应该改变
	if csm.Size() != len(keys) {
		t.Errorf("Expected size %d, got %d", len(keys), csm.Size())
	}
}

// ============ 性能基准测试 ============

func BenchmarkSkipListMapPut(b *testing.B) {
	csm := NewConcurrentSkipListMap[int, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		csm.Put(i, i)
	}
}

func BenchmarkSkipListMapGet(b *testing.B) {
	csm := NewConcurrentSkipListMap[int, int]()
	for i := 0; i < 10000; i++ {
		csm.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		csm.Get(i % 10000)
	}
}

func BenchmarkSkipListMapRemove(b *testing.B) {
	csm := NewConcurrentSkipListMap[int, int]()
	for i := 0; i < b.N; i++ {
		csm.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		csm.Remove(i)
	}
}

func BenchmarkSkipListMapRange(b *testing.B) {
	csm := NewConcurrentSkipListMap[int, int]()
	for i := 0; i < 1000; i++ {
		csm.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		csm.Range(func(key int, value int) bool {
			return true
		})
	}
}

func BenchmarkSkipListMapConcurrentPut(b *testing.B) {
	csm := NewConcurrentSkipListMap[int, int]()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			csm.Put(i, i)
			i++
		}
	})
}

func BenchmarkSkipListMapConcurrentGet(b *testing.B) {
	csm := NewConcurrentSkipListMap[int, int]()
	for i := 0; i < 10000; i++ {
		csm.Put(i, i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			csm.Get(i % 10000)
			i++
		}
	})
}

func BenchmarkSkipListMapConcurrentMixed(b *testing.B) {
	csm := NewConcurrentSkipListMap[int, int]()
	for i := 0; i < 1000; i++ {
		csm.Put(i, i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			op := i % 3
			key := i % 1000
			switch op {
			case 0:
				csm.Put(key, i)
			case 1:
				csm.Get(key)
			case 2:
				csm.ContainsKey(key)
			}
			i++
		}
	})
}

// ============ 与HashMap性能对比测试 ============

func TestComparePerformance(t *testing.T) {
	const size = 100000
	timer := dateUtil.NewTimer()

	// 测试ConcurrentHashMap
	chm := NewConcurrentHashMap[int, int]()
	for i := 0; i < size; i++ {
		chm.Put(i, i)
	}
	chmTime := timer.Interval()

	// 测试ConcurrentSkipListMap
	csm := NewConcurrentSkipListMap[int, int]()
	for i := 0; i < size; i++ {
		csm.Put(i, i)
	}
	csmTime := timer.Interval()

	t.Logf("ConcurrentHashMap Put %d items: %d ms", size, chmTime)
	t.Logf("ConcurrentSkipListMap Put %d items: %d ms", size, csmTime)

	// 测试随机读取性能
	timer = dateUtil.NewTimer()
	for i := 0; i < size; i++ {
		_ = chm.Get(rand.Intn(size))
	}
	chmGetTime := timer.Interval()

	for i := 0; i < size; i++ {
		_ = csm.Get(rand.Intn(size))
	}
	csmGetTime := timer.Interval()

	t.Logf("ConcurrentHashMap Random Get: %d ms", chmGetTime)
	t.Logf("ConcurrentSkipListMap Random Get: %d ms", csmGetTime)

	// 测试有序遍历(SkipListMap的优势)
	timer = dateUtil.NewTimer()
	keys := csm.Keys()
	csmKeysTime := timer.Interval()

	// 验证有序性
	isSorted := sort.IntsAreSorted(keys)
	t.Logf("ConcurrentSkipListMap Keys() (sorted): %d ms, is sorted: %v", csmKeysTime, isSorted)

	// HashMap的Keys()是无序的
	timer = dateUtil.NewTimer()
	_ = chm.Keys()
	chmKeysTime := timer.Interval()
	t.Logf("ConcurrentHashMap Keys() (unsorted): %d ms", chmKeysTime)
}

// ============ 特殊场景测试 ============

func TestSkipListMapWithNegativeKeys(t *testing.T) {
	csm := NewConcurrentSkipListMap[int, string]()

	// 插入负数键
	keys := []int{-5, -2, -8, -1, 0, 3, 7, 1, -9}
	for _, k := range keys {
		csm.Put(k, fmt.Sprintf("v%d", k))
	}

	// 验证有序性
	resultKeys := csm.Keys()
	for i := 1; i < len(resultKeys); i++ {
		if resultKeys[i] <= resultKeys[i-1] {
			t.Errorf("Keys not in ascending order: %d after %d", resultKeys[i], resultKeys[i-1])
		}
	}

	// 验证FirstKey和LastKey
	firstKey, _ := csm.FirstKey()
	lastKey, _ := csm.LastKey()

	if firstKey != -9 {
		t.Errorf("Expected first key -9, got %d", firstKey)
	}
	if lastKey != 7 {
		t.Errorf("Expected last key 7, got %d", lastKey)
	}
}

func TestSkipListMapRangeModification(t *testing.T) {
	csm := NewConcurrentSkipListMap[int, string]()
	for i := 1; i <= 10; i++ {
		csm.Put(i, fmt.Sprintf("v%d", i))
	}

	// 在Range中不应该修改原map(因为Range复制了数据)
	// 但我们可以在Range之外修改
	count := 0
	csm.Range(func(key int, value string) bool {
		count++
		return true
	})

	if count != 10 {
		t.Errorf("Expected to iterate 10 items, got %d", count)
	}

	// 验证在Range过程中添加的元素不会影响当前迭代
	iterations := 0
	var wg sync.WaitGroup
	wg.Add(1)
	csm.Range(func(key int, value string) bool {
		if iterations == 0 {
			// 在另一个goroutine中添加元素
			go func() {
				defer wg.Done()
				csm.Put(100, "v100")
			}()
		}
		iterations++
		return true
	})

	// 等待goroutine完成
	wg.Wait()

	// 验证元素已添加
	if !csm.ContainsKey(100) {
		t.Error("New element should be added")
	}
}

func TestSkipListMapStressTest(t *testing.T) {
	csm := NewConcurrentSkipListMap[int, int]()
	const numRoutines = 20
	const numOps = 1000

	var wg sync.WaitGroup
	wg.Add(numRoutines * 3)

	// 混合操作压力测试
	for i := 0; i < numRoutines; i++ {
		// Put
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				key := rand.Intn(numOps)
				csm.Put(key, id*numOps+j)
			}
		}(i)

		// Get
		go func() {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				key := rand.Intn(numOps)
				_ = csm.Get(key)
			}
		}()

		// Remove
		go func() {
			defer wg.Done()
			for j := 0; j < numOps/2; j++ {
				key := rand.Intn(numOps)
				csm.Remove(key)
			}
		}()
	}

	wg.Wait()

	// 验证map仍然可用
	size := csm.Size()
	t.Logf("After stress test, map size: %d", size)

	// 验证有序性
	keys := csm.Keys()
	for i := 1; i < len(keys); i++ {
		if keys[i] <= keys[i-1] {
			t.Error("Keys not in ascending order after stress test")
			break
		}
	}
}
