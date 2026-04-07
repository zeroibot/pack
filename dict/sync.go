package dict

import (
	"maps"
	"sync"
)

// SyncMap is a concurrency-safe generic map
type SyncMap[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

// NewSyncMap creates a new SyncMap
func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return new(SyncMap[K, V]{data: make(map[K]V)})
}

// NewSyncMapFrom creates a new SyncMap from existing map
func NewSyncMapFrom[K comparable, V any](items map[K]V) *SyncMap[K, V] {
	return new(SyncMap[K, V]{data: items})
}

// Get returns the value for given key, and a boolean indicating if the key exists
func (sm *SyncMap[K, V]) Get(key K) (V, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	value, ok := sm.data[key]
	return value, ok
}

// Set sets the value for given key
func (sm *SyncMap[K, V]) Set(key K, value V) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data[key] = value
}

// SetIf sets the value for given key if the isValid functions passes for the current value
func (sm *SyncMap[K, V]) SetIf(key K, value V, isValid func(V) bool) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	currValue, ok := sm.data[key]
	if ok && !isValid(currValue) {
		return false
	}
	sm.data[key] = value
	return true
}

// Delete deletes the value for given key
func (sm *SyncMap[K, V]) Delete(key K) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.data, key)
}

// DeleteKeys deletes the values for given keys
func (sm *SyncMap[K, V]) DeleteKeys(keys []K) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	for _, key := range keys {
		delete(sm.data, key)
	}
}

// Clear deletes all values from the map
func (sm *SyncMap[K, V]) Clear() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data = make(map[K]V)
}

// Len returns the number of items
func (sm *SyncMap[K, V]) Len() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.data)
}

// Map returns the underlying map
func (sm *SyncMap[K, V]) Map() map[K]V {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	data := make(map[K]V)
	maps.Copy(data, sm.data)
	return data
}

// ClearMap copies the underlying map and clears the SyncMap items
func (sm *SyncMap[K, V]) ClearMap() map[K]V {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	data := make(map[K]V)
	maps.Copy(data, sm.data)
	sm.data = make(map[K]V)
	return data
}

// Keys returns the SyncMap keys in arbitrary order
func (sm *SyncMap[K, V]) Keys() []K {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return Keys(sm.data)
}

// Values returns the SyncMap values in arbitrary order
func (sm *SyncMap[K, V]) Values() []V {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return Values(sm.data)
}
