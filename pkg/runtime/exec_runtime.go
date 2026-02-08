package runtime

type ExecRuntime struct {
	kubeContext string
	timeoutSec  int

	k Kubectl
	h Helm
	l Logger
}

func NewExecRuntime(kubeContext string, timeoutSeconds int) *ExecRuntime {
	r := &ExecRuntime{
		kubeContext: kubeContext,
		timeoutSec:  timeoutSeconds,
	}

	r.k = KubectlExec(kubeContext, timeoutSeconds)
	r.h = NewHelmExec(kubeContext, timeoutSeconds)
	r.l = NewStdLogger()

	return r
}

func (r *ExecRuntime) Kubectl() Kubectl {
	return r.k
}

func (r *ExecRuntime) Helm() Helm {
	return r.h
}

func (r *ExecRuntime) Logger() Logger {
	return r.l
}
