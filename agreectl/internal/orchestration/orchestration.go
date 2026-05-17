package orchestration

import (
	"agreectl/internal/cluster"
	"agreectl/internal/files"
	"agreectl/internal/opts"
)

type Orchestration struct {
	cluster K8sClient
	files   ConfigWriter
}

type K8sClient interface {
	GetSecret(namespace, name string) (*cluster.Secret, error)
	NodeIP() (string, error)
}

type ConfigWriter interface {
	WriteJSON(path string, v any) error
}

func New(cluster K8sClient, files ConfigWriter) *Orchestration {
	return &Orchestration{
		cluster: cluster,
		files:   files,
	}
}

func WithMocks(overrides ...any) *Orchestration {
	defaultK8s := cluster.WithSecret(cluster.AnySecret())
	defaultCW := &files.CapturingConfigWriter{}
	var k8s K8sClient = defaultK8s
	var cw ConfigWriter = defaultCW

	for _, o := range overrides {
		switch v := o.(type) {
		case K8sClient:
			k8s = v
		case ConfigWriter:
			cw = v
		}
	}

	return New(k8s, cw)
}

type StubK8sClient struct {
	Secret     *cluster.Secret
	RetNodeIP  string
	NodeErr    error
	Calls      struct {
		Namespace string
		Secret    string
	}
}

func (s *StubK8sClient) GetSecret(namespace, name string) (*cluster.Secret, error) {
	s.Calls.Namespace = namespace
	s.Calls.Secret = name
	return s.Secret, s.NodeErr
}

func (s *StubK8sClient) NodeIP() (string, error) {
	return s.RetNodeIP, s.NodeErr
}

func (o *Orchestration) Postgres(in opts.Opts) error {
	secret, err := o.cluster.GetSecret(in.Namespace, in.DBSecret)
	if err != nil {
		return err
	}

	host := in.Host
	if host == "" {
		host, err = o.cluster.NodeIP()
		if err != nil {
			return err
		}
	}

	return o.files.WriteJSON(files.PostgresConfigPath, files.PostgresConfig{
		Host:     host,
		Port:     in.Port,
		User:     secret.User(),
		Password: secret.Password(),
		DBName:   secret.DBName(),
	})
}