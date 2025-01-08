package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/sosumecho/modules/exception"
	"github.com/sosumecho/modules/exceptions"
	"github.com/sosumecho/modules/i18n"
	"github.com/sosumecho/modules/logger"
	"github.com/sosumecho/modules/response"
	"go.uber.org/zap"
)

type Handler[T any, U any] func(ctx *gin.Context, param T) (U, exception.Exception)

func ResponseWrapper[T any, U any](handler Handler[T, U], hasQuery bool, locale *i18n.I18N, logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := response.New(c, locale, logger)
		var param T
		if err := c.ShouldBind(&param); err != nil {
			logger.Warn("params validate", zap.Error(err))
			resp.Fail(exception.NewParamsError(exceptions.ParamsError))
			return
		}
		if hasQuery {
			if err := c.ShouldBind(&param); err != nil {
				logger.Warn("params validate", zap.Error(err))
				resp.Fail(exception.NewParamsError(exceptions.ParamsError))
				return
			}
		}
		result, err := handler(c, param)
		if err != nil {
			resp.Fail(err)
			return
		}
		resp.Data(result)
	}
}
