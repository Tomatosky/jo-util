package mapUtil

import (
	"encoding/json"
	"sync"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

// TestNewConcurrentHashMap 测试构造函数
func TestNewConcurrentHashMap(t *testing.T) {
	t.Run("创建空map", func(t *testing.T) {
		cm := NewConcurrentHashMap[string, int]()
		if cm == nil {
			t.Fatal("期望创建非空ConcurrentHashMap")
		}
		if cm.Size() != 0 {
			t.Errorf("期望size为0, 实际为 %d", cm.Size())
		}
	})

	t.Run("使用初始map创建", func(t *testing.T) {
		initMap := map[string]int{
			"one":   1,
			"two":   2,
			"three": 3,
		}
		cm := NewConcurrentHashMap[string, int](initMap)
		if cm.Size() != 3 {
			t.Errorf("期望size为3, 实际为 %d", cm.Size())
		}
		if cm.Get("one") != 1 {
			t.Errorf("期望one的值为1, 实际为 %d", cm.Get("one"))
		}
		if cm.Get("two") != 2 {
			t.Errorf("期望two的值为2, 实际为 %d", cm.Get("two"))
		}
	})

	t.Run("使用nil初始map创建", func(t *testing.T) {
		cm := NewConcurrentHashMap[string, int](nil)
		if cm.Size() != 0 {
			t.Errorf("期望size为0, 实际为 %d", cm.Size())
		}
	})
}

// TestPutAndGet 测试Put和Get方法
func TestPutAndGet(t *testing.T) {
	cm := NewConcurrentHashMap[string, string]()

	t.Run("插入并获取单个元素", func(t *testing.T) {
		cm.Put("key1", "value1")
		if value := cm.Get("key1"); value != "value1" {
			t.Errorf("期望value1, 实际为 %s", value)
		}
	})

	t.Run("覆盖已存在的key", func(t *testing.T) {
		cm.Put("key1", "newValue")
		if value := cm.Get("key1"); value != "newValue" {
			t.Errorf("期望newValue, 实际为 %s", value)
		}
	})

	t.Run("获取不存在的key", func(t *testing.T) {
		value := cm.Get("nonexistent")
		if value != "" {
			t.Errorf("期望空字符串, 实际为 %s", value)
		}
	})
}

// TestRemove 测试Remove方法
func TestRemove(t *testing.T) {
	cm := NewConcurrentHashMap[string, int]()
	cm.Put("key1", 100)
	cm.Put("key2", 200)

	t.Run("移除已存在的key", func(t *testing.T) {
		cm.Remove("key1")
		if cm.ContainsKey("key1") {
			t.Error("期望key1已被移除")
		}
		if cm.Size() != 1 {
			t.Errorf("期望size为1, 实际为 %d", cm.Size())
		}
	})

	t.Run("移除不存在的key", func(t *testing.T) {
		cm.Remove("nonexistent")
		if cm.Size() != 1 {
			t.Errorf("期望size保持为1, 实际为 %d", cm.Size())
		}
	})
}

// TestSize 测试Size方法
func TestSize(t *testing.T) {
	cm := NewConcurrentHashMap[int, string]()

	if cm.Size() != 0 {
		t.Errorf("期望初始size为0, 实际为 %d", cm.Size())
	}

	cm.Put(1, "one")
	if cm.Size() != 1 {
		t.Errorf("期望size为1, 实际为 %d", cm.Size())
	}

	cm.Put(2, "two")
	cm.Put(3, "three")
	if cm.Size() != 3 {
		t.Errorf("期望size为3, 实际为 %d", cm.Size())
	}

	cm.Remove(2)
	if cm.Size() != 2 {
		t.Errorf("期望size为2, 实际为 %d", cm.Size())
	}

	cm.Clear()
	if cm.Size() != 0 {
		t.Errorf("期望清空后size为0, 实际为 %d", cm.Size())
	}
}

// TestConcurrentHashMapContainsKey 测试ContainsKey方法
func TestConcurrentHashMapContainsKey(t *testing.T) {
	cm := NewConcurrentHashMap[string, int]()
	cm.Put("exists", 1)

	t.Run("检查存在的key", func(t *testing.T) {
		if !cm.ContainsKey("exists") {
			t.Error("期望key存在")
		}
	})

	t.Run("检查不存在的key", func(t *testing.T) {
		if cm.ContainsKey("nonexistent") {
			t.Error("期望key不存在")
		}
	})

	t.Run("移除后检查key", func(t *testing.T) {
		cm.Remove("exists")
		if cm.ContainsKey("exists") {
			t.Error("期望key已被移除")
		}
	})
}

// TestClear 测试Clear方法
func TestClear(t *testing.T) {
	cm := NewConcurrentHashMap[string, int]()
	cm.Put("key1", 1)
	cm.Put("key2", 2)
	cm.Put("key3", 3)

	cm.Clear()

	if cm.Size() != 0 {
		t.Errorf("期望清空后size为0, 实际为 %d", cm.Size())
	}

	if cm.ContainsKey("key1") || cm.ContainsKey("key2") || cm.ContainsKey("key3") {
		t.Error("期望清空后所有key都不存在")
	}

	// 清空后应该可以继续使用
	cm.Put("newKey", 100)
	if cm.Get("newKey") != 100 {
		t.Error("清空后无法正常插入新元素")
	}
}

// TestConcurrentHashMapKeys 测试Keys方法
func TestConcurrentHashMapKeys(t *testing.T) {
	cm := NewConcurrentHashMap[string, int]()

	t.Run("空map的keys", func(t *testing.T) {
		keys := cm.Keys()
		if len(keys) != 0 {
			t.Errorf("期望keys长度为0, 实际为 %d", len(keys))
		}
	})

	t.Run("获取所有keys", func(t *testing.T) {
		cm.Put("a", 1)
		cm.Put("b", 2)
		cm.Put("c", 3)

		keys := cm.Keys()
		if len(keys) != 3 {
			t.Errorf("期望keys长度为3, 实际为 %d", len(keys))
		}

		keyMap := make(map[string]bool)
		for _, k := range keys {
			keyMap[k] = true
		}

		if !keyMap["a"] || !keyMap["b"] || !keyMap["c"] {
			t.Error("返回的keys不完整")
		}
	})
}

// TestConcurrentHashMapValues 测试Values方法
func TestConcurrentHashMapValues(t *testing.T) {
	cm := NewConcurrentHashMap[int, string]()

	t.Run("空map的values", func(t *testing.T) {
		values := cm.Values()
		if len(values) != 0 {
			t.Errorf("期望values长度为0, 实际为 %d", len(values))
		}
	})

	t.Run("获取所有values", func(t *testing.T) {
		cm.Put(1, "one")
		cm.Put(2, "two")
		cm.Put(3, "three")

		values := cm.Values()
		if len(values) != 3 {
			t.Errorf("期望values长度为3, 实际为 %d", len(values))
		}

		valueMap := make(map[string]bool)
		for _, v := range values {
			valueMap[v] = true
		}

		if !valueMap["one"] || !valueMap["two"] || !valueMap["three"] {
			t.Error("返回的values不完整")
		}
	})
}

// TestConcurrentHashMapPutIfAbsent 测试PutIfAbsent方法
func TestConcurrentHashMapPutIfAbsent(t *testing.T) {
	cm := NewConcurrentHashMap[string, int]()

	t.Run("插入新key", func(t *testing.T) {
		existing, loaded := cm.PutIfAbsent("key1", 100)
		if loaded {
			t.Error("期望loaded为false")
		}
		if existing != 100 {
			t.Errorf("期望返回插入的值100, 实际为 %d", existing)
		}
		if cm.Get("key1") != 100 {
			t.Error("插入失败")
		}
	})

	t.Run("尝试插入已存在的key", func(t *testing.T) {
		existing, loaded := cm.PutIfAbsent("key1", 200)
		if !loaded {
			t.Error("期望loaded为true")
		}
		if existing != 100 {
			t.Errorf("期望返回现有值100, 实际为 %d", existing)
		}
		if cm.Get("key1") != 100 {
			t.Error("已存在的值被意外修改")
		}
	})
}

// TestConcurrentHashMapGetOrDefault 测试GetOrDefault方法
func TestConcurrentHashMapGetOrDefault(t *testing.T) {
	cm := NewConcurrentHashMap[string, int]()
	cm.Put("exists", 100)

	t.Run("获取存在的key", func(t *testing.T) {
		value := cm.GetOrDefault("exists", 999)
		if value != 100 {
			t.Errorf("期望返回100, 实际为 %d", value)
		}
	})

	t.Run("获取不存在的key", func(t *testing.T) {
		value := cm.GetOrDefault("nonexistent", 999)
		if value != 999 {
			t.Errorf("期望返回默认值999, 实际为 %d", value)
		}
	})
}

// TestToMap 测试ToMap方法
func TestToMap(t *testing.T) {
	cm := NewConcurrentHashMap[string, int]()
	cm.Put("a", 1)
	cm.Put("b", 2)
	cm.Put("c", 3)

	m := cm.ToMap()

	if len(m) != 3 {
		t.Errorf("期望map长度为3, 实际为 %d", len(m))
	}

	if m["a"] != 1 || m["b"] != 2 || m["c"] != 3 {
		t.Error("ToMap返回的map内容不正确")
	}
}

// TestRange 测试Range方法
func TestRange(t *testing.T) {
	cm := NewConcurrentHashMap[string, int]()
	cm.Put("a", 1)
	cm.Put("b", 2)
	cm.Put("c", 3)

	t.Run("正常遍历所有元素", func(t *testing.T) {
		sum := 0
		count := 0
		cm.Range(func(key string, value int) bool {
			sum += value
			count++
			return true
		})

		if count != 3 {
			t.Errorf("期望遍历3个元素, 实际遍历 %d", count)
		}
		if sum != 6 {
			t.Errorf("期望sum为6, 实际为 %d", sum)
		}
	})

	t.Run("提前终止遍历", func(t *testing.T) {
		count := 0
		cm.Range(func(key string, value int) bool {
			count++
			return count < 2 // 只遍历2个元素
		})

		if count != 2 {
			t.Errorf("期望遍历2个元素后终止, 实际遍历 %d", count)
		}
	})

	t.Run("空map遍历", func(t *testing.T) {
		emptyCm := NewConcurrentHashMap[string, int]()
		count := 0
		emptyCm.Range(func(key string, value int) bool {
			count++
			return true
		})

		if count != 0 {
			t.Errorf("期望遍历0个元素, 实际遍历 %d", count)
		}
	})
}

// TestConcurrentHashMapToString 测试ToString方法
func TestConcurrentHashMapToString(t *testing.T) {
	cm := NewConcurrentHashMap[string, int]()
	cm.Put("a", 1)
	cm.Put("b", 2)

	str := cm.ToString()
	if str == "" {
		t.Error("期望返回非空字符串")
	}

	// 验证是否为合法JSON
	var m map[string]int
	if err := json.Unmarshal([]byte(str), &m); err != nil {
		t.Errorf("返回的字符串不是合法的JSON: %v", err)
	}

	if m["a"] != 1 || m["b"] != 2 {
		t.Error("JSON解析后的内容不正确")
	}
}

// TestJSONSerialization 测试JSON序列化和反序列化
func TestJSONSerialization(t *testing.T) {
	t.Run("序列化", func(t *testing.T) {
		cm := NewConcurrentHashMap[string, int]()
		cm.Put("x", 10)
		cm.Put("y", 20)

		data, err := json.Marshal(cm)
		if err != nil {
			t.Fatalf("序列化失败: %v", err)
		}

		var m map[string]int
		if err := json.Unmarshal(data, &m); err != nil {
			t.Fatalf("反序列化为map失败: %v", err)
		}

		if m["x"] != 10 || m["y"] != 20 {
			t.Error("序列化后的内容不正确")
		}
	})

	t.Run("反序列化", func(t *testing.T) {
		jsonData := []byte(`{"name":"test","count":42}`)

		cm := NewConcurrentHashMap[string, interface{}]()
		if err := json.Unmarshal(jsonData, cm); err != nil {
			t.Fatalf("反序列化失败: %v", err)
		}

		if cm.Size() != 2 {
			t.Errorf("期望size为2, 实际为 %d", cm.Size())
		}

		if !cm.ContainsKey("name") || !cm.ContainsKey("count") {
			t.Error("反序列化后的key不完整")
		}
	})

	t.Run("空map序列化", func(t *testing.T) {
		cm := NewConcurrentHashMap[string, int]()
		data, err := json.Marshal(cm)
		if err != nil {
			t.Fatalf("序列化失败: %v", err)
		}

		if string(data) != "{}" {
			t.Errorf("期望序列化为{}, 实际为 %s", string(data))
		}
	})
}

// TestBSONSerialization 测试BSON序列化和反序列化
func TestBSONSerialization(t *testing.T) {
	t.Run("序列化", func(t *testing.T) {
		cm := NewConcurrentHashMap[string, int]()
		cm.Put("x", 10)
		cm.Put("y", 20)

		data, err := bson.Marshal(cm)
		if err != nil {
			t.Fatalf("BSON序列化失败: %v", err)
		}

		var m map[string]int
		if err := bson.Unmarshal(data, &m); err != nil {
			t.Fatalf("BSON反序列化为map失败: %v", err)
		}

		if m["x"] != 10 || m["y"] != 20 {
			t.Error("BSON序列化后的内容不正确")
		}
	})

	t.Run("反序列化", func(t *testing.T) {
		originalMap := map[string]int{"a": 1, "b": 2}
		bsonData, err := bson.Marshal(originalMap)
		if err != nil {
			t.Fatalf("准备BSON数据失败: %v", err)
		}

		cm := NewConcurrentHashMap[string, int]()
		if err := bson.Unmarshal(bsonData, cm); err != nil {
			t.Fatalf("BSON反序列化失败: %v", err)
		}

		if cm.Size() != 2 {
			t.Errorf("期望size为2, 实际为 %d", cm.Size())
		}

		if cm.Get("a") != 1 || cm.Get("b") != 2 {
			t.Error("BSON反序列化后的内容不正确")
		}
	})
}

// TestConcurrentAccess 测试并发访问的安全性
func TestConcurrentAccess(t *testing.T) {
	cm := NewConcurrentHashMap[int, int]()
	concurrency := 100
	iterations := 1000

	t.Run("并发写入", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(concurrency)

		for i := 0; i < concurrency; i++ {
			go func(id int) {
				defer wg.Done()
				for j := 0; j < iterations; j++ {
					key := id*iterations + j
					cm.Put(key, key*2)
				}
			}(i)
		}

		wg.Wait()

		expectedSize := concurrency * iterations
		if cm.Size() != expectedSize {
			t.Errorf("期望size为%d, 实际为 %d", expectedSize, cm.Size())
		}
	})

	t.Run("并发读写混合", func(t *testing.T) {
		cm.Clear()
		// 先插入一些数据
		for i := 0; i < 1000; i++ {
			cm.Put(i, i)
		}

		var wg sync.WaitGroup
		wg.Add(concurrency * 3)

		// 并发读
		for i := 0; i < concurrency; i++ {
			go func() {
				defer wg.Done()
				for j := 0; j < iterations; j++ {
					cm.Get(j % 1000)
				}
			}()
		}

		// 并发写
		for i := 0; i < concurrency; i++ {
			go func(id int) {
				defer wg.Done()
				for j := 0; j < iterations; j++ {
					cm.Put(j%1000, id)
				}
			}(i)
		}

		// 并发删除
		for i := 0; i < concurrency; i++ {
			go func() {
				defer wg.Done()
				for j := 0; j < iterations; j++ {
					cm.Remove(j % 1000)
				}
			}()
		}

		wg.Wait()
		// 测试通过意味着没有竞态条件或死锁
	})

	t.Run("并发Range操作", func(t *testing.T) {
		cm.Clear()
		for i := 0; i < 100; i++ {
			cm.Put(i, i)
		}

		var wg sync.WaitGroup
		wg.Add(concurrency)

		for i := 0; i < concurrency; i++ {
			go func() {
				defer wg.Done()
				cm.Range(func(key, value int) bool {
					return true
				})
			}()
		}

		wg.Wait()
	})
}

