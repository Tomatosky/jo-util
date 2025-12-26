package mapUtil

import (
	"context"
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
	// 使用 t.Skip() 来阻止正常 go test 时运行
	// 如果要运行此测试,需要显式运行: go test -run TestRun
	if testing.Short() {
		t.Skip("跳过长时间运行的性能测试,使用 go test -run TestRun 来运行")
	}

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

	// 记录测试开始时间
	testStartTime := time.Now()
	fmt.Printf("测试开始时间: %s\n", testStartTime.Format("2006-01-02 15:04:05"))
	fmt.Println()

	// 运行基准测试
	fmt.Println("正在运行基准测试，这可能需要几分钟...")
	fmt.Println()

	// 创建可取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "test", "-bench=.", "-benchmem", "-benchtime=3s", "-timeout=60m")
	cmd.Dir = dir

	// 设置实时输出
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 启动命令
	err = cmd.Start()
	if err != nil {
		fmt.Printf("启动测试失败: %v\n", err)
		return
	}

	// 等待命令完成
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("运行测试失败: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Println("基准测试完成！")
	fmt.Println()

	// 生成报告 - 重新运行一次基准测试以捕获输出
	fmt.Println("正在生成性能分析报告...")
	fmt.Println("重新捕获测试输出以生成报告...")

	cmdReport := exec.Command("go", "test", "-bench=.", "-benchmem", "-benchtime=3s", "-timeout=60m")
	cmdReport.Dir = dir
	output, err := cmdReport.CombinedOutput()
	if err != nil {
		fmt.Printf("重新运行测试以生成报告失败: %v\n", err)
		// 尝试使用空输出生成基础报告
		output = []byte("")
	}

	reportPath := filepath.Join(dir, "Map性能测试分析报告.md")
	err = generateReport(reportPath, string(output), testStartTime)
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

func generateReport(reportPath, benchOutput string, testStartTime time.Time) error {
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
	report.WriteString(fmt.Sprintf("**测试开始时间**: %s\n\n", testStartTime.Format("2006-01-02 15:04:05")))
	report.WriteString(fmt.Sprintf("**报告生成时间**: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

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
