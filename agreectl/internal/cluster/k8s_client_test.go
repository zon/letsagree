package cluster

import (
	"testing"
)

func TestSecret_User(t *testing.T) {
	s := &Secret{Data: map[string][]byte{"user": []byte("admin")}}
	if got := s.User(); got != "admin" {
		t.Errorf("User() = %q, want %q", got, "admin")
	}
}

func TestSecret_UserEmpty(t *testing.T) {
	s := &Secret{Data: nil}
	if got := s.User(); got != "" {
		t.Errorf("User() = %q, want %q", got, "")
	}
}

func TestSecret_UserMissing(t *testing.T) {
	s := &Secret{Data: map[string][]byte{}}
	if got := s.User(); got != "" {
		t.Errorf("User() = %q, want %q", got, "")
	}
}

func TestSecret_Password(t *testing.T) {
	s := &Secret{Data: map[string][]byte{"password": []byte("secret123")}}
	if got := s.Password(); got != "secret123" {
		t.Errorf("Password() = %q, want %q", got, "secret123")
	}
}

func TestSecret_PasswordEmpty(t *testing.T) {
	s := &Secret{Data: nil}
	if got := s.Password(); got != "" {
		t.Errorf("Password() = %q, want %q", got, "")
	}
}

func TestSecret_PasswordMissing(t *testing.T) {
	s := &Secret{Data: map[string][]byte{}}
	if got := s.Password(); got != "" {
		t.Errorf("Password() = %q, want %q", got, "")
	}
}

func TestSecret_DBName(t *testing.T) {
	s := &Secret{Data: map[string][]byte{"dbname": []byte("mydb")}}
	if got := s.DBName(); got != "mydb" {
		t.Errorf("DBName() = %q, want %q", got, "mydb")
	}
}

func TestSecret_DBNameEmpty(t *testing.T) {
	s := &Secret{Data: nil}
	if got := s.DBName(); got != "" {
		t.Errorf("DBName() = %q, want %q", got, "")
	}
}

func TestSecret_DBNameMissing(t *testing.T) {
	s := &Secret{Data: map[string][]byte{}}
	if got := s.DBName(); got != "" {
		t.Errorf("DBName() = %q, want %q", got, "")
	}
}