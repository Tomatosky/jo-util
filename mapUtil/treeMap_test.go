package mapUtil

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

// TestTreeMapNew 测试TreeMap的创建
func TestTreeMapNew(t *testing.T) {
	tm := NewTreeMap[int, string](func(a, b int) bool {
		return a < b
	})
	if tm == nil {
		t.Fatal("NewTreeMap返回nil")
	}
	if tm.Size() != 0 {
		t.Errorf("新创建的TreeMap大小应为0，实际为%d", tm.Size())
	}
}

// TestTreeMapPutAndGet 测试基本的Put和Get操作
func TestTreeMapPutAndGet(t *testing.T) {
	tm := NewTreeMap[int, string](func(a, b int) bool {
		return a < b
	})

	// 测试插入单个元素
	tm.Put(1, "one")
	if tm.Size() != 1 {
		t.Errorf("插入一个元素后，大小应为1，实际为%d", tm.Size())
	}

	// 测试获取存在的元素
	val := tm.Get(1)
	if val != "one" {
		t.Errorf("Get(1)应返回'one'，实际为'%s'", val)
	}

	// 测试获取不存在的元素
	val = tm.Get(2)
	if val != "" {
		t.Errorf("Get(2)应返回空字符串，实际为'%s'", val)
	}
}

// TestTreeMapPutUpdate 测试更新已存在的键
func TestTreeMapPutUpdate(t *testing.T) {
	tm := NewTreeMap[int, string](func(a, b int) bool {
		return a < b
	})

	tm.Put(1, "one")
	tm.Put(1, "ONE")

	if tm.Size() != 1 {
		t.Errorf("更新已存在的键后，大小应为1，实际为%d", tm.Size())
	}

	val := tm.Get(1)
	if val != "ONE" {
		t.Errorf("Get(1)应返回更新后的值'ONE'，实际为'%s'", val)
	}
}

// TestTreeMapMultiplePuts 测试插入多个元素
func TestTreeMapMultiplePuts(t *testing.T) {
	tm := NewTreeMap[int, string](func(a, b int) bool {
		return a < b
	})

	// 插入多个元素
	for i := 1; i <= 100; i++ {
		tm.Put(i, fmt.Sprintf("value%d", i))
	}

	if tm.Size() != 100 {
		t.Errorf("插入100个元素后，大小应为100，实际为%d", tm.Size())
	}

	// 验证所有元素
	for i := 1; i <= 100; i++ {
		expected := fmt.Sprintf("value%d", i)
		val := tm.Get(i)
		if val != expected {
			t.Errorf("Get(%d)应返回'%s'，实际为'%s'", i, expected, val)
		}
	}
}

// TestTreeMapRemove 测试删除操作
func TestTreeMapRemove(t *testing.T) {
	tm := NewTreeMap[int, string](func(a, b int) bool {
		return a < b
	})

	// 插入元素
	tm.Put(1, "one")
	tm.Put(2, "two")
	tm.Put(3, "three")

	// 删除存在的元素
	tm.Remove(2)
	if tm.Size() != 2 {
		t.Errorf("删除一个元素后，大小应为2，实际为%d", tm.Size())
	}

	if tm.ContainsKey(2) {
		t.Error("删除后，ContainsKey(2)应返回false")
	}

	// 删除不存在的元素
	tm.Remove(10)
	if tm.Size() != 2 {
		t.Errorf("删除不存在的元素后，大小应保持2，实际为%d", tm.Size())
	}

	// 验证剩余元素
	if tm.Get(1) != "one" || tm.Get(3) != "three" {
		t.Error("删除后，剩余元素的值不正确")
	}
}

// TestTreeMapRemoveAll 测试删除所有元素
func TestTreeMapRemoveAll(t *testing.T) {
	tm := NewTreeMap[int, string](func(a, b int) bool {
		return a < b
	})

	// 插入元素
	for i := 1; i <= 10; i++ {
		tm.Put(i, fmt.Sprintf("value%d", i))
	}

	// 删除所有元素
	for i := 1; i <= 10; i++ {
		tm.Remove(i)
	}

	if tm.Size() != 0 {
		t.Errorf("删除所有元素后，大小应为0，实际为%d", tm.Size())
	}

	// 验证树为空
	keys := tm.Keys()
	if len(keys) != 0 {
		t.Errorf("删除所有元素后，Keys()应返回空切片，实际长度为%d", len(keys))
	}
}

