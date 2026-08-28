package mapUtil

import (
	"encoding/json"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestConcurrentHashMap2CoreAPI(t *testing.T) {
	initial := map[string]int{"one": 1, "zero": 0}
	cm := NewConcurrentHashMap2(initial)
	initial["one"] = 99

	if got := cm.Get("one"); got != 1 {
		t.Fatalf("constructor did not copy initial map: Get(one) = %d, want 1", got)
	}
	if got := cm.Get("missing"); got != 0 {
		t.Errorf("Get(missing) = %d, want zero value", got)
	}
	if !cm.ContainsKey("zero") || cm.GetOrDefault("zero", 10) != 0 {
		t.Error("an existing zero value must be distinguishable through ContainsKey/GetOrDefault")
	}
	if got := cm.GetOrDefault("missing", 10); got != 10 {
		t.Errorf("GetOrDefault(missing) = %d, want 10", got)
	}

	cm.Put("two", 2)
	cm.Put("two", 22)
	if got := cm.Get("two"); got != 22 || cm.Size() != 3 {
		t.Errorf("Put overwrite produced value=%d size=%d, want value=22 size=3", got, cm.Size())
	}
	if existing, loaded := cm.PutIfAbsent("three", 3); loaded || existing != 3 {
		t.Errorf("PutIfAbsent(new) = (%d, %v), want (3, false)", existing, loaded)
	}
	if existing, loaded := cm.PutIfAbsent("three", 33); !loaded || existing != 3 {
		t.Errorf("PutIfAbsent(existing) = (%d, %v), want (3, true)", existing, loaded)
	}

	want := map[string]int{"one": 1, "zero": 0, "two": 22, "three": 3}
	if got := cm.ToMap(); !reflect.DeepEqual(got, want) {
		t.Errorf("ToMap() = %v, want %v", got, want)
	} else {
		got["one"] = 100
		if cm.Get("one") != 1 {
			t.Error("mutating ToMap result changed ConcurrentHashMap2")
		}
	}

	keys := make(map[string]int)
	for _, key := range cm.Keys() {
		keys[key]++
	}
	if !reflect.DeepEqual(keys, map[string]int{"one": 1, "zero": 1, "two": 1, "three": 1}) {
		t.Errorf("Keys() returned wrong multiset: %v", keys)
	}
	values := make(map[int]int)
	for _, value := range cm.Values() {
		values[value]++
	}
	if !reflect.DeepEqual(values, map[int]int{0: 1, 1: 1, 22: 1, 3: 1}) {
		t.Errorf("Values() returned wrong multiset: %v", values)
	}

	visited := make(map[string]int)
	cm.Range(func(key string, value int) bool {
		visited[key] = value
		return true
	})
	if !reflect.DeepEqual(visited, want) {
		t.Errorf("Range() visited %v, want %v", visited, want)
	}
	var stringMap map[string]int
	if err := json.Unmarshal([]byte(cm.ToString()), &stringMap); err != nil || !reflect.DeepEqual(stringMap, want) {
		t.Errorf("ToString() parsed to %v, err=%v, want %v", stringMap, err, want)
	}
	count := 0
	cm.Range(func(string, int) bool {
		count++
		return false
	})
	if count != 1 {
		t.Errorf("Range early termination visited %d entries, want 1", count)
	}

	cm.Remove("two")
	cm.Remove("missing")
	if cm.ContainsKey("two") || cm.Size() != 3 {
		t.Errorf("Remove produced contains(two)=%v size=%d, want false and 3", cm.ContainsKey("two"), cm.Size())
	}
	cm.Clear()
	if cm.Size() != 0 || len(cm.Keys()) != 0 || len(cm.Values()) != 0 {
		t.Errorf("Clear left entries: size=%d keys=%v values=%v", cm.Size(), cm.Keys(), cm.Values())
	}
	cm.Put("reused", 1)
	if cm.Get("reused") != 1 {
		t.Error("map is not reusable after Clear")
	}
}

func TestConcurrentHashMap2Serialization(t *testing.T) {
	t.Run("JSON replaces existing state", func(t *testing.T) {
		cm := NewConcurrentHashMap2(map[string]int{"stale": 9})
		if err := json.Unmarshal([]byte(`{"a":1,"b":2}`), cm); err != nil {
			t.Fatalf("json.Unmarshal() error: %v", err)
		}
		if got, want := cm.ToMap(), map[string]int{"a": 1, "b": 2}; !reflect.DeepEqual(got, want) {
			t.Errorf("JSON result = %v, want %v", got, want)
		}
		data, err := json.Marshal(cm)
		if err != nil {
			t.Fatalf("json.Marshal() error: %v", err)
		}
		var got map[string]int
		if err := json.Unmarshal(data, &got); err != nil || !reflect.DeepEqual(got, map[string]int{"a": 1, "b": 2}) {
			t.Errorf("JSON round trip = %v, err=%v", got, err)
		}
		if err := json.Unmarshal([]byte(`{"a":`), cm); err == nil {
			t.Error("invalid JSON expected an error")
		}
		if got := cm.ToMap(); !reflect.DeepEqual(got, map[string]int{"a": 1, "b": 2}) {
			t.Errorf("invalid JSON mutated receiver: %v", got)
		}
	})

	t.Run("BSON replaces existing state", func(t *testing.T) {
		input, err := bson.Marshal(map[string]int{"a": 1, "b": 2})
		if err != nil {
			t.Fatalf("preparing BSON: %v", err)
		}
		cm := NewConcurrentHashMap2(map[string]int{"stale": 9})
		if err := bson.Unmarshal(input, cm); err != nil {
			t.Fatalf("bson.Unmarshal() error: %v", err)
		}
		if got, want := cm.ToMap(), map[string]int{"a": 1, "b": 2}; !reflect.DeepEqual(got, want) {
			t.Errorf("BSON result = %v, want %v", got, want)
		}
		if _, err := bson.Marshal(cm); err != nil {
			t.Errorf("bson.Marshal() error: %v", err)
		}
		if err := bson.Unmarshal([]byte{0, 1, 2}, cm); err == nil {
			t.Error("invalid BSON expected an error")
		}
		if got, want := cm.ToMap(), map[string]int{"a": 1, "b": 2}; !reflect.DeepEqual(got, want) {
			t.Errorf("invalid BSON mutated receiver: %v", got)
		}
	})
}

func TestConcurrentHashMap2ConcurrentOperations(t *testing.T) {
	const (
		workers = 16
		per     = 200
	)
	cm := NewConcurrentHashMap2[int, int]()
	var wg sync.WaitGroup
	wg.Add(workers)
	for worker := 0; worker < workers; worker++ {
		go func(worker int) {
			defer wg.Done()
			for offset := 0; offset < per; offset++ {
				key := worker*per + offset
				cm.Put(key, key*2)
			}
		}(worker)
	}
	wg.Wait()

	if got, want := cm.Size(), workers*per; got != want {
		t.Fatalf("concurrent Put size = %d, want %d", got, want)
	}
	for key := 0; key < workers*per; key++ {
		if got := cm.Get(key); got != key*2 {
			t.Fatalf("Get(%d) = %d, want %d", key, got, key*2)
		}
	}

	wg.Add(workers * 2)
	for worker := 0; worker < workers; worker++ {
		go func(worker int) {
			defer wg.Done()
			for offset := 0; offset < per; offset++ {
				key := worker*per + offset
				_ = cm.Get(key)
				_ = cm.ContainsKey(key)
			}
		}(worker)
		go func(worker int) {
			defer wg.Done()
			for offset := 0; offset < per; offset += 2 {
				cm.Remove(worker*per + offset)
			}
		}(worker)
	}
	wg.Wait()

	if got, want := cm.Size(), workers*per/2; got != want {
		t.Fatalf("concurrent read/remove size = %d, want %d", got, want)
	}
	for key := 0; key < workers*per; key++ {
		if got := cm.ContainsKey(key); got != (key%2 == 1) {
			t.Errorf("ContainsKey(%d) = %v, want %v", key, got, key%2 == 1)
		}
	}

	var inserted atomic.Int32
	wg.Add(workers)
	for worker := 0; worker < workers; worker++ {
		go func(value int) {
			defer wg.Done()
			if _, loaded := cm.PutIfAbsent(-1, value); !loaded {
				inserted.Add(1)
			}
		}(worker)
	}
	wg.Wait()
	if inserted.Load() != 1 || !cm.ContainsKey(-1) {
		t.Errorf("concurrent PutIfAbsent inserted=%d contains=%v, want 1 and true", inserted.Load(), cm.ContainsKey(-1))
	}
}

func TestConcurrentHashMap2ConcurrentClear(t *testing.T) {
	cm := NewConcurrentHashMap2[int, int]()
	var wg sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		wg.Add(2)
		go func(worker int) {
			defer wg.Done()
			for i := 0; i < 200; i++ {
				cm.Put(worker*200+i, i)
			}
		}(worker)
		go func() {
			defer wg.Done()
			for i := 0; i < 20; i++ {
				cm.Clear()
			}
		}()
	}
	wg.Wait()

	cm.Clear()
	if cm.Size() != 0 {
		t.Fatalf("final Clear size = %d, want 0", cm.Size())
	}
	cm.Put(1, 10)
	if got := cm.ToMap(); !reflect.DeepEqual(got, map[int]int{1: 10}) {
		t.Errorf("map unusable after concurrent Clear: %v", got)
	}
}
