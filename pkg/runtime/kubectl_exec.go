package runtime

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
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

func (k *KubectlExec) run(args []string) ([]byte, error) {
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

func (k *KubectlExec) Apply(namespace, path string) error {
	args := []string{"--context", k.context}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	args = append(args, "apply", "-f", path)
	_, err := k.run(args)
	return err
}

func (k *KubectlExec) ApplyStdin(namespace string, yaml string) error {
	args := []string{"--context", k.context}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	args = append(args, "apply", "-f", "-")
	_, err := k.run(args)
	return err
}

func (k *KubectlExec) GetJSONPath(namespace string, args []string, jsonpath string) (string, error) {
	base := []string{"--context", k.context}
	if namespace != "" {
		base = append(base, "-n", namespace)
	}
	// args aqui é algo tipo: []string{"get", "providers"} etc.
	base = append(base, args...)
	base = append(base, "-o", "jsonpath="+jsonpath)

	out, err := k.run(base)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func (k *KubectlExec) WaitDeploymentReady(namespace, name string, timeoutSeconds int) error {
	args := []string{"--context", k.context}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	args = append(args, "rollout", "status", "deployment/"+name, fmt.Sprintf("--timeout=%ds", timeoutSeconds))
	_, err := k.run(args)
	return err
}

func (k *KubectlExec) WaitCRDEstablished(crdName string, timeoutSeconds int) error {
	// espera o CRD aparecer no API server
	args := []string{"--context", k.context, "wait", "--for=condition=Established", "crd/" + crdName, fmt.Sprintf("--timeout=%ds", timeoutSeconds)}
	_, err := k.run(args)
	return err
}

func (k *KubectlExec) WaitPodsReady(namespace, selector string, timeoutSeconds int) error {
	args := []string{"--context", k.context}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	args = append(args, "wait", "--for=condition=Ready", "pod", "-l", selector, fmt.Sprintf("--timeout=%ds", timeoutSeconds))
	_, err := k.run(args)
	return err
}

// func (k *KubectlExec) run(args []string, stdin string) (string, error) {
// 	to := time.Duration(k.timeoutSeconds) * time.Second
// 	if to <= 0 {
// 		to = 60 * time.Second
// 	}

// 	cmd := exec.Command("kubectl", args...)

// 	var stdout, stderr bytes.Buffer
// 	cmd.Stdout = &stdout
// 	cmd.Stderr = &stderr
// 	if stdin != "" {
// 		cmd.Stdin = strings.NewReader(stdin)
// 	}

// 	// timeout simples (sem goroutine): use `exec.CommandContext` se quiser hard-timeout por ctx
// 	if err := cmd.Run(); err != nil {
// 		return "", fmt.Errorf("kubectl %v failed: %w\nstderr: %s", args, err, strings.TrimSpace(stderr.String()))
// 	}

// 	return stdout.String(), nil
// }
