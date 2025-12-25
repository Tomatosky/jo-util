package mapUtil

import (
	"cmp"
	"encoding/json"
	"math/rand"
	"runtime/debug"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
)

var _ IMap[string, int] = (*ConcurrentSkipListMap[string, int])(nil)

const (
	maxLevel     = 32   // 跳表最大层级
	probability  = 0.25 // 晋升概率
)

// skipListNode 跳表节点
type skipListNode[K cmp.Ordered, V any] struct {
	key     K
	value   V
	forward []*skipListNode[K, V] // 每一层的前向指针
}

// ConcurrentSkipListMap 并发安全的有序映射,基于跳表实现
type ConcurrentSkipListMap[K cmp.Ordered, V any] struct {
	mu     sync.RWMutex
	head   *skipListNode[K, V]
	level  int // 当前跳表的最高层级
	length int // 元素个数
}

// NewConcurrentSkipListMap 创建一个新的ConcurrentSkipListMap
func NewConcurrentSkipListMap[K cmp.Ordered, V any](initMap ...map[K]V) *ConcurrentSkipListMap[K, V] {
	head := &skipListNode[K, V]{
		forward: make([]*skipListNode[K, V], maxLevel),
	}

	csm := &ConcurrentSkipListMap[K, V]{
		head:  head,
		level: 1,
	}

	// 如果有传入初始化map
	if len(initMap) > 0 && initMap[0] != nil {
		for k, v := range initMap[0] {
			csm.Put(k, v)
		}
	}

	return csm
}

// randomLevel 随机生成节点层级
func (csm *ConcurrentSkipListMap[K, V]) randomLevel() int {
	level := 1
	for rand.Float64() < probability && level < maxLevel {
		level++
	}
	return level
}

// findPredecessors 查找每一层的前驱节点
func (csm *ConcurrentSkipListMap[K, V]) findPredecessors(key K) []*skipListNode[K, V] {
	update := make([]*skipListNode[K, V], maxLevel)
	current := csm.head

	for i := csm.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && cmp.Compare(current.forward[i].key, key) < 0 {
			current = current.forward[i]
		}
		update[i] = current
	}

	return update
}

// Get 获取指定键的值
func (csm *ConcurrentSkipListMap[K, V]) Get(key K) V {
	csm.mu.RLock()
	defer csm.mu.RUnlock()

	current := csm.head
	for i := csm.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && cmp.Compare(current.forward[i].key, key) < 0 {
			current = current.forward[i]
		}
	}

	current = current.forward[0]
	if current != nil && cmp.Compare(current.key, key) == 0 {
		return current.value
	}

	var zero V
	return zero
}

// Put 插入或更新键值对
func (csm *ConcurrentSkipListMap[K, V]) Put(key K, value V) {
	csm.mu.Lock()
	defer csm.mu.Unlock()

	update := csm.findPredecessors(key)
	current := update[0].forward[0]

	// 如果key已存在,更新值
	if current != nil && cmp.Compare(current.key, key) == 0 {
		current.value = value
		return
	}

	// 插入新节点
	newLevel := csm.randomLevel()
	if newLevel > csm.level {
		for i := csm.level; i < newLevel; i++ {
			update[i] = csm.head
		}
		csm.level = newLevel
	}

	newNode := &skipListNode[K, V]{
		key:     key,
		value:   value,
		forward: make([]*skipListNode[K, V], newLevel),
	}

	for i := 0; i < newLevel; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}

	csm.length++
}

// Remove 删除指定键
func (csm *ConcurrentSkipListMap[K, V]) Remove(key K) {
	csm.mu.Lock()
	defer csm.mu.Unlock()

	update := csm.findPredecessors(key)
	current := update[0].forward[0]

	if current == nil || cmp.Compare(current.key, key) != 0 {
		return
	}

	// 从各层中删除节点
	for i := 0; i < csm.level; i++ {
		if update[i].forward[i] != current {
			break
		}
		update[i].forward[i] = current.forward[i]
	}

	// 更新level
	for csm.level > 1 && csm.head.forward[csm.level-1] == nil {
		csm.level--
	}

	csm.length--
}

// Size 返回元素个数
func (csm *ConcurrentSkipListMap[K, V]) Size() int {
	csm.mu.RLock()
	defer csm.mu.RUnlock()
	return csm.length
}

