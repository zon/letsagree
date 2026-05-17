package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteJSON_createsDirectoryIfMissing(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "backend/config/postgres.json")

	err := WriteJSON(path, PostgresConfig{Host: "localhost", Port: 30432, User: "app", Password: "secret", DBName: "app"})
	if err != nil {
		t.Fatalf("WriteJSON: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tmp, "backend/config")); os.IsNotExist(err) {
		t.Error("backend/config directory was not created")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("postgres.json was not created")
	}

	var got PostgresConfig
	if err := readJSONFile(path, &got); err != nil {
		t.Fatalf("readJSONFile: %v", err)
	}

	if got.Host != "localhost" {
		t.Errorf("Host = %q, want %q", got.Host, "localhost")
	}
	if got.Port != 30432 {
		t.Errorf("Port = %d, want %d", got.Port, 30432)
	}
	if got.User != "app" {
		t.Errorf("User = %q, want %q", got.User, "app")
	}
	if got.Password != "secret" {
		t.Errorf("Password = %q, want %q", got.Password, "secret")
	}
	if got.DBName != "app" {
		t.Errorf("DBName = %q, want %q", got.DBName, "app")
	}
}