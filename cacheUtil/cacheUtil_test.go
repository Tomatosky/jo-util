package cacheUtil

import (
	"sync"
	"testing"
	"time"
)

// TestNew 测试基本的缓存创建
func TestNew(t *testing.T) {
	cache1 := New[string, int](5 * time.Second)
	if cache1 == nil {
		t.Fatal("New() returned nil")
	}
	if cache1.cache == nil {
		t.Fatal("cache1.cache is nil")
	}
	if cache1.expiration != 5*time.Second {
		t.Errorf("Expected expiration 5s, got %v", cache1.expiration)
	}
	if cache1.accessExpire {
		t.Error("accessExpire should be false for New()")
	}
}

// TestNewAccessExpire 测试访问过期模式的缓存创建
func TestNewAccessExpire(t *testing.T) {
	cache1 := NewAccessExpire[string, int](5 * time.Second)
	if cache1 == nil {
		t.Fatal("NewAccessExpire() returned nil")
	}
	if !cache1.accessExpire {
		t.Error("accessExpire should be true for NewAccessExpire()")
	}
}

// TestSetAndGet 测试基本的设置和获取功能
func TestSetAndGet(t *testing.T) {
	cache1 := New[string, string](5 * time.Second)

	// 测试设置和获取
	cache1.Set("key1", "value1")
	val, found := cache1.Get("key1")
	if !found {
		t.Error("Expected to find key1")
	}
	if val != "value1" {
		t.Errorf("Expected value1, got %s", val)
	}

	// 测试不存在的键
	_, found = cache1.Get("nonexistent")
	if found {
		t.Error("Should not find nonexistent key")
	}
}

// TestSetWithCustomExpiration 测试自定义过期时间
func TestSetWithCustomExpiration(t *testing.T) {
	cache1 := New[string, string](5 * time.Second)

	// 使用自定义的短过期时间
	cache1.Set("key1", "value1", 100*time.Millisecond)

	// 立即获取应该存在
	val, found := cache1.Get("key1")
	if !found || val != "value1" {
		t.Error("Key should exist immediately after setting")
	}

	// 等待过期
	time.Sleep(150 * time.Millisecond)
	_, found = cache1.Get("key1")
	if found {
		t.Error("Key should have expired")
	}
}

// TestSetIfAbsent 测试 SetIfAbsent 功能
func TestSetIfAbsent(t *testing.T) {
	cache1 := New[string, int](5 * time.Second)

	// 第一次设置应该成功
	success := cache1.SetIfAbsent("key1", 100)
	if !success {
		t.Error("First SetIfAbsent should succeed")
	}

	val, found := cache1.Get("key1")
	if !found || val != 100 {
		t.Error("Value should be set to 100")
	}

	// 第二次设置应该失败
	success = cache1.SetIfAbsent("key1", 200)
	if success {
		t.Error("Second SetIfAbsent should fail")
	}

	val, found = cache1.Get("key1")
	if !found || val != 100 {
		t.Error("Value should still be 100")
	}
}

// TestSetIfAbsentWithExpiredKey 测试 SetIfAbsent 在键过期后的行为
func TestSetIfAbsentWithExpiredKey(t *testing.T) {
	cache1 := New[string, int](100 * time.Millisecond)

	// 设置一个会过期的键
	cache1.Set("key1", 100)

	// 等待过期
	time.Sleep(150 * time.Millisecond)

	// 过期后应该可以再次设置
	success := cache1.SetIfAbsent("key1", 200)
	if !success {
		t.Error("SetIfAbsent should succeed on expired key")
	}

	val, found := cache1.Get("key1")
	if !found || val != 200 {
		t.Error("Value should be updated to 200")
	}
}

// TestGetWithExpiration 测试获取带过期时间的功能
func TestGetWithExpiration(t *testing.T) {
	cache1 := New[string, string](5 * time.Second)

	cache1.Set("key1", "value1")
	val, found, expTime := cache1.GetWithExpiration("key1")

	if !found {
		t.Error("Key should be found")
	}
	if val != "value1" {
		t.Errorf("Expected value1, got %s", val)
	}
	if expTime.IsZero() {
		t.Error("Expiration time should not be zero")
	}

	// 检查过期时间是否在合理范围内（大约5秒后）
	expectedTime := time.Now().Add(5 * time.Second)
	diff := expTime.Sub(expectedTime).Abs()
	if diff > time.Second {
		t.Errorf("Expiration time difference too large: %v", diff)
	}
}

// TestAccessExpire 测试访问过期模式
func TestAccessExpire(t *testing.T) {
	cache1 := NewAccessExpire[string, string](200 * time.Millisecond)

	cache1.Set("key1", "value1")

	// 在过期前多次访问，每次访问都会重置过期时间
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		val, found := cache1.Get("key1")
		if !found || val != "value1" {
			t.Errorf("Iteration %d: Key should still be valid due to access", i)
		}
	}

	// 停止访问，等待过期
	time.Sleep(250 * time.Millisecond)
	_, found := cache1.Get("key1")
	if found {
		t.Error("Key should have expired after no access")
	}
}

