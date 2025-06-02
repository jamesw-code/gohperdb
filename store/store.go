package store

import (
	"sync"
	"time"
)

type entry struct {
	value     string
	expiresAt time.Time
}

type InMemoryStore struct {
	data map[string]entry
	mu   sync.RWMutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		data: make(map[string]entry),
	}
}

// Set stores a key-value pair
func (s *InMemoryStore) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = entry{value: value}
}

// Get retrieves a value by key
func (s *InMemoryStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val.value, ok
}

// Delete removes a key
func (s *InMemoryStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}
func (s *InMemoryStore) SetEX(key, value string, ttlSeconds int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	exp := time.Now().Add(time.Duration(ttlSeconds) * time.Second)
	s.data[key] = entry{
		value:     value,
		expiresAt: exp,
	}
}

func (s *InMemoryStore) CleanupExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for k, v := range s.data {
		if !v.expiresAt.IsZero() && now.After(v.expiresAt) {
			delete(s.data, k)
		}
	}
}
