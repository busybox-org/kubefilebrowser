package filebrowser

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xmapst/kubefilebrowser/handlers"
)

// CreateFile
// @Summary CreateFile
// @description 容器文件浏览器 - 创建文件
// @Tags FileBrowser
// @Param namespace query Query true "namespace"
// @Param pods query Query true "Pod名称"
// @Param container query Query true "容器名称"
// @Param path query Query true "路径"
// @Param content body string false "内容"
// @Success 200 {object} handlers.JSONResult
// @Failure 500 {object} handlers.JSONResult
// @Router /api/filebrowser/createfile [post]
func CreateFile(c *gin.Context) {
	render := handlers.Gin{Context: c}
	var query = &Query{
		Path: "/",
	}
	if err := c.ShouldBindQuery(query); err != nil {
		logrus.Error(err)
		render.SetError(handlers.CodeErrParam, err)
		return
	}
	query.Command = []string{"/kftools", "touch", query.Path}
	query.Stdin = c.Request.Body
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

// CreateDir
// @Summary CreateDir
// @description 容器文件浏览器 - 创建目录
// @Tags FileBrowser
// @Param namespace query Query true "namespace"
// @Param pod query Query true "Pod名称"
// @Param container query Query true "容器名称"
// @Param path query Query true "路径"
// @Success 200 {object} handlers.JSONResult
// @Failure 500 {object} handlers.JSONResult
// @Router /api/filebrowser/createdir [post]
func CreateDir(c *gin.Context) {
	render := handlers.Gin{Context: c}
	var query = &Query{}
	if err := c.ShouldBindQuery(query); err != nil {
		logrus.Error(err)
		render.SetError(handlers.CodeErrParam, err)
		return
	}
	query.Command = []string{"/kftools", "mkdir", query.Path}
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
