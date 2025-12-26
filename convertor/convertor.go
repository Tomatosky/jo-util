package convertor

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/Tomatosky/jo-util/logger"
)

func ToBool(value any) bool {
	s := ToString(value)
	parseBool, err := strconv.ParseBool(s)
	if err != nil {
		logger.Log.Fatal(fmt.Sprintf("%v", err))
	}
	return parseBool
}

// ToInt 转换为 int
func ToInt(value any) int {
	s := ToString(value)
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		logger.Log.Fatal(fmt.Sprintf("%v", err))
	}
	return int(f)
}

// ToInt32 转换为 int32
func ToInt32(value any) int32 {
	s := ToString(value)
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		logger.Log.Fatal(fmt.Sprintf("%v", err))
	}
	return int32(f)
}

// ToInt64 转换为 int64
func ToInt64(value any) int64 {
	s := ToString(value)
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		logger.Log.Fatal(fmt.Sprintf("%v", err))
	}
	return int64(f)
}

// ToFloat32 转换为 float32
func ToFloat32(value any) float32 {
	s := ToString(value)
	f32, err := strconv.ParseFloat(s, 32)
	if err != nil {
		logger.Log.Fatal(fmt.Sprintf("%v", err))
	}
	return float32(f32)
}

// ToFloat64 转换为 float64
func ToFloat64(value any) float64 {
	s := ToString(value)
	f64, err := strconv.ParseFloat(s, 64)
	if err != nil {
		logger.Log.Fatal(fmt.Sprintf("%v", err))
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

func ToBytes(value any) ([]byte, error) {
	v := reflect.ValueOf(value)

	switch value.(type) {
	case int, int8, int16, int32, int64:
		number := v.Int()
		buf := bytes.NewBuffer([]byte{})
		buf.Reset()
		err := binary.Write(buf, binary.BigEndian, number)
		return buf.Bytes(), err
	case uint, uint8, uint16, uint32, uint64:
		number := v.Uint()
		buf := bytes.NewBuffer([]byte{})
		buf.Reset()
		err := binary.Write(buf, binary.BigEndian, number)
		return buf.Bytes(), err
	case float32:
		number := float32(v.Float())
		bits := math.Float32bits(number)
		bytes2 := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes2, bits)
		return bytes2, nil
	case float64:
		number := v.Float()
		bits := math.Float64bits(number)
		bytes2 := make([]byte, 8)
		binary.BigEndian.PutUint64(bytes2, bits)
		return bytes2, nil
	case bool:
		return strconv.AppendBool([]byte{}, v.Bool()), nil
	case string:
		return []byte(v.String()), nil
	case []byte:
		return v.Bytes(), nil
	default:
		newValue, err := json.Marshal(value)
		return newValue, err
	}
}
