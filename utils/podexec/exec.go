package podexec

import (
	"io"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type PodExec struct {
	K8sClient     *kubernetes.Clientset
	RESTConfig    *rest.Config
	Namespace     string
	PodName       string
	ContainerName string
	Command       []string
	Stdin         io.Reader
	Stdout        io.Writer
	Stderr        io.Writer
	Tty           bool
	NoPreserve    bool
}

func NewPodExec(namespace, podName, containerName string, restConfig *rest.Config, k8sClient *kubernetes.Clientset) *PodExec {
	return &PodExec{
		Namespace:     namespace,
		PodName:       podName,
		ContainerName: containerName,
		RESTConfig:    restConfig,
		K8sClient:     k8sClient,
	}
}

// Exec 在给定容器中执行命令
func (p *PodExec) Exec() error {
	req := p.K8sClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(p.PodName).
		Namespace(p.Namespace).
		SubResource("exec").
		VersionedParams(&coreV1.PodExecOptions{
			Command:   p.Command,
			Container: p.ContainerName,
			Stdin:     p.Stdin != nil,
			Stdout:    p.Stdout != nil,
			Stderr:    p.Stderr != nil,
			TTY:       p.Tty,
		}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(p.RESTConfig, "POST", req.URL())
	if err != nil {
		return err
	}
	var sizeQueue remotecommand.TerminalSizeQueue
	return exec.Stream(remotecommand.StreamOptions{
		Stdin:             p.Stdin,
		Stdout:            p.Stdout,
		Stderr:            p.Stderr,
		Tty:               p.Tty,
		TerminalSizeQueue: sizeQueue,
	})
}
