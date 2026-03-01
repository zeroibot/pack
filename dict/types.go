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

// NewCounter creates a new counter, with each item initialized to count = 0
func NewCounter[L ~[]T, T comparable](items L) Counter[T] {
	count := make(Counter[T], len(items))
	for _, item := range items {
		count[item] = 0
	}
	return count
}

// CounterFunc creates a counter, using the keys produced from keyFn
func CounterFunc[L ~[]T, T any, K comparable](items L, keyFn func(T) K) Counter[K] {
	count := make(Counter[K], len(items))
	for _, item := range items {
		count[keyFn(item)] += 1
	}
	return count
}

// UpdateCounter updates the counter in place with incoming items
func UpdateCounter[L ~[]T, T comparable](counter Counter[T], items L) {
	for _, item := range items {
		counter[item] += 1
	}
}

// NewFlags creates a new flags map, with each item initialized to given flag
func NewFlags[L ~[]T, T comparable](items L, flag bool) Flags[T] {
	flags := make(Flags[T], len(items))
	for _, item := range items {
		flags[item] = flag
	}
	return flags
}
