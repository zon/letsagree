package opts

type Opts struct {
	Context         string
	Namespace       string
	DBSecret        string
	DBHost          string
	DBPort          int
	RalphNamespace  string
}

var defaultOpts = Opts{
	Context:        "microk8s",
	Namespace:      "letsagree",
	DBSecret:       "letsagree-app",
	DBPort:         30432,
	RalphNamespace: "ralph-letsagree",
}

func Any() Opts {
	return defaultOpts
}

func WithDBHost(host string) Opts {
	o := defaultOpts
	o.DBHost = host
	return o
}

func AnyDBPort() int {
	return defaultOpts.DBPort
}

func WithDBPort(port int) Opts {
	o := defaultOpts
	o.DBPort = port
	return o
}

func WithRalphNamespace(ns string) Opts {
	o := defaultOpts
	o.RalphNamespace = ns
	return o
}

func WithContext(context string) Opts {
	o := defaultOpts
	o.Context = context
	return o
}

func WithNamespace(namespace string) Opts {
	o := defaultOpts
	o.Namespace = namespace
	return o
}

func WithDBSecret(secret string) Opts {
	o := defaultOpts
	o.DBSecret = secret
	return o
}