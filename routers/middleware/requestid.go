package middleware

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

const xRequestIDKey = "X-Request-ID"

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		u4 := uuid.NewV4()
		xRequestID := u4.String()
		c.Request.Header.Set(xRequestIDKey, xRequestID)
		c.Writer.Header().Set(xRequestIDKey, xRequestID)
		c.Set(xRequestIDKey, xRequestID)
		c.Next()
	}
}
