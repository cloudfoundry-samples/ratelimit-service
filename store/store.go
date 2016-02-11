package store

import "sync"

type Store interface {
	Increment(string) int
	ExpiresIn(int, string)
}

type InMemoryStore struct {
	storage map[string]int
	sync.Mutex
}

func NewStore() Store {
	return &InMemoryStore{
		storage: make(map[string]int),
	}
}

func (s *InMemoryStore) Increment(key string) int {
	s.Lock()
	defer s.Unlock()
	i, ok := s.storage[key]
	if !ok {
		i = 0
	}
	i++
	s.storage[key] = i
	return i
}

func (s *InMemoryStore) ExpiresIn(secs int, key string) {
}
