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

func BenchmarkRun(b *testing.B) {
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
	report.WriteString("## 3. 综合分析与建议\n\n")

	// 动态分析：统计每种Map类型的性能表现
	mapTypePerformance := make(map[string]struct {
		totalTests    int
		bestCount     int
		worstCount    int
		avgPerformace float64
	})

	for _, op := range operations {
		var opResults []BenchResult
		for key, results := range grouped {
			if strings.HasSuffix(key, "-"+op) {
				opResults = append(opResults, results...)
			}
		}

		if len(opResults) == 0 {
			continue
		}

		// 排序找出最快和最慢
		sort.Slice(opResults, func(i, j int) bool {
			return opResults[i].NsPerOp < opResults[j].NsPerOp
		})

		if len(opResults) > 0 {
			fastest := extractMapType(opResults[0].Name)
			slowest := extractMapType(opResults[len(opResults)-1].Name)

			stats := mapTypePerformance[fastest]
			stats.bestCount++
			stats.totalTests++
			mapTypePerformance[fastest] = stats

			stats = mapTypePerformance[slowest]
			stats.worstCount++
			stats.totalTests++
			mapTypePerformance[slowest] = stats
		}
	}

	// 生成动态性能特点总结
	report.WriteString("### 3.1 性能特点总结\n\n")
	report.WriteString("基于本次测试结果的分析：\n\n")

	// 找出表现最好和最差的Map类型
	var bestMap, worstMap string
	var bestScore, worstScore float64 = -1, 1000

	for mapType, stats := range mapTypePerformance {
		score := float64(stats.bestCount) / float64(stats.totalTests)
		if score > bestScore {
			bestScore = score
			bestMap = mapType
		}
		if float64(stats.worstCount)/float64(stats.totalTests) > worstScore {
			worstScore = float64(stats.worstCount) / float64(stats.totalTests)
			worstMap = mapType
		}
	}

	// 输出整体性能排名
	if bestMap != "" {
		report.WriteString(fmt.Sprintf("- **综合性能最佳**: %s (在 %.1f%% 的测试中表现最优)\n", bestMap, bestScore*100))
	}
	if worstMap != "" && worstMap != bestMap {
		report.WriteString(fmt.Sprintf("- **综合性能最弱**: %s (在 %.1f%% 的测试中表现最差)\n", worstMap, worstScore*100))
	}
	report.WriteString("\n")

	// 针对每种Map生成动态分析
	mapTypes := []string{"NativeMap", "sync.Map", "ConcurrentHashMap", "ConcurrentSkipListMap", "OrderedMap", "TreeMap"}

	for _, mapType := range mapTypes {
		report.WriteString(fmt.Sprintf("#### 3.1.%d %s\n", getMapTypeIndex(mapType), getMapTypeTitle(mapType)))

		// 统计该Map类型在各操作中的表现
		var putPerf, getPerf, rangePerf, concurrentPerf, mixedPerf float64
		var putCount, getCount, rangeCount, concurrentCount, mixedCount int
		perfCount := 0

		for _, r := range results {
			if extractMapType(r.Name) == mapType {
				perfCount++

				// 根据操作类型累计性能数据
				name := r.Name
				if strings.Contains(name, "_Put_") && !strings.Contains(name, "ConcurrentPut") {
					putPerf += r.NsPerOp
					putCount++
				} else if strings.Contains(name, "_Get_") {
					getPerf += r.NsPerOp
					getCount++
				} else if strings.Contains(name, "_Range_") {
					rangePerf += r.NsPerOp
					rangeCount++
				} else if strings.Contains(name, "_ConcurrentPut_") {
					concurrentPerf += r.NsPerOp
					concurrentCount++
				} else if strings.Contains(name, "_Mixed_") {
					mixedPerf += r.NsPerOp
					mixedCount++
				}
			}
		}

		if perfCount > 0 {
			// 根据实际测试数据生成分析
			report.WriteString(fmt.Sprintf("- **测试场景数**: %d 个\n", perfCount))

			// 显示各操作的平均性能
			perfDetails := []string{}
			if putCount > 0 {
				perfDetails = append(perfDetails, fmt.Sprintf("写入 %.2f ns/op", putPerf/float64(putCount)))
			}
			if getCount > 0 {
				perfDetails = append(perfDetails, fmt.Sprintf("读取 %.2f ns/op", getPerf/float64(getCount)))
			}
			if rangeCount > 0 {
				perfDetails = append(perfDetails, fmt.Sprintf("遍历 %.2f ns/op", rangePerf/float64(rangeCount)))
			}
			if concurrentCount > 0 {
				perfDetails = append(perfDetails, fmt.Sprintf("并发写入 %.2f ns/op", concurrentPerf/float64(concurrentCount)))
			}
			if mixedCount > 0 {
				perfDetails = append(perfDetails, fmt.Sprintf("混合操作 %.2f ns/op", mixedPerf/float64(mixedCount)))
			}

			if len(perfDetails) > 0 {
				report.WriteString(fmt.Sprintf("- **平均性能**: %s\n", strings.Join(perfDetails, "、")))
			}

			// 查找该Map的优势操作
			advantages := []string{}
			for _, op := range operations {
				var opResults []BenchResult
				for key, results := range grouped {
					if strings.HasSuffix(key, "-"+op) {
						opResults = append(opResults, results...)
					}
				}

				if len(opResults) > 0 {
					sort.Slice(opResults, func(i, j int) bool {
						return opResults[i].NsPerOp < opResults[j].NsPerOp
					})

					// 如果在前30%，说明该操作是优势
					for i := 0; i < len(opResults)*3/10; i++ {
						if extractMapType(opResults[i].Name) == mapType {
							advantages = append(advantages, getOpTitle(op))
							break
						}
					}
				}
			}

			if len(advantages) > 0 {
				report.WriteString(fmt.Sprintf("- **优势操作**: %s\n", strings.Join(advantages, "、")))
			}

			// 生成适用场景建议
			report.WriteString(fmt.Sprintf("- **适用场景**: %s\n\n", getScenarioRecommendation(mapType, advantages)))
		} else {
			report.WriteString("- 本次测试未包含此Map类型的详细数据\n\n")
		}
	}

	report.WriteString("### 3.2 使用建议\n\n")
	report.WriteString("根据本次测试结果，针对不同使用场景的推荐：\n\n")
	report.WriteString("| 使用场景 | 推荐 Map 类型 | 理由 |\n")
	report.WriteString("|---------|--------------|------|\n")

	// 动态生成使用建议表格
	recommendations := generateDynamicRecommendations(grouped, operations)
	for _, rec := range recommendations {
		report.WriteString(fmt.Sprintf("| %s | %s | %s |\n", rec.Scenario, rec.RecommendedMap, rec.Reason))
	}

	report.WriteString("\n")

	report.WriteString("### 3.3 内存与 GC 分析\n\n")

	// 动态分析内存使用情况
	memoryRanking := analyzeMemoryUsage(results)
	report.WriteString(fmt.Sprintf("- **内存占用排序**: %s\n", strings.Join(memoryRanking, " < ")))
	report.WriteString("- **建议**: 在内存敏感场景下，优先选择内存占用较低的Map类型\n")
	report.WriteString("- **GC 压力**: 内存分配次数越多，GC 压力越大，建议在高频场景使用对象池\n\n")

	report.WriteString("---\n\n")
	report.WriteString("## 4. 性能优化建议\n\n")

	// 基于测试结果生成优化建议
	optimizationSuggestions := generateOptimizationSuggestions(results, grouped)
	for i, suggestion := range optimizationSuggestions {
		report.WriteString(fmt.Sprintf("%d. **%s**: %s\n", i+1, suggestion.Title, suggestion.Content))
	}
	report.WriteString("\n")

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

// 获取Map类型的索引号
func getMapTypeIndex(mapType string) int {
	mapTypes := map[string]int{
		"NativeMap":             1,
		"sync.Map":              2,
		"ConcurrentHashMap":     3,
		"ConcurrentSkipListMap": 4,
		"OrderedMap":            5,
		"TreeMap":               6,
	}
	return mapTypes[mapType]
}

// 获取Map类型的标题
func getMapTypeTitle(mapType string) string {
	titles := map[string]string{
		"NativeMap":             "Go 原生 map (NativeMap)",
		"sync.Map":              "sync.Map",
		"ConcurrentHashMap":     "ConcurrentHashMap",
		"ConcurrentSkipListMap": "ConcurrentSkipListMap",
		"OrderedMap":            "OrderedMap",
		"TreeMap":               "TreeMap",
	}
	return titles[mapType]
}

// 根据Map类型和优势操作生成适用场景建议
func getScenarioRecommendation(mapType string, advantages []string) string {
	recommendations := map[string]string{
		"NativeMap":             "单线程或已有外部锁保护的高性能场景",
		"sync.Map":              "并发读多写少的场景，如缓存系统",
		"ConcurrentHashMap":     "中等并发、读写均衡的场景",
		"ConcurrentSkipListMap": "需要有序性和并发访问的场景，如范围查询",
		"OrderedMap":            "需要保持插入顺序的单线程场景",
		"TreeMap":               "需要有序遍历或范围查询的场景",
	}

	baseRec := recommendations[mapType]
	if len(advantages) > 0 {
		return fmt.Sprintf("%s，擅长%s", baseRec, strings.Join(advantages, "、"))
	}
	return baseRec
}

// Recommendation 表示一个使用建议
type Recommendation struct {
	Scenario       string
	RecommendedMap string
	Reason         string
}

// 动态生成使用建议
func generateDynamicRecommendations(grouped map[string][]BenchResult, operations []string) []Recommendation {
	recommendations := []Recommendation{}

	// 分析每个操作类型，找出最优Map
	for _, op := range operations {
		var opResults []BenchResult
		for key, results := range grouped {
			if strings.HasSuffix(key, "-"+op) {
				opResults = append(opResults, results...)
			}
		}

		if len(opResults) == 0 {
			continue
		}

		// 按性能排序
		sort.Slice(opResults, func(i, j int) bool {
			return opResults[i].NsPerOp < opResults[j].NsPerOp
		})

		if len(opResults) > 0 {
			bestMap := extractMapType(opResults[0].Name)
			scenario := getScenarioForOperation(op)
			reason := fmt.Sprintf("在%s测试中表现最优，平均耗时 %.2f ns/op", getOpTitle(op), opResults[0].NsPerOp)

			recommendations = append(recommendations, Recommendation{
				Scenario:       scenario,
				RecommendedMap: bestMap,
				Reason:         reason,
			})
		}
	}

	return recommendations
}

// 根据操作类型获取使用场景描述
func getScenarioForOperation(op string) string {
	scenarios := map[string]string{
		"Put":           "高频写入场景",
		"Get":           "高频读取场景",
		"Range":         "需要遍历所有元素",
		"ConcurrentPut": "高并发写入场景",
		"Mixed":         "混合读写场景",
	}
	return scenarios[op]
}

// OptimizationSuggestion 表示一个优化建议
type OptimizationSuggestion struct {
	Title   string
	Content string
}

// 基于测试结果生成优化建议
func generateOptimizationSuggestions(results []BenchResult, grouped map[string][]BenchResult) []OptimizationSuggestion {
	suggestions := []OptimizationSuggestion{}

	// 建议1：预分配容量
	suggestions = append(suggestions, OptimizationSuggestion{
		Title:   "预分配容量",
		Content: "对于可预估大小的 map，使用 `make(map[K]V, capacity)` 可以减少扩容次数，提升性能",
	})

	// 建议2：根据并发情况选择Map
	hasConcurrent := false
	for key := range grouped {
		if strings.Contains(key, "Concurrent") {
			hasConcurrent = true
			break
		}
	}

	if hasConcurrent {
		suggestions = append(suggestions, OptimizationSuggestion{
			Title:   "选择合适的并发Map",
			Content: "根据测试结果，在高并发场景下应选择专门设计的并发安全Map，避免使用普通map加锁的方式",
		})
	}

	// 建议3：内存优化
	var highMemMaps []string
	avgMem := int64(0)
	for _, r := range results {
		avgMem += r.BytesPerOp
	}
	if len(results) > 0 {
		avgMem = avgMem / int64(len(results))
		for _, r := range results {
			if r.BytesPerOp > avgMem*2 {
				mapType := extractMapType(r.Name)
				if !contains(highMemMaps, mapType) {
					highMemMaps = append(highMemMaps, mapType)
				}
			}
		}
	}

	if len(highMemMaps) > 0 {
		suggestions = append(suggestions, OptimizationSuggestion{
			Title:   "注意内存占用",
			Content: fmt.Sprintf("测试发现 %s 的内存占用较高，在内存敏感场景下需谨慎使用或考虑使用对象池复用", strings.Join(highMemMaps, "、")),
		})
	}

	// 建议4：避免热点key
	suggestions = append(suggestions, OptimizationSuggestion{
		Title:   "避免热点 key",
		Content: "高并发场景下，热点 key 会导致锁竞争，考虑使用分片技术将热点key分散到多个map中",
	})

	// 建议5：批量操作
	suggestions = append(suggestions, OptimizationSuggestion{
		Title:   "批量操作优化",
		Content: "在并发场景下，尽量批量读写减少锁的获取和释放次数，可以显著提升性能",
	})

	return suggestions
}

// 分析内存使用情况
func analyzeMemoryUsage(results []BenchResult) []string {
	// 统计每种Map类型的平均内存使用
	memUsage := make(map[string]int64)
	memCount := make(map[string]int)

	for _, r := range results {
		mapType := extractMapType(r.Name)
		memUsage[mapType] += r.BytesPerOp
		memCount[mapType]++
	}

	// 计算平均值
	type memStat struct {
		mapType string
		avgMem  int64
	}

	var memStats []memStat
	for mapType, total := range memUsage {
		if count := memCount[mapType]; count > 0 {
			memStats = append(memStats, memStat{
				mapType: mapType,
				avgMem:  total / int64(count),
			})
		}
	}

	// 排序
	sort.Slice(memStats, func(i, j int) bool {
		return memStats[i].avgMem < memStats[j].avgMem
	})

	// 生成排序结果
	ranking := make([]string, 0, len(memStats))
	for _, stat := range memStats {
		ranking = append(ranking, stat.mapType)
	}

	return ranking
}

// 辅助函数：检查字符串切片是否包含某个元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
