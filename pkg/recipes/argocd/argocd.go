package argocd

import "github.com/luisfelipegodoi/clusterforge/pkg/steps"

type Options struct {
	Namespace      string
	ManifestPath   string
	WaitTimeoutSec int
}

func Recipe(opts Options) []steps.Step {
	ns := opts.Namespace
	if ns == "" {
		ns = "argocd"
	}
	timeout := opts.WaitTimeoutSec
	if timeout == 0 {
		timeout = 240
	}

	return []steps.Step{
		steps.EnsureNamespace{Namespace: ns},
		steps.KubectlApply{Namespace: ns, Path: opts.ManifestPath},

		// espera os principais pods ficarem Ready
		steps.WaitPodsReady{Namespace: ns, Selector: "app.kubernetes.io/name=argocd-server", TimeoutSeconds: timeout},
		steps.WaitPodsReady{Namespace: ns, Selector: "app.kubernetes.io/name=argocd-repo-server", TimeoutSeconds: timeout},
		steps.WaitPodsReady{Namespace: ns, Selector: "app.kubernetes.io/name=argocd-redis", TimeoutSeconds: timeout},
		steps.WaitPodsReady{Namespace: ns, Selector: "app.kubernetes.io/name=argocd-dex-server", TimeoutSeconds: timeout},
	}
}
