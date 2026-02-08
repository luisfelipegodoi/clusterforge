package steps

import (
	"context"
	"fmt"

	"clusterforge/pkg/runtime"
)

type HelmUpgradeInstall struct {
	Release   string
	Chart     string
	Namespace string
	Args      []string
}

func (s HelmUpgradeInstall) Name() string {
	return fmt.Sprintf("helm-upgrade-install(%s:%s)", s.Namespace, s.Release)
}

func (s HelmUpgradeInstall) Run(ctx context.Context, rt runtime.Runtime) error {
	return rt.Helm().UpgradeInstall(s.Release, s.Chart, s.Namespace, s.Args)
}
