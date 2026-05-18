package files

import (
	"encoding/json"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type ConfigWriter interface {
	WriteJSON(path string, v any) error
	WriteYAML(path string, v any) error
	ParseHPEnv(path string) (HPCredentials, error)
}

var _ ConfigWriter = (*realConfigWriter)(nil)

type HPCredentials struct {
	ClientID     string
	ClientSecret string
	PublicKey    string
}

type HumanityProtocolConfig struct {
	ClientID     string `yaml:"clientID"`
	ClientSecret string `yaml:"clientSecret"`
	PublicKey    string `yaml:"publicKey"`
	IssuerURL    string `yaml:"issuerURL"`
	RedirectURL  string `yaml:"redirectURL"`
}

const HumanityProtocolConfigPath = "backend/config/humanity-protocol.yaml"

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

const PostgresConfigPath = "backend/config/postgres.json"

func (p PostgresConfig) ToSecretData() map[string]string {
	data, _ := json.Marshal(p)
	return map[string]string{"postgres.json": string(data)}
}

func (h HumanityProtocolConfig) ToSecretData() map[string]string {
	data, _ := yaml.Marshal(h)
	return map[string]string{"humanity-protocol.yaml": string(data)}
}

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

func WriteYAML(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func ParseHPEnv(path string) (HPCredentials, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return HPCredentials{}, err
	}
	var creds HPCredentials
	var missing []string
	for _, line := range splitLines(string(data)) {
		line = trimSpace(line)
		if line == "" || hasPrefix(line, "#") {
			continue
		}
		key, value, ok := parseLine(line)
		if !ok {
			continue
		}
		switch key {
		case "HUMANITY_CLIENT_ID":
			creds.ClientID = value
		case "HUMANITY_CLIENT_SECRET":
			creds.ClientSecret = value
		case "HUMANITY_PUBLIC_KEY":
			creds.PublicKey = value
		}
	}
	if creds.ClientID == "" {
		missing = append(missing, "HUMANITY_CLIENT_ID")
	}
	if creds.ClientSecret == "" {
		missing = append(missing, "HUMANITY_CLIENT_SECRET")
	}
	if creds.PublicKey == "" {
		missing = append(missing, "HUMANITY_PUBLIC_KEY")
	}
	if len(missing) > 0 {
		return HPCredentials{}, &missingKeyError{keys: missing}
	}
	return creds, nil
}

type missingKeyError struct {
	keys []string
}

func (m *missingKeyError) Error() string {
	return "missing keys: " + joinStrings(m.keys)
}

func splitLines(s string) []string {
	var lines []string
	for _, line := range splitString(s, "\n") {
		lines = append(lines, line)
	}
	return lines
}

func splitString(s, sep string) []string {
	if s == "" {
		return nil
	}
	var result []string
	start := 0
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}

func parseLine(line string) (key, value string, ok bool) {
	for i := 0; i < len(line); i++ {
		if line[i] == '=' {
			return line[:i], line[i+1:], true
		}
	}
	return "", "", false
}

func joinStrings(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += ", " + strs[i]
	}
	return result
}

type realConfigWriter struct{}

func NewConfigWriter() ConfigWriter {
	return &realConfigWriter{}
}

func (r *realConfigWriter) WriteJSON(path string, v any) error {
	return WriteJSON(path, v)
}

func (r *realConfigWriter) WriteYAML(path string, v any) error {
	return WriteYAML(path, v)
}

func (r *realConfigWriter) ParseHPEnv(path string) (HPCredentials, error) {
	return ParseHPEnv(path)
}