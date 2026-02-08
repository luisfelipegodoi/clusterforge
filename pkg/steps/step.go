package steps

import (
	"clusterforge/pkg/runtime"
	"context"
	"fmt"
	"time"
)

type Step interface {
	Name() string
	Run(ctx context.Context, rt runtime.Runtime) error
}

type StepFunc struct {
	StepName string
	Fn       func(ctx context.Context, rt runtime.Runtime) error
}

func (s StepFunc) Name() string { return s.StepName }
func (s StepFunc) Run(ctx context.Context, rt runtime.Runtime) error {
	return s.Fn(ctx, rt)
}

type Retry struct {
	Inner       Step
	Attempts    int
	Sleep       time.Duration
	Description string
}

func (r Retry) Name() string {
	if r.Description != "" {
		return r.Description
	}
	return fmt.Sprintf("retry(%s)", r.Inner.Name())
}

func (r Retry) Run(ctx context.Context, rt runtime.Runtime) error {
	attempts := r.Attempts
	if attempts <= 0 {
		attempts = 3
	}
	sleep := r.Sleep
	if sleep <= 0 {
		sleep = 2 * time.Second
	}

	var lastErr error
	for i := 1; i <= attempts; i++ {
		if err := r.Inner.Run(ctx, rt); err != nil {
			lastErr = err
			rt.Logger().Warnf("step failed (attempt %d/%d): %v", i, attempts, err)
			if i < attempts {
				time.Sleep(sleep)
			}
			continue
		}
		return nil
	}
	return lastErr
}
