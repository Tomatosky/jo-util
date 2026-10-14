package cacheUtil

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func newTestCache[K comparable, V any](t *testing.T, expiration time.Duration) *Cache[K, V] {
	t.Helper()
	c := New[K, V](expiration)
	if c.janitor != nil {
		t.Cleanup(func() {
			stopJanitor(c.janitor)
		})
	}
	return c
}

func newAccessTestCache[K comparable, V any](t *testing.T, expiration time.Duration) *Cache[K, V] {
	t.Helper()
	c := NewAccessExpire[K, V](expiration)
	if c.janitor != nil {
		t.Cleanup(func() {
			stopJanitor(c.janitor)
		})
	}
	return c
}

func expireTestKey[K comparable, V any](c *Cache[K, V], key K) {
	c.mu.Lock()
	c.items[key].Expiration = time.Now().Add(-time.Nanosecond).UnixNano()
	c.mu.Unlock()
}

// TestNew 测试基本的缓存创建
func TestNew(t *testing.T) {
	cache1 := newTestCache[string, int](t, 5*time.Second)
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
	cache1 := newAccessTestCache[string, int](t, 5*time.Second)
	if cache1 == nil {
		t.Fatal("NewAccessExpire() returned nil")
	}
	if !cache1.accessExpire {
		t.Error("accessExpire should be true for NewAccessExpire()")
	}
}

func TestStopJanitorIsIdempotent(t *testing.T) {
	j := &janitor[string, int]{stop: make(chan struct{})}
	stopJanitor(j)
	stopJanitor(j)

	select {
	case <-j.stop:
	default:
		t.Fatal("stopJanitor() did not close the stop channel")
	}
}

// TestSetAndGet 测试基本的设置和获取功能
func TestSetAndGet(t *testing.T) {
	cache1 := newTestCache[string, string](t, 5*time.Second)

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
	cache1 := newTestCache[string, string](t, 5*time.Second)

	// 使用自定义的短过期时间
	beforeSet := time.Now()
	cache1.Set("key1", "value1", time.Minute)
	afterSet := time.Now()
	gotExpiration := time.Unix(0, cache1.items["key1"].Expiration)
	if gotExpiration.Before(beforeSet.Add(time.Minute)) || gotExpiration.After(afterSet.Add(time.Minute)) {
		t.Fatalf("custom expiration = %v, want within [%v, %v]", gotExpiration, beforeSet.Add(time.Minute), afterSet.Add(time.Minute))
	}

	// 立即获取应该存在
	val, found := cache1.Get("key1")
	if !found || val != "value1" {
		t.Error("Key should exist immediately after setting")
	}

	// 直接将条目推进到过期状态，避免依赖调度时序。
	expireTestKey(cache1, "key1")
	_, found = cache1.Get("key1")
	if found {
		t.Error("Key should have expired")
	}
}

// TestSetIfAbsent 测试 SetIfAbsent 功能
func TestSetIfAbsent(t *testing.T) {
	cache1 := newTestCache[string, int](t, 5*time.Second)

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
	cache1 := newTestCache[string, int](t, time.Minute)

	// 设置一个会过期的键
	cache1.Set("key1", 100)

	expireTestKey(cache1, "key1")

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
	cache1 := newTestCache[string, string](t, 5*time.Second)

	beforeSet := time.Now()
	cache1.Set("key1", "value1")
	afterSet := time.Now()
	val, found, expTime := cache1.GetWithExpiration("key1")

	if !found {
		t.Error("Key should be found")
	}
	if val != "value1" {
		t.Errorf("Expected value1, got %s", val)
	}
	if expTime.Before(beforeSet.Add(5*time.Second)) || expTime.After(afterSet.Add(5*time.Second)) {
		t.Errorf("expiration = %v, want within [%v, %v]", expTime, beforeSet.Add(5*time.Second), afterSet.Add(5*time.Second))
	}
}

// TestAccessExpire 测试访问过期模式
func TestAccessExpire(t *testing.T) {
	const expiration = time.Minute
	cache1 := newAccessTestCache[string, string](t, expiration)

	cache1.Set("key1", "value1")
	beforeGet := time.Now()
	val, found, refreshedExpiration := cache1.GetWithExpiration("key1")
	afterGet := time.Now()
	if !found || val != "value1" {
		t.Fatalf("GetWithExpiration() = (%q, %v), want (%q, true)", val, found, "value1")
	}
	if refreshedExpiration.Before(beforeGet.Add(expiration)) || refreshedExpiration.After(afterGet.Add(expiration)) {
		t.Errorf("refreshed expiration = %v, want within [%v, %v]", refreshedExpiration, beforeGet.Add(expiration), afterGet.Add(expiration))
	}

	expireTestKey(cache1, "key1")
	_, found = cache1.Get("key1")
	if found {
		t.Error("Key should have expired after no access")
	}
}

