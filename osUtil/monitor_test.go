package osUtil

import (
	"os"
	"testing"
	"time"

	"github.com/Tomatosky/jo-util/alertUtil"
	"github.com/Tomatosky/jo-util/logger"
)

func TestMonitor_SetDefaultAlert(t *testing.T) {
	monitor := NewMonitor("test")
	monitor.SetCPU(5, 10*time.Second)
	//monitor.SetMemory(20, 10*time.Second)
	//monitor.SetDisk(10, 10*time.Second)
	monitor.SetAlertInterval(60 * time.Second)
	_ = monitor.Start()
	logger.Log.Info("monitor started")

	interrupt := make(chan os.Signal, 1)
	<-interrupt
}

func TestMonitor_SetGotifyAlert(t *testing.T) {
	monitor := NewMonitor("test")
	monitor.SetMemory(5, 10*time.Second)
	monitor.SetAlertInterval(30 * time.Second)
	monitor.AddAlert(&alertUtil.GotifyAlert{
		Host:  "https://abc.gotify.com",
		Token: "123456",
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
	monitor.AddAlert(&alertUtil.DingdingAlert{
		Secret:      "abcdefg",
		AccessToken: "abcdefg",
	})
	_ = monitor.Start()
	logger.Log.Info("monitor started")

	interrupt := make(chan os.Signal, 1)
	<-interrupt
}
