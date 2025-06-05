package listUtil

import (
	"encoding/json"
	"math/rand"
	"time"
)

func Contain[T comparable](slice []T, target T) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}

	return false
}

// Unique 定义一个泛型函数，返回切片中不重复的元素
func Unique[T comparable](slice []T) []T {
	// 使用 map 来记录已经出现过的元素
	seen := make(map[T]bool)
	result := make([]T, 0)

	// 遍历切片，将不重复的元素添加到结果中
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

func ToString[T comparable](slice []T) string {
	marshal, err := json.Marshal(slice)
	if err != nil {
		panic(err)
	}
	return string(marshal)
}

func Reverse[T any](slice []T) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func Shuffle[T any](slice []T) []T {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})

	return slice
}

func AddIfAbsent[T comparable](slice *[]T, item T) {
	for _, v := range *slice {
		if v == item {
			return
		}
	}
	*slice = append(*slice, item)
}

func Remove[T comparable](slice []T, target T, all bool) []T {
	if !all {
		// 只删除第一个匹配项
		for i, item := range slice {
			if item == target {
				return append(slice[:i], slice[i+1:]...)
			}
		}
	} else {
		// 删除所有匹配项
		result := make([]T, 0, len(slice))
		for _, item := range slice {
			if item != target {
				result = append(result, item)
			}
		}
		return result
	}
	return slice
}
