package mapUtil

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

// BenchmarkResult 存储单个基准测试结果
type BenchmarkResult struct {
	MapType           string
	Operation         string
	DataSize          int
	Concurrency       int
	KeyType           string
	ValueType         string
	OperationsPerSec  float64
	AvgLatencyNs      int64
	MemAllocBytes     uint64
	MemAllocObjects   uint64
	GCCount           uint32
	GCPauseNs         uint64
	TotalMemBytes     uint64
}

// 全局结果收集器
var benchmarkResults []BenchmarkResult

// 记录基准测试结果
func recordResult(result BenchmarkResult) {
	benchmarkResults = append(benchmarkResults, result)
}

// ==================== 单线程性能测试 ====================

// ==================== sync.Map 性能测试 ====================

func BenchmarkSyncMap_Put_1w(b *testing.B) {
	benchmarkSyncMapPut(b, 10000, "string", "int")
}

func BenchmarkSyncMap_Put_10w(b *testing.B) {
	benchmarkSyncMapPut(b, 100000, "string", "int")
}

func BenchmarkSyncMap_Put_100w(b *testing.B) {
	benchmarkSyncMapPut(b, 1000000, "string", "int")
}

func benchmarkSyncMapPut(b *testing.B, size int, keyType, valueType string) {
	fmt.Printf("  → 执行 Put 操作测试: sync.Map | 数据规模: %d\n", size)
	var memStart, memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)

	b.ResetTimer()
	start := time.Now()

	for i := 0; i < b.N; i++ {
		var sm sync.Map
		for j := 0; j < size; j++ {
			sm.Store(fmt.Sprintf("key_%d", j), j)
		}
	}

	elapsed := time.Since(start)
	b.StopTimer()

	runtime.ReadMemStats(&memEnd)

	recordResult(BenchmarkResult{
		MapType:          "sync.Map",
		Operation:        "Put",
		DataSize:         size,
		Concurrency:      1,
		KeyType:          keyType,
		ValueType:        valueType,
		OperationsPerSec: float64(b.N*size) / elapsed.Seconds(),
		AvgLatencyNs:     elapsed.Nanoseconds() / int64(b.N*size),
		MemAllocBytes:    memEnd.TotalAlloc - memStart.TotalAlloc,
		MemAllocObjects:  memEnd.Mallocs - memStart.Mallocs,
		GCCount:          memEnd.NumGC - memStart.NumGC,
		GCPauseNs:        memEnd.PauseTotalNs - memStart.PauseTotalNs,
		TotalMemBytes:    memEnd.Alloc,
	})
}

// BenchmarkNativeMap_Put_1w 原生map写入1万条
func BenchmarkNativeMap_Put_1w(b *testing.B) {
	benchmarkNativeMapPut(b, 10000, "string", "int")
}

func BenchmarkNativeMap_Put_10w(b *testing.B) {
	benchmarkNativeMapPut(b, 100000, "string", "int")
}

func BenchmarkNativeMap_Put_100w(b *testing.B) {
	benchmarkNativeMapPut(b, 1000000, "string", "int")
}

func benchmarkNativeMapPut(b *testing.B, size int, keyType, valueType string) {
	fmt.Printf("  → 执行 Put 操作测试: NativeMap | 数据规模: %d\n", size)
	var memStart, memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)

	b.ResetTimer()
	start := time.Now()

	for i := 0; i < b.N; i++ {
		m := make(map[string]int)
		for j := 0; j < size; j++ {
			m[fmt.Sprintf("key_%d", j)] = j
		}
	}

	elapsed := time.Since(start)
	b.StopTimer()

	runtime.ReadMemStats(&memEnd)

	recordResult(BenchmarkResult{
		MapType:          "NativeMap",
		Operation:        "Put",
		DataSize:         size,
		Concurrency:      1,
		KeyType:          keyType,
		ValueType:        valueType,
		OperationsPerSec: float64(b.N*size) / elapsed.Seconds(),
		AvgLatencyNs:     elapsed.Nanoseconds() / int64(b.N*size),
		MemAllocBytes:    memEnd.TotalAlloc - memStart.TotalAlloc,
		MemAllocObjects:  memEnd.Mallocs - memStart.Mallocs,
		GCCount:          memEnd.NumGC - memStart.NumGC,
		GCPauseNs:        memEnd.PauseTotalNs - memStart.PauseTotalNs,
		TotalMemBytes:    memEnd.Alloc,
	})
}

// BenchmarkConcurrentHashMap_Put 系列
func BenchmarkConcurrentHashMap_Put_1w(b *testing.B) {
	benchmarkConcurrentHashMapPut(b, 10000, "string", "int")
}

func BenchmarkConcurrentHashMap_Put_10w(b *testing.B) {
	benchmarkConcurrentHashMapPut(b, 100000, "string", "int")
}

func BenchmarkConcurrentHashMap_Put_100w(b *testing.B) {
	benchmarkConcurrentHashMapPut(b, 1000000, "string", "int")
}

