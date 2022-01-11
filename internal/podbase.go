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
	osType = "linux"
	arch = "amd64"
	// get pod system arch and type
	node, err := configs.RestClient.CoreV1().Nodes().
		Get(context.TODO(), nodeName, metaV1.GetOptions{})
	if err == nil {
		if node.Labels["beta.kubernetes.io/os"] != "" {
			osType = node.Labels["beta.kubernetes.io/os"]
		} else if node.Labels["kubernetes.io/os"] != "" {
			osType = node.Labels["kubernetes.io/os"]
		}
		if node.Labels["beta.kubernetes.io/arch"] != "" {
			arch = node.Labels["beta.kubernetes.io/arch"]
		} else if node.Labels["kubernetes.io/arch"] != "" {
			arch = node.Labels["kubernetes.io/arch"]
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
	osType, arch := p.OsAndArch(pod.Name)
	kfToolsPath := fmt.Sprintf("/kftools_%s_%s", osType, arch)
	if osType == "windows" {
		kfToolsPath = fmt.Sprintf("/kftools_%s_%s.exe", osType, arch)
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
		err := utils.TarKFTools(kfToolsPath, writer)
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
