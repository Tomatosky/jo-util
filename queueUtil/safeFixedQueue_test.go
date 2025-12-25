package queueUtil

import (
	"sync"
	"testing"
)

// TestNewSafeFixedQueue 测试队列创建
func TestNewSafeFixedQueue(t *testing.T) {
	q := NewSafeFixedQueue[int](10)
	if q == nil {
		t.Fatal("NewSafeFixedQueue should not return nil")
	}
	if q.Len() != 0 {
		t.Errorf("New queue should have length 0, got %d", q.Len())
	}
	if q.maxSize != 10 {
		t.Errorf("Queue max size should be 10, got %d", q.maxSize)
	}
	if !q.IsEmpty() {
		t.Error("New queue should be empty")
	}
	if q.IsFull() {
		t.Error("New queue should not be full")
	}
}

// TestSafeFixedQueueEnqueue 测试普通入队操作
func TestSafeFixedQueueEnqueue(t *testing.T) {
	q := NewSafeFixedQueue[int](3)

	// 测试入队成功
	if !q.Enqueue(1) {
		t.Error("First enqueue should succeed")
	}
	if q.Len() != 1 {
		t.Errorf("After one enqueue, length should be 1, got %d", q.Len())
	}

	// 继续入队
	if !q.Enqueue(2) {
		t.Error("Second enqueue should succeed")
	}
	if !q.Enqueue(3) {
		t.Error("Third enqueue should succeed")
	}

	if q.Len() != 3 {
		t.Errorf("After three enqueues, length should be 3, got %d", q.Len())
	}

	// 队列已满,入队应该失败
	if q.Enqueue(4) {
		t.Error("Enqueue to full queue should fail")
	}
	if q.Len() != 3 {
		t.Errorf("Length should remain 3 after failed enqueue, got %d", q.Len())
	}
}

// TestSafeFixedQueueEnqueueForce 测试强制入队
func TestSafeFixedQueueEnqueueForce(t *testing.T) {
	q := NewSafeFixedQueue[int](3)

	// 先填满队列
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	if !q.IsFull() {
		t.Error("Queue should be full")
	}

	// 强制入队,应该淘汰队首元素1
	q.EnqueueForce(4)

	if q.Len() != 3 {
		t.Errorf("After force enqueue, length should still be 3, got %d", q.Len())
	}

	// 验证队首元素已被淘汰,现在应该是2
	val, ok := q.Dequeue()
	if !ok {
		t.Fatal("Dequeue should succeed")
	}
	if val != 2 {
		t.Errorf("Expected 2 (first element should be evicted), got %d", val)
	}

	val, ok = q.Dequeue()
	if !ok {
		t.Fatal("Dequeue should succeed")
	}
	if val != 3 {
		t.Errorf("Expected 3, got %d", val)
	}

	val, ok = q.Dequeue()
	if !ok {
		t.Fatal("Dequeue should succeed")
	}
	if val != 4 {
		t.Errorf("Expected 4, got %d", val)
	}
}

// TestSafeFixedQueueEnqueueForceMultiple 测试多次强制入队
func TestSafeFixedQueueEnqueueForceMultiple(t *testing.T) {
	q := NewSafeFixedQueue[int](3)

	// 填满队列: 1, 2, 3
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	// 多次强制入队
	q.EnqueueForce(4) // 队列: 2, 3, 4
	q.EnqueueForce(5) // 队列: 3, 4, 5
	q.EnqueueForce(6) // 队列: 4, 5, 6

	if q.Len() != 3 {
		t.Errorf("Length should be 3, got %d", q.Len())
	}

	// 验证队列内容
	expected := []int{4, 5, 6}
	for i, exp := range expected {
		val, ok := q.Dequeue()
		if !ok {
			t.Fatalf("Dequeue %d should succeed", i)
		}
		if val != exp {
			t.Errorf("Position %d: expected %d, got %d", i, exp, val)
		}
	}
}

