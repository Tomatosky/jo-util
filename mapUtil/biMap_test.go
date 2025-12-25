package mapUtil

import (
	"encoding/json"
	"sync"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

// TestNewBiMap_Empty 测试创建空的BiMap
func TestNewBiMap_Empty(t *testing.T) {
	bm := NewBiMap[string, int]()
	if bm == nil {
		t.Fatal("NewBiMap returned nil")
	}
	if bm.Size() != 0 {
		t.Errorf("Expected size 0, got %d", bm.Size())
	}
}

// TestNewBiMap_WithInitMap 测试用初始map创建BiMap
func TestNewBiMap_WithInitMap(t *testing.T) {
	initMap := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	bm := NewBiMap(initMap)

	if bm.Size() != 3 {
		t.Errorf("Expected size 3, got %d", bm.Size())
	}

	if bm.Get("one") != 1 {
		t.Errorf("Expected value 1, got %d", bm.Get("one"))
	}

	if bm.GetKey(2) != "two" {
		t.Errorf("Expected key 'two', got '%s'", bm.GetKey(2))
	}
}

// TestNewBiMap_WithDuplicateValues 测试初始map中有重复value的情况
func TestNewBiMap_WithDuplicateValues(t *testing.T) {
	// 当有重复value时，后面的会覆盖前面的
	initMap := map[string]int{
		"a": 1,
		"b": 1, // 重复的value
	}
	bm := NewBiMap(initMap)

	// 由于map遍历顺序不确定，只检查最终状态
	if bm.Size() != 1 {
		t.Errorf("Expected size 1 (duplicate values should result in one mapping), got %d", bm.Size())
	}

	if !bm.ContainsValue(1) {
		t.Error("Expected value 1 to be present")
	}
}

// TestPut_Basic 测试基本的Put操作
func TestPut_Basic(t *testing.T) {
	bm := NewBiMap[string, int]()

	bm.Put("key1", 100)
	if bm.Get("key1") != 100 {
		t.Errorf("Expected value 100, got %d", bm.Get("key1"))
	}
	if bm.GetKey(100) != "key1" {
		t.Errorf("Expected key 'key1', got '%s'", bm.GetKey(100))
	}
}

// TestPut_OverwriteKey 测试覆盖已存在的key
func TestPut_OverwriteKey(t *testing.T) {
	bm := NewBiMap[string, int]()

	bm.Put("key1", 100)
	bm.Put("key1", 200) // 覆盖

	if bm.Get("key1") != 200 {
		t.Errorf("Expected value 200, got %d", bm.Get("key1"))
	}
	if bm.Size() != 1 {
		t.Errorf("Expected size 1, got %d", bm.Size())
	}
	// 旧的value 100应该不存在了
	if bm.ContainsValue(100) {
		t.Error("Old value 100 should not exist")
	}
}

// TestPut_DuplicateValue 测试Put已存在的value
func TestPut_DuplicateValue(t *testing.T) {
	bm := NewBiMap[string, int]()

	bm.Put("key1", 100)
	bm.Put("key2", 100) // 相同的value

	// key1应该被移除，因为value 100现在映射到key2
	if bm.ContainsKey("key1") {
		t.Error("key1 should have been removed")
	}
	if bm.GetKey(100) != "key2" {
		t.Errorf("Expected key2 for value 100, got %s", bm.GetKey(100))
	}
	if bm.Size() != 1 {
		t.Errorf("Expected size 1, got %d", bm.Size())
	}
}

// TestGet 测试Get操作
func TestGet(t *testing.T) {
	bm := NewBiMap[string, int]()
	bm.Put("key1", 100)

	value := bm.Get("key1")
	if value != 100 {
		t.Errorf("Expected 100, got %d", value)
	}

	// 不存在的key应该返回零值
	value = bm.Get("nonexistent")
	if value != 0 {
		t.Errorf("Expected zero value for nonexistent key, got %d", value)
	}
}

// TestGetKey 测试GetKey操作
func TestGetKey(t *testing.T) {
	bm := NewBiMap[string, int]()
	bm.Put("key1", 100)

	key := bm.GetKey(100)
	if key != "key1" {
		t.Errorf("Expected 'key1', got '%s'", key)
	}

	// 不存在的value应该返回零值
	key = bm.GetKey(999)
	if key != "" {
		t.Errorf("Expected empty string for nonexistent value, got '%s'", key)
	}
}

// TestBiMap_Remove 测试Remove操作
func TestBiMap_Remove(t *testing.T) {
	bm := NewBiMap[string, int]()
	bm.Put("key1", 100)
	bm.Put("key2", 200)

	bm.Remove("key1")

	if bm.ContainsKey("key1") {
		t.Error("key1 should have been removed")
	}
	if bm.ContainsValue(100) {
		t.Error("value 100 should have been removed")
	}
	if bm.Size() != 1 {
		t.Errorf("Expected size 1, got %d", bm.Size())
	}
}

// TestBiMap_Remove_Nonexistent 测试移除不存在的key
func TestBiMap_Remove_Nonexistent(t *testing.T) {
	bm := NewBiMap[string, int]()
	bm.Put("key1", 100)

	// 移除不存在的key不应该报错
	bm.Remove("nonexistent")

	if bm.Size() != 1 {
		t.Errorf("Expected size 1, got %d", bm.Size())
	}
}

// TestRemoveValue 测试RemoveValue操作
func TestRemoveValue(t *testing.T) {
	bm := NewBiMap[string, int]()
	bm.Put("key1", 100)
	bm.Put("key2", 200)

	bm.RemoveValue(100)

	if bm.ContainsValue(100) {
		t.Error("value 100 should have been removed")
	}
	if bm.ContainsKey("key1") {
		t.Error("key1 should have been removed")
	}
	if bm.Size() != 1 {
		t.Errorf("Expected size 1, got %d", bm.Size())
	}
}

// TestRemoveValue_Nonexistent 测试移除不存在的value
func TestRemoveValue_Nonexistent(t *testing.T) {
	bm := NewBiMap[string, int]()
	bm.Put("key1", 100)

	// 移除不存在的value不应该报错
	bm.RemoveValue(999)

	if bm.Size() != 1 {
		t.Errorf("Expected size 1, got %d", bm.Size())
	}
}

// TestBiMap_ContainsKey 测试ContainsKey操作
func TestBiMap_ContainsKey(t *testing.T) {
	bm := NewBiMap[string, int]()
	bm.Put("key1", 100)

	if !bm.ContainsKey("key1") {
		t.Error("Expected key1 to exist")
	}
	if bm.ContainsKey("key2") {
		t.Error("Expected key2 to not exist")
	}
}

// TestContainsValue 测试ContainsValue操作
func TestContainsValue(t *testing.T) {
	bm := NewBiMap[string, int]()
	bm.Put("key1", 100)

	if !bm.ContainsValue(100) {
		t.Error("Expected value 100 to exist")
	}
	if bm.ContainsValue(200) {
		t.Error("Expected value 200 to not exist")
	}
}

// TestBiMap_Clear 测试Clear操作
func TestBiMap_Clear(t *testing.T) {
	bm := NewBiMap[string, int]()
	bm.Put("key1", 100)
	bm.Put("key2", 200)
	bm.Put("key3", 300)

	bm.Clear()

	if bm.Size() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", bm.Size())
	}
	if bm.ContainsKey("key1") {
		t.Error("Expected no keys after clear")
	}
	if bm.ContainsValue(100) {
		t.Error("Expected no values after clear")
	}
}

