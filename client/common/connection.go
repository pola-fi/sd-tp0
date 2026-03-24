package common

import (
	"net"
)

type Connection interface {
	Connect(address string) error
	Write(data []byte) (int, error)
	Read(buffer []byte) (int, error)
	Close() error
}

type TCPConnection struct {
	conn net.Conn
}

func NewTCPConnection() *TCPConnection {
	return &TCPConnection{}
}

func (tcp *TCPConnection) Connect(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	tcp.conn = conn
	return nil
}

func (tcp *TCPConnection) Write(data []byte) (int, error) {
	if tcp.conn == nil {
		return 0, net.ErrClosed
	}
	return tcp.conn.Write(data)
}

func (tcp *TCPConnection) Read(buffer []byte) (int, error) {
	if tcp.conn == nil {
		return 0, net.ErrClosed
	}
	return tcp.conn.Read(buffer)
}

func (tcp *TCPConnection) Close() error {
	if tcp.conn == nil {
		return nil
	}
	err := tcp.conn.Close()
	tcp.conn = nil
	return err
}
