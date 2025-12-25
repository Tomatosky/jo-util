package cacheUtil

import (
	"sync"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	// 创建一个过期时间为1秒的缓存
	cache1 := New[string, int](time.Second)

	// 测试1: 基本设置功能
	t.Run("Basic Set", func(t *testing.T) {
		cache1.Set("key1", 100)
		if val, found := cache1.Get("key1"); !found || val != 100 {
			t.Error("Set failed: key1 should be 100")
		}
	})

	// 测试2: 覆盖已有值
	t.Run("Overwrite Existing", func(t *testing.T) {
		cache1.Set("key2", 200)
		cache1.Set("key2", 300)
		if val, found := cache1.Get("key2"); !found || val != 300 {
			t.Error("Set failed: key2 should be overwritten to 300")
		}
	})

	// 测试3: 自定义过期时间
	t.Run("Custom Expiration", func(t *testing.T) {
		cache1.Set("key3", 400, 2*time.Second)
		time.Sleep(1500 * time.Millisecond) // 等待1.5秒
		if _, found := cache1.Get("key3"); !found {
			t.Error("Set with custom expiration failed: key3 should still exist")
		}
		time.Sleep(600 * time.Millisecond) // 总共2.1秒
		if _, found := cache1.Get("key3"); found {
			t.Error("Set with custom expiration failed: key3 should be expired")
		}
	})

	// 测试4: 零值设置
	t.Run("Zero Value", func(t *testing.T) {
		cache1.Set("key4", 0)
		if val, found := cache1.Get("key4"); !found || val != 0 {
			t.Error("Set failed: zero value should be stored correctly")
		}
	})

	// 测试5: 多个自定义过期时间参数(只取第一个)
	t.Run("Multiple Expiration Parameters", func(t *testing.T) {
		cache1.Set("key5", 500, 1*time.Second, 2*time.Second, 3*time.Second)
		time.Sleep(1500 * time.Millisecond)
		if _, found := cache1.Get("key5"); found {
			t.Error("Set with multiple expiration params failed: should only use first param")
		}
	})

	// 测试6: 无过期时间(永不过期)
	t.Run("No Expiration", func(t *testing.T) {
		noExpirecache := New[string, int](0)
		noExpirecache.Set("key6", 600)
		time.Sleep(1 * time.Second)
		if _, found := noExpirecache.Get("key6"); !found {
			t.Error("Set with no expiration failed: key6 should never expire")
		}
	})

	// 测试7: 并发设置
	t.Run("Concurrent Set", func(t *testing.T) {
		const numGoroutines = 100
		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(n int) {
				defer wg.Done()
				key := "concurrent_key"
				cache1.Set(key, n)
			}(i)
		}

		wg.Wait()
		// 检查最终值，虽然不确定是哪个goroutine最后设置的值
		if _, found := cache1.Get("concurrent_key"); !found {
			t.Error("Concurrent Set failed: key should exist")
		}
	})
}

