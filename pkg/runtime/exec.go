package runtime

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type Exec struct {
	Timeout time.Duration
	Logger  Logger
}

func (e Exec) Run(name string, args []string, stdin string) (string, error) {
	timeout := e.Timeout
	if timeout <= 0 {
		timeout = 2 * time.Minute
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}

	if e.Logger != nil {
		e.Logger.Infof("exec: %s %s", name, strings.Join(args, " "))
	}

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return out.String(), fmt.Errorf("command timeout: %s %s", name, strings.Join(args, " "))
	}
	if err != nil {
		return out.String(), fmt.Errorf("command error: %v; stderr=%s", err, stderr.String())
	}
	return out.String(), nil
}
