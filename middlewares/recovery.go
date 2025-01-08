package middlewares

import (
	"github.com/sosumecho/modules/exception"
	"github.com/sosumecho/modules/i18n"
	"github.com/sosumecho/modules/logger"
	"github.com/sosumecho/modules/response"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery(logf *logger.Logger, locale *i18n.I18N) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				httpRequest, _ := httputil.DumpRequest(c.Request, true)

				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						errStr := strings.ToLower(se.Error())
						if strings.Contains(errStr, "broken pipe") || strings.Contains(errStr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				if brokenPipe {
					logf.Error(c.Request.URL.Path,
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Error(err.(error))
					c.Abort()
					return
				}

				logf.Error("recovery from panic",
					zap.Time("time", time.Now()),
					zap.Any("error", err),
					zap.String("request", string(httpRequest)),
					zap.Stack("stacktrace"),
				)
				response.New(c, locale, logf).Error(http.StatusInternalServerError, exception.NewSystemError(nil))
			}
		}()
		c.Next()
	}
}
