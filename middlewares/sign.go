package middlewares

//import (
//	"bytes"
//	"io"
//	"github.com/sosumecho/modules/exception"
//	"github.com/sosumecho/modules/i18n"
//	"github.com/sosumecho/modules/logger"
//	"github.com/sosumecho/modules/response"
//	"go.uber.org/zap"
//	"log"
//	"github.com/sosumecho/modules/sign"
//	"seal/config"
//
//	"errors"
//
//	"io/ioutil"
//	"strconv"
//	"strings"
//
//	"github.com/sirupsen/logrus"
//
//	jsoniter "github.com/json-iterator/go"
//
//	"github.com/gin-gonic/gin"
//)
//
//// CheckSign 校验签名
//func CheckSign(logf *logger.Logger, locale *i18n.I18N, exclude ...string) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// var err error
//		// 得到平台
//		platform := c.GetHeader("x-platform")
//		// 得到设备号
//		deviceID := c.GetHeader("x-device-id")
//		// 得到过期时间
//		timestamp := c.GetHeader("x-timestamp")
//		// 得到手机型号和系统版本
//		deviceModel := c.GetHeader("x-device-model")
//		//wantSign := c.GetHeader("x-want-sign")
//		// 签名字符串
//		signStr := c.GetHeader("sign")
//		versionStr := c.GetHeader("x-version")
//		bundleStr := c.GetHeader("x-bundle-id")
//		channelStr := c.GetHeader("x-channel")
//		language := strings.ToLower(c.GetHeader("x-language"))
//		data := make(map[string]interface{})
//
//		contentType := c.Request.Header.Get("Content-Type")
//		// log.New().Info("content-type: ", contentType)
//		if strings.Contains(contentType, "json") {
//			b, err := io.ReadAll(c.Request.Body)
//			c.Request.Body = io.NopCloser(bytes.NewBuffer(b))
//			if err != nil {
//				logf.Error("读取json请求体中body内容失败,", zap.Error(err))
//				c.Abort()
//				return
//			}
//			if len(b) > 0 {
//				var d map[string]interface{}
//				if err = jsoniter.Unmarshal(b, &d); err != nil {
//					logf.Error("反序列化body内容失败, ", zap.Error(err))
//					c.Abort()
//					return
//				}
//				for k, v := range d {
//					data[k] = v
//				}
//			}
//
//		} else {
//			c.Request.ParseMultipartForm(32 << 20)
//			for k, v := range c.Request.PostForm {
//				//fmt.Println(k, ":", v)
//				if len(v) > 0 {
//					data[k] = v[0]
//				}
//			}
//		}
//
//		//if deviceID == "" || platform == "" || deviceModel == "" || timestamp == "" || signStr == "" {
//		//	response.New().SetErrorCode(response.CodeParamsError).APIFail(c, "params error")
//		//	c.Abort()
//		//	return
//		//}
//		data["x-timestamp"] = timestamp
//		data["x-device-id"] = deviceID
//		data["x-platform"] = platform
//		data["x-device-model"] = deviceModel
//		data[config.Config.GetString("sign.default_key_name")] = signStr
//		//log.Println("签名头上的参数 : ", data)
//
//		// TODO: 这里后面要打开
//		// log.New().WithFields(logrus.Fields{
//		// 	"sign": sign.New(sign.CustomSignerType).SetKey(config.Config.GetString("sign.custom.key")).SetKeyName(config.Config.GetString("sign.custom.key_name")).Sign(data),
//		// 	"data": data,
//		// }).Debug("sign 加密")
//		isExclude := false
//		if len(exclude) > 0 {
//			for _, item := range exclude {
//				if item == c.Request.URL.Path {
//					isExclude = true
//					break
//				}
//			}
//		}
//		if config.Env != "dev" {
//			//	signPlatform := platform
//			//	if wantSign != "" {
//			//		signPlatform = wantSign
//			//	}
//			//key := services.GetEncryptKey(models.Sign, models.Platform(signPlatform))
//			////if key == "" {
//			key := config.Config.GetString("sign.custom.key")
//			//}
//			if !isExclude && !sign.New(sign.CustomSignerType).SetKey(key).SetKeyName(config.Config.GetString("sign.custom.key_name")).Validate(data) {
//				log.Logger.WithFields(logrus.Fields{
//					"sign": sign.New(sign.CustomSignerType).SetKey(key).SetKeyName(config.Config.GetString("sign.custom.key_name")).Sign(data),
//					"data": data,
//				}).Error("sign error")
//				response.New(c).Fail(exception.NewParamsError(errors.New("invalid key")))
//				c.Abort()
//				return
//			}
//		}
//		//}
//
//		c.Set("device_id", deviceID)
//		c.Set("platform", platform)
//		c.Set("device_model", deviceModel)
//		v, _ := strconv.Atoi(versionStr)
//		c.Set("version", v)
//		c.Set("bundle_id", bundleStr)
//		c.Set("channel", channelStr)
//		c.Set("language", language)
//		//device.NewVersionCache(models.Platform(platform), deviceID).Set(versionStr)
//
//		c.Next()
//	}
//}
