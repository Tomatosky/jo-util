package eventUtil

import (
	"fmt"
	"github.com/Tomatosky/jo-util/logger"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	logger.InitLog(nil)

	pool := NewPool(100)
	logger.Log.Info("")
	for i := 0; i < 3; i++ {
		pool.Submit(func() {
			time.Sleep(2 * time.Second)
			logger.Log.Info(fmt.Sprintf("%d", i))
		})
	}

	pool.Release(10 * time.Second)
	logger.Log.Info("")
}