func benchmarkConcurrentHashMapPut(b *testing.B, size int, keyType, valueType string) {
	fmt.Printf("  → 执行 Put 操作测试: ConcurrentHashMap | 数据规模: %d\n", size)
	var memStart, memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)

	b.ResetTimer()
	start := time.Now()

	for i := 0; i < b.N; i++ {
		cm := NewConcurrentHashMap[string, int]()
		for j := 0; j < size; j++ {
			cm.Put(fmt.Sprintf("key_%d", j), j)
		}
	}

	elapsed := time.Since(start)
	b.StopTimer()

	runtime.ReadMemStats(&memEnd)

	recordResult(BenchmarkResult{
		MapType:          "ConcurrentHashMap",
		Operation:        "Put",
		DataSize:         size,
		Concurrency:      1,
		KeyType:          keyType,
		ValueType:        valueType,
		OperationsPerSec: float64(b.N*size) / elapsed.Seconds(),
		AvgLatencyNs:     elapsed.Nanoseconds() / int64(b.N*size),
		MemAllocBytes:    memEnd.TotalAlloc - memStart.TotalAlloc,
		MemAllocObjects:  memEnd.Mallocs - memStart.Mallocs,
		GCCount:          memEnd.NumGC - memStart.NumGC,
		GCPauseNs:        memEnd.PauseTotalNs - memStart.PauseTotalNs,
		TotalMemBytes:    memEnd.Alloc,
	})
}

// BenchmarkConcurrentSkipListMap_Put 系列
func BenchmarkConcurrentSkipListMap_Put_1w(b *testing.B) {
	benchmarkConcurrentSkipListMapPut(b, 10000, "int", "int")
}

func BenchmarkConcurrentSkipListMap_Put_10w(b *testing.B) {
	benchmarkConcurrentSkipListMapPut(b, 100000, "int", "int")
}

func BenchmarkConcurrentSkipListMap_Put_100w(b *testing.B) {
	benchmarkConcurrentSkipListMapPut(b, 1000000, "int", "int")
}

func benchmarkConcurrentSkipListMapPut(b *testing.B, size int, keyType, valueType string) {
	fmt.Printf("  → 执行 Put 操作测试: ConcurrentSkipListMap | 数据规模: %d\n", size)
	var memStart, memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)

	b.ResetTimer()
	start := time.Now()

	for i := 0; i < b.N; i++ {
		csm := NewConcurrentSkipListMap[int, int]()
		for j := 0; j < size; j++ {
			csm.Put(j, j)
		}
	}

	elapsed := time.Since(start)
	b.StopTimer()

	runtime.ReadMemStats(&memEnd)

	recordResult(BenchmarkResult{
		MapType:          "ConcurrentSkipListMap",
		Operation:        "Put",
		DataSize:         size,
		Concurrency:      1,
		KeyType:          keyType,
		ValueType:        valueType,
		OperationsPerSec: float64(b.N*size) / elapsed.Seconds(),
		AvgLatencyNs:     elapsed.Nanoseconds() / int64(b.N*size),
		MemAllocBytes:    memEnd.TotalAlloc - memStart.TotalAlloc,
		MemAllocObjects:  memEnd.Mallocs - memStart.Mallocs,
		GCCount:          memEnd.NumGC - memStart.NumGC,
		GCPauseNs:        memEnd.PauseTotalNs - memStart.PauseTotalNs,
		TotalMemBytes:    memEnd.Alloc,
	})
}

// BenchmarkOrderedMap_Put 系列
func BenchmarkOrderedMap_Put_1w(b *testing.B) {
	benchmarkOrderedMapPut(b, 10000, "string", "int")
}

func BenchmarkOrderedMap_Put_10w(b *testing.B) {
	benchmarkOrderedMapPut(b, 100000, "string", "int")
}

func BenchmarkOrderedMap_Put_100w(b *testing.B) {
	benchmarkOrderedMapPut(b, 1000000, "string", "int")
}

func benchmarkOrderedMapPut(b *testing.B, size int, keyType, valueType string) {
	fmt.Printf("  → 执行 Put 操作测试: OrderedMap | 数据规模: %d\n", size)
	var memStart, memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)

	b.ResetTimer()
	start := time.Now()

	for i := 0; i < b.N; i++ {
		om := NewOrderedMap[string, int]()
		for j := 0; j < size; j++ {
			om.Put(fmt.Sprintf("key_%d", j), j)
		}
	}

	elapsed := time.Since(start)
	b.StopTimer()

	runtime.ReadMemStats(&memEnd)

	recordResult(BenchmarkResult{
		MapType:          "OrderedMap",
		Operation:        "Put",
		DataSize:         size,
		Concurrency:      1,
		KeyType:          keyType,
		ValueType:        valueType,
		OperationsPerSec: float64(b.N*size) / elapsed.Seconds(),
		AvgLatencyNs:     elapsed.Nanoseconds() / int64(b.N*size),
		MemAllocBytes:    memEnd.TotalAlloc - memStart.TotalAlloc,
		MemAllocObjects:  memEnd.Mallocs - memStart.Mallocs,
		GCCount:          memEnd.NumGC - memStart.NumGC,
		GCPauseNs:        memEnd.PauseTotalNs - memStart.PauseTotalNs,
		TotalMemBytes:    memEnd.Alloc,
	})
}

