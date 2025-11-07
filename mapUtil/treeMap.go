package mapUtil

import (
	"encoding/json"
	"runtime/debug"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	RED   = true  // 红黑树红色节点标识
	BLACK = false // 红黑树黑色节点标识
)

var _ IMap[string, int] = (*TreeMap[string, int])(nil)

// TreeMap 是一个基于红黑树实现的并发安全的有序映射
type TreeMap[K comparable, V any] struct {
	mu     sync.RWMutex      // 读写锁保证并发安全
	root   *node[K, V]       // 红黑树根节点
	size   int               // 映射中元素数量
	less   func(a, b K) bool // 键比较函数
	emptyK K                 // 键类型零值
	emptyV V                 // 值类型零值
}

// node 是红黑树的节点结构
type node[K comparable, V any] struct {
	key    K           // 节点键
	value  V           // 节点值
	left   *node[K, V] // 左子节点
	right  *node[K, V] // 右子节点
	parent *node[K, V] // 父节点
	color  bool        // 节点颜色(RED/BLACK)
}

// NewTreeMap 创建一个新的TreeMap实例
// less: 用于比较键的函数，确定键的顺序
func NewTreeMap[K comparable, V any](less func(a, b K) bool) *TreeMap[K, V] {
	return &TreeMap[K, V]{
		less: less,
	}
}

// Get 获取指定键对应的值
// key: 要查找的键
// 返回值: 找到的值和是否存在的布尔值
func (tm *TreeMap[K, V]) Get(key K) V {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	n := tm.getNode(key)
	if n == nil {
		return tm.emptyV
	}
	return n.value
}

// getNode 内部方法，根据键查找节点
func (tm *TreeMap[K, V]) getNode(key K) *node[K, V] {
	current := tm.root
	for current != nil {
		if tm.less(key, current.key) {
			current = current.left
		} else if tm.less(current.key, key) {
			current = current.right
		} else {
			return current
		}
	}
	return nil
}

// Put 向映射中添加或更新键值对
// key: 要添加的键
// value: 要添加的值
func (tm *TreeMap[K, V]) Put(key K, value V) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.put(key, value)
}

func (tm *TreeMap[K, V]) put(key K, value V) {
	if tm.root == nil {
		tm.root = &node[K, V]{
			key:   key,
			value: value,
			color: BLACK,
		}
		tm.size++
		return
	}

	var parent *node[K, V]
	current := tm.root
	for current != nil {
		parent = current
		if tm.less(key, current.key) {
			current = current.left
		} else if tm.less(current.key, key) {
			current = current.right
		} else {
			// 键已存在，更新值
			current.value = value
			return
		}
	}

	newNode := &node[K, V]{
		key:    key,
		value:  value,
		parent: parent,
		color:  RED,
	}

	if tm.less(key, parent.key) {
		parent.left = newNode
	} else {
		parent.right = newNode
	}

	tm.fixAfterInsertion(newNode)
	tm.size++
}

// Remove 删除指定键的键值对
// key: 要删除的键
func (tm *TreeMap[K, V]) Remove(key K) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	n := tm.getNode(key)
	if n == nil {
		return
	}

	tm.deleteNode(n)
	tm.size--
}

// deleteNode 内部方法，删除指定节点
func (tm *TreeMap[K, V]) deleteNode(p *node[K, V]) {
	// 如果节点有两个子节点，用后继节点替换
	if p.left != nil && p.right != nil {
		s := tm.successor(p)
		p.key = s.key
		p.value = s.value
		p = s
	}

	var replacement *node[K, V]
	if p.left != nil {
		replacement = p.left
	} else {
		replacement = p.right
	}

	if replacement != nil {
		replacement.parent = p.parent
		if p.parent == nil {
			tm.root = replacement
		} else if p == p.parent.left {
			p.parent.left = replacement
		} else {
			p.parent.right = replacement
		}

		p.left = nil
		p.right = nil
		p.parent = nil

		if p.color == BLACK {
			tm.fixAfterDeletion(replacement)
		}
	} else if p.parent == nil {
		tm.root = nil
	} else {
		if p.color == BLACK {
			tm.fixAfterDeletion(p)
		}

		if p.parent != nil {
			if p == p.parent.left {
				p.parent.left = nil
			} else if p == p.parent.right {
				p.parent.right = nil
			}
			p.parent = nil
		}
	}
}

