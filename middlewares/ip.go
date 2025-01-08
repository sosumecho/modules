package middlewares

import (
	"github.com/gin-gonic/gin"
)

func IP() gin.HandlerFunc {
	return func(c *gin.Context) {
		//ip := c.ClientIP()
		//if !config.NewService().IsChinaCanVisit() && utils.IP2Country(ip) == "中国" {
		//	response.New().SetErrorCode(response.CodeSystemError).APIFail(c, "Not allowed in this area")
		//	c.Abort()
		//	return
		//}
		c.Next()
	}
}
