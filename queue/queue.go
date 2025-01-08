package queue

import (
	"context"
	"github.com/sosumecho/modules/drivers/nsq"

	"sync"
)

var (
	nsqQueue *NsqQueue
	once     = &sync.Once{}
)

// NsqQueue nsq队列
type NsqQueue struct {
	ctx      context.Context
	mutex    *sync.Mutex
	handlers map[string]nsq.MessageHandler
}

// Register 注册
func (n *NsqQueue) Register(handler nsq.MessageHandler) *NsqQueue {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.handlers[handler.Topic()] = handler
	return n
}

// Run 运行
func Run(conf *nsq.NsqConfig) {
	for _, item := range nsqQueue.handlers {
		go nsq.Consume(conf, nsqQueue.ctx, item)
	}
}

// New 新建
func New(ctx context.Context) *NsqQueue {
	if nsqQueue == nil {
		once.Do(func() {
			nsqQueue = &NsqQueue{
				ctx:      ctx,
				mutex:    &sync.Mutex{},
				handlers: make(map[string]nsq.MessageHandler),
			}
		})
	}
	return nsqQueue
}
