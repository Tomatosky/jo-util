package cacheUtil

import (
	"runtime"
	"sync"
	"time"
)

type Item[V any] struct {
	Object     V
	Expiration int64
}

func (item Item[V]) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

const (
	cleanupInterval = 1 * time.Minute // 固定清理间隔
)

type Cache[K comparable, V any] struct {
	*cache[K, V]
}

type cache[K comparable, V any] struct {
	expiration time.Duration
	items      map[K]Item[V]
	mu         sync.RWMutex
	janitor    *janitor[K, V]
}

func (c *cache[K, V]) Set(k K, x V) {
	c.mu.Lock()
	c.items[k] = Item[V]{
		Object:     x,
		Expiration: time.Now().Add(c.expiration).UnixNano(),
	}
	c.mu.Unlock()
}

func (c *cache[K, V]) set(k K, x V) {
	c.items[k] = Item[V]{
		Object:     x,
		Expiration: time.Now().Add(c.expiration).UnixNano(),
	}
}

func (c *cache[K, V]) SetIfAbsent(k K, x V) bool {
	c.mu.Lock()
	_, found := c.get(k)
	if found {
		c.mu.Unlock()
		return false
	}
	c.set(k, x)
	c.mu.Unlock()
	return true
}

func (c *cache[K, V]) Get(k K) (V, bool) {
	c.mu.RLock()
	item, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		var zero V
		return zero, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.RUnlock()
			var zero V
			return zero, false
		}
	}
	c.mu.RUnlock()
	return item.Object, true
}

func (c *cache[K, V]) get(k K) (V, bool) {
	item, found := c.items[k]
	if !found {
		var zero V
		return zero, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			var zero V
			return zero, false
		}
	}
	return item.Object, true
}

func (c *cache[K, V]) Delete(k K) {
	c.mu.Lock()
	c.delete(k)
}

func (c *cache[K, V]) delete(k K) {
	delete(c.items, k)
}

func (c *cache[K, V]) deleteExpired() {
	now := time.Now().UnixNano()
	c.mu.Lock()
	for k, v := range c.items {
		if v.Expiration > 0 && now > v.Expiration {
			c.delete(k)
		}
	}
	c.mu.Unlock()
}

func (c *cache[K, V]) Items() map[K]Item[V] {
	c.mu.RLock()
	defer c.mu.RUnlock()
	m := make(map[K]Item[V], len(c.items))
	now := time.Now().UnixNano()
	for k, v := range c.items {
		if v.Expiration > 0 {
			if now > v.Expiration {
				continue
			}
		}
		m[k] = v
	}
	return m
}

func (c *cache[K, V]) ItemCount() int {
	c.mu.RLock()
	n := len(c.items)
	c.mu.RUnlock()
	return n
}

func (c *cache[K, V]) Flush() {
	c.mu.Lock()
	c.items = make(map[K]Item[V])
	c.mu.Unlock()
}

type janitor[K comparable, V any] struct {
	Interval time.Duration
	stop     chan bool
}

func (j *janitor[K, V]) Run(c *cache[K, V]) {
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			c.deleteExpired()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

func stopJanitor[K comparable, V any](c *Cache[K, V]) {
	c.janitor.stop <- true
}

func runJanitor[K comparable, V any](c *cache[K, V]) {
	j := &janitor[K, V]{
		Interval: cleanupInterval, // 使用固定间隔
		stop:     make(chan bool),
	}
	c.janitor = j
	go j.Run(c)
}

func newCache[K comparable, V any](de time.Duration, m map[K]Item[V]) *cache[K, V] {
	if de == 0 {
		de = -1
	}
	return &cache[K, V]{
		expiration: de,
		items:      m,
	}
}

func newCacheWithJanitor[K comparable, V any](expiration time.Duration, m map[K]Item[V]) *Cache[K, V] {
	c := newCache(expiration, m)
	C := &Cache[K, V]{c}
	if expiration > 0 {
		runJanitor(c) // 自动启用janitor
		runtime.SetFinalizer(C, stopJanitor[K, V])
	}
	return C
}

func New[K comparable, V any](expiration time.Duration) *Cache[K, V] {
	items := make(map[K]Item[V])
	return newCacheWithJanitor[K, V](expiration, items)
}
