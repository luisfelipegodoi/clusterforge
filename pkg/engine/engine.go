package engine

import (
	"context"
	"fmt"

	"github.com/luisfelipegodoi/clusterforge/pkg/runtime"
)

type Engine struct{}

func New() Engine { return Engine{} }

func (e Engine) Run(ctx context.Context, plan Plan, runtimeByCluster map[string]runtime.Runtime) error {

	for _, cluster := range plan.Clusters {
		rt, ok := runtimeByCluster[cluster.ClusterName]
		if !ok {
			return fmt.Errorf("missing runtime for cluster: %s", cluster.ClusterName)
		}

		rt.Logger().Infof("=== cluster: %s ===", cluster.ClusterName)

		for _, step := range cluster.Steps {
			rt.Logger().Infof("-> step: %s", step.Name())
			if err := step.Run(ctx, rt); err != nil {
				rt.Logger().Errorf("step failed: %s: %v", step.Name(), err)
				return fmt.Errorf("cluster=%s step=%s: %w", cluster.ClusterName, step.Name(), err)
			}
			rt.Logger().Infof("<- ok: %s", step.Name())
		}
	}
	return nil
}