// TestSafeFixedQueueDequeue 测试出队操作
func TestSafeFixedQueueDequeue(t *testing.T) {
	q := NewSafeFixedQueue[int](5)

	// 测试空队列出队
	_, ok := q.Dequeue()
	if ok {
		t.Error("Dequeue from empty queue should return false")
	}

	// 入队后出队
	q.Enqueue(10)
	q.Enqueue(20)
	q.Enqueue(30)

	val, ok := q.Dequeue()
	if !ok {
		t.Error("Dequeue should succeed")
	}
	if val != 10 {
		t.Errorf("Expected 10, got %d", val)
	}

	if q.Len() != 2 {
		t.Errorf("After one dequeue, length should be 2, got %d", q.Len())
	}

	val, ok = q.Dequeue()
	if !ok {
		t.Error("Dequeue should succeed")
	}
	if val != 20 {
		t.Errorf("Expected 20, got %d", val)
	}

	val, ok = q.Dequeue()
	if !ok {
		t.Error("Dequeue should succeed")
	}
	if val != 30 {
		t.Errorf("Expected 30, got %d", val)
	}

	if !q.IsEmpty() {
		t.Error("Queue should be empty after dequeuing all elements")
	}
}

// TestSafeFixedQueueFIFO 测试先进先出特性
func TestSafeFixedQueueFIFO(t *testing.T) {
	q := NewSafeFixedQueue[int](10)

	// 入队 1-10
	for i := 1; i <= 10; i++ {
		if !q.Enqueue(i) {
			t.Fatalf("Enqueue %d failed", i)
		}
	}

	if !q.IsFull() {
		t.Error("Queue should be full")
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

	if !q.IsEmpty() {
		t.Error("Queue should be empty")
	}
}

// TestSafeFixedQueueCircular 测试循环队列特性
func TestSafeFixedQueueCircular(t *testing.T) {
	q := NewSafeFixedQueue[int](3)

	// 第一轮:填满队列
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	// 出队两个
	q.Dequeue() // 移除1
	q.Dequeue() // 移除2

	// 再入队两个,测试循环
	q.Enqueue(4)
	q.Enqueue(5)

	// 现在队列应该是: 3, 4, 5
	expected := []int{3, 4, 5}
	for i, exp := range expected {
		val, ok := q.Dequeue()
		if !ok {
			t.Fatalf("Dequeue %d should succeed", i)
		}
		if val != exp {
			t.Errorf("Position %d: expected %d, got %d", i, exp, val)
		}
	}
}

// TestSafeFixedQueueIsFullIsEmpty 测试满和空状态检查
func TestSafeFixedQueueIsFullIsEmpty(t *testing.T) {
	q := NewSafeFixedQueue[int](2)

	// 初始状态
	if !q.IsEmpty() {
		t.Error("New queue should be empty")
	}
	if q.IsFull() {
		t.Error("New queue should not be full")
	}

	// 入队一个元素
	q.Enqueue(1)
	if q.IsEmpty() {
		t.Error("Queue with one element should not be empty")
	}
	if q.IsFull() {
		t.Error("Queue with one element should not be full (capacity is 2)")
	}

	// 填满队列
	q.Enqueue(2)
	if q.IsEmpty() {
		t.Error("Full queue should not be empty")
	}
	if !q.IsFull() {
		t.Error("Queue should be full")
	}

	// 出队一个元素
	q.Dequeue()
	if q.IsEmpty() {
		t.Error("Queue with one element should not be empty")
	}
	if q.IsFull() {
		t.Error("Queue with one element should not be full")
	}

	// 清空队列
	q.Dequeue()
	if !q.IsEmpty() {
		t.Error("Queue should be empty")
	}
	if q.IsFull() {
		t.Error("Empty queue should not be full")
	}
}

// TestSafeFixedQueueLen 测试长度方法
func TestSafeFixedQueueLen(t *testing.T) {
	q := NewSafeFixedQueue[int](5)

	lengths := []int{0, 1, 2, 3, 4, 5}
	for _, expected := range lengths[:len(lengths)-1] {
		if q.Len() != expected {
			t.Errorf("Expected length %d, got %d", expected, q.Len())
		}
		q.Enqueue(expected)
	}

	// 最后检查填满后的长度
	if q.Len() != 5 {
		t.Errorf("Expected length 5, got %d", q.Len())
	}

	// 测试出队后长度变化
	for i := 5; i > 0; i-- {
		if q.Len() != i {
			t.Errorf("Expected length %d, got %d", i, q.Len())
		}
		q.Dequeue()
	}

	if q.Len() != 0 {
		t.Errorf("Expected length 0, got %d", q.Len())
	}
}

// TestSafeFixedQueueConcurrentEnqueue 测试并发入队
func TestSafeFixedQueueConcurrentEnqueue(t *testing.T) {
	q := NewSafeFixedQueue[int](1000)
	var wg sync.WaitGroup
	goroutines := 10
	itemsPerGoroutine := 100

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			for j := 0; j < itemsPerGoroutine; j++ {
				q.Enqueue(start*itemsPerGoroutine + j)
			}
		}(i)
	}

	wg.Wait()

	// 验证所有元素都已入队
	if q.Len() != goroutines*itemsPerGoroutine {
		t.Errorf("Expected %d elements, got %d", goroutines*itemsPerGoroutine, q.Len())
	}
}

