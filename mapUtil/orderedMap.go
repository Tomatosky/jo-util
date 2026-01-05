package mapUtil

import (
	"encoding/json"
	"fmt"
	"iter"

	"github.com/Tomatosky/jo-util/logger"
	"go.mongodb.org/mongo-driver/bson"
)

var _ IMap[string, int] = (*OrderedMap[string, int])(nil)

type OrderedMap[K comparable, V any] struct {
	kv     map[K]*Element[K, V]
	ll     list[K, V]
	emptyV V
}

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		kv: make(map[K]*Element[K, V]),
	}
}

// NewOrderedMapWithCapacity creates a map with enough pre-allocated space to
// hold the specified number of elements.
func NewOrderedMapWithCapacity[K comparable, V any](capacity int) *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		kv: make(map[K]*Element[K, V], capacity),
	}
}

func NewOrderedMapWithElements[K comparable, V any](els ...*Element[K, V]) *OrderedMap[K, V] {
	om := NewOrderedMapWithCapacity[K, V](len(els))
	for _, el := range els {
		om.Put(el.Key, el.Value)
	}
	return om
}

// Get returns the value for a key
func (m *OrderedMap[K, V]) Get(key K) V {
	v, ok := m.kv[key]
	if ok {
		return v.Value
	}
	return m.emptyV
}

// Put will set (or replace) a value for a key
func (m *OrderedMap[K, V]) Put(key K, value V) {
	_, alreadyExist := m.kv[key]
	if alreadyExist {
		m.kv[key].Value = value
		return
	}

	element := m.ll.PushBack(key, value)
	m.kv[key] = element
	return
}

// ReplaceKey replaces an existing key with a new key while preserving order of
// the value. This function will return true if the operation was successful, or
// false if 'originalKey' is not found OR 'newKey' already exists (which would be an overwrite).
func (m *OrderedMap[K, V]) ReplaceKey(originalKey, newKey K) bool {
	element, originalExists := m.kv[originalKey]
	_, newKeyExists := m.kv[newKey]
	if originalExists && !newKeyExists {
		delete(m.kv, originalKey)
		m.kv[newKey] = element
		element.Key = newKey
		return true
	}
	return false
}

func (m *OrderedMap[K, V]) Clear() {
	m.kv = make(map[K]*Element[K, V])
	m.ll = list[K, V]{}
}

func (m *OrderedMap[K, V]) PutIfAbsent(key K, value V) (existing V, loaded bool) {
	el, ok := m.kv[key]
	if ok {
		return el.Value, true
	}
	element := m.ll.PushBack(key, value)
	m.kv[key] = element
	return value, false
}

// GetOrDefault returns the value for a key. If the key does not exist, returns
// the default value instead.
func (m *OrderedMap[K, V]) GetOrDefault(key K, defaultValue V) V {
	if value, ok := m.kv[key]; ok {
		return value.Value
	}

	return defaultValue
}

// GetElement returns the element for a key. If the key does not exist, the
// pointer will be nil.
func (m *OrderedMap[K, V]) GetElement(key K) *Element[K, V] {
	element, ok := m.kv[key]
	if ok {
		return element
	}

	return nil
}

// Size returns the number of elements in the map.
func (m *OrderedMap[K, V]) Size() int {
	return len(m.kv)
}

// AllFromFront returns an iterator that yields all elements in the map starting
// at the front (oldest Set element).
func (m *OrderedMap[K, V]) AllFromFront() iter.Seq2[K, V] {
	return func(yield func(key K, value V) bool) {
		for el := m.Front(); el != nil; el = el.Next() {
			if !yield(el.Key, el.Value) {
				return
			}
		}
	}
}

// AllFromBack returns an iterator that yields all elements in the map starting
// at the back (most recent Set element).
func (m *OrderedMap[K, V]) AllFromBack() iter.Seq2[K, V] {
	return func(yield func(key K, value V) bool) {
		for el := m.Back(); el != nil; el = el.Prev() {
			if !yield(el.Key, el.Value) {
				return
			}
		}
	}
}

func (m *OrderedMap[K, V]) Keys() []K {
	keys := make([]K, 0, m.Size())
	for el := m.Front(); el != nil; el = el.Next() {
		keys = append(keys, el.Key)
	}
	return keys
}

func (m *OrderedMap[K, V]) Values() []V {
	values := make([]V, 0, m.Size())
	for el := m.Front(); el != nil; el = el.Next() {
		values = append(values, el.Value)
	}
	return values
}

// Remove will remove a key from the map
func (m *OrderedMap[K, V]) Remove(key K) {
	element, ok := m.kv[key]
	if ok {
		m.ll.Remove(element)
		delete(m.kv, key)
	}
}

// Front will return the element that is the first (oldest Set element). If
// there are no elements this will return nil.
func (m *OrderedMap[K, V]) Front() *Element[K, V] {
	return m.ll.Front()
}

// Back will return the element that is the last (most recent Set element). If
// there are no elements this will return nil.
func (m *OrderedMap[K, V]) Back() *Element[K, V] {
	return m.ll.Back()
}

// Copy returns a new OrderedMap with the same elements.
// Using Copy while there are concurrent writes may mangle the result.
func (m *OrderedMap[K, V]) Copy() *OrderedMap[K, V] {
	m2 := NewOrderedMapWithCapacity[K, V](m.Size())
	for el := m.Front(); el != nil; el = el.Next() {
		m2.Put(el.Key, el.Value)
	}
	return m2
}

// ContainsKey checks if a key exists in the map.
func (m *OrderedMap[K, V]) ContainsKey(key K) bool {
	_, exists := m.kv[key]
	return exists
}

func (m *OrderedMap[K, V]) ToMap() map[K]V {
	m2 := make(map[K]V, m.Size())
	for el := m.Front(); el != nil; el = el.Next() {
		m2[el.Key] = el.Value
	}
	return m2
}

func (m *OrderedMap[K, V]) Range(f func(key K, value V) bool) {
	for el := m.Front(); el != nil; el = el.Next() {
		if !f(el.Key, el.Value) {
			return
		}
	}
}

func (m *OrderedMap[K, V]) ToString() string {
	toMap := m.ToMap()
	bytes, err := json.Marshal(toMap)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("%v", err))
		panic(err)
	}
	return string(bytes)
}

func (m *OrderedMap[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.ToMap())
}

func (m *OrderedMap[K, V]) UnmarshalJSON(data []byte) error {
	var toMap map[K]V
	err := json.Unmarshal(data, &toMap)
	if err != nil {
		return err
	}
	m.kv = make(map[K]*Element[K, V])
	for k, v := range toMap {
		m.Put(k, v)
	}
	return nil
}

func (m *OrderedMap[K, V]) MarshalBSON() ([]byte, error) {
	return bson.Marshal(m.ToMap())
}

func (m *OrderedMap[K, V]) UnmarshalBSON(bytes []byte) error {
	var toMap map[K]V
	err := bson.Unmarshal(bytes, &toMap)
	if err != nil {
		return err
	}
	m.kv = make(map[K]*Element[K, V])
	for k, v := range toMap {
		m.Put(k, v)
	}
	return nil
}