// TestEdgeCases 测试边界情况
func TestEdgeCases(t *testing.T) {
	t.Run("零值测试", func(t *testing.T) {
		cm := NewConcurrentHashMap[string, int]()
		cm.Put("zero", 0)

		if !cm.ContainsKey("zero") {
			t.Error("存储零值失败")
		}

		if cm.Get("zero") != 0 {
			t.Error("获取零值失败")
		}
	})

	t.Run("空字符串作为key", func(t *testing.T) {
		cm := NewConcurrentHashMap[string, int]()
		cm.Put("", 100)

		if !cm.ContainsKey("") {
			t.Error("空字符串作为key失败")
		}

		if cm.Get("") != 100 {
			t.Error("获取空字符串key的值失败")
		}
	})

	t.Run("大量数据", func(t *testing.T) {
		cm := NewConcurrentHashMap[int, int]()
		count := 10000

		for i := 0; i < count; i++ {
			cm.Put(i, i*2)
		}

		if cm.Size() != count {
			t.Errorf("期望size为%d, 实际为 %d", count, cm.Size())
		}

		// 验证部分数据
		for i := 0; i < count; i += 1000 {
			if cm.Get(i) != i*2 {
				t.Errorf("key %d 的值不正确", i)
			}
		}
	})

	t.Run("结构体作为值", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		cm := NewConcurrentHashMap[string, Person]()
		cm.Put("person1", Person{Name: "Alice", Age: 30})

		person := cm.Get("person1")
		if person.Name != "Alice" || person.Age != 30 {
			t.Error("存储结构体失败")
		}
	})

	t.Run("指针作为值", func(t *testing.T) {
		cm := NewConcurrentHashMap[string, *int]()
		value := 42
		cm.Put("ptr", &value)

		ptr := cm.Get("ptr")
		if ptr == nil || *ptr != 42 {
			t.Error("存储指针失败")
		}
	})
}