// TestSafeFixedQueueConcurrentDequeue 测试并发出队
func TestSafeFixedQueueConcurrentDequeue(t *testing.T) {
	q := NewSafeFixedQueue[int](1000)
	totalItems := 1000

	// 先填充队列
	for i := 0; i < totalItems; i++ {
		q.Enqueue(i)
	}

	var wg sync.WaitGroup
	goroutines := 10
	successCount := make([]int, goroutines)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			count := 0
			for {
				_, ok := q.Dequeue()
				if !ok {
					break
				}
				count++
			}
			successCount[id] = count
		}(i)
	}

	wg.Wait()

	// 验证所有元素都已出队
	total := 0
	for _, count := range successCount {
		total += count
	}
	if total != totalItems {
		t.Errorf("Expected %d items dequeued, got %d", totalItems, total)
	}

	if !q.IsEmpty() {
		t.Errorf("Queue should be empty, but has %d elements", q.Len())
	}
}

// TestSafeFixedQueueConcurrentMixed 测试并发混合操作
func TestSafeFixedQueueConcurrentMixed(t *testing.T) {
	q := NewSafeFixedQueue[int](100)
	var wg sync.WaitGroup
	iterations := 1000

	// 启动多个生产者
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				q.Enqueue(start*iterations + j)
			}
		}(i)
	}

	// 启动多个消费者
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				q.Dequeue()
			}
		}()
	}

	wg.Wait()

	// 最终队列应该是空的(相同数量的入队和出队)
	if !q.IsEmpty() {
		t.Logf("Queue has %d elements remaining (this is normal in concurrent scenario)", q.Len())
	}
}

// TestSafeFixedQueueConcurrentForceEnqueue 测试并发强制入队
func TestSafeFixedQueueConcurrentForceEnqueue(t *testing.T) {
	q := NewSafeFixedQueue[int](50)
	var wg sync.WaitGroup
	goroutines := 10
	itemsPerGoroutine := 100

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			for j := 0; j < itemsPerGoroutine; j++ {
				q.EnqueueForce(start*itemsPerGoroutine + j)
			}
		}(i)
	}

	wg.Wait()

	// 队列应该是满的
	if !q.IsFull() {
		t.Error("Queue should be full after concurrent force enqueue")
	}

	// 长度应该等于容量
	if q.Len() != 50 {
		t.Errorf("Expected length 50, got %d", q.Len())
	}
}

// TestSafeFixedQueueZeroValue 测试零值
func TestSafeFixedQueueZeroValue(t *testing.T) {
	q := NewSafeFixedQueue[int](5)

	// 入队零值
	q.Enqueue(0)
	q.Enqueue(1)
	q.Enqueue(0)

	val, ok := q.Dequeue()
	if !ok {
		t.Fatal("Dequeue should succeed")
	}
	if val != 0 {
		t.Errorf("Expected 0, got %d", val)
	}

	val, ok = q.Dequeue()
	if !ok {
		t.Fatal("Dequeue should succeed")
	}
	if val != 1 {
		t.Errorf("Expected 1, got %d", val)
	}

	val, ok = q.Dequeue()
	if !ok {
		t.Fatal("Dequeue should succeed")
	}
	if val != 0 {
		t.Errorf("Expected 0, got %d", val)
	}
}

