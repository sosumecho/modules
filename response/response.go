package response

import (
	"github.com/sosumecho/modules/exception"
	"github.com/sosumecho/modules/i18n"
	"github.com/sosumecho/modules/logger"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"net/http"

	"github.com/gin-gonic/gin"
)

const Ok = "ok"

type Result struct {
	Code exception.Code `json:"code"`
	Msg  string         `json:"msg"`
	Data interface{}    `json:"data"`
}

type Response struct {
	c      *gin.Context
	params map[string]interface{}
	locale *i18n.I18N
	logger *logger.Logger
}

func (r *Response) WithParams(params map[string]interface{}) *Response {
	r.params = params
	return r
}

func (r *Response) Data(data interface{}) {
	r.Json(http.StatusOK, Result{
		Code: 0,
		Msg:  r.locale.Tr(GetAcceptLang(r.c, r.locale), "success", nil),
		Data: data,
	})
}

func (r *Response) Json(status int, data interface{}) {
	rs, err := i18n.NewJSON(r.logger, GetAcceptLang(r.c, r.locale)...).MarshalToString(data)
	if err != nil {
		r.logger.Error("json", zap.Error(err))
	}
	r.c.Header("Content-Type", "application/json")
	r.c.Status(status)
	_, err = r.c.Writer.Write([]byte(rs))
	if err != nil {
		r.logger.Error("json", zap.Error(err))
	}
}

func GetAcceptLang(c *gin.Context, n *i18n.I18N) []language.Tag {
	lang := c.Query("lang")
	if lang != "" {
		tag, err := language.Parse(lang)
		if err == nil {
			return []language.Tag{tag}
		}
	}
	return n.GetRootLanguages(c.GetHeader("Accept-language"), language.English)
}

func (r *Response) Fail(msg exception.Exception) {
	r.Json(msg.HttpStatus(), Result{
		Code: msg.Code(),
		Msg:  r.locale.Tr(GetAcceptLang(r.c, r.locale), msg.Error(), r.params),
	})
}

func (r *Response) Error(statusCode int, msg exception.Exception) {
	r.Json(statusCode, Result{
		Code: msg.Code(),
		Msg:  r.locale.Tr(GetAcceptLang(r.c, r.locale), msg.Error(), r.params),
	})
}

func New(c *gin.Context, locale *i18n.I18N, log *logger.Logger) *Response {
	return &Response{
		c:      c,
		locale: locale,
		logger: log,
	}
}
