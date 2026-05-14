//go:build linux || darwin

package ipc

import (
	"net"
	"os"
)

const defaultSockPath = "/tmp/ipc.sock"

func Listen(addr string) (net.Listener, error) {
	if err := os.Remove(addr); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return net.Listen("unix", addr)
}

func Dial(addr string) (net.Conn, error) {
	return net.Dial("unix", addr)
}
