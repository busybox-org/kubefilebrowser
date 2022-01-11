package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xmapst/kubefilebrowser/configs"
	"github.com/xmapst/kubefilebrowser/handlers"
	"github.com/xmapst/kubefilebrowser/utils"
	"github.com/xmapst/kubefilebrowser/utils/denyip"
)

var (
	checker *denyip.Checker
	err     error
)

func init() {
	if !utils.InSliceString("*", configs.Config.IPWhiteList) && len(configs.Config.IPWhiteList) != 0 {
		checker, err = denyip.NewChecker(configs.Config.IPWhiteList)
	}
	if err != nil {
		logrus.Fatal(err)
	}
}

func DenyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		render := handlers.Gin{Context: c}
		reqIPAddr := denyip.GetRemoteIP(c.Request)
		if !utils.InSliceString("*", configs.Config.IPWhiteList) && len(configs.Config.IPWhiteList) != 0 {
			reeIPadLenOffset := len(reqIPAddr) - 1
			for i := reeIPadLenOffset; i >= 0; i-- {
				err = checker.IsAuthorized(reqIPAddr[i])
				if err != nil {
					logrus.Error(err)
					render.SetError(handlers.CodeErrNoPriv, err)
					return
				}
			}
		}
		c.Next()
	}
}
