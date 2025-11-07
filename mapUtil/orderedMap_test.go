package mapUtil

import (
	"encoding/json"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestOrderedMap(t *testing.T) {
	// 测试新建空map
	t.Run("NewOrderedMap", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		if m.Size() != 0 {
			t.Error("新建的map大小应该为0")
		}
		if m.Front() != nil || m.Back() != nil {
			t.Error("新建的map前后指针应该为nil")
		}
	})

	// 测试带容量的新建
	t.Run("NewOrderedMapWithCapacity", func(t *testing.T) {
		m := NewOrderedMapWithCapacity[string, int](10)
		if m.Size() != 0 {
			t.Error("带容量新建的map大小应该为0")
		}
	})

	// 测试Put和Get
	t.Run("PutAndGet", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		m.Put("a", 1)
		m.Put("b", 2)

		if val := m.Get("a"); val != 1 {
			t.Error("Get获取的值不正确")
		}
		if val := m.Get("b"); val != 2 {
			t.Error("Get获取的值不正确")
		}
		if val := m.Get("c"); val != 0 { // int的零值
			t.Error("不存在的key应该返回零值")
		}
	})

	// 测试ReplaceKey
	t.Run("ReplaceKey", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		m.Put("a", 1)
		m.Put("b", 2)

		if !m.ReplaceKey("a", "c") {
			t.Error("替换key应该成功")
		}
		if m.ReplaceKey("a", "d") {
			t.Error("替换不存在的key应该失败")
		}
		if m.ReplaceKey("c", "b") {
			t.Error("替换为已存在的key应该失败")
		}
		if val := m.Get("c"); val != 1 {
			t.Error("替换key后值应该保持不变")
		}
	})

	// 测试PutIfAbsent
	t.Run("PutIfAbsent", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		m.Put("a", 1)

		val, loaded := m.PutIfAbsent("a", 2)
		if !loaded || val != 1 {
			t.Error("已存在的key应该返回原值")
		}

		val, loaded = m.PutIfAbsent("b", 2)
		if loaded || val != 2 {
			t.Error("不存在的key应该设置新值")
		}
	})

	// 测试GetOrDefault
	t.Run("GetOrDefault", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		m.Put("a", 1)

		if val := m.GetOrDefault("a", 0); val != 1 {
			t.Error("存在的key应该返回对应值")
		}
		if val := m.GetOrDefault("b", 2); val != 2 {
			t.Error("不存在的key应该返回默认值")
		}
	})

	// 测试Remove
	t.Run("Remove", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		m.Put("a", 1)
		m.Put("b", 2)

		m.Remove("a")
		if m.Size() != 1 {
			t.Error("删除后map大小不正确")
		}
		if m.ContainsKey("a") {
			t.Error("删除后key应该不存在")
		}
		if m.Front().Key != "b" {
			t.Error("删除后顺序应该保持")
		}
	})

	// 测试Clear
	t.Run("Clear", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		m.Put("a", 1)
		m.Put("b", 2)

		m.Clear()
		if m.Size() != 0 {
			t.Error("清空后map大小应该为0")
		}
		if m.Front() != nil || m.Back() != nil {
			t.Error("清空后前后指针应该为nil")
		}
	})

	// 测试顺序相关功能
	t.Run("Order", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		m.Put("a", 1)
		m.Put("b", 2)
		m.Put("c", 3)

		// 测试Front和Back
		if m.Front().Key != "a" || m.Front().Value != 1 {
			t.Error("Front返回的元素不正确")
		}
		if m.Back().Key != "c" || m.Back().Value != 3 {
			t.Error("Back返回的元素不正确")
		}

		// 测试Keys和Values顺序
		keys := m.Keys()
		if len(keys) != 3 || keys[0] != "a" || keys[1] != "b" || keys[2] != "c" {
			t.Error("Keys返回的顺序不正确")
		}
		values := m.Values()
		if len(values) != 3 || values[0] != 1 || values[1] != 2 || values[2] != 3 {
			t.Error("Values返回的顺序不正确")
		}
	})

	// 测试Copy
	t.Run("Copy", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		m.Put("a", 1)
		m.Put("b", 2)

		m2 := m.Copy()
		if m2.Size() != m.Size() {
			t.Error("复制后的map大小应该相同")
		}
		if m2.Get("a") != 1 || m2.Get("b") != 2 {
			t.Error("复制后的map值应该相同")
		}
		m2.Put("c", 3)
		if m.ContainsKey("c") {
			t.Error("修改复制后的map不应该影响原map")
		}
	})

	// 测试JSON序列化和反序列化
	t.Run("JSON", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		m.Put("a", 1)
		m.Put("b", 2)

		data, err := json.Marshal(m)
		if err != nil {
			t.Error("JSON序列化失败:", err)
		}

		m2 := NewOrderedMap[string, int]()
		err = json.Unmarshal(data, &m2)
		if err != nil {
			t.Error("JSON反序列化失败:", err)
		}

		if m2.Size() != 2 || m2.Get("a") != 1 || m2.Get("b") != 2 {
			t.Error("反序列化后的map内容不正确")
		}
	})

	// 测试ToString
	t.Run("ToString", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		m.Put("a", 1)
		m.Put("b", 2)

		str := m.ToString()
		if str != `{"a":1,"b":2}` && str != `{"b":2,"a":1}` {
			t.Error("ToString结果不正确:", str)
		}
	})

	// 测试Range
	t.Run("Range", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		m.Put("a", 1)
		m.Put("b", 2)
		m.Put("c", 3)

		count := 0
		m.Range(func(key string, value int) bool {
			count++
			if count >= 2 {
				return false
			}
			return true
		})

		if count != 2 {
			t.Error("Range应该能提前终止")
		}
	})

	// 测试BSON序列化和反序列化
	t.Run("BSON", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		m.Put("a", 1)
		m.Put("b", 2)
		// 测试MarshalBSON
		bsonData, err := bson.Marshal(m)
		if err != nil {
			t.Error("BSON序列化失败:", err)
		}
		// 测试UnmarshalBSON
		m2 := NewOrderedMap[string, int]()
		err = bson.Unmarshal(bsonData, &m2)
		if err != nil {
			t.Error("BSON反序列化失败:", err)
		}
		if m2.Size() != 2 {
			t.Error("反序列化后的map大小不正确")
		}
		if val := m2.Get("a"); val != 1 {
			t.Error("反序列化后的值不正确")
		}
		if val := m2.Get("b"); val != 2 {
			t.Error("反序列化后的值不正确")
		}
		// 测试空map的BSON处理
		emptyMap := NewOrderedMap[string, int]()
		emptyBson, err := bson.Marshal(emptyMap)
		if err != nil {
			t.Error("空map的BSON序列化失败:", err)
		}
		emptyMap2 := NewOrderedMap[string, int]()
		err = bson.Unmarshal(emptyBson, &emptyMap2)
		if err != nil {
			t.Error("空map的BSON反序列化失败:", err)
		}
		if emptyMap2.Size() != 0 {
			t.Error("反序列化后的空map大小应该为0")
		}
		// 测试无效BSON数据
		invalidData := []byte{0x01, 0x02, 0x03} // 无效的BSON数据
		invalidMap := NewOrderedMap[string, int]()
		err = bson.Unmarshal(invalidData, &invalidMap)
		if err == nil {
			t.Error("无效的BSON数据应该返回错误")
		}
	})

}
