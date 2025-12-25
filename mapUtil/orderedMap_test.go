package mapUtil

import (
	"encoding/json"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

// TestOrderedMapNew 测试创建新的OrderedMap
func TestOrderedMapNew(t *testing.T) {
	om := NewOrderedMap[string, int]()
	if om == nil {
		t.Fatal("NewOrderedMap返回了nil")
	}
	if om.Size() != 0 {
		t.Errorf("新创建的OrderedMap大小应该为0，实际为%d", om.Size())
	}
}

// TestOrderedMapNewWithCapacity 测试带容量创建OrderedMap
func TestOrderedMapNewWithCapacity(t *testing.T) {
	om := NewOrderedMapWithCapacity[string, int](10)
	if om == nil {
		t.Fatal("NewOrderedMapWithCapacity返回了nil")
	}
	if om.Size() != 0 {
		t.Errorf("新创建的OrderedMap大小应该为0，实际为%d", om.Size())
	}
}

// TestOrderedMapNewWithElements 测试用元素创建OrderedMap
func TestOrderedMapNewWithElements(t *testing.T) {
	el1 := &Element[string, int]{Key: "a", Value: 1}
	el2 := &Element[string, int]{Key: "b", Value: 2}
	om := NewOrderedMapWithElements(el1, el2)

	if om.Size() != 2 {
		t.Errorf("OrderedMap大小应该为2，实际为%d", om.Size())
	}
	if om.Get("a") != 1 {
		t.Errorf("期望值为1，实际为%d", om.Get("a"))
	}
	if om.Get("b") != 2 {
		t.Errorf("期望值为2，实际为%d", om.Get("b"))
	}
}

// TestOrderedMapPutAndGet 测试基本的Put和Get操作
func TestOrderedMapPutAndGet(t *testing.T) {
	om := NewOrderedMap[string, int]()

	om.Put("key1", 100)
	if om.Get("key1") != 100 {
		t.Errorf("期望值为100，实际为%d", om.Get("key1"))
	}

	// 测试覆盖
	om.Put("key1", 200)
	if om.Get("key1") != 200 {
		t.Errorf("期望值为200，实际为%d", om.Get("key1"))
	}
	if om.Size() != 1 {
		t.Errorf("覆盖后大小应该为1，实际为%d", om.Size())
	}
}

// TestOrderedMapGetNonExistentKey 测试获取不存在的键
func TestOrderedMapGetNonExistentKey(t *testing.T) {
	om := NewOrderedMap[string, int]()
	value := om.Get("nonexistent")
	if value != 0 {
		t.Errorf("不存在的键应该返回零值，实际为%d", value)
	}
}

// TestOrderedMapSize 测试Size方法
func TestOrderedMapSize(t *testing.T) {
	om := NewOrderedMap[string, int]()

	if om.Size() != 0 {
		t.Errorf("空map的大小应该为0，实际为%d", om.Size())
	}

	om.Put("a", 1)
	if om.Size() != 1 {
		t.Errorf("添加一个元素后大小应该为1，实际为%d", om.Size())
	}

	om.Put("b", 2)
	om.Put("c", 3)
	if om.Size() != 3 {
		t.Errorf("添加三个元素后大小应该为3，实际为%d", om.Size())
	}

	om.Remove("b")
	if om.Size() != 2 {
		t.Errorf("删除一个元素后大小应该为2，实际为%d", om.Size())
	}
}

// TestOrderedMapContainsKey 测试ContainsKey方法
func TestOrderedMapContainsKey(t *testing.T) {
	om := NewOrderedMap[string, int]()

	if om.ContainsKey("key1") {
		t.Error("空map不应该包含任何键")
	}

	om.Put("key1", 100)
	if !om.ContainsKey("key1") {
		t.Error("map应该包含key1")
	}

	if om.ContainsKey("key2") {
		t.Error("map不应该包含key2")
	}
}

// TestOrderedMapRemove 测试Remove方法
func TestOrderedMapRemove(t *testing.T) {
	om := NewOrderedMap[string, int]()

	om.Put("a", 1)
	om.Put("b", 2)
	om.Put("c", 3)

	om.Remove("b")
	if om.Size() != 2 {
		t.Errorf("删除后大小应该为2，实际为%d", om.Size())
	}
	if om.ContainsKey("b") {
		t.Error("删除后不应该包含b键")
	}

	// 测试删除不存在的键
	om.Remove("nonexistent")
	if om.Size() != 2 {
		t.Errorf("删除不存在的键后大小应该保持2，实际为%d", om.Size())
	}
}

// TestOrderedMapClear 测试Clear方法
func TestOrderedMapClear(t *testing.T) {
	om := NewOrderedMap[string, int]()

	om.Put("a", 1)
	om.Put("b", 2)
	om.Put("c", 3)

	om.Clear()
	if om.Size() != 0 {
		t.Errorf("Clear后大小应该为0，实际为%d", om.Size())
	}
	if om.ContainsKey("a") {
		t.Error("Clear后不应该包含任何键")
	}
}

// TestOrderedMapPutIfAbsent 测试PutIfAbsent方法
func TestOrderedMapPutIfAbsent(t *testing.T) {
	om := NewOrderedMap[string, int]()

	// 添加新键
	existing, loaded := om.PutIfAbsent("key1", 100)
	if loaded {
		t.Error("新键不应该已存在")
	}
	if existing != 100 {
		t.Errorf("期望返回值为100，实际为%d", existing)
	}

	// 尝试添加已存在的键
	existing, loaded = om.PutIfAbsent("key1", 200)
	if !loaded {
		t.Error("已存在的键应该返回loaded=true")
	}
	if existing != 100 {
		t.Errorf("期望返回已存在的值100，实际为%d", existing)
	}
	if om.Get("key1") != 100 {
		t.Errorf("值不应该被覆盖，期望100，实际为%d", om.Get("key1"))
	}
}

// TestOrderedMapGetOrDefault 测试GetOrDefault方法
func TestOrderedMapGetOrDefault(t *testing.T) {
	om := NewOrderedMap[string, int]()

	om.Put("key1", 100)

	// 获取存在的键
	value := om.GetOrDefault("key1", 999)
	if value != 100 {
		t.Errorf("期望值为100，实际为%d", value)
	}

	// 获取不存在的键
	value = om.GetOrDefault("nonexistent", 999)
	if value != 999 {
		t.Errorf("期望默认值为999，实际为%d", value)
	}
}

// TestOrderedMapGetElement 测试GetElement方法
func TestOrderedMapGetElement(t *testing.T) {
	om := NewOrderedMap[string, int]()

	om.Put("key1", 100)

	el := om.GetElement("key1")
	if el == nil {
		t.Fatal("GetElement不应该返回nil")
	}
	if el.Key != "key1" || el.Value != 100 {
		t.Errorf("期望Key=key1, Value=100，实际Key=%s, Value=%d", el.Key, el.Value)
	}

	// 获取不存在的键
	el = om.GetElement("nonexistent")
	if el != nil {
		t.Error("不存在的键应该返回nil")
	}
}

// TestOrderedMapReplaceKey 测试ReplaceKey方法
func TestOrderedMapReplaceKey(t *testing.T) {
	om := NewOrderedMap[string, int]()

	om.Put("a", 1)
	om.Put("b", 2)
	om.Put("c", 3)

	// 成功替换
	if !om.ReplaceKey("b", "d") {
		t.Error("ReplaceKey应该成功")
	}
	if om.ContainsKey("b") {
		t.Error("旧键b应该不存在")
	}
	if !om.ContainsKey("d") {
		t.Error("新键d应该存在")
	}
	if om.Get("d") != 2 {
		t.Errorf("新键d的值应该为2，实际为%d", om.Get("d"))
	}

	// 替换不存在的键
	if om.ReplaceKey("nonexistent", "e") {
		t.Error("替换不存在的键应该失败")
	}

	// 新键已存在
	if om.ReplaceKey("d", "a") {
		t.Error("新键已存在时应该失败")
	}
}

// TestOrderedMapOrderPreservation 测试插入顺序保持
func TestOrderedMapOrderPreservation(t *testing.T) {
	om := NewOrderedMap[string, int]()

	om.Put("a", 1)
	om.Put("b", 2)
	om.Put("c", 3)
	om.Put("d", 4)

	keys := om.Keys()
	expectedKeys := []string{"a", "b", "c", "d"}
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("期望顺序为%v，实际为%v", expectedKeys, keys)
	}

	values := om.Values()
	expectedValues := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(values, expectedValues) {
		t.Errorf("期望值顺序为%v，实际为%v", expectedValues, values)
	}
}

