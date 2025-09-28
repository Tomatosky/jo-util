package poolUtil

import (
	"context"
	"github.com/panjf2000/ants/v2"
	"sync"
	"time"
)

var _ IPool = (*AntsPool)(nil)

type AntsPool struct {
	pool *ants.Pool
	wg   sync.WaitGroup
}

func NewAntsPool(size int) *AntsPool {
	pool, _ := ants.NewPool(size)
	return &AntsPool{pool: pool}
}

func (p *AntsPool) SubmitWithId(id any, task func()) {
	p.Submit(task)
}

func (p *AntsPool) Submit(task func()) {
	if task == nil {
		panic("task cannot be nil")
	}
	p.wg.Add(1)
	_ = p.pool.Submit(func() {
		defer p.wg.Done()
		task()
	})
}

// ScheduleAtFixedRate 类似于Java的scheduleAtFixedRate
// 以固定的频率执行任务，不考虑任务执行时间
// 返回一个函数，调用它可以停止调度
func (p *AntsPool) ScheduleAtFixedRate(initialDelay, period time.Duration, task func()) (stop func()) {
	done := make(chan struct{})
	// 初始延迟
	time.AfterFunc(initialDelay, func() {
		p.Submit(task)

		ticker := time.NewTicker(period)
		go func() {
			for {
				select {
				case <-ticker.C:
					p.Submit(task)
				case <-done:
					ticker.Stop()
					return
				}
			}
		}()
	})
	return func() { close(done) }
}

// ScheduleWithFixedDelay 类似于Java的scheduleWithFixedDelay
// 在上一次任务完成后，固定延迟时间后执行下一次任务
// 返回一个函数，调用它可以停止调度
func (p *AntsPool) ScheduleWithFixedDelay(initialDelay, delay time.Duration, task func()) (stop func()) {
	done := make(chan struct{})
	var once sync.Once
	// 初始延迟
	time.AfterFunc(initialDelay, func() {
		p.Submit(func() {
			task()
			p.scheduleNextWithDelay(delay, task, done, &once)
		})
	})
	return func() { once.Do(func() { close(done) }) }
}

// 递归调用来实现固定延迟调度
func (p *AntsPool) scheduleNextWithDelay(delay time.Duration, task func(), done <-chan struct{}, once *sync.Once) {
	select {
	case <-done:
		return
	case <-time.After(delay):
		p.Submit(func() {
			task()
			p.scheduleNextWithDelay(delay, task, done, once)
		})
	}
}

func (p *AntsPool) Shutdown(timeout time.Duration) (isTimeout bool) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	defer p.pool.Release()

	// 创建一个通道用于通知等待完成
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	// 使用select实现超时控制
	select {
	case <-done:
		return false
	case <-ctx.Done():
		return true
	}
}