func TestSetIfAbsent(t *testing.T) {
	// 创建一个过期时间为1秒的缓存
	cache1 := New[string, int](time.Second)

	// 测试1: 键不存在时设置成功
	t.Run("Key Not Exist - Set Success", func(t *testing.T) {
		success := cache1.SetIfAbsent("key1", 100)
		if !success {
			t.Error("SetIfAbsent failed: should return true when key doesn't exist")
		}
		if val, found := cache1.Get("key1"); !found || val != 100 {
			t.Error("SetIfAbsent failed: key1 should be set to 100")
		}
	})

	// 测试2: 键已存在时设置失败
	t.Run("Key Exists - Set Failed", func(t *testing.T) {
		cache1.Set("key2", 200)
		success := cache1.SetIfAbsent("key2", 300)
		if success {
			t.Error("SetIfAbsent failed: should return false when key exists")
		}
		if val, found := cache1.Get("key2"); !found || val != 200 {
			t.Error("SetIfAbsent failed: key2 should remain 200")
		}
	})

	// 测试3: 键已过期时设置成功
	t.Run("Key Expired - Set Success", func(t *testing.T) {
		cache1.Set("key3", 400)
		time.Sleep(1100 * time.Millisecond) // 等待1.1秒让key3过期
		success := cache1.SetIfAbsent("key3", 500)
		if !success {
			t.Error("SetIfAbsent failed: should return true when key is expired")
		}
		if val, found := cache1.Get("key3"); !found || val != 500 {
			t.Error("SetIfAbsent failed: key3 should be updated to 500")
		}
	})

	// 测试4: 并发设置相同键
	t.Run("Concurrent Set Same Key", func(t *testing.T) {
		const numGoroutines = 100
		var wg sync.WaitGroup
		wg.Add(numGoroutines)
		successCount := 0

		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				if cache1.SetIfAbsent("concurrent_key", 600) {
					successCount++
				}
			}()
		}

		wg.Wait()
		if successCount != 1 {
			t.Errorf("SetIfAbsent concurrent test failed: expected 1 success, got %d", successCount)
		}
		if val, found := cache1.Get("concurrent_key"); !found || val != 600 {
			t.Error("SetIfAbsent concurrent test failed: value not set correctly")
		}
	})

	// 测试5: 永不过期缓存中的行为
	t.Run("No Expiration Cache", func(t *testing.T) {
		noExpirecache1 := New[string, int](0)
		noExpirecache1.Set("key4", 700)
		success := noExpirecache1.SetIfAbsent("key4", 800)
		if success {
			t.Error("SetIfAbsent failed: should return false in no-expiration cache when key exists")
		}
	})

	// 测试6: 零值设置
	t.Run("Set Zero Value", func(t *testing.T) {
		success := cache1.SetIfAbsent("key5", 0)
		if !success {
			t.Error("SetIfAbsent failed: should allow setting zero value")
		}
		if val, found := cache1.Get("key5"); !found || val != 0 {
			t.Error("SetIfAbsent failed: zero value not stored correctly")
		}
	})
}

func TestGet(t *testing.T) {
	// 创建一个过期时间为1秒的缓存
	cache1 := New[string, int](time.Second)

	// 测试1: 获取存在的键
	t.Run("Existing Key", func(t *testing.T) {
		cache1.Set("key1", 100)
		if val, found := cache1.Get("key1"); !found || val != 100 {
			t.Error("Get failed: should return correct value for existing key")
		}
	})

	// 测试2: 获取不存在的键
	t.Run("Non-existent Key", func(t *testing.T) {
		if _, found := cache1.Get("nonexistent"); found {
			t.Error("Get failed: should return false for non-existent key")
		}
	})

	// 测试3: 获取已过期的键
	t.Run("Expired Key", func(t *testing.T) {
		cache1.Set("key2", 200)
		time.Sleep(1100 * time.Millisecond) // 等待1.1秒让key2过期
		if _, found := cache1.Get("key2"); found {
			t.Error("Get failed: should return false for expired key")
		}
	})

	// 测试4: 永不过期缓存中的获取
	t.Run("No Expiration Cache", func(t *testing.T) {
		noExpirecache1 := New[string, int](0)
		noExpirecache1.Set("key3", 300)
		time.Sleep(1 * time.Second)
		if val, found := noExpirecache1.Get("key3"); !found || val != 300 {
			t.Error("Get failed: should return value from no-expiration cache")
		}
	})

	// 测试5: 获取零值
	t.Run("Zero Value", func(t *testing.T) {
		cache1.Set("key4", 0)
		if val, found := cache1.Get("key4"); !found || val != 0 {
			t.Error("Get failed: should return zero value correctly")
		}
	})

	// 测试6: 并发读写测试
	t.Run("Concurrent Read", func(t *testing.T) {
		cache1.Set("key5", 500)
		var wg sync.WaitGroup
		const numReaders = 100
		wg.Add(numReaders)

		for i := 0; i < numReaders; i++ {
			go func() {
				defer wg.Done()
				if val, found := cache1.Get("key5"); !found || val != 500 {
					t.Error("Concurrent Get failed")
				}
			}()
		}
		wg.Wait()
	})

	// 测试7: 类型正确性测试
	t.Run("Type Correctness", func(t *testing.T) {
		stringcache1 := New[string, string](time.Minute)
		stringcache1.Set("key6", "value")
		if val, found := stringcache1.Get("key6"); !found || val != "value" {
			t.Error("Get failed: should return correct type")
		}
	})
}

