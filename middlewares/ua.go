package middlewares

import (
	"errors"
	"github.com/sosumecho/modules/exception"
	"github.com/sosumecho/modules/i18n"
	"github.com/sosumecho/modules/logger"
	"github.com/sosumecho/modules/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ForceUA(locale *i18n.I18N, logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("User-Agent") == "" {
			response.New(c, locale, logger).Error(http.StatusBadRequest, exception.NewUAError(errors.New("invalid user-agent")))
			return
		}
		c.Next()
	}
}
