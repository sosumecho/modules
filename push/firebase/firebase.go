package firebase

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"github.com/sosumecho/modules/push"
	"github.com/sosumecho/modules/utils"
	"google.golang.org/api/option"
	"os"
	"strings"
	"sync"
)

const (
	Name = "firebase"
)

var (
	firebaseObj *Firebase
	once        sync.Once
)

// Firebase firebase
type Firebase struct {
	ctx    context.Context
	client *firebase.App
}

// Push 推送
func (f *Firebase) Push(message push.Message) error {
	client, err := f.client.Messaging(f.ctx)
	if err != nil {
		return err
	}
	// 如果有平台
	if message.Topic != "" {
		msg := messaging.Message{
			Data: message.Data,
			Notification: &messaging.Notification{
				Title: message.Title,
				Body:  message.Body,
			},
			APNS:  &messaging.APNSConfig{},
			Topic: message.Topic,
		}
		_, err := client.Send(f.ctx, &msg)
		if err != nil {
			return err
		}
	} else if len(message.Token) == 0 {
	} else {
		br, err := client.SendEachForMulticast(f.ctx, &messaging.MulticastMessage{
			Tokens: message.Token,
			Notification: &messaging.Notification{
				Title: message.Title,
				Body:  message.Body,
			},
			Data: message.Data,
		})
		if err != nil {
			return err
		}
		if br.FailureCount > 0 {
			for _, resp := range br.Responses {
				if !resp.Success {
					return resp.Error
				}
			}
		}
	}
	return nil
}

func (f *Firebase) PushMany(message push.Message) ([]string, error) {
	client, err := f.client.Messaging(f.ctx)
	if err != nil {
		return nil, err
	}

	br, err := client.SendEachForMulticast(f.ctx, &messaging.MulticastMessage{
		Tokens: message.Token,
		Notification: &messaging.Notification{
			Title: message.Title,
			Body:  message.Body,
		},
		Data: message.Data,
	})
	if err != nil {
		return nil, err
	}
	var failToken = make([]string, 0)
	if br.FailureCount > 0 {
		for index, resp := range br.Responses {
			if !resp.Success && resp.Error != nil && strings.Contains(resp.Error.Error(), "Requested entity was not found") {
				failToken = append(failToken, message.Token[index])
			}
		}
	}
	return failToken, err
}

func New() *Firebase {
	if firebaseObj == nil {
		once.Do(func() {
			var (
				err error
				app *firebase.App
			)
			ctx := context.Background()
			if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
				opt := option.WithCredentialsFile(fmt.Sprintf("%s/configs/%s", utils.GetAbsDir(), fmt.Sprintf("%s", "firebase.json")))
				app, err = firebase.NewApp(ctx, nil, opt)
			} else {
				app, err = firebase.NewApp(ctx, nil)
			}

			if err != nil {
				panic("创建firebase失败")
			}
			firebaseObj = &Firebase{
				ctx:    ctx,
				client: app,
			}
		})
	}
	return firebaseObj
}