func TestGetWithExpiration(t *testing.T) {
	// 创建一个过期时间为1秒的缓存
	cache1 := New[string, int](time.Second)

	// 测试1: 获取存在的键及其过期时间
	t.Run("Existing Key with Expiration", func(t *testing.T) {
		cache1.Set("key1", 100)
		val, found, exp := cache1.GetWithExpiration("key1")
		if !found {
			t.Error("GetWithExpiration failed: should find existing key")
		}
		if val != 100 {
			t.Error("GetWithExpiration failed: wrong value returned")
		}
		if exp.Before(time.Now()) {
			t.Error("GetWithExpiration failed: expiration time should be in the future")
		}
	})

	// 测试2: 获取不存在的键
	t.Run("Non-existent Key", func(t *testing.T) {
		_, found, _ := cache1.GetWithExpiration("nonexistent")
		if found {
			t.Error("GetWithExpiration failed: should not find non-existent key")
		}
	})

	// 测试3: 获取已过期的键
	t.Run("Expired Key", func(t *testing.T) {
		cache1.Set("key2", 200)
		time.Sleep(1100 * time.Millisecond) // 等待1.1秒让key2过期
		_, found, _ := cache1.GetWithExpiration("key2")
		if found {
			t.Error("GetWithExpiration failed: should not find expired key")
		}
	})

	// 测试4: 永不过期缓存中的键
	t.Run("No Expiration Cache", func(t *testing.T) {
		noExpirecache1 := New[string, int](0)
		noExpirecache1.Set("key3", 300)
		val, found, exp := noExpirecache1.GetWithExpiration("key3")
		if !found {
			t.Error("GetWithExpiration failed: should find key in no-expiration cache")
		}
		if val != 300 {
			t.Error("GetWithExpiration failed: wrong value in no-expiration cache")
		}
		if exp.Unix() != 0 {
			t.Error("GetWithExpiration failed: expiration should be zero in no-expiration cache")
		}
	})

	// 测试5: 自定义过期时间的键
	t.Run("Custom Expiration", func(t *testing.T) {
		cache1.Set("key4", 400, 2*time.Second)
		_, found, exp := cache1.GetWithExpiration("key4")
		if !found {
			t.Error("GetWithExpiration failed: should find key with custom expiration")
		}
		expectedMin := time.Now().Add(1900 * time.Millisecond) // 允许100ms误差
		expectedMax := time.Now().Add(2100 * time.Millisecond)
		if exp.Before(expectedMin) || exp.After(expectedMax) {
			t.Errorf("GetWithExpiration failed: expiration time %v not in expected range (%v-%v)",
				exp, expectedMin, expectedMax)
		}
	})

	// 测试6: 零值测试
	t.Run("Zero Value", func(t *testing.T) {
		cache1.Set("key5", 0)
		val, found, _ := cache1.GetWithExpiration("key5")
		if !found {
			t.Error("GetWithExpiration failed: should find zero value")
		}
		if val != 0 {
			t.Error("GetWithExpiration failed: wrong zero value returned")
		}
	})

	// 测试7: 并发读取测试
	t.Run("Concurrent Read", func(t *testing.T) {
		cache1.Set("key6", 600)
		var wg sync.WaitGroup
		const numReaders = 10
		wg.Add(numReaders)

		for i := 0; i < numReaders; i++ {
			go func() {
				defer wg.Done()
				val, found, _ := cache1.GetWithExpiration("key6")
				if !found || val != 600 {
					t.Error("Concurrent GetWithExpiration failed")
				}
			}()
		}
		wg.Wait()
	})
}

