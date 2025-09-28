package sliceUtil

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"runtime/debug"
	"sync"
)

// CopyOnWriteSlice 线程安全的动态数组，写时复制
type CopyOnWriteSlice[T comparable] struct {
	mu   sync.RWMutex // 读写锁
	data []T          // 实际存储数据的数组
}

// NewCopyOnWriteSlice 创建新的线程安全动态数组
func NewCopyOnWriteSlice[T comparable]() *CopyOnWriteSlice[T] {
	return &CopyOnWriteSlice[T]{
		data: make([]T, 0),
	}
}

// Add 添加元素到末尾（写操作需要复制整个数组）
func (c *CopyOnWriteSlice[T]) Add(element T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 创建新数组并追加元素
	newData := append(make([]T, 0, len(c.data)+1), c.data...)
	newData = append(newData, element)
	c.data = newData
}

func (c *CopyOnWriteSlice[T]) AddAll(elements ...T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 创建新数组并追加所有元素
	newData := append(make([]T, 0, len(c.data)+len(elements)), c.data...)
	newData = append(newData, elements...)
	c.data = newData
}

// Insert 在指定索引插入元素
func (c *CopyOnWriteSlice[T]) Insert(index int, element T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if index < 0 || index > len(c.data) {
		debug.PrintStack()
		panic("index out of range")
	}

	// 创建长度+1的新数组
	newData := make([]T, len(c.data)+1)
	// 复制前半部分
	copy(newData[:index], c.data[:index])
	// 插入元素
	newData[index] = element
	// 复制后半部分
	copy(newData[index+1:], c.data[index:])
	c.data = newData
}

// Get 获取元素（读操作使用读锁）
func (c *CopyOnWriteSlice[T]) Get(index int) T {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if index < 0 {
		index += len(c.data)
	}
	if index < 0 || index >= len(c.data) {
		debug.PrintStack()
		panic("index out of range")
	}
	return c.data[index]
}

// Remove 移除指定索引元素
func (c *CopyOnWriteSlice[T]) Remove(index int) T {
	c.mu.Lock()
	defer c.mu.Unlock()

	if index < 0 || index >= len(c.data) {
		debug.PrintStack()
		panic("index out of range")
	}

	removed := c.data[index]
	// 创建新数组并跳过指定元素
	newData := make([]T, len(c.data)-1)
	copy(newData[:index], c.data[:index])
	copy(newData[index:], c.data[index+1:])
	c.data = newData
	return removed
}

// Size 当前元素数量（读操作）
func (c *CopyOnWriteSlice[T]) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.data)
}

// Range 安全遍历（遍历过程中数组不会被修改）
func (c *CopyOnWriteSlice[T]) Range(f func(int, T) bool) {
	c.mu.RLock()
	snapshot := append(make([]T, 0, len(c.data)), c.data...) // 创建快照
	c.mu.RUnlock()

	for i, v := range snapshot {
		if !f(i, v) {
			break
		}
	}
}

// Contains 检查元素是否存在（读操作）
func (c *CopyOnWriteSlice[T]) Contains(element T) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, v := range c.data {
		if v == element {
			return true
		}
	}
	return false
}

// RemoveObject 删除所有匹配元素
func (c *CopyOnWriteSlice[T]) RemoveObject(element T) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 收集保留元素
	newData := make([]T, 0, len(c.data))
	count := 0
	for _, v := range c.data {
		if v != element {
			newData = append(newData, v)
		} else {
			count++
		}
	}
	c.data = newData
	return count
}

// ToSlice 返回数组副本（读操作）
func (c *CopyOnWriteSlice[T]) ToSlice() []T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return append(make([]T, 0, len(c.data)), c.data...)
}

// ToString 返回JSON格式字符串（读操作）
func (c *CopyOnWriteSlice[T]) ToString() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	bytes, err := json.Marshal(c.data)
	if err != nil {
		debug.PrintStack()
		panic(err)
	}
	return string(bytes)
}

// MarshalJSON 实现 json.Marshaler 接口
func (c *CopyOnWriteSlice[T]) MarshalJSON() ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return json.Marshal(c.data[:])
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (c *CopyOnWriteSlice[T]) UnmarshalJSON(data []byte) error {
	var tmp []T
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	c.data = []T{}
	c.AddAll(tmp...)
	return nil
}

func (c *CopyOnWriteSlice[T]) MarshalBSONValue() (bsontype.Type, []byte, error) {
	elements := c.ToSlice()
	return bson.MarshalValue(elements)
}

func (c *CopyOnWriteSlice[T]) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	var elements []T
	if err := bson.UnmarshalValue(t, data, &elements); err != nil {
		return err
	}
	c.data = []T{}
	c.AddAll(elements...)
	return nil
}
