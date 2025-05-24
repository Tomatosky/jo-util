package cacheUtil

import (
	"testing"
	"time"
)

func TestCache_SetAndGet(t *testing.T) {
	cacheO := New[string, int](time.Minute)

	// 测试基本设置和获取
	cacheO.Set("key1", 100)
	if val, found := cacheO.Get("key1"); !found || val != 100 {
		t.Errorf("Expected to find key1 with value 100, got %v, %v", val, found)
	}

	// 测试不存在的键
	if _, found := cacheO.Get("nonexistent"); found {
		t.Error("Expected not to find nonexistent key")
	}
}

func TestCache_SetWithCustomExpiration(t *testing.T) {
	cacheO := New[string, string](time.Minute)

	// 测试自定义过期时间
	cacheO.Set("temp", "value", time.Millisecond*100)
	if val, found := cacheO.Get("temp"); !found || val != "value" {
		t.Error("Expected to find temp key immediately after setting")
	}

	// 等待过期
	time.Sleep(time.Millisecond * 150)
	if _, found := cacheO.Get("temp"); found {
		t.Error("Expected temp key to be expired")
	}
}

func TestCache_SetIfAbsent(t *testing.T) {
	cacheO := New[string, float64](time.Minute)

	// 测试首次设置
	if !cacheO.SetIfAbsent("pi", 3.14) {
		t.Error("Expected SetIfAbsent to return true for new key")
	}

	// 测试重复设置
	if cacheO.SetIfAbsent("pi", 3.14159) {
		t.Error("Expected SetIfAbsent to return false for existing key")
	}

	// 验证值未被覆盖
	if val, _ := cacheO.Get("pi"); val != 3.14 {
		t.Errorf("Expected pi to remain 3.14, got %v", val)
	}
}

func TestCache_Delete(t *testing.T) {
	cacheO := New[int, string](time.Minute)

	cacheO.Set(42, "answer")
	cacheO.Delete(42)

	if _, found := cacheO.Get(42); found {
		t.Error("Expected key 42 to be deleted")
	}
}

func TestCache_Items(t *testing.T) {
	cacheO := New[string, bool](time.Minute)

	cacheO.Set("a", true)
	cacheO.Set("b", false)

	items := cacheO.Items()
	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}

	if val, ok := items["a"]; !ok || val.Object != true {
		t.Error("Expected item a to be true")
	}
	if val, ok := items["b"]; !ok || val.Object != false {
		t.Error("Expected item b to be false")
	}
}

func TestCache_Flush(t *testing.T) {
	cacheO := New[string, int](time.Minute)

	cacheO.Set("one", 1)
	cacheO.Set("two", 2)
	cacheO.Flush()

	if len(cacheO.Items()) != 0 {
		t.Error("Expected cacheO to be empty after flush")
	}
}

func TestCache_Expiration(t *testing.T) {
	// 使用很短的过期时间测试自动过期
	cacheO := New[string, int](time.Millisecond * 50)

	cacheO.Set("short", 123)
	time.Sleep(time.Millisecond * 10)
	if _, found := cacheO.Get("short"); !found {
		t.Error("Expected key to exist before expiration")
	}

	time.Sleep(time.Millisecond * 50)
	if _, found := cacheO.Get("short"); found {
		t.Error("Expected key to be expired")
	}
}

func TestCache_NoExpiration(t *testing.T) {
	cacheO := New[string, int](0) // 永不过期

	cacheO.Set("forever", 999)
	time.Sleep(time.Millisecond * 10)
	if _, found := cacheO.Get("forever"); !found {
		t.Error("Expected key to exist with no expiration")
	}
}

func TestCache_JanitorCleanup(t *testing.T) {
	// 测试janitor自动清理
	cacheO := New[string, int](time.Millisecond * 50)

	cacheO.Set("temp1", 1)
	cacheO.Set("temp2", 2, time.Millisecond*150) // 自定义更长的过期时间

	// 等待第一次清理
	time.Sleep(time.Millisecond * 120)

	items := cacheO.Items()
	if len(items) != 1 {
		t.Errorf("Expected 1 item after cleanup, got %d", len(items))
	}

	if _, found := items["temp2"]; !found {
		t.Error("Expected temp2 to still exist")
	}
}
