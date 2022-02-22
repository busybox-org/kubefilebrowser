package kubeapiproxy

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
    "github.com/xmapst/kubefilebrowser/handlers"
    "github.com/xmapst/kubefilebrowser/internal"
    "io"
    "strconv"
    "strings"
    "time"
)

type DownloadQuery struct {
	Namespace string   `json:"namespace" form:"namespace" binding:"required"`
	Pod       string   `json:"pod" form:"pod" binding:"required"`
	Container string   `json:"container" form:"container" binding:"required"`
	DestPaths []string `json:"dest_paths" form:"dest_paths" binding:"required"`
	Style     string   `json:"style" form:"style"`
}

// Download
// @Summary Download
// @description 从容器下载到本地
// @Tags Kubernetes Api Proxy
// @Accept json
// @Param namespace query DownloadQuery true "namespace" default(default)
// @Param pod query DownloadQuery true "pod" default(nginx-test-76996486df)
// @Param container query DownloadQuery true "container" default(nginx-0)
// @Param dest_paths query DownloadQuery true "dest_paths" default(/root)
// @Param style query DownloadQuery true "style" default(rar)
// @Success 200 {object} handlers.JSONResult
// @Failure 500 {object} handlers.JSONResult
// @Router /api/kubeapiproxy/download [get]
func Download(c *gin.Context) {
	render := handlers.Gin{Context: c}
	var q = DownloadQuery{
		Style: "rar",
	}
	if err := c.ShouldBindQuery(&q); err != nil {
		logrus.Error(err)
		render.SetError(handlers.CodeErrApp, err)
		return
	}

	fileName := fmt.Sprintf("%s.%s", strconv.FormatInt(time.Now().UnixNano(), 10), q.Style)
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("X-File-Name", fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	pb := internal.PodBase{
		Namespace: q.Namespace,
		PodName:   q.Pod,
		Container: q.Container,
	}
	cp := pb.NewPodExec()
	cp.Stdout = render.Writer
	switch q.Style {
	case "tar":
		err := cp.FromPodContainer(q.DestPaths, q.Style)
		if err != nil {
			logrus.Error(err)
			render.SetError(handlers.CodeErrApp, err)
			return
		}
	case "zip":
		reader, writer := io.Pipe()
		cp.Stdin = reader
		go func() {
			<-c.Request.Context().Done()
			_, _ = writer.Write([]byte("close\n"))
		}()
		err := cp.FromPodContainer(q.DestPaths, q.Style)
		if err != nil {
			if strings.Contains(err.Error(), "kftools") ||
				err.Error() == "command terminated with exit code 126" {
				err = pb.InstallKFTools()
				if err != nil {
					logrus.Error(err)
					render.SetError(handlers.CodeErrApp, err)
					return
				}
			}
			err = cp.FromPodContainer(q.DestPaths, q.Style)
			if err != nil {
				logrus.Error(err)
				render.SetError(handlers.CodeErrApp, err)
				return
			}
		}
	default:
		render.SetError(handlers.CodeErrMsg, fmt.Errorf("no matching compression type found"))
		return
	}
}
