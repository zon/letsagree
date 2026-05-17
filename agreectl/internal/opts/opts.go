package opts

type Opts struct {
	Context   string
	Namespace string
	DBSecret  string
	Host      string
	Port      int
}

var defaultOpts = Opts{
	Context:   "microk8s",
	Namespace: "letsagree",
	DBSecret:  "letsagree-app",
	Port:      30432,
}

func Any() Opts {
	return defaultOpts
}

func WithHost(host string) Opts {
	o := defaultOpts
	o.Host = host
	return o
}

func AnyPort() int {
	return defaultOpts.Port
}

func WithPort(port int) Opts {
	o := defaultOpts
	o.Port = port
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