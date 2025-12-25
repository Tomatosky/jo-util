package queueUtil

import (
	"testing"
)

// TestNewQueue 测试队列创建
func TestNewQueue(t *testing.T) {
	q := NewQueue[int]()
	if q == nil {
		t.Fatal("NewQueue should not return nil")
	}
	if q.Len() != 0 {
		t.Errorf("New queue should have length 0, got %d", q.Len())
	}
}

// TestEnqueue 测试入队操作
func TestEnqueue(t *testing.T) {
	q := NewQueue[int]()

	// 测试单个元素入队
	q.Enqueue(1)
	if q.Len() != 1 {
		t.Errorf("After one enqueue, length should be 1, got %d", q.Len())
	}

	// 测试多个元素入队
	q.Enqueue(2)
	q.Enqueue(3)
	if q.Len() != 3 {
		t.Errorf("After three enqueues, length should be 3, got %d", q.Len())
	}
}

// TestDequeue 测试出队操作
func TestDequeue(t *testing.T) {
	q := NewQueue[int]()

	// 测试空队列出队
	_, ok := q.Dequeue()
	if ok {
		t.Error("Dequeue from empty queue should return false")
	}

	// 测试正常出队
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	val, ok := q.Dequeue()
	if !ok {
		t.Error("Dequeue should succeed")
	}
	if val != 1 {
		t.Errorf("Expected 1, got %d", val)
	}
	if q.Len() != 2 {
		t.Errorf("After one dequeue, length should be 2, got %d", q.Len())
	}

	val, ok = q.Dequeue()
	if !ok {
		t.Error("Dequeue should succeed")
	}
	if val != 2 {
		t.Errorf("Expected 2, got %d", val)
	}
}

// TestFIFO 测试先进先出特性
func TestFIFO(t *testing.T) {
	q := NewQueue[int]()

	// 入队 1-10
	for i := 1; i <= 10; i++ {
		q.Enqueue(i)
	}

	// 出队并验证顺序
	for i := 1; i <= 10; i++ {
		val, ok := q.Dequeue()
		if !ok {
			t.Fatalf("Dequeue failed at position %d", i)
		}
		if val != i {
			t.Errorf("Expected %d, got %d", i, val)
		}
	}

	// 验证队列为空
	if q.Len() != 0 {
		t.Errorf("Queue should be empty, got length %d", q.Len())
	}
}

// TestShrink 测试缩容机制
func TestShrink(t *testing.T) {
	q := NewQueue[int]()

	// 入队 200 个元素
	for i := 0; i < 200; i++ {
		q.Enqueue(i)
	}

	// 出队 150 个元素,应该触发缩容
	for i := 0; i < 150; i++ {
		val, ok := q.Dequeue()
		if !ok {
			t.Fatalf("Dequeue failed at position %d", i)
		}
		if val != i {
			t.Errorf("Expected %d, got %d", i, val)
		}
	}

	// 验证剩余元素数量
	if q.Len() != 50 {
		t.Errorf("Expected 50 remaining items, got %d", q.Len())
	}

	// 验证剩余元素的正确性
	for i := 150; i < 200; i++ {
		val, ok := q.Dequeue()
		if !ok {
			t.Fatalf("Dequeue failed at position %d", i)
		}
		if val != i {
			t.Errorf("Expected %d, got %d", i, val)
		}
	}
}

// TestMixedOperations 测试混合操作
func TestMixedOperations(t *testing.T) {
	q := NewQueue[int]()

	// 交替入队和出队
	q.Enqueue(1)
	q.Enqueue(2)

	val, ok := q.Dequeue()
	if !ok || val != 1 {
		t.Errorf("Expected 1, got %d, ok=%v", val, ok)
	}

	q.Enqueue(3)
	q.Enqueue(4)

	val, ok = q.Dequeue()
	if !ok || val != 2 {
		t.Errorf("Expected 2, got %d, ok=%v", val, ok)
	}

	val, ok = q.Dequeue()
	if !ok || val != 3 {
		t.Errorf("Expected 3, got %d, ok=%v", val, ok)
	}

	val, ok = q.Dequeue()
	if !ok || val != 4 {
		t.Errorf("Expected 4, got %d, ok=%v", val, ok)
	}

	// 队列应该为空
	if q.Len() != 0 {
		t.Errorf("Queue should be empty, got length %d", q.Len())
	}
}

