package randomUtil

import (
	"testing"
)

func TestRandomInt(t *testing.T) {
	tests := []struct {
		name    string
		start   int
		end     int
		wantErr bool
	}{
		{"valid range", 1, 10, false},
		{"invalid range", 10, 1, true},
		{"same values", 5, 5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tt.wantErr {
					t.Errorf("RandomInt() panicked unexpectedly: %v", r)
				}
			}()

			got := RandomInt(tt.start, tt.end)
			if tt.wantErr {
				t.Error("Expected panic but none occurred")
				return
			}

			if got < tt.start || got >= tt.end {
				t.Errorf("RandomInt() = %v, want in range [%v, %v)", got, tt.start, tt.end)
			}
		})
	}
}

func TestRandomEle(t *testing.T) {
	t.Run("non-empty slice", func(t *testing.T) {
		slice := []string{"a", "b", "c", "d"}
		got := RandomEle(slice)
		found := false
		for _, v := range slice {
			if v == got {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("RandomEle() = %v, not found in slice %v", got, slice)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for empty slice but none occurred")
			}
		}()
		RandomEle([]int{})
	})
}

func TestRandomEleSet(t *testing.T) {
	tests := []struct {
		name   string
		slice  []int
		n      int
		length int
	}{
		{"n < len", []int{1, 2, 3, 4, 5}, 3, 3},
		{"n > len", []int{1, 2}, 5, 2},
		{"n = len", []int{1, 2, 3}, 3, 3},
		{"n = 0", []int{1, 2, 3}, 0, 0},
		{"n < 0", []int{1, 2, 3}, -1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RandomEleSet(tt.slice, tt.n)
			if len(got) != tt.length {
				t.Errorf("RandomEleSet() length = %v, want %v", len(got), tt.length)
			}

			// Check for duplicates
			seen := make(map[int]bool)
			for _, v := range got {
				if seen[v] {
					t.Errorf("RandomEleSet() returned duplicate value: %v", v)
				}
				seen[v] = true
			}
		})
	}

	t.Run("empty slice", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for empty slice but none occurred")
			}
		}()
		RandomEleSet([]string{}, 1)
	})
}

func TestRandomWeightedKey(t *testing.T) {
	tests := []struct {
		name    string
		weights map[string]int
		wantErr bool
	}{
		{"valid weights", map[string]int{"a": 1, "b": 2, "c": 3}, false},
		{"all zero weights", map[string]int{"a": 0, "b": 0, "c": 0}, true},
		{"some zero weights", map[string]int{"a": 1, "b": 0, "c": 2}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tt.wantErr {
					t.Errorf("RandomWeightedKey() panicked unexpectedly: %v", r)
				}
			}()

			if tt.wantErr {
				func() {
					defer func() {
						if r := recover(); r == nil {
							t.Error("Expected panic but none occurred")
						}
					}()
					RandomWeightedKey(tt.weights)
				}()
				return
			}

			got := RandomWeightedKey(tt.weights)
			if _, exists := tt.weights[got]; !exists {
				t.Errorf("RandomWeightedKey() = %v, not found in weights map", got)
			}
		})
	}
}

func TestRandomString(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"length 5", 5},
		{"length 10", 10},
		{"length 0", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RandomString(tt.length)
			if len(got) != tt.length {
				t.Errorf("RandomString() length = %v, want %v", len(got), tt.length)
			}
		})
	}
}

func TestRandomNumbers(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"length 5", 5},
		{"length 10", 10},
		{"length 0", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RandomNumbers(tt.length)
			if len(got) != tt.length {
				t.Errorf("RandomNumbers() length = %v, want %v", len(got), tt.length)
			}

			for _, c := range got {
				if c < '0' || c > '9' {
					t.Errorf("RandomNumbers() contains non-digit character: %c", c)
				}
			}
		})
	}
}
