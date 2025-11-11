package queueUtil

import (
	"testing"
)

func TestSafeFixedQueue(t *testing.T) {
	// 测试1: 创建队列
	t.Run("创建队列", func(t *testing.T) {
		capacity := 5
		q := NewSafeFixedQueue[int](capacity)

		if q == nil {
			t.Error("队列创建失败")
		}
		if q.maxSize != capacity {
			t.Errorf("期望容量为 %d，实际为 %d", capacity, q.maxSize)
		}
		if q.Len() != 0 {
			t.Errorf("新队列长度应为0，实际为 %d", q.Len())
		}
		if !q.IsEmpty() {
			t.Error("新队列应为空")
		}
		if q.IsFull() {
			t.Error("新队列不应为满")
		}
	})

	// 测试2: 正常入队出队操作
	t.Run("正常入队出队", func(t *testing.T) {
		q := NewSafeFixedQueue[int](3)

		// 入队测试
		if !q.Enqueue(1) {
			t.Error("入队失败，期望成功")
		}
		if q.Len() != 1 {
			t.Errorf("队列长度应为1，实际为 %d", q.Len())
		}

		if !q.Enqueue(2) {
			t.Error("入队失败，期望成功")
		}
		if q.Len() != 2 {
			t.Errorf("队列长度应为2，实际为 %d", q.Len())
		}

		// 出队测试
		item, ok := q.Dequeue()
		if !ok {
			t.Error("出队失败，期望成功")
		}
		if item != 1 {
			t.Errorf("出队元素应为1，实际为 %d", item)
		}
		if q.Len() != 1 {
			t.Errorf("队列长度应为1，实际为 %d", q.Len())
		}

		item, ok = q.Dequeue()
		if !ok {
			t.Error("出队失败，期望成功")
		}
		if item != 2 {
			t.Errorf("出队元素应为2，实际为 %d", item)
		}
		if !q.IsEmpty() {
			t.Error("队列应为空")
		}
	})

	// 测试3: 队列满时入队
	t.Run("队列满测试", func(t *testing.T) {
		q := NewSafeFixedQueue[int](2)

		// 填满队列
		if !q.Enqueue(1) {
			t.Error("第一次入队失败")
		}
		if !q.Enqueue(2) {
			t.Error("第二次入队失败")
		}
		if !q.IsFull() {
			t.Error("队列应为满")
		}

		// 尝试再次入队（应该失败）
		if q.Enqueue(3) {
			t.Error("队列已满时应入队失败")
		}
		if q.Len() != 2 {
			t.Errorf("队列长度应为2，实际为 %d", q.Len())
		}
	})

	// 测试4: 强制入队（淘汰队首）
	t.Run("强制入队测试", func(t *testing.T) {
		q := NewSafeFixedQueue[int](2)

		// 填满队列
		q.Enqueue(1)
		q.Enqueue(2)

		// 强制入队，应该淘汰1
		q.EnqueueForce(3)

		if q.Len() != 2 {
			t.Errorf("队列长度应为2，实际为 %d", q.Len())
		}
		if !q.IsFull() {
			t.Error("队列应为满")
		}

		// 验证淘汰机制
		item, ok := q.Dequeue()
		if !ok {
			t.Error("出队失败")
		}
		if item != 2 {
			t.Errorf("第一个出队元素应为2（1被淘汰），实际为 %d", item)
		}

		item, ok = q.Dequeue()
		if !ok {
			t.Error("出队失败")
		}
		if item != 3 {
			t.Errorf("第二个出队元素应为3，实际为 %d", item)
		}
	})

	// 测试5: 空队列出队
	t.Run("空队列出队", func(t *testing.T) {
		q := NewSafeFixedQueue[string](3)

		item, ok := q.Dequeue()
		if ok {
			t.Error("空队列出队应返回false")
		}
		if item != "" {
			t.Errorf("空队列出队应返回零值，实际为 %s", item)
		}
	})

	// 测试6: 循环队列功能测试
	t.Run("循环队列功能", func(t *testing.T) {
		q := NewSafeFixedQueue[int](3)

		// 入队3个元素
		q.Enqueue(1)
		q.Enqueue(2)
		q.Enqueue(3)

		// 出队2个元素
		q.Dequeue() // 出队1
		q.Dequeue() // 出队2

		// 再入队2个元素（应该使用循环空间）
		if !q.Enqueue(4) {
			t.Error("循环入队失败")
		}
		if !q.Enqueue(5) {
			t.Error("循环入队失败")
		}

		if q.Len() != 3 {
			t.Errorf("队列长度应为3，实际为 %d", q.Len())
		}

		// 验证出队顺序
		item, _ := q.Dequeue()
		if item != 3 {
			t.Errorf("第一个出队元素应为3，实际为 %d", item)
		}
		item, _ = q.Dequeue()
		if item != 4 {
			t.Errorf("第二个出队元素应为4，实际为 %d", item)
		}
		item, _ = q.Dequeue()
		if item != 5 {
			t.Errorf("第三个出队元素应为5，实际为 %d", item)
		}
	})

	// 测试7: 边界条件测试
	t.Run("边界条件", func(t *testing.T) {
		// 测试容量为0的队列
		q := NewSafeFixedQueue[int](0)
		if q.maxSize != 0 {
			t.Error("容量为0的队列创建失败")
		}
		if !q.IsFull() {
			t.Error("容量为0的队列应为满")
		}
		if !q.IsEmpty() {
			t.Error("容量为0的队列应为空")
		}
		if q.Enqueue(1) {
			t.Error("容量为0的队列入队应失败")
		}

		// 测试容量为1的队列
		q1 := NewSafeFixedQueue[int](1)
		if !q1.Enqueue(10) {
			t.Error("容量为1的队列第一次入队应成功")
		}
		if q1.Enqueue(20) {
			t.Error("容量为1的队列第二次入队应失败")
		}
	})

	// 测试8: 并发安全测试（简单版本）
	t.Run("并发安全", func(t *testing.T) {
		q := NewSafeFixedQueue[int](100)

		// 简单并发测试
		done := make(chan bool)

		// 启动多个goroutine同时操作队列
		for i := 0; i < 10; i++ {
			go func(val int) {
				for j := 0; j < 10; j++ {
					q.Enqueue(val*10 + j)
				}
				done <- true
			}(i)
		}

		// 等待所有goroutine完成
		for i := 0; i < 10; i++ {
			<-done
		}

		// 验证队列状态
		if q.Len() != 100 {
			t.Errorf("并发操作后队列长度应为100，实际为 %d", q.Len())
		}
		if !q.IsFull() {
			t.Error("并发操作后队列应为满")
		}
	})
}
