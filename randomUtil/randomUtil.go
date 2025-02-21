package randomUtil

import (
	"math/rand"
)

// RandomInt 生成指定范围 [start, end) 的随机整数
func RandomInt(start, end int) int {
	if start >= end {
		panic("invalid range: start >= end")
	}
	return rand.Intn(end-start) + start
}

// RandomBytes 生成指定长度的随机字节切片
func RandomBytes(length int) []byte {
	if length < 0 {
		panic("length cannot be negative")
	}
	b := make([]byte, length)
	_, _ = rand.Read(b) // 忽略错误，因为 math/rand 的 Read 不会返回错误
	return b
}

// RandomEle 从切片中随机选择一个元素
func RandomEle[T any](slice []T) T {
	if len(slice) == 0 {
		panic("slice is empty")
	}
	return slice[rand.Intn(len(slice))]
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
	indices := rand.Perm(length)
	result := make([]T, n)
	for i := 0; i < n; i++ {
		result[i] = slice[indices[i]]
	}
	return result
}

// RandomString 生成包含数字和字母的随机字符串
func RandomString(length int) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// RandomNumbers 生成只包含数字的随机字符串
func RandomNumbers(length int) string {
	const charset = "0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
