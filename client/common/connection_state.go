package common

import (
	"net"
	"sync"
)

type connectionState struct {
	mu   sync.Mutex
	conn net.Conn
}

func (s *connectionState) set(conn net.Conn) {
	s.mu.Lock()
	s.conn = conn
	s.mu.Unlock()
}

func (s *connectionState) connect(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	s.set(conn)
	return nil
}

func (s *connectionState) get() net.Conn {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.conn
}

func (s *connectionState) close(clientID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.conn == nil {
		return
	}

	_ = s.conn.Close()
	s.conn = nil
	log.Infof("action: close_socket | result: success | resource: client_socket | client_id: %s", clientID)
}

