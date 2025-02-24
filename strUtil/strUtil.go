package strUtil

import (
	"strconv"
)

// ToInt 将字符串转换为 int
func ToInt(s string) int {
	atoi, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return atoi
}

// ToInt32 将字符串转换为 int32
func ToInt32(s string) int32 {
	i64, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		panic(err)
	}
	return int32(i64)
}

// ToInt64 将字符串转换为 int64
func ToInt64(s string) int64 {
	i64, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		panic(err)
	}
	return i64
}

// ToFloat32 将字符串转换为 float32
func ToFloat32(s string) float32 {
	f64, err := strconv.ParseFloat(s, 32)
	if err != nil {
		panic(err)
	}
	return float32(f64)
}

// ToFloat64 将字符串转换为 float64
func ToFloat64(s string) float64 {
	f64, err := strconv.ParseFloat(s, 32)
	if err != nil {
		panic(err)
	}
	return f64
}