// TestBiMap_Keys 测试Keys操作
func TestBiMap_Keys(t *testing.T) {
	bm := NewBiMap[string, int]()
	bm.Put("key1", 100)
	bm.Put("key2", 200)
	bm.Put("key3", 300)

	keys := bm.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	// 检查所有key都存在
	keyMap := make(map[string]bool)
	for _, k := range keys {
		keyMap[k] = true
	}
	if !keyMap["key1"] || !keyMap["key2"] || !keyMap["key3"] {
		t.Error("Expected keys key1, key2, key3")
	}
}

// TestBiMap_Values 测试Values操作
func TestBiMap_Values(t *testing.T) {
	bm := NewBiMap[string, int]()
	bm.Put("key1", 100)
	bm.Put("key2", 200)
	bm.Put("key3", 300)

	values := bm.Values()
	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}

	// 检查所有value都存在
	valueMap := make(map[int]bool)
	for _, v := range values {
		valueMap[v] = true
	}
	if !valueMap[100] || !valueMap[200] || !valueMap[300] {
		t.Error("Expected values 100, 200, 300")
	}
}

// TestPutIfAbsent_New 测试PutIfAbsent添加新key
func TestPutIfAbsent_New(t *testing.T) {
	bm := NewBiMap[string, int]()

	existing, loaded := bm.PutIfAbsent("key1", 100)

	if loaded {
		t.Error("Expected loaded to be false for new key")
	}
	if existing != 100 {
		t.Errorf("Expected existing value to be 100, got %d", existing)
	}
	if bm.Get("key1") != 100 {
		t.Error("Expected key1 to be added")
	}
}

