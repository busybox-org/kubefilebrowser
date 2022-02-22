package filebrowser

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xmapst/kubefilebrowser/handlers"
	"strings"
)

// Remove
// @Summary Remove
// @description 容器文件浏览器 - 删除
// @Tags FileBrowser
// @Param namespace query Query true "namespace"
// @Param pod query Query true "Pod名称"
// @Param container query Query true "容器名称"
// @Param path query Query true "路径"
// @Success 200 {object} handlers.JSONResult
// @Failure 500 {object} handlers.JSONResult
// @Router /api/filebrowser/remove [post]
func Remove(c *gin.Context) {
	render := handlers.Gin{Context: c}
	var query = &Query{}
	if err := c.ShouldBindQuery(query); err != nil {
		logrus.Error(err)
		render.SetError(handlers.CodeErrParam, err)
		return
	}
	query.Command = append([]string{"/kftools", "rm"}, strings.Split(query.Path, ",")...)
	bs, err := query.fileBrowser()
	if err != nil {
		render.SetError(handlers.CodeErrParam, err)
		return
	}
	if len(string(bs)) != 0 {
		render.SetJson(string(bs))
		return
	}
	render.SetJson("success")
}
