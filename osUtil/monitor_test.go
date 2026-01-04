package osUtil

import (
	"os"
	"testing"
	"time"

	"github.com/Tomatosky/jo-util/logger"
)

func TestMonitor_SetDefaultAlert(t *testing.T) {
	monitor := NewMonitor("test")
	monitor.SetCPU(5, 10*time.Second)
	//monitor.SetMemory(20, 10*time.Second)
	//monitor.SetDisk(10, 10*time.Second)
	monitor.SetAlertInterval(60 * time.Second)
	monitor.SetAlert(&defaultAlert{})
	_ = monitor.Start()
	logger.Log.Info("monitor started")

	interrupt := make(chan os.Signal, 1)
	<-interrupt
}

func TestMonitor_SetGotifyAlert(t *testing.T) {
	monitor := NewMonitor("test")
	monitor.SetMemory(5, 10*time.Second)
	monitor.SetAlertInterval(30 * time.Second)
	monitor.SetAlert(&GotifyAlert{
		host:  "https://abc.gotify.com",
		token: "123456",
	})
	_ = monitor.Start()
	logger.Log.Info("monitor started")

	interrupt := make(chan os.Signal, 1)
	<-interrupt
}

func TestMonitor_SetDingdingAlert(t *testing.T) {
	monitor := NewMonitor("test")
	monitor.SetMemory(5, 10*time.Second)
	monitor.SetAlertInterval(30 * time.Second)
	monitor.SetAlert(&DingdingAlert{
		secret:      "abcdefg",
		accessToken: "abcdefg",
	})
	_ = monitor.Start()
	logger.Log.Info("monitor started")

	interrupt := make(chan os.Signal, 1)
	<-interrupt
}