// TestDelete 测试删除功能
func TestDelete(t *testing.T) {
	cache1 := New[string, string](5 * time.Second)

	cache1.Set("key1", "value1")
	cache1.Set("key2", "value2")

	// 验证键存在
	_, found := cache1.Get("key1")
	if !found {
		t.Error("key1 should exist")
	}

	// 删除键
	cache1.Delete("key1")

	// 验证键已删除
	_, found = cache1.Get("key1")
	if found {
		t.Error("key1 should be deleted")
	}

	// 验证其他键不受影响
	val, found := cache1.Get("key2")
	if !found || val != "value2" {
		t.Error("key2 should still exist")
	}
}

// TestItems 测试获取所有项
func TestItems(t *testing.T) {
	cache1 := New[string, int](5 * time.Second)

	cache1.Set("key1", 1)
	cache1.Set("key2", 2)
	cache1.Set("key3", 3)

	items := cache1.Items()
	if len(items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(items))
	}

	if items["key1"] != 1 || items["key2"] != 2 || items["key3"] != 3 {
		t.Error("Items contain incorrect values")
	}
}

// TestItemsWithExpiredKeys 测试 Items 方法过滤过期键
func TestItemsWithExpiredKeys(t *testing.T) {
	cache1 := New[string, int](100 * time.Millisecond)

	cache1.Set("key1", 1)
	cache1.Set("key2", 2)

	// 等待部分键过期
	time.Sleep(150 * time.Millisecond)

	// 添加新键
	cache1.Set("key3", 3)

	items := cache1.Items()
	// key1 和 key2 应该过期，只剩 key3
	if len(items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(items))
	}

	if items["key3"] != 3 {
		t.Error("key3 should be the only valid item")
	}
}

// TestFlush 测试清空缓存
func TestFlush(t *testing.T) {
	cache1 := New[string, int](5 * time.Second)

	cache1.Set("key1", 1)
	cache1.Set("key2", 2)
	cache1.Set("key3", 3)

	// 清空缓存
	cache1.Flush()

	// 验证所有键都被删除
	items := cache1.Items()
	if len(items) != 0 {
		t.Errorf("Expected 0 items after flush, got %d", len(items))
	}

	_, found := cache1.Get("key1")
	if found {
		t.Error("key1 should not exist after flush")
	}
}

// TestExpiration 测试过期机制
func TestExpiration(t *testing.T) {
	cache1 := New[string, string](200 * time.Millisecond)

	cache1.Set("key1", "value1")

	// 立即获取应该存在
	val, found := cache1.Get("key1")
	if !found || val != "value1" {
		t.Error("Key should exist immediately")
	}

	// 等待过期
	time.Sleep(250 * time.Millisecond)

	// 获取应该失败
	_, found = cache1.Get("key1")
	if found {
		t.Error("Key should have expired")
	}
}

// TestJanitorCleanup 测试自动清理功能
func TestJanitorCleanup(t *testing.T) {
	cache1 := New[string, string](100 * time.Millisecond)

	// 设置多个会过期的键
	for i := 0; i < 10; i++ {
		cache1.Set(string(rune('a'+i)), "value")
	}

	// 等待过期
	time.Sleep(150 * time.Millisecond)

	// 检查内部 map 是否被清理
	cache1.mu.RLock()
	itemCount := len(cache1.Items())
	cache1.mu.RUnlock()

	if itemCount != 0 {
		t.Errorf("Expected janitor to clean up expired items, found %d items", itemCount)
	}
}

// TestConcurrentAccess 测试并发访问
func TestConcurrentAccess(t *testing.T) {
	cache1 := New[int, int](5 * time.Second)
	var wg sync.WaitGroup

	// 并发写入
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			cache1.Set(val, val*2)
		}(i)
	}

	// 并发读取
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			cache1.Get(key)
		}(i)
	}

	// 并发删除
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			cache1.Delete(key)
		}(i)
	}

	wg.Wait()

	// 验证缓存仍然可用
	cache1.Set(999, 999)
	val, found := cache1.Get(999)
	if !found || val != 999 {
		t.Error("Cache should still be functional after concurrent access")
	}
}