// TestZeroValue 测试零值类型
func TestZeroValue(t *testing.T) {
	q := NewQueue[int]()

	// 入队零值
	q.Enqueue(0)

	val, ok := q.Dequeue()
	if !ok {
		t.Error("Dequeue should succeed")
	}
	if val != 0 {
		t.Errorf("Expected 0, got %d", val)
	}
}

// TestStringQueue 测试字符串类型队列
func TestStringQueue(t *testing.T) {
	q := NewQueue[string]()

	q.Enqueue("hello")
	q.Enqueue("world")

	val, ok := q.Dequeue()
	if !ok || val != "hello" {
		t.Errorf("Expected 'hello', got '%s', ok=%v", val, ok)
	}

	val, ok = q.Dequeue()
	if !ok || val != "world" {
		t.Errorf("Expected 'world', got '%s', ok=%v", val, ok)
	}
}

// TestStructQueue 测试结构体类型队列
func TestStructQueue(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	q := NewQueue[Person]()

	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Bob", Age: 25}

	q.Enqueue(p1)
	q.Enqueue(p2)

	val, ok := q.Dequeue()
	if !ok {
		t.Fatal("Dequeue should succeed")
	}
	if val.Name != "Alice" || val.Age != 30 {
		t.Errorf("Expected Alice, 30, got %s, %d", val.Name, val.Age)
	}

	val, ok = q.Dequeue()
	if !ok {
		t.Fatal("Dequeue should succeed")
	}
	if val.Name != "Bob" || val.Age != 25 {
		t.Errorf("Expected Bob, 25, got %s, %d", val.Name, val.Age)
	}
}

// TestLargeQueue 测试大量数据
func TestLargeQueue(t *testing.T) {
	q := NewQueue[int]()

	const size = 10000

	// 入队大量数据
	for i := 0; i < size; i++ {
		q.Enqueue(i)
	}

	if q.Len() != size {
		t.Errorf("Expected length %d, got %d", size, q.Len())
	}

	// 出队并验证
	for i := 0; i < size; i++ {
		val, ok := q.Dequeue()
		if !ok {
			t.Fatalf("Dequeue failed at position %d", i)
		}
		if val != i {
			t.Errorf("Expected %d, got %d", i, val)
		}
	}

	if q.Len() != 0 {
		t.Errorf("Queue should be empty, got length %d", q.Len())
	}
}

// TestEmptyQueueLen 测试空队列的长度
func TestEmptyQueueLen(t *testing.T) {
	q := NewQueue[int]()

	if q.Len() != 0 {
		t.Errorf("Empty queue should have length 0, got %d", q.Len())
	}

	// 入队后出队
	q.Enqueue(1)
	q.Dequeue()

	if q.Len() != 0 {
		t.Errorf("Queue should be empty after dequeue, got length %d", q.Len())
	}
}

// TestDequeueAll 测试出队所有元素后继续出队
func TestDequeueAll(t *testing.T) {
	q := NewQueue[int]()

	q.Enqueue(1)
	q.Enqueue(2)

	q.Dequeue()
	q.Dequeue()

	// 尝试从空队列出队
	_, ok := q.Dequeue()
	if ok {
		t.Error("Dequeue from empty queue should return false")
	}

	// 再次入队应该正常工作
	q.Enqueue(3)
	val, ok := q.Dequeue()
	if !ok || val != 3 {
		t.Errorf("Expected 3, got %d, ok=%v", val, ok)
	}
}

// TestPointerQueue 测试指针类型队列
func TestPointerQueue(t *testing.T) {
	q := NewQueue[*int]()

	val1 := 42
	val2 := 100

	q.Enqueue(&val1)
	q.Enqueue(&val2)

	p1, ok := q.Dequeue()
	if !ok {
		t.Fatal("Dequeue should succeed")
	}
	if *p1 != 42 {
		t.Errorf("Expected 42, got %d", *p1)
	}

	p2, ok := q.Dequeue()
	if !ok {
		t.Fatal("Dequeue should succeed")
	}
	if *p2 != 100 {
		t.Errorf("Expected 100, got %d", *p2)
	}
}

// BenchmarkEnqueue 性能测试:入队
func BenchmarkEnqueue(b *testing.B) {
	q := NewQueue[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
}

// BenchmarkDequeue 性能测试:出队
func BenchmarkDequeue(b *testing.B) {
	q := NewQueue[int]()
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Dequeue()
	}
}

// BenchmarkMixed 性能测试:混合操作
func BenchmarkMixed(b *testing.B) {
	q := NewQueue[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
		if i%2 == 0 {
			q.Dequeue()
		}
	}
}
