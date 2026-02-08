package runtime

type Runtime interface {
	Kubectl() Kubectl
	Helm() Helm
	Logger() Logger
}

type Logger interface {
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
}

type Kubectl interface {
	Apply(namespace, path string) error
	ApplyStdin(namespace string, yaml string) error
	WaitPodsReady(namespace, selector string, timeoutSeconds int) error
	WaitDeploymentReady(namespace, name string, timeoutSeconds int) error
	WaitCRDEstablished(crdName string, timeoutSeconds int) error
	GetJSONPath(namespace string, args []string, jsonpath string) (string, error)
}

type Helm interface {
	UpgradeInstall(release, chart, namespace string, args []string) error
}