func TestConcurrentAccessExpire(t *testing.T) {
	const (
		expiration = time.Minute
		readers    = 100
	)
	cache1 := newAccessTestCache[string, string](t, expiration)
	cache1.Set("key", "value")

	type result struct {
		value      string
		found      bool
		expiration time.Time
	}
	results := make(chan result, readers)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < readers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			value, found, refreshedExpiration := cache1.GetWithExpiration("key")
			results <- result{value: value, found: found, expiration: refreshedExpiration}
		}()
	}
	before := time.Now()
	close(start)
	wg.Wait()
	after := time.Now()
	close(results)

	for got := range results {
		if !got.found || got.value != "value" {
			t.Errorf("GetWithExpiration() = (%q, %v), want (%q, true)", got.value, got.found, "value")
		}
		if got.expiration.Before(before.Add(expiration)) || got.expiration.After(after.Add(expiration)) {
			t.Errorf("refreshed expiration = %v, want within [%v, %v]", got.expiration, before.Add(expiration), after.Add(expiration))
		}
	}
}

// TestDelete 测试删除功能
func TestDelete(t *testing.T) {
	cache1 := newTestCache[string, string](t, 5*time.Second)

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
	cache1 := newTestCache[string, int](t, 5*time.Second)

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
	cache1 := newTestCache[string, int](t, time.Minute)

	cache1.Set("key1", 1)
	cache1.Set("key2", 2)

	expireTestKey(cache1, "key1")
	expireTestKey(cache1, "key2")

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
	cache1 := newTestCache[string, int](t, 5*time.Second)

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
	cache1 := newTestCache[string, string](t, time.Minute)

	cache1.Set("key1", "value1")

	// 立即获取应该存在
	val, found := cache1.Get("key1")
	if !found || val != "value1" {
		t.Error("Key should exist immediately")
	}

	expireTestKey(cache1, "key1")

	// 获取应该失败
	_, found = cache1.Get("key1")
	if found {
		t.Error("Key should have expired")
	}
}

// TestJanitorCleanup 测试自动清理功能
func TestDeleteExpired(t *testing.T) {
	cache1 := newTestCache[string, string](t, time.Minute)

	// 设置多个会过期的键
	for i := 0; i < 10; i++ {
		cache1.Set(string(rune('a'+i)), "value")
	}

	cache1.mu.Lock()
	for _, item := range cache1.items {
		item.Expiration = time.Now().Add(-time.Nanosecond).UnixNano()
	}
	cache1.mu.Unlock()
	cache1.deleteExpired()

	cache1.mu.RLock()
	itemCount := len(cache1.items)
	cache1.mu.RUnlock()

	if itemCount != 0 {
		t.Errorf("deleteExpired() left %d expired items in the backing map", itemCount)
	}
}

// TestConcurrentAccess 测试并发访问
func TestConcurrentAccess(t *testing.T) {
	cache1 := newTestCache[int, int](t, 5*time.Second)
	var wg sync.WaitGroup

	// 并发写入
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			cache1.Set(val, val*2)
		}(i)
	}

	wg.Wait()
	for i := 0; i < 100; i++ {
		if got, found := cache1.Get(i); !found || got != i*2 {
			t.Fatalf("after concurrent writes Get(%d) = (%d, %v), want (%d, true)", i, got, found, i*2)
		}
	}

	// 并发读取已知存在的键。
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			if got, found := cache1.Get(key); !found || got != key*2 {
				t.Errorf("concurrent Get(%d) = (%d, %v), want (%d, true)", key, got, found, key*2)
			}
		}(i)
	}
	wg.Wait()

	// 并发删除
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			cache1.Delete(key)
		}(i)
	}

	wg.Wait()
	for i := 0; i < 100; i++ {
		got, found := cache1.Get(i)
		if i < 50 && found {
			t.Errorf("Get(%d) found deleted value %d", i, got)
		}
		if i >= 50 && (!found || got != i*2) {
			t.Errorf("Get(%d) = (%d, %v), want (%d, true)", i, got, found, i*2)
		}
	}
}

// TestConcurrentSetIfAbsent 测试并发 SetIfAbsent
func TestConcurrentSetIfAbsent(t *testing.T) {
	cache1 := newTestCache[string, int](t, 5*time.Second)
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
	value, found := cache1.Get("key")
	if !found || value < 0 || value >= 10 {
		t.Errorf("Get(key) = (%d, %v), want one submitted value in [0, 10)", value, found)
	}
}

// TestZeroExpiration 测试不过期的缓存
func TestZeroExpiration(t *testing.T) {
	cache1 := New[string, string](0)

	cache1.Set("key1", "value1")
	if expiration := cache1.items["key1"].Expiration; expiration != 0 {
		t.Fatalf("zero-expiration item stored expiration %d, want 0", expiration)
	}
	cache1.deleteExpired()

	// 键应该仍然存在
	val, found := cache1.Get("key1")
	if !found || val != "value1" {
		t.Error("Key should never expire with zero expiration")
	}
}