// Size 返回映射中元素的数量
func (tm *TreeMap[K, V]) Size() int {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.size
}

// ContainsKey 检查映射中是否包含指定键
// key: 要检查的键
// 返回值: 是否包含的布尔值
func (tm *TreeMap[K, V]) ContainsKey(key K) bool {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.getNode(key) != nil
}

// Clear 清空映射
func (tm *TreeMap[K, V]) Clear() {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.root = nil
	tm.size = 0
}

// Keys 返回映射中所有键的切片，按升序排列
func (tm *TreeMap[K, V]) Keys() []K {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	keys := make([]K, 0, tm.size)
	tm.inOrder(tm.root, func(n *node[K, V]) {
		keys = append(keys, n.key)
	})
	return keys
}

// Values 返回映射中所有值的切片，按键的升序排列
func (tm *TreeMap[K, V]) Values() []V {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	values := make([]V, 0, tm.size)
	tm.inOrder(tm.root, func(n *node[K, V]) {
		values = append(values, n.value)
	})
	return values
}

// inOrder 内部方法，中序遍历树
func (tm *TreeMap[K, V]) inOrder(n *node[K, V], f func(*node[K, V])) {
	if n == nil {
		return
	}
	tm.inOrder(n.left, f)
	f(n)
	tm.inOrder(n.right, f)
}

// PutIfAbsent 只有当键不存在时才放入值
// key: 要放入的键
// value: 要放入的值
// 返回值: 已存在的值(如果有)和是否已存在的布尔值
func (tm *TreeMap[K, V]) PutIfAbsent(key K, value V) (existing V, loaded bool) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	n := tm.getNode(key)
	if n != nil {
		return n.value, true
	}

	if tm.root == nil {
		tm.root = &node[K, V]{
			key:   key,
			value: value,
			color: BLACK,
		}
		tm.size++
		return tm.emptyV, false
	}

	var parent *node[K, V]
	current := tm.root
	for current != nil {
		parent = current
		if tm.less(key, current.key) {
			current = current.left
		} else if tm.less(current.key, key) {
			current = current.right
		}
	}

	newNode := &node[K, V]{
		key:    key,
		value:  value,
		parent: parent,
		color:  RED,
	}

	if tm.less(key, parent.key) {
		parent.left = newNode
	} else {
		parent.right = newNode
	}

	tm.fixAfterInsertion(newNode)
	tm.size++
	return tm.emptyV, false
}

// GetOrDefault 获取指定键对应的值，如果不存在则返回默认值
// key: 要查找的键
// defaultValue: 默认值
// 返回值: 找到的值或默认值
func (tm *TreeMap[K, V]) GetOrDefault(key K, defaultValue V) V {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	n := tm.getNode(key)
	if n == nil {
		return defaultValue
	}
	return n.value
}

// successor 内部方法，查找指定节点的后继节点
func (tm *TreeMap[K, V]) successor(t *node[K, V]) *node[K, V] {
	if t == nil {
		return nil
	} else if t.right != nil {
		p := t.right
		for p.left != nil {
			p = p.left
		}
		return p
	} else {
		p := t.parent
		ch := t
		for p != nil && ch == p.right {
			ch = p
			p = p.parent
		}
		return p
	}
}

