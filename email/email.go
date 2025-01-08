package email

import (
	"bytes"
	"github.com/sosumecho/modules/captcha/conf"
	"gopkg.in/gomail.v2"
)

func SendCustom(config *conf.CustomEmailConfig, account, subject, content string) error {
	d := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	m := gomail.NewMessage()
	m.SetHeader("From", config.Username)
	m.SetHeader("To", account)
	m.SetHeader("Subject", subject)
	var buf bytes.Buffer
	buf.WriteString(content)
	m.SetBody("text/html", buf.String())
	return d.DialAndSend(m)
}