// BenchmarkTreeMap_Put 系列
func BenchmarkTreeMap_Put_1w(b *testing.B) {
	benchmarkTreeMapPut(b, 10000, "int", "int")
}

func BenchmarkTreeMap_Put_10w(b *testing.B) {
	benchmarkTreeMapPut(b, 100000, "int", "int")
}

func BenchmarkTreeMap_Put_100w(b *testing.B) {
	benchmarkTreeMapPut(b, 1000000, "int", "int")
}

func benchmarkTreeMapPut(b *testing.B, size int, keyType, valueType string) {
	fmt.Printf("  → 执行 Put 操作测试: TreeMap | 数据规模: %d\n", size)
	var memStart, memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)

	b.ResetTimer()
	start := time.Now()

	for i := 0; i < b.N; i++ {
		tm := NewTreeMap[int, int](func(a, b int) bool { return a < b })
		for j := 0; j < size; j++ {
			tm.Put(j, j)
		}
	}

	elapsed := time.Since(start)
	b.StopTimer()

	runtime.ReadMemStats(&memEnd)

	recordResult(BenchmarkResult{
		MapType:          "TreeMap",
		Operation:        "Put",
		DataSize:         size,
		Concurrency:      1,
		KeyType:          keyType,
		ValueType:        valueType,
		OperationsPerSec: float64(b.N*size) / elapsed.Seconds(),
		AvgLatencyNs:     elapsed.Nanoseconds() / int64(b.N*size),
		MemAllocBytes:    memEnd.TotalAlloc - memStart.TotalAlloc,
		MemAllocObjects:  memEnd.Mallocs - memStart.Mallocs,
		GCCount:          memEnd.NumGC - memStart.NumGC,
		GCPauseNs:        memEnd.PauseTotalNs - memStart.PauseTotalNs,
		TotalMemBytes:    memEnd.Alloc,
	})
}

// ==================== Get 性能测试 ====================

func BenchmarkSyncMap_Get_10w(b *testing.B) {
	fmt.Printf("  → 执行 Get 操作测试: sync.Map | 数据规模: 100000\n")
	var sm sync.Map
	for j := 0; j < 100000; j++ {
		sm.Store(fmt.Sprintf("key_%d", j), j)
	}

	var memStart, memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)

	b.ResetTimer()
	start := time.Now()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 100000; j++ {
			_, _ = sm.Load(fmt.Sprintf("key_%d", j))
		}
	}

	elapsed := time.Since(start)
	b.StopTimer()

	runtime.ReadMemStats(&memEnd)

	recordResult(BenchmarkResult{
		MapType:          "sync.Map",
		Operation:        "Get",
		DataSize:         100000,
		Concurrency:      1,
		KeyType:          "string",
		ValueType:        "int",
		OperationsPerSec: float64(b.N*100000) / elapsed.Seconds(),
		AvgLatencyNs:     elapsed.Nanoseconds() / int64(b.N*100000),
		MemAllocBytes:    memEnd.TotalAlloc - memStart.TotalAlloc,
		MemAllocObjects:  memEnd.Mallocs - memStart.Mallocs,
		GCCount:          memEnd.NumGC - memStart.NumGC,
		GCPauseNs:        memEnd.PauseTotalNs - memStart.PauseTotalNs,
		TotalMemBytes:    memEnd.Alloc,
	})
}

func BenchmarkNativeMap_Get_10w(b *testing.B) {
	benchmarkNativeMapGet(b, 100000)
}

func benchmarkNativeMapGet(b *testing.B, size int) {
	fmt.Printf("  → 执行 Get 操作测试: NativeMap | 数据规模: %d\n", size)
	m := make(map[string]int)
	for j := 0; j < size; j++ {
		m[fmt.Sprintf("key_%d", j)] = j
	}

	var memStart, memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)

	b.ResetTimer()
	start := time.Now()

	for i := 0; i < b.N; i++ {
		for j := 0; j < size; j++ {
			_ = m[fmt.Sprintf("key_%d", j)]
		}
	}

	elapsed := time.Since(start)
	b.StopTimer()

	runtime.ReadMemStats(&memEnd)

	recordResult(BenchmarkResult{
		MapType:          "NativeMap",
		Operation:        "Get",
		DataSize:         size,
		Concurrency:      1,
		KeyType:          "string",
		ValueType:        "int",
		OperationsPerSec: float64(b.N*size) / elapsed.Seconds(),
		AvgLatencyNs:     elapsed.Nanoseconds() / int64(b.N*size),
		MemAllocBytes:    memEnd.TotalAlloc - memStart.TotalAlloc,
		MemAllocObjects:  memEnd.Mallocs - memStart.Mallocs,
		GCCount:          memEnd.NumGC - memStart.NumGC,
		GCPauseNs:        memEnd.PauseTotalNs - memStart.PauseTotalNs,
		TotalMemBytes:    memEnd.Alloc,
	})
}

