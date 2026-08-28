package idUtil

import (
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
)

func TestUUIDFormats(t *testing.T) {
	standard := RandomUUID()
	parsed, err := uuid.Parse(standard)
	if err != nil {
		t.Fatalf("RandomUUID() = %q: %v", standard, err)
	}
	if parsed.Version() != uuid.Version(4) || parsed.Variant() != uuid.RFC4122 {
		t.Errorf("RandomUUID() version/variant = (%v, %v), want (4, RFC4122)", parsed.Version(), parsed.Variant())
	}
	if len(standard) != 36 || strings.Count(standard, "-") != 4 {
		t.Errorf("RandomUUID() format = %q", standard)
	}

	simple := SimpleUUID()
	if len(simple) != 32 || strings.Contains(simple, "-") {
		t.Errorf("SimpleUUID() format = %q", simple)
	}
	parsed, err = uuid.Parse(simple)
	if err != nil {
		t.Fatalf("SimpleUUID() = %q: %v", simple, err)
	}
	if parsed.Version() != uuid.Version(4) || parsed.Variant() != uuid.RFC4122 {
		t.Errorf("SimpleUUID() version/variant = (%v, %v), want (4, RFC4122)", parsed.Version(), parsed.Variant())
	}
}

func TestUUIDConcurrentUniqueness(t *testing.T) {
	const goroutines = 20
	const perGoroutine = 100
	values := make(chan string, goroutines*perGoroutine*2)
	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < perGoroutine; j++ {
				values <- RandomUUID()
				values <- SimpleUUID()
			}
		}()
	}
	wg.Wait()
	close(values)

	seen := make(map[string]struct{}, goroutines*perGoroutine*2)
	for value := range values {
		if _, exists := seen[value]; exists {
			t.Fatalf("duplicate UUID %q", value)
		}
		seen[value] = struct{}{}
	}
	if got, want := len(seen), goroutines*perGoroutine*2; got != want {
		t.Errorf("unique UUID count = %d, want %d", got, want)
	}
}

func TestSnowflakeIDs(t *testing.T) {
	previous := GetSnowflakeNextId()
	if previous <= 0 {
		t.Fatalf("GetSnowflakeNextId() = %d, want positive", previous)
	}
	for i := 0; i < 1000; i++ {
		current := GetSnowflakeNextId()
		if current <= previous {
			t.Fatalf("snowflake IDs not strictly increasing: previous=%d current=%d", previous, current)
		}
		previous = current
	}

	stringID := GetSnowflakeNextIdStr()
	parsed, err := strconv.ParseInt(stringID, 10, 64)
	if err != nil {
		t.Fatalf("GetSnowflakeNextIdStr() = %q: %v", stringID, err)
	}
	if parsed <= previous {
		t.Errorf("string snowflake ID %d is not newer than previous %d", parsed, previous)
	}
}

func TestSnowflakeConcurrentUniqueness(t *testing.T) {
	const total = 5000
	values := make(chan int64, total)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < total/20; j++ {
				values <- GetSnowflakeNextId()
			}
		}()
	}
	wg.Wait()
	close(values)

	seen := make(map[int64]struct{}, total)
	for value := range values {
		if value <= 0 {
			t.Fatalf("non-positive snowflake ID %d", value)
		}
		if _, exists := seen[value]; exists {
			t.Fatalf("duplicate snowflake ID %d", value)
		}
		seen[value] = struct{}{}
	}
	if len(seen) != total {
		t.Errorf("unique snowflake count = %d, want %d", len(seen), total)
	}
}
