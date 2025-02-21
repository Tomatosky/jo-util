package listutil

import "encoding/json"

type ArrayList[T comparable] struct {
	data []T
}

// NewArrayList 创建初始容量为 capacity 的动态数组
func NewArrayList[T comparable]() *ArrayList[T] {
	return &ArrayList[T]{
		data: make([]T, 0),
	}
}

// Add 添加元素到末尾
func (a *ArrayList[T]) Add(element T) {
	a.data = append(a.data, element)
}

// Insert 在指定索引插入元素
func (a *ArrayList[T]) Insert(index int, element T) {
	if index < 0 || index > len(a.data) {
		panic("index out of range")
	}
	a.data = append(a.data, element)
	copy(a.data[index+1:], a.data[index:])
	a.data[index] = element
}

// Get 获取元素（支持负数索引，-1表示最后一个元素）
func (a *ArrayList[T]) Get(index int) T {
	if index < 0 {
		index += len(a.data)
	}
	if index < 0 || index >= len(a.data) {
		panic("index out of range")
	}
	return a.data[index]
}

// Remove 移除指定索引元素
func (a *ArrayList[T]) Remove(index int) T {
	if index < 0 || index >= len(a.data) {
		panic("index out of range")
	}
	removed := a.data[index]
	a.data = append(a.data[:index], a.data[index+1:]...)
	return removed
}

// Size 当前元素数量
func (a *ArrayList[T]) Size() int {
	return len(a.data)
}

// Range 遍历元素（返回false可终止遍历）
func (a *ArrayList[T]) Range(f func(int, T) bool) {
	for i, v := range a.data {
		if !f(i, v) {
			break
		}
	}
}

// Contains 检查元素是否存在
func (a *ArrayList[T]) Contains(element T) bool {
	for _, v := range a.data {
		if v == element {
			return true
		}
	}
	return false
}

// RemoveObject 删除匹配的元素
func (a *ArrayList[T]) RemoveObject(element T) int {
	count := 0
	// 反向遍历避免索引错位
	for i := len(a.data) - 1; i >= 0; i-- {
		if a.data[i] == element {
			a.data = append(a.data[:i], a.data[i+1:]...)
			count++
		}
	}
	return count // 返回实际删除数量
}

func (a *ArrayList[T]) ToSlice() []T {
	return a.data
}

func (a *ArrayList[T]) ToString() string {
	bytes, err := json.Marshal(a.data)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