// TestPutIfAbsent_Existing 测试PutIfAbsent对已存在的key
func TestPutIfAbsent_Existing(t *testing.T) {
	bm := NewBiMap[string, int]()
	bm.Put("key1", 100)

	existing, loaded := bm.PutIfAbsent("key1", 200)

	if !loaded {
		t.Error("Expected loaded to be true for existing key")
	}
	if existing != 100 {
		t.Errorf("Expected existing value to be 100, got %d", existing)
	}
	if bm.Get("key1") != 100 {
		t.Error("Expected key1 value to remain 100")
	}
}

// TestPutIfAbsent_DuplicateValue 测试PutIfAbsent时value已存在
func TestPutIfAbsent_DuplicateValue(t *testing.T) {
	bm := NewBiMap[string, int]()
	bm.Put("key1", 100)

	existing, loaded := bm.PutIfAbsent("key2", 100)

	if loaded {
		t.Error("Expected loaded to be false")
	}
	if existing != 100 {
		t.Errorf("Expected existing value to be 100, got %d", existing)
	}
	// key1应该被移除，因为value 100现在映射到key2
	if bm.ContainsKey("key1") {
		t.Error("key1 should have been removed")
	}
	if bm.GetKey(100) != "key2" {
		t.Error("value 100 should map to key2")
	}
}

// TestGetOrDefault_Exists 测试GetOrDefault对存在的key
func TestGetOrDefault_Exists(t *testing.T) {
	bm := NewBiMap[string, int]()
	bm.Put("key1", 100)

	value := bm.GetOrDefault("key1", 999)
	if value != 100 {
		t.Errorf("Expected 100, got %d", value)
	}
}

// TestGetOrDefault_NotExists 测试GetOrDefault对不存在的key
func TestGetOrDefault_NotExists(t *testing.T) {
	bm := NewBiMap[string, int]()

	value := bm.GetOrDefault("nonexistent", 999)
	if value != 999 {
		t.Errorf("Expected default value 999, got %d", value)
	}
}

// TestBiMap_ToMap 测试ToMap操作
func TestBiMap_ToMap(t *testing.T) {
	bm := NewBiMap[string, int]()
	bm.Put("key1", 100)
	bm.Put("key2", 200)

	m := bm.ToMap()

	if len(m) != 2 {
		t.Errorf("Expected map size 2, got %d", len(m))
	}
	if m["key1"] != 100 {
		t.Errorf("Expected m[key1]=100, got %d", m["key1"])
	}
	if m["key2"] != 200 {
		t.Errorf("Expected m[key2]=200, got %d", m["key2"])
	}

	// 修改返回的map不应该影响BiMap
	m["key1"] = 999
	if bm.Get("key1") != 100 {
		t.Error("ToMap should return a copy, not affect original BiMap")
	}
}

