package listUtil

import (
	"encoding/json"
	"fmt"
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

// GetByIndex 根据索引返回slice中的元素，支持负数索引
// 负数索引表示从末尾开始计数，例如-1表示最后一个元素
func GetByIndex[T any](slice []T, index int) (T, error) {
	var zero T
	if len(slice) == 0 {
		return zero, fmt.Errorf("slice is empty")
	}

	// 处理负数索引
	if index < 0 {
		index = len(slice) + index
	}

	if index < 0 || index >= len(slice) {
		return zero, fmt.Errorf("index out of range")
	}

	return slice[index], nil
}

// InsertByIndex 将元素插入指定的索引位置，支持负数索引
// 负数索引表示从末尾开始计数，例如-1表示最后一个元素的位置
func InsertByIndex[T any](slice []T, index int, value T) ([]T, error) {
	// 处理负数索引
	if index < 0 {
		index = len(slice) + index + 1 // +1 因为插入位置在负数索引之后
	}

	if index < 0 || index > len(slice) {
		return nil, fmt.Errorf("index out of range")
	}

	// 扩展slice
	slice = append(slice, value)
	// 将元素移动到正确位置
	copy(slice[index+1:], slice[index:])
	slice[index] = value

	return slice, nil
}

// IndexOf 返回元素在切片中的索引，如果不存在则返回-1
func IndexOf[T comparable](slice []T, target T) int {
	for i, v := range slice {
		if v == target {
			return i
		}
	}
	return -1
}

// ToMap Struct Slice 转 Map
func ToMap[K comparable, T any](slice []T, getKey func(T) K) map[K]T {
	result := make(map[K]T, len(slice))
	for _, item := range slice {
		key := getKey(item)
		result[key] = item
	}
	return result
}

// FieldExtractor 用于从结构体切片中提取特定字段
func FieldExtractor[T any, F any](slice []T, getField func(T) F) []F {
	result := make([]F, len(slice))
	for i, item := range slice {
		result[i] = getField(item)
	}
	return result
}

// ContainAll 检查集合中是否包含所有指定元素
func ContainAll[T comparable](in []T, elements ...T) bool {
	// 创建元素查找map
	elementMap := make(map[T]struct{})
	for _, v := range elements {
		elementMap[v] = struct{}{}
	}
	// 检查集合中的每个元素
	for _, item := range in {
		if _, exists := elementMap[item]; exists {
			delete(elementMap, item)
		}
	}
	// 如果所有元素都被找到，map应该为空
	return len(elementMap) == 0
}

// ContainOne 检查集合中是否包含指定元素中的一个
func ContainOne[T comparable](in []T, elements ...T) bool {
	// 创建一个map用于快速查找
	elementSet := make(map[T]struct{}, len(elements))
	for _, e := range elements {
		elementSet[e] = struct{}{}
	}

	// 检查切片中的每个元素
	for _, v := range in {
		if _, ok := elementSet[v]; ok {
			return true
		}
	}

	return false
}

// Fill 用给定的值填充切片
func Fill[T any](in []T, fillValue T) []T {
	for i := range in {
		in[i] = fillValue
	}
	return in
}

// Intersection 获取两个slice的交集
func Intersection[T comparable](a, b []T) []T {
	set := make(map[T]bool)
	for _, item := range a {
		set[item] = true
	}
	var result []T
	for _, item := range b {
		if set[item] {
			result = append(result, item)
			set[item] = false
		}
	}
	if result == nil {
		return make([]T, 0)
	}
	return result
}

// Union 获取两个slice的并集
func Union[T comparable](a, b []T) []T {
	set := make(map[T]bool)
	var result []T
	// 遍历第一个 slice，存入 map 并记录结果
	for _, item := range a {
		if !set[item] {
			set[item] = true
			result = append(result, item)
		}
	}
	// 遍历第二个 slice，只添加未重复的元素
	for _, item := range b {
		if !set[item] {
			set[item] = true
			result = append(result, item)
		}
	}
	return result
}
