package kubeapiproxy

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
    "github.com/xmapst/kubefilebrowser/configs"
    "github.com/xmapst/kubefilebrowser/handlers"
    "github.com/xmapst/kubefilebrowser/internal"
    "os"
    "path/filepath"
    "strings"
    "sync"
    "time"
)

type MultiUploadQuery struct {
	Namespace string   `json:"namespace" form:"namespace" binding:"required"`
	Pods      []string `json:"pods" form:"pods" binding:"required"`
	DestPath  string   `json:"dest_path" form:"dest_path" binding:"required"`
}

// MultiUpload
// @Summary MultiUpload
// @description 批量上传到容器
// @Tags Kubernetes Api Proxy
// @Accept multipart/form-data
// @Param namespace query MultiUploadQuery true "namespace" default(default)
// @Param pods query MultiUploadQuery true "pods" default(nginx-test-76996486df)
// @Param dest_path query MultiUploadQuery false "dest_path" default(/root/)
// @Param files formData file true "files"
// @Success 200 {object} handlers.JSONResult
// @Failure 500 {object} handlers.JSONResult
// @Router /api/kubeapiproxy/multiupload [post]
func MultiUpload(c *gin.Context) {
	render := handlers.Gin{Context: c}
	var q MultiUploadQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		logrus.Error(err)
		render.SetError(handlers.CodeErrParam, err)
		return
	}

	// strip trailing slash (if any)
	if q.DestPath != "/" && strings.HasSuffix(string(q.DestPath[len(q.DestPath)-1]), "/") {
		q.DestPath = q.DestPath[:len(q.DestPath)-1]
	}

	var srcPath = filepath.Join(configs.TmpPath, fmt.Sprintf("%d", time.Now().UnixNano()))
	defer func() {
		_ = os.RemoveAll(srcPath)
	}()

	if err := render.SaveToTarFile(srcPath); err != nil {
		logrus.Error(err)
		render.SetError(handlers.CodeErrMsg, err)
		return
	}

	var res []*copyRes
	var wg = new(sync.WaitGroup)
	var lock = new(sync.Mutex)
	for _, pod := range q.Pods {
		wg.Add(1)
		go func(podName string) {
			defer wg.Done()
			var containerSlice []string
			pb := internal.PodBase{
				Namespace: q.Namespace,
				PodName:   podName,
			}
			_pod, err := pb.PodInfo()
			if err != nil {
				logrus.Error(err)
				return
			}
			for _, container := range _pod.Spec.Containers {
				containerSlice = append(containerSlice, container.Name)
			}
			_copy := &copyReq{
				namespace: q.Namespace,
				pod:       podName,
				srcPath:   srcPath,
				destPath:  q.DestPath,
			}
			lock.Lock()
			defer lock.Unlock()
			res = append(res, _copy.toPodContainers(containerSlice)...)
		}(pod)
	}
	wg.Wait()
	render.SetJson(res)
}

type UploadQuery struct {
	Namespace  string   `json:"namespace" form:"namespace" binding:"required"`
	Pod        string   `json:"pod" form:"pod" binding:"required"`
	Containers []string `json:"containers" form:"containers"`
	DestPath   string   `json:"dest_path" form:"dest_path" binding:"required"`
}

// Upload
// @Summary Upload
// @description 上传到容器
// @Tags Kubernetes Api Proxy
// @Accept multipart/form-data
// @Param namespace query UploadQuery true "namespace" default(default)
// @Param pod query UploadQuery true "pod" default(nginx-test-76996486df)
// @Param containers query UploadQuery true "containers" default(nginx-0)
// @Param dest_path query UploadQuery false "dest_path" default(/root/)
// @Param files formData file true "files"
// @Success 200 {object} handlers.JSONResult
// @Failure 500 {object} handlers.JSONResult
// @Router /api/kubeapiproxy/upload [post]
func Upload(c *gin.Context) {
	render := handlers.Gin{Context: c}
	var q UploadQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		logrus.Error(err)
		render.SetError(handlers.CodeErrParam, err)
		return
	}
	// strip trailing slash (if any)
	if q.DestPath != "/" && strings.HasSuffix(string(q.DestPath[len(q.DestPath)-1]), "/") {
		q.DestPath = q.DestPath[:len(q.DestPath)-1]
	}
	var srcPath = filepath.Join(configs.TmpPath, fmt.Sprintf("%d", time.Now().UnixNano()))
	defer func() {
		_ = os.RemoveAll(srcPath)
	}()

	if err := render.SaveToTarFile(srcPath); err != nil {
		logrus.Error(err)
		render.SetError(handlers.CodeErrMsg, err)
		return
	}

	var containerSlice []string
	if len(q.Containers) != 0 {
		containerSlice = q.Containers
	} else {
		pb := internal.PodBase{
			Namespace: q.Namespace,
			PodName:   q.Pod,
		}
		res, err := pb.PodInfo()
		if err != nil {
			logrus.Error(err)
			render.SetError(handlers.CodeErrApp, err)
			return
		}
		for _, container := range res.Spec.Containers {
			containerSlice = append(containerSlice, container.Name)
		}
	}
	_copy := &copyReq{
		namespace: q.Namespace,
		pod:       q.Pod,
		srcPath:   srcPath,
		destPath:  q.DestPath,
	}
	res := _copy.toPodContainers(containerSlice)
	render.SetJson(res)
}
