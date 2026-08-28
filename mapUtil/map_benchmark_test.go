package mapUtil

import "testing"

const benchmarkMapEntries = 1024

var benchmarkMapSink int

func BenchmarkMapGet(b *testing.B) {
	b.Run("native", func(b *testing.B) {
		m := make(map[int]int, benchmarkMapEntries)
		for i := 0; i < benchmarkMapEntries; i++ {
			m[i] = i
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchmarkMapSink = m[i%benchmarkMapEntries]
		}
	})

	b.Run("concurrent_hash", func(b *testing.B) {
		m := NewConcurrentHashMap[int, int]()
		for i := 0; i < benchmarkMapEntries; i++ {
			m.Put(i, i)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchmarkMapSink = m.Get(i % benchmarkMapEntries)
		}
	})

	b.Run("sync_map", func(b *testing.B) {
		m := NewConcurrentHashMap2[int, int]()
		for i := 0; i < benchmarkMapEntries; i++ {
			m.Put(i, i)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchmarkMapSink = m.Get(i % benchmarkMapEntries)
		}
	})

	b.Run("skip_list", func(b *testing.B) {
		m := NewConcurrentSkipListMap[int, int]()
		for i := 0; i < benchmarkMapEntries; i++ {
			m.Put(i, i)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchmarkMapSink = m.Get(i % benchmarkMapEntries)
		}
	})

	b.Run("ordered", func(b *testing.B) {
		m := NewOrderedMap[int, int]()
		for i := 0; i < benchmarkMapEntries; i++ {
			m.Put(i, i)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchmarkMapSink = m.Get(i % benchmarkMapEntries)
		}
	})

	b.Run("tree", func(b *testing.B) {
		m := NewTreeMap[int, int](func(a, b int) bool { return a < b })
		for i := 0; i < benchmarkMapEntries; i++ {
			m.Put(i, i)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchmarkMapSink = m.Get(i % benchmarkMapEntries)
		}
	})
}

func BenchmarkMapPut(b *testing.B) {
	b.Run("native", func(b *testing.B) {
		m := make(map[int]int)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m[i] = i
		}
	})
	b.Run("concurrent_hash", func(b *testing.B) {
		m := NewConcurrentHashMap[int, int]()
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m.Put(i, i)
		}
	})
	b.Run("sync_map", func(b *testing.B) {
		m := NewConcurrentHashMap2[int, int]()
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m.Put(i, i)
		}
	})
	b.Run("skip_list", func(b *testing.B) {
		m := NewConcurrentSkipListMap[int, int]()
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m.Put(i, i)
		}
	})
	b.Run("ordered", func(b *testing.B) {
		m := NewOrderedMap[int, int]()
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m.Put(i, i)
		}
	})
	b.Run("tree", func(b *testing.B) {
		m := NewTreeMap[int, int](func(a, b int) bool { return a < b })
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m.Put(i, i)
		}
	})
}
