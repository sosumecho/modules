package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/sosumecho/modules/i18n"
	"golang.org/x/text/language"
)

// I18N 设置语言中间件
func I18N(n *i18n.I18N) gin.HandlerFunc {
	return func(c *gin.Context) {
		n.SetLang(GetAcceptLang(c, n))
		c.Next()
	}
}

func GetAcceptLang(c *gin.Context, n *i18n.I18N) language.Tag {
	lang := c.Query("lang")
	if lang != "" {
		tag, err := language.Parse(lang)
		if err == nil {
			return tag
		}
	}
	return n.GetRootLanguage(c.GetHeader("Accept-language"), language.English)
}