// TestBiMap_Range 测试Range遍历
func TestBiMap_Range(t *testing.T) {
	bm := NewBiMap[string, int]()
	bm.Put("key1", 100)
	bm.Put("key2", 200)
	bm.Put("key3", 300)

	count := 0
	bm.Range(func(key string, value int) bool {
		count++
		if !bm.ContainsKey(key) {
			t.Errorf("Unexpected key in range: %s", key)
		}
		if bm.Get(key) != value {
			t.Errorf("Value mismatch for key %s", key)
		}
		return true
	})

	if count != 3 {
		t.Errorf("Expected to iterate 3 times, got %d", count)
	}
}

// TestRange_EarlyTermination 测试Range提前终止
func TestRange_EarlyTermination(t *testing.T) {
	bm := NewBiMap[string, int]()
	bm.Put("key1", 100)
	bm.Put("key2", 200)
	bm.Put("key3", 300)

	count := 0
	bm.Range(func(key string, value int) bool {
		count++
		return count < 2 // 只遍历2次
	})

	if count != 2 {
		t.Errorf("Expected to iterate 2 times, got %d", count)
	}
}

// TestBiMap_ToString 测试ToString序列化
func TestBiMap_ToString(t *testing.T) {
	bm := NewBiMap[string, int]()
	bm.Put("key1", 100)
	bm.Put("key2", 200)

	jsonStr := bm.ToString()
	if jsonStr == "" {
		t.Error("ToString returned empty string")
	}

	// 解析JSON验证
	var m map[string]int
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		t.Errorf("Failed to parse JSON: %v", err)
	}
	if len(m) != 2 {
		t.Errorf("Expected 2 entries in JSON, got %d", len(m))
	}
}

// TestMarshalJSON 测试JSON序列化
func TestMarshalJSON(t *testing.T) {
	bm := NewBiMap[string, int]()
	bm.Put("key1", 100)
	bm.Put("key2", 200)

	data, err := json.Marshal(bm)
	if err != nil {
		t.Errorf("Failed to marshal: %v", err)
	}

	var m map[string]int
	err = json.Unmarshal(data, &m)
	if err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}
	if m["key1"] != 100 || m["key2"] != 200 {
		t.Error("Unmarshaled data doesn't match")
	}
}

// TestUnmarshalJSON 测试JSON反序列化
func TestUnmarshalJSON(t *testing.T) {
	jsonData := `{"key1":100,"key2":200}`

	bm := &BiMap[string, int]{}
	err := json.Unmarshal([]byte(jsonData), bm)
	if err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if bm.Size() != 2 {
		t.Errorf("Expected size 2, got %d", bm.Size())
	}
	if bm.Get("key1") != 100 {
		t.Error("key1 not correctly unmarshaled")
	}
	if bm.GetKey(200) != "key2" {
		t.Error("Inverse mapping not correctly built")
	}
}

// TestUnmarshalJSON_DuplicateValues 测试JSON反序列化时有重复value
func TestUnmarshalJSON_DuplicateValues(t *testing.T) {
	// 在JSON中模拟重复value的情况是困难的，因为map本身不允许
	// 但我们可以测试反序列化后再添加重复value
	jsonData := `{"key1":100,"key2":200}`

	bm := &BiMap[string, int]{}
	err := json.Unmarshal([]byte(jsonData), bm)
	if err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	// 验证双向映射正确建立
	if bm.GetKey(100) != "key1" {
		t.Error("Inverse mapping for 100 not correct")
	}
}

