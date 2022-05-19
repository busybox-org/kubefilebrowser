package kubeapiproxy

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xmapst/kubefilebrowser/configs"
	"github.com/xmapst/kubefilebrowser/handlers"
	"github.com/xmapst/kubefilebrowser/internal"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PvcsQuery struct {
	Namespace     string   `json:"namespace" form:"namespace"`
	Pvcs          []string `json:"pvcs" form:"pvcs"`
	FieldSelector string   `json:"field_selector" form:"field_selector"`
	LabelSelector string   `json:"label_selector" form:"label_selector"`
	Limit         int64    `json:"limit" form:"limit"`
	Continue      string   `json:"continue" form:"continue"`
}

// Pvcs
// @Summary Pvcs
// @description 命名空间下Pvc资源列表
// @Tags Kubernetes Api Proxy
// @Param namespace query PvcsQuery false "namespace" default(test)
// @Param pvcs query PvcsQuery false "pvcs"
// @Param field_selector query PvcsQuery false "field_selector"
// @Param label_selector query PvcsQuery false "label_selector"
// @Success 200 {object} handlers.JSONResult
// @Failure 500 {object} handlers.JSONResult
// @Router /api/kubeapiproxy/pvcs [get]
func Pvcs(c *gin.Context) {
	render := handlers.Gin{Context: c}
	var q PvcsQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		logrus.Error(err)
		render.SetError(handlers.CodeErrParam, err)
		return
	}
	listOptions := metaV1.ListOptions{
		LabelSelector: q.LabelSelector,
		FieldSelector: q.FieldSelector,
		Limit:         q.Limit,
		Continue:      q.Continue,
	}

	// 所有命名空间
	if q.Namespace == "" {
		res, err := configs.RestClient.CoreV1().
			PersistentVolumeClaims(metaV1.NamespaceAll).
			List(context.TODO(), listOptions)
		if err != nil {
			logrus.Error(err)
			render.SetError(handlers.CodeErrApp, err)
			return
		}
		var wg = new(sync.WaitGroup)
		for k, v := range res.Items {
			wg.Add(1)
			go func(k int, v coreV1.PersistentVolumeClaim) {
				defer wg.Done()
				if res.Items[k].ObjectMeta.Annotations == nil {
					res.Items[k].ObjectMeta.Annotations = make(map[string]string)
				}
			}(k, v)
		}
		wg.Wait()
		render.SetJson(res)
		return
	}
	// 指定pvcs
	if len(q.Pvcs) != 0 {
		var res []*coreV1.PersistentVolumeClaim
		var lock = new(sync.Mutex)
		var wg = new(sync.WaitGroup)
		for _, pvcName := range q.Pvcs {
			wg.Add(1)
			go func(pvcName string) {
				defer wg.Done()
				pb := internal.PvcBase{
					Namespace: q.Namespace,
					PvcName:   pvcName,
				}
				pvc, err := pb.PvcInfo()
				if err != nil {
					logrus.Error(err)
					return
				}
				lock.Lock()
				defer lock.Unlock()
				res = append(res, pvc)
			}(pvcName)
		}
		wg.Wait()
		render.SetJson(res)
	}
	// 指定命名空间
	res, err := configs.RestClient.CoreV1().
		PersistentVolumeClaims(q.Namespace).
		List(context.TODO(), listOptions)
	if err != nil {
		logrus.Error(err)
		render.SetError(handlers.CodeErrApp, err)
		return
	}
	var wg = new(sync.WaitGroup)
	for k, v := range res.Items {
		wg.Add(1)
		go func(k int, v coreV1.PersistentVolumeClaim) {
			defer wg.Done()
			if res.Items[k].ObjectMeta.Annotations == nil {
				res.Items[k].ObjectMeta.Annotations = make(map[string]string)
			}
		}(k, v)
	}
	wg.Wait()
	render.SetJson(res)
}
