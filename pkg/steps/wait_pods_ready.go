package steps

import (
	"context"
	"fmt"

	"clusterforge/pkg/runtime"
)

type WaitPodsReady struct {
	Namespace      string
	Selector       string
	TimeoutSeconds int
}

func (s WaitPodsReady) Name() string {
	return fmt.Sprintf("wait-pods-ready(ns=%s,sel=%s)", s.Namespace, s.Selector)
}

func (s WaitPodsReady) Run(ctx context.Context, rt runtime.Runtime) error {
	return rt.Kubectl().WaitPodsReady(s.Namespace, s.Selector, s.TimeoutSeconds)
}
