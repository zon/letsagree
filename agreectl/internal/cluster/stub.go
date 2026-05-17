package cluster

type stubK8sClient struct {
	secret  *Secret
	nodeIP  string
	nodeErr error
	secErr  error
}

func (s *stubK8sClient) GetSecret(namespace, name string) (*Secret, error) {
	return s.secret, s.secErr
}

func (s *stubK8sClient) NodeIP() (string, error) {
	return s.nodeIP, s.nodeErr
}

func AnySecret() *Secret {
	return &Secret{Data: map[string][]byte{
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
	return &stubK8sClient{nodeIP: ip, secret: &Secret{Data: map[string][]byte{
		"user":     []byte("app"),
		"password": []byte("secret"),
		"dbname":   []byte("app"),
	}}}
}

func ThatFailsOnNodeIP() K8sClient {
	return &stubK8sClient{nodeErr: assertNeverError{}, secret: &Secret{Data: map[string][]byte{
		"user":     []byte("app"),
		"password": []byte("secret"),
		"dbname":   []byte("app"),
	}}}
}

type assertNeverError struct{}

func (assertNeverError) Error() string {
	return "NodeIP should not be called"
}

var capturedCall struct {
	namespace string
	secret    string
}

type capturingK8sClient struct {
	secret  *Secret
	nodeIP  string
	nodeErr error
}

func (c *capturingK8sClient) GetSecret(namespace, name string) (*Secret, error) {
	capturedCall.namespace = namespace
	capturedCall.secret = name
	return c.secret, nil
}

func (c *capturingK8sClient) NodeIP() (string, error) {
	return c.nodeIP, c.nodeErr
}

func WithSecretAndCapturing(secret *Secret) K8sClient {
	return &capturingK8sClient{
		secret: secret,
		nodeIP: "10.0.0.1",
	}
}

func CapturedGetSecret() (namespace, name string) {
	return capturedCall.namespace, capturedCall.secret
}