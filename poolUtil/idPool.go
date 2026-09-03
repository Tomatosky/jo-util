package poolUtil

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Tomatosky/jo-util/convertor"
	"github.com/Tomatosky/jo-util/logger"
	"github.com/Tomatosky/jo-util/randomUtil"
)

var _ IPool = (*IdPool)(nil)

type IdPool struct {
	workers      []*worker
	idTaskCounts map[int64]int32
	cores        int64
	running      bool           // 控制服务运行状态
	wg           sync.WaitGroup // 用于等待所有worker退出
	stateMu      sync.Mutex     // 保护生命周期和任务计数
	terminated   chan struct{}
	poolName     string
}

type worker struct {
	idPool *IdPool          // 反向引用 IdPool
	queue  chan *customTask // 任务通道
	done   chan struct{}    // 关闭信号
}

type customTask struct {
	id   int64
	task func()
}

type IdPoolOpt struct {
	PoolSize  int64
	QueueSize int
	PoolName  string
}

func NewIdPool(opt *IdPoolOpt) *IdPool {
	if opt.PoolSize <= 0 || opt.QueueSize <= 0 {
		logger.Log.Error(fmt.Sprintf("%v", "pool size and queue size must be greater than 0"))
		panic("pool size and queue size must be greater than 0")
	}

	idPool := &IdPool{
		cores:        opt.PoolSize,
		workers:      make([]*worker, opt.PoolSize),
		idTaskCounts: make(map[int64]int32),
		terminated:   make(chan struct{}),
		running:      true,
		poolName:     opt.PoolName,
	}
	// 初始化 workers
	for i := int64(0); i < opt.PoolSize; i++ {
		idPool.workers[i] = newWorker(idPool, opt.QueueSize)
		idPool.wg.Add(1) // 为每个worker增加计数
		go func(w *worker) {
			defer idPool.wg.Done() // worker退出时减少计数
			w.run()
		}(idPool.workers[i])
	}
	return idPool
}

func (i *IdPool) Submit(task func()) {
	i.SubmitWithId(int32(randomUtil.RandomInt(0, 100000)), task)
}

// SubmitWithId 添加任务
func (i *IdPool) SubmitWithId(id any, task func()) {
	idInt64 := convertor.ToInt64(id)

	i.stateMu.Lock()
	defer i.stateMu.Unlock()
	if !i.running {
		return
	}

	// 更新任务计数
	i.idTaskCounts[idInt64]++
	// 选择 worker（哈希取模）
	workerIndex := idInt64 % i.cores
	if workerIndex < 0 {
		workerIndex += i.cores
	}
	w := i.workers[workerIndex]
	// 发送任务
	select {
	case w.queue <- &customTask{id: idInt64, task: task}:
	default:
		i.idTaskCounts[idInt64]--
		if i.idTaskCounts[idInt64] == 0 {
			delete(i.idTaskCounts, idInt64)
		}
		logger.Log.Warn(fmt.Sprintf("%s queue is full", i.poolName))
	}
}

// GetTaskCount 获取任务计数
func (i *IdPool) GetTaskCount(id any) int32 {
	idInt64 := convertor.ToInt64(id)
	i.stateMu.Lock()
	defer i.stateMu.Unlock()
	return i.idTaskCounts[idInt64]
}

// MaxQueue 最大worker队列长度
func (i *IdPool) MaxQueue() int {
	num := 0
	for _, v := range i.workers {
		if len(v.queue) > num {
			num = len(v.queue)
		}
	}
	return num
}

// Shutdown 关闭服务
func (i *IdPool) Shutdown(timeout time.Duration) (isTimeout bool) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 停止接收新任务，并且只关闭 worker 一次。
	i.stateMu.Lock()
	if i.running {
		i.running = false
		for _, w := range i.workers {
			close(w.done)
		}
		go func() {
			i.wg.Wait()
			close(i.terminated)
		}()
	}
	i.stateMu.Unlock()

	// 等待所有 worker 退出或上下文取消
	select {
	case <-i.terminated:
		return false
	case <-ctx.Done():
		return true
	}
}

func newWorker(i *IdPool, queueSize int) *worker {
	return &worker{
		idPool: i,
		queue:  make(chan *customTask, queueSize), // 带缓冲的任务队列
		done:   make(chan struct{}),
	}
}

// worker 运行循环
func (w *worker) run() {
	for {
		select {
		case task := <-w.queue:
			w.processTask(task)
		case <-w.done:
			// 处理剩余任务
			w.drainQueue()
			return
		}
	}
}

// 排空剩余任务
func (w *worker) drainQueue() {
	for {
		select {
		case task, ok := <-w.queue:
			if !ok { // 通道已关闭时安全退出
				return
			}
			w.processTask(task)
		default: // 队列为空时立即退出
			return
		}
	}
}

func (w *worker) processTask(task *customTask) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Log.Error(fmt.Sprintf("err=%v", err))
		}

		// 与新任务提交串行化，避免计数丢失。
		w.idPool.stateMu.Lock()
		defer w.idPool.stateMu.Unlock()
		w.idPool.idTaskCounts[task.id]--
		if w.idPool.idTaskCounts[task.id] == 0 {
			delete(w.idPool.idTaskCounts, task.id)
		}
	}()

	// 执行任务
	task.task()
}
