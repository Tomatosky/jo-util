package queueUtil

import "testing"

func TestQueue(t *testing.T) {
	// 测试1: 创建新队列
	q := NewQueue[int]()
	if q == nil {
		t.Fatal("创建队列失败")
	}
	if q.Len() != 0 {
		t.Fatal("新创建的队列长度应该为0")
	}

	// 测试2: 入队操作
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	if q.Len() != 3 {
		t.Errorf("入队后队列长度应该为3，实际为%d", q.Len())
	}

	// 测试3: 出队操作
	item, ok := q.Dequeue()
	if !ok {
		t.Error("出队操作应该成功")
	}
	if item != 1 {
		t.Errorf("出队元素应该为1，实际为%d", item)
	}
	if q.Len() != 2 {
		t.Errorf("出队后队列长度应该为2，实际为%d", q.Len())
	}

	// 测试4: 连续出队
	item, ok = q.Dequeue()
	if !ok || item != 2 {
		t.Error("第二次出队失败")
	}
	item, ok = q.Dequeue()
	if !ok || item != 3 {
		t.Error("第三次出队失败")
	}
	if q.Len() != 0 {
		t.Errorf("全部出队后队列长度应该为0，实际为%d", q.Len())
	}

	// 测试5: 空队列出队
	item, ok = q.Dequeue()
	if ok {
		t.Error("空队列出队应该失败")
	}
	var zero int
	if item != zero {
		t.Error("空队列出队应该返回零值")
	}

	// 测试6: 测试字符串类型
	strQueue := NewQueue[string]()
	strQueue.Enqueue("hello")
	strQueue.Enqueue("world")
	if strQueue.Len() != 2 {
		t.Error("字符串队列长度错误")
	}
	str, ok := strQueue.Dequeue()
	if !ok || str != "hello" {
		t.Error("字符串队列出队失败")
	}

	// 测试7: 测试结构体类型
	type Person struct {
		Name string
		Age  int
	}
	personQueue := NewQueue[Person]()
	personQueue.Enqueue(Person{"Alice", 25})
	personQueue.Enqueue(Person{"Bob", 30})
	if personQueue.Len() != 2 {
		t.Error("结构体队列长度错误")
	}
	person, ok := personQueue.Dequeue()
	if !ok || person.Name != "Alice" || person.Age != 25 {
		t.Error("结构体队列出队失败")
	}

	// 测试8: 测试缩容机制
	// 创建大量元素触发缩容
	largeQueue := NewQueue[int]()
	for i := 0; i < 200; i++ {
		largeQueue.Enqueue(i)
	}
	// 出队101个元素，触发缩容条件
	for i := 0; i < 101; i++ {
		largeQueue.Dequeue()
	}
	// 检查队列长度是否正确
	if largeQueue.Len() != 99 {
		t.Errorf("缩容后队列长度应该为99，实际为%d", largeQueue.Len())
	}
	// 检查下一个出队元素是否正确
	item, ok = largeQueue.Dequeue()
	if !ok || item != 101 {
		t.Errorf("缩容后出队元素应该为101，实际为%d", item)
	}

	// 测试9: 边界情况测试
	emptyQueue := NewQueue[bool]()
	if emptyQueue.Len() != 0 {
		t.Error("空队列长度应该为0")
	}

	// 测试10: 混合操作测试
	mixedQueue := NewQueue[int]()
	for i := 0; i < 10; i++ {
		mixedQueue.Enqueue(i)
	}
	for i := 0; i < 5; i++ {
		mixedQueue.Dequeue()
	}
	for i := 10; i < 15; i++ {
		mixedQueue.Enqueue(i)
	}
	if mixedQueue.Len() != 10 {
		t.Errorf("混合操作后队列长度应该为10，实际为%d", mixedQueue.Len())
	}
}
