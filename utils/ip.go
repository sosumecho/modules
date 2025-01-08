package utils

import "github.com/gin-gonic/gin"

// GetClientIP 得到客户端IP
func GetClientIP(c *gin.Context) string {
	cfIP := c.GetHeader("Cf-Connecting-Ip")
	if cfIP != "" {
		return cfIP
	}
	realIP := c.GetHeader("X-Real-Ip")
	if realIP != "" {
		return realIP
	}
	return c.ClientIP()
}
