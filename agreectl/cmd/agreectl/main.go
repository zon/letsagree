package main

import (
	"agreectl/internal/cluster"
	"agreectl/internal/files"
	"agreectl/internal/opts"
	"agreectl/internal/orchestration"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

var cli struct {
	Set SetCmd `cmd:"" help:"Write backend configuration from Kubernetes."`
}

type SetCmd struct {
	Config SetConfig `cmd:"" help:"Write backend configuration from Kubernetes."`
}

type SetConfig struct {
	Context        string `name:"context" default:"microk8s" help:"Kubernetes context"`
	Namespace      string `name:"namespace" default:"letsagree" help:"Kubernetes namespace"`
	DBSecret       string `name:"db-secret" default:"letsagree-app" help:"Secret name for DB credentials"`
	DBPort         int    `name:"db-port" default:"30432" help:"NodePort for local postgres access"`
	PostgresSecret string `name:"postgres-secret" default:"postgres" help:"Secret name for postgres config in Ralph's namespace"`
	HPSecret       string `name:"hp-secret" default:"humanity-protocol" help:"Secret name for Humanity Protocol config in Ralph's namespace"`
	RalphNamespace string `name:"ralph-namespace" default:"ralph-letsagree" help:"Namespace where Ralph is deployed"`
	HPEnvFile      string `name:"hp-env" help:"Path to env file with Humanity Protocol credentials"`
	OIDCIssuer     string `name:"oidc-issuer" default:"https://api.sandbox.humanity.org/v2" help:"OIDC issuer URL"`
	OIDCRedirect   string `name:"oidc-redirect" help:"OIDC redirect URL"`
}

func (c *SetConfig) Run() error {
	return c.RunWith(cluster.NewK8sClient, files.NewConfigWriter())
}

func (c *SetConfig) RunWith(newK8sClient func(string) (cluster.K8sClient, error), cw files.ConfigWriter) error {
	o := opts.Opts{
		Context:        c.Context,
		Namespace:      c.Namespace,
		DBSecret:       c.DBSecret,
		DBPort:         c.DBPort,
		PostgresSecret: c.PostgresSecret,
		HPSecret:       c.HPSecret,
		RalphNamespace: c.RalphNamespace,
		HPEnvFile:      c.HPEnvFile,
		OIDCIssuer:     c.OIDCIssuer,
		OIDCRedirect:   c.OIDCRedirect,
	}

	k8s, err := newK8sClient(o.Context)
	if err != nil {
		return err
	}

	svc := orchestration.New(k8s, cw)
	if err := svc.Postgres(o); err != nil {
		return err
	}
	return svc.HumanityProtocol(o)
}

func main() {
	ctx := kong.Parse(&cli)
	if err := ctx.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}