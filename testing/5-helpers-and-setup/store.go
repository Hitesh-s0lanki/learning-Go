// Package store is a tiny file-backed key/value store used to demonstrate
// test helpers, per-test temp dirs, cleanup, and package-level setup.
package store

import (
	"encoding/json"
	"errors"
	"os"
)

// Store persists string key/value pairs to a JSON file on disk.
type Store struct {
	path string
}

// New returns a Store backed by the file at path (created on first Set).
func New(path string) *Store {
	return &Store{path: path}
}

// Set writes (or overwrites) a key and persists the whole store.
func (s *Store) Set(key, value string) error {
	data, err := s.load()
	if err != nil {
		return err
	}
	data[key] = value
	return s.save(data)
}

// Get returns the value for key and whether it was present.
func (s *Store) Get(key string) (string, bool) {
	data, err := s.load()
	if err != nil {
		return "", false
	}
	v, ok := data[key]
	return v, ok
}

func (s *Store) load() (map[string]string, error) {
	b, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return map[string]string{}, nil // empty store until first write
	}
	if err != nil {
		return nil, err
	}
	m := map[string]string{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Store) save(data map[string]string) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0644)
}
