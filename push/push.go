package push

import "sync"

var (
	pushMap = make(map[string]Push)
	mutex   = &sync.Mutex{}
)

// Push 推送
type Push interface {
	// Push 推送消息
	Push(message Message) error
}

// Message 消息
type Message struct {
	Token []string
	Data  map[string]string
	Title string
	Body  string
	Topic string
}

// Register 注册
func Register(name string, push Push) {
	mutex.Lock()
	defer mutex.Unlock()
	pushMap[name] = push
}

// New 新建
func New(name string) Push {
	if _, ok := pushMap[name]; ok {
		return pushMap[name]
	}
	return nil
}
