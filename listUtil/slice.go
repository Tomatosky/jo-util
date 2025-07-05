package listUtil

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
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
func ToMap[K comparable, T any](fieldName string, slice []T) (map[K]T, error) {
	result := make(map[K]T)
	for _, item := range slice {
		// 使用反射获取字段值
		val := reflect.ValueOf(item)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		if val.Kind() != reflect.Struct {
			return nil, fmt.Errorf("expected struct type, got %v", val.Kind())
		}
		field := val.FieldByName(fieldName)
		if !field.IsValid() {
			return nil, fmt.Errorf("field %s not found in struct", fieldName)
		}
		// 将字段值转换为 K 类型
		key, ok := field.Interface().(K)
		if !ok {
			return nil, fmt.Errorf("field %s is not of type %T", fieldName, *new(K))
		}
		result[key] = item
	}
	return result, nil
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
