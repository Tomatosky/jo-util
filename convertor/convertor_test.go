package convertor

import (
	"testing"
)

func TestConversionFunctions(t *testing.T) {
	tests := []struct {
		name     string
		function string
		input    any
		want     any
	}{
		{"ToInt with string", "ToInt", "123", 123},
		{"ToInt with int", "ToInt", 456, 456},
		{"ToInt with float", "ToInt", 789.0, 789},
		{"ToInt32 with string", "ToInt32", "123", int32(123)},
		{"ToInt32 with int", "ToInt32", 456, int32(456)},
		{"ToInt32 with float", "ToInt32", 789.0, int32(789)},
		{"ToInt64 with string", "ToInt64", "123", int64(123)},
		{"ToInt64 with int", "ToInt64", 456, int64(456)},
		{"ToInt64 with float", "ToInt64", 789.0, int64(789)},
		{"ToFloat32 with string", "ToFloat32", "123.5", float32(123.5)},
		{"ToFloat32 with int", "ToFloat32", 456, float32(456)},
		{"ToFloat32 with float", "ToFloat32", 789.5, float32(789.5)},
		{"ToFloat64 with string", "ToFloat64", "123.5", 123.5},
		{"ToFloat64 with int", "ToFloat64", 456, float64(456)},
		{"ToFloat64 with float", "ToFloat64", 789.5, 789.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%s panicked with %v", tt.function, r)
				}
			}()

			var got any
			switch tt.function {
			case "ToInt":
				got = ToInt(tt.input)
			case "ToInt32":
				got = ToInt32(tt.input)
			case "ToInt64":
				got = ToInt64(tt.input)
			case "ToFloat32":
				got = ToFloat32(tt.input)
			case "ToFloat64":
				got = ToFloat64(tt.input)
			default:
				t.Errorf("Unknown function: %s", tt.function)
				return
			}

			if got != tt.want {
				t.Errorf("%s(%v) = %v (%T), want %v (%T)",
					tt.function, tt.input, got, got, tt.want, tt.want)
			}
		})
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{"nil", nil, ""},
		{"float32", float32(3.14), "3.14"},
		{"float64", 3.141592653589793, "3.141592653589793"},
		{"int", 42, "42"},
		{"int8", int8(8), "8"},
		{"int16", int16(16), "16"},
		{"int32", int32(32), "32"},
		{"int64", int64(64), "64"},
		{"uint", uint(42), "42"},
		{"uint8", uint8(8), "8"},
		{"uint16", uint16(16), "16"},
		{"uint32", uint32(32), "32"},
		{"uint64", uint64(64), "64"},
		{"string", "hello", "hello"},
		{"[]byte", []byte("world"), "world"},
		{"struct", struct{ Name string }{Name: "test"}, `{"Name":"test"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToString(tt.input)
			if got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
