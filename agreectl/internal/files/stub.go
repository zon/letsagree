package files

import (
	"encoding/json"
	"os"
	"sync"
	"testing"

	"gopkg.in/yaml.v3"
)

func AnyHPCredentials() HPCredentials {
	return HPCredentials{
		ClientID:     "hp-client-id",
		ClientSecret: "hp-client-secret",
		PublicKey:    "hp-public-key",
	}
}

func WithHPEnv(creds HPCredentials) ConfigWriter {
	return &hpEnvConfigWriter{creds: creds}
}

type hpEnvConfigWriter struct {
	creds HPCredentials
}

func (h *hpEnvConfigWriter) WriteJSON(path string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	lastWritten.mu.Lock()
	defer lastWritten.mu.Unlock()
	lastWritten.path = path
	lastWritten.data = data
	return nil
}

func (h *hpEnvConfigWriter) WriteYAML(path string, v any) error {
	data, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	lastWritten.mu.Lock()
	defer lastWritten.mu.Unlock()
	lastWritten.path = path
	lastWritten.data = data
	return nil
}

func (h *hpEnvConfigWriter) ParseHPEnv(path string) (HPCredentials, error) {
	return h.creds, nil
}

func WrittenYAMLAt(t testing.TB, path string, v any) *HumanityProtocolConfig {
	t.Helper()
	lastWritten.mu.Lock()
	defer lastWritten.mu.Unlock()
	if lastWritten.path != path {
		t.Fatalf("WrittenYAMLAt: path mismatch: got %q, want %q", lastWritten.path, path)
	}
	if err := yaml.Unmarshal(lastWritten.data, v); err != nil {
		t.Fatalf("WrittenYAMLAt: %v", err)
	}
	return v.(*HumanityProtocolConfig)
}

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

func (c *CapturingConfigWriter) WriteYAML(path string, v any) error {
	lastWritten.mu.Lock()
	defer lastWritten.mu.Unlock()
	data, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	lastWritten.path = path
	lastWritten.data = data
	return nil
}

func (c *CapturingConfigWriter) ParseHPEnv(path string) (HPCredentials, error) {
	return ParseHPEnv(path)
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