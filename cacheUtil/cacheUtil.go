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

const (
	cleanupInterval = 30 * time.Second // 固定清理间隔
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

// Set 方法用于向缓存中设置一个键值对，并可选择性地指定该键值对的过期时间。
// 如果未提供过期时间，则使用缓存实例的默认过期时间。
// 参数 k 为缓存的键，类型为 K。
// 参数 x 为缓存的值，类型为 V。
// 参数 expiration 为可选的过期时间，可传入 0 个或 1 个 time.Duration 类型的值。
func (c *cache[K, V]) Set(k K, x V, expiration ...time.Duration) {
	c.mu.Lock()
	exp := time.Now().Add(c.expiration).UnixNano()
	if len(expiration) > 0 {
		exp = time.Now().Add(expiration[0]).UnixNano()
	}
	c.items[k] = Item[V]{
		Object:     x,
		Expiration: exp,
	}
	c.mu.Unlock()
}

// set 是一个辅助方法，用于向缓存中设置一个键值对，使用缓存实例的默认过期时间。
// 该方法不会加锁，调用时需要确保已经获取了写锁，避免并发修改问题。
// 参数 k 为缓存的键，类型为 K。
// 参数 x 为缓存的值，类型为 V。
func (c *cache[K, V]) set(k K, x V) {
	c.items[k] = Item[V]{
		Object:     x,
		Expiration: time.Now().Add(c.expiration).UnixNano(),
	}
}

// SetIfAbsent 方法用于在缓存中键不存在时设置键值对。若键已存在，该方法不会修改缓存，直接返回 false。
// 参数 k 为要检查和设置的缓存键，类型为 K。
// 参数 x 为要设置的缓存值，类型为 V。
// 返回值为 bool 类型，若成功设置键值对返回 true，若键已存在返回 false。
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

// Get 根据键从缓存中获取对应的值。
// 若缓存项存在且未过期，返回该值和 true；若缓存项不存在或已过期，返回对应类型的零值和 false。
// 参数 k 为要查找的缓存键。
// 返回值依次为缓存项的值、缓存项是否存在且未过期的标志。
func (c *cache[K, V]) Get(k K) (V, bool) {
	c.mu.RLock()
	item, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		var zero V
		return zero, false
	}
	if c.expiration > 0 && time.Now().UnixNano() > item.Expiration {
		c.mu.RUnlock()
		var zero V
		return zero, false
	}
	c.mu.RUnlock()
	return item.Object, true
}

// get 是一个辅助方法，用于根据键从缓存中获取对应的值。
// 该方法不会加锁，调用时需要确保已经获取了读锁，避免并发修改问题。
// 若缓存项存在且未过期，返回该值和 true；若缓存项不存在或已过期，返回对应类型的零值和 false。
// 参数 k 为要查找的缓存键。
// 返回值依次为缓存项的值、缓存项是否存在且未过期的标志。
func (c *cache[K, V]) get(k K) (V, bool) {
	item, found := c.items[k]
	if !found {
		var zero V
		return zero, false
	}
	if c.expiration > 0 && time.Now().UnixNano() > item.Expiration {
		var zero V
		return zero, false
	}
	return item.Object, true
}

// GetWithExpiration 根据键获取缓存项，并返回缓存项的值、过期时间以及是否存在的标志。
// 若缓存项不存在或已过期，将返回对应类型的零值、零时间和 false。
// 参数 k 为要查找的缓存键。
// 返回值依次为缓存项的值、缓存项的过期时间、缓存项是否存在且未过期的标志。
func (c *cache[K, V]) GetWithExpiration(k K) (V, time.Time, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[k]
	if !found {
		var zero V
		return zero, time.Time{}, false
	}

	if c.expiration > 0 && time.Now().UnixNano() > item.Expiration {
		var zero V
		return zero, time.Time{}, false
	}

	// 将纳秒时间戳转换为time.Time
	expirationTime := time.Unix(0, item.Expiration)
	if c.expiration == 0 {
		expirationTime = time.Unix(0, 0)
	}
	return item.Object, expirationTime, true
}

// Delete 方法用于从缓存中删除指定键对应的缓存项。
// 该方法会加写锁，确保在删除操作过程中不会被其他 goroutine 并发修改。
// 参数 k 为要删除的缓存键，类型为 K。
func (c *cache[K, V]) Delete(k K) {
	c.mu.Lock()
	c.delete(k)
	c.mu.Unlock()
}

// delete 是一个辅助方法，用于从缓存中删除指定键对应的缓存项。
// 该方法不会加锁，调用时需要确保已经获取了写锁，避免并发修改问题。
// 参数 k 为要删除的缓存键，类型为 K。
func (c *cache[K, V]) delete(k K) {
	delete(c.items, k)
}

// deleteExpired 方法用于删除缓存中所有已过期的键值对。
// 若缓存没有设置过期时间，该方法将直接返回，不进行任何操作。
func (c *cache[K, V]) deleteExpired() {
	if c.expiration <= 0 {
		return
	}
	now := time.Now().UnixNano()
	c.mu.Lock()
	for k, v := range c.items {
		if v.Expiration > 0 && now > v.Expiration {
			c.delete(k)
		}
	}
	c.mu.Unlock()
}

// Items 方法用于获取缓存中所有未过期的键值对。
// 该方法会加读锁，确保在遍历缓存时不会被其他 goroutine 并发修改，同时允许多个读操作并发执行。
// 返回一个包含所有未过期键值对的新 map。
func (c *cache[K, V]) Items() map[K]Item[V] {
	c.mu.RLock()
	defer c.mu.RUnlock()
	m := make(map[K]Item[V], len(c.items))
	now := time.Now().UnixNano()
	for k, v := range c.items {
		if c.expiration > 0 && v.Expiration > 0 {
			if now > v.Expiration {
				continue
			}
		}
		m[k] = v
	}
	return m
}

// Flush 方法用于清空缓存中的所有键值对。
// 该方法会加写锁，确保在清空操作过程中不会被其他 goroutine 并发修改，
// 操作完成后再释放写锁。
func (c *cache[K, V]) Flush() {
	c.mu.Lock()
	c.items = make(map[K]Item[V])
	c.mu.Unlock()
}

type janitor[K comparable, V any] struct {
	Interval time.Duration
	stop     chan bool
}

// run 方法用于启动一个定时任务，定期清理缓存中已过期的键值对。
// 该方法会在一个独立的 goroutine 中运行，通过定时器按指定间隔触发清理操作。
// 参数 c 为需要清理的缓存实例。
func (j *janitor[K, V]) run(c *cache[K, V]) {
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
	go j.run(c)
}

// New 函数用于创建一个新的缓存实例。
// 该函数接受一个过期时间作为参数，若过期时间大于 0，会自动启动一个定时任务来清理过期的缓存项。
// 参数 expiration 为缓存项的默认过期时间，类型为 time.Duration。
// 返回值为指向 Cache[K, V] 类型的指针，代表新创建的缓存实例。
func New[K comparable, V any](expiration time.Duration) *Cache[K, V] {
	c := &cache[K, V]{
		expiration: expiration,
		items:      make(map[K]Item[V]),
	}
	C := &Cache[K, V]{c}
	if expiration > 0 {
		runJanitor(c) // 自动启用janitor
		runtime.SetFinalizer(C, stopJanitor[K, V])
	}
	return C
}
