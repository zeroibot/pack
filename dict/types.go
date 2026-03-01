package dict

type (
	Object                = map[string]any
	IntMap                = map[string]int
	UintMap               = map[string]uint
	BoolMap               = map[string]bool
	StringMap             = map[string]string
	StringListMap         = map[string][]string
	StringCounter         = map[string]int
	IntCounter            = map[int]int
	UintCounter           = map[uint]int
	Counter[T comparable] = map[T]int
	Flags[T comparable]   = map[T]bool
)

// Entry represents a Key-Value pair
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// Tuple returns the Key, Value of Entry
func (e Entry[K, V]) Tuple() (K, V) {
	return e.Key, e.Value
}

// NewCounter creates a new counter, with each item initialized to count = 0
func NewCounter[T comparable](items []T) Counter[T] {
	count := make(Counter[T], len(items))
	for _, item := range items {
		count[item] = 0
	}
	return count
}

// CounterFunc creates a counter, using the keys produced from keyFn
func CounterFunc[T any, K comparable](items []T, keyFn func(T) K) Counter[K] {
	count := make(Counter[K], len(items))
	for _, item := range items {
		count[keyFn(item)] += 1
	}
	return count
}

// UpdateCounter updates the counter in place with incoming items
func UpdateCounter[T comparable](counter Counter[T], items []T) {
	for _, item := range items {
		counter[item] += 1
	}
}

// NewFlags creates a new flags map, with each item initialized to given flag
func NewFlags[T comparable](items []T, flag bool) Flags[T] {
	flags := make(Flags[T], len(items))
	for _, item := range items {
		flags[item] = flag
	}
	return flags
}
