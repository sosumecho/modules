package formatter

import (
	"fmt"
	"github.com/sosumecho/modules/exception"
	"github.com/sosumecho/modules/file/upload"
	"github.com/sosumecho/modules/i18n"
	"github.com/sosumecho/modules/logger"
	"github.com/sosumecho/modules/response"

	"github.com/gin-gonic/gin"
)

type CustomFormatter struct {
	logger *logger.Logger
}

func NewCustomFormatter(logger *logger.Logger) *CustomFormatter {
	return &CustomFormatter{logger: logger}
}

type CustomResponse struct {
	UID    int    `json:"uid"`
	Name   string `json:"name"`
	URL    string `json:"url"`
	Status string `json:"status"`
}

func (c2 CustomResponse) GetURL() string {
	return c2.URL
}

func (c2 CustomFormatter) RespFormat(c *gin.Context, locale *i18n.I18N, group string, domain string, fileInfo *upload.FileInfo) {
	rs := c2.Format(group, domain, fileInfo)
	response.New(c, locale, c2.logger).Data(rs)
}

func (c2 CustomFormatter) Format(group string, domain string, fileInfo *upload.FileInfo) upload.URLFormatter {
	return CustomResponse{
		UID:    0,
		Name:   fileInfo.OriginName,
		URL:    fmt.Sprintf("%s%s", domain, fileInfo.FilePath),
		Status: "done",
	}
}

func (c2 CustomFormatter) DryFormat(group string, domain string, fileInfo *upload.FileInfo) string {
	return fmt.Sprintf("%s%s", domain, fileInfo.FilePath)
}

func (c2 CustomFormatter) RespFormatMulti(c *gin.Context, locale *i18n.I18N, group string, domain string, filenames map[string]*upload.FileInfo) {
	resp := response.New(c, locale, c2.logger)
	rs, err := c2.FormatMulti(group, domain, filenames)
	if err != nil {
		resp.Fail(exception.NewSystemError(err))
		return
	}
	resp.Data(rs)
}

func (c2 CustomFormatter) FormatMulti(group string, domain string, filenames map[string]*upload.FileInfo) ([]upload.URLFormatter, error) {
	var rs = make([]upload.URLFormatter, 0, len(filenames))
	count := 0
	for k, v := range filenames {
		count++
		rs = append(rs, &CustomResponse{
			UID:    count,
			Name:   k,
			URL:    fmt.Sprintf("%s%s", domain, v.FilePath),
			Status: "done",
		})
	}
	return rs, nil
}
