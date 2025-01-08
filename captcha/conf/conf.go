package conf

import "github.com/sosumecho/modules/sms"

type CustomEmailConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Product  bool   `mapstructure:"product"`
}

type Config struct {
	Sms         *sms.SmsConfig     `mapstructure:"sms"`
	Email       *EmailConfig       `mapstructure:"email"`
	CustomEmail *CustomEmailConfig `mapstructure:"custom_email"`
}

type EmailConfig struct {
	AppKey        string `mapstructure:"app-key"`
	AppSecret     string `mapstructure:"app-secret"`
	Host          string `mapstructure:"host"`
	Username      string `mapstructure:"username"`
	UsernameAlias string `mapstructure:"username-alias"`
	Product       bool   `mapstructure:"product"`
}
