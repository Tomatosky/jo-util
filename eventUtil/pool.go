package eventUtil

import (
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