// TestTreeMapContainsKey 测试ContainsKey方法
func TestTreeMapContainsKey(t *testing.T) {
	tm := NewTreeMap[int, string](func(a, b int) bool {
		return a < b
	})

	tm.Put(1, "one")
	tm.Put(2, "two")

	if !tm.ContainsKey(1) {
		t.Error("ContainsKey(1)应返回true")
	}

	if tm.ContainsKey(3) {
		t.Error("ContainsKey(3)应返回false")
	}
}

// TestTreeMapClear 测试Clear方法
func TestTreeMapClear(t *testing.T) {
	tm := NewTreeMap[int, string](func(a, b int) bool {
		return a < b
	})

	for i := 1; i <= 10; i++ {
		tm.Put(i, fmt.Sprintf("value%d", i))
	}

	tm.Clear()

	if tm.Size() != 0 {
		t.Errorf("Clear后，大小应为0，实际为%d", tm.Size())
	}

	if tm.ContainsKey(1) {
		t.Error("Clear后，应不包含任何键")
	}
}

// TestTreeMapKeys 测试Keys方法和排序
func TestTreeMapKeys(t *testing.T) {
	tm := NewTreeMap[int, string](func(a, b int) bool {
		return a < b
	})

	// 乱序插入
	order := []int{5, 2, 8, 1, 9, 3, 7, 4, 6}
	for _, v := range order {
		tm.Put(v, fmt.Sprintf("value%d", v))
	}

	keys := tm.Keys()

	// 验证键已排序
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	if len(keys) != len(expected) {
		t.Fatalf("Keys()返回的长度不正确，期望%d，实际%d", len(expected), len(keys))
	}

	for i, key := range keys {
		if key != expected[i] {
			t.Errorf("索引%d处的键应为%d，实际为%d", i, expected[i], key)
		}
	}
}

// TestTreeMapValues 测试Values方法
func TestTreeMapValues(t *testing.T) {
	tm := NewTreeMap[int, string](func(a, b int) bool {
		return a < b
	})

	// 乱序插入
	order := []int{5, 2, 8, 1, 9}
	for _, v := range order {
		tm.Put(v, fmt.Sprintf("value%d", v))
	}

	values := tm.Values()

	// 验证值按键的顺序排列
	expected := []string{"value1", "value2", "value5", "value8", "value9"}
	if len(values) != len(expected) {
		t.Fatalf("Values()返回的长度不正确，期望%d，实际%d", len(expected), len(values))
	}

	for i, val := range values {
		if val != expected[i] {
			t.Errorf("索引%d处的值应为'%s'，实际为'%s'", i, expected[i], val)
		}
	}
}

// TestTreeMapPutIfAbsent 测试PutIfAbsent方法
func TestTreeMapPutIfAbsent(t *testing.T) {
	tm := NewTreeMap[int, string](func(a, b int) bool {
		return a < b
	})

	// 测试插入新键
	existing, loaded := tm.PutIfAbsent(1, "one")
	if loaded {
		t.Error("PutIfAbsent插入新键时，loaded应为false")
	}
	if existing != "" {
		t.Errorf("PutIfAbsent插入新键时，existing应为空字符串，实际为'%s'", existing)
	}

	// 测试插入已存在的键
	existing, loaded = tm.PutIfAbsent(1, "ONE")
	if !loaded {
		t.Error("PutIfAbsent插入已存在的键时，loaded应为true")
	}
	if existing != "one" {
		t.Errorf("PutIfAbsent插入已存在的键时，existing应为'one'，实际为'%s'", existing)
	}

	// 验证值未被更新
	val := tm.Get(1)
	if val != "one" {
		t.Errorf("PutIfAbsent不应更新已存在的值，期望'one'，实际为'%s'", val)
	}
}