// TestOrderedMapFrontAndBack 测试Front和Back方法
func TestOrderedMapFrontAndBack(t *testing.T) {
	om := NewOrderedMap[string, int]()

	// 空map
	if om.Front() != nil {
		t.Error("空map的Front应该为nil")
	}
	if om.Back() != nil {
		t.Error("空map的Back应该为nil")
	}

	om.Put("a", 1)
	om.Put("b", 2)
	om.Put("c", 3)

	front := om.Front()
	if front == nil || front.Key != "a" || front.Value != 1 {
		t.Error("Front应该返回第一个元素")
	}

	back := om.Back()
	if back == nil || back.Key != "c" || back.Value != 3 {
		t.Error("Back应该返回最后一个元素")
	}
}

// TestOrderedMapIterationFromFront 测试从前向后遍历
func TestOrderedMapIterationFromFront(t *testing.T) {
	om := NewOrderedMap[string, int]()

	om.Put("a", 1)
	om.Put("b", 2)
	om.Put("c", 3)

	expectedKeys := []string{"a", "b", "c"}
	expectedValues := []int{1, 2, 3}
	i := 0

	for key, value := range om.AllFromFront() {
		if key != expectedKeys[i] {
			t.Errorf("索引%d: 期望键为%s，实际为%s", i, expectedKeys[i], key)
		}
		if value != expectedValues[i] {
			t.Errorf("索引%d: 期望值为%d，实际为%d", i, expectedValues[i], value)
		}
		i++
	}

	if i != 3 {
		t.Errorf("期望遍历3个元素，实际遍历%d个", i)
	}
}