// fixAfterInsertion 内部方法，插入节点后调整红黑树平衡
func (tm *TreeMap[K, V]) fixAfterInsertion(x *node[K, V]) {
	x.color = RED

	for x != nil && x != tm.root && x.parent.color == RED {
		if x.parent == x.parent.parent.left {
			y := x.parent.parent.right
			if y != nil && y.color == RED {
				x.parent.color = BLACK
				y.color = BLACK
				x.parent.parent.color = RED
				x = x.parent.parent
			} else {
				if x == x.parent.right {
					x = x.parent
					tm.rotateLeft(x)
				}
				x.parent.color = BLACK
				x.parent.parent.color = RED
				tm.rotateRight(x.parent.parent)
			}
		} else {
			y := x.parent.parent.left
			if y != nil && y.color == RED {
				x.parent.color = BLACK
				y.color = BLACK
				x.parent.parent.color = RED
				x = x.parent.parent
			} else {
				if x == x.parent.left {
					x = x.parent
					tm.rotateRight(x)
				}
				x.parent.color = BLACK
				x.parent.parent.color = RED
				tm.rotateLeft(x.parent.parent)
			}
		}
	}
	tm.root.color = BLACK
}

// fixAfterDeletion 内部方法，删除节点后调整红黑树平衡
func (tm *TreeMap[K, V]) fixAfterDeletion(x *node[K, V]) {
	for x != tm.root && colorOf(x) == BLACK {
		if x == leftOf(parentOf(x)) {
			sib := rightOf(parentOf(x))

			if colorOf(sib) == RED {
				setColor(sib, BLACK)
				setColor(parentOf(x), RED)
				tm.rotateLeft(parentOf(x))
				sib = rightOf(parentOf(x))
			}

			if colorOf(leftOf(sib)) == BLACK && colorOf(rightOf(sib)) == BLACK {
				setColor(sib, RED)
				x = parentOf(x)
			} else {
				if colorOf(rightOf(sib)) == BLACK {
					setColor(leftOf(sib), BLACK)
					setColor(sib, RED)
					tm.rotateRight(sib)
					sib = rightOf(parentOf(x))
				}
				setColor(sib, colorOf(parentOf(x)))
				setColor(parentOf(x), BLACK)
				setColor(rightOf(sib), BLACK)
				tm.rotateLeft(parentOf(x))
				x = tm.root
			}
		} else {
			sib := leftOf(parentOf(x))

			if colorOf(sib) == RED {
				setColor(sib, BLACK)
				setColor(parentOf(x), RED)
				tm.rotateRight(parentOf(x))
				sib = leftOf(parentOf(x))
			}

			if colorOf(rightOf(sib)) == BLACK && colorOf(leftOf(sib)) == BLACK {
				setColor(sib, RED)
				x = parentOf(x)
			} else {
				if colorOf(leftOf(sib)) == BLACK {
					setColor(rightOf(sib), BLACK)
					setColor(sib, RED)
					tm.rotateLeft(sib)
					sib = leftOf(parentOf(x))
				}
				setColor(sib, colorOf(parentOf(x)))
				setColor(parentOf(x), BLACK)
				setColor(leftOf(sib), BLACK)
				tm.rotateRight(parentOf(x))
				x = tm.root
			}
		}
	}
	setColor(x, BLACK)
}

// rotateLeft 内部方法，左旋转
func (tm *TreeMap[K, V]) rotateLeft(p *node[K, V]) {
	if p != nil {
		r := p.right
		p.right = r.left
		if r.left != nil {
			r.left.parent = p
		}
		r.parent = p.parent
		if p.parent == nil {
			tm.root = r
		} else if p.parent.left == p {
			p.parent.left = r
		} else {
			p.parent.right = r
		}
		r.left = p
		p.parent = r
	}
}

// rotateRight 内部方法，右旋转
func (tm *TreeMap[K, V]) rotateRight(p *node[K, V]) {
	if p != nil {
		l := p.left
		p.left = l.right
		if l.right != nil {
			l.right.parent = p
		}
		l.parent = p.parent
		if p.parent == nil {
			tm.root = l
		} else if p.parent.right == p {
			p.parent.right = l
		} else {
			p.parent.left = l
		}
		l.right = p
		p.parent = l
	}
}

