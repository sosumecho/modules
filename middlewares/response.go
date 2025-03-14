package middlewares

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/sosumecho/modules/encrypt"
	"github.com/sosumecho/modules/exception"
	"github.com/sosumecho/modules/exceptions"
	"github.com/sosumecho/modules/i18n"
	"github.com/sosumecho/modules/logger"
	"github.com/sosumecho/modules/response"
	"go.uber.org/zap"
	"reflect"
)

type Handler func(ctx *gin.Context, param any) (any, exception.Exception)

// ResponseWrapper 函数，保持通用性并自动推导参数类型
func ResponseWrapper(locale *i18n.I18N, logger *logger.Logger, isEncrypt bool, secretKey string) func(handler interface{}, hasQuery bool) gin.HandlerFunc {
	return func(handler interface{}, hasQuery bool) gin.HandlerFunc {
		// 获取 handler 的类型并推导出参数类型
		handlerValue := reflect.ValueOf(handler)
		if handlerValue.Kind() != reflect.Func {
			logger.Error("Expected a function handler")
			return nil
		}

		// 获取函数的参数类型
		handlerType := handlerValue.Type()
		paramType := handlerType.In(1) // 获取第2个参数的类型（即 T 类型）

		// 返回 gin.HandlerFunc
		return func(c *gin.Context) {
			resp := response.New(c, locale, logger)

			// 动态创建 param 类型的实例
			param := reflect.New(paramType).Interface()

			// 绑定请求体参数（表单或 JSON）
			if err := c.ShouldBind(param); err != nil {
				logger.Warn("params validate", zap.Error(err))
				resp.Fail(exception.NewParamsError(exceptions.ParamsError))
				return
			}

			// 如果需要绑定查询参数
			if hasQuery {
				if err := c.ShouldBindQuery(param); err != nil {
					logger.Warn("params validate", zap.Error(err))
					resp.Fail(exception.NewParamsError(exceptions.ParamsError))
					return
				}
			}

			// 构造传递给 handler 的参数
			args := []reflect.Value{reflect.ValueOf(c), reflect.ValueOf(param).Elem()} // 使用 .Elem() 获取指针的元素

			// 调用 handler 函数
			resultValues := handlerValue.Call(args)

			// 处理返回值
			if len(resultValues) == 2 {
				result := resultValues[0].Interface() // 第一个返回值
				err := resultValues[1].Interface()    // 第二个返回值（错误）

				// 判断错误
				if err != nil {
					resp.Fail(err.(exception.Exception))
					return
				}

				if isEncrypt {
					var e exception.Exception
					result, e = EncryptResponse(result, secretKey)
					if err != nil {
						resp.Fail(e)
					}
				}

				resp.Data(result)
			} else {
				resp.Fail(exception.NewSystemError(exceptions.SystemError))
			}
		}
	}
}

func EncryptResponse(rs any, key string) (string, exception.Exception) {
	b, e := jsoniter.Marshal(rs)
	if e != nil {
		return "", exception.NewSystemError(e)
	}
	data := encrypt.NewXOR().SetKey([]byte(key)).Encrypt(b)
	return base64.StdEncoding.EncodeToString(data), nil
}