// TestSafeFixedQueueString 测试字符串类型
func TestSafeFixedQueueString(t *testing.T) {
	q := NewSafeFixedQueue[string](3)

	q.Enqueue("hello")
	q.Enqueue("world")
	q.Enqueue("!")

	// 队列满时入队应该失败
	if q.Enqueue("overflow") {
		t.Error("Enqueue to full queue should fail")
	}

	val, ok := q.Dequeue()
	if !ok || val != "hello" {
		t.Errorf("Expected 'hello', got '%s', ok=%v", val, ok)
	}

	// 现在队列有2个元素: world, !
	// 强制入队不会淘汰元素(因为队列未满)
	q.EnqueueForce("forced")

	// 验证顺序: world, !, forced
	expected := []string{"world", "!", "forced"}
	for i, exp := range expected {
		val, ok := q.Dequeue()
		if !ok {
			t.Fatalf("Dequeue %d should succeed", i)
		}
		if val != exp {
			t.Errorf("Position %d: expected '%s', got '%s'", i, exp, val)
		}
	}
}

// TestSafeFixedQueueStruct 测试结构体类型
func TestSafeFixedQueueStruct(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	q := NewSafeFixedQueue[Person](3)

	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Bob", Age: 25}
	p3 := Person{Name: "Charlie", Age: 35}

	q.Enqueue(p1)
	q.Enqueue(p2)
	q.Enqueue(p3)

	if !q.IsFull() {
		t.Error("Queue should be full")
	}

	val, ok := q.Dequeue()
	if !ok {
		t.Fatal("Dequeue should succeed")
	}
	if val.Name != "Alice" || val.Age != 30 {
		t.Errorf("Expected Alice, 30, got %s, %d", val.Name, val.Age)
	}

	// 现在队列有2个元素: Bob, Charlie
	// 强制入队第四个人,不会淘汰(因为队列未满)
	p4 := Person{Name: "David", Age: 40}
	q.EnqueueForce(p4)

	// 验证顺序: Bob, Charlie, David
	val, ok = q.Dequeue()
	if !ok {
		t.Fatal("Dequeue should succeed")
	}
	if val.Name != "Bob" || val.Age != 25 {
		t.Errorf("Expected Bob, 25, got %s, %d", val.Name, val.Age)
	}

	val, ok = q.Dequeue()
	if !ok {
		t.Fatal("Dequeue should succeed")
	}
	if val.Name != "Charlie" || val.Age != 35 {
		t.Errorf("Expected Charlie, 35, got %s, %d", val.Name, val.Age)
	}

	val, ok = q.Dequeue()
	if !ok {
		t.Fatal("Dequeue should succeed")
	}
	if val.Name != "David" || val.Age != 40 {
		t.Errorf("Expected David, 40, got %s, %d", val.Name, val.Age)
	}
}

// TestSafeFixedQueuePointer 测试指针类型
func TestSafeFixedQueuePointer(t *testing.T) {
	q := NewSafeFixedQueue[*int](3)

	val1 := 42
	val2 := 100
	val3 := 200

	q.Enqueue(&val1)
	q.Enqueue(&val2)
	q.Enqueue(&val3)

	p1, ok := q.Dequeue()
	if !ok {
		t.Fatal("Dequeue should succeed")
	}
	if *p1 != 42 {
		t.Errorf("Expected 42, got %d", *p1)
	}

	// 修改原始值,验证指针引用
	val1 = 99
	val4 := 300
	// 现在队列有2个元素: &val2, &val3
	// 强制入队不会淘汰(因为队列未满)
	q.EnqueueForce(&val4)

	// 验证顺序: &val2(100), &val3(200), &val4(300)
	p2, ok := q.Dequeue()
	if !ok {
		t.Fatal("Dequeue should succeed")
	}
	if *p2 != 100 {
		t.Errorf("Expected 100, got %d", *p2)
	}

	p3, ok := q.Dequeue()
	if !ok {
		t.Fatal("Dequeue should succeed")
	}
	if *p3 != 200 {
		t.Errorf("Expected 200, got %d", *p3)
	}

	p4, ok := q.Dequeue()
	if !ok {
		t.Fatal("Dequeue should succeed")
	}
	if *p4 != 300 {
		t.Errorf("Expected 300, got %d", *p4)
	}
}

