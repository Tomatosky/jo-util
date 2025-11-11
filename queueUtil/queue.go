package queueUtil

type Queue[T any] struct {
	items []T
	head  int // 队首指针
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		items: make([]T, 0),
		head:  0,
	}
}

// Enqueue 入队
func (q *Queue[T]) Enqueue(item T) {
	q.items = append(q.items, item)
}

// Dequeue 出队
func (q *Queue[T]) Dequeue() (T, bool) {
	if q.head >= len(q.items) {
		var zero T
		return zero, false
	}
	item := q.items[q.head]
	q.head++

	// 定期缩容，避免内存泄漏
	if q.head > 100 && q.head > len(q.items)/2 {
		q.items = q.items[q.head:]
		q.head = 0
	}

	return item, true
}

// Len 队列长度
func (q *Queue[T]) Len() int {
	return len(q.items) - q.head
}
