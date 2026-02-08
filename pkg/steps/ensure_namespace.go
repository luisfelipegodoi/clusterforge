package steps

import (
	"context"
	"fmt"

	"clusterforge/pkg/runtime"
)

type EnsureNamespace struct {
	Namespace string
}

func (s EnsureNamespace) Name() string { return fmt.Sprintf("ensure-namespace(%s)", s.Namespace) }

func (s EnsureNamespace) Run(ctx context.Context, rt runtime.Runtime) error {
	if s.Namespace == "" {
		return nil
	}
	yaml := fmt.Sprintf(`apiVersion: v1
kind: Namespace
metadata:
  name: %s
`, s.Namespace)

	return rt.Kubectl().ApplyStdin("", yaml)
}
