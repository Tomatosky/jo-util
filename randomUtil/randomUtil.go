package randomUtil

import (
	"math/rand"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// Number 约束，限制为所有整数类型
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// RandomInt 生成指定范围 [start, end) 的随机整数
func RandomInt[T Number](start, end T) T {
	if start >= end {
		panic("invalid range: start >= end")
	}
	return T(rng.Int63n(int64(end-start))) + start
}

// RandomEle 从切片中随机选择一个元素
func RandomEle[T any](slice []T) T {
	if len(slice) == 0 {
		panic("slice is empty")
	}
	return slice[rng.Intn(len(slice))]
}

// RandomEleSet 从切片中随机选择 n 个不重复的元素
func RandomEleSet[T any](slice []T, n int) []T {
	if n <= 0 {
		return nil
	}
	length := len(slice)
	if length == 0 {
		panic("slice is empty")
	}
	if n > length {
		n = length
	}
	indices := rng.Perm(length)
	result := make([]T, n)
	for i := 0; i < n; i++ {
		result[i] = slice[indices[i]]
	}
	return result
}

// RandomWeightedKey 根据权重随机选择一个键
func RandomWeightedKey[K comparable, V Number](weights map[K]V) K {
	// 计算总权重
	var sum int
	for _, w := range weights {
		sum += int(w)
	}

	// 处理无效权重的情况
	if sum == 0 {
		panic("所有权重值总和不能为0")
	}

	// 生成随机数
	r := rng.Intn(sum)

	// 查找对应的键
	var runningTotal int
	for key, weight := range weights {
		runningTotal += int(weight)
		if runningTotal > r {
			return key
		}
	}

	// 理论上不会执行到这里（因为sum > 0）
	panic("未找到有效键")
}

// RandomString 生成包含数字和字母的随机字符串
func RandomString(length int) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rng.Intn(len(charset))]
	}
	return string(b)
}

// RandomNumbers 生成只包含数字的随机字符串
func RandomNumbers(length int) string {
	const charset = "0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rng.Intn(len(charset))]
	}
	return string(b)
}
