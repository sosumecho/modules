package ipc

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type Client struct {
	conn    Transport
	mu      sync.Mutex
	pending map[uint32]chan Message
	idGen   atomic.Uint32
	opt     *ClientOptions
}

type ClientOptions struct {
	timeout time.Duration
}

type Option func(*ClientOptions)

func WithTimeout(timeout time.Duration) Option {
	return func(options *ClientOptions) {
		options.timeout = timeout
	}
}

func (c *Client) readLoop() {
	defer func() {
		c.conn.Close()

		c.mu.Lock()
		for id, ch := range c.pending {
			close(ch)
			delete(c.pending, id)
		}
		c.mu.Unlock()
	}()
	for {
		msg, err := ReadMsg(c.conn)
		if err != nil {
			return
		}
		c.mu.Lock()
		ch, ok := c.pending[msg.ID]
		if ok {
			delete(c.pending, msg.ID)
		}
		c.mu.Unlock()
		if ok {
			ch <- msg
		}

	}
}

func (c *Client) Send(msgType MsgType, payload []byte) (Message, error) {
	id := c.idGen.Add(1)
	ch := make(chan Message, 1)

	c.mu.Lock()
	c.pending[id] = ch
	c.mu.Unlock()

	err := WriteMsg(c.conn, Message{
		Type:    msgType,
		ID:      id,
		Payload: payload,
	})
	if err != nil {
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
		return Message{}, err
	}
	timer := time.NewTimer(c.opt.timeout)
	defer timer.Stop()
	select {
	case resp := <-ch:
		return resp, nil
	case <-timer.C:
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
		return Message{}, errors.New("timeout")
	}
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func NewClient(addr string, options ...Option) (*Client, error) {
	conn, err := Dial(addr)
	if err != nil {
		return nil, err
	}

	c := &Client{
		conn:    conn,
		pending: make(map[uint32]chan Message),
		opt: &ClientOptions{
			timeout: 5 * time.Second,
		},
	}

	for _, opt := range options {
		opt(c.opt)
	}
	go c.readLoop()
	return c, nil
}