// TestMarshalBSON 测试BSON序列化
func TestMarshalBSON(t *testing.T) {
	bm := NewBiMap[string, int]()
	bm.Put("key1", 100)
	bm.Put("key2", 200)

	data, err := bson.Marshal(bm)
	if err != nil {
		t.Errorf("Failed to marshal BSON: %v", err)
	}

	var m map[string]int
	err = bson.Unmarshal(data, &m)
	if err != nil {
		t.Errorf("Failed to unmarshal BSON: %v", err)
	}
	if m["key1"] != 100 || m["key2"] != 200 {
		t.Error("BSON unmarshaled data doesn't match")
	}
}

// TestUnmarshalBSON 测试BSON反序列化
func TestUnmarshalBSON(t *testing.T) {
	originalMap := map[string]int{
		"key1": 100,
		"key2": 200,
	}
	bsonData, err := bson.Marshal(originalMap)
	if err != nil {
		t.Fatalf("Failed to create BSON data: %v", err)
	}

	bm := &BiMap[string, int]{}
	err = bson.Unmarshal(bsonData, bm)
	if err != nil {
		t.Errorf("Failed to unmarshal BSON: %v", err)
	}

	if bm.Size() != 2 {
		t.Errorf("Expected size 2, got %d", bm.Size())
	}
	if bm.Get("key1") != 100 {
		t.Error("key1 not correctly unmarshaled from BSON")
	}
	if bm.GetKey(200) != "key2" {
		t.Error("Inverse mapping not correctly built from BSON")
	}
}

// TestConcurrency 测试并发安全性
func TestConcurrency(t *testing.T) {
	bm := NewBiMap[int, int]()
	var wg sync.WaitGroup

	// 并发写入
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			bm.Put(n, n*10)
		}(i)
	}

	// 并发读取
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			_ = bm.Get(n)
			_ = bm.GetKey(n * 10)
		}(i)
	}

	// 并发删除
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			bm.Remove(n)
		}(i)
	}

	wg.Wait()

	// 验证没有panic，并且数据一致性
	size := bm.Size()
	if size < 0 || size > 100 {
		t.Errorf("Unexpected size after concurrent operations: %d", size)
	}

	// 验证双向映射的一致性
	for _, key := range bm.Keys() {
		value := bm.Get(key)
		retrievedKey := bm.GetKey(value)
		if retrievedKey != key {
			t.Errorf("Inconsistent bidirectional mapping: key=%d, value=%d, retrievedKey=%d", key, value, retrievedKey)
		}
	}
}

// TestBiMap_Size 测试Size方法
func TestBiMap_Size(t *testing.T) {
	bm := NewBiMap[string, int]()

	if bm.Size() != 0 {
		t.Errorf("Expected initial size 0, got %d", bm.Size())
	}

	bm.Put("key1", 100)
	if bm.Size() != 1 {
		t.Errorf("Expected size 1, got %d", bm.Size())
	}

	bm.Put("key2", 200)
	if bm.Size() != 2 {
		t.Errorf("Expected size 2, got %d", bm.Size())
	}

	bm.Remove("key1")
	if bm.Size() != 1 {
		t.Errorf("Expected size 1 after remove, got %d", bm.Size())
	}
}

// TestBidirectionalConsistency 测试双向映射的一致性
func TestBidirectionalConsistency(t *testing.T) {
	bm := NewBiMap[string, int]()

	// 添加多个映射
	bm.Put("a", 1)
	bm.Put("b", 2)
	bm.Put("c", 3)

	// 验证双向一致性
	for _, key := range bm.Keys() {
		value := bm.Get(key)
		retrievedKey := bm.GetKey(value)
		if retrievedKey != key {
			t.Errorf("Bidirectional inconsistency: key=%s maps to value=%d, but value=%d maps back to key=%s",
				key, value, value, retrievedKey)
		}
	}

	// 更新后再次验证
	bm.Put("a", 10)
	for _, key := range bm.Keys() {
		value := bm.Get(key)
		retrievedKey := bm.GetKey(value)
		if retrievedKey != key {
			t.Errorf("Bidirectional inconsistency after update: key=%s maps to value=%d, but value=%d maps back to key=%s",
				key, value, value, retrievedKey)
		}
	}
}

