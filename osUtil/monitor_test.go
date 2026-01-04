package osUtil

import (
	"os"
	"testing"
	"time"

	"github.com/Tomatosky/jo-util/logger"
)

func TestMonitor_SetAlert(t *testing.T) {
	monitor := NewMonitor()
	monitor.SetMemory(20, 10*time.Second)
	monitor.SetAlertInterval(5 * time.Second)
	monitor.SetAlert(&defaultAlert{})
	_ = monitor.Start()
	logger.Log.Info("monitor started")

	interrupt := make(chan os.Signal, 1)
	<-interrupt
}
