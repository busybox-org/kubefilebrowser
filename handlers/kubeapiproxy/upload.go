package kubeapiproxy

import (
	"context"
	"fmt"
	"github.com/avast/retry-go/v4"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xmapst/kubefilebrowser/configs"
	"github.com/xmapst/kubefilebrowser/handlers"
	"github.com/xmapst/kubefilebrowser/internal"
	coreV1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// UploadPods
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
func UploadPods(c *gin.Context) {
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

type UploadPVCQuery struct {
	Namespace string `json:"namespace" form:"namespace" binding:"required"`
	Pvc       string `json:"pvc" form:"pvc" binding:"required"`
}

var (
	newPodNode      = ""
	containerName   = "upload"
	containerImange = "busybox"
	destVolume      = "tmp-upload"
	destPath        = "/tmp-upload"
)

// UploadPVC
// @Summary Upload
// @description 上传到PVC
// @Tags Kubernetes Api Proxy
// @Accept multipart/form-data
// @Param namespace query UploadQuery true "namespace" default(default)
// @Param pvc query UploadQuery true "pvc" default(nginx-test-76996486df)
// @Param files formData file true "files"
// @Success 200 {object} handlers.JSONResult
// @Failure 500 {object} handlers.JSONResult
// @Router /api/kubeapiproxy/upload [post]
func UploadPVC(c *gin.Context) {
	var err error
	render := handlers.Gin{Context: c}
	var q UploadPVCQuery
	if err = c.ShouldBindQuery(&q); err != nil {
		logrus.Error(err)
		render.SetError(handlers.CodeErrParam, err)
		return
	}

	var newPodName = fmt.Sprintf("tmp-upload-pvc-%v", time.Now().Unix())
	//如果pvc的状态是readOnlyMany，则不允许上传
	pvc, err := configs.RestClient.CoreV1().PersistentVolumeClaims(q.Namespace).Get(context.TODO(), q.Pvc, v1.GetOptions{})
	if err != nil {
		logrus.Error(err)
		render.SetError(handlers.CodeErrMsg, err)
		return
	}
	if len(pvc.Spec.AccessModes) > 0 && pvc.Spec.AccessModes[0] == coreV1.ReadOnlyMany {
		err = fmt.Errorf("ReadOnlyMany pvc: Permission denied")
		logrus.Error(err)
		render.SetError(handlers.CodeErrMsg, err)
		return
	}

	//如果pvc已经挂载，找到挂载节点
	pods, err := configs.RestClient.CoreV1().Pods(q.Namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		logrus.Error(err)
		render.SetError(handlers.CodeErrMsg, err)
		return
	}
	for _, pod := range pods.Items {
		for _, v := range pod.Spec.Volumes {
			if v.PersistentVolumeClaim != nil && v.PersistentVolumeClaim.ClaimName == q.Pvc {
				osType, _ := internal.GetOsAndArch(pod.Spec.NodeName)
				if osType == "windows" {
					containerImange = "mcr.microsoft.com/windows/nanoserver"
				}
				newPodNode = pod.Spec.NodeName
				break
			}
		}
	}

	var srcPath = filepath.Join(configs.TmpPath, fmt.Sprintf("%d", time.Now().UnixNano()))
	defer func() {
		_ = os.RemoveAll(srcPath)
	}()

	if err = render.SaveToTarFile(srcPath); err != nil {
		logrus.Error(err)
		render.SetError(handlers.CodeErrMsg, err)
		return
	}
	//创建upload pod
	newPod := &coreV1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name:      newPodName,
			Namespace: q.Namespace,
		},
		Spec: coreV1.PodSpec{
			Containers: []coreV1.Container{
				{
					Name:            containerName,
					Image:           containerImange,
					ImagePullPolicy: coreV1.PullIfNotPresent,
					TTY:             true,
					VolumeMounts: []coreV1.VolumeMount{
						{
							Name:      destVolume,
							MountPath: destPath,
						},
					},
				},
			},
			Volumes: []coreV1.Volume{
				{
					Name: destVolume,
					VolumeSource: coreV1.VolumeSource{
						PersistentVolumeClaim: &coreV1.PersistentVolumeClaimVolumeSource{
							ClaimName: q.Pvc,
						},
					},
				},
			},
		},
	}

	if newPodNode != "" {
		newPod.Spec.NodeName = newPodNode
	}

	_, err = configs.RestClient.CoreV1().Pods(q.Namespace).Create(context.TODO(), newPod, v1.CreateOptions{})
	if err != nil {
		logrus.Errorf("create upload pod err: %s", err.Error())
		render.SetError(handlers.CodeErrMsg, err)
		return
	}

	// 等待创建pod ready
	err = retry.Do(
		func() error {
			pb := internal.PodBase{
				Namespace: q.Namespace,
				PodName:   newPodName,
			}
			res, err := pb.PodInfo()
			if err != nil {
				return err
			}
			if res.Status.Phase == coreV1.PodRunning {
				return nil
			}
			return fmt.Errorf("pod not ready")
		},
		retry.Attempts(7),
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			max := time.Duration(n)
			if max > 8 {
				max = 8
			}
			duration := time.Second * max * max
			return duration
		}),
		retry.MaxDelay(time.Second*60),
	)
	if err != nil {
		err = fmt.Errorf("upload pod status not running: %s", err.Error())
		logrus.Error(err)
		errDel := configs.RestClient.CoreV1().Pods(q.Namespace).Delete(context.TODO(), newPodName, v1.DeleteOptions{})
		if errDel != nil {
			logrus.Errorf("delete upload pod err: %s", err.Error())
		}
		render.SetError(handlers.CodeErrApp, err)
		return
	}
	// copy
	_copy := &copyReq{
		namespace: q.Namespace,
		pod:       newPodName,
		srcPath:   srcPath,
		destPath:  destPath,
	}
	res := _copy.toPodContainer(containerName)
	//pod销毁
	err = configs.RestClient.CoreV1().Pods(q.Namespace).Delete(context.TODO(), newPodName, v1.DeleteOptions{})
	if err != nil {
		logrus.Errorf("delete upload pod err: %s", err.Error())
		render.SetError(handlers.CodeErrMsg, err)
		return
	}

	render.SetJson(res)
}
