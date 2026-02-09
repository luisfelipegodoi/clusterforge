package runtime

// New creates a default executable runtime using kubectl + helm binaries.
func New(kubeContext string, timeoutSeconds int) Runtime {
	return &ExecRuntime{
		k: NewKubectlExec(kubeContext, timeoutSeconds),
		h: NewHelmExec(kubeContext, timeoutSeconds),
		l: NewStdLogger(),
	}
}