// TestConcurrentSetIfAbsent 测试并发 SetIfAbsent
func TestConcurrentSetIfAbsent(t *testing.T) {
	cache1 := New[string, int](5 * time.Second)
	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	// 多个 goroutine 尝试设置同一个键
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			if cache1.SetIfAbsent("key", val) {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	// 只应该有一个成功
	if successCount != 1 {
		t.Errorf("Expected 1 successful SetIfAbsent, got %d", successCount)
	}

	// 验证值被正确设置
	_, found := cache1.Get("key")
	if !found {
		t.Error("Key should exist")
	}
}

// TestZeroExpiration 测试不过期的缓存
func TestZeroExpiration(t *testing.T) {
	cache1 := New[string, string](0)

	cache1.Set("key1", "value1")

	// 等待一段时间
	time.Sleep(100 * time.Millisecond)

	// 键应该仍然存在
	val, found := cache1.Get("key1")
	if !found || val != "value1" {
		t.Error("Key should never expire with zero expiration")
	}
}

// TestDifferentTypes 测试不同的数据类型
func TestDifferentTypes(t *testing.T) {
	// 测试 int 键和 string 值
	cache1 := New[int, string](5 * time.Second)
	cache1.Set(1, "one")
	val1, found := cache1.Get(1)
	if !found || val1 != "one" {
		t.Error("int/string cache failed")
	}

	// 测试 string 键和结构体值
	type Person struct {
		Name string
		Age  int
	}
	cache2 := New[string, Person](5 * time.Second)
	cache2.Set("john", Person{Name: "John", Age: 30})
	val2, found := cache2.Get("john")
	if !found || val2.Name != "John" || val2.Age != 30 {
		t.Error("string/struct cache failed")
	}

	// 测试 string 键和 slice 值
	cache3 := New[string, []int](5 * time.Second)
	cache3.Set("numbers", []int{1, 2, 3, 4, 5})
	val3, found := cache3.Get("numbers")
	if !found || len(val3) != 5 || val3[0] != 1 {
		t.Error("string/slice cache failed")
	}
}

// TestUpdateExistingKey 测试更新已存在的键
func TestUpdateExistingKey(t *testing.T) {
	cache1 := New[string, int](5 * time.Second)

	cache1.Set("key1", 100)
	val, _ := cache1.Get("key1")
	if val != 100 {
		t.Error("Initial value should be 100")
	}

	// 更新值
	cache1.Set("key1", 200)
	val, _ = cache1.Get("key1")
	if val != 200 {
		t.Error("Updated value should be 200")
	}
}

// TestMultipleExpirationsInSameCache 测试同一缓存中的不同过期时间
func TestMultipleExpirationsInSameCache(t *testing.T) {
	cache1 := New[string, string](5 * time.Second)

	// 使用默认过期时间
	cache1.Set("key1", "value1")

	// 使用短过期时间
	cache1.Set("key2", "value2", 100*time.Millisecond)

	// 使用长过期时间
	cache1.Set("key3", "value3", 10*time.Second)

	// 等待一段时间
	time.Sleep(150 * time.Millisecond)

	// key1 应该仍然存在（5秒过期）
	_, found := cache1.Get("key1")
	if !found {
		t.Error("key1 should still exist")
	}

	// key2 应该过期（100ms过期）
	_, found = cache1.Get("key2")
	if found {
		t.Error("key2 should have expired")
	}

	// key3 应该仍然存在（10秒过期）
	_, found = cache1.Get("key3")
	if !found {
		t.Error("key3 should still exist")
	}
}

// TestEmptyCache 测试空缓存操作
func TestEmptyCache(t *testing.T) {
	cache1 := New[string, string](5 * time.Second)

	// 从空缓存获取
	_, found := cache1.Get("nonexistent")
	if found {
		t.Error("Should not find anything in empty cache")
	}

	// 从空缓存删除
	cache1.Delete("nonexistent") // 不应该 panic

	// 获取空缓存的所有项
	items := cache1.Items()
	if len(items) != 0 {
		t.Error("Empty cache should return empty items map")
	}

	// 清空空缓存
	cache1.Flush() // 不应该 panic
}

// BenchmarkSet 基准测试 Set 操作
func BenchmarkSet(b *testing.B) {
	cache1 := New[int, int](5 * time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache1.Set(i, i*2)
	}
}

// BenchmarkGet 基准测试 Get 操作
func BenchmarkGet(b *testing.B) {
	cache1 := New[int, int](5 * time.Second)
	for i := 0; i < 1000; i++ {
		cache1.Set(i, i*2)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache1.Get(i % 1000)
	}
}

// BenchmarkSetIfAbsent 基准测试 SetIfAbsent 操作
func BenchmarkSetIfAbsent(b *testing.B) {
	cache1 := New[int, int](5 * time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache1.SetIfAbsent(i%100, i)
	}
}

// BenchmarkConcurrentSet 基准测试并发 Set 操作
func BenchmarkConcurrentSet(b *testing.B) {
	cache1 := New[int, int](5 * time.Second)
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache1.Set(i, i*2)
			i++
		}
	})
}

// BenchmarkConcurrentGet 基准测试并发 Get 操作
func BenchmarkConcurrentGet(b *testing.B) {
	cache1 := New[int, int](5 * time.Second)
	for i := 0; i < 1000; i++ {
		cache1.Set(i, i*2)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache1.Get(i % 1000)
			i++
		}
	})
}
