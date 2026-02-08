package steps

import (
	"context"
	"fmt"

	"github.com/luisfelipegodoi/clusterforge/pkg/runtime"
)

type WaitDeploymentReady struct {
	Namespace      string
	DeployName     string
	TimeoutSeconds int
}

func (s WaitDeploymentReady) Name() string {
	return fmt.Sprintf("wait-deploy-ready(ns=%s,name=%s)", s.Namespace, s.Name)
}

func (s WaitDeploymentReady) Run(ctx context.Context, rt runtime.Runtime) error {
	return rt.Kubectl().WaitDeploymentReady(s.Namespace, s.DeployName, s.TimeoutSeconds)
}
