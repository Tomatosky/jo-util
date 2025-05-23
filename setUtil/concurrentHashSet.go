package setUtil

import (
	"encoding/json"
	"github.com/Tomatosky/jo-util/mapUtil"
)

// ConcurrentHashSet 基于 ConcurrentHashMap 实现的并发安全集合
type ConcurrentHashSet[T comparable] struct {
	m *mapUtil.ConcurrentHashMap[T, struct{}] // 使用空结构体作为值类型
}

// NewConcurrentHashSet 构造函数
func NewConcurrentHashSet[T comparable](elements ...T) *ConcurrentHashSet[T] {
	set := &ConcurrentHashSet[T]{
		m: mapUtil.NewConcurrentHashMap[T, struct{}](), // 初始化底层Map
	}
	set.AddAll(elements...)
	return set
}

// Add 添加元素
func (s *ConcurrentHashSet[T]) Add(element T) {
	s.m.Put(element, struct{}{})
}

// AddAll 批量添加元素
func (s *ConcurrentHashSet[T]) AddAll(elements ...T) {
	for _, e := range elements {
		s.Add(e)
	}
}

// Remove 移除元素
func (s *ConcurrentHashSet[T]) Remove(element T) {
	s.m.Remove(element)
}

// Contains 检查元素存在性
func (s *ConcurrentHashSet[T]) Contains(element T) bool {
	return s.m.ContainsKey(element)
}

// Size 获取元素数量
func (s *ConcurrentHashSet[T]) Size() int {
	return s.m.Size()
}

// Clear 清空集合
func (s *ConcurrentHashSet[T]) Clear() {
	s.m.Clear()
}

// ToSlice 转换为切片
func (s *ConcurrentHashSet[T]) ToSlice() []T {
	return s.m.Keys()
}

// IsEmpty 判断是否为空
func (s *ConcurrentHashSet[T]) IsEmpty() bool {
	return s.Size() == 0
}

func (s *ConcurrentHashSet[T]) ToString() string {
	bytes, err := json.Marshal(s.ToSlice())
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (s *ConcurrentHashSet[T]) Range(f func(T) bool) {
	s.m.ra
}
