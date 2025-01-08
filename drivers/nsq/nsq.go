package nsq

import (
	"context"
	"github.com/nsqio/go-nsq"
	"sync"
	"time"
)

var (
	producer *nsq.Producer
	once     = &sync.Once{}
)

type NsqConfig struct {
	Lookups []string `mapstructure:"lookups"`
	NSQD    string   `mapstructure:"nsqd"`
}

// MessageHandler
type MessageHandler interface {
	nsq.Handler
	Topic() string
	Channel() string
	Workers() int
}

// Consume 消费
func Consume(conf *NsqConfig, ctx context.Context, handler MessageHandler) {
	consumerConf := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(handler.Topic(), handler.Channel(), consumerConf)
	if err != nil {
		return
	}
	consumer.AddConcurrentHandlers(handler, handler.Workers())

	if err = consumer.ConnectToNSQD(conf.NSQD); err != nil {
		return
	}
	select {
	case <-ctx.Done():
		return
	}
}

// Produce 生产
func Produce(conf *NsqConfig, topic string, message []byte) error {
	delayInit(conf)
	if err := producer.Publish(topic, message); err != nil {
		return err
	}
	return nil
}

// DelayProduce 延迟生产
func DelayProduce(conf *NsqConfig, topic string, message []byte, delay time.Duration) error {
	delayInit(conf)
	if err := producer.DeferredPublish(topic, delay, message); err != nil {
		return err
	}
	return nil
}

func delayInit(nsqClientConf *NsqConfig) {
	if producer == nil {
		once.Do(func() {
			var err error
			producerConf := nsq.NewConfig()
			producer, err = nsq.NewProducer(nsqClientConf.NSQD, producerConf)
			if err != nil {

			}
		})
	}
}
