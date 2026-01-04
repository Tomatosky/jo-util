package osUtil

import (
	"fmt"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

// Alert 报警接口，由用户自行实现
type Alert interface {
	Alert(resourceType string, value float64, threshold float64, duration time.Duration)
}

// ResourceType 资源类型
type ResourceType string

const (
	CPU    ResourceType = "CPU"
	Memory ResourceType = "Memory"
	Disk   ResourceType = "Disk"
)

// thresholdConfig 阈值配置
type thresholdConfig struct {
	enabled   bool
	threshold float64 // 阈值百分比
	duration  time.Duration
	startTime time.Time
}

// Monitor 资源监控器
type Monitor struct {
	cpu    thresholdConfig
	memory thresholdConfig
	disk   thresholdConfig

	alert Alert

	stopChan chan struct{}
	wg      sync.WaitGroup
	mu      sync.Mutex
	running bool
}

// defaultAlert 默认报警实现
type defaultAlert struct{}

func (d *defaultAlert) Alert(resourceType string, value float64, threshold float64, duration time.Duration) {
	fmt.Printf("[资源报警] %s 当前值: %.2f%% 阈值: %.2f%% 持续时间: %v\n",
		resourceType, value, threshold, duration)
}

// NewMonitor 创建一个新的监控器，使用默认的 fmt.Printf 报警
func NewMonitor() *Monitor {
	return &Monitor{
		cpu: thresholdConfig{
			enabled:   false,
			threshold: 0,
			duration:  0,
		},
		memory: thresholdConfig{
			enabled:   false,
			threshold: 0,
			duration:  0,
		},
		disk: thresholdConfig{
			enabled:   false,
			threshold: 0,
			duration:  0,
		},
		alert:    &defaultAlert{},
		stopChan: make(chan struct{}),
		running:  false,
	}
}

// SetAlert 设置自定义报警
func (m *Monitor) SetAlert(alert Alert) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.alert = alert
}

// SetCPU 设置 CPU 监控阈值
// threshold: 阈值百分比 (0-100)
// duration: 超过阈值后持续的时间，达到此时间才报警
func (m *Monitor) SetCPU(threshold float64, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cpu.threshold = threshold
	m.cpu.duration = duration
	m.cpu.enabled = true
	m.cpu.startTime = time.Time{}
}

// SetMemory 设置内存监控阈值
// threshold: 阈值百分比 (0-100)
// duration: 超过阈值后持续的时间，达到此时间才报警
func (m *Monitor) SetMemory(threshold float64, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.memory.threshold = threshold
	m.memory.duration = duration
	m.memory.enabled = true
	m.memory.startTime = time.Time{}
}

// SetDisk 设置硬盘监控阈值
// threshold: 阈值百分比 (0-100)
// duration: 超过阈值后持续的时间，达到此时间才报警
func (m *Monitor) SetDisk(threshold float64, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.disk.threshold = threshold
	m.disk.duration = duration
	m.disk.enabled = true
	m.disk.startTime = time.Time{}
}

// Start 启动监控
func (m *Monitor) Start() error {
	m.mu.Lock()
	if m.running {
		m.mu.Unlock()
		return fmt.Errorf("monitor is already running")
	}
	m.running = true
	m.mu.Unlock()

	m.wg.Add(1)
	go m.monitorLoop()
	return nil
}

// Stop 停止监控
func (m *Monitor) Stop() {
	m.mu.Lock()
	if !m.running {
		m.mu.Unlock()
		return
	}
	m.running = false
	m.mu.Unlock()

	close(m.stopChan)
	m.wg.Wait()
	m.stopChan = make(chan struct{})
}

// monitorLoop 监控循环
func (m *Monitor) monitorLoop() {
	defer m.wg.Done()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopChan:
			return
		case <-ticker.C:
			m.checkResources()
		}
	}
}

// checkResources 检查资源占用情况
func (m *Monitor) checkResources() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查 CPU
	if m.cpu.enabled {
		cpuUsage := m.getCPUUsage()
		m.checkThreshold(&m.cpu, CPU, cpuUsage)
	}

	// 检查内存
	if m.memory.enabled {
		memoryUsage := m.getMemoryUsage()
		m.checkThreshold(&m.memory, Memory, memoryUsage)
	}

	// 检查硬盘
	if m.disk.enabled {
		diskUsage := m.getDiskUsage()
		m.checkThreshold(&m.disk, Disk, diskUsage)
	}
}

// checkThreshold 检查是否超过阈值
func (m *Monitor) checkThreshold(config *thresholdConfig, resourceType ResourceType, value float64) {
	now := time.Now()

	if value >= config.threshold {
		// 超过阈值
		if config.startTime.IsZero() {
			// 第一次超过，记录开始时间
			config.startTime = now
		} else if now.Sub(config.startTime) >= config.duration {
			// 持续超过阈值时间，触发报警
			m.alert.Alert(string(resourceType), value, config.threshold, config.duration)

			// 重置开始时间，避免重复报警
			// 如果希望持续报警，可以注释掉下面这行
			config.startTime = now
		}
	} else {
		// 未超过阈值，重置开始时间
		config.startTime = time.Time{}
	}
}

// getCPUUsage 获取 CPU 使用率（百分比）
func (m *Monitor) getCPUUsage() float64 {
	percent, err := cpu.Percent(0, false)
	if err != nil || len(percent) == 0 {
		return 0
	}
	return percent[0]
}

// getMemoryUsage 获取内存使用率（百分比）
func (m *Monitor) getMemoryUsage() float64 {
	memStat, err := mem.VirtualMemory()
	if err != nil {
		return 0
	}
	return memStat.UsedPercent
}

// getDiskUsage 获取磁盘使用率（百分比）
func (m *Monitor) getDiskUsage() float64 {
	// 获取当前目录所在磁盘的使用率
	diskStat, err := disk.Usage(".")
	if err != nil {
		return 0
	}
	return diskStat.UsedPercent
}
