package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sosumecho/modules/exception"
	"github.com/sosumecho/modules/i18n"
	"github.com/sosumecho/modules/logger"
	"github.com/sosumecho/modules/response"
)

// ParseRequest 解析请求
func ParseRequest(locale *i18n.I18N, logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// var err error
		// 得到平台
		platform := c.GetHeader("x-platform")
		// 得到设备号
		deviceID := c.GetHeader("x-device-id")
		// 得到手机型号和系统版本
		deviceModel := c.GetHeader("x-device-model")
		versionStr := c.GetHeader("x-version")
		bundleStr := c.GetHeader("x-bundle-id")

		if deviceID == "" || platform == "" || deviceModel == "" {
			err := errors.New("deviceID or platform or deviceModel is empty")
			paramError := exception.NewParamsError(err)
			response.New(c, locale, logger).Fail(paramError)
			c.Abort()
			return
		}

		c.Set("device_id", deviceID)
		c.Set("platform", platform)
		c.Set("device_model", deviceModel)
		c.Set("version", versionStr)
		c.Set("bundle_id", bundleStr)

		c.Next()
	}
}
