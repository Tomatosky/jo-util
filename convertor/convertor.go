package convertor

import (
	"encoding/json"
	"strconv"
)

// ToInt 转换为 int
func ToInt(value any) int {
	s := ToString(value)
	atoi, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return atoi
}

// ToInt32 转换为 int32
func ToInt32(value any) int32 {
	s := ToString(value)
	i32, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		panic(err)
	}
	return int32(i32)
}

// ToInt64 转换为 int64
func ToInt64(value any) int64 {
	s := ToString(value)
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return i64
}

// ToFloat32 转换为 float32
func ToFloat32(value any) float32 {
	s := ToString(value)
	f32, err := strconv.ParseFloat(s, 32)
	if err != nil {
		panic(err)
	}
	return float32(f32)
}

// ToFloat64 转换为 float64
func ToFloat64(value any) float64 {
	s := ToString(value)
	f64, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f64
}

func ToString(value any) string {
	if value == nil {
		return ""
	}

	switch val := value.(type) {
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case int:
		return strconv.FormatInt(int64(val), 10)
	case int8:
		return strconv.FormatInt(int64(val), 10)
	case int16:
		return strconv.FormatInt(int64(val), 10)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case int64:
		return strconv.FormatInt(val, 10)
	case uint:
		return strconv.FormatUint(uint64(val), 10)
	case uint8:
		return strconv.FormatUint(uint64(val), 10)
	case uint16:
		return strconv.FormatUint(uint64(val), 10)
	case uint32:
		return strconv.FormatUint(uint64(val), 10)
	case uint64:
		return strconv.FormatUint(val, 10)
	case string:
		return val
	case []byte:
		return string(val)
	default:
		b, err := json.Marshal(val)
		if err != nil {
			return ""
		}
		return string(b)
	}
}
