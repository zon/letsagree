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

func (o *Orchestration) Postgres(in opts.Opts) error {
	secret, err := o.cluster.GetSecret(in.Namespace, in.DBSecret)
	if err != nil {
		return err
	}

	host := in.DBHost
	if host == "" {
		host, err = o.cluster.NodeIP()
		if err != nil {
			return err
		}
	}

	return o.files.WriteJSON(files.PostgresConfigPath, files.PostgresConfig{
		Host:     host,
		Port:     in.DBPort,
		User:     secret.User(),
		Password: secret.Password(),
		DBName:   secret.DBName(),
	})
}