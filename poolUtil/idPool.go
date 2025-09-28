package poolUtil

import (
	"context"
	"fmt"
	"github.com/Tomatosky/jo-util/convertor"
	"github.com/Tomatosky/jo-util/idUtil"
	"github.com/Tomatosky/jo-util/mapUtil"
	"github.com/Tomatosky/jo-util/randomUtil"
	"go.uber.org/zap"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

var _ IPool = (*IdPool)(nil)

type IdPool struct {
	workers      []*worker
	taskIdMap    *mapUtil.ConcurrentHashMap[string, int64]        // key: taskID(string), value: id
	idTaskCounts *mapUtil.ConcurrentHashMap[int64, *atomic.Int32] // key: id, value: *atomic.Int32
	cores        int64
	running      atomic.Bool    // 控制服务运行状态
	wg           sync.WaitGroup // 用于等待所有worker退出
	logger       *zap.Logger
	poolName     string
}

type worker struct {
	idPool *IdPool          // 反向引用 IdPool
	queue  chan *customTask // 任务通道
	done   chan struct{}    // 关闭信号
}

type customTask struct {
	taskID string
	task   func()
}

type IdPoolOpt struct {
	PoolSize  int64
	QueueSize int
	Logger    *zap.Logger
	PoolName  string
}

func NewIdPool(opt *IdPoolOpt) *IdPool {
	if opt.PoolSize <= 0 || opt.QueueSize <= 0 {
		debug.PrintStack()
		panic("pool size and queue size must be greater than 0")
	}

	idPool := &IdPool{
		cores:        opt.PoolSize,
		workers:      make([]*worker, opt.PoolSize),
		taskIdMap:    mapUtil.NewConcurrentHashMap[string, int64](),
		idTaskCounts: mapUtil.NewConcurrentHashMap[int64, *atomic.Int32](),
	}
	if opt.Logger != nil {
		idPool.logger = opt.Logger
	}
	idPool.running.Store(true)
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
	if !i.running.Load() {
		return
	}
	idInt64 := convertor.ToInt64(id)
	// 生成唯一任务ID
	taskID := idUtil.RandomUUID()
	// 更新任务计数
	v, _ := i.idTaskCounts.PutIfAbsent(idInt64, &atomic.Int32{})
	v.Add(1)
	// 记录任务映射
	i.taskIdMap.Put(taskID, idInt64)
	// 选择 worker（哈希取模）
	w := i.workers[idInt64%i.cores]
	// 发送任务
	select {
	case w.queue <- &customTask{taskID: taskID, task: task}:
	default:
		if i.logger != nil {
			i.logger.Warn(fmt.Sprintf("%s queue is full", i.poolName))
		} else {
			fmt.Println(fmt.Sprintf("%s queue is full", i.poolName))
		}
	}
}

// GetTaskCount 获取任务计数
func (i *IdPool) GetTaskCount(id any) int32 {
	idInt64 := convertor.ToInt64(id)
	v := i.idTaskCounts.GetOrDefault(idInt64, &atomic.Int32{})
	return v.Load()
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

	// 停止接收新任务
	i.running.Store(false)

	// 通知所有 worker 停止
	for _, w := range i.workers {
		close(w.done) // 发送关闭信号
	}

	// 创建一个 channel 用于等待 WaitGroup
	waitCh := make(chan struct{})
	go func() {
		i.wg.Wait()
		close(waitCh)
	}()

	// 等待所有 worker 退出或上下文取消
	select {
	case <-waitCh:
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
			if w.idPool.logger != nil {
				w.idPool.logger.Error(fmt.Sprintf("err=%v", err))
			} else {
				fmt.Println(fmt.Sprintf("err=%v", err))
				debug.PrintStack()
			}
		}

		// 清理任务映射并减少计数
		id := w.idPool.taskIdMap.Get(task.taskID)
		w.idPool.taskIdMap.Remove(task.taskID)
		v := w.idPool.idTaskCounts.Get(id)
		if v != nil {
			v.Add(-1)
		}
	}()

	// 执行任务
	task.task()
}
