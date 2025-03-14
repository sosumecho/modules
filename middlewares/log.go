package middlewares

import (
	"bytes"
	"encoding/base64"
	"github.com/sosumecho/modules/encrypt"
	"github.com/sosumecho/modules/logger"
	"github.com/sosumecho/modules/utils"
	"io"
	"strings"

	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type ResponseWriter interface {
	gin.ResponseWriter
	GetBody() *bytes.Buffer
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r responseBodyWriter) GetBody() *bytes.Buffer {
	return r.body
}

type encryptWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
	key  string
}

func (r encryptWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	data := encrypt.NewXOR().SetKey([]byte(r.key)).Encrypt(b)
	return r.ResponseWriter.Write([]byte(base64.StdEncoding.EncodeToString(data)))
}

func (r encryptWriter) GetBody() *bytes.Buffer {
	return r.body
}

func Logger(logf *logger.Logger, key string, excludes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var w ResponseWriter
		isEncrypt := c.GetHeader("x-encrypt")
		for _, exclude := range excludes {
			if strings.HasPrefix(strings.TrimPrefix(c.Request.URL.Path, "/api"), strings.TrimPrefix(exclude, "/api")) {
				isEncrypt = "false"
				break
			}
		}
		if isEncrypt == "false" || key == "" {
			w = &responseBodyWriter{
				body:           &bytes.Buffer{},
				ResponseWriter: c.Writer,
			}
		} else {
			w = &encryptWriter{
				body:           &bytes.Buffer{},
				ResponseWriter: c.Writer,
				key:            key,
			}
			isEncrypt = "true"
		}

		c.Writer = w

		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		start := time.Now()

		c.Header("x-encrypt", isEncrypt)
		c.Next()

		cost := time.Since(start)
		responseStatus := c.Writer.Status()

		logFields := []zap.Field{
			zap.Int("status", responseStatus),
			zap.String("request", fmt.Sprintf("[%s] %s", c.Request.Method, c.Request.URL)),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", utils.GetClientIP(c)),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.String("time", utils.MicrosecondsStr(cost)),
			zap.Any("headers", c.Request.Header),
		}

		if utils.InArray([]string{"POST", "PUT", "DELETE"}, c.Request.Method) && c.Request.MultipartForm == nil {
			if !(strings.Contains(c.Request.RequestURI, "upload") || strings.Contains(c.Request.RequestURI, "tinymce")) {
				logFields = append(logFields, zap.String("body", string(requestBody)))
			}
		}

		if strings.Contains(c.Request.Header.Get("Content-Type"), "application/json") {
			logFields = append(logFields, zap.Any("response", w.GetBody()))
		}

		if responseStatus > 400 && responseStatus < 500 {
			logf.Warn("HTTP Warning "+cast.ToString(responseStatus), logFields...)
		} else if responseStatus > 500 && responseStatus <= 599 {
			logf.Error("HTTP Error "+cast.ToString(responseStatus), logFields...)
		} else {
			logf.Debug("HTTP Access Log", logFields...)
		}

	}
}
