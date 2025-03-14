package dbUtil

import (
	"github.com/Tomatosky/jo-util/idUtil"
	"github.com/Tomatosky/jo-util/mapUtil"
	"math"
	"sync/atomic"
)

type DbService struct {
	workers      []*Worker
	taskIdMap    *mapUtil.ConcurrentHashMap[string, int]        // key: taskID(string), value: id
	idTaskCounts *mapUtil.ConcurrentHashMap[int, *atomic.Int32] // key: id, value: *atomic.Int32
	cores        int
	running      atomic.Bool // 控制服务运行状态
}

type Worker struct {
	db    *DbService      // 反向引用 DbService
	queue chan CustomTask // 任务通道
	done  chan struct{}   // 关闭信号
}

type CustomTask struct {
	taskID string
	task   func()
}

func NewDbService(cores int) *DbService {
	db := &DbService{
		cores:        cores,
		workers:      make([]*Worker, cores),
		taskIdMap:    mapUtil.NewConcurrentHashMap[string, int](),
		idTaskCounts: mapUtil.NewConcurrentHashMap[int, *atomic.Int32](),
	}
	db.running.Store(true)
	// 初始化 workers
	for i := 0; i < cores; i++ {
		db.workers[i] = NewWorker(db)
		go db.workers[i].run()
	}
	return db
}

func NewWorker(db *DbService) *Worker {
	return &Worker{
		db:    db,
		queue: make(chan CustomTask, math.MaxInt), // 带缓冲的任务队列
		done:  make(chan struct{}),
	}
}

// AddTask 添加任务（线程安全）
func (db *DbService) AddTask(id int, task func()) {
	if !db.running.Load() {
		return
	}
	// 生成唯一任务ID
	taskID := idUtil.RandomUUID()
	// 更新任务计数
	v, _ := db.idTaskCounts.PutIfAbsent(id, &atomic.Int32{})
	v.Add(1)
	// 记录任务映射
	db.taskIdMap.Put(taskID, id)
	// 选择 worker（哈希取模）
	worker := db.workers[id%db.cores]
	// 发送任务（非阻塞发送，若队列满则可能丢失任务）
	select {
	case worker.queue <- CustomTask{taskID: taskID, task: task}:
	}
}

// GetTaskCount 获取任务计数
func (db *DbService) GetTaskCount(id int) int32 {
	v := db.idTaskCounts.GetOrDefault(id, &atomic.Int32{})
	return v.Load()
}

// Worker 运行循环
func (w *Worker) run() {
	for {
		select {
		case task := <-w.queue:
			// 执行任务
			task.task()
			// 清理任务映射并减少计数
			id := w.db.taskIdMap.Get(task.taskID)
			w.db.taskIdMap.Remove(task.taskID)
			v := w.db.idTaskCounts.Get(id)
			if v != nil {
				v.Add(-1)
			}
		case <-w.done:
			// 处理剩余任务
			w.drainQueue()
			return
		}
	}
}

// 排空剩余任务
func (w *Worker) drainQueue() {
	for {
		select {
		case task := <-w.queue:
			task.task()
			id := w.db.taskIdMap.Get(task.taskID)
			w.db.taskIdMap.Remove(task.taskID)
			v := w.db.idTaskCounts.Get(id)
			if v != nil {
				v.Add(-1)
			}
		default:
			return
		}
	}
}

// Shutdown 关闭服务
func (db *DbService) Shutdown() {
	// 停止接收新任务
	db.running.Store(false)
	// 通知所有 worker 停止
	for _, worker := range db.workers {
		close(worker.done) // 发送关闭信号
	}
	// 等待所有 worker 退出（可通过 sync.WaitGroup 改进）
	// 此处简化处理，实际生产环境需要更严谨的退出机制
}