// ContainsKey 检查是否包含指定键
func (csm *ConcurrentSkipListMap[K, V]) ContainsKey(key K) bool {
	csm.mu.RLock()
	defer csm.mu.RUnlock()

	current := csm.head
	for i := csm.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && cmp.Compare(current.forward[i].key, key) < 0 {
			current = current.forward[i]
		}
	}

	current = current.forward[0]
	return current != nil && cmp.Compare(current.key, key) == 0
}

// Clear 清空所有元素
func (csm *ConcurrentSkipListMap[K, V]) Clear() {
	csm.mu.Lock()
	defer csm.mu.Unlock()

	csm.head = &skipListNode[K, V]{
		forward: make([]*skipListNode[K, V], maxLevel),
	}
	csm.level = 1
	csm.length = 0
}

// Keys 返回所有键(有序)
func (csm *ConcurrentSkipListMap[K, V]) Keys() []K {
	csm.mu.RLock()
	defer csm.mu.RUnlock()

	keys := make([]K, 0, csm.length)
	current := csm.head.forward[0]

	for current != nil {
		keys = append(keys, current.key)
		current = current.forward[0]
	}

	return keys
}

// Values 返回所有值(按键的顺序)
func (csm *ConcurrentSkipListMap[K, V]) Values() []V {
	csm.mu.RLock()
	defer csm.mu.RUnlock()

	values := make([]V, 0, csm.length)
	current := csm.head.forward[0]

	for current != nil {
		values = append(values, current.value)
		current = current.forward[0]
	}

	return values
}

// PutIfAbsent 如果键不存在则插入
func (csm *ConcurrentSkipListMap[K, V]) PutIfAbsent(key K, value V) (existing V, loaded bool) {
	csm.mu.Lock()
	defer csm.mu.Unlock()

	update := csm.findPredecessors(key)
	current := update[0].forward[0]

	// 如果key已存在,返回现有值
	if current != nil && cmp.Compare(current.key, key) == 0 {
		return current.value, true
	}

	// 插入新节点
	newLevel := csm.randomLevel()
	if newLevel > csm.level {
		for i := csm.level; i < newLevel; i++ {
			update[i] = csm.head
		}
		csm.level = newLevel
	}

	newNode := &skipListNode[K, V]{
		key:     key,
		value:   value,
		forward: make([]*skipListNode[K, V], newLevel),
	}

	for i := 0; i < newLevel; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}

	csm.length++

	return value, false
}

// GetOrDefault 获取值或返回默认值
func (csm *ConcurrentSkipListMap[K, V]) GetOrDefault(key K, defaultValue V) V {
	csm.mu.RLock()
	defer csm.mu.RUnlock()

	current := csm.head
	for i := csm.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && cmp.Compare(current.forward[i].key, key) < 0 {
			current = current.forward[i]
		}
	}

	current = current.forward[0]
	if current != nil && cmp.Compare(current.key, key) == 0 {
		return current.value
	}

	return defaultValue
}

// ToMap 转换为普通map
func (csm *ConcurrentSkipListMap[K, V]) ToMap() map[K]V {
	csm.mu.RLock()
	defer csm.mu.RUnlock()

	result := make(map[K]V, csm.length)
	current := csm.head.forward[0]

	for current != nil {
		result[current.key] = current.value
		current = current.forward[0]
	}

	return result
}

// Range 遍历元素(按键的顺序,返回false可提前终止)
func (csm *ConcurrentSkipListMap[K, V]) Range(f func(key K, value V) bool) {
	csm.mu.RLock()
	// 复制所有键值对
	items := make([]struct {
		key   K
		value V
	}, 0, csm.length)

	current := csm.head.forward[0]
	for current != nil {
		items = append(items, struct {
			key   K
			value V
		}{key: current.key, value: current.value})
		current = current.forward[0]
	}
	csm.mu.RUnlock()

	// 遍历复制后的数据,不需要锁
	for _, item := range items {
		if !f(item.key, item.value) {
			break
		}
	}
}

// ToString 转换为JSON字符串
func (csm *ConcurrentSkipListMap[K, V]) ToString() string {
	csm.mu.RLock()
	defer csm.mu.RUnlock()

	m := make(map[K]V, csm.length)
	current := csm.head.forward[0]

	for current != nil {
		m[current.key] = current.value
		current = current.forward[0]
	}

	bytes, err := json.Marshal(m)
	if err != nil {
		debug.PrintStack()
		panic(err)
	}
	return string(bytes)
}

