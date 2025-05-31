package store

import (
	"sync"
)

var (
	data = make(map[string]string)
	mu   sync.RWMutex
)

// Set stores a key-value pair
func Set(key, value string) {
	mu.Lock()
	defer mu.Unlock()
	data[key] = value
}

// Get retrieves a value by key
func Get(key string) (string, bool) {
	mu.RLock()
	defer mu.RUnlock()
	val, ok := data[key]
	return val, ok
}

// Delete removes a key
func Delete(key string) {
	mu.Lock()
	defer mu.Unlock()
	delete(data, key)
}
