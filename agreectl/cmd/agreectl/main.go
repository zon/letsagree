package main

import (
	"agreectl/internal/cluster"
	"agreectl/internal/files"
	"agreectl/internal/opts"
	"agreectl/internal/orchestration"
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
	Context         string `name:"context" default:"microk8s" help:"Kubernetes context"`
	Namespace       string `name:"namespace" default:"letsagree" help:"Kubernetes namespace"`
	DBSecret        string `name:"db-secret" default:"letsagree-app" help:"Secret name for DB credentials"`
	DBHost          string `name:"db-host" help:"Database host (optional, auto-detected if not set)"`
	DBPort          int    `name:"db-port" default:"30432" help:"Database port"`
	RalphNamespace  string `name:"ralph-namespace" default:"ralph-letsagree" help:"Namespace where Ralph is deployed"`
	HPSecret        string `name:"hp-secret" default:"humanity-protocol" help:"Secret name for Humanity Protocol credentials"`
	HPEnvFile       string `name:"hp-env" help:"Path to env file with Humanity Protocol credentials"`
	OIDCIssuer      string `name:"oidc-issuer" default:"https://api.sandbox.humanity.org/v2" help:"OIDC issuer URL"`
	OIDCRedirect    string `name:"oidc-redirect" help:"OIDC redirect URL"`
}

func (c *SetConfig) Run() error {
	return c.RunWith(cluster.NewK8sClient, files.NewConfigWriter())
}

func (c *SetConfig) RunWith(newK8sClient func(string) (cluster.K8sClient, error), cw files.ConfigWriter) error {
	o := opts.Opts{
		Context:        c.Context,
		Namespace:      c.Namespace,
		DBSecret:       c.DBSecret,
		DBHost:         c.DBHost,
		DBPort:         c.DBPort,
		RalphNamespace: c.RalphNamespace,
		HPSecret:       c.HPSecret,
		HPEnvFile:      c.HPEnvFile,
		OIDCIssuer:     c.OIDCIssuer,
		OIDCRedirect:   c.OIDCRedirect,
	}

	k8s, err := newK8sClient(o.Context)
	if err != nil {
		return err
	}

	svc := orchestration.New(k8s, cw)
	return svc.HumanityProtocol(o)
}

func main() {
	ctx := kong.Parse(&cli)
	if err := ctx.Run(); err != nil {
		os.Exit(1)
	}
}