package internal

import (
	"context"

	"github.com/xmapst/kubefilebrowser/configs"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PvcBase struct {
	Namespace string
	PvcName   string
}

func (p *PvcBase) PvcInfo() (*coreV1.PersistentVolumeClaim, error) {
	return configs.RestClient.CoreV1().PersistentVolumeClaims(p.Namespace).
		Get(context.TODO(), p.PvcName, metaV1.GetOptions{})
}
