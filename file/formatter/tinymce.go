package formatter

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sosumecho/modules/exception"
	"github.com/sosumecho/modules/file/upload"
	"github.com/sosumecho/modules/i18n"
	"github.com/sosumecho/modules/logger"
	"github.com/sosumecho/modules/response"
	"net/http"
)

type TinymceFormatter struct {
	logger *logger.Logger
}

func (t TinymceFormatter) RespFormat(c *gin.Context, locale *i18n.I18N, group string, domain string, fileInfo *upload.FileInfo) {
	response.New(c, locale, t.logger).Json(http.StatusOK, t.Format(group, domain, fileInfo))
}

func (t TinymceFormatter) Format(group string, domain string, fileInfo *upload.FileInfo) upload.URLFormatter {
	return tinymceResponse{
		Name: fileInfo.OriginName,
		URL:  fmt.Sprintf("%s%s", domain, fileInfo.FilePath),
	}
}

func (t TinymceFormatter) DryFormat(group string, domain string, fileInfo *upload.FileInfo) string {
	return fmt.Sprintf("%s%s", domain, fileInfo.FilePath)
}

func (t TinymceFormatter) RespFormatMulti(c *gin.Context, locale *i18n.I18N, group string, domain string, filenames map[string]*upload.FileInfo) {
	resp := response.New(c, locale, t.logger)
	rs, err := t.FormatMulti(group, domain, filenames)
	if err != nil {
		resp.Fail(exception.NewSystemError(err))
		return
	}
	resp.Data(rs)
}

func (t TinymceFormatter) FormatMulti(group string, domain string, filenames map[string]*upload.FileInfo) ([]upload.URLFormatter, error) {
	var rs = make([]upload.URLFormatter, 0, len(filenames))
	count := 0
	for k, v := range filenames {
		count++
		rs = append(rs, tinymceResponse{
			Name: k,
			URL:  fmt.Sprintf("%s%s", domain, v.FilePath),
		})
	}
	return rs, nil
}

func NewTinymceFormatter(logger *logger.Logger) *TinymceFormatter {
	return &TinymceFormatter{logger: logger}
}

type tinymceResponse struct {
	Name string `json:"name"`
	URL  string `json:"location"`
}

func (t tinymceResponse) GetURL() string {
	return t.URL
}

//
//func (c2 TinymceFormatter) Format(
//	c *gin.Context,
//	locale *i18n.I18N,
//	group string,
//	domain string,
//	fileInfo *upload.FileInfo,
//) {
//	response.New(c, locale, c2.logger).Json(http.StatusOK, tinymceResponse{
//		Name: fileInfo.OriginName,
//		URL:  fmt.Sprintf("%s%s", domain, fileInfo.FilePath),
//	})
//}
//
//func (c2 TinymceFormatter) DryFormat(
//	c *gin.Context,
//	locale *i18n.I18N,
//	group string,
//	domain string,
//	fileInfo *upload.FileInfo,
//) string {
//	return fmt.Sprintf("%s%s", domain, fileInfo.FilePath)
//}
//
//func (c2 TinymceFormatter) FormatMulti(
//	c *gin.Context,
//	locale *i18n.I18N,
//	group string,
//	domain string,
//	filenames map[string]upload.FileInfo,
//) {
//	var rs = make([]tinymceResponse, 0, len(filenames))
//	count := 0
//	for k, v := range filenames {
//		count++
//		rs = append(rs, tinymceResponse{
//			Name: k,
//			URL:  fmt.Sprintf("%s%s", domain, v.FilePath),
//		})
//	}
//	response.New(c, locale, c2.logger).Json(http.StatusOK, rs)
//}
