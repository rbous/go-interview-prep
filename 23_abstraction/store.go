package abstraction

import "fmt"

// Storage abstraction exercise.
//
// You're building a key-value store system where different backends
// (in-memory, logged) can be swapped in without changing the caller's code.
//
// There are several bugs preventing the tests from passing. Find and fix them.
//
// Rules:
// - Do NOT modify the test file.

// Store represents a generic key-value store.
type Store interface {
	Put(key, value string)
	Get(key string) (string, bool)
	Delete(key string)
}

// MemoryStore is an in-memory implementation of Store.
type MemoryStore struct {
	data map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (m *MemoryStore) Put(key, value string) {
	m.data[key] = value
}

func (m *MemoryStore) Get(key string) (string, bool) {
	v, ok := m.data[key]
	return v, ok
}

func (m *MemoryStore) Delete(key string) {
	delete(m.data, key)
}

// LoggedStore wraps any Store and logs all operations.
type LoggedStore struct {
	inner *MemoryStore
	log   []string
}

// NewLoggedStore creates a LoggedStore wrapping the given store.
func NewLoggedStore(s *MemoryStore) *LoggedStore {
	return &LoggedStore{inner: s}
}

func (l *LoggedStore) Put(key, value string) {
	l.log = append(l.log, fmt.Sprintf("PUT %s=%s", key, value))
	l.inner.Put(key, value)
}

func (l *LoggedStore) Get(key string) (string, bool) {
	l.log = append(l.log, fmt.Sprintf("GET %s", key))
	return l.inner.Get(key)
}

func (l *LoggedStore) Delete(key string) {
	l.log = append(l.log, fmt.Sprintf("DELETE %s", key))
}

func (l *LoggedStore) Log() []string {
	return l.log
}

// GetOrDefault retrieves a value from a Store, returning defaultVal if the key is missing.
// This function works with ANY Store thanks to the interface.
func GetOrDefault(s Store, key, defaultVal string) string {
	val, ok := s.Get(key)
	if !ok {
		return val
	}
	return val
}
