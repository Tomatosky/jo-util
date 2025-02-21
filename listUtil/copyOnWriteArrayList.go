package listUtil

import (
	"encoding/json"
	"sync"
)

// CopyOnWriteArrayList 线程安全的动态数组，写时复制
type CopyOnWriteArrayList[T comparable] struct {
	mu   sync.RWMutex // 读写锁
	data []T          // 实际存储数据的数组
}

// NewCopyOnWriteArrayList 创建新的线程安全动态数组
func NewCopyOnWriteArrayList[T comparable]() *CopyOnWriteArrayList[T] {
	return &CopyOnWriteArrayList[T]{
		data: make([]T, 0),
	}
}

// Add 添加元素到末尾（写操作需要复制整个数组）
func (c *CopyOnWriteArrayList[T]) Add(element T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 创建新数组并追加元素
	newData := append([]T(nil), c.data...)
	newData = append(newData, element)
	c.data = newData
}

// Insert 在指定索引插入元素
func (c *CopyOnWriteArrayList[T]) Insert(index int, element T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if index < 0 || index > len(c.data) {
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
func (c *CopyOnWriteArrayList[T]) Get(index int) T {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if index < 0 {
		index += len(c.data)
	}
	if index < 0 || index >= len(c.data) {
		panic("index out of range")
	}
	return c.data[index]
}

// Remove 移除指定索引元素
func (c *CopyOnWriteArrayList[T]) Remove(index int) T {
	c.mu.Lock()
	defer c.mu.Unlock()

	if index < 0 || index >= len(c.data) {
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
func (c *CopyOnWriteArrayList[T]) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.data)
}

// Range 安全遍历（遍历过程中数组不会被修改）
func (c *CopyOnWriteArrayList[T]) Range(f func(int, T) bool) {
	c.mu.RLock()
	snapshot := append([]T(nil), c.data...) // 创建快照
	c.mu.RUnlock()

	for i, v := range snapshot {
		if !f(i, v) {
			break
		}
	}
}

// Contains 检查元素是否存在（读操作）
func (c *CopyOnWriteArrayList[T]) Contains(element T) bool {
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
func (c *CopyOnWriteArrayList[T]) RemoveObject(element T) int {
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
func (c *CopyOnWriteArrayList[T]) ToSlice() []T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return append([]T(nil), c.data...)
}

// ToString 返回JSON格式字符串（读操作）
func (c *CopyOnWriteArrayList[T]) ToString() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	bytes, err := json.Marshal(c.data)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
