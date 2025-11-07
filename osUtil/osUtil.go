package osUtil

import (
	"fmt"
	"runtime/metrics"
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
