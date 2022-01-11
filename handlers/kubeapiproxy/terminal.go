package kubeapiproxy

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xmapst/kubefilebrowser/configs"
	"github.com/xmapst/kubefilebrowser/handlers"
	"github.com/xmapst/kubefilebrowser/utils"
	"github.com/xmapst/kubefilebrowser/utils/terminal"
	"k8s.io/client-go/tools/remotecommand"
)

type TerminalQuery struct {
	Namespace string `json:"namespace" form:"namespace" binding:"required"`
	Pod       string `json:"pod" form:"pod" binding:"required"`
	Container string `json:"container" form:"container" binding:"required"`
	Shell     string `json:"shell" form:"shell"`
}

// Terminal
// @Summary Container terminal
// @description pod 中容器的终端
// @Tags Kubernetes Api Proxy
// @Param namespace query TerminalQuery true "namespace" default(default)
// @Param pod query TerminalQuery true "Pod名称"
// @Param container query TerminalQuery true "容器名称"
// @Param shell query TerminalQuery false "shell" default(sh[bash/sh/cmd])
// @Success 200 {object} handlers.JSONResult
// @Failure 500 {object} handlers.JSONResult
// @Router /api/kubeapiproxy/terminal [get]
func Terminal(c *gin.Context) {
	render := handlers.Gin{Context: c}
	var q = &TerminalQuery{
		Shell: "sh",
	}
	if err := c.ShouldBindQuery(q); err != nil {
		render.SetError(handlers.CodeErrParam, err)
		return
	}

	wsConn, err := utils.InitWebsocket(c.Writer, c.Request)
	if utils.WsHandleError(wsConn, err) {
		logrus.Error(err)
		return
	}
	defer wsConn.WsClose()
	webTerminal := terminal.WebTerminal{
		K8sClient: configs.RestClient,
		Namespace: q.Namespace,
		PodName:   q.Pod,
		Container: q.Container,
		Shell:     q.Shell,
	}
	SshSPDYExecutor := webTerminal.NewSshSPDYExecutor()
	executor, err := remotecommand.NewSPDYExecutor(configs.KuBeResConf, "POST", SshSPDYExecutor.URL())
	if utils.WsHandleError(wsConn, err) {
		logrus.Error(err)
		return
	}
	handler := &terminal.StreamHandler{
		SessionID:   c.Request.Header.Get("X-Request-Id"),
		WsConn:      wsConn,
		ResizeEvent: make(chan remotecommand.TerminalSize),
		Shell:       q.Shell,
	}
	go handler.CommandRecordChan()
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             handler,
		Stdout:            handler,
		Stderr:            handler,
		Tty:               true,
		TerminalSizeQueue: handler,
	})
	if utils.WsHandleError(wsConn, err) {
		logrus.Error(err)
	}
}