// TestOrderedMapIterationFromBack 测试从后向前遍历
func TestOrderedMapIterationFromBack(t *testing.T) {
	om := NewOrderedMap[string, int]()

	om.Put("a", 1)
	om.Put("b", 2)
	om.Put("c", 3)

	expectedKeys := []string{"c", "b", "a"}
	expectedValues := []int{3, 2, 1}
	i := 0

	for key, value := range om.AllFromBack() {
		if key != expectedKeys[i] {
			t.Errorf("索引%d: 期望键为%s，实际为%s", i, expectedKeys[i], key)
		}
		if value != expectedValues[i] {
			t.Errorf("索引%d: 期望值为%d，实际为%d", i, expectedValues[i], value)
		}
		i++
	}

	if i != 3 {
		t.Errorf("期望遍历3个元素，实际遍历%d个", i)
	}
}

// TestOrderedMapIterationBreak 测试迭代器的中断
func TestOrderedMapIterationBreak(t *testing.T) {
	om := NewOrderedMap[string, int]()

	om.Put("a", 1)
	om.Put("b", 2)
	om.Put("c", 3)

	count := 0
	for range om.AllFromFront() {
		count++
		if count == 2 {
			break
		}
	}

	if count != 2 {
		t.Errorf("期望中断后count为2，实际为%d", count)
	}
}

// TestOrderedMapRange 测试Range方法
func TestOrderedMapRange(t *testing.T) {
	om := NewOrderedMap[string, int]()

	om.Put("a", 1)
	om.Put("b", 2)
	om.Put("c", 3)

	// 完整遍历
	count := 0
	om.Range(func(key string, value int) bool {
		count++
		return true
	})
	if count != 3 {
		t.Errorf("期望遍历3个元素，实际遍历%d个", count)
	}

	// 中断遍历
	count = 0
	om.Range(func(key string, value int) bool {
		count++
		return count < 2
	})
	if count != 2 {
		t.Errorf("期望中断后count为2，实际为%d", count)
	}
}

// TestOrderedMapCopy 测试Copy方法
func TestOrderedMapCopy(t *testing.T) {
	om := NewOrderedMap[string, int]()

	om.Put("a", 1)
	om.Put("b", 2)
	om.Put("c", 3)

	om2 := om.Copy()

	if om2.Size() != om.Size() {
		t.Errorf("复制后大小应该相同，原始:%d，复制:%d", om.Size(), om2.Size())
	}

	keys1 := om.Keys()
	keys2 := om2.Keys()
	if !reflect.DeepEqual(keys1, keys2) {
		t.Errorf("复制后键顺序应该相同，原始:%v，复制:%v", keys1, keys2)
	}

	// 修改复制的map不应该影响原始map
	om2.Put("d", 4)
	if om.ContainsKey("d") {
		t.Error("修改复制的map不应该影响原始map")
	}
}

// TestOrderedMapToMap 测试ToMap方法
func TestOrderedMapToMap(t *testing.T) {
	om := NewOrderedMap[string, int]()

	om.Put("a", 1)
	om.Put("b", 2)
	om.Put("c", 3)

	m := om.ToMap()

	if len(m) != 3 {
		t.Errorf("期望map长度为3，实际为%d", len(m))
	}
	if m["a"] != 1 || m["b"] != 2 || m["c"] != 3 {
		t.Error("ToMap的值不正确")
	}
}

// TestOrderedMapToString 测试ToString方法
func TestOrderedMapToString(t *testing.T) {
	om := NewOrderedMap[string, int]()

	om.Put("a", 1)
	om.Put("b", 2)

	str := om.ToString()
	if str == "" {
		t.Error("ToString不应该返回空字符串")
	}

	// 验证是否为有效的JSON
	var m map[string]int
	err := json.Unmarshal([]byte(str), &m)
	if err != nil {
		t.Errorf("ToString应该返回有效的JSON: %v", err)
	}
}

// TestOrderedMapJSONMarshalUnmarshal 测试JSON序列化和反序列化
func TestOrderedMapJSONMarshalUnmarshal(t *testing.T) {
	om := NewOrderedMap[string, int]()

	om.Put("a", 1)
	om.Put("b", 2)
	om.Put("c", 3)

	// Marshal
	data, err := json.Marshal(om)
	if err != nil {
		t.Fatalf("JSON Marshal失败: %v", err)
	}

	// Unmarshal
	om2 := NewOrderedMap[string, int]()
	err = json.Unmarshal(data, om2)
	if err != nil {
		t.Fatalf("JSON Unmarshal失败: %v", err)
	}

	if om2.Size() != 3 {
		t.Errorf("反序列化后大小应该为3，实际为%d", om2.Size())
	}
	if om2.Get("a") != 1 || om2.Get("b") != 2 || om2.Get("c") != 3 {
		t.Error("反序列化后的值不正确")
	}
}

// TestOrderedMapBSONMarshalUnmarshal 测试BSON序列化和反序列化
func TestOrderedMapBSONMarshalUnmarshal(t *testing.T) {
	om := NewOrderedMap[string, int]()

	om.Put("a", 1)
	om.Put("b", 2)
	om.Put("c", 3)

	// Marshal
	data, err := bson.Marshal(om)
	if err != nil {
		t.Fatalf("BSON Marshal失败: %v", err)
	}

	// Unmarshal
	om2 := NewOrderedMap[string, int]()
	err = bson.Unmarshal(data, om2)
	if err != nil {
		t.Fatalf("BSON Unmarshal失败: %v", err)
	}

	if om2.Size() != 3 {
		t.Errorf("反序列化后大小应该为3，实际为%d", om2.Size())
	}
	if om2.Get("a") != 1 || om2.Get("b") != 2 || om2.Get("c") != 3 {
		t.Error("反序列化后的值不正确")
	}
}

// TestOrderedMapEmptyOperations 测试空map的各种操作
func TestOrderedMapEmptyOperations(t *testing.T) {
	om := NewOrderedMap[string, int]()

	if om.Size() != 0 {
		t.Error("空map的Size应该为0")
	}
	if om.Front() != nil {
		t.Error("空map的Front应该为nil")
	}
	if om.Back() != nil {
		t.Error("空map的Back应该为nil")
	}

	keys := om.Keys()
	if len(keys) != 0 {
		t.Error("空map的Keys应该为空切片")
	}

	values := om.Values()
	if len(values) != 0 {
		t.Error("空map的Values应该为空切片")
	}

	m := om.ToMap()
	if len(m) != 0 {
		t.Error("空map的ToMap应该为空map")
	}
}

// TestOrderedMapElementNavigation 测试Element的Next和Prev方法
func TestOrderedMapElementNavigation(t *testing.T) {
	om := NewOrderedMap[string, int]()

	om.Put("a", 1)
	om.Put("b", 2)
	om.Put("c", 3)

	el := om.Front()
	if el == nil || el.Key != "a" {
		t.Fatal("Front应该返回第一个元素")
	}

	el = el.Next()
	if el == nil || el.Key != "b" {
		t.Error("Next应该返回下一个元素")
	}

	el = el.Next()
	if el == nil || el.Key != "c" {
		t.Error("Next应该返回下一个元素")
	}

	el = el.Next()
	if el != nil {
		t.Error("最后一个元素的Next应该为nil")
	}

	// 反向遍历
	el = om.Back()
	if el == nil || el.Key != "c" {
		t.Fatal("Back应该返回最后一个元素")
	}

	el = el.Prev()
	if el == nil || el.Key != "b" {
		t.Error("Prev应该返回前一个元素")
	}

	el = el.Prev()
	if el == nil || el.Key != "a" {
		t.Error("Prev应该返回前一个元素")
	}

	el = el.Prev()
	if el != nil {
		t.Error("第一个元素的Prev应该为nil")
	}
}

// TestOrderedMapUpdateValuePreservesOrder 测试更新值保持顺序
func TestOrderedMapUpdateValuePreservesOrder(t *testing.T) {
	om := NewOrderedMap[string, int]()

	om.Put("a", 1)
	om.Put("b", 2)
	om.Put("c", 3)

	// 更新中间元素的值
	om.Put("b", 200)

	keys := om.Keys()
	expectedKeys := []string{"a", "b", "c"}
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("更新值后顺序应该保持不变，期望%v，实际%v", expectedKeys, keys)
	}

	if om.Get("b") != 200 {
		t.Errorf("期望b的值为200，实际为%d", om.Get("b"))
	}
}

// TestOrderedMapRemovePreservesOrder 测试删除保持顺序
func TestOrderedMapRemovePreservesOrder(t *testing.T) {
	om := NewOrderedMap[string, int]()

	om.Put("a", 1)
	om.Put("b", 2)
	om.Put("c", 3)
	om.Put("d", 4)

	om.Remove("b")

	keys := om.Keys()
	expectedKeys := []string{"a", "c", "d"}
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("删除后顺序应该保持不变，期望%v，实际%v", expectedKeys, keys)
	}
}

// TestOrderedMapReplaceKeyPreservesOrder 测试ReplaceKey保持顺序
func TestOrderedMapReplaceKeyPreservesOrder(t *testing.T) {
	om := NewOrderedMap[string, int]()

	om.Put("a", 1)
	om.Put("b", 2)
	om.Put("c", 3)

	om.ReplaceKey("b", "x")

	keys := om.Keys()
	expectedKeys := []string{"a", "x", "c"}
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("ReplaceKey后顺序应该保持不变，期望%v，实际%v", expectedKeys, keys)
	}

	if om.Get("x") != 2 {
		t.Errorf("期望x的值为2，实际为%d", om.Get("x"))
	}
}

// TestOrderedMapLargeDataSet 测试大数据集
func TestOrderedMapLargeDataSet(t *testing.T) {
	om := NewOrderedMap[int, int]()

	n := 10000
	for i := 0; i < n; i++ {
		om.Put(i, i*2)
	}

	if om.Size() != n {
		t.Errorf("期望大小为%d，实际为%d", n, om.Size())
	}

	// 验证随机访问
	for i := 0; i < 100; i++ {
		idx := i * 100
		if om.Get(idx) != idx*2 {
			t.Errorf("索引%d的值不正确", idx)
		}
	}

	// 验证顺序
	i := 0
	for key, value := range om.AllFromFront() {
		if key != i || value != i*2 {
			t.Errorf("索引%d的键或值不正确", i)
			break
		}
		i++
		if i >= 100 { // 只验证前100个
			break
		}
	}
}

// TestOrderedMapDifferentTypes 测试不同类型的键值
func TestOrderedMapDifferentTypes(t *testing.T) {
	// 测试int键和string值
	om1 := NewOrderedMap[int, string]()
	om1.Put(1, "one")
	om1.Put(2, "two")
	if om1.Get(1) != "one" {
		t.Error("int键string值类型测试失败")
	}

	// 测试string键和struct值
	type Person struct {
		Name string
		Age  int
	}
	om2 := NewOrderedMap[string, Person]()
	om2.Put("alice", Person{"Alice", 30})
	om2.Put("bob", Person{"Bob", 25})
	if om2.Get("alice").Name != "Alice" {
		t.Error("string键struct值类型测试失败")
	}
}

// TestOrderedMapConcurrentReadsSafety 测试并发读取的基本场景
func TestOrderedMapConcurrentReadsSafety(t *testing.T) {
	om := NewOrderedMap[string, int]()

	// 先填充数据
	for i := 0; i < 100; i++ {
		om.Put(string(rune('a'+i%26))+string(rune('0'+i/26)), i)
	}

	// 多个goroutine同时读取
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				_ = om.Size()
				_ = om.Keys()
				_ = om.Get("a0")
			}
			done <- true
		}()
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestOrderedMapIMapInterface 测试IMap接口实现
func TestOrderedMapIMapInterface(t *testing.T) {
	var imap IMap[string, int] = NewOrderedMap[string, int]()

	imap.Put("a", 1)
	imap.Put("b", 2)

	if imap.Size() != 2 {
		t.Error("IMap接口Size方法测试失败")
	}
	if imap.Get("a") != 1 {
		t.Error("IMap接口Get方法测试失败")
	}
	if !imap.ContainsKey("a") {
		t.Error("IMap接口ContainsKey方法测试失败")
	}

	imap.Remove("a")
	if imap.Size() != 1 {
		t.Error("IMap接口Remove方法测试失败")
	}

	keys := imap.Keys()
	if len(keys) != 1 || keys[0] != "b" {
		t.Error("IMap接口Keys方法测试失败")
	}

	imap.Clear()
	if imap.Size() != 0 {
		t.Error("IMap接口Clear方法测试失败")
	}
}

// TestOrderedMapSingleElement 测试单个元素的操作
func TestOrderedMapSingleElement(t *testing.T) {
	om := NewOrderedMap[string, int]()
	om.Put("only", 1)

	if om.Front() != om.Back() {
		t.Error("单个元素的map，Front和Back应该是同一个元素")
	}

	if om.Front().Next() != nil {
		t.Error("单个元素的Next应该为nil")
	}

	if om.Back().Prev() != nil {
		t.Error("单个元素的Prev应该为nil")
	}

	om.Remove("only")
	if om.Size() != 0 {
		t.Error("删除唯一元素后，大小应该为0")
	}
}

// TestOrderedMapMultipleReplaces 测试多次替换操作
func TestOrderedMapMultipleReplaces(t *testing.T) {
	om := NewOrderedMap[string, int]()
	om.Put("a", 1)
	om.Put("b", 2)
	om.Put("c", 3)

	om.Put("a", 10)
	om.Put("b", 20)
	om.Put("c", 30)

	if om.Size() != 3 {
		t.Errorf("多次替换后大小应该为3，实际为%d", om.Size())
	}

	if om.Get("a") != 10 || om.Get("b") != 20 || om.Get("c") != 30 {
		t.Error("多次替换后值不正确")
	}

	keys := om.Keys()
	expectedKeys := []string{"a", "b", "c"}
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("多次替换后顺序应该保持不变，期望%v，实际%v", expectedKeys, keys)
	}
}

// TestOrderedMapRemoveFirstAndLast 测试删除首尾元素
func TestOrderedMapRemoveFirstAndLast(t *testing.T) {
	om := NewOrderedMap[string, int]()
	om.Put("a", 1)
	om.Put("b", 2)
	om.Put("c", 3)

	// 删除第一个
	om.Remove("a")
	if om.Front().Key != "b" {
		t.Error("删除第一个元素后，Front应该是b")
	}

	om.Put("a", 1) // 重新添加
	om.Put("d", 4)

	// 删除最后一个
	om.Remove("d")
	if om.Back().Key != "a" {
		t.Error("删除最后一个元素后，Back应该是a")
	}
}

// TestOrderedMapZeroValues 测试零值
func TestOrderedMapZeroValues(t *testing.T) {
	om := NewOrderedMap[string, int]()
	om.Put("zero", 0)

	if !om.ContainsKey("zero") {
		t.Error("应该包含键zero")
	}

	if om.Get("zero") != 0 {
		t.Error("零值应该被正确存储和获取")
	}

	// 测试string零值
	om2 := NewOrderedMap[int, string]()
	om2.Put(1, "")

	if !om2.ContainsKey(1) {
		t.Error("应该包含键1")
	}

	if om2.Get(1) != "" {
		t.Error("空字符串应该被正确存储和获取")
	}
}
