package dingtalk

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/sosumecho/modules/push"
)

const (
	WebHooKURL = "https://oapi.dingtalk.com/robot/send?access_token="
)

type DingTalk struct {
	mobiles []string
	secret  string
	keyword string
	msgType string
}

type Conf struct {
	Secret string `mapstructure:"secret"`
}

func NewDingTalk(conf *Conf) *DingTalk {
	return &DingTalk{
		secret:  conf.Secret,
		mobiles: make([]string, 0),
	}
}

type Message struct {
	MsgType  string       `json:"msgtype"`
	Text     *TextContent `json:"text,omitempty"`
	MarkDown *MarkDown    `json:"markdown,omitempty"`
	At       At           `json:"at"`
}

func NewMessage(msgType string) *Message {
	return &Message{MsgType: msgType}
}

func (m *Message) SetContent(content push.Message) *Message {
	switch m.MsgType {
	case "text":
		m.Text = &TextContent{Content: fmt.Sprintf("[%s] %s", content.Title, content.Body)}
	case "markdown":
		m.MarkDown = &MarkDown{
			Title: content.Title,
			Text:  content.Body,
		}
	}
	return m
}

func (m *Message) SetMobile(mobile []string) *Message {
	m.At.AtMobiles = mobile
	return m
}

type At struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type TextContent struct {
	Content string `json:"content"`
}

type MarkDown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func (d *DingTalk) SetMobiles(mobiles []string) *DingTalk {
	d.mobiles = mobiles
	return d
}

func (d *DingTalk) SetMessageType(msgType string) *DingTalk {
	d.msgType = msgType
	return d
}

func (d *DingTalk) SetKeyword(kw string) *DingTalk {
	d.keyword = kw
	return d
}

func (d *DingTalk) Push(message push.Message) error {
	if message.Title == "" {
		message.Title = d.keyword
	}
	m := NewMessage(d.msgType).SetContent(message).SetMobile(d.mobiles)
	b, err := jsoniter.Marshal(m)
	if err != nil {
		return err
	}
	_, err = resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(b).
		Post(fmt.Sprintf("%s%s", WebHooKURL, d.secret))
	if err != nil {
		return err
	}
	return nil
}