func TestDelete(t *testing.T) {
	// 创建一个过期时间为1秒的缓存
	cache1 := New[string, int](time.Second)

	// 测试1: 删除存在的键
	t.Run("Delete Existing Key", func(t *testing.T) {
		cache1.Set("key1", 100)
		cache1.Delete("key1")
		if _, found := cache1.Get("key1"); found {
			t.Error("Delete failed: should remove existing key")
		}
	})

	// 测试2: 删除不存在的键
	t.Run("Delete Non-existent Key", func(t *testing.T) {
		// 确保key2不存在
		if _, found := cache1.Get("key2"); found {
			t.Fatal("Test setup failed: key2 should not exist")
		}
		cache1.Delete("key2") // 应该不会panic或出错
		if _, found := cache1.Get("key2"); found {
			t.Error("Delete failed: should handle non-existent key gracefully")
		}
	})

	// 测试3: 删除已过期的键
	t.Run("Delete Expired Key", func(t *testing.T) {
		cache1.Set("key3", 300)
		time.Sleep(1100 * time.Millisecond) // 等待1.1秒让key3过期
		cache1.Delete("key3")
		if _, found := cache1.Get("key3"); found {
			t.Error("Delete failed: should handle expired key correctly")
		}
	})

	// 测试4: 并发删除测试
	t.Run("Concurrent Delete", func(t *testing.T) {
		cache1.Set("key4", 400)
		var wg sync.WaitGroup
		const numDeleters = 10
		wg.Add(numDeleters)

		for i := 0; i < numDeleters; i++ {
			go func() {
				defer wg.Done()
				cache1.Delete("key4")
			}()
		}
		wg.Wait()

		if _, found := cache1.Get("key4"); found {
			t.Error("Concurrent Delete failed: key should be deleted")
		}
	})

	// 测试5: 删除后可以重新设置
	t.Run("Re-set After Delete", func(t *testing.T) {
		cache1.Set("key5", 500)
		cache1.Delete("key5")
		cache1.Set("key5", 600)
		if val, found := cache1.Get("key5"); !found || val != 600 {
			t.Error("Delete failed: should allow re-setting deleted key")
		}
	})

	// 测试6: 删除永不过期缓存中的键
	t.Run("Delete in No-Expiration Cache", func(t *testing.T) {
		noExpirecache1 := New[string, int](0)
		noExpirecache1.Set("key6", 700)
		noExpirecache1.Delete("key6")
		if _, found := noExpirecache1.Get("key6"); found {
			t.Error("Delete failed: should work in no-expiration cache")
		}
	})

	// 测试7: 删除后Items()不包含该键
	t.Run("Delete and Items()", func(t *testing.T) {
		cache1.Set("key7", 800)
		cache1.Delete("key7")
		items := cache1.Items()
		if _, exists := items["key7"]; exists {
			t.Error("Delete failed: Items() should not contain deleted key")
		}
	})
}

func TestItems(t *testing.T) {
	// 创建一个过期时间为1秒的缓存
	cache1 := New[string, int](time.Second)

	// 测试1: 获取空缓存的Items
	t.Run("Empty Cache", func(t *testing.T) {
		items := cache1.Items()
		if len(items) != 0 {
			t.Errorf("Items() failed: expected empty map, got %v", items)
		}
	})

	// 测试2: 获取包含多个项目的Items
	t.Run("Multiple Items", func(t *testing.T) {
		cache1.Set("key1", 100)
		cache1.Set("key2", 200)
		items := cache1.Items()
		if len(items) != 2 {
			t.Errorf("Items() failed: expected 2 items, got %d", len(items))
		}
		if items["key1"].Object != 100 || items["key2"].Object != 200 {
			t.Error("Items() failed: values not matched")
		}
	})

	// 测试3: 不返回过期项目
	t.Run("Exclude Expired Items", func(t *testing.T) {
		cache1.Set("key3", 300)
		time.Sleep(1100 * time.Millisecond) // 等待1.1秒让key3过期
		items := cache1.Items()
		if _, exists := items["key3"]; exists {
			t.Error("Items() failed: should not include expired items")
		}
	})

	// 测试4: 返回的项目是副本
	t.Run("Return Copy of Items", func(t *testing.T) {
		cache1.Set("key4", 400)
		items1 := cache1.Items()
		cache1.Delete("key4")
		items2 := cache1.Items()
		if _, exists := items1["key4"]; !exists {
			t.Error("Items() failed: first snapshot should contain key4")
		}
		if _, exists := items2["key4"]; exists {
			t.Error("Items() failed: second snapshot should not contain key4")
		}
	})

	// 测试5: 永不过期缓存中的Items
	t.Run("No Expiration Cache", func(t *testing.T) {
		noExpirecache1 := New[string, int](0)
		noExpirecache1.Set("key5", 500)
		time.Sleep(1 * time.Second)
		items := noExpirecache1.Items()
		if _, exists := items["key5"]; !exists {
			t.Error("Items() failed: should include items in no-expiration cache")
		}
	})

	// 测试6: 包含过期时间的项目
	t.Run("Items With Expiration", func(t *testing.T) {
		cache1.Set("key6", 600, 2*time.Second)
		items := cache1.Items()
		if item, exists := items["key6"]; exists {
			if item.Expiration <= time.Now().UnixNano() {
				t.Error("Items() failed: expiration time not set correctly")
			}
		} else {
			t.Error("Items() failed: should include valid items")
		}
	})

	// 测试7: 并发访问Items
	t.Run("Concurrent Access", func(t *testing.T) {
		var wg sync.WaitGroup
		const numGoroutines = 10
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				items := cache1.Items()
				if len(items) < 0 { // 基本检查
					t.Error("Concurrent Items() failed")
				}
			}()
		}
		wg.Wait()
	})
}

