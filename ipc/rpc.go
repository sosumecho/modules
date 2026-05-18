package ipc

import "encoding/json"

type RPCRequest struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type RPCResponse struct {
	Error  string          `json:"error"`
	Result json.RawMessage `json:"result"`
}

type RPCHandlerFunc func(conn *Conn, params json.RawMessage) (any, error)

type RPCEvent struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}
