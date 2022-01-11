package filebrowser

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xmapst/kubefilebrowser/handlers"
	"github.com/xmapst/kubefilebrowser/utils"
)

// ListFile
// @Summary ListFile
// @description 容器文件浏览器 - 文件列表
// @Tags FileBrowser
// @Param namespace query Query true "namespace"
// @Param pod query Query true "Pod名称"
// @Param container query Query true "容器名称"
// @Param path query Query true "路径"
// @Success 200 {object} handlers.JSONResult
// @Failure 500 {object} handlers.JSONResult
// @Router /api/filebrowser/list [get]
func ListFile(c *gin.Context) {
	render := handlers.Gin{Context: c}
	var query = &Query{}
	if err := c.ShouldBindQuery(query); err != nil {
		logrus.Error(err)
		render.SetError(handlers.CodeErrParam, err)
		return
	}
	query.Command = []string{"/kftools", "ls", query.Path}
	bs, err := query.fileBrowser()
	if err != nil {
		render.SetError(handlers.CodeErrParam, err)
		return
	}
	var res []utils.File
	if err := json.Unmarshal(bs, &res); err != nil {
		logrus.Error(err)
		render.SetError(handlers.CodeErrApp, err)
		return
	}
	render.SetJson(res)
}
