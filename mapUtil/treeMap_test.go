package mapUtil

import (
	"encoding/json"
	"testing"
)

func TestTreeMap(t *testing.T) {
	// 创建一个用于测试的TreeMap
	less := func(a, b string) bool { return a < b }
	tm := NewTreeMap[string, int](less)

	// 测试Put和Get
	t.Run("PutAndGet", func(t *testing.T) {
		tm.Put("a", 1)
		tm.Put("b", 2)
		tm.Put("c", 3)

		if val := tm.Get("a"); val != 1 {
			t.Error("Get(a) expected 1, got", val)
		}
		if val := tm.Get("b"); val != 2 {
			t.Error("Get(b) expected 2, got", val)
		}
		if val := tm.Get("c"); val != 3 {
			t.Error("Get(c) expected 3, got", val)
		}
		if val := tm.Get("d"); val != 0 { // 0是int的零值
			t.Error("Get(d) expected 0, got", val)
		}
	})

	// 测试Size
	t.Run("Size", func(t *testing.T) {
		if size := tm.Size(); size != 3 {
			t.Error("Size expected 3, got", size)
		}
	})

	// 测试ContainsKey
	t.Run("ContainsKey", func(t *testing.T) {
		if !tm.ContainsKey("a") {
			t.Error("ContainsKey(a) expected true, got false")
		}
		if tm.ContainsKey("d") {
			t.Error("ContainsKey(d) expected false, got true")
		}
	})

	// 测试Remove
	t.Run("Remove", func(t *testing.T) {
		tm.Remove("b")
		if tm.ContainsKey("b") {
			t.Error("After Remove(b), ContainsKey(b) should be false")
		}
		if size := tm.Size(); size != 2 {
			t.Error("After Remove(b), Size expected 2, got", size)
		}
	})

	// 测试Clear
	t.Run("Clear", func(t *testing.T) {
		tm.Clear()
		if size := tm.Size(); size != 0 {
			t.Error("After Clear, Size expected 0, got", size)
		}
		if tm.ContainsKey("a") {
			t.Error("After Clear, ContainsKey(a) should be false")
		}
	})

	// 测试Keys和Values
	t.Run("KeysAndValues", func(t *testing.T) {
		tm.Put("z", 26)
		tm.Put("a", 1)
		tm.Put("m", 13)

		keys := tm.Keys()
		if len(keys) != 3 {
			t.Error("Keys length expected 3, got", len(keys))
		}
		if keys[0] != "a" || keys[1] != "m" || keys[2] != "z" {
			t.Error("Keys not in sorted order, got", keys)
		}

		values := tm.Values()
		if len(values) != 3 {
			t.Error("Values length expected 3, got", len(values))
		}
		if values[0] != 1 || values[1] != 13 || values[2] != 26 {
			t.Error("Values not in key-sorted order, got", values)
		}
	})

	// 测试PutIfAbsent
	t.Run("PutIfAbsent", func(t *testing.T) {
		existing, loaded := tm.PutIfAbsent("a", 100)
		if !loaded || existing != 1 {
			t.Error("PutIfAbsent(a) expected (1, true), got", existing, loaded)
		}

		existing, loaded = tm.PutIfAbsent("b", 2)
		if loaded || existing != 0 {
			t.Error("PutIfAbsent(b) expected (0, false), got", existing, loaded)
		}
		if val := tm.Get("b"); val != 2 {
			t.Error("After PutIfAbsent(b), Get(b) expected 2, got", val)
		}
	})

	// 测试GetOrDefault
	t.Run("GetOrDefault", func(t *testing.T) {
		if val := tm.GetOrDefault("a", 100); val != 1 {
			t.Error("GetOrDefault(a) expected 1, got", val)
		}
		if val := tm.GetOrDefault("x", 24); val != 24 {
			t.Error("GetOrDefault(x) expected 24, got", val)
		}
	})

	// 测试FirstKey和LastKey
	t.Run("FirstAndLastKey", func(t *testing.T) {
		first, ok := tm.FirstKey()
		if !ok || first != "a" {
			t.Error("FirstKey expected (a, true), got", first, ok)
		}

		last, ok := tm.LastKey()
		if !ok || last != "z" {
			t.Error("LastKey expected (z, true), got", last, ok)
		}

		tm.Clear()
		_, ok = tm.FirstKey()
		if ok {
			t.Error("Empty map FirstKey should return false")
		}
		_, ok = tm.LastKey()
		if ok {
			t.Error("Empty map LastKey should return false")
		}
	})

	// 测试Range
	t.Run("Range", func(t *testing.T) {
		tm.Clear()
		tm.Put("a", 1)
		tm.Put("b", 2)
		tm.Put("c", 3)

		var keys []string
		var values []int
		tm.Range(func(key string, value int) bool {
			keys = append(keys, key)
			values = append(values, value)
			return true
		})

		if len(keys) != 3 || keys[0] != "a" || keys[1] != "b" || keys[2] != "c" {
			t.Error("Range keys expected [a b c], got", keys)
		}
		if len(values) != 3 || values[0] != 1 || values[1] != 2 || values[2] != 3 {
			t.Error("Range values expected [1 2 3], got", values)
		}

		// 测试提前终止
		count := 0
		tm.Range(func(key string, value int) bool {
			count++
			return count < 2
		})
		if count != 2 {
			t.Error("Range should stop after 2 iterations, got", count)
		}
	})

	// 测试JSON序列化和反序列化
	t.Run("JSON", func(t *testing.T) {
		tm.Clear()
		tm.Put("a", 1)
		tm.Put("b", 2)

		data, err := json.Marshal(tm)
		if err != nil {
			t.Error("Marshal error:", err)
		}

		newTm := NewTreeMap[string, int](less)
		err = json.Unmarshal(data, &newTm)
		if err != nil {
			t.Error("Unmarshal error:", err)
		}

		if newTm.Size() != 2 {
			t.Error("After Unmarshal, Size expected 2, got", newTm.Size())
		}
		if val := newTm.Get("a"); val != 1 {
			t.Error("After Unmarshal, Get(a) expected 1, got", val)
		}
		if val := newTm.Get("b"); val != 2 {
			t.Error("After Unmarshal, Get(b) expected 2, got", val)
		}
	})

	// 测试ToString
	t.Run("ToString", func(t *testing.T) {
		tm.Clear()
		tm.Put("x", 24)
		tm.Put("y", 25)

		str := tm.ToString()
		expected := `{"x":24,"y":25}`
		if str != expected {
			t.Errorf("ToString expected %s, got %s", expected, str)
		}
	})

	// 测试并发安全性
	t.Run("Concurrency", func(t *testing.T) {
		tm.Clear()
		const num = 1000
		done := make(chan bool)

		// 并发写入
		for i := 0; i < num; i++ {
			go func(i int) {
				tm.Put(string(rune(i)), i)
				done <- true
			}(i)
		}

		// 等待所有写入完成
		for i := 0; i < num; i++ {
			<-done
		}

		if size := tm.Size(); size != num {
			t.Error("After concurrent writes, Size expected", num, "got", size)
		}

		// 并发读取
		for i := 0; i < num; i++ {
			go func(i int) {
				val := tm.Get(string(rune(i)))
				if val != i {
					t.Errorf("Concurrent Get(%d) expected %d, got %d", i, i, val)
				}
				done <- true
			}(i)
		}

		// 等待所有读取完成
		for i := 0; i < num; i++ {
			<-done
		}
	})

	// 测试BSON序列化和反序列化
	t.Run("BSON", func(t *testing.T) {
		tm.Clear()
		tm.Put("a", 1)
		tm.Put("b", 2)

		data, err := tm.MarshalBSON()
		if err != nil {
			t.Error("MarshalBSON error:", err)
		}

		newTm := NewTreeMap[string, int](less)
		err = newTm.UnmarshalBSON(data)
		if err != nil {
			t.Error("UnmarshalBSON error:", err)
		}

		if newTm.Size() != 2 {
			t.Error("After UnmarshalBSON, Size expected 2, got", newTm.Size())
		}
		if val := newTm.Get("a"); val != 1 {
			t.Error("After UnmarshalBSON, Get(a) expected 1, got", val)
		}
		if val := newTm.Get("b"); val != 2 {
			t.Error("After UnmarshalBSON, Get(b) expected 2, got", val)
		}
	})

	// 测试空树的BSON序列化
	t.Run("EmptyBSON", func(t *testing.T) {
		tm.Clear()

		data, err := tm.MarshalBSON()
		if err != nil {
			t.Error("MarshalBSON on empty map error:", err)
		}

		newTm := NewTreeMap[string, int](less)
		err = newTm.UnmarshalBSON(data)
		if err != nil {
			t.Error("UnmarshalBSON empty map error:", err)
		}

		if newTm.Size() != 0 {
			t.Error("After UnmarshalBSON empty map, Size expected 0, got", newTm.Size())
		}
	})

	// 测试BSON反序列化错误情况
	t.Run("BSONError", func(t *testing.T) {
		invalidData := []byte("invalid bson data")
		newTm := NewTreeMap[string, int](less)
		err := newTm.UnmarshalBSON(invalidData)
		if err == nil {
			t.Error("Expected error for invalid BSON data, got nil")
		}
	})

}
