package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cors 允许跨域中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"))
		c.Header("Access-Control-Allow-Headers", "AccessToken,X-CSRF-Token, Authorization, Token, x-app-id, X-Requested-With, Content-Type, x-version, x-device-id,x-platform, x-channel, x-web, x-encrypt,x-mr")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		//放行所有OPTIONS方法，因为有的模板是要请求两次的
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}
