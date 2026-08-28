package osUtil

import (
	"regexp"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestNewMonitorDefaults(t *testing.T) {
	monitor := NewMonitor("service")
	if monitor == nil || monitor.name != "service" || monitor.stopChan == nil {
		t.Fatalf("NewMonitor() = %#v", monitor)
	}
	if monitor.running {
		t.Error("new monitor is running")
	}
	for name, config := range map[string]thresholdConfig{
		"cpu":    monitor.cpu,
		"memory": monitor.memory,
		"disk":   monitor.disk,
	} {
		if config.enabled || config.threshold != 0 || config.duration != 0 || config.alertInterval != time.Minute || !config.startTime.IsZero() || !config.lastAlertTime.IsZero() {
			t.Errorf("%s default config = %+v", name, config)
		}
	}
}

func TestMonitorConfiguration(t *testing.T) {
	monitor := NewMonitor("service")
	monitor.cpu.startTime = time.Now()
	monitor.cpu.lastAlertTime = time.Now()
	monitor.SetCPU(80.5, 30*time.Second)
	if !monitor.cpu.enabled || monitor.cpu.threshold != 80.5 || monitor.cpu.duration != 30*time.Second || !monitor.cpu.startTime.IsZero() || !monitor.cpu.lastAlertTime.IsZero() {
		t.Errorf("SetCPU() config = %+v", monitor.cpu)
	}

	monitor.SetMemory(70, time.Minute)
	if !monitor.memory.enabled || monitor.memory.threshold != 70 || monitor.memory.duration != time.Minute {
		t.Errorf("SetMemory() config = %+v", monitor.memory)
	}
	monitor.SetDisk(60, 2*time.Minute)
	if !monitor.disk.enabled || monitor.disk.threshold != 60 || monitor.disk.duration != 2*time.Minute {
		t.Errorf("SetDisk() config = %+v", monitor.disk)
	}

	monitor.SetAll(50, 5*time.Second)
	for name, config := range map[string]thresholdConfig{
		"cpu":    monitor.cpu,
		"memory": monitor.memory,
		"disk":   monitor.disk,
	} {
		if !config.enabled || config.threshold != 50 || config.duration != 5*time.Second {
			t.Errorf("%s config after SetAll() = %+v", name, config)
		}
	}

	monitor.SetAlertInterval(15 * time.Minute)
	if monitor.cpu.alertInterval != 15*time.Minute || monitor.memory.alertInterval != 15*time.Minute || monitor.disk.alertInterval != 15*time.Minute {
		t.Error("SetAlertInterval() did not update all resource configs")
	}
}

func TestMonitorLifecycle(t *testing.T) {
	monitor := NewMonitor("service")
	monitor.Stop()
	if monitor.running {
		t.Error("Stop() changed a stopped monitor to running")
	}

	if err := monitor.Start(); err != nil {
		t.Fatalf("Start(): %v", err)
	}
	if !monitor.running {
		t.Error("Start() did not set running")
	}
	if err := monitor.Start(); err == nil || err.Error() != "monitor is already running" {
		t.Errorf("second Start() error = %v", err)
	}

	oldStopChan := monitor.stopChan
	monitor.Stop()
	if monitor.running {
		t.Error("Stop() did not clear running")
	}
	select {
	case <-oldStopChan:
	default:
		t.Error("Stop() did not close the active stop channel")
	}
	select {
	case <-monitor.stopChan:
		t.Error("Stop() left the replacement stop channel closed")
	default:
	}

	if err := monitor.Start(); err != nil {
		t.Fatalf("restart: %v", err)
	}
	monitor.Stop()
	monitor.Stop()
}

func TestMonitorConcurrentLifecycle(t *testing.T) {
	monitor := NewMonitor("service")
	const callers = 20
	for cycle := 0; cycle < 20; cycle++ {
		start := make(chan struct{})
		results := make(chan error, callers)
		var wg sync.WaitGroup
		wg.Add(callers)
		for i := 0; i < callers; i++ {
			go func() {
				defer wg.Done()
				<-start
				results <- monitor.Start()
			}()
		}
		close(start)
		wg.Wait()
		close(results)

		successes := 0
		for err := range results {
			if err == nil {
				successes++
				continue
			}
			if err.Error() != "monitor is already running" {
				t.Errorf("Start() error = %v", err)
			}
		}
		if successes != 1 {
			t.Fatalf("concurrent Start successes = %d, want 1", successes)
		}

		wg.Add(callers)
		for i := 0; i < callers; i++ {
			go func() {
				defer wg.Done()
				monitor.Stop()
			}()
		}
		wg.Wait()
		if monitor.running {
			t.Fatalf("monitor still running after concurrent Stop in cycle %d", cycle)
		}
	}

	if err := monitor.Start(); err != nil {
		t.Fatalf("monitor unusable after concurrent lifecycle cycles: %v", err)
	}
	monitor.Stop()
}

func TestCheckThresholdStateTransitions(t *testing.T) {
	monitor := NewMonitor("service")
	config := thresholdConfig{
		enabled:       true,
		threshold:     80,
		duration:      time.Minute,
		alertInterval: time.Hour,
		startTime:     time.Now().Add(-time.Minute),
	}
	monitor.checkThreshold(&config, CPU, 79.99)
	if !config.startTime.IsZero() {
		t.Errorf("below-threshold value did not reset start time: %v", config.startTime)
	}

	before := time.Now()
	monitor.checkThreshold(&config, CPU, 80)
	after := time.Now()
	if config.startTime.Before(before) || config.startTime.After(after) || !config.lastAlertTime.IsZero() {
		t.Errorf("first threshold crossing state = %+v", config)
	}

	started := time.Now().Add(-2 * time.Minute)
	config.startTime = started
	before = time.Now()
	monitor.checkThreshold(&config, CPU, 95)
	after = time.Now()
	if config.lastAlertTime.Before(before) || config.lastAlertTime.After(after) {
		t.Errorf("alert time = %v, want within [%v, %v]", config.lastAlertTime, before, after)
	}
	firstAlert := config.lastAlertTime
	monitor.checkThreshold(&config, CPU, 96)
	if !config.lastAlertTime.Equal(firstAlert) {
		t.Errorf("alert interval was ignored: first=%v second=%v", firstAlert, config.lastAlertTime)
	}

	monitor.checkThreshold(&config, CPU, 0)
	if !config.startTime.IsZero() || !config.lastAlertTime.Equal(firstAlert) {
		t.Errorf("recovery state = %+v, want zero start and preserved last alert", config)
	}
}

func TestMemoryStringFormats(t *testing.T) {
	format := regexp.MustCompile(`^\d+\.\d{2}$`)
	for name, value := range map[string]string{
		"KB": MemUseKBStr(),
		"MB": MemUseMBStr(),
		"GB": MemUseGBStr(),
	} {
		if !format.MatchString(value) {
			t.Errorf("MemUse%sStr() = %q, want non-negative number with two decimals", name, value)
		}
		parsed, err := strconv.ParseFloat(value, 64)
		if err != nil || parsed < 0 {
			t.Errorf("MemUse%sStr() = %q, parse error %v", name, value, err)
		}
	}
}
