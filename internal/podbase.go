package internal

import (
	"bytes"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xmapst/kubefilebrowser/configs"
	"github.com/xmapst/kubefilebrowser/utils"
	"github.com/xmapst/kubefilebrowser/utils/podexec"
	"io"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodBase struct {
	Namespace string
	PodName   string
	Container string
}

func (p *PodBase) NewPodExec() *podexec.PodExec {
	return podexec.NewPodExec(p.Namespace, p.PodName, p.Container, configs.KuBeResConf, configs.RestClient)
}

func (p *PodBase) PodInfo() (*coreV1.Pod, error) {
	pod, err := configs.RestClient.CoreV1().Pods(p.Namespace).
		Get(context.TODO(), p.PodName, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if pod.Status.Phase == coreV1.PodSucceeded || pod.Status.Phase == coreV1.PodFailed {
		return nil, fmt.Errorf("cannot exec into a container in a completed pod; current phase is %s", pod.Status.Phase)
	}
	return pod, nil
}

func (p *PodBase) OsAndArch(nodeName string) (osType string, arch string) {
	// get pod system arch and type
	node, err := configs.RestClient.CoreV1().Nodes().
		Get(context.TODO(), nodeName, metaV1.GetOptions{})
	if err == nil {
		var ok bool
		osType, ok = node.Labels["beta.kubernetes.io/os"]
		if !ok {
			osType, ok = node.Labels["kubernetes.io/os"]
			if !ok {
				osType = "linux"
			}
		}
		arch, ok = node.Labels["beta.kubernetes.io/arch"]
		if !ok {
			arch, ok = node.Labels["kubernetes.io/arch"]
			if !ok {
				arch = "amd64"
			}
		}
	}
	return
}

func (p *PodBase) Exec(stdin io.Reader, command ...string) ([]byte, error) {
	var stdout, stderr bytes.Buffer
	exec := p.NewPodExec()
	exec.Command = command
	exec.Tty = false
	exec.Stdin = stdin
	exec.Stdout = &stdout
	exec.Stderr = &stderr
	err := exec.Exec()
	if err != nil {
		if len(stderr.Bytes()) != 0 {
			return nil, fmt.Errorf("%s %s", err.Error(), stderr.String())
		}
		return nil, err
	}
	if len(stderr.Bytes()) != 0 {
		return nil, fmt.Errorf(stderr.String())
	}
	return stdout.Bytes(), nil
}

func (p *PodBase) InstallKFTools() error {
	pod, err := p.PodInfo()
	if err != nil {
		return err
	}
	osType, arch := p.OsAndArch(pod.Spec.NodeName)
	kfToolsPath := fmt.Sprintf("/kftools_%s_%s", osType, arch)
	if osType == "windows" {
		kfToolsPath+=".exe"
	}

	reader, writer := io.Pipe()
	exec := p.NewPodExec()
	exec.Stdin = reader

	go func() {
		defer func() {
			{
				_ = writer.Close()
			}
		}()
		err = utils.TarKFTools(kfToolsPath, writer)
		if err != nil {
			logrus.Error(err)
		}
	}()
	err = exec.ToPodContainer("/")
	if err != nil {
		return err
	}
	if osType != "windows" {
		chmodCmd := []string{"chmod", "+x", "/kftools"}
		exec.Command = chmodCmd
		var stderr bytes.Buffer
		exec.Stderr = &stderr
		err = exec.Exec()
		if err != nil {
			return fmt.Errorf(err.Error(), stderr)
		}
	}
	return nil
}
