package filebrowser

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xmapst/kubefilebrowser/handlers"
)

// OpenFile
// @Summary OpenFile
// @description 容器文件浏览器 - 打开文件
// @Tags FileBrowser
// @Param namespace query Query true "namespace"
// @Param pod query Query true "Pod名称"
// @Param container query Query true "容器名称"
// @Param path query Query true "路径"
// @Success 200 {object} handlers.JSONResult
// @Failure 500 {object} handlers.JSONResult
// @Router /api/filebrowser/open [get]
func OpenFile(c *gin.Context) {
	render := handlers.Gin{Context: c}
	var query = &Query{}
	if err := c.ShouldBindQuery(query); err != nil {
		logrus.Error(err)
		render.SetError(handlers.CodeErrParam, err)
		return
	}
	query.Command = []string{"/kftools", "cat", query.Path}
	bs, err := query.fileBrowser()
	if err != nil {
		render.SetError(handlers.CodeErrParam, err)
		return
	}
	render.SetJson(string(bs))
}
