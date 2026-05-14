package ipc

import "net"

type Server struct {
	listener net.Listener
	handlers map[MsgType]HandlerFunc
	mu       sync.RWMutex
}

type HandlerFunc func(msg Message) ([]byte, error)

func NewServer(addr string) (*Server, error) {
	ln, err := Listen(addr)
	if err != nil {
		return nil, err
	}
	return &Server{
		listener: ln,
		handlers: make(map[MsgType]HandlerFunc),
	}, nil
}

func (s *Server) Handle(msg MsgType, handler HandlerFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[msg] = handler
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
	defer conn.Close()
	for {
		msg, err := ReadMsg(conn)
		if err != nil {
			return // 连接断开
		}

		go s.handleMsg(conn, msg)
	}
}

func (s *Server) handleMsg(conn net.Conn, msg Message) {
	var (
		err     error
		payload []byte
	)
	resp := Message{
		Type: MsgTypeResponse,
		ID:   msg.ID,
	}
	handler, ok := s.handlers[msg.Type]
	if !ok {
		resp.Payload, err = MarshalPayload(map[string]string{"error": "unknown msg type"})
		if err != nil {
			return
		}
		WriteMsg(conn, resp)
		return
	}
	payload, err = handler(msg)
	if err != nil {
		resp.Payload, _ = MarshalPayload(map[string]string{
			"error": err.Error(),
		})
	} else {
		resp.Payload = payload
	}
	WriteMsg(conn, resp)
}