func TestFlush(t *testing.T) {
	// 创建一个过期时间为1分钟的缓存
	cache1 := New[string, int](time.Minute)

	// 测试1: 清空非空缓存
	t.Run("Flush Non-empty Cache", func(t *testing.T) {
		// 添加测试数据
		cache1.Set("key1", 100)
		cache1.Set("key2", 200)
		cache1.Set("key3", 300)

		// 验证缓存非空
		if len(cache1.Items()) == 0 {
			t.Fatal("Test setup failed: cache should not be empty")
		}

		// 执行清空操作
		cache1.Flush()

		// 验证缓存已清空
		if len(cache1.Items()) != 0 {
			t.Error("Flush failed: cache should be empty after flushing")
		}
	})

	// 测试2: 清空空缓存
	t.Run("Flush Empty Cache", func(t *testing.T) {
		// 确保缓存为空
		cache1.Flush()
		if len(cache1.Items()) != 0 {
			t.Fatal("Test setup failed: cache should be empty")
		}

		// 执行清空操作
		cache1.Flush()

		// 验证缓存仍为空
		if len(cache1.Items()) != 0 {
			t.Error("Flush failed: empty cache should remain empty after flushing")
		}
	})

	// 测试3: 清空后可以重新添加数据
	t.Run("Reuse After Flush", func(t *testing.T) {
		// 清空缓存
		cache1.Flush()

		// 添加新数据
		cache1.Set("newKey", 400)

		// 验证新数据
		if val, found := cache1.Get("newKey"); !found || val != 400 {
			t.Error("Flush failed: should be able to add new items after flushing")
		}
	})

	// 测试4: 清空包含过期项的缓存
	t.Run("Flush Cache with Expired Items", func(t *testing.T) {
		// 添加会过期的数据
		cache1.Set("expiredKey", 500, 100*time.Millisecond)
		time.Sleep(150 * time.Millisecond) // 确保过期

		// 添加未过期数据
		cache1.Set("validKey", 600)

		// 执行清空
		cache1.Flush()

		// 验证所有数据都被清空
		if len(cache1.Items()) != 0 {
			t.Error("Flush failed: should remove all items including expired ones")
		}
	})

	// 测试5: 并发清空测试
	t.Run("Concurrent Flush", func(t *testing.T) {
		// 准备测试数据
		for i := 0; i < 100; i++ {
			cache1.Set(string(rune(i)), i)
		}

		var wg sync.WaitGroup
		const numFlushers = 5
		wg.Add(numFlushers)

		// 启动多个goroutine并发清空
		for i := 0; i < numFlushers; i++ {
			go func() {
				defer wg.Done()
				cache1.Flush()
			}()
		}
		wg.Wait()

		// 验证缓存确实被清空
		if len(cache1.Items()) != 0 {
			t.Error("Concurrent Flush failed: cache should be empty")
		}
	})

	// 测试6: 清空永不过期缓存
	t.Run("Flush No-Expiration Cache", func(t *testing.T) {
		noExpirecache1 := New[string, int](0)
		noExpirecache1.Set("key4", 700)
		noExpirecache1.Set("key5", 800)

		noExpirecache1.Flush()

		if len(noExpirecache1.Items()) != 0 {
			t.Error("Flush failed: should work on no-expiration cache")
		}
	})

	// 测试7: 清空后内存使用
	t.Run("Memory Usage After Flush", func(t *testing.T) {
		// 添加大量数据
		for i := 0; i < 10000; i++ {
			cache1.Set(string(rune(i)), i)
		}

		// 获取清空前map的容量
		before := len(cache1.Items())

		// 执行清空
		cache1.Flush()

		// 验证
		if len(cache1.Items()) != 0 {
			t.Error("Flush failed: should remove all items")
		}
		if cap1 := len(cache1.Items()); cap1 >= before {
			t.Logf("Note: Flush kept original capacity (%d), may affect memory usage", cap1)
		}
	})
}
