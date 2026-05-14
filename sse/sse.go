package sse

import (
	"sync"
)

type MessageType string

type Message struct {
	Type MessageType `json:"type"`
	Data string      `json:"data"`
}

var client *Broker

type Broker struct {
	clients map[string]map[chan Message]bool // 用户ID -> 连接映射
	mu      sync.RWMutex
}

func NewBroker() *Broker {
	if client == nil {
		client = &Broker{
			clients: make(map[string]map[chan Message]bool),
		}
	}
	return client
}

// 添加客户端连接
func (b *Broker) AddClient(userID string, ch chan Message) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, exists := b.clients[userID]; !exists {
		b.clients[userID] = make(map[chan Message]bool)
	}
	b.clients[userID][ch] = true
}

// 移除客户端连接
func (b *Broker) RemoveClient(userID string, ch chan Message) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if userClients, exists := b.clients[userID]; exists {
		delete(userClients, ch)
		if len(userClients) == 0 {
			delete(b.clients, userID)
		}
	}
}

// 向指定用户发送消息
func (b *Broker) SendToUser(userID string, msg Message) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if userClients, exists := b.clients[userID]; exists {
		for ch := range userClients {
			select {
			case ch <- msg:
			default:
				// 防止阻塞，如果通道已满则跳过
			}
		}
	}
}

func (b *Broker) SendMessage(userID string, action, msg Message) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if userClients, exists := b.clients[userID]; exists {
		for ch := range userClients {
			select {
			case ch <- msg:
			default:
				// 防止阻塞，如果通道已满则跳过
			}
		}
	}
}

// 广播给所有用户
func (b *Broker) Broadcast(msg Message) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, userClients := range b.clients {
		for ch := range userClients {
			select {
			case ch <- msg:
			default:
				// 防止阻塞
			}
		}
	}
}
