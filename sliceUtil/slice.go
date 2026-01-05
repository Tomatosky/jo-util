package sliceUtil

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Tomatosky/jo-util/convertor"
	"github.com/Tomatosky/jo-util/logger"
	"github.com/Tomatosky/jo-util/numberUtil"
)

// Contain 切片是否包含某个元素
func Contain[T comparable](slice []T, target T) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}

	return false
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
		logger.Log.Error(fmt.Sprintf("%v", err))
		panic(err)
	}
	return string(marshal)
}

// Reverse 反转切片
func Reverse[T any](slice []T) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// Shuffle 打乱切片
func Shuffle[T any](slice []T) []T {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})

	return slice
}

// AddIfAbsent 要是切片中不存在元素，则添加
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

// ToMap 函数用于将结构体切片转换为 map。
// 它接收一个泛型切片和一个用于从切片元素中提取键的函数，返回一个以提取的键为键、切片元素为值的 map。
// 参数 K 是 map 的键类型，需满足 comparable 约束，即键必须是可比较的类型。
// 参数 T 是切片元素的类型，可以是任意类型。
// 参数 slice 是要转换的结构体切片。
// 参数 getKey 是一个函数，用于从切片的每个元素中提取键。
// 返回值是一个以 K 为键、T 为值的 map。
func ToMap[K comparable, T any](slice []T, getKey func(T) K) map[K]T {
	result := make(map[K]T, len(slice))
	for _, item := range slice {
		key := getKey(item)
		result[key] = item
	}
	return result
}

// FieldExtractor 用于从结构体切片中提取特定字段，返回一个包含提取字段值的切片。
// 这是一个泛型函数，支持任意类型的结构体切片和提取的字段类型。
// 参数 T 表示结构体切片中元素的类型。
// 参数 F 表示要提取的字段的类型。
// 参数 slice 是包含结构体元素的切片。
// 参数 getField 是一个函数，用于从结构体元素中提取指定字段的值。
// 返回值为一个包含提取字段值的切片。
func FieldExtractor[T any, F any](slice []T, getField func(T) F) []F {
	result := make([]F, len(slice))
	for i, item := range slice {
		result[i] = getField(item)
	}
	return result
}

// Fill 用给定的值填充切片中的所有元素，并返回填充后的切片。
// 该函数会修改传入的切片本身。
// 参数 in 是要填充的切片，支持任意类型。
// 参数 fillValue 是用于填充切片元素的值，类型与切片元素类型一致。
// 返回值为填充后的切片，实际上与传入的切片是同一个实例。
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

func Max[T numberUtil.Number](nums []T) T {
	if len(nums) < 1 {
		logger.Log.Error(fmt.Sprintf("%v", "mathutil.Max: empty list"))
		panic("mathutil.Max: empty list")
	}
	max2 := nums[0]
	for _, v := range nums {
		if max2 < v {
			max2 = v
		}
	}
	return max2
}

func Min[T numberUtil.Number](nums []T) T {
	if len(nums) < 1 {
		logger.Log.Error(fmt.Sprintf("%v", "mathutil.min: empty list"))
		panic("mathutil.min: empty list")
	}
	min2 := nums[0]
	for _, v := range nums {
		if min2 > v {
			min2 = v
		}
	}
	return min2
}

func Sum[T numberUtil.Number](nums ...T) T {
	var sum T
	for _, v := range nums {
		sum += v
	}
	return sum
}

// Join 将切片中的元素连接成一个字符串
// 例: Join([]int{1, 2, 3}, ",") // "1,2,3"
func Join[T any](slice []T, sep string) string {
	strs := make([]string, len(slice))
	for i, v := range slice {
		strs[i] = convertor.ToString(v)
	}
	return strings.Join(strs, sep)
}
