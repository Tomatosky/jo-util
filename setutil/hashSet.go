package setutil

import "encoding/json"

// HashSet 非并发安全的哈希集合实现
type HashSet[T comparable] struct {
	m map[T]struct{} // 使用空结构体作为值类型（0内存占用）
}

// NewHashSet 构造函数
func NewHashSet[T comparable](elements ...T) *HashSet[T] {
	set := &HashSet[T]{
		m: make(map[T]struct{}),
	}
	set.AddAll(elements...)
	return set
}

// Add 添加元素
func (s *HashSet[T]) Add(element T) {
	s.m[element] = struct{}{}
}

// AddAll 批量添加元素
func (s *HashSet[T]) AddAll(elements ...T) {
	for _, e := range elements {
		s.Add(e)
	}
}

// Remove 移除元素
func (s *HashSet[T]) Remove(element T) {
	delete(s.m, element)
}

// Contains 检查元素存在性
func (s *HashSet[T]) Contains(element T) bool {
	_, exists := s.m[element]
	return exists
}

// Size 获取元素数量
func (s *HashSet[T]) Size() int {
	return len(s.m)
}

// Clear 清空集合
func (s *HashSet[T]) Clear() {
	s.m = make(map[T]struct{})
}

// Range 遍历元素（返回false可提前终止）
func (s *HashSet[T]) Range(f func(T) bool) {
	for k := range s.m {
		if !f(k) {
			break
		}
	}
}

// ToSlice 转换为切片
func (s *HashSet[T]) ToSlice() []T {
	slice := make([]T, 0, len(s.m))
	for k := range s.m {
		slice = append(slice, k)
	}
	return slice
}

// IsEmpty 判断是否为空
func (s *HashSet[T]) IsEmpty() bool {
	return len(s.m) == 0
}

func (s *HashSet[T]) ToString() string {
	bytes, err := json.Marshal(s.ToSlice())
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
