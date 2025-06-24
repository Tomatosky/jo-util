package mapUtil

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"sync"
	"testing"
)

func TestNewConcurrentHashMap(t *testing.T) {
	// 测试空初始化
	emptyMap := NewConcurrentHashMap[string, int]()
	if emptyMap.Size() != 0 {
		t.Errorf("Expected empty map size 0, got %d", emptyMap.Size())
	}

	// 测试带初始值的初始化
	initData := map[string]int{"a": 1, "b": 2}
	initMap := NewConcurrentHashMap(initData)
	if initMap.Size() != 2 {
		t.Errorf("Expected initialized map size 2, got %d", initMap.Size())
	}
	if val := initMap.Get("a"); val != 1 {
		t.Errorf("Expected value 1 for key 'a', got %d", val)
	}
}

func TestGetPutRemove(t *testing.T) {
	cm := NewConcurrentHashMap[string, string]()

	// 测试Put和Get
	cm.Put("key1", "value1")
	if val := cm.Get("key1"); val != "value1" {
		t.Errorf("Expected 'value1', got '%s'", val)
	}

	// 测试Remove
	cm.Remove("key1")
	if val := cm.Get("key1"); val != "" {
		t.Errorf("Expected empty value after remove, got '%s'", val)
	}
}

func TestSizeAndContainsKey(t *testing.T) {
	cm := NewConcurrentHashMap[int, bool]()

	// 初始大小
	if cm.Size() != 0 {
		t.Errorf("Expected size 0, got %d", cm.Size())
	}

	// 添加元素后测试
	cm.Put(1, true)
	if cm.Size() != 1 {
		t.Errorf("Expected size 1, got %d", cm.Size())
	}
	if !cm.ContainsKey(1) {
		t.Error("Expected to contain key 1")
	}
	if cm.ContainsKey(2) {
		t.Error("Should not contain key 2")
	}
}

func TestClear(t *testing.T) {
	cm := NewConcurrentHashMap[string, float64]()
	cm.Put("pi", 3.14)
	cm.Put("e", 2.71)

	cm.Clear()
	if cm.Size() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", cm.Size())
	}
}

func TestKeysAndValues(t *testing.T) {
	cm := NewConcurrentHashMap[string, int]()
	cm.Put("a", 1)
	cm.Put("b", 2)
	cm.Put("c", 3)

	keys := cm.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	values := cm.Values()
	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}
}

func TestPutIfAbsent2(t *testing.T) {
	cm := NewConcurrentHashMap[string, int]()

	// 第一次放入
	existing, loaded := cm.PutIfAbsent("a", 1)
	if loaded {
		t.Error("Expected key 'a' not to exist")
	}
	if existing != 0 {
		t.Errorf("Expected existing value 0, got %d", existing)
	}

	// 第二次尝试放入相同key
	existing, loaded = cm.PutIfAbsent("a", 2)
	if !loaded {
		t.Error("Expected key 'a' to exist")
	}
	if existing != 1 {
		t.Errorf("Expected existing value 1, got %d", existing)
	}
	if val := cm.Get("a"); val != 1 {
		t.Errorf("Expected value to remain 1, got %d", val)
	}
}

func TestGetOrDefault2(t *testing.T) {
	cm := NewConcurrentHashMap[string, string]()

	// 测试不存在的key
	defaultVal := cm.GetOrDefault("missing", "default")
	if defaultVal != "default" {
		t.Errorf("Expected 'default', got '%s'", defaultVal)
	}

	// 测试存在的key
	cm.Put("exists", "value")
	val := cm.GetOrDefault("exists", "default")
	if val != "value" {
		t.Errorf("Expected 'value', got '%s'", val)
	}
}

func TestToMap2(t *testing.T) {
	cm := NewConcurrentHashMap[int, string]()
	cm.Put(1, "one")
	cm.Put(2, "two")

	m := cm.ToMap()
	if len(m) != 2 {
		t.Errorf("Expected map size 2, got %d", len(m))
	}
	if m[1] != "one" {
		t.Errorf("Expected 'one' for key 1, got '%s'", m[1])
	}
}

func TestToString2(t *testing.T) {
	cm := NewConcurrentHashMap[string, int]()
	cm.Put("a", 1)
	cm.Put("b", 2)

	str := cm.ToString()
	var m map[string]int
	err := json.Unmarshal([]byte(str), &m)
	if err != nil {
		t.Errorf("Failed to unmarshal string: %v", err)
	}
	if len(m) != 2 {
		t.Errorf("Expected unmarshaled map size 2, got %d", len(m))
	}
	if m["a"] != 1 {
		t.Errorf("Expected 'a' to be 1, got %d", m["a"])
	}
}

func TestConcurrentHashMap_Range(t *testing.T) {
	// 准备测试数据
	testData := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	cm := NewConcurrentHashMap(testData)
	// 测试完整遍历
	visited := make(map[string]int)
	cm.Range(func(key string, value int) bool {
		visited[key] = value
		return true
	})
	if len(visited) != len(testData) {
		t.Errorf("Expected %d items, got %d", len(testData), len(visited))
	}
	for k, v := range testData {
		if visited[k] != v {
			t.Errorf("Expected value %d for key %s, got %d", v, k, visited[k])
		}
	}
	// 测试提前终止
	count := 0
	cm.Range(func(key string, value int) bool {
		count++
		return count < 2 // 只遍历前两个元素
	})
	if count != 2 {
		t.Errorf("Expected to stop after 2 items, but processed %d", count)
	}
}

func TestConcurrentAccess(t *testing.T) {
	cm := NewConcurrentHashMap[int, int]()
	const numRoutines = 100
	const numIterations = 1000

	var wg sync.WaitGroup
	wg.Add(numRoutines * 2)

	// 并发写入
	for i := 0; i < numRoutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				cm.Put(j, j)
			}
		}()
	}

	// 并发读取
	for i := 0; i < numRoutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				_ = cm.Get(j)
			}
		}()
	}

	wg.Wait()

	// 验证最终结果
	for j := 0; j < numIterations; j++ {
		if val := cm.Get(j); val != j {
			t.Errorf("Expected %d, got %d", j, val)
		}
	}
}

func TestJSONMarshalUnmarshal(t *testing.T) {
	cm := NewConcurrentHashMap[string, int]()
	cm.Put("x", 10)
	cm.Put("y", 20)
	// 测试序列化
	data, err := json.Marshal(cm)
	if err != nil {
		t.Error("JSON序列化失败:", err)
	}
	// 测试反序列化
	var newCm ConcurrentHashMap[string, int]
	if err := json.Unmarshal(data, &newCm); err != nil {
		t.Error("JSON反序列化失败:", err)
	}
	if newCm.Size() != 2 {
		t.Error("反序列化后的map大小不正确")
	}
	if newCm.Get("x") != 10 {
		t.Error("反序列化后的值不正确")
	}
}

func TestBSONMarshalUnmarshal(t *testing.T) {
	cm := NewConcurrentHashMap[string, float64]()
	cm.Put("pi", 3.14)
	cm.Put("e", 2.718)
	// 测试BSON序列化
	data, err := bson.Marshal(cm)
	if err != nil {
		t.Error("BSON序列化失败:", err)
	}
	// 测试BSON反序列化
	var newCm ConcurrentHashMap[string, float64]
	if err := bson.Unmarshal(data, &newCm); err != nil {
		t.Error("BSON反序列化失败:", err)
	}
	if newCm.Size() != 2 {
		t.Error("反序列化后的map大小不正确")
	}
	if newCm.Get("pi") != 3.14 {
		t.Error("反序列化后的值不正确")
	}
}
