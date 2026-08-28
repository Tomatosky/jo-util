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
	pool    *ants.Pool
	wg      sync.WaitGroup
	stateMu sync.Mutex
	closed  bool
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
	done := make(chan struct{})
	var once sync.Once
	// 初始延迟
	timer := time.AfterFunc(initialDelay, func() {
		select {
		case <-done:
			return
		default:
			p.Submit(func() {
				select {
				case <-done:
					return
				default:
					task()
				}
			})
		}

		ticker := time.NewTicker(period)
		go func() {
			for {
				select {
				case <-ticker.C:
					p.Submit(func() {
						select {
						case <-done:
							return
						default:
							task()
						}
					})
				case <-done:
					ticker.Stop()
					return
				}
			}
		}()
	})
	return func() {
		once.Do(func() {
			close(done)
			timer.Stop()
		})
	}
}

// ScheduleWithFixedDelay 类似于Java的scheduleWithFixedDelay
// 在上一次任务完成后，固定延迟时间后执行下一次任务
// 返回一个函数，调用它可以停止调度
func (p *AntsPool) ScheduleWithFixedDelay(initialDelay, delay time.Duration, task func()) (stop func()) {
	done := make(chan struct{})
	var once sync.Once
	// 初始延迟
	timer := time.AfterFunc(initialDelay, func() {
		select {
		case <-done:
			return
		default:
		}
		p.Submit(func() {
			select {
			case <-done:
				return
			default:
			}
			task()
			p.scheduleNextWithDelay(delay, task, done, &once)
		})
	})
	return func() {
		once.Do(func() {
			close(done)
			timer.Stop()
		})
	}
}

// 递归调用来实现固定延迟调度
func (p *AntsPool) scheduleNextWithDelay(delay time.Duration, task func(), done <-chan struct{}, once *sync.Once) {
	select {
	case <-done:
		return
	case <-time.After(delay):
		select {
		case <-done:
			return
		default:
		}
		p.Submit(func() {
			select {
			case <-done:
				return
			default:
			}
			task()
			p.scheduleNextWithDelay(delay, task, done, once)
		})
	}
}

func (p *AntsPool) Shutdown(timeout time.Duration) (isTimeout bool) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	p.stateMu.Lock()
	p.closed = true
	p.stateMu.Unlock()

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