// TestTreeMapGetOrDefault 测试GetOrDefault方法
func TestTreeMapGetOrDefault(t *testing.T) {
	tm := NewTreeMap[int, string](func(a, b int) bool {
		return a < b
	})

	tm.Put(1, "one")

	// 测试获取存在的键
	val := tm.GetOrDefault(1, "default")
	if val != "one" {
		t.Errorf("GetOrDefault获取存在的键应返回'one'，实际为'%s'", val)
	}

	// 测试获取不存在的键
	val = tm.GetOrDefault(2, "default")
	if val != "default" {
		t.Errorf("GetOrDefault获取不存在的键应返回'default'，实际为'%s'", val)
	}
}

// TestTreeMapFirstKey 测试FirstKey方法
func TestTreeMapFirstKey(t *testing.T) {
	tm := NewTreeMap[int, string](func(a, b int) bool {
		return a < b
	})

	// 测试空树
	_, ok := tm.FirstKey()
	if ok {
		t.Error("空树的FirstKey应返回false")
	}

	// 插入元素
	order := []int{5, 2, 8, 1, 9}
	for _, v := range order {
		tm.Put(v, fmt.Sprintf("value%d", v))
	}

	// 测试非空树
	key, ok := tm.FirstKey()
	if !ok {
		t.Error("非空树的FirstKey应返回true")
	}
	if key != 1 {
		t.Errorf("FirstKey应返回1，实际为%d", key)
	}
}

// TestTreeMapLastKey 测试LastKey方法
func TestTreeMapLastKey(t *testing.T) {
	tm := NewTreeMap[int, string](func(a, b int) bool {
		return a < b
	})

	// 测试空树
	_, ok := tm.LastKey()
	if ok {
		t.Error("空树的LastKey应返回false")
	}

	// 插入元素
	order := []int{5, 2, 8, 1, 9}
	for _, v := range order {
		tm.Put(v, fmt.Sprintf("value%d", v))
	}

	// 测试非空树
	key, ok := tm.LastKey()
	if !ok {
		t.Error("非空树的LastKey应返回true")
	}
	if key != 9 {
		t.Errorf("LastKey应返回9，实际为%d", key)
	}
}

// TestTreeMapRange 测试Range方法
func TestTreeMapRange(t *testing.T) {
	tm := NewTreeMap[int, string](func(a, b int) bool {
		return a < b
	})

	for i := 1; i <= 5; i++ {
		tm.Put(i, fmt.Sprintf("value%d", i))
	}

	// 测试完整遍历
	count := 0
	lastKey := 0
	tm.Range(func(key int, value string) bool {
		count++
		if key <= lastKey {
			t.Errorf("Range遍历的键应按升序排列，当前键%d <= 上一个键%d", key, lastKey)
		}
		lastKey = key

		expected := fmt.Sprintf("value%d", key)
		if value != expected {
			t.Errorf("键%d的值应为'%s'，实际为'%s'", key, expected, value)
		}
		return true
	})

	if count != 5 {
		t.Errorf("Range应遍历5个元素，实际遍历%d个", count)
	}

	// 测试提前终止
	count = 0
	tm.Range(func(key int, value string) bool {
		count++
		return key < 3 // 在键为3时停止
	})

	if count != 3 {
		t.Errorf("Range提前终止应遍历3个元素，实际遍历%d个", count)
	}
}

// TestTreeMapToString 测试ToString方法
func TestTreeMapToString(t *testing.T) {
	tm := NewTreeMap[int, string](func(a, b int) bool {
		return a < b
	})

	tm.Put(1, "one")
	tm.Put(2, "two")

	str := tm.ToString()
	if str == "" {
		t.Error("ToString不应返回空字符串")
	}

	// 验证是否为有效的JSON
	var m map[int]string
	err := json.Unmarshal([]byte(str), &m)
	if err != nil {
		t.Errorf("ToString返回的字符串应为有效的JSON: %v", err)
	}

	if m[1] != "one" || m[2] != "two" {
		t.Error("ToString返回的JSON内容不正确")
	}
}

