package leowebgin

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// 中间件

// timeout middleware wraps the request context with a timeout
func (engine *Engine) TimeoutMiddleware(timeout time.Duration) func(c *Context) {
	return func(cont *Context) {
		c := cont.context
		// wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer func() {
			// check if context timeout was reached
			if ctx.Err() == context.DeadlineExceeded {
				// write response and abort the request
				c.Writer.WriteHeader(http.StatusGatewayTimeout)
				engine.log(engine.isProduction(), formatRequestLog(c))
				c.Abort()
			}
			// cancel to clear resources after finished
			cancel()
		}()
		// replace request with context wrapped request
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func (engine *Engine) recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		c.AbortWithStatus(http.StatusInternalServerError)
		engine.log(engine.isProduction(),
			fmt.Sprintf("%s err: %v |", formatRequestLog(c), err))
	})
}

func (engine *Engine) middleLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		if engine.log == nil {
			return
		}
		engine.log(engine.runMode == Pro, formatRequestLog(c))
		c.Next()
		engine.log(engine.runMode == Pro, formatRequestLog(c))
	}
}

func cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,"+
			"X-CSRF-Token,Authorization,Token,X-Token,"+
			"Referer,X-Requested-With")
		context.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
		context.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin"+
			",Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			//context.AbortWithStatus(http.StatusNoContent)
			context.JSON(http.StatusOK, "options request")
		}
		context.Next()
	}
}
