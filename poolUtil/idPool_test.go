package poolUtil

import (
	"github.com/Tomatosky/jo-util/randomUtil"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go.uber.org/zap/zaptest"
)

func TestIdPool(t *testing.T) {
	// 测试用例1: 基本功能测试
	t.Run("BasicFunctionality", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		pool := NewIdPool[int32](&IdPoolOpt{
			PoolSize:  4,
			QueueSize: 10,
			Logger:    logger,
			PoolName:  "TestPool",
		})

		// 测试任务提交和计数
		var counter int32
		pool.SubmitWithId(1, func() {
			atomic.AddInt32(&counter, 1)
		})
		pool.SubmitWithId(1, func() {
			atomic.AddInt32(&counter, 1)
		})

		// 等待任务执行完成
		time.Sleep(100 * time.Millisecond)

		if got := pool.GetTaskCount(1); got != 0 {
			t.Errorf("Expected task count 0, got %d", got)
		}
		if atomic.LoadInt32(&counter) != 2 {
			t.Errorf("Expected counter 2, got %d", counter)
		}

		// 测试关闭
		pool.Shutdown(time.Second)
	})

	// 测试用例2: 并发任务提交
	t.Run("ConcurrentSubmission", func(t *testing.T) {
		pool := NewIdPool[int32](&IdPoolOpt{
			PoolSize:  4,
			QueueSize: 2048,
		})

		var wg sync.WaitGroup
		var counter int32
		const numTasks = 1000

		// 并发提交任务
		for i := 0; i < numTasks; i++ {
			wg.Add(1)
			go func(id int32) {
				defer wg.Done()
				pool.SubmitWithId(id%4, func() {
					atomic.AddInt32(&counter, 1)
				})
			}(int32(i))
		}

		wg.Wait()
		time.Sleep(500 * time.Millisecond) // 等待所有任务执行完成

		if atomic.LoadInt32(&counter) != numTasks {
			t.Errorf("Expected counter %d, got %d", numTasks, counter)
		}

		// 检查所有ID的任务计数是否清零
		for i := 0; i < 4; i++ {
			if count := pool.GetTaskCount(int32(i)); count != 0 {
				t.Errorf("Expected task count 0 for id %d, got %d", i, count)
			}
		}

		pool.Shutdown(time.Second)
	})

	// 测试用例3: 大并发压力测试
	t.Run("HighConcurrencyStressTest", func(t *testing.T) {
		pool := NewIdPool[int32](&IdPoolOpt{
			PoolSize:  16,
			QueueSize: 1000,
		})

		var wg sync.WaitGroup
		var counter int32
		const numTasks = 10000
		const numWorkers = 100

		// 使用多个worker并发提交任务
		for w := 0; w < numWorkers; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := 0; i < numTasks/numWorkers; i++ {
					pool.SubmitWithId(int32(i%16), func() {
						atomic.AddInt32(&counter, 1)
					})
				}
			}()
		}

		wg.Wait()
		time.Sleep(2 * time.Second) // 给足够时间执行所有任务

		if atomic.LoadInt32(&counter) != numTasks {
			t.Errorf("Expected counter %d, got %d", numTasks, counter)
		}

		pool.Shutdown(5 * time.Second)
	})

	// 测试用例4: 队列满的情况
	t.Run("QueueFull", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		pool := NewIdPool[int32](&IdPoolOpt{
			PoolSize:  2,
			QueueSize: 2,
			Logger:    logger,
			PoolName:  "SmallQueuePool",
		})

		// 填满队列
		for i := 0; i < 2; i++ {
			pool.SubmitWithId(0, func() {
				time.Sleep(100 * time.Millisecond)
			})
		}

		// 尝试提交更多任务（应该会触发队列满警告）
		pool.SubmitWithId(0, func() {})
		pool.SubmitWithId(0, func() {})

		pool.Shutdown(time.Second)
	})

	// 测试用例5: 关闭时处理剩余任务
	t.Run("ShutdownWithPendingTasks", func(t *testing.T) {
		pool := NewIdPool[int32](&IdPoolOpt{
			PoolSize:  2,
			QueueSize: 4096,
		})

		var counter int32
		const numTasks = 5

		// 提交一些长时间运行的任务
		for i := 0; i < numTasks; i++ {
			pool.SubmitWithId(0, func() {
				time.Sleep(200 * time.Millisecond)
				atomic.AddInt32(&counter, 1)
			})
		}

		// 立即关闭
		pool.Shutdown(3 * time.Second)

		if atomic.LoadInt32(&counter) != numTasks {
			t.Errorf("Expected all %d tasks to complete, got %d", numTasks, counter)
		}
	})

	// 测试用例6: 关闭超时
	t.Run("ShutdownTimeout", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		pool := NewIdPool[int32](&IdPoolOpt{
			PoolSize:  2,
			QueueSize: 10,
			Logger:    logger,
			PoolName:  "TimeoutPool",
		})

		// 提交一个长时间运行的任务
		pool.SubmitWithId(0, func() {
			time.Sleep(2 * time.Second)
		})

		// 设置很短的超时时间
		pool.Shutdown(100 * time.Millisecond)

		// 预期会触发超时警告
	})

	// 测试用例7: 任务ID映射清理
	t.Run("TaskIdMapCleanup", func(t *testing.T) {
		pool := NewIdPool[int32](&IdPoolOpt{
			PoolSize:  8,
			QueueSize: 1024,
		})

		const numTasks = 100
		for i := 0; i < numTasks; i++ {
			pool.SubmitWithId(int32(randomUtil.RandomInt(0, 32)), func() {})
		}

		time.Sleep(100 * time.Millisecond)

		// 检查任务ID映射是否被清理
		if size := pool.taskIdMap.Size(); size != 0 {
			t.Errorf("Expected taskIdMap size 0, got %d", size)
		}

		pool.Shutdown(time.Second)
	})

	// 测试用例8: 任务计数准确性
	t.Run("TaskCountAccuracy", func(t *testing.T) {
		pool := NewIdPool[int32](&IdPoolOpt{
			PoolSize:  4,
			QueueSize: 100,
		})

		const numTasksPerID = 50
		const numIDs = 4

		// 提交任务
		for id := 0; id < numIDs; id++ {
			for i := 0; i < numTasksPerID; i++ {
				pool.SubmitWithId(int32(id), func() {
					time.Sleep(10 * time.Millisecond)
				})
			}
		}

		// 立即检查任务计数（可能还在队列中）
		for id := 0; id < numIDs; id++ {
			count := pool.GetTaskCount(int32(id))
			if count <= 0 || count > numTasksPerID {
				t.Errorf("Unexpected task count %d for id %d", count, id)
			}
		}

		// 等待所有任务完成
		time.Sleep(600 * time.Millisecond)

		// 再次检查任务计数（应该为0）
		for id := 0; id < numIDs; id++ {
			if count := pool.GetTaskCount(int32(id)); count != 0 {
				t.Errorf("Expected task count 0 for id %d, got %d", id, count)
			}
		}

		pool.Shutdown(time.Second)
	})
}