// MarshalJSON 实现 json.Marshaler 接口
func (csm *ConcurrentSkipListMap[K, V]) MarshalJSON() ([]byte, error) {
	csm.mu.RLock()
	defer csm.mu.RUnlock()

	m := make(map[K]V, csm.length)
	current := csm.head.forward[0]

	for current != nil {
		m[current.key] = current.value
		current = current.forward[0]
	}

	return json.Marshal(m)
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (csm *ConcurrentSkipListMap[K, V]) UnmarshalJSON(data []byte) error {
	csm.mu.Lock()
	defer csm.mu.Unlock()

	m := make(map[K]V)
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	// 重新初始化跳表
	csm.head = &skipListNode[K, V]{
		forward: make([]*skipListNode[K, V], maxLevel),
	}
	csm.level = 1
	csm.length = 0

	// 插入所有元素
	for k, v := range m {
		csm.putWithoutLock(k, v)
	}

	return nil
}

// MarshalBSON 实现 bson.Marshaler 接口
func (csm *ConcurrentSkipListMap[K, V]) MarshalBSON() ([]byte, error) {
	csm.mu.RLock()
	defer csm.mu.RUnlock()

	m := make(map[K]V, csm.length)
	current := csm.head.forward[0]

	for current != nil {
		m[current.key] = current.value
		current = current.forward[0]
	}

	return bson.Marshal(m)
}

// UnmarshalBSON 实现 bson.Unmarshaler 接口
func (csm *ConcurrentSkipListMap[K, V]) UnmarshalBSON(data []byte) error {
	csm.mu.Lock()
	defer csm.mu.Unlock()

	m := make(map[K]V)
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	// 重新初始化跳表
	csm.head = &skipListNode[K, V]{
		forward: make([]*skipListNode[K, V], maxLevel),
	}
	csm.level = 1
	csm.length = 0

	// 插入所有元素
	for k, v := range m {
		csm.putWithoutLock(k, v)
	}

	return nil
}

// putWithoutLock 内部使用的不加锁的put方法(用于反序列化)
func (csm *ConcurrentSkipListMap[K, V]) putWithoutLock(key K, value V) {
	update := csm.findPredecessors(key)
	current := update[0].forward[0]

	// 如果key已存在,更新值
	if current != nil && cmp.Compare(current.key, key) == 0 {
		current.value = value
		return
	}

	// 插入新节点
	newLevel := csm.randomLevel()
	if newLevel > csm.level {
		for i := csm.level; i < newLevel; i++ {
			update[i] = csm.head
		}
		csm.level = newLevel
	}

	newNode := &skipListNode[K, V]{
		key:     key,
		value:   value,
		forward: make([]*skipListNode[K, V], newLevel),
	}

	for i := 0; i < newLevel; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}

	csm.length++
}

// FirstKey 返回第一个(最小的)键
func (csm *ConcurrentSkipListMap[K, V]) FirstKey() (K, bool) {
	csm.mu.RLock()
	defer csm.mu.RUnlock()

	first := csm.head.forward[0]
	if first != nil {
		return first.key, true
	}

	var zero K
	return zero, false
}

// LastKey 返回最后一个(最大的)键
func (csm *ConcurrentSkipListMap[K, V]) LastKey() (K, bool) {
	csm.mu.RLock()
	defer csm.mu.RUnlock()

	if csm.head.forward[0] == nil {
		var zero K
		return zero, false
	}

	current := csm.head
	for i := csm.level - 1; i >= 0; i-- {
		for current.forward[i] != nil {
			current = current.forward[i]
		}
	}

	return current.key, true
}

// FirstEntry 返回第一个(最小的)键值对
func (csm *ConcurrentSkipListMap[K, V]) FirstEntry() (K, V, bool) {
	csm.mu.RLock()
	defer csm.mu.RUnlock()

	first := csm.head.forward[0]
	if first != nil {
		return first.key, first.value, true
	}

	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

// LastEntry 返回最后一个(最大的)键值对
func (csm *ConcurrentSkipListMap[K, V]) LastEntry() (K, V, bool) {
	csm.mu.RLock()
	defer csm.mu.RUnlock()

	if csm.head.forward[0] == nil {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}

	current := csm.head
	for i := csm.level - 1; i >= 0; i-- {
		for current.forward[i] != nil {
			current = current.forward[i]
		}
	}

	return current.key, current.value, true
}
