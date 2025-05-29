package ws

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sosumecho/modules/logger"
	"go.uber.org/zap"
)

type Message struct {
	Type    string `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Response struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Data    any `json:"data"`
}

type MessageHandler func(clientID string, payload json.RawMessage) (*Response, bool) 

type Server struct {
	upgrader websocket.Upgrader
	clients map[string]*websocket.Conn
	handlers map[string]MessageHandler
	clientMutex sync.Mutex
	handlerMutex sync.Mutex
	logger *logger.Logger
}

func NewServer(l *logger.Logger) *Server { 
	return &Server{
		upgrader: websocket.Upgrader{
			ReadBufferSize: 1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		clients: make(map[string]*websocket.Conn),
		handlers: make(map[string]MessageHandler),
		logger: l,
	}
}

// RegisterHandler 注册消息处理函数
func (s *Server) RegisterHandler(messageType string, handler MessageHandler) { 
	s.handlerMutex.Lock()
	defer s.handlerMutex.Unlock()
	s.handlers[messageType] = handler
}

func (s *Server) HandleConnection(w http.ResponseWriter, r *http.Request) { 
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error("升级WebSocket连接失败", zap.Error(err))
		return
	}
	clientID := uuid.New().String()
	s.clientMutex.Lock()
	s.clients[clientID] = conn
	s.clientMutex.Unlock()

	s.logger.Info("客户端连接成功", zap.String("clientID", clientID))

	defer func ()  {
		s.clientMutex.Lock()
		delete(s.clients, clientID)
		s.clientMutex.Unlock()
		conn.Close()
		s.logger.Info("客户端断开连接", zap.String("clientID", clientID))
	}()

	s.readMessages(clientID, conn)
}

func (s *Server) readMessages(clientID string, conn *websocket.Conn) { 
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			s.logger.Error("读取消息失败", zap.Error(err))
			return
		}
		var msg Message
		if err := json.Unmarshal(message, &msg); err!= nil {
			s.logger.Error("解析消息失败", zap.Error(err))
			continue
		}

		s.handlerMutex.Lock()
		handler, ok := s.handlers[msg.Type]
		s.handlerMutex.Unlock()

		if !ok {
			s.logger.Warn("未注册的消息类型", zap.String("messageType", msg.Type))
			continue
		}

		go func ()  {
			response, broadcast := handler(clientID, msg.Payload)

			if response != nil {
				responseJSON, err := json.Marshal(response)
				if err != nil {
					s.logger.Error("序列化响应消息失败", zap.Error(err))
					return
				}
				if broadcast {
					s.Broadcast(responseJSON)
				} else {
					s.Send(clientID, responseJSON)
				}
			}
		}()
	}
}

func (s *Server) Send(clientID string, message []byte) error {
	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()

	if conn, ok := s.clients[clientID]; ok {
		return conn.WriteMessage(websocket.TextMessage, message)
	}

	return nil
}

func (s *Server) Broadcast(message []byte) {
	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()

	for _, conn := range s.clients {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			s.logger.Error("广播消息失败", zap.Error(err))
		}
	}
	
}