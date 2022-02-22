package kubeapiproxy

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xmapst/kubefilebrowser/configs"
	"github.com/xmapst/kubefilebrowser/handlers"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceSQuery struct {
	Namespace     string `json:"namespace" form:"namespace"`
	FieldSelector string `json:"field_selector" form:"field_selector"`
	LabelSelector string `json:"label_selector" form:"label_selector"`
}

// Namespace
// @Summary Namespace
// @description 命名空间列表
// @Tags Kubernetes Api Proxy
// @Param namespace query NamespaceSQuery false "namespace"
// @Param field_selector query NamespaceSQuery false "field_selector"
// @Param label_selector query NamespaceSQuery false "label_selector"
// @Success 200 {object} handlers.JSONResult
// @Failure 500 {object} handlers.JSONResult
// @Router /api/kubeapiproxy/namespace [get]
func Namespace(c *gin.Context) {
	render := handlers.Gin{Context: c}
	var q NamespaceSQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		logrus.Error(err)
		render.SetError(handlers.CodeErrParam, err)
		return
	}
	if _, ok := c.GetQuery("namespace"); !ok {
		res, err := configs.RestClient.CoreV1().Namespaces().
			List(context.TODO(), metaV1.ListOptions{
				LabelSelector: q.LabelSelector,
				FieldSelector: q.FieldSelector,
			})
		if err != nil {
			logrus.Error(err)
			render.SetError(handlers.CodeErrApp, err)
			return
		}
		render.SetJson(res)
		return
	}
	res, err := configs.RestClient.CoreV1().Namespaces().
		Get(context.TODO(), q.Namespace, metaV1.GetOptions{})
	if err != nil {
		logrus.Error(err)
		render.SetError(handlers.CodeErrApp, err)
		return
	}
	render.SetJson(res)
}
