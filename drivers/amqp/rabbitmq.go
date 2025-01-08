package amqp

import (
	"fmt"
	"github.com/sosumecho/modules/logger"
	"go.uber.org/zap"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ 消息队列
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	done    chan bool
	logger  *logger.Logger
}

type RabbitMQConf struct {
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
	Host     string `json:"host" mapstructure:"host"`
	Port     int    `json:"port" mapstructure:"port"`
}

// SendMessage 发送消息
func (q *RabbitMQ) SendMessage(topic, message string) error {
	queue, err := q.channel.QueueDeclare(topic, true, false, false, false, nil)
	if err != nil {
		return err
	}
	return q.channel.Publish("", queue.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         []byte(message),
	})
}

// ConsumeMessage 消费消息
func (q *RabbitMQ) ConsumeMessage(topic string, callback func(data []byte) error) error {
	queue, err := q.channel.QueueDeclare(topic, true, false, false, false, nil)
	if err != nil {
		return err
	}
	messages, err := q.channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	for {
		select {
		case msg := <-messages:
			q.logger.Debug("msg", zap.Any("msg", msg.Body))
			if len(msg.Body) > 0 {
				err := callback(msg.Body)
				if err == nil {
					_ = msg.Ack(true)
				} else {
				}
			} else {
				time.Sleep(1 * time.Second)
				_ = msg.Ack(true)
			}
		case <-q.done:
			q.logger.Debug("done")
			return nil
		}
	}
}

func (q *RabbitMQ) Pop(topic string) ([]byte, error) {
	_ = q.channel.Qos(1, 0, false)
	queue, err := q.channel.QueueDeclare(topic, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	messages, err := q.channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	select {
	case msg := <-messages:
		q.logger.Debug("msg", zap.Any("msg", msg.Body))
		if len(msg.Body) > 0 {
			_ = msg.Ack(false)
			return msg.Body, nil
		}
	case <-time.After(time.Second * 12):
		q.logger.Debug("done")
	}
	return nil, err
}

// Close 关闭
func (q *RabbitMQ) Close() {
	err := q.channel.Close()
	if err != nil {
	}
	err = q.conn.Close()
	if err != nil {
	}
}

// New 初始化一个rabbitmq客户端
func New(conf *RabbitMQConf, logf *logger.Logger) *RabbitMQ {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d", conf.Username, conf.Password, conf.Host, conf.Port))
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
		done:    make(chan bool),
		logger:  logf,
	}
}
