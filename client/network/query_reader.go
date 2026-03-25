package network

import (
	"io"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/network/protocol"
)

func (s *BetMessageSender) sendAndReceiveQueryResponse(serverAddress string, data []byte) ([]byte, error) {
	return s.executeRequest(serverAddress, data, s.readQueryResponse)
}

func (s *BetMessageSender) readQueryResponse(reader io.Reader) ([]byte, error) {
	header := make([]byte, protocol.QueryHeaderSize)
	if _, err := io.ReadFull(reader, header); err != nil {
		return nil, ErrReadResponseFailed
	}

	dniCSVLen := protocol.DecodeUint16BE(header[4], header[5])
	buf := make([]byte, protocol.QueryHeaderSize+dniCSVLen)
	copy(buf, header)
	if dniCSVLen > 0 {
		if _, err := io.ReadFull(reader, buf[protocol.QueryHeaderSize:]); err != nil {
			return nil, ErrReadResponseFailed
		}
	}

	return buf, nil
}

