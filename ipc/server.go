package ipc

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
)

type Server struct {
	listener    net.Listener
	handlers    map[MsgType]HandlerFunc
	rpcHandlers map[string]RPCHandlerFunc
	mu          sync.RWMutex
	rpcMu       sync.RWMutex
}

type HandlerFunc func(msg Message) ([]byte, error)

func NewServer(addr string) (*Server, error) {
	ln, err := Listen(addr)
	if err != nil {
		return nil, err
	}
	return &Server{
		listener:    ln,
		handlers:    make(map[MsgType]HandlerFunc),
		rpcHandlers: make(map[string]RPCHandlerFunc),
	}, nil
}

func (s *Server) Handle(msg MsgType, handler HandlerFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[msg] = handler
}

func (s *Server) HandleRPC(method string, handler RPCHandlerFunc) {
	s.rpcMu.Lock()
	defer s.rpcMu.Unlock()
	s.rpcHandlers[method] = handler
}

type Conn struct {
	net.Conn
	send chan Message

	closed chan struct{}
	once   sync.Once
}

func (c *Conn) writeLoop() {
	defer c.Close()
	for {
		select {
		case msg := <-c.send:
			if err := WriteMsg(c.Conn, msg); err != nil {
				return
			}
		case <-c.closed:
			return
		}
	}
}

func (c *Conn) Send(msg Message) bool {
	select {
	case c.send <- msg:
		return true
	case <-c.closed:
		return false
	}
}

func (c *Conn) Close() {
	c.once.Do(func() {
		close(c.closed)
		c.Conn.Close()
	})
}

func (s *Server) Serve() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	c := &Conn{
		Conn:   conn,
		send:   make(chan Message, 100),
		closed: make(chan struct{}),
	}
	defer c.Close()
	go c.writeLoop()
	for {
		msg, err := ReadMsg(c.Conn)
		if err != nil {
			return // 连接断开
		}

		go s.handleMsg(c, msg)
	}
}

func (s *Server) handleMsg(conn *Conn, msg Message) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("failed to handle Msg", err)
		}
	}()

	var (
		err     error
		payload []byte
	)

	switch msg.Type {
	case MsgTypePing:
		conn.Send(Message{
			Type: MsgTypePong,
			ID:   msg.ID,
		})
		return
	case MsgTypeRequest:
		resp := Message{
			Type: MsgTypeResponse,
			ID:   msg.ID,
		}
		var req RPCRequest
		if err = json.Unmarshal(msg.Payload, &req); err != nil {
			SendError(conn, msg, err)
			return
		}
		s.rpcMu.RLock()
		handler, ok := s.rpcHandlers[req.Method]
		s.rpcMu.RUnlock()
		if !ok {
			SendError(conn, msg, fmt.Errorf("method not found"))
			return
		}
		rpcResponse, handlerErr := handler(req.Params)

		rs := RPCResponse{}
		if handlerErr != nil {
			rs.Error = handlerErr.Error()
		} else {
			rpcRawResponse, err := json.Marshal(rpcResponse)
			if err != nil {
				SendError(conn, msg, err)
				return
			}
			rs.Result = rpcRawResponse
		}
		resp.Payload, err = json.Marshal(rs)
		if err != nil {
			SendError(conn, msg, err)
			return
		}
		conn.Send(resp)
		return
	default:
		resp := Message{
			Type: MsgTypeResponse,
			ID:   msg.ID,
		}
		s.mu.RLock()
		handler, ok := s.handlers[msg.Type]
		s.mu.RUnlock()
		if !ok && msg.Type != MsgTypeEvent {
			SendError(conn, msg, fmt.Errorf("method not found"))
			return
		}
		payload, err = handler(msg)
		if err != nil && msg.Type != MsgTypeEvent {
			SendError(conn, msg, err)
			return
		}
		if msg.Type != MsgTypeEvent {
			resp.Payload, _ = json.Marshal(RPCResponse{
				Result: payload,
			})

			conn.Send(resp)
		}
	}
}

func SendError(conn *Conn, msg Message, err error) {
	resp := Message{
		Type: MsgTypeResponse,
		ID:   msg.ID,
	}
	resp.Payload, err = MarshalPayload(RPCResponse{Error: err.Error()})
	if err != nil {
		return
	}
	conn.Send(resp)
}
