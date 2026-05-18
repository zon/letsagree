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

func TestSecret_StringData(t *testing.T) {
	s := &Secret{Data: map[string][]byte{"user": []byte("admin"), "password": []byte("secret123")}}
	got := s.StringData()
	want := map[string]string{"user": "admin", "password": "secret123"}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("StringData()[%q] = %q, want %q", k, got[k], v)
		}
	}
}

func TestSecret_StringDataEmpty(t *testing.T) {
	s := &Secret{Data: nil}
	if got := s.StringData(); got != nil {
		t.Errorf("StringData() = %v, want nil", got)
	}
}

func TestSecret_StringDataMissing(t *testing.T) {
	s := &Secret{Data: map[string][]byte{}}
	if got := s.StringData(); len(got) != 0 {
		t.Errorf("StringData() = %v, want empty map", got)
	}
}

func TestUpsertedSecretData(t *testing.T) {
	capturingClient := &capturingK8sClient{
		secret: &Secret{Data: map[string][]byte{"user": []byte("test")}},
	}
	capturingClient.UpsertSecret("ns", "name", map[string]string{"key": "value"})

	data := upsertedSecretData()
	if data["key"] != "value" {
		t.Errorf("upsertedSecretData() = %v, want %v", data, map[string]string{"key": "value"})
	}
}

func TestUpsertedSecretDataNotCalled(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when UpsertSecret not called")
		}
	}()
	capturedUpsertedSecretData = nil
	upsertedSecretData()
}

func TestThatFailsOnUpsert(t *testing.T) {
	client := ThatFailsOnUpsert()
	err := client.UpsertSecret("ns", "name", map[string]string{"key": "value"})
	if err == nil {
		t.Error("expected error from ThatFailsOnUpsert")
	}
}

func upsertedSecretData() map[string]string {
	if capturedUpsertedSecretData == nil {
		panic("UpsertSecret was not called")
	}
	return capturedUpsertedSecretData
}