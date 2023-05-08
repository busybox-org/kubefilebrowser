package kubefilebrowser

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"path"
	"strings"
)

//go:embed static/*
var staticFS embed.FS

const fsBase = "static"

// StaticFile 使用go.16新的特性embed 到包前端编译后的代码. 替代nginx.   one binary rules them all
func StaticFile(urlPrefix string) gin.HandlerFunc {
	const indexHtml = "index.html"
	return func(c *gin.Context) {
		urlPath := strings.TrimSpace(c.Request.URL.Path)
		if urlPath == urlPrefix {
			urlPath = path.Join(urlPrefix, indexHtml)
		}
		urlPath = path.Join(fsBase, urlPath)
		f, err := staticFS.Open(urlPath)
		if err != nil {
			//NoRoute
			bs, err := staticFS.ReadFile(path.Join(fsBase, "/", indexHtml))
			if err != nil {
				logrus.Error(err, "embed fs")
				return
			}
			c.Status(200)
			_, _ = c.Writer.Write(bs)
			c.Abort()
			return
		}
		fi, err := f.Stat()
		if strings.HasSuffix(urlPath, ".html") {
			c.Header("Cache-Control", "no-cache")
			c.Header("Content-Type", "text/html; charset=utf-8")
		}

		if strings.HasSuffix(urlPath, ".js") {
			c.Header("Content-Type", "text/javascript; charset=utf-8")
		}
		if strings.HasSuffix(urlPath, ".css") {
			c.Header("Content-Type", "text/css; charset=utf-8")
		}

		if err != nil || !fi.IsDir() {
			bs, err := staticFS.ReadFile(urlPath)
			if err != nil {
				logrus.Error(err, "embed fs")
				return
			}
			c.Status(200)
			_, _ = c.Writer.Write(bs)
			c.Abort()
		}
	}
}
