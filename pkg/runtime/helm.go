package runtime

type HelmCLI struct {
	KubeContext string
	Exec        Exec
}

func (h HelmCLI) UpgradeInstall(release, chart, namespace string, args []string) error {
	cmdArgs := []string{}
	if h.KubeContext != "" {
		cmdArgs = append(cmdArgs, "--kube-context", h.KubeContext)
	}
	cmdArgs = append(cmdArgs,
		"upgrade", "--install",
		release, chart,
	)
	if namespace != "" {
		cmdArgs = append(cmdArgs, "-n", namespace, "--create-namespace")
	}
	cmdArgs = append(cmdArgs, args...)

	_, err := h.Exec.Run("helm", cmdArgs, "")
	return err
}