// TestTreeMapMarshalJSON 测试JSON序列化
func TestTreeMapMarshalJSON(t *testing.T) {
	tm := NewTreeMap[int, string](func(a, b int) bool {
		return a < b
	})

	tm.Put(1, "one")
	tm.Put(2, "two")
	tm.Put(3, "three")

	data, err := json.Marshal(tm)
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	var m map[int]string
	err = json.Unmarshal(data, &m)
	if err != nil {
		t.Fatalf("JSON反序列化失败: %v", err)
	}

	if len(m) != 3 {
		t.Errorf("序列化后的map大小应为3，实际为%d", len(m))
	}

	if m[1] != "one" || m[2] != "two" || m[3] != "three" {
		t.Error("序列化后的数据不正确")
	}
}

// TestTreeMapUnmarshalJSON 测试JSON反序列化
func TestTreeMapUnmarshalJSON(t *testing.T) {
	jsonStr := `{"1":"one","2":"two","3":"three"}`

	tm := NewTreeMap[string, string](func(a, b string) bool {
		return a < b
	})

	err := json.Unmarshal([]byte(jsonStr), tm)
	if err != nil {
		t.Fatalf("JSON反序列化失败: %v", err)
	}

	if tm.Size() != 3 {
		t.Errorf("反序列化后的TreeMap大小应为3，实际为%d", tm.Size())
	}

	if tm.Get("1") != "one" || tm.Get("2") != "two" || tm.Get("3") != "three" {
		t.Error("反序列化后的数据不正确")
	}

	// 验证键已排序
	keys := tm.Keys()
	expected := []string{"1", "2", "3"}
	for i, key := range keys {
		if key != expected[i] {
			t.Errorf("索引%d处的键应为'%s'，实际为'%s'", i, expected[i], key)
		}
	}
}

// TestTreeMapMarshalBSON 测试BSON序列化
func TestTreeMapMarshalBSON(t *testing.T) {
	tm := NewTreeMap[string, int](func(a, b string) bool {
		return a < b
	})

	tm.Put("one", 1)
	tm.Put("two", 2)
	tm.Put("three", 3)

	data, err := bson.Marshal(tm)
	if err != nil {
		t.Fatalf("BSON序列化失败: %v", err)
	}

	var m map[string]int
	err = bson.Unmarshal(data, &m)
	if err != nil {
		t.Fatalf("BSON反序列化失败: %v", err)
	}

	if len(m) != 3 {
		t.Errorf("序列化后的map大小应为3，实际为%d", len(m))
	}

	if m["one"] != 1 || m["two"] != 2 || m["three"] != 3 {
		t.Error("BSON序列化后的数据不正确")
	}
}

// TestTreeMapUnmarshalBSON 测试BSON反序列化
func TestTreeMapUnmarshalBSON(t *testing.T) {
	originalMap := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	data, err := bson.Marshal(originalMap)
	if err != nil {
		t.Fatalf("BSON序列化原始map失败: %v", err)
	}

	tm := NewTreeMap[string, int](func(a, b string) bool {
		return a < b
	})

	err = bson.Unmarshal(data, tm)
	if err != nil {
		t.Fatalf("BSON反序列化失败: %v", err)
	}

	if tm.Size() != 3 {
		t.Errorf("反序列化后的TreeMap大小应为3，实际为%d", tm.Size())
	}

	if tm.Get("one") != 1 || tm.Get("two") != 2 || tm.Get("three") != 3 {
		t.Error("BSON反序列化后的数据不正确")
	}
}

// TestTreeMapConcurrency 测试并发安全性
func TestTreeMapConcurrency(t *testing.T) {
	tm := NewTreeMap[int, int](func(a, b int) bool {
		return a < b
	})

	const goroutines = 10
	const operations = 100

	var wg sync.WaitGroup
	wg.Add(goroutines * 3) // 3种操作：Put, Get, Remove

	// 并发Put
	for i := 0; i < goroutines; i++ {
		go func(start int) {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				tm.Put(start*operations+j, j)
			}
		}(i)
	}

	// 并发Get
	for i := 0; i < goroutines; i++ {
		go func(start int) {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				tm.Get(start*operations + j)
			}
		}(i)
	}

	// 并发Remove
	for i := 0; i < goroutines; i++ {
		go func(start int) {
			defer wg.Done()
			for j := 0; j < operations/2; j++ {
				tm.Remove(start*operations + j)
			}
		}(i)
	}

	wg.Wait()

	// 验证TreeMap仍然有效
	size := tm.Size()
	keys := tm.Keys()

	if size != len(keys) {
		t.Errorf("Size()返回%d，但Keys()返回%d个键", size, len(keys))
	}

	// 验证键已排序
	if !sort.IntsAreSorted(keys) {
		t.Error("并发操作后，Keys()返回的键应保持有序")
	}
}

