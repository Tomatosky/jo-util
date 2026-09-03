package poolUtil

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Tomatosky/jo-util/logger"
	"github.com/panjf2000/ants/v2"
)

var _ IPool = (*AntsPool)(nil)

type AntsPool struct {
	pool       *ants.Pool
	wg         sync.WaitGroup
	stateMu    sync.Mutex
	closed     bool
	shutdown   chan struct{}
	terminated chan struct{}
}

func NewAntsPool(size int) *AntsPool {
	pool, _ := ants.NewPool(size)
	return &AntsPool{
		pool:       pool,
		shutdown:   make(chan struct{}),
		terminated: make(chan struct{}),
	}
}

func (p *AntsPool) SubmitWithId(id any, task func()) {
	p.Submit(task)
}

func (p *AntsPool) Submit(task func()) {
	if task == nil {
		logger.Log.Error(fmt.Sprintf("%v", "task cannot be nil"))
		panic("task cannot be nil")
	}
	p.stateMu.Lock()
	if p.closed {
		p.stateMu.Unlock()
		return
	}
	p.wg.Add(1)
	p.stateMu.Unlock()

	err := p.pool.Submit(func() {
		defer p.wg.Done()
		task()
	})
	if err != nil {
		p.wg.Done()
		logger.Log.Warn(fmt.Sprintf("submit task failed: %v", err))
	}
}

// ScheduleAtFixedRate 类似于Java的scheduleAtFixedRate
// 以固定的频率执行任务，不考虑任务执行时间
// 返回一个函数，调用它可以停止调度
func (p *AntsPool) ScheduleAtFixedRate(initialDelay, period time.Duration, task func()) (stop func()) {
	if task == nil {
		panic("task cannot be nil")
	}
	if period <= 0 {
		panic("period must be greater than 0")
	}
	done := make(chan struct{})
	var once sync.Once
	started := p.startSchedule(func() {
		if !p.waitScheduleDelay(initialDelay, done) {
			return
		}
		p.submitScheduled(done, task)

		ticker := time.NewTicker(period)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				p.submitScheduled(done, task)
			case <-done:
				return
			case <-p.shutdown:
				return
			}
		}
	})
	if !started {
		return func() {}
	}
	return func() {
		once.Do(func() {
			close(done)
		})
	}
}

// ScheduleWithFixedDelay 类似于Java的scheduleWithFixedDelay
// 在上一次任务完成后，固定延迟时间后执行下一次任务
// 返回一个函数，调用它可以停止调度
func (p *AntsPool) ScheduleWithFixedDelay(initialDelay, delay time.Duration, task func()) (stop func()) {
	if task == nil {
		panic("task cannot be nil")
	}
	if delay <= 0 {
		panic("delay must be greater than 0")
	}
	done := make(chan struct{})
	var once sync.Once
	started := p.startSchedule(func() {
		if !p.waitScheduleDelay(initialDelay, done) {
			return
		}

		for {
			completed := make(chan bool, 1)
			p.Submit(func() {
				succeeded := false
				defer func() {
					if recovered := recover(); recovered != nil {
						completed <- false
						panic(recovered)
					}
					completed <- succeeded
				}()

				select {
				case <-done:
					return
				case <-p.shutdown:
					return
				default:
					task()
					succeeded = true
				}
			})

			select {
			case succeeded := <-completed:
				if !succeeded {
					return
				}
			case <-done:
				return
			case <-p.shutdown:
				return
			}

			if !p.waitScheduleDelay(delay, done) {
				return
			}
		}
	})
	if !started {
		return func() {}
	}
	return func() {
		once.Do(func() {
			close(done)
		})
	}
}

func (p *AntsPool) startSchedule(run func()) bool {
	p.stateMu.Lock()
	defer p.stateMu.Unlock()
	if p.closed {
		return false
	}
	go run()
	return true
}

func (p *AntsPool) waitScheduleDelay(delay time.Duration, done <-chan struct{}) bool {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-done:
		return false
	case <-p.shutdown:
		return false
	case <-timer.C:
		return true
	}
}

func (p *AntsPool) submitScheduled(done <-chan struct{}, task func()) {
	p.Submit(func() {
		select {
		case <-done:
			return
		case <-p.shutdown:
			return
		default:
			task()
		}
	})
}

func (p *AntsPool) Shutdown(timeout time.Duration) (isTimeout bool) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	p.stateMu.Lock()
	if !p.closed {
		p.closed = true
		close(p.shutdown)
		go func() {
			p.wg.Wait()
			close(p.terminated)
		}()
	}
	p.stateMu.Unlock()
	defer p.pool.Release()

	// 使用select实现超时控制
	select {
	case <-p.terminated:
		return false
	case <-ctx.Done():
		return true
	}
}
