package store

import (
	"sync"
	"time"
)

type Store interface {
	Increment(string) int
	ExpiresIn(time.Duration, string)
	CountFor(string) int
}

type InMemoryStore struct {
	storage map[string]*entry
	sync.Mutex
}

type entry struct {
	count     int
	createdAt time.Time
	expiry    time.Duration
}

func NewStore() Store {
	return &InMemoryStore{
		storage: make(map[string]*entry),
	}
}

func newEntry() *entry {
	return &entry{
		createdAt: time.Now(),
		count:     0,
	}
}

func (s *InMemoryStore) Increment(key string) int {
	v, ok := s.get(key)
	if !ok {
		v = newEntry()
	}
	v.count++
	s.set(key, v)
	return v.count
}

func (s *InMemoryStore) get(key string) (*entry, bool) {
	s.Lock()
	defer s.Unlock()
	v, ok := s.storage[key]
	return v, ok
}

func (s *InMemoryStore) set(key string, value *entry) {
	s.Lock()
	defer s.Unlock()
	s.storage[key] = value
}

func (s *InMemoryStore) ExpiresIn(expireIn time.Duration, key string) {
	v, ok := s.get(key)
	if !ok {
		v = newEntry()
	}
	v.expiry = expireIn
	s.set(key, v)
}

func (s *InMemoryStore) CountFor(key string) int {
	v, ok := s.get(key)
	if !ok {
		return 0
	}
	return v.count
}