// TestEmptyBiMapOperations 测试空BiMap的各种操作
func TestEmptyBiMapOperations(t *testing.T) {
	bm := NewBiMap[string, int]()

	// 所有操作都应该正常工作，不panic
	_ = bm.Get("any")
	_ = bm.GetKey(123)
	bm.Remove("any")
	bm.RemoveValue(123)
	bm.Clear()

	if bm.ContainsKey("any") {
		t.Error("Empty BiMap should not contain any key")
	}
	if bm.ContainsValue(123) {
		t.Error("Empty BiMap should not contain any value")
	}

	keys := bm.Keys()
	if len(keys) != 0 {
		t.Error("Empty BiMap should return empty keys slice")
	}

	values := bm.Values()
	if len(values) != 0 {
		t.Error("Empty BiMap should return empty values slice")
	}

	m := bm.ToMap()
	if len(m) != 0 {
		t.Error("Empty BiMap should return empty map")
	}
}

// TestDifferentTypes 测试不同类型的BiMap
func TestDifferentTypes(t *testing.T) {
	// 测试int->string
	bm1 := NewBiMap[int, string]()
	bm1.Put(1, "one")
	if bm1.Get(1) != "one" {
		t.Error("int->string BiMap failed")
	}
	if bm1.GetKey("one") != 1 {
		t.Error("int->string BiMap reverse lookup failed")
	}

	// 测试string->string
	bm2 := NewBiMap[string, string]()
	bm2.Put("en", "hello")
	bm2.Put("cn", "你好")
	if bm2.Get("en") != "hello" {
		t.Error("string->string BiMap failed")
	}
	if bm2.GetKey("你好") != "cn" {
		t.Error("string->string BiMap with unicode failed")
	}
}

// TestUnmarshalJSON_InvalidJSON 测试反序列化无效的JSON
func TestUnmarshalJSON_InvalidJSON(t *testing.T) {
	invalidJSON := `{"key1": invalid}`

	bm := &BiMap[string, int]{}
	err := json.Unmarshal([]byte(invalidJSON), bm)
	if err == nil {
		t.Error("Expected error when unmarshaling invalid JSON")
	}
}

// TestUnmarshalBSON_InvalidBSON 测试反序列化无效的BSON
func TestUnmarshalBSON_InvalidBSON(t *testing.T) {
	invalidBSON := []byte{0x00, 0x01, 0x02} // 无效的BSON数据

	bm := &BiMap[string, int]{}
	err := bson.Unmarshal(invalidBSON, bm)
	if err == nil {
		t.Error("Expected error when unmarshaling invalid BSON")
	}
}

// BenchmarkBiMap_Put 性能测试：Put操作
func BenchmarkBiMap_Put(b *testing.B) {
	bm := NewBiMap[int, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bm.Put(i, i*10)
	}
}

// BenchmarkBiMap_Get 性能测试：Get操作
func BenchmarkBiMap_Get(b *testing.B) {
	bm := NewBiMap[int, int]()
	for i := 0; i < 1000; i++ {
		bm.Put(i, i*10)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bm.Get(i % 1000)
	}
}

// BenchmarkGetKey 性能测试：GetKey操作
func BenchmarkGetKey(b *testing.B) {
	bm := NewBiMap[int, int]()
	for i := 0; i < 1000; i++ {
		bm.Put(i, i*10)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bm.GetKey((i % 1000) * 10)
	}
}

// BenchmarkBiMap_ConcurrentReadWrite 性能测试：并发读写
func BenchmarkBiMap_ConcurrentReadWrite(b *testing.B) {
	bm := NewBiMap[int, int]()
	for i := 0; i < 100; i++ {
		bm.Put(i, i*10)
	}

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				bm.Put(i%100, i*10)
			} else {
				_ = bm.Get(i % 100)
			}
			i++
		}
	})
}
