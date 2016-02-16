package store

import (
	"fmt"
	"sync"
	"time"
)

// store represents a key/value storage with counter and expiry of keys
type Store interface {
	Increment(string) int
	ExpiresIn(time.Duration, string)
	CountFor(string) int
	Current() map[string]int
}

// in-memory store is a simple map-backed key store
type InMemoryStore struct {
	storage map[string]*Entry
	sync.RWMutex
}

type Entry struct {
	count      int
	Expirable  bool
	ExpiryTime time.Time
}

func (e *Entry) Expired() bool {
	if e.Expirable && time.Now().After(e.ExpiryTime) {
		return true
	}
	return false
}

func NewStore() Store {
	store := &InMemoryStore{
		storage: make(map[string]*Entry),
	}
	store.expiryCycle()

	return store
}

func NewEntry() *Entry {
	return &Entry{}
}

func (s *InMemoryStore) Increment(key string) int {
	v, ok := s.get(key)
	if !ok {
		v = NewEntry()
	}
	v.count++
	s.set(key, v)
	return v.count
}

func (s *InMemoryStore) get(key string) (*Entry, bool) {
	s.RLock()
	defer s.RUnlock()
	v, ok := s.storage[key]
	return v, ok
}

func (s *InMemoryStore) set(key string, value *Entry) {
	s.Lock()
	defer s.Unlock()
	s.storage[key] = value
}

func (s *InMemoryStore) ExpiresIn(expireIn time.Duration, key string) {
	v, ok := s.get(key)
	if !ok {
		v = NewEntry()
	}
	v.Expirable = true
	v.ExpiryTime = time.Now().Add(expireIn)
	s.set(key, v)
}

func (s *InMemoryStore) expiryCycle() {
	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for _ = range ticker.C {
			s.Lock()
			for k, v := range s.storage {
				if v.Expired() {
					fmt.Printf("removing expired key [%s]\n", k)
					delete(s.storage, k)
				}
			}
			s.Unlock()
		}
	}()
}

func (s *InMemoryStore) CountFor(key string) int {
	v, ok := s.get(key)
	if !ok {
		return 0
	}
	return v.count
}

func (s *InMemoryStore) Current() map[string]int {
	m := make(map[string]int)
	s.Lock()
	for k, v := range s.storage {
		m[k] = v.count
	}
	s.Unlock()
	return m
}
