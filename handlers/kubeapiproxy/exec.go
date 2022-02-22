package kubeapiproxy

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xmapst/kubefilebrowser/handlers"
	"github.com/xmapst/kubefilebrowser/internal"
)

type ExecQuery struct {
	Namespace string   `json:"namespace" form:"namespace" binding:"required"`
	Pod       string   `json:"pod" form:"pod" binding:"required"`
	Container string   `json:"container" form:"container" binding:"required"`
	Stdout    bool     `json:"stdout" form:"stdout"`
	Stderr    bool     `json:"stderr" form:"stderr"`
	Tty       bool     `json:"tty" form:"tty"`
	Commands  []string `json:"commands" form:"commands"`
}

// Exec
// @Summary Exec
// @description 在pod的容器中执行
// @Tags Kubernetes Api Proxy
// @Param namespace query ExecQuery true "namespace"
// @Param pods query ExecQuery true "Pod名称"
// @Param container query ExecQuery true "容器名称"
// @Param command query ExecQuery true "命令"
// @Param stdout query ExecQuery false "标准输出"
// @Param stderr query ExecQuery false "错误输出"
// @Param tty query ExecQuery false "终端"
// @Success 200 {object} handlers.JSONResult
// @Failure 500 {object} handlers.JSONResult
// @Router /api/kubeapiproxy/exec [get]
func Exec(c *gin.Context) {
	render := handlers.Gin{Context: c}
	var q = &ExecQuery{
		Stdout:   true,
		Stderr:   true,
		Tty:      false,
		Commands: []string{"ls", "-lQ", "--color=never", "--full-time", "/"},
	}
	if err := c.ShouldBindQuery(q); err != nil {
		render.SetError(handlers.CodeErrParam, err)
		return
	}
	pb := internal.PodBase{
		Namespace: q.Namespace,
		PodName:   q.Pod,
		Container: q.Container,
	}
	exec := pb.NewPodExec()
	exec.Command = q.Commands
	exec.Tty = q.Tty
	var stdout, stderr bytes.Buffer
	if q.Stdout {
		exec.Stdout = &stdout
	}
	if q.Stderr {
		exec.Stderr = &stderr
	}
	if err := exec.Exec(); err != nil {
		logrus.Error(err)
		render.SetError(handlers.CodeErrApp, err)
		return
	}
	render.SetJson(map[string]interface{}{
		"err": stderr.String(),
		"out": stdout.String(),
	})
}
