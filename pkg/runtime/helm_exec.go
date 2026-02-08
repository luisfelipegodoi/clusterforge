package runtime

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"
)

type HelmExec struct {
	context    string
	timeoutSec int
}

func NewHelmExec(context string, timeoutSeconds int) *HelmExec {
	return &HelmExec{
		context:    context,
		timeoutSec: timeoutSeconds,
	}
}

func (h *HelmExec) UpgradeInstall(release, chart, namespace string, args []string) error {

	baseArgs := []string{
		"--kube-context", h.context,
		"upgrade",
		"--install",
		release,
		chart,
		"-n", namespace,
		"--create-namespace",
		"--wait",
	}

	baseArgs = append(baseArgs, args...)

	cmd := exec.Command("helm", baseArgs...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	done := make(chan error, 1)

	go func() {
		done <- cmd.Run()
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("helm upgrade/install failed: %s", stderr.String())
		}
		return nil

	case <-time.After(time.Duration(h.timeoutSec) * time.Second):
		_ = cmd.Process.Kill()
		return fmt.Errorf("helm timeout installing %s", release)
	}
}
