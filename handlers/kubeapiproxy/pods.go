package kubeapiproxy

import (
    "context"
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
    "github.com/xmapst/kubefilebrowser/configs"
    "github.com/xmapst/kubefilebrowser/handlers"
    "github.com/xmapst/kubefilebrowser/internal"
    coreV1 "k8s.io/api/core/v1"
    metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "sync"
)

type PodsQuery struct {
	Namespace     string   `json:"namespace" form:"namespace"`
	Pods          []string `json:"pods" form:"pods"`
	FieldSelector string   `json:"field_selector" form:"field_selector"`
	LabelSelector string   `json:"label_selector" form:"label_selector"`
	Limit         int64    `json:"limit" form:"limit"`
	Continue      string   `json:"continue" form:"continue"`
}

// Pods
// @Summary Pods
// @description 命名空间下Pod资源列表
// @Tags Kubernetes Api Proxy
// @Param namespace query PodsQuery false "namespace" default(test)
// @Param pods query PodsQuery false "pods"
// @Param field_selector query PodsQuery false "field_selector"
// @Param label_selector query PodsQuery false "label_selector"
// @Success 200 {object} handlers.JSONResult
// @Failure 500 {object} handlers.JSONResult
// @Router /api/kubeapiproxy/pods [get]
func Pods(c *gin.Context) {
	render := handlers.Gin{Context: c}
	var q PodsQuery
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
			Pods(metaV1.NamespaceAll).
			List(context.TODO(), listOptions)
		if err != nil {
			logrus.Error(err)
			render.SetError(handlers.CodeErrApp, err)
			return
		}
		var wg = new(sync.WaitGroup)
		for k, v := range res.Items {
			wg.Add(1)
			go func(k int, v coreV1.Pod) {
				defer wg.Done()
				if res.Items[k].ObjectMeta.Annotations == nil {
					res.Items[k].ObjectMeta.Annotations = make(map[string]string)
				}
				pb := internal.PodBase{
					Namespace: q.Namespace,
					PodName:   v.Name,
				}
				res.Items[k].ObjectMeta.Annotations["os"], res.Items[k].ObjectMeta.Annotations["arch"] =
					pb.OsAndArch(v.Spec.NodeName)
			}(k, v)
		}
		wg.Wait()
		render.SetJson(res)
		return
	}
	// 指定pods
	if len(q.Pods) != 0 {
		var res []*coreV1.Pod
		var lock = new(sync.Mutex)
		var wg = new(sync.WaitGroup)
		for _, podName := range q.Pods {
			wg.Add(1)
			go func(podName string) {
				defer wg.Done()
				pb := internal.PodBase{
					Namespace: q.Namespace,
					PodName:   podName,
				}
				pod, err := pb.PodInfo()
				if err != nil {
					logrus.Error(err)
					return
				}
				pod.ObjectMeta.Annotations["os"], pod.ObjectMeta.Annotations["arch"] = pb.OsAndArch(pod.Name)
				lock.Lock()
				defer lock.Unlock()
				res = append(res, pod)
			}(podName)
		}
		wg.Wait()
		render.SetJson(res)
	}
	// 指定命名空间
	res, err := configs.RestClient.CoreV1().
		Pods(q.Namespace).
		List(context.TODO(), listOptions)
	if err != nil {
		logrus.Error(err)
		render.SetError(handlers.CodeErrApp, err)
		return
	}
	var wg = new(sync.WaitGroup)
	for k, v := range res.Items {
		wg.Add(1)
		go func(k int, v coreV1.Pod) {
			defer wg.Done()
			if res.Items[k].ObjectMeta.Annotations == nil {
				res.Items[k].ObjectMeta.Annotations = make(map[string]string)
			}
			pb := internal.PodBase{
				Namespace: q.Namespace,
				PodName:   v.Name,
			}
			res.Items[k].ObjectMeta.Annotations["os"], res.Items[k].ObjectMeta.Annotations["arch"] =
				pb.OsAndArch(v.Spec.NodeName)
		}(k, v)
	}
	wg.Wait()
	render.SetJson(res)
}
