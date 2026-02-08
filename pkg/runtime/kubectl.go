package runtime

import (
	"fmt"
	"strings"
)

type KubectlCLI struct {
	Context string
	Exec    Exec
}

func (k KubectlCLI) baseArgs() []string {
	if k.Context == "" {
		return []string{}
	}
	return []string{"--context", k.Context}
}

func (k KubectlCLI) Apply(namespace, path string) error {
	args := append(k.baseArgs(), "apply")
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	args = append(args, "-f", path)
	_, err := k.Exec.Run("kubectl", args, "")
	return err
}

func (k KubectlCLI) ApplyStdin(namespace string, yaml string) error {
	args := append(k.baseArgs(), "apply")
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	args = append(args, "-f", "-")
	_, err := k.Exec.Run("kubectl", args, yaml)
	return err
}

func (k KubectlCLI) WaitPodsReady(namespace, selector string, timeoutSeconds int) error {
	if timeoutSeconds <= 0 {
		timeoutSeconds = 180
	}
	args := append(k.baseArgs(), "wait")
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	args = append(args,
		"--for=condition=Ready",
		"pod",
		"-l", selector,
		fmt.Sprintf("--timeout=%ds", timeoutSeconds),
	)
	_, err := k.Exec.Run("kubectl", args, "")
	return err
}

func (k KubectlCLI) WaitDeploymentReady(namespace, name string, timeoutSeconds int) error {
	if timeoutSeconds <= 0 {
		timeoutSeconds = 240
	}
	args := append(k.baseArgs(), "rollout", "status", "deployment/"+name)
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	args = append(args, fmt.Sprintf("--timeout=%ds", timeoutSeconds))
	_, err := k.Exec.Run("kubectl", args, "")
	return err
}

func (k KubectlCLI) WaitCRDEstablished(crdName string, timeoutSeconds int) error {
	if timeoutSeconds <= 0 {
		timeoutSeconds = 240
	}
	// Usa kubectl wait no CRD:
	// kubectl wait --for=condition=Established crd/<name> --timeout=...
	args := append(k.baseArgs(),
		"wait",
		"--for=condition=Established",
		"crd/"+crdName,
		fmt.Sprintf("--timeout=%ds", timeoutSeconds),
	)
	_, err := k.Exec.Run("kubectl", args, "")
	return err
}

func (k KubectlCLI) GetJSONPath(namespace string, args []string, jsonpath string) (string, error) {
	// Ex: args = []{"get","pods","-l","app=xyz"}
	cmdArgs := append(k.baseArgs(), args...)
	if namespace != "" {
		// tenta inserir -n <ns> após "get"
		// se o caller já passou -n, não repete
		if !contains(cmdArgs, "-n") && !contains(cmdArgs, "--namespace") {
			// heurística simples
			cmdArgs = append(cmdArgs, "-n", namespace)
		}
	}
	cmdArgs = append(cmdArgs, "-o", "jsonpath="+jsonpath)

	out, err := k.Exec.Run("kubectl", cmdArgs, "")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func contains(xs []string, v string) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}
