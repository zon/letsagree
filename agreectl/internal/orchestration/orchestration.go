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
	UpsertSecret(namespace, name string, data map[string]string) error
	NodeIP() (string, error)
}

type ConfigWriter interface {
	WriteJSON(path string, v any) error
	WriteYAML(path string, v any) error
	ParseHPEnv(path string) (files.HPCredentials, error)
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

	if err := o.cluster.UpsertSecret(in.RalphNamespace, in.DBSecret, secret.Data()); err != nil {
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

func (o *Orchestration) HumanityProtocol(in opts.Opts) error {
	var creds files.HPCredentials

	if in.HPEnvFile != "" {
		parsed, err := o.files.ParseHPEnv(in.HPEnvFile)
		if err != nil {
			return err
		}
		if err := o.cluster.UpsertSecret(in.RalphNamespace, in.HPSecret, parsed.ToSecretData()); err != nil {
			return err
		}
		creds = parsed
	} else {
		secret, err := o.cluster.GetSecret(in.RalphNamespace, in.HPSecret)
		if err != nil {
			return err
		}
		creds = files.HPCredentials{
			ClientID:     secret.ClientID(),
			ClientSecret: secret.ClientSecret(),
			PublicKey:    secret.PublicKey(),
		}
	}

	return o.files.WriteYAML(files.HumanityProtocolConfigPath, files.HumanityProtocolConfig{
		ClientID:    creds.ClientID,
		ClientSecret: creds.ClientSecret,
		PublicKey:    creds.PublicKey,
		IssuerURL:    in.OIDCIssuer,
		RedirectURL:  in.OIDCRedirect,
	})
}