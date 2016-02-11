package store

type Store interface {
	Increment(string) int
	ExpiresIn(int, string)
}

type InMemoryStore struct {
}

func NewStore() Store {
	return &InMemoryStore{}
}

func (s *InMemoryStore) Increment(key string) int {
	return 1
}
func (s *InMemoryStore) ExpiresIn(secs int, key string) {
}
