package convertor

import (
	"bytes"
	"encoding/binary"
	"math"
	"strconv"
	"testing"
)

func requirePanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Error("function did not panic")
		}
	}()
	f()
}

func TestToBool(t *testing.T) {
	tests := []struct {
		input any
		want  bool
	}{
		{input: true, want: true},
		{input: false, want: false},
		{input: "TRUE", want: true},
		{input: "1", want: true},
		{input: "0", want: false},
	}
	for _, tt := range tests {
		if got := ToBool(tt.input); got != tt.want {
			t.Errorf("ToBool(%v) = %v, want %v", tt.input, got, tt.want)
		}
	}
	requirePanic(t, func() { ToBool("not-a-bool") })
}

func TestNumericConversions(t *testing.T) {
	tests := []struct {
		name string
		got  func() any
		want any
	}{
		{name: "int from string", got: func() any { return ToInt("123.9") }, want: 123},
		{name: "int from negative float", got: func() any { return ToInt(-12.9) }, want: -12},
		{name: "int32", got: func() any { return ToInt32("123") }, want: int32(123)},
		{name: "int32 from decimal", got: func() any { return ToInt32("-12.9") }, want: int32(-12)},
		{name: "int64", got: func() any { return ToInt64(int32(-456)) }, want: int64(-456)},
		{name: "int64 from decimal", got: func() any { return ToInt64("123.9") }, want: int64(123)},
		{name: "float32", got: func() any { return ToFloat32("123.5") }, want: float32(123.5)},
		{name: "float64", got: func() any { return ToFloat64(-789.25) }, want: -789.25},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.got(); got != tt.want {
				t.Errorf("conversion = %v (%T), want %v (%T)", got, got, tt.want, tt.want)
			}
		})
	}

	for name, f := range map[string]func(){
		"ToInt":     func() { ToInt("invalid") },
		"ToInt32":   func() { ToInt32("invalid") },
		"ToInt64":   func() { ToInt64("invalid") },
		"ToFloat32": func() { ToFloat32("invalid") },
		"ToFloat64": func() { ToFloat64("invalid") },
	} {
		t.Run(name+" invalid input", func(t *testing.T) {
			requirePanic(t, f)
		})
	}
}

func TestIntegerConversionsPreserveLargeValues(t *testing.T) {
	if got := ToInt64(int64(math.MaxInt64)); got != math.MaxInt64 {
		t.Errorf("ToInt64(MaxInt64) = %d, want %d", got, int64(math.MaxInt64))
	}
	if got := ToInt64("-9223372036854775808"); got != math.MinInt64 {
		t.Errorf("ToInt64(MinInt64 string) = %d, want %d", got, int64(math.MinInt64))
	}
	if got := ToInt64(uint64(math.MaxInt64)); got != math.MaxInt64 {
		t.Errorf("ToInt64(uint64(MaxInt64)) = %d, want %d", got, int64(math.MaxInt64))
	}
	if strconv.IntSize == 64 {
		if got := ToInt(int64(math.MaxInt64)); int64(got) != math.MaxInt64 {
			t.Errorf("ToInt(MaxInt64) = %d, want %d", got, int64(math.MaxInt64))
		}
	}
}

func TestIntegerConversionsRejectOverflowAndNonFiniteValues(t *testing.T) {
	tests := map[string]func(){
		"int32 overflow":         func() { ToInt32("2147483648") },
		"int32 decimal overflow": func() { ToInt32("2147483648.0") },
		"int64 overflow":         func() { ToInt64("9223372036854775808") },
		"int64 decimal overflow": func() { ToInt64("9223372036854775808.0") },
		"uint64 overflow":        func() { ToInt64(uint64(math.MaxUint64)) },
		"NaN":                    func() { ToInt64(math.NaN()) },
		"positive Inf":           func() { ToInt64(math.Inf(1)) },
	}
	for name, f := range tests {
		t.Run(name, func(t *testing.T) {
			requirePanic(t, f)
		})
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{name: "nil", input: nil, want: ""},
		{name: "float32", input: float32(3.14), want: "3.14"},
		{name: "float64", input: 3.141592653589793, want: "3.141592653589793"},
		{name: "int", input: -42, want: "-42"},
		{name: "int8", input: int8(-8), want: "-8"},
		{name: "int16", input: int16(-16), want: "-16"},
		{name: "int32", input: int32(-32), want: "-32"},
		{name: "int64", input: int64(-64), want: "-64"},
		{name: "uint", input: uint(42), want: "42"},
		{name: "uint8", input: uint8(8), want: "8"},
		{name: "uint16", input: uint16(16), want: "16"},
		{name: "uint32", input: uint32(32), want: "32"},
		{name: "uint64", input: uint64(64), want: "64"},
		{name: "string", input: "hello", want: "hello"},
		{name: "bytes", input: []byte("world"), want: "world"},
		{name: "bool JSON fallback", input: true, want: "true"},
		{name: "struct JSON fallback", input: struct{ Name string }{Name: "test"}, want: `{"Name":"test"}`},
		{name: "unsupported JSON value", input: make(chan int), want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToString(tt.input); got != tt.want {
				t.Errorf("ToString(%v) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToBytes(t *testing.T) {
	intBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(intBytes, uint64(42))
	negativeIntBytes := make([]byte, 8)
	negative := int64(-2)
	binary.BigEndian.PutUint64(negativeIntBytes, uint64(negative))
	uintBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(uintBytes, 42)
	float32Bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(float32Bytes, math.Float32bits(1.5))
	float64Bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(float64Bytes, math.Float64bits(-2.25))

	tests := []struct {
		name  string
		input any
		want  []byte
	}{
		{name: "int", input: 42, want: intBytes},
		{name: "int8 encoded as int64", input: int8(-2), want: negativeIntBytes},
		{name: "uint16 encoded as uint64", input: uint16(42), want: uintBytes},
		{name: "float32", input: float32(1.5), want: float32Bytes},
		{name: "float64", input: -2.25, want: float64Bytes},
		{name: "bool", input: true, want: []byte("true")},
		{name: "string", input: "hello", want: []byte("hello")},
		{name: "bytes", input: []byte{0, 1, 255}, want: []byte{0, 1, 255}},
		{name: "struct JSON fallback", input: struct{ N int }{N: 3}, want: []byte(`{"N":3}`)},
		{name: "nil JSON fallback", input: nil, want: []byte("null")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToBytes(tt.input)
			if err != nil {
				t.Fatalf("ToBytes(%v): %v", tt.input, err)
			}
			if !bytes.Equal(got, tt.want) {
				t.Errorf("ToBytes(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}

	got, err := ToBytes(func() {})
	if err == nil || got != nil {
		t.Errorf("ToBytes(unsupported) = (%v, %v), want (nil, error)", got, err)
	}
}