func BenchmarkConcurrentHashMap_Get_10w(b *testing.B) {
	fmt.Printf("  → 执行 Get 操作测试: ConcurrentHashMap | 数据规模: 100000\n")
	cm := NewConcurrentHashMap[string, int]()
	for j := 0; j < 100000; j++ {
		cm.Put(fmt.Sprintf("key_%d", j), j)
	}

	var memStart, memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)

	b.ResetTimer()
	start := time.Now()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 100000; j++ {
			_ = cm.Get(fmt.Sprintf("key_%d", j))
		}
	}

	elapsed := time.Since(start)
	b.StopTimer()

	runtime.ReadMemStats(&memEnd)

	recordResult(BenchmarkResult{
		MapType:          "ConcurrentHashMap",
		Operation:        "Get",
		DataSize:         100000,
		Concurrency:      1,
		KeyType:          "string",
		ValueType:        "int",
		OperationsPerSec: float64(b.N*100000) / elapsed.Seconds(),
		AvgLatencyNs:     elapsed.Nanoseconds() / int64(b.N*100000),
		MemAllocBytes:    memEnd.TotalAlloc - memStart.TotalAlloc,
		MemAllocObjects:  memEnd.Mallocs - memStart.Mallocs,
		GCCount:          memEnd.NumGC - memStart.NumGC,
		GCPauseNs:        memEnd.PauseTotalNs - memStart.PauseTotalNs,
		TotalMemBytes:    memEnd.Alloc,
	})
}

func BenchmarkConcurrentSkipListMap_Get_10w(b *testing.B) {
	fmt.Printf("  → 执行 Get 操作测试: ConcurrentSkipListMap | 数据规模: 100000\n")
	csm := NewConcurrentSkipListMap[int, int]()
	for j := 0; j < 100000; j++ {
		csm.Put(j, j)
	}

	var memStart, memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)

	b.ResetTimer()
	start := time.Now()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 100000; j++ {
			_ = csm.Get(j)
		}
	}

	elapsed := time.Since(start)
	b.StopTimer()

	runtime.ReadMemStats(&memEnd)

	recordResult(BenchmarkResult{
		MapType:          "ConcurrentSkipListMap",
		Operation:        "Get",
		DataSize:         100000,
		Concurrency:      1,
		KeyType:          "int",
		ValueType:        "int",
		OperationsPerSec: float64(b.N*100000) / elapsed.Seconds(),
		AvgLatencyNs:     elapsed.Nanoseconds() / int64(b.N*100000),
		MemAllocBytes:    memEnd.TotalAlloc - memStart.TotalAlloc,
		MemAllocObjects:  memEnd.Mallocs - memStart.Mallocs,
		GCCount:          memEnd.NumGC - memStart.NumGC,
		GCPauseNs:        memEnd.PauseTotalNs - memStart.PauseTotalNs,
		TotalMemBytes:    memEnd.Alloc,
	})
}

func BenchmarkOrderedMap_Get_10w(b *testing.B) {
	fmt.Printf("  → 执行 Get 操作测试: OrderedMap | 数据规模: 100000\n")
	om := NewOrderedMap[string, int]()
	for j := 0; j < 100000; j++ {
		om.Put(fmt.Sprintf("key_%d", j), j)
	}

	var memStart, memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)

	b.ResetTimer()
	start := time.Now()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 100000; j++ {
			_ = om.Get(fmt.Sprintf("key_%d", j))
		}
	}

	elapsed := time.Since(start)
	b.StopTimer()

	runtime.ReadMemStats(&memEnd)

	recordResult(BenchmarkResult{
		MapType:          "OrderedMap",
		Operation:        "Get",
		DataSize:         100000,
		Concurrency:      1,
		KeyType:          "string",
		ValueType:        "int",
		OperationsPerSec: float64(b.N*100000) / elapsed.Seconds(),
		AvgLatencyNs:     elapsed.Nanoseconds() / int64(b.N*100000),
		MemAllocBytes:    memEnd.TotalAlloc - memStart.TotalAlloc,
		MemAllocObjects:  memEnd.Mallocs - memStart.Mallocs,
		GCCount:          memEnd.NumGC - memStart.NumGC,
		GCPauseNs:        memEnd.PauseTotalNs - memStart.PauseTotalNs,
		TotalMemBytes:    memEnd.Alloc,
	})
}

func BenchmarkTreeMap_Get_10w(b *testing.B) {
	fmt.Printf("  → 执行 Get 操作测试: TreeMap | 数据规模: 100000\n")
	tm := NewTreeMap[int, int](func(a, b int) bool { return a < b })
	for j := 0; j < 100000; j++ {
		tm.Put(j, j)
	}

	var memStart, memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)

	b.ResetTimer()
	start := time.Now()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 100000; j++ {
			_ = tm.Get(j)
		}
	}

	elapsed := time.Since(start)
	b.StopTimer()

	runtime.ReadMemStats(&memEnd)

	recordResult(BenchmarkResult{
		MapType:          "TreeMap",
		Operation:        "Get",
		DataSize:         100000,
		Concurrency:      1,
		KeyType:          "int",
		ValueType:        "int",
		OperationsPerSec: float64(b.N*100000) / elapsed.Seconds(),
		AvgLatencyNs:     elapsed.Nanoseconds() / int64(b.N*100000),
		MemAllocBytes:    memEnd.TotalAlloc - memStart.TotalAlloc,
		MemAllocObjects:  memEnd.Mallocs - memStart.Mallocs,
		GCCount:          memEnd.NumGC - memStart.NumGC,
		GCPauseNs:        memEnd.PauseTotalNs - memStart.PauseTotalNs,
		TotalMemBytes:    memEnd.Alloc,
	})
}

// ==================== 并发写入测试 ====================

func BenchmarkSyncMap_ConcurrentPut_10Goroutines(b *testing.B) {
	benchmarkSyncMapConcurrentPut(b, 10000, 10)
}

func BenchmarkSyncMap_ConcurrentPut_100Goroutines(b *testing.B) {
	benchmarkSyncMapConcurrentPut(b, 10000, 100)
}

func benchmarkSyncMapConcurrentPut(b *testing.B, size, goroutines int) {
	fmt.Printf("  → 执行 ConcurrentPut 操作测试: sync.Map | 数据规模: %d | 并发度: %d\n", size, goroutines)
	var memStart, memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)

	b.ResetTimer()
	start := time.Now()

	for i := 0; i < b.N; i++ {
		var sm sync.Map
		var wg sync.WaitGroup
		wg.Add(goroutines)

		for g := 0; g < goroutines; g++ {
			go func(gid int) {
				defer wg.Done()
				for j := 0; j < size/goroutines; j++ {
					sm.Store(fmt.Sprintf("key_%d_%d", gid, j), j)
				}
			}(g)
		}
		wg.Wait()
	}

	elapsed := time.Since(start)
	b.StopTimer()

	runtime.ReadMemStats(&memEnd)

	recordResult(BenchmarkResult{
		MapType:          "sync.Map",
		Operation:        "ConcurrentPut",
		DataSize:         size,
		Concurrency:      goroutines,
		KeyType:          "string",
		ValueType:        "int",
		OperationsPerSec: float64(b.N*size) / elapsed.Seconds(),
		AvgLatencyNs:     elapsed.Nanoseconds() / int64(b.N*size),
		MemAllocBytes:    memEnd.TotalAlloc - memStart.TotalAlloc,
		MemAllocObjects:  memEnd.Mallocs - memStart.Mallocs,
		GCCount:          memEnd.NumGC - memStart.NumGC,
		GCPauseNs:        memEnd.PauseTotalNs - memStart.PauseTotalNs,
		TotalMemBytes:    memEnd.Alloc,
	})
}

func BenchmarkConcurrentHashMap_ConcurrentPut_10Goroutines(b *testing.B) {
	benchmarkConcurrentHashMapConcurrentPut(b, 10000, 10)
}

func BenchmarkConcurrentHashMap_ConcurrentPut_100Goroutines(b *testing.B) {
	benchmarkConcurrentHashMapConcurrentPut(b, 10000, 100)
}

func benchmarkConcurrentHashMapConcurrentPut(b *testing.B, size, goroutines int) {
	fmt.Printf("  → 执行 ConcurrentPut 操作测试: ConcurrentHashMap | 数据规模: %d | 并发度: %d\n", size, goroutines)
	var memStart, memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)

	b.ResetTimer()
	start := time.Now()

	for i := 0; i < b.N; i++ {
		cm := NewConcurrentHashMap[string, int]()
		var wg sync.WaitGroup
		wg.Add(goroutines)

		for g := 0; g < goroutines; g++ {
			go func(gid int) {
				defer wg.Done()
				for j := 0; j < size/goroutines; j++ {
					cm.Put(fmt.Sprintf("key_%d_%d", gid, j), j)
				}
			}(g)
		}
		wg.Wait()
	}

	elapsed := time.Since(start)
	b.StopTimer()

	runtime.ReadMemStats(&memEnd)

	recordResult(BenchmarkResult{
		MapType:          "ConcurrentHashMap",
		Operation:        "ConcurrentPut",
		DataSize:         size,
		Concurrency:      goroutines,
		KeyType:          "string",
		ValueType:        "int",
		OperationsPerSec: float64(b.N*size) / elapsed.Seconds(),
		AvgLatencyNs:     elapsed.Nanoseconds() / int64(b.N*size),
		MemAllocBytes:    memEnd.TotalAlloc - memStart.TotalAlloc,
		MemAllocObjects:  memEnd.Mallocs - memStart.Mallocs,
		GCCount:          memEnd.NumGC - memStart.NumGC,
		GCPauseNs:        memEnd.PauseTotalNs - memStart.PauseTotalNs,
		TotalMemBytes:    memEnd.Alloc,
	})
}

func BenchmarkConcurrentSkipListMap_ConcurrentPut_10Goroutines(b *testing.B) {
	benchmarkConcurrentSkipListMapConcurrentPut(b, 10000, 10)
}

func BenchmarkConcurrentSkipListMap_ConcurrentPut_100Goroutines(b *testing.B) {
	benchmarkConcurrentSkipListMapConcurrentPut(b, 10000, 100)
}

func benchmarkConcurrentSkipListMapConcurrentPut(b *testing.B, size, goroutines int) {
	fmt.Printf("  → 执行 ConcurrentPut 操作测试: ConcurrentSkipListMap | 数据规模: %d | 并发度: %d\n", size, goroutines)
	var memStart, memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)

	b.ResetTimer()
	start := time.Now()

	for i := 0; i < b.N; i++ {
		csm := NewConcurrentSkipListMap[int, int]()
		var wg sync.WaitGroup
		wg.Add(goroutines)

		for g := 0; g < goroutines; g++ {
			go func(gid int) {
				defer wg.Done()
				for j := 0; j < size/goroutines; j++ {
					csm.Put(gid*size+j, j)
				}
			}(g)
		}
		wg.Wait()
	}

	elapsed := time.Since(start)
	b.StopTimer()

	runtime.ReadMemStats(&memEnd)

	recordResult(BenchmarkResult{
		MapType:          "ConcurrentSkipListMap",
		Operation:        "ConcurrentPut",
		DataSize:         size,
		Concurrency:      goroutines,
		KeyType:          "int",
		ValueType:        "int",
		OperationsPerSec: float64(b.N*size) / elapsed.Seconds(),
		AvgLatencyNs:     elapsed.Nanoseconds() / int64(b.N*size),
		MemAllocBytes:    memEnd.TotalAlloc - memStart.TotalAlloc,
		MemAllocObjects:  memEnd.Mallocs - memStart.Mallocs,
		GCCount:          memEnd.NumGC - memStart.NumGC,
		GCPauseNs:        memEnd.PauseTotalNs - memStart.PauseTotalNs,
		TotalMemBytes:    memEnd.Alloc,
	})
}

func BenchmarkTreeMap_ConcurrentPut_10Goroutines(b *testing.B) {
	benchmarkTreeMapConcurrentPut(b, 10000, 10)
}

func benchmarkTreeMapConcurrentPut(b *testing.B, size, goroutines int) {
	fmt.Printf("  → 执行 ConcurrentPut 操作测试: TreeMap | 数据规模: %d | 并发度: %d\n", size, goroutines)
	var memStart, memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)

	b.ResetTimer()
	start := time.Now()

	for i := 0; i < b.N; i++ {
		tm := NewTreeMap[int, int](func(a, b int) bool { return a < b })
		var wg sync.WaitGroup
		wg.Add(goroutines)

		for g := 0; g < goroutines; g++ {
			go func(gid int) {
				defer wg.Done()
				for j := 0; j < size/goroutines; j++ {
					tm.Put(gid*size+j, j)
				}
			}(g)
		}
		wg.Wait()
	}

	elapsed := time.Since(start)
	b.StopTimer()

	runtime.ReadMemStats(&memEnd)

	recordResult(BenchmarkResult{
		MapType:          "TreeMap",
		Operation:        "ConcurrentPut",
		DataSize:         size,
		Concurrency:      goroutines,
		KeyType:          "int",
		ValueType:        "int",
		OperationsPerSec: float64(b.N*size) / elapsed.Seconds(),
		AvgLatencyNs:     elapsed.Nanoseconds() / int64(b.N*size),
		MemAllocBytes:    memEnd.TotalAlloc - memStart.TotalAlloc,
		MemAllocObjects:  memEnd.Mallocs - memStart.Mallocs,
		GCCount:          memEnd.NumGC - memStart.NumGC,
		GCPauseNs:        memEnd.PauseTotalNs - memStart.PauseTotalNs,
		TotalMemBytes:    memEnd.Alloc,
	})
}

// ==================== 遍历性能测试 ====================

func BenchmarkSyncMap_Range_10w(b *testing.B) {
	fmt.Printf("  → 执行 Range 操作测试: sync.Map | 数据规模: 100000\n")
	var sm sync.Map
	for j := 0; j < 100000; j++ {
		sm.Store(fmt.Sprintf("key_%d", j), j)
	}

	var memStart, memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)

	b.ResetTimer()
	start := time.Now()

	for i := 0; i < b.N; i++ {
		count := 0
		sm.Range(func(key, value interface{}) bool {
			count++
			return true
		})
	}

	elapsed := time.Since(start)
	b.StopTimer()

	runtime.ReadMemStats(&memEnd)

	recordResult(BenchmarkResult{
		MapType:          "sync.Map",
		Operation:        "Range",
		DataSize:         100000,
		Concurrency:      1,
		KeyType:          "string",
		ValueType:        "int",
		OperationsPerSec: float64(b.N*100000) / elapsed.Seconds(),
		AvgLatencyNs:     elapsed.Nanoseconds() / int64(b.N*100000),
		MemAllocBytes:    memEnd.TotalAlloc - memStart.TotalAlloc,
		MemAllocObjects:  memEnd.Mallocs - memStart.Mallocs,
		GCCount:          memEnd.NumGC - memStart.NumGC,
		GCPauseNs:        memEnd.PauseTotalNs - memStart.PauseTotalNs,
		TotalMemBytes:    memEnd.Alloc,
	})
}

func BenchmarkNativeMap_Range_10w(b *testing.B) {
	fmt.Printf("  → 执行 Range 操作测试: NativeMap | 数据规模: 100000\n")
	m := make(map[string]int)
	for j := 0; j < 100000; j++ {
		m[fmt.Sprintf("key_%d", j)] = j
	}

	var memStart, memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)

	b.ResetTimer()
	start := time.Now()

	for i := 0; i < b.N; i++ {
		count := 0
		for range m {
			count++
		}
	}

	elapsed := time.Since(start)
	b.StopTimer()

	runtime.ReadMemStats(&memEnd)

	recordResult(BenchmarkResult{
		MapType:          "NativeMap",
		Operation:        "Range",
		DataSize:         100000,
		Concurrency:      1,
		KeyType:          "string",
		ValueType:        "int",
		OperationsPerSec: float64(b.N*100000) / elapsed.Seconds(),
		AvgLatencyNs:     elapsed.Nanoseconds() / int64(b.N*100000),
		MemAllocBytes:    memEnd.TotalAlloc - memStart.TotalAlloc,
		MemAllocObjects:  memEnd.Mallocs - memStart.Mallocs,
		GCCount:          memEnd.NumGC - memStart.NumGC,
		GCPauseNs:        memEnd.PauseTotalNs - memStart.PauseTotalNs,
		TotalMemBytes:    memEnd.Alloc,
	})
}

func BenchmarkConcurrentHashMap_Range_10w(b *testing.B) {
	fmt.Printf("  → 执行 Range 操作测试: ConcurrentHashMap | 数据规模: 100000\n")
	cm := NewConcurrentHashMap[string, int]()
	for j := 0; j < 100000; j++ {
		cm.Put(fmt.Sprintf("key_%d", j), j)
	}

	var memStart, memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)

	b.ResetTimer()
	start := time.Now()

	for i := 0; i < b.N; i++ {
		count := 0
		cm.Range(func(key string, value int) bool {
			count++
			return true
		})
	}

	elapsed := time.Since(start)
	b.StopTimer()

	runtime.ReadMemStats(&memEnd)

	recordResult(BenchmarkResult{
		MapType:          "ConcurrentHashMap",
		Operation:        "Range",
		DataSize:         100000,
		Concurrency:      1,
		KeyType:          "string",
		ValueType:        "int",
		OperationsPerSec: float64(b.N*100000) / elapsed.Seconds(),
		AvgLatencyNs:     elapsed.Nanoseconds() / int64(b.N*100000),
		MemAllocBytes:    memEnd.TotalAlloc - memStart.TotalAlloc,
		MemAllocObjects:  memEnd.Mallocs - memStart.Mallocs,
		GCCount:          memEnd.NumGC - memStart.NumGC,
		GCPauseNs:        memEnd.PauseTotalNs - memStart.PauseTotalNs,
		TotalMemBytes:    memEnd.Alloc,
	})
}

func BenchmarkOrderedMap_Range_10w(b *testing.B) {
	fmt.Printf("  → 执行 Range 操作测试: OrderedMap | 数据规模: 100000\n")
	om := NewOrderedMap[string, int]()
	for j := 0; j < 100000; j++ {
		om.Put(fmt.Sprintf("key_%d", j), j)
	}

	var memStart, memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)

	b.ResetTimer()
	start := time.Now()

	for i := 0; i < b.N; i++ {
		count := 0
		om.Range(func(key string, value int) bool {
			count++
			return true
		})
	}

	elapsed := time.Since(start)
	b.StopTimer()

	runtime.ReadMemStats(&memEnd)

	recordResult(BenchmarkResult{
		MapType:          "OrderedMap",
		Operation:        "Range",
		DataSize:         100000,
		Concurrency:      1,
		KeyType:          "string",
		ValueType:        "int",
		OperationsPerSec: float64(b.N*100000) / elapsed.Seconds(),
		AvgLatencyNs:     elapsed.Nanoseconds() / int64(b.N*100000),
		MemAllocBytes:    memEnd.TotalAlloc - memStart.TotalAlloc,
		MemAllocObjects:  memEnd.Mallocs - memStart.Mallocs,
		GCCount:          memEnd.NumGC - memStart.NumGC,
		GCPauseNs:        memEnd.PauseTotalNs - memStart.PauseTotalNs,
		TotalMemBytes:    memEnd.Alloc,
	})
}

// ==================== 混合操作测试 ====================

func BenchmarkSyncMap_Mixed_10w(b *testing.B) {
	fmt.Printf("  → 执行 Mixed 操作测试: sync.Map | 数据规模: 100000\n")
	var memStart, memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)

	b.ResetTimer()
	start := time.Now()

	size := 100000
	for i := 0; i < b.N; i++ {
		var sm sync.Map

		// 70% 写，20% 读，10% 删
		for j := 0; j < size; j++ {
			op := rand.Intn(10)
			key := fmt.Sprintf("key_%d", rand.Intn(size))

			if op < 7 {
				sm.Store(key, j)
			} else if op < 9 {
				_, _ = sm.Load(key)
			} else {
				sm.Delete(key)
			}
		}
	}

	elapsed := time.Since(start)
	b.StopTimer()

	runtime.ReadMemStats(&memEnd)

	recordResult(BenchmarkResult{
		MapType:          "sync.Map",
		Operation:        "Mixed",
		DataSize:         size,
		Concurrency:      1,
		KeyType:          "string",
		ValueType:        "int",
		OperationsPerSec: float64(b.N*size) / elapsed.Seconds(),
		AvgLatencyNs:     elapsed.Nanoseconds() / int64(b.N*size),
		MemAllocBytes:    memEnd.TotalAlloc - memStart.TotalAlloc,
		MemAllocObjects:  memEnd.Mallocs - memStart.Mallocs,
		GCCount:          memEnd.NumGC - memStart.NumGC,
		GCPauseNs:        memEnd.PauseTotalNs - memStart.PauseTotalNs,
		TotalMemBytes:    memEnd.Alloc,
	})
}

func BenchmarkNativeMap_Mixed_10w(b *testing.B) {
	fmt.Printf("  → 执行 Mixed 操作测试: NativeMap | 数据规模: 100000\n")
	var memStart, memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)

	b.ResetTimer()
	start := time.Now()

	size := 100000
	for i := 0; i < b.N; i++ {
		m := make(map[string]int)

		// 70% 写，20% 读，10% 删
		for j := 0; j < size; j++ {
			op := rand.Intn(10)
			key := fmt.Sprintf("key_%d", rand.Intn(size))

			if op < 7 {
				m[key] = j
			} else if op < 9 {
				_ = m[key]
			} else {
				delete(m, key)
			}
		}
	}

	elapsed := time.Since(start)
	b.StopTimer()

	runtime.ReadMemStats(&memEnd)

	recordResult(BenchmarkResult{
		MapType:          "NativeMap",
		Operation:        "Mixed",
		DataSize:         size,
		Concurrency:      1,
		KeyType:          "string",
		ValueType:        "int",
		OperationsPerSec: float64(b.N*size) / elapsed.Seconds(),
		AvgLatencyNs:     elapsed.Nanoseconds() / int64(b.N*size),
		MemAllocBytes:    memEnd.TotalAlloc - memStart.TotalAlloc,
		MemAllocObjects:  memEnd.Mallocs - memStart.Mallocs,
		GCCount:          memEnd.NumGC - memStart.NumGC,
		GCPauseNs:        memEnd.PauseTotalNs - memStart.PauseTotalNs,
		TotalMemBytes:    memEnd.Alloc,
	})
}

func BenchmarkConcurrentHashMap_Mixed_10w(b *testing.B) {
	fmt.Printf("  → 执行 Mixed 操作测试: ConcurrentHashMap | 数据规模: 100000\n")
	var memStart, memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)

	b.ResetTimer()
	start := time.Now()

	size := 100000
	for i := 0; i < b.N; i++ {
		cm := NewConcurrentHashMap[string, int]()

		for j := 0; j < size; j++ {
			op := rand.Intn(10)
			key := fmt.Sprintf("key_%d", rand.Intn(size))

			if op < 7 {
				cm.Put(key, j)
			} else if op < 9 {
				_ = cm.Get(key)
			} else {
				cm.Remove(key)
			}
		}
	}

	elapsed := time.Since(start)
	b.StopTimer()

	runtime.ReadMemStats(&memEnd)

	recordResult(BenchmarkResult{
		MapType:          "ConcurrentHashMap",
		Operation:        "Mixed",
		DataSize:         size,
		Concurrency:      1,
		KeyType:          "string",
		ValueType:        "int",
		OperationsPerSec: float64(b.N*size) / elapsed.Seconds(),
		AvgLatencyNs:     elapsed.Nanoseconds() / int64(b.N*size),
		MemAllocBytes:    memEnd.TotalAlloc - memStart.TotalAlloc,
		MemAllocObjects:  memEnd.Mallocs - memStart.Mallocs,
		GCCount:          memEnd.NumGC - memStart.NumGC,
		GCPauseNs:        memEnd.PauseTotalNs - memStart.PauseTotalNs,
		TotalMemBytes:    memEnd.Alloc,
	})
}
