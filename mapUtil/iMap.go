package mapUtil

type IMap[K comparable, V any] interface {
	Get(key K) V
	Put(key K, value V)
	Remove(key K)
	Size() int
	ContainsKey(key K) bool
	Clear()
	Keys() []K
	Values() []V
	PutIfAbsent(key K, value V) (existing V, loaded bool)
	GetOrDefault(key K, defaultValue V) V
	ToMap() map[K]V
	Range(f func(key K, value V) bool)
	ToString() string
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
	MarshalBSON() ([]byte, error)
	UnmarshalBSON(data []byte) error
}
