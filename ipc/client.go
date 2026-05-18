package ipc

import (
	"encoding/json"
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

	closeChan chan struct{}
	once      sync.Once

	eventHandlers map[string]EventHandlerFunc
	eventMu       sync.RWMutex
}

type EventHandlerFunc func(data json.RawMessage)

type ClientOptions struct {
	timeout      time.Duration
	pingInterval time.Duration
}

type Option func(*ClientOptions)

func WithTimeout(timeout time.Duration) Option {
	return func(options *ClientOptions) {
		options.timeout = timeout
	}
}

func WithPingInterval(pingInterval time.Duration) Option {
	return func(options *ClientOptions) {
		options.pingInterval = pingInterval
	}
}

func (c *Client) readLoop() {
	defer func() {
		c.Close()

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
		} else if msg.Type == MsgTypeEvent {
			var evt RPCEvent
			if err := json.Unmarshal(msg.Payload, &evt); err == nil {
				c.eventMu.RLock()
				handler, ok := c.eventHandlers[evt.Event]
				c.eventMu.RUnlock()
				if ok {
					go handler(evt.Data)
				}
			}
		}

	}
}

func (c *Client) Call(msgType MsgType, payload []byte) (Message, error) {
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

func (c *Client) Send(method string, params any, resp any) error {
	b, err := json.Marshal(params)
	if err != nil {
		return err
	}
	req := RPCRequest{
		Method: method,
		Params: b,
	}
	b, err = json.Marshal(req)
	if err != nil {
		return err
	}
	rs, err := c.Call(MsgTypeRequest, b)
	if err != nil {
		return err
	}
	var rpcResp RPCResponse
	err = json.Unmarshal(rs.Payload, &rpcResp)
	if err != nil {
		return err
	}
	if rpcResp.Error != "" {
		return errors.New(rpcResp.Error)
	}
	return json.Unmarshal(rpcResp.Result, resp)

}

func (c *Client) Publish(payload []byte) error {
	return WriteMsg(c.conn, Message{
		Type:    MsgTypeEvent,
		ID:      0,
		Payload: payload,
	})
}

func (c *Client) HandleEvent(event string, handler EventHandlerFunc) {
	c.eventMu.Lock()
	defer c.eventMu.Unlock()
	c.eventHandlers[event] = handler
}

func (c *Client) Close() error {
	c.once.Do(func() {
		c.closeChan <- struct{}{}
		c.conn.Close()
	})
	return nil
}

func (c *Client) keepAlive() {
	if c.opt.pingInterval <= 0 {
		return
	}
	ticker := time.NewTicker(c.opt.pingInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			_, _ = c.Call(MsgTypePing, nil)
		case <-c.closeChan:
			return
		}
	}
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
			timeout:      5 * time.Second,
			pingInterval: time.Second * 30,
		},
		eventHandlers: make(map[string]EventHandlerFunc),
		closeChan: make(chan struct{}),
	}

	for _, opt := range options {
		opt(c.opt)
	}
	go c.readLoop()
	go c.keepAlive()
	return c, nil
}
