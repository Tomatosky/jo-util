package queueUtil

import "sync"

type SafeFixedQueue[T any] struct {
	items   []T
	head    int
	tail    int
	size    int
	maxSize int
	mutex   sync.Mutex // 互斥锁
}

func NewSafeFixedQueue[T any](capacity int) *SafeFixedQueue[T] {
	return &SafeFixedQueue[T]{
		items:   make([]T, capacity),
		head:    0,
		tail:    0,
		size:    0,
		maxSize: capacity,
	}
}

// Enqueue 安全入队（如果队列已满，返回false）
func (q *SafeFixedQueue[T]) Enqueue(item T) bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.size == q.maxSize {
		return false
	}
	q.items[q.tail] = item
	q.tail = (q.tail + 1) % q.maxSize
	q.size++
	return true
}

// EnqueueForce 安全强制入队（队列满时自动淘汰队首）
func (q *SafeFixedQueue[T]) EnqueueForce(item T) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.size == q.maxSize {
		q.head = (q.head + 1) % q.maxSize
		q.size--
	}
	q.items[q.tail] = item
	q.tail = (q.tail + 1) % q.maxSize
	q.size++
}

// Dequeue 安全出队
func (q *SafeFixedQueue[T]) Dequeue() (T, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.size == 0 {
		var zero T
		return zero, false
	}
	item := q.items[q.head]
	q.head = (q.head + 1) % q.maxSize
	q.size--
	return item, true
}

// Len 安全获取当前元素数量
func (q *SafeFixedQueue[T]) Len() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.size
}

// IsFull 安全检查是否已满
func (q *SafeFixedQueue[T]) IsFull() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.size == q.maxSize
}

// IsEmpty 安全检查是否为空
func (q *SafeFixedQueue[T]) IsEmpty() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.size == 0
}
