package engine

import "github.com/luisfelipegodoi/clusterforge/pkg/steps"

type ClusterPlan struct {
	ClusterName string
	Steps       []steps.Step
}

type Plan struct {
	Clusters []ClusterPlan
}

func NewPlan() Plan { return Plan{} }

func (p *Plan) AddCluster(name string) *ClusterPlan {
	p.Clusters = append(p.Clusters, ClusterPlan{ClusterName: name})
	return &p.Clusters[len(p.Clusters)-1]
}

func (c *ClusterPlan) Add(step steps.Step) *ClusterPlan {
	c.Steps = append(c.Steps, step)
	return c
}
