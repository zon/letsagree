package cluster

import "testing"

type GetSecretCall struct {
	Namespace string
	Secret    string
}

type StubK8sClient struct {
	Secret    *Secret
	RetNodeIP string
	SecErr    error
	UpsertErr error
	Calls     []GetSecretCall
}

func (s *StubK8sClient) GetSecret(namespace, name string) (*Secret, error) {
	s.Calls = append(s.Calls, GetSecretCall{Namespace: namespace, Secret: name})
	return s.Secret, s.SecErr
}

func (s *StubK8sClient) UpsertSecret(namespace, name string, data map[string]string) error {
	return s.UpsertErr
}

func (s *StubK8sClient) NodeIP() (string, error) {
	return s.RetNodeIP, nil
}

type stubK8sClient struct {
	secret    *Secret
	nodeIP    string
	nodeErr   error
	secErr    error
	upsertErr error
}

func (s *stubK8sClient) GetSecret(namespace, name string) (*Secret, error) {
	return s.secret, s.secErr
}

func (s *stubK8sClient) UpsertSecret(namespace, name string, data map[string]string) error {
	capturedUpsertedSecretData = data
	return s.upsertErr
}

func (s *stubK8sClient) NodeIP() (string, error) {
	return s.nodeIP, s.nodeErr
}

func SecretFromStringData(data map[string]string) *Secret {
	bytes := make(map[string][]byte, len(data))
	for k, v := range data {
		bytes[k] = []byte(v)
	}
	return &Secret{Bytes: bytes}
}

func AnySecret() *Secret {
	return &Secret{Bytes: map[string][]byte{
		"host":     []byte("letsagree-rw"),
		"port":     []byte("5432"),
		"user":     []byte("app"),
		"password": []byte("secret"),
		"dbname":   []byte("app"),
	}}
}

func WithSecret(secret *Secret) K8sClient {
	return &stubK8sClient{secret: secret}
}

func AnyNodeIP() string {
	return "192.168.1.10"
}

func WithNodeIP(ip string) K8sClient {
	return &stubK8sClient{nodeIP: ip, secret: AnySecret()}
}

func ThatFailsOnNodeIP() K8sClient {
	return &stubK8sClient{nodeErr: assertNeverNodeIP{}, secret: AnySecret()}
}

func ThatFailsOnUpsert() K8sClient {
	return &stubK8sClient{upsertErr: assertNeverUpsert{}, secret: AnySecret()}
}

func ThatFailsOnGetSecret() K8sClient {
	return &stubK8sClient{secErr: assertNeverGetSecret{}}
}

type assertNeverNodeIP struct{}

func (assertNeverNodeIP) Error() string {
	return "NodeIP should not be called"
}

type assertNeverUpsert struct{}

func (assertNeverUpsert) Error() string {
	return "UpsertSecret should not be called"
}

type assertNeverGetSecret struct{}

func (assertNeverGetSecret) Error() string {
	return "GetSecret should not be called"
}

var capturedCall struct {
	namespace string
	secret    string
}

var capturedUpsertedSecretData map[string]string

type capturingK8sClient struct {
	secret *Secret
	nodeIP string
}

func (c *capturingK8sClient) GetSecret(namespace, name string) (*Secret, error) {
	capturedCall.namespace = namespace
	capturedCall.secret = name
	return c.secret, nil
}

func (c *capturingK8sClient) UpsertSecret(namespace, name string, data map[string]string) error {
	capturedUpsertedSecretData = data
	return nil
}

func (c *capturingK8sClient) NodeIP() (string, error) {
	return c.nodeIP, nil
}

func WithSecretAndCapturing(secret *Secret) K8sClient {
	return &capturingK8sClient{secret: secret, nodeIP: "10.0.0.1"}
}

func CapturedGetSecret() (namespace, name string) {
	return capturedCall.namespace, capturedCall.secret
}

func UpsertedSecretData(t testing.TB) map[string]string {
	if capturedUpsertedSecretData == nil {
		t.Fatal("UpsertSecret was not called")
	}
	return capturedUpsertedSecretData
}