// TestTreeMapConcurrentRange 测试并发Range
func TestTreeMapConcurrentRange(t *testing.T) {
	tm := NewTreeMap[int, int](func(a, b int) bool {
		return a < b
	})

	// 初始化数据
	for i := 0; i < 100; i++ {
		tm.Put(i, i)
	}

	var wg sync.WaitGroup
	wg.Add(3)

	// 并发Range
	go func() {
		defer wg.Done()
		tm.Range(func(key, value int) bool {
			return true
		})
	}()

	// 并发修改
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			tm.Put(i, i*2)
		}
	}()

	// 并发删除
	go func() {
		defer wg.Done()
		for i := 50; i < 75; i++ {
			tm.Remove(i)
		}
	}()

	wg.Wait()

	// 验证TreeMap仍然有效
	keys := tm.Keys()
	if !sort.IntsAreSorted(keys) {
		t.Error("并发Range后，Keys()返回的键应保持有序")
	}
}

// TestTreeMapStringKeys 测试字符串键的TreeMap
func TestTreeMapStringKeys(t *testing.T) {
	tm := NewTreeMap[string, int](func(a, b string) bool {
		return a < b
	})

	words := []string{"banana", "apple", "cherry", "date", "elderberry"}
	for i, word := range words {
		tm.Put(word, i)
	}

	// 验证排序
	keys := tm.Keys()
	expected := []string{"apple", "banana", "cherry", "date", "elderberry"}

	if len(keys) != len(expected) {
		t.Fatalf("Keys()返回的长度不正确，期望%d，实际%d", len(expected), len(keys))
	}

	for i, key := range keys {
		if key != expected[i] {
			t.Errorf("索引%d处的键应为'%s'，实际为'%s'", i, expected[i], key)
		}
	}
}

// TestTreeMapComplexStruct 测试复杂结构体作为值
func TestTreeMapComplexStruct(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	tm := NewTreeMap[int, Person](func(a, b int) bool {
		return a < b
	})

	tm.Put(1, Person{"Alice", 30})
	tm.Put(2, Person{"Bob", 25})
	tm.Put(3, Person{"Charlie", 35})

	person := tm.Get(2)
	if person.Name != "Bob" || person.Age != 25 {
		t.Errorf("Get(2)应返回Bob(25)，实际为%s(%d)", person.Name, person.Age)
	}

	values := tm.Values()
	if len(values) != 3 {
		t.Errorf("Values()应返回3个元素，实际为%d", len(values))
	}
}

// TestTreeMapEmptyOperations 测试空TreeMap的各种操作
func TestTreeMapEmptyOperations(t *testing.T) {
	tm := NewTreeMap[int, string](func(a, b int) bool {
		return a < b
	})

	// 测试空树的Get
	val := tm.Get(1)
	if val != "" {
		t.Errorf("空树的Get应返回零值，实际为'%s'", val)
	}

	// 测试空树的Remove
	tm.Remove(1) // 不应panic

	// 测试空树的ContainsKey
	if tm.ContainsKey(1) {
		t.Error("空树的ContainsKey应返回false")
	}

	// 测试空树的Keys
	keys := tm.Keys()
	if len(keys) != 0 {
		t.Errorf("空树的Keys()应返回空切片，实际长度为%d", len(keys))
	}

	// 测试空树的Values
	values := tm.Values()
	if len(values) != 0 {
		t.Errorf("空树的Values()应返回空切片，实际长度为%d", len(values))
	}

	// 测试空树的Range
	count := 0
	tm.Range(func(key int, value string) bool {
		count++
		return true
	})
	if count != 0 {
		t.Errorf("空树的Range不应遍历任何元素，实际遍历%d个", count)
	}

	// 测试空树的ToString
	str := tm.ToString()
	if str != "{}" {
		t.Errorf("空树的ToString应返回'{}'，实际为'%s'", str)
	}
}

// TestTreeMapLargeDataset 测试大数据集
func TestTreeMapLargeDataset(t *testing.T) {
	tm := NewTreeMap[int, int](func(a, b int) bool {
		return a < b
	})

	const size = 10000

	// 插入大量数据
	for i := 0; i < size; i++ {
		tm.Put(i, i*2)
	}

	if tm.Size() != size {
		t.Errorf("插入%d个元素后，大小应为%d，实际为%d", size, size, tm.Size())
	}

	// 验证所有元素
	for i := 0; i < size; i++ {
		val := tm.Get(i)
		if val != i*2 {
			t.Errorf("Get(%d)应返回%d，实际为%d", i, i*2, val)
		}
	}

	// 验证排序
	keys := tm.Keys()
	if !sort.IntsAreSorted(keys) {
		t.Error("大数据集的Keys()应保持有序")
	}

	// 删除一半数据
	for i := 0; i < size/2; i++ {
		tm.Remove(i)
	}

	if tm.Size() != size/2 {
		t.Errorf("删除一半元素后，大小应为%d，实际为%d", size/2, tm.Size())
	}
}

// TestTreeMapDescendingOrder 测试降序比较函数
func TestTreeMapDescendingOrder(t *testing.T) {
	// 使用降序比较函数
	tm := NewTreeMap[int, string](func(a, b int) bool {
		return a > b
	})

	order := []int{5, 2, 8, 1, 9}
	for _, v := range order {
		tm.Put(v, fmt.Sprintf("value%d", v))
	}

	// 验证键按降序排列
	keys := tm.Keys()
	expected := []int{9, 8, 5, 2, 1}

	if len(keys) != len(expected) {
		t.Fatalf("Keys()返回的长度不正确，期望%d，实际%d", len(expected), len(keys))
	}

	for i, key := range keys {
		if key != expected[i] {
			t.Errorf("索引%d处的键应为%d，实际为%d", i, expected[i], key)
		}
	}
}

// TestTreeMapRedBlackTreeProperties 测试红黑树性质（间接测试）
func TestTreeMapRedBlackTreeProperties(t *testing.T) {
	tm := NewTreeMap[int, string](func(a, b int) bool {
		return a < b
	})

	// 插入有序数据（最坏情况）
	for i := 1; i <= 100; i++ {
		tm.Put(i, fmt.Sprintf("value%d", i))
	}

	// 如果是普通二叉搜索树，有序插入会导致退化成链表
	// 红黑树应该保持平衡，所以所有操作应该很快完成
	// 这里我们验证树的功能性

	// 验证所有元素都能正确获取
	for i := 1; i <= 100; i++ {
		val := tm.Get(i)
		expected := fmt.Sprintf("value%d", i)
		if val != expected {
			t.Errorf("Get(%d)应返回'%s'，实际为'%s'", i, expected, val)
		}
	}

	// 验证键有序
	keys := tm.Keys()
	if !sort.IntsAreSorted(keys) {
		t.Error("Keys()应返回有序的键")
	}
}

// BenchmarkTreeMapPut 基准测试Put操作
func BenchmarkTreeMapPut(b *testing.B) {
	tm := NewTreeMap[int, int](func(a, b int) bool {
		return a < b
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tm.Put(i, i)
	}
}

// BenchmarkTreeMapGet 基准测试Get操作
func BenchmarkTreeMapGet(b *testing.B) {
	tm := NewTreeMap[int, int](func(a, b int) bool {
		return a < b
	})

	for i := 0; i < 1000; i++ {
		tm.Put(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tm.Get(i % 1000)
	}
}

// BenchmarkTreeMapRemove 基准测试Remove操作
func BenchmarkTreeMapRemove(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		tm := NewTreeMap[int, int](func(a, b int) bool {
			return a < b
		})
		for j := 0; j < 1000; j++ {
			tm.Put(j, j)
		}

		b.StartTimer()
		tm.Remove(i % 1000)
		b.StopTimer()
	}
}
