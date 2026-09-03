package randomUtil

import (
	"fmt"
	"math"
	"math/rand/v2"

	"github.com/Tomatosky/jo-util/logger"
	"github.com/Tomatosky/jo-util/numberUtil"
)

// RandomInt 生成指定范围 [start, end) 的随机整数
func RandomInt[T numberUtil.Number](start, end T) T {
	if start >= end {
		logger.Log.Error(fmt.Sprintf("%v", "invalid range: start >= end"))
		panic("invalid range: start >= end")
	}
	rangeSize := uint64(end) - uint64(start)
	return start + T(rand.Uint64N(rangeSize))
}

// RandomEle 从切片中随机选择一个元素
func RandomEle[T any](slice []T) T {
	if len(slice) == 0 {
		logger.Log.Error(fmt.Sprintf("%v", "slice is empty"))
		panic("slice is empty")
	}
	return slice[rand.IntN(len(slice))]
}

// RandomEleSet 从切片中随机选择 n 个不重复的元素
func RandomEleSet[T any](slice []T, n int) []T {
	if n <= 0 {
		return nil
	}
	length := len(slice)
	if length == 0 {
		logger.Log.Error(fmt.Sprintf("%v", "slice is empty"))
		panic("slice is empty")
	}
	if n > length {
		n = length
	}
	indices := rand.Perm(length)
	result := make([]T, n)
	for i := 0; i < n; i++ {
		result[i] = slice[indices[i]]
	}
	return result
}

// RandomWeightedKey 根据权重随机选择一个键
func RandomWeightedKey[K comparable, V numberUtil.Number](weights map[K]V) K {
	// 计算总权重
	var sum uint64
	for _, w := range weights {
		if w < 0 {
			const message = "权重值不能为负数"
			logger.Log.Error(message)
			panic(message)
		}
		weight := uint64(w)
		if math.MaxUint64-sum < weight {
			const message = "权重值总和溢出"
			logger.Log.Error(message)
			panic(message)
		}
		sum += weight
	}

	// 处理无效权重的情况
	if sum == 0 {
		logger.Log.Error(fmt.Sprintf("%v", "所有权重值总和不能为0"))
		panic("所有权重值总和不能为0")
	}

	// 生成随机数
	r := rand.Uint64N(sum)

	// 查找对应的键
	var runningTotal uint64
	for key, weight := range weights {
		runningTotal += uint64(weight)
		if runningTotal > r {
			return key
		}
	}

	// 理论上不会执行到这里（因为sum > 0）
	logger.Log.Error(fmt.Sprintf("%v", "未找到有效键"))
	var zero K
	return zero
}

// RandomString 生成包含数字和字母的随机字符串
func RandomString(length int) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return string(b)
}

// RandomNumbers 生成只包含数字的随机字符串
func RandomNumbers(length int) string {
	const charset = "0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return string(b)
}
