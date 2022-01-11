package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xmapst/kubefilebrowser"
	"github.com/xmapst/kubefilebrowser/configs"
	"github.com/xmapst/kubefilebrowser/routers"
	_ "github.com/xmapst/kubefilebrowser/routers"
	"github.com/xmapst/kubefilebrowser/utils"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var webTimeout time.Duration

func init() {
	// log format init
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&utils.ConsoleFormatter{})
	kingpin.Version(kubefilebrowser.VersionIfo())
	kingpin.HelpFlag.Short('h')
	kingpin.Flag(
		"timeout",
		"Timeout for calling endpoints on the engine",
	).Default("600000s").DurationVar(&webTimeout)
	kingpin.Parse()
}

// @title Kube-FileBrowser Swagger
// @version 1.0
// @description Kubernetes FileBrowser
// @BasePath /
// @query.collection.format multi
func main() {
	kubefilebrowser.PrintHeadInfo()
	configs.LoadConfig()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	gin.SetMode(configs.Config.RunMode)
	logrus.Info("start up...")
	r := routers.Router()
	r.Use(kubefilebrowser.StaticFile("/"))
	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%s", configs.Config.HTTPAddr, configs.Config.HTTPPort),
		// Good practice to set timeouts to avoid Solaris attacks.
		WriteTimeout: webTimeout,
		ReadTimeout:  webTimeout,
		IdleTimeout:  webTimeout,
		Handler:      r, // Pass our instance of gin in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		logrus.Infof("listen address [%s:%s]", configs.Config.HTTPAddr, configs.Config.HTTPPort)
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				logrus.Fatalln(err)
			}
		}
	}()
	logrus.Info("kubernetes file browser is running ...")

	<-signals
	logrus.Info("shutdown server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	_ = srv.Shutdown(ctx)
	os.Exit(0)
}
