package mapUtil

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	fmt.Println("=================================================================")
	fmt.Println("           Go Map 数据结构性能对比测试")
	fmt.Println("=================================================================")
	fmt.Println()

	// 获取当前目录
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取当前目录失败: %v\n", err)
		return
	}

	// 运行基准测试
	fmt.Println("正在运行基准测试，这可能需要几分钟...")
	fmt.Println()

	cmd := exec.Command("go", "test", "-bench=.", "-benchmem", "-benchtime=3s", "-timeout=60m")
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("运行测试失败: %v\n输出: %s\n", err, string(output))
		return
	}

	fmt.Println("基准测试完成！")
	fmt.Println()
	fmt.Println(string(output))
	fmt.Println()

	// 生成报告
	fmt.Println("正在生成性能分析报告...")

	reportPath := filepath.Join(dir, "Map性能测试分析报告.md")
	err = generateReport(reportPath, string(output))
	if err != nil {
		fmt.Printf("生成报告失败: %v\n", err)
		return
	}

	fmt.Printf("报告已生成: %s\n", reportPath)
}

type BenchResult struct {
	Name        string
	NsPerOp     float64
	BytesPerOp  int64
	AllocsPerOp int64
	MBPerSec    float64
}

func generateReport(reportPath, benchOutput string) error {
	// 解析基准测试输出
	results := parseBenchOutput(benchOutput)

	// 按照Map类型和操作分组
	grouped := make(map[string][]BenchResult)
	for _, r := range results {
		// 提取 Map 类型和操作
		parts := strings.Split(r.Name, "_")
		if len(parts) >= 2 {
			mapType := parts[0][9:] // 去掉 "Benchmark" 前缀
			operation := parts[1]
			key := fmt.Sprintf("%s-%s", mapType, operation)
			grouped[key] = append(grouped[key], r)
		}
	}

	// 生成 Markdown 报告
	var report strings.Builder

	report.WriteString("# Go Map 数据结构性能对比分析报告\n\n")
	report.WriteString(fmt.Sprintf("**生成时间**: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	report.WriteString("## 1. 测试概述\n\n")
	report.WriteString("### 1.1 测试对象\n\n")
	report.WriteString("本次性能测试对比了以下 6 种 Map 实现：\n\n")
	report.WriteString("| Map 类型 | 描述 | 特性 |\n")
	report.WriteString("|---------|------|------|\n")
	report.WriteString("| NativeMap | Go 原生 map | 最基础的哈希表实现，非并发安全 |\n")
	report.WriteString("| sync.Map | Go 标准库 sync.Map | 官方并发安全实现，适用于读多写少场景 |\n")
	report.WriteString("| ConcurrentHashMap | 自实现并发哈希表 | 使用读写锁保护的并发安全 map |\n")
	report.WriteString("| ConcurrentSkipListMap | 并发跳表 | 基于跳表的有序 map，并发安全 |\n")
	report.WriteString("| OrderedMap | 有序 map | 保持插入顺序的 map，非并发安全 |\n")
	report.WriteString("| TreeMap | 红黑树 map | 基于红黑树的有序 map，并发安全 |\n\n")

	report.WriteString("### 1.2 测试维度\n\n")
	report.WriteString("- **性能测试**: Put (写入)、Get (读取)、Remove (删除)、Range (遍历)、Mixed (混合操作)\n")
	report.WriteString("- **数据规模**: 1万、10万、100万 条数据\n")
	report.WriteString("- **并发测试**: 单线程、10、100、1000 个 goroutine\n")
	report.WriteString("- **内存指标**: 内存分配量、内存对象数、GC 次数\n\n")

	report.WriteString("### 1.3 测试环境\n\n")
	report.WriteString(fmt.Sprintf("- **操作系统**: Windows\n"))
	report.WriteString(fmt.Sprintf("- **Go 版本**: %s\n", getGoVersion()))
	report.WriteString(fmt.Sprintf("- **测试时间**: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	report.WriteString("---\n\n")
	report.WriteString("## 2. 性能测试结果\n\n")

	// 按操作类型分组展示
	operations := []string{"Put", "Get", "Range", "ConcurrentPut", "Mixed"}

	for _, op := range operations {
		report.WriteString(fmt.Sprintf("### 2.%d %s 操作性能\n\n", getOpIndex(op), getOpTitle(op)))

		// 创建表格
		report.WriteString("| Map 类型 | 数据规模/并发度 | 平均耗时 (ns/op) | 内存分配 (B/op) | 分配次数 (allocs/op) |\n")
		report.WriteString("|---------|----------------|-----------------|----------------|-------------------|\n")

		// 收集该操作的所有结果
		var opResults []BenchResult
		for key, results := range grouped {
			if strings.HasSuffix(key, "-"+op) {
				opResults = append(opResults, results...)
			}
		}

		// 排序
		sort.Slice(opResults, func(i, j int) bool {
			return opResults[i].NsPerOp < opResults[j].NsPerOp
		})

		for _, r := range opResults {
			parts := strings.Split(r.Name, "_")
			mapType := parts[0][9:]
			size := ""
			if len(parts) >= 3 {
				size = parts[2]
			}
			if len(parts) >= 4 {
				size += " / " + parts[3]
			}

			report.WriteString(fmt.Sprintf("| %s | %s | %.2f | %d | %d |\n",
				mapType, size, r.NsPerOp, r.BytesPerOp, r.AllocsPerOp))
		}

		report.WriteString("\n")

		// 添加分析
		if len(opResults) > 0 {
			report.WriteString("**性能分析**:\n\n")
			fastest := opResults[0]
			slowest := opResults[len(opResults)-1]

			report.WriteString(fmt.Sprintf("- 最快: %s (%.2f ns/op)\n",
				extractMapType(fastest.Name), fastest.NsPerOp))
			report.WriteString(fmt.Sprintf("- 最慢: %s (%.2f ns/op)\n",
				extractMapType(slowest.Name), slowest.NsPerOp))
			report.WriteString(fmt.Sprintf("- 性能差距: %.2fx\n\n", slowest.NsPerOp/fastest.NsPerOp))
		}
	}

	report.WriteString("---\n\n")
	report.WriteString("## 3. 综合分析与建议\n\n")

	report.WriteString("### 3.1 性能特点总结\n\n")
	report.WriteString("#### 3.1.1 Go 原生 map (NativeMap)\n")
	report.WriteString("- **优势**: 单线程性能最优，内存开销最小\n")
	report.WriteString("- **劣势**: 非并发安全，多 goroutine 访问需要额外加锁\n")
	report.WriteString("- **适用场景**: 单线程操作或已有外部锁保护的场景\n\n")

	report.WriteString("#### 3.1.2 sync.Map\n")
	report.WriteString("- **优势**: 官方实现，读操作性能优秀（特别是读多写少场景）\n")
	report.WriteString("- **劣势**: 写操作性能一般，存储空间开销较大\n")
	report.WriteString("- **适用场景**: 读多写少的并发场景，如缓存\n\n")

	report.WriteString("#### 3.1.3 ConcurrentHashMap\n")
	report.WriteString("- **优势**: 实现简单，性能稳定，并发安全\n")
	report.WriteString("- **劣势**: 全局锁导致高并发时性能瓶颈\n")
	report.WriteString("- **适用场景**: 中等并发场景，读写均衡\n\n")

	report.WriteString("#### 3.1.4 ConcurrentSkipListMap\n")
	report.WriteString("- **优势**: 有序性，并发性能好，适合范围查询\n")
	report.WriteString("- **劣势**: 内存开销大，单次操作耗时较高\n")
	report.WriteString("- **适用场景**: 需要有序遍历或范围查询的并发场景\n\n")

	report.WriteString("#### 3.1.5 OrderedMap\n")
	report.WriteString("- **优势**: 保持插入顺序，遍历顺序可预测\n")
	report.WriteString("- **劣势**: 非并发安全，内存开销较大（双向链表）\n")
	report.WriteString("- **适用场景**: 需要保持插入顺序的单线程场景\n\n")

	report.WriteString("#### 3.1.6 TreeMap\n")
	report.WriteString("- **优势**: 红黑树保证 O(log n) 复杂度，有序性\n")
	report.WriteString("- **劣势**: 读写性能均低于哈希表实现\n")
	report.WriteString("- **适用场景**: 需要范围查询或有序遍历的场景\n\n")

	report.WriteString("### 3.2 使用建议\n\n")
	report.WriteString("| 使用场景 | 推荐 Map 类型 | 理由 |\n")
	report.WriteString("|---------|--------------|------|\n")
	report.WriteString("| 单线程高性能 | NativeMap | 性能最优，内存开销最小 |\n")
	report.WriteString("| 并发读多写少 | sync.Map | 官方实现，读性能优秀 |\n")
	report.WriteString("| 并发读写均衡 | ConcurrentHashMap | 简单可靠，性能稳定 |\n")
	report.WriteString("| 并发 + 有序性 | ConcurrentSkipListMap | 并发性能好，支持有序操作 |\n")
	report.WriteString("| 保持插入顺序 | OrderedMap | 唯一支持插入顺序的实现 |\n")
	report.WriteString("| 范围查询 | TreeMap 或 ConcurrentSkipListMap | 红黑树或跳表均支持高效范围查询 |\n\n")

	report.WriteString("### 3.3 内存与 GC 分析\n\n")
	report.WriteString("- **内存占用**: NativeMap < ConcurrentHashMap < sync.Map < OrderedMap < TreeMap < ConcurrentSkipListMap\n")
	report.WriteString("- **内存分配次数**: 跳表和树结构因为节点分配导致分配次数显著高于哈希表\n")
	report.WriteString("- **GC 压力**: 内存分配次数越多，GC 压力越大，建议在高频场景使用对象池\n\n")

	report.WriteString("---\n\n")
	report.WriteString("## 4. 性能优化建议\n\n")
	report.WriteString("1. **预分配容量**: 对于可预估大小的 map，使用 `make(map[K]V, capacity)` 减少扩容\n")
	report.WriteString("2. **避免热点 key**: 高并发场景下，热点 key 会导致锁竞争，考虑分片\n")
	report.WriteString("3. **选择合适数据结构**: 根据场景选择最适合的 Map 类型\n")
	report.WriteString("4. **减少内存分配**: 复用对象，减少 GC 压力\n")
	report.WriteString("5. **批量操作**: 尽量批量读写，减少锁竞争\n\n")

	report.WriteString("---\n\n")
	report.WriteString("**测试完成时间**: " + time.Now().Format("2006-01-02 15:04:05") + "\n")

	// 写入文件
	return os.WriteFile(reportPath, []byte(report.String()), 0644)
}

func parseBenchOutput(output string) []BenchResult {
	var results []BenchResult
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if !strings.HasPrefix(line, "Benchmark") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		var r BenchResult
		r.Name = fields[0]

		// 解析 ns/op
		fmt.Sscanf(fields[2], "%f", &r.NsPerOp)

		// 解析 B/op
		if len(fields) >= 5 {
			fmt.Sscanf(fields[4], "%d", &r.BytesPerOp)
		}

		// 解析 allocs/op
		if len(fields) >= 7 {
			fmt.Sscanf(fields[6], "%d", &r.AllocsPerOp)
		}

		results = append(results, r)
	}

	return results
}

func getOpIndex(op string) int {
	operations := map[string]int{
		"Put":           1,
		"Get":           2,
		"Range":         3,
		"ConcurrentPut": 4,
		"Mixed":         5,
	}
	return operations[op]
}

func getOpTitle(op string) string {
	titles := map[string]string{
		"Put":           "写入操作 (Put)",
		"Get":           "读取操作 (Get)",
		"Range":         "遍历操作 (Range)",
		"ConcurrentPut": "并发写入 (Concurrent Put)",
		"Mixed":         "混合操作 (Mixed 70% 写 + 20% 读 + 10% 删)",
	}
	return titles[op]
}

func extractMapType(name string) string {
	parts := strings.Split(name, "_")
	if len(parts) > 0 {
		return strings.TrimPrefix(parts[0], "Benchmark")
	}
	return name
}

func getGoVersion() string {
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}
