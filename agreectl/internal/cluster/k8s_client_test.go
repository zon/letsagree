package cluster

import (
	"testing"
)

func TestSecret_User(t *testing.T) {
	s := &Secret{Bytes: map[string][]byte{"user": []byte("admin")}}
	if got := s.User(); got != "admin" {
		t.Errorf("User() = %q, want %q", got, "admin")
	}
}

func TestSecret_UserEmpty(t *testing.T) {
	s := &Secret{Bytes: nil}
	if got := s.User(); got != "" {
		t.Errorf("User() = %q, want %q", got, "")
	}
}

func TestSecret_UserMissing(t *testing.T) {
	s := &Secret{Bytes: map[string][]byte{}}
	if got := s.User(); got != "" {
		t.Errorf("User() = %q, want %q", got, "")
	}
}

func TestSecret_Password(t *testing.T) {
	s := &Secret{Bytes: map[string][]byte{"password": []byte("secret123")}}
	if got := s.Password(); got != "secret123" {
		t.Errorf("Password() = %q, want %q", got, "secret123")
	}
}

func TestSecret_PasswordEmpty(t *testing.T) {
	s := &Secret{Bytes: nil}
	if got := s.Password(); got != "" {
		t.Errorf("Password() = %q, want %q", got, "")
	}
}

func TestSecret_PasswordMissing(t *testing.T) {
	s := &Secret{Bytes: map[string][]byte{}}
	if got := s.Password(); got != "" {
		t.Errorf("Password() = %q, want %q", got, "")
	}
}

func TestSecret_DBName(t *testing.T) {
	s := &Secret{Bytes: map[string][]byte{"dbname": []byte("mydb")}}
	if got := s.DBName(); got != "mydb" {
		t.Errorf("DBName() = %q, want %q", got, "mydb")
	}
}

func TestSecret_DBNameEmpty(t *testing.T) {
	s := &Secret{Bytes: nil}
	if got := s.DBName(); got != "" {
		t.Errorf("DBName() = %q, want %q", got, "")
	}
}

func TestSecret_DBNameMissing(t *testing.T) {
	s := &Secret{Bytes: map[string][]byte{}}
	if got := s.DBName(); got != "" {
		t.Errorf("DBName() = %q, want %q", got, "")
	}
}

func TestSecret_ClientID(t *testing.T) {
	s := &Secret{Bytes: map[string][]byte{"clientId": []byte("my-client-id")}}
	if got := s.ClientID(); got != "my-client-id" {
		t.Errorf("ClientID() = %q, want %q", got, "my-client-id")
	}
}

func TestSecret_ClientIDEmpty(t *testing.T) {
	s := &Secret{Bytes: nil}
	if got := s.ClientID(); got != "" {
		t.Errorf("ClientID() = %q, want %q", got, "")
	}
}

func TestSecret_ClientIDMissing(t *testing.T) {
	s := &Secret{Bytes: map[string][]byte{}}
	if got := s.ClientID(); got != "" {
		t.Errorf("ClientID() = %q, want %q", got, "")
	}
}

func TestSecret_ClientSecret(t *testing.T) {
	s := &Secret{Bytes: map[string][]byte{"clientSecret": []byte("my-client-secret")}}
	if got := s.ClientSecret(); got != "my-client-secret" {
		t.Errorf("ClientSecret() = %q, want %q", got, "my-client-secret")
	}
}

func TestSecret_ClientSecretEmpty(t *testing.T) {
	s := &Secret{Bytes: nil}
	if got := s.ClientSecret(); got != "" {
		t.Errorf("ClientSecret() = %q, want %q", got, "")
	}
}

func TestSecret_ClientSecretMissing(t *testing.T) {
	s := &Secret{Bytes: map[string][]byte{}}
	if got := s.ClientSecret(); got != "" {
		t.Errorf("ClientSecret() = %q, want %q", got, "")
	}
}

func TestSecret_PublicKey(t *testing.T) {
	s := &Secret{Bytes: map[string][]byte{"publicKey": []byte("my-public-key")}}
	if got := s.PublicKey(); got != "my-public-key" {
		t.Errorf("PublicKey() = %q, want %q", got, "my-public-key")
	}
}

func TestSecret_PublicKeyEmpty(t *testing.T) {
	s := &Secret{Bytes: nil}
	if got := s.PublicKey(); got != "" {
		t.Errorf("PublicKey() = %q, want %q", got, "")
	}
}

func TestSecret_PublicKeyMissing(t *testing.T) {
	s := &Secret{Bytes: map[string][]byte{}}
	if got := s.PublicKey(); got != "" {
		t.Errorf("PublicKey() = %q, want %q", got, "")
	}
}

func TestSecret_StringData(t *testing.T) {
	s := &Secret{Bytes: map[string][]byte{"user": []byte("admin"), "password": []byte("secret123")}}
	got := s.StringData()
	want := map[string]string{"user": "admin", "password": "secret123"}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("StringData()[%q] = %q, want %q", k, got[k], v)
		}
	}
}

func TestSecret_StringDataEmpty(t *testing.T) {
	s := &Secret{Bytes: nil}
	if got := s.StringData(); got != nil {
		t.Errorf("StringData() = %v, want nil", got)
	}
}

func TestSecret_StringDataMissing(t *testing.T) {
	s := &Secret{Bytes: map[string][]byte{}}
	if got := s.StringData(); len(got) != 0 {
		t.Errorf("StringData() = %v, want empty map", got)
	}
}

func TestUpsertedSecretBytes(t *testing.T) {
	capturingClient := &capturingK8sClient{
		secret: &Secret{Bytes: map[string][]byte{"user": []byte("test")}},
	}
	capturingClient.UpsertSecret("ns", "name", map[string]string{"key": "value"})

	data := upsertedSecretBytes()
	if data["key"] != "value" {
		t.Errorf("upsertedSecretBytes() = %v, want %v", data, map[string]string{"key": "value"})
	}
}

func TestUpsertedSecretBytesNotCalled(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when UpsertSecret not called")
		}
	}()
	capturedUpsertedSecretData = nil
	upsertedSecretBytes()
}

func TestThatFailsOnUpsert(t *testing.T) {
	client := ThatFailsOnUpsert()
	err := client.UpsertSecret("ns", "name", map[string]string{"key": "value"})
	if err == nil {
		t.Error("expected error from ThatFailsOnUpsert")
	}
}

func upsertedSecretBytes() map[string]string {
	if capturedUpsertedSecretData == nil {
		panic("UpsertSecret was not called")
	}
	return capturedUpsertedSecretData
}