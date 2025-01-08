package middlewares

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/sosumecho/modules/encrypt"
	"github.com/sosumecho/modules/logger"
	"go.uber.org/zap"
)

func Encrypt(key string, logf *logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		isEncrypt := ctx.GetHeader("x-encrypt")
		if isEncrypt == "false" {
			ctx.Next()
			return
		}
		w := &responseBodyWriter{
			body:           &bytes.Buffer{},
			ResponseWriter: ctx.Writer,
		}
		ctx.Writer = w
		ctx.Next()
		// 取出body后重新加密
		logf.Debug("encrypt", zap.String("body", w.body.String()))
		data := encrypt.NewXOR().SetKey([]byte(key)).Encrypt(w.body.Bytes())
		ctx.Header("x-encrypt", "true")
		_, err := ctx.Writer.Write(data)
		if err != nil {
			logf.Error("encrypt body", zap.Error(err))
			ctx.Abort()
		}
	}
}
