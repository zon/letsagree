package files

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ConfigWriter interface {
	WriteJSON(path string, v any) error
}

var _ ConfigWriter = (*realConfigWriter)(nil)

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

const PostgresConfigPath = "backend/config/postgres.json"

func WriteJSON(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

type realConfigWriter struct{}

func NewConfigWriter() ConfigWriter {
	return &realConfigWriter{}
}

func (r *realConfigWriter) WriteJSON(path string, v any) error {
	return WriteJSON(path, v)
}