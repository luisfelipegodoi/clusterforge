package localstack

import "github.com/luisfelipegodoi/clusterforge/pkg/steps"

type Options struct {
	Namespace      string
	ChartPath      string
	WaitTimeoutSec int
}

func Recipe(opts Options) []steps.Step {
	ns := opts.Namespace
	if ns == "" {
		ns = "localstack"
	}
	timeout := opts.WaitTimeoutSec
	if timeout == 0 {
		timeout = 240
	}

	return []steps.Step{
		steps.EnsureNamespace{Namespace: ns},
		steps.KubectlApply{Namespace: ns, Path: opts.ChartPath},

		// espera os principais pods ficarem Ready
		steps.WaitPodsReady{Namespace: ns, Selector: "app.kubernetes.io/name=localstack-server", TimeoutSeconds: timeout},
		steps.WaitPodsReady{Namespace: ns, Selector: "app.kubernetes.io/name=localstack-repo-server", TimeoutSeconds: timeout},
		steps.WaitPodsReady{Namespace: ns, Selector: "app.kubernetes.io/name=localstack-redis", TimeoutSeconds: timeout},
		steps.WaitPodsReady{Namespace: ns, Selector: "app.kubernetes.io/name=localstack-dex-server", TimeoutSeconds: timeout},
	}
}
