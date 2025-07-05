package strUtil

import (
	"testing"
)

func TestToInt(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{"valid int", "123", 123, false},
		{"invalid int", "abc", 0, true},
		{"empty string", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tt.wantErr {
					t.Errorf("ToInt() panicked unexpectedly: %v", r)
				}
			}()

			got := ToInt(tt.input)
			if got != tt.want && !tt.wantErr {
				t.Errorf("ToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt32(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int32
		wantErr bool
	}{
		{"valid int32", "12345", 12345, false},
		{"out of range", "2147483648", 0, true}, // 超过int32最大值
		{"invalid int32", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tt.wantErr {
					t.Errorf("ToInt32() panicked unexpectedly: %v", r)
				}
			}()

			got := ToInt32(tt.input)
			if got != tt.want && !tt.wantErr {
				t.Errorf("ToInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt64(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		{"valid int64", "9223372036854775807", 9223372036854775807, false}, // int64最大值
		{"invalid int64", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tt.wantErr {
					t.Errorf("ToInt64() panicked unexpectedly: %v", r)
				}
			}()

			got := ToInt64(tt.input)
			if got != tt.want && !tt.wantErr {
				t.Errorf("ToInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToFloat32(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    float32
		wantErr bool
	}{
		{"valid float32", "3.14", 3.14, false},
		{"invalid float32", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tt.wantErr {
					t.Errorf("ToFloat32() panicked unexpectedly: %v", r)
				}
			}()

			got := ToFloat32(tt.input)
			if got != tt.want && !tt.wantErr {
				t.Errorf("ToFloat32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToFloat64(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    float64
		wantErr bool
	}{
		{"valid float64", "3.141592653589793", 3.141592653589793, false},
		{"invalid float64", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tt.wantErr {
					t.Errorf("ToFloat64() panicked unexpectedly: %v", r)
				}
			}()

			got := ToFloat64(tt.input)
			if got != tt.want && !tt.wantErr {
				t.Errorf("ToFloat64() = %v, want %v", got, tt.want)
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

// 测试结构体
type testUser struct {
	ID   int
	Name string
	Age  int
}

// 测试指针结构体
type testProduct struct {
	Code  string
	Price float64
}

func TestSlice2Map(t *testing.T) {
	// 测试用例1：正常情况，使用int类型作为key
	t.Run("Normal case with int key", func(t *testing.T) {
		users := []testUser{
			{ID: 1, Name: "Alice", Age: 25},
			{ID: 2, Name: "Bob", Age: 30},
		}

		result, err := Slice2Map[int, testUser]("ID", users)
		if err != nil {
			t.Error("Unexpected error:", err)
		}

		if len(result) != 2 {
			t.Error("Expected map length 2, got", len(result))
		}

		if result[1].Name != "Alice" || result[2].Name != "Bob" {
			t.Error("Map content incorrect")
		}
	})

	// 测试用例2：正常情况，使用string类型作为key
	t.Run("Normal case with string key", func(t *testing.T) {
		users := []testUser{
			{ID: 1, Name: "Alice", Age: 25},
			{ID: 2, Name: "Bob", Age: 30},
		}

		result, err := Slice2Map[string, testUser]("Name", users)
		if err != nil {
			t.Error("Unexpected error:", err)
		}

		if result["Alice"].ID != 1 || result["Bob"].Age != 30 {
			t.Error("Map content incorrect")
		}
	})

	// 测试用例3：测试指针结构体
	t.Run("Pointer struct case", func(t *testing.T) {
		products := []*testProduct{
			{Code: "P001", Price: 9.99},
			{Code: "P002", Price: 19.99},
		}

		result, err := Slice2Map[string, *testProduct]("Code", products)
		if err != nil {
			t.Error("Unexpected error:", err)
		}

		if result["P001"].Price != 9.99 || result["P002"].Price != 19.99 {
			t.Error("Map content incorrect")
		}
	})

	// 测试用例4：字段不存在的情况
	t.Run("Field not found case", func(t *testing.T) {
		users := []testUser{
			{ID: 1, Name: "Alice", Age: 25},
		}

		_, err := Slice2Map[int, testUser]("NonExistentField", users)
		if err == nil {
			t.Error("Expected error for non-existent field, got nil")
		}
	})

	// 测试用例5：字段类型不匹配的情况
	t.Run("Field type mismatch case", func(t *testing.T) {
		users := []testUser{
			{ID: 1, Name: "Alice", Age: 25},
		}

		_, err := Slice2Map[string, testUser]("ID", users)
		if err == nil {
			t.Error("Expected error for type mismatch, got nil")
		}
	})

	// 测试用例6：非结构体类型的情况
	t.Run("Non-struct type case", func(t *testing.T) {
		numbers := []int{1, 2, 3}

		_, err := Slice2Map[int, int]("Field", numbers)
		if err == nil {
			t.Error("Expected error for non-struct type, got nil")
		}
	})

	// 测试用例7：空切片的情况
	t.Run("Empty slice case", func(t *testing.T) {
		var users []testUser

		result, err := Slice2Map[int, testUser]("ID", users)
		if err != nil {
			t.Error("Unexpected error:", err)
		}

		if len(result) != 0 {
			t.Error("Expected empty map, got", len(result))
		}
	})

	// 测试用例8：重复key的情况（应该覆盖）
	t.Run("Duplicate key case", func(t *testing.T) {
		users := []testUser{
			{ID: 1, Name: "Alice", Age: 25},
			{ID: 1, Name: "Bob", Age: 30}, // 相同ID
		}

		result, err := Slice2Map[int, testUser]("ID", users)
		if err != nil {
			t.Error("Unexpected error:", err)
		}

		if len(result) != 1 {
			t.Error("Expected map length 1 due to duplicate keys, got", len(result))
		}

		// 应该保留最后一个值
		if result[1].Name != "Bob" {
			t.Error("Expected last value to be kept for duplicate key")
		}
	})
}
