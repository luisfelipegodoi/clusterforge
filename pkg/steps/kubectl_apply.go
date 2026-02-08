package steps

import (
	"context"
	"fmt"

	"github.com/luisfelipegodoi/clusterforge/pkg/runtime"
)

type KubectlApply struct {
	Namespace string
	Path      string
}

func (s KubectlApply) Name() string {
	return fmt.Sprintf("kubectl-apply(ns=%s,path=%s)", s.Namespace, s.Path)
}

func (s KubectlApply) Run(ctx context.Context, rt runtime.Runtime) error {
	return rt.Kubectl().Apply(s.Namespace, s.Path)
}
