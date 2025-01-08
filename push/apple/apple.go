package apple

import (
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
	"github.com/sosumecho/modules/push"
	"sync"
)

const (
	Name = "apns"
)

var (
	applePush *APNs
	once      sync.Once
)

type APNsConf struct {
	P12File     string `mapstructure:"p12_file"`
	P12Password string `mapstructure:"p12_password"`
	Production  bool   `mapstructrue:"production"`
	Topic       string `mapstructure:"topic"`
}

type APNs struct {
	conf   *APNsConf
	client *apns2.Client
}

func (a *APNs) Push(message push.Message) error {
	if len(message.Token) > 0 {
		for _, token := range message.Token {
			notification := &apns2.Notification{}
			notification.DeviceToken = token
			notification.Topic = a.conf.Topic
			payloads := payload.NewPayload().Alert(map[string]string{
				"title": message.Title,
				"body":  message.Body,
			}).Badge(1).SoundName("default")
			for k, v := range message.Data {
				payloads = payloads.Custom(k, v)
			}
			notification.Payload = payloads
			rs, err := a.client.Push(notification)
			if err != nil {
				continue
			}
			if rs.StatusCode != 200 {
				continue
			}
		}

	}
	return nil
}

func New(conf *APNsConf) push.Push {
	if applePush == nil {
		once.Do(func() {
			cert, err := certificate.FromP12File(conf.P12File, conf.P12Password)
			if err != nil {
				panic(err)
			}
			client := apns2.NewClient(cert)
			if conf.Production {
				client = client.Production()
			} else {
				client = client.Development()
			}
			applePush = &APNs{
				conf:   conf,
				client: client,
			}
		})
	}
	return applePush
}