// TestIMapInterface 测试IMap接口实现
func TestIMapInterface(t *testing.T) {
	var imap IMap[string, int] = NewConcurrentHashMap[string, int]()

	imap.Put("key1", 100)
	if imap.Get("key1") != 100 {
		t.Error("通过接口操作失败")
	}

	if imap.Size() != 1 {
		t.Error("接口Size方法失败")
	}

	if !imap.ContainsKey("key1") {
		t.Error("接口ContainsKey方法失败")
	}

	imap.Remove("key1")
	if imap.ContainsKey("key1") {
		t.Error("接口Remove方法失败")
	}
}

// BenchmarkPut 基准测试: Put操作
func BenchmarkPut(b *testing.B) {
	cm := NewConcurrentHashMap[int, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cm.Put(i, i)
	}
}

// BenchmarkGet 基准测试: Get操作
func BenchmarkGet(b *testing.B) {
	cm := NewConcurrentHashMap[int, int]()
	for i := 0; i < 10000; i++ {
		cm.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cm.Get(i % 10000)
	}
}

// BenchmarkConcurrentReadWrite 基准测试: 并发读写
func BenchmarkConcurrentReadWrite(b *testing.B) {
	cm := NewConcurrentHashMap[int, int]()
	for i := 0; i < 1000; i++ {
		cm.Put(i, i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				cm.Get(i % 1000)
			} else {
				cm.Put(i%1000, i)
			}
			i++
		}
	})
}

// BenchmarkRange 基准测试: Range操作
func BenchmarkRange(b *testing.B) {
	cm := NewConcurrentHashMap[int, int]()
	for i := 0; i < 1000; i++ {
		cm.Put(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cm.Range(func(key, value int) bool {
			return true
		})
	}
}

// BenchmarkPutIfAbsent 基准测试: PutIfAbsent操作
func BenchmarkPutIfAbsent(b *testing.B) {
	cm := NewConcurrentHashMap[int, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cm.PutIfAbsent(i%1000, i)
	}
}
