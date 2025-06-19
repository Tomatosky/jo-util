package eventUtil

import (
	"errors"
	"github.com/panjf2000/ants/v2"
	"sync"
	"time"
)

type Pool struct {
	pool *ants.Pool
	wg   sync.WaitGroup
}

func NewPool(size int) *Pool {
	pool, _ := ants.NewPool(size)
	return &Pool{pool: pool}
}

func (p *Pool) Submit(task func()) error {
	if task == nil {
		return errors.New("task cannot be nil")
	}
	p.wg.Add(1)
	err := p.pool.Submit(func() {
		defer p.wg.Done()
		task()
	})
	if err != nil {
		p.wg.Done()
		return err
	}
	return nil
}

// ScheduleAtFixedRate 类似于Java的scheduleAtFixedRate
// 以固定的频率执行任务，不考虑任务执行时间
// 返回一个函数，调用它可以停止调度
func (p *Pool) ScheduleAtFixedRate(initialDelay, period time.Duration, task func()) (stop func()) {
	done := make(chan struct{})
	// 初始延迟
	time.AfterFunc(initialDelay, func() {
		_ = p.Submit(task)

		ticker := time.NewTicker(period)
		go func() {
			for {
				select {
				case <-ticker.C:
					_ = p.Submit(task)
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
func (p *Pool) ScheduleWithFixedDelay(initialDelay, delay time.Duration, task func()) (stop func()) {
	done := make(chan struct{})
	var once sync.Once
	// 初始延迟
	time.AfterFunc(initialDelay, func() {
		_ = p.Submit(func() {
			task()
			p.scheduleNextWithDelay(delay, task, done, &once)
		})
	})
	return func() { once.Do(func() { close(done) }) }
}

// 递归调用来实现固定延迟调度
func (p *Pool) scheduleNextWithDelay(delay time.Duration, task func(), done <-chan struct{}, once *sync.Once) {
	select {
	case <-done:
		return
	case <-time.After(delay):
		_ = p.Submit(func() {
			task()
			p.scheduleNextWithDelay(delay, task, done, once)
		})
	}
}

func (p *Pool) Release(timeout time.Duration) (isTimeout bool) {
	defer p.pool.Release()

	// 创建一个通道用于通知等待完成
	done := make(chan struct{})

	// 启动一个goroutine来等待任务完成
	go func() {
		p.wg.Wait()
		close(done)
	}()

	// 使用select实现超时控制
	select {
	case <-done:
		// 所有任务正常完成
		return false
	case <-time.After(timeout):
		// 超时后直接返回，不等待剩余任务
		return true
	}
}
