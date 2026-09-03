package randomUtil

import (
	"math"
	"reflect"
	"strings"
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

func requirePanicValue(t *testing.T, want any, f func()) {
	t.Helper()
	defer func() {
		if got := recover(); got != want {
			t.Errorf("panic = %v, want %v", got, want)
		}
	}()
	f()
}

func TestRandomInt(t *testing.T) {
	for i := 0; i < 500; i++ {
		if got := RandomInt(-10, 10); got < -10 || got >= 10 {
			t.Fatalf("RandomInt(-10, 10) = %d", got)
		}
		if got := RandomInt[int64](5, 6); got != 5 {
			t.Fatalf("RandomInt[int64](5, 6) = %d, want 5", got)
		}
		if got := RandomInt[uint](1, 4); got < 1 || got >= 4 {
			t.Fatalf("RandomInt[uint](1, 4) = %v", got)
		}
	}
	requirePanic(t, func() { RandomInt(10, 1) })
	requirePanic(t, func() { RandomInt(5, 5) })
	if got := RandomInt[uint64](0, math.MaxUint64); got >= math.MaxUint64 {
		t.Fatalf("RandomInt[uint64](0, MaxUint64) = %d", got)
	}
	if got := RandomInt[int64](math.MinInt64, math.MaxInt64); got >= math.MaxInt64 {
		t.Fatalf("RandomInt[int64](MinInt64, MaxInt64) = %d", got)
	}
}

func TestRandomEle(t *testing.T) {
	if got := RandomEle([]string{"only"}); got != "only" {
		t.Errorf("RandomEle(single element) = %q, want only", got)
	}
	input := []string{"a", "b", "c", "d"}
	allowed := map[string]bool{"a": true, "b": true, "c": true, "d": true}
	for i := 0; i < 200; i++ {
		if got := RandomEle(input); !allowed[got] {
			t.Fatalf("RandomEle() = %q, not in input", got)
		}
	}
	requirePanic(t, func() { RandomEle([]int{}) })
}

func TestRandomEleSet(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want int
	}{
		{name: "subset", n: 3, want: 3},
		{name: "clamped to input", n: 10, want: 5},
		{name: "entire input", n: 5, want: 5},
		{name: "zero", n: 0, want: 0},
		{name: "negative", n: -1, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := []int{1, 2, 3, 4, 5}
			original := append([]int(nil), input...)
			got := RandomEleSet(input, tt.n)
			if len(got) != tt.want {
				t.Fatalf("length = %d, want %d", len(got), tt.want)
			}
			if !reflect.DeepEqual(input, original) {
				t.Errorf("input mutated: got %v, want %v", input, original)
			}
			seen := make(map[int]bool)
			for _, value := range got {
				if value < 1 || value > 5 {
					t.Errorf("result value %d is not from input", value)
				}
				if seen[value] {
					t.Errorf("duplicate result value %d", value)
				}
				seen[value] = true
			}
			if tt.want == len(input) && len(seen) != len(input) {
				t.Errorf("full selection = %v, want every input element", got)
			}
		})
	}
	requirePanic(t, func() { RandomEleSet([]string{}, 1) })
	if got := RandomEleSet([]string{}, 0); got != nil {
		t.Errorf("RandomEleSet(empty, 0) = %#v, want nil", got)
	}
}

func TestRandomWeightedKey(t *testing.T) {
	weights := map[string]int{"zero-a": 0, "selected": 10, "zero-b": 0}
	for i := 0; i < 500; i++ {
		if got := RandomWeightedKey(weights); got != "selected" {
			t.Fatalf("RandomWeightedKey() = %q, want selected", got)
		}
	}
	requirePanicValue(t, "所有权重值总和不能为0", func() {
		RandomWeightedKey(map[string]int{"a": 0, "b": 0})
	})
	requirePanicValue(t, "所有权重值总和不能为0", func() {
		RandomWeightedKey(map[string]int{})
	})
	for _, weights := range []map[string]int{
		{"negative": -1},
		{"negative": -1, "positive": 10},
		{"negative": -10, "positive": 1},
	} {
		requirePanicValue(t, "权重值不能为负数", func() {
			RandomWeightedKey(weights)
		})
	}
	if got := RandomWeightedKey(map[string]uint64{"selected": math.MaxUint64}); got != "selected" {
		t.Errorf("RandomWeightedKey(MaxUint64 weight) = %q, want selected", got)
	}
	requirePanicValue(t, "权重值总和溢出", func() {
		RandomWeightedKey(map[string]uint64{"a": math.MaxUint64, "b": 1})
	})
}

func TestRandomString(t *testing.T) {
	const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for _, length := range []int{0, 1, 32, 256} {
		got := RandomString(length)
		if len(got) != length {
			t.Errorf("RandomString(%d) length = %d", length, len(got))
		}
		for _, char := range got {
			if !strings.ContainsRune(alphabet, char) {
				t.Errorf("RandomString(%d) contains %q outside alphabet", length, char)
			}
		}
	}
	requirePanic(t, func() { RandomString(-1) })
}

func TestRandomNumbers(t *testing.T) {
	for _, length := range []int{0, 1, 32, 256} {
		got := RandomNumbers(length)
		if len(got) != length {
			t.Errorf("RandomNumbers(%d) length = %d", length, len(got))
		}
		for _, char := range got {
			if char < '0' || char > '9' {
				t.Errorf("RandomNumbers(%d) contains non-digit %q", length, char)
			}
		}
	}
	requirePanic(t, func() { RandomNumbers(-1) })
}
