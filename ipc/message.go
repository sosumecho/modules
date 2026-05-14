package ipc

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

const (
	headerSize = 10 // type 2 + ID 4 + payloadLen 4
	maxMsgLen  = 64 << 20
)

type MsgType uint16

const (
	MsgTypeRequest MsgType = iota
	MsgTypeResponse
	MsgTypeEvent
	MsgTypePing
	MsgTypePong
)

type Message struct {
	Type    MsgType
	ID      uint32
	Payload []byte
}

func WriteMsg(w io.Writer, msg Message) error {
	data := make([]byte, headerSize+len(msg.Payload))
	binary.BigEndian.PutUint16(data[0:2], uint16(msg.Type))
	binary.BigEndian.PutUint32(data[2:6], msg.ID)
	binary.BigEndian.PutUint32(data[6:10], uint32(len(msg.Payload)))
	if len(msg.Payload) > 0 {
		copy(data[headerSize:], msg.Payload)
	}
	_, err := w.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func ReadMsg(r io.Reader) (Message, error) {
	header := make([]byte, headerSize)
	_, err := io.ReadFull(r, header)
	if err != nil {
		return Message{}, err
	}

	msgType := MsgType(binary.BigEndian.Uint16(header[:2]))
	msgID := binary.BigEndian.Uint32(header[2:6])
	payloadLen := binary.BigEndian.Uint32(header[6:10])
	if payloadLen > maxMsgLen {
		return Message{}, fmt.Errorf("payload too large")
	}
	payload := make([]byte, payloadLen)
	_, err = io.ReadFull(r, payload)
	if err != nil {
		return Message{}, err
	}
	return Message{
		Type:    msgType,
		ID:      msgID,
		Payload: payload,
	}, nil
}

func MarshalPayload(v any) ([]byte, error) {
	return json.Marshal(v)
}

func UnmarshalPayload(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