func TestCustomExpirationWithZeroDefault(t *testing.T) {
	cache1 := newTestCache[string, int](t, 0)
	cache1.Set("get", 1, time.Minute)
	cache1.Set("items", 2, time.Minute)
	cache1.Set("janitor", 3, time.Minute)
	cache1.Set("permanent", 4)

	for _, key := range []string{"get", "items", "janitor"} {
		expireTestKey(cache1, key)
	}
	if got, found := cache1.Get("get"); found || got != 0 {
		t.Errorf("Get(expired custom item) = (%d, %v), want (0, false)", got, found)
	}
	if got := cache1.Items(); !reflect.DeepEqual(got, map[string]int{"permanent": 4}) {
		t.Errorf("Items() = %v, want only permanent item", got)
	}

	cache1.deleteExpired()
	cache1.mu.RLock()
	_, getExists := cache1.items["get"]
	_, itemsExists := cache1.items["items"]
	_, janitorExists := cache1.items["janitor"]
	permanent := cache1.items["permanent"]
	cache1.mu.RUnlock()
	if getExists || itemsExists || janitorExists {
		t.Errorf("deleteExpired left custom-expired items: get=%v items=%v janitor=%v", getExists, itemsExists, janitorExists)
	}
	if permanent == nil || permanent.Expiration != 0 || permanent.Object != 4 {
		t.Errorf("permanent item after cleanup = %#v", permanent)
	}
}

func TestExplicitZeroExpirationOverridesDefault(t *testing.T) {
	cache1 := newTestCache[string, int](t, time.Minute)
	cache1.Set("permanent", 1, 0)

	value, found, expiration := cache1.GetWithExpiration("permanent")
	if !found || value != 1 || !expiration.IsZero() {
		t.Errorf("GetWithExpiration(permanent) = (%d, %v, %v), want (1, true, zero)", value, found, expiration)
	}
	cache1.deleteExpired()
	if value, found := cache1.Get("permanent"); !found || value != 1 {
		t.Errorf("Get(permanent) after cleanup = (%d, %v), want (1, true)", value, found)
	}
}

func TestGetWithExpirationExpired(t *testing.T) {
	cache1 := newTestCache[string, int](t, time.Minute)
	cache1.Set("expired", 1)
	expireTestKey(cache1, "expired")
	if value, found, expiration := cache1.GetWithExpiration("expired"); value != 0 || found || !expiration.IsZero() {
		t.Errorf("GetWithExpiration(expired) = (%d, %v, %v), want zero values", value, found, expiration)
	}
}

// TestDifferentTypes 测试不同的数据类型
func TestDifferentTypes(t *testing.T) {
	// 测试 int 键和 string 值
	cache1 := newTestCache[int, string](t, 5*time.Second)
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
	cache2 := newTestCache[string, Person](t, 5*time.Second)
	cache2.Set("john", Person{Name: "John", Age: 30})
	val2, found := cache2.Get("john")
	if !found || val2.Name != "John" || val2.Age != 30 {
		t.Error("string/struct cache failed")
	}

	// 测试 string 键和 slice 值
	cache3 := newTestCache[string, []int](t, 5*time.Second)
	cache3.Set("numbers", []int{1, 2, 3, 4, 5})
	val3, found := cache3.Get("numbers")
	if !found || !reflect.DeepEqual(val3, []int{1, 2, 3, 4, 5}) {
		t.Error("string/slice cache failed")
	}
}

// TestUpdateExistingKey 测试更新已存在的键
func TestUpdateExistingKey(t *testing.T) {
	cache1 := newTestCache[string, int](t, 5*time.Second)

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
	cache1 := newTestCache[string, string](t, 5*time.Second)

	// 使用默认过期时间
	cache1.Set("key1", "value1")

	// 使用短过期时间
	cache1.Set("key2", "value2", 100*time.Millisecond)

	// 使用长过期时间
	cache1.Set("key3", "value3", 10*time.Second)

	expireTestKey(cache1, "key2")

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
	cache1 := newTestCache[string, string](t, 5*time.Second)

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
	cache1 := New[int, int](0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache1.Set(i, i*2)
	}
}

// BenchmarkGet 基准测试 Get 操作
func BenchmarkGet(b *testing.B) {
	cache1 := New[int, int](0)
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
	cache1 := New[int, int](0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache1.SetIfAbsent(i%100, i)
	}
}

// BenchmarkConcurrentSet 基准测试并发 Set 操作
func BenchmarkConcurrentSet(b *testing.B) {
	cache1 := New[int, int](0)
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
	cache1 := New[int, int](0)
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