// TestSafeFixedQueueCapacityOne 测试容量为1的队列
func TestSafeFixedQueueCapacityOne(t *testing.T) {
	q := NewSafeFixedQueue[int](1)

	if !q.IsEmpty() {
		t.Error("New queue should be empty")
	}

	// 入队
	if !q.Enqueue(1) {
		t.Error("First enqueue should succeed")
	}

	if !q.IsFull() {
		t.Error("Queue should be full")
	}

	// 再次入队应该失败
	if q.Enqueue(2) {
		t.Error("Enqueue to full queue should fail")
	}

	// 强制入队应该淘汰第一个元素
	q.EnqueueForce(3)

	val, ok := q.Dequeue()
	if !ok {
		t.Fatal("Dequeue should succeed")
	}
	if val != 3 {
		t.Errorf("Expected 3, got %d", val)
	}

	if !q.IsEmpty() {
		t.Error("Queue should be empty")
	}
}

// TestSafeFixedQueueMixedOperations 测试混合操作
func TestSafeFixedQueueMixedOperations(t *testing.T) {
	q := NewSafeFixedQueue[int](5)

	// 场景1: 正常入队出队
	q.Enqueue(1)
	q.Enqueue(2)
	val, _ := q.Dequeue()
	if val != 1 {
		t.Errorf("Expected 1, got %d", val)
	}

	// 场景2: 填满队列
	q.Enqueue(3)
	q.Enqueue(4)
	q.Enqueue(5)
	q.Enqueue(6)

	if !q.IsFull() {
		t.Error("Queue should be full")
	}

	// 场景3: 强制入队
	q.EnqueueForce(7)

	// 验证: 应该是 3, 4, 5, 6, 7 (2被淘汰了)
	expected := []int{3, 4, 5, 6, 7}
	for i, exp := range expected {
		val, ok := q.Dequeue()
		if !ok {
			t.Fatalf("Dequeue %d should succeed", i)
		}
		if val != exp {
			t.Errorf("Position %d: expected %d, got %d", i, exp, val)
		}
	}

	// 场景4: 队列再次为空
	if !q.IsEmpty() {
		t.Error("Queue should be empty")
	}

	// 场景5: 重新使用队列
	q.Enqueue(10)
	val, ok := q.Dequeue()
	if !ok || val != 10 {
		t.Errorf("Expected 10, got %d, ok=%v", val, ok)
	}
}

// TestSafeFixedQueueEmptyDequeue 测试连续从空队列出队
func TestSafeFixedQueueEmptyDequeue(t *testing.T) {
	q := NewSafeFixedQueue[int](3)

	// 多次从空队列出队
	for i := 0; i < 10; i++ {
		_, ok := q.Dequeue()
		if ok {
			t.Errorf("Iteration %d: dequeue from empty queue should fail", i)
		}
	}

	if q.Len() != 0 {
		t.Errorf("Length should be 0, got %d", q.Len())
	}
}

// TestSafeFixedQueueFullEnqueue 测试连续向满队列入队
func TestSafeFixedQueueFullEnqueue(t *testing.T) {
	q := NewSafeFixedQueue[int](3)

	// 填满队列
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	// 多次向满队列入队
	for i := 4; i < 10; i++ {
		if q.Enqueue(i) {
			t.Errorf("Iteration %d: enqueue to full queue should fail", i)
		}
	}

	if q.Len() != 3 {
		t.Errorf("Length should be 3, got %d", q.Len())
	}
}

// BenchmarkSafeFixedQueueEnqueue 性能测试:入队
func BenchmarkSafeFixedQueueEnqueue(b *testing.B) {
	q := NewSafeFixedQueue[int](b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
}

// BenchmarkSafeFixedQueueEnqueueForce 性能测试:强制入队
func BenchmarkSafeFixedQueueEnqueueForce(b *testing.B) {
	q := NewSafeFixedQueue[int](1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.EnqueueForce(i)
	}
}

// BenchmarkSafeFixedQueueDequeue 性能测试:出队
func BenchmarkSafeFixedQueueDequeue(b *testing.B) {
	q := NewSafeFixedQueue[int](b.N)
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Dequeue()
	}
}

// BenchmarkSafeFixedQueueMixed 性能测试:混合操作
func BenchmarkSafeFixedQueueMixed(b *testing.B) {
	q := NewSafeFixedQueue[int](1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
		if i%2 == 0 {
			q.Dequeue()
		}
	}
}

// BenchmarkSafeFixedQueueConcurrent 性能测试:并发操作
func BenchmarkSafeFixedQueueConcurrent(b *testing.B) {
	q := NewSafeFixedQueue[int](10000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				q.Enqueue(i)
			} else {
				q.Dequeue()
			}
			i++
		}
	})
}