// colorOf 内部方法，获取节点颜色
func colorOf[K comparable, V any](n *node[K, V]) bool {
	if n == nil {
		return BLACK
	}
	return n.color
}

// setColor 内部方法，设置节点颜色
func setColor[K comparable, V any](n *node[K, V], color bool) {
	if n != nil {
		n.color = color
	}
}

// parentOf 内部方法，获取父节点
func parentOf[K comparable, V any](n *node[K, V]) *node[K, V] {
	if n == nil {
		return nil
	}
	return n.parent
}

// leftOf 内部方法，获取左子节点
func leftOf[K comparable, V any](n *node[K, V]) *node[K, V] {
	if n == nil {
		return nil
	}
	return n.left
}

// rightOf 内部方法，获取右子节点
func rightOf[K comparable, V any](n *node[K, V]) *node[K, V] {
	if n == nil {
		return nil
	}
	return n.right
}

// FirstKey 返回映射中第一个键
func (tm *TreeMap[K, V]) FirstKey() (K, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	if tm.root == nil {
		return tm.emptyK, false
	}
	n := tm.firstNode(tm.root)
	return n.key, true
}

// firstNode 内部方法，查找最小节点
func (tm *TreeMap[K, V]) firstNode(n *node[K, V]) *node[K, V] {
	for n.left != nil {
		n = n.left
	}
	return n
}

// LastKey 返回映射中最后一个键
func (tm *TreeMap[K, V]) LastKey() (K, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	if tm.root == nil {
		return tm.emptyK, false
	}
	n := tm.lastNode(tm.root)
	return n.key, true
}

// lastNode 内部方法，查找最大节点
func (tm *TreeMap[K, V]) lastNode(n *node[K, V]) *node[K, V] {
	for n.right != nil {
		n = n.right
	}
	return n
}

// Range 遍历映射中的键值对
// f: 遍历函数，返回false可提前终止遍历
func (tm *TreeMap[K, V]) Range(f func(key K, value V) bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	tm.inOrderFunc(tm.root, func(n *node[K, V]) bool {
		return f(n.key, n.value)
	})
}

// inOrderFunc 内部方法，带中断功能的中序遍历
func (tm *TreeMap[K, V]) inOrderFunc(n *node[K, V], f func(*node[K, V]) bool) bool {
	if n == nil {
		return true
	}
	if !tm.inOrderFunc(n.left, f) {
		return false
	}
	if !f(n) {
		return false
	}
	return tm.inOrderFunc(n.right, f)
}

// ToString 将映射转换为JSON格式字符串
func (tm *TreeMap[K, V]) ToString() string {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	bytes, err := json.Marshal(tm.ToMap())
	if err != nil {
		debug.PrintStack()
		panic(err)
	}
	return string(bytes)
}

// ToMap 内部方法，将树转换为map
func (tm *TreeMap[K, V]) ToMap() map[K]V {
	m := make(map[K]V, tm.size)
	tm.inOrder(tm.root, func(n *node[K, V]) {
		m[n.key] = n.value
	})
	return m
}

// MarshalJSON 实现JSON序列化接口
func (tm *TreeMap[K, V]) MarshalJSON() ([]byte, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return json.Marshal(tm.ToMap())
}

// UnmarshalJSON 实现JSON反序列化接口
func (tm *TreeMap[K, V]) UnmarshalJSON(data []byte) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	var m map[K]V
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	tm.root = nil
	tm.size = 0
	for k, v := range m {
		tm.put(k, v)
	}
	return nil
}

// MarshalBSON 实现BSON序列化接口
func (tm *TreeMap[K, V]) MarshalBSON() ([]byte, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return bson.Marshal(tm.ToMap())
}

// UnmarshalBSON 实现BSON反序列化接口
func (tm *TreeMap[K, V]) UnmarshalBSON(data []byte) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	var m map[K]V
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	tm.root = nil
	tm.size = 0
	for k, v := range m {
		tm.put(k, v)
	}
	return nil
}
