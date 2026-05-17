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