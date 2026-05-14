//go:build windows

package ipc

import (
	"net"

	"github.com/Microsoft/go-winio"
)

const defaultSockPath = `\\.\pipe\ipc`

func Listen(addr string) (net.Listener, error) {
	return winio.ListenPipe(addr, nil)
}

func Dial(addr string) (net.Conn, error) {
	return winio.DialPipe(addr, nil)
}
