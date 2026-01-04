package store

import "sync"

type Store[K comparable, V any] struct {
	mu sync.RWMutex
	data map[K]V
}

func NewStore[K comparable, V any]() *Store[K, V] {
	return &Store[K, V]{
		data:  make(map[K]V),
	}
}

func (s *Store[K, V]) Set(k K, v V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[k] = v
}
func (s *Store[K, V]) Get(k K) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, ok := s.data[k]
	return v, ok
}

func (s *Store[K, V]) Delete(k K) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[k]
	if ok {
		delete(s.data, k)
	}

	return ok
}
func (s *Store[K, V]) Snapshot() map[K]V {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cp := make(map[K]V, len(s.data))
	for k, v := range s.data {
		cp[k] = v
	}

	return cp
}

func (s *Store[K, V]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}
