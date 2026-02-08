package steps

import (
	"context"
	"fmt"

	"clusterforge/pkg/runtime"
)

type WaitCRDEstablished struct {
	CRDName        string
	TimeoutSeconds int
}

func (s WaitCRDEstablished) Name() string {
	return fmt.Sprintf("wait-crd-established(%s)", s.CRDName)
}

func (s WaitCRDEstablished) Run(ctx context.Context, rt runtime.Runtime) error {
	return rt.Kubectl().WaitCRDEstablished(s.CRDName, s.TimeoutSeconds)
}
