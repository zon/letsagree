package files

import (
	"encoding/json"
	"os"
	"sync"
	"testing"
)

type stubConfigWriter struct {
	mu   sync.Mutex
	data map[string][]byte
}

func (s *stubConfigWriter) WriteJSON(path string, v any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.data == nil {
		s.data = make(map[string][]byte)
	}
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	s.data[path] = data
	return nil
}

func (s *stubConfigWriter) Written(path string, v any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, ok := s.data[path]
	if !ok {
		return os.ErrNotExist
	}
	return json.Unmarshal(data, v)
}

var lastWritten struct {
	mu   sync.Mutex
	path string
	data []byte
}

type CapturingConfigWriter struct{}

func (c *CapturingConfigWriter) WriteJSON(path string, v any) error {
	lastWritten.mu.Lock()
	defer lastWritten.mu.Unlock()
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	lastWritten.path = path
	lastWritten.data = data
	return nil
}

func WrittenAt(t testing.TB, path string, v any) *PostgresConfig {
	t.Helper()
	lastWritten.mu.Lock()
	defer lastWritten.mu.Unlock()
	if lastWritten.path != path {
		t.Fatalf("WrittenAt: path mismatch: got %q, want %q", lastWritten.path, path)
	}
	if err := json.Unmarshal(lastWritten.data, v); err != nil {
		t.Fatalf("WrittenAt: %v", err)
	}
	return v.(*PostgresConfig)
}

func readJSONFile(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}