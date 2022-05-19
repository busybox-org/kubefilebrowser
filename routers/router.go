package routers

import (
	_ "embed"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/xmapst/kubefilebrowser/configs"
	_ "github.com/xmapst/kubefilebrowser/docs"
	"github.com/xmapst/kubefilebrowser/handlers/filebrowser"
	"github.com/xmapst/kubefilebrowser/handlers/kubeapiproxy"
	"github.com/xmapst/kubefilebrowser/routers/middleware"
)

func Router() *gin.Engine {
	router := gin.New()
	// 设置文件上传大小限制为8G
	router.MaxMultipartMemory = 32 << 20
	// middleware
	router.Use(
		cors.Default(),
		gin.Recovery(),
		gzip.Gzip(gzip.DefaultCompression),
		middleware.NoCache(),
		middleware.DenyMiddleware(),
		middleware.RequestIDMiddleware(),
	)

	// prometheus metrics
	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")
	m.SetSlowTime(10)
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	m.Use(router)

	// swagger doc
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// api
	apiGroup := router.Group("/api")
	{
		k8sGroup := apiGroup.Group("/kubeapiproxy")
		{
			k8sGroup.GET("/namespace", kubeapiproxy.Namespace)
			k8sGroup.GET("/pods", kubeapiproxy.Pods)
			k8sGroup.GET("/pvcs", kubeapiproxy.Pvcs)
			k8sGroup.POST("/upload", kubeapiproxy.Upload)
			k8sGroup.POST("/multiupload", kubeapiproxy.MultiUpload)
			k8sGroup.POST("/uploadpvc", kubeapiproxy.UploadPVC)
			k8sGroup.GET("/download", kubeapiproxy.Download)
			k8sGroup.GET("/terminal", kubeapiproxy.Terminal)
			k8sGroup.GET("/exec", kubeapiproxy.Exec)
		}
		fileBrowserGroup := apiGroup.Group("/filebrowser")
		{
			fileBrowserGroup.GET("/list", filebrowser.ListFile)
			fileBrowserGroup.GET("/open", filebrowser.OpenFile)
			fileBrowserGroup.POST("/createfile", filebrowser.CreateFile)
			fileBrowserGroup.POST("/createdir", filebrowser.CreateDir)
			fileBrowserGroup.POST("/rename", filebrowser.Rename)
			fileBrowserGroup.POST("/remove", filebrowser.Remove)
		}
	}
	if configs.Config.RunMode != gin.DebugMode {
		return router
	}
	// debug
	pprof.Register(router)
	return router
}
