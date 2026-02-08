package runtime

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"
)

type KubectlExec struct {
	context    string
	timeoutSec int
}

func NewKubectlExec(context string, timeoutSeconds int) *KubectlExec {
	return &KubectlExec{
		context:    context,
		timeoutSec: timeoutSeconds,
	}
}

func (k *KubectlExec) run(args ...string) ([]byte, error) {
	ctxArgs := []string{"--context", k.context}
	ctxArgs = append(ctxArgs, args...)

	cmd := exec.Command("kubectl", ctxArgs...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	done := make(chan error, 1)

	go func() {
		done <- cmd.Run()
	}()

	select {
	case err := <-done:
		if err != nil {
			return nil, fmt.Errorf("kubectl %v failed: %s", args, stderr.String())
		}
		return stdout.Bytes(), nil

	case <-time.After(time.Duration(k.timeoutSec) * time.Second):
		_ = cmd.Process.Kill()
		return nil, fmt.Errorf("kubectl timeout: %v", args)
	}
}
