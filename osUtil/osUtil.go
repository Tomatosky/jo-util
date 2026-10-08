package osUtil

import (
	"fmt"
	"os"
	"runtime/metrics"

	"github.com/Tomatosky/jo-util/logger"
)

func MemUse() uint64 {
	m := []metrics.Sample{{Name: "/memory/classes/heap/objects:bytes"}}
	metrics.Read(m)
	return m[0].Value.Uint64()
}

func MemUseKB() float32 {
	return float32(MemUse()) / 1024
}

func MemUseKBStr() string {
	return fmt.Sprintf("%.2f", MemUseKB())
}

func MemUseMB() float32 {
	return float32(MemUse()) / 1024 / 1024
}

func MemUseMBStr() string {
	return fmt.Sprintf("%.2f", MemUseMB())
}

func MemUseGB() float32 {
	return float32(MemUse()) / 1024 / 1024 / 1024
}

func MemUseGBStr() string {
	return fmt.Sprintf("%.2f", MemUseGB())
}

func PackageDateTime() string {
	buildTime := "unknown"
	if executable, err := os.Executable(); err != nil {
		logger.Log.Warn(fmt.Sprintf("获取可执行文件路径失败: %v", err))
	} else if info, err := os.Stat(executable); err != nil {
		logger.Log.Warn(fmt.Sprintf("读取可执行文件修改时间失败: %v", err))
	} else {
		buildTime = info.ModTime().Format("2006-01-02 15:04:05")
	}
	return buildTime
}
