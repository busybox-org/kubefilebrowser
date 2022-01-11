package filebrowser

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
    "github.com/xmapst/kubefilebrowser/handlers"
)

// Rename
// @Summary Rename
// @description 容器文件浏览器 - 重命名
// @Tags FileBrowser
// @Param namespace query Query true "namespace"
// @Param pod query Query true "Pod名称"
// @Param container query Query true "容器名称"
// @Param path query Query true "新路径"
// @Param old_path query Query true "旧路径"
// @Success 200 {object} handlers.JSONResult
// @Failure 500 {object} handlers.JSONResult
// @Router /api/filebrowser/rename [post]
func Rename(c *gin.Context) {
	render := handlers.Gin{Context: c}
	var query = &Query{}
	if err := c.ShouldBindQuery(query); err != nil {
		logrus.Error(err)
		render.SetError(handlers.CodeErrParam, err)
		return
	}
	if query.OldPath == "" {
		render.SetError(handlers.CodeErrParam, fmt.Errorf("file path does not exist"))
		return
	}
	query.Command = []string{"/kftools", "mv", query.OldPath, query.Path}
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
