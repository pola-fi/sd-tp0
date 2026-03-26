package network

import "io"

func (s *BetMessageSender) sendAndReceiveFixedResponse(serverAddress string, data []byte, responseSize int) ([]byte, error) {
	return s.executeRequest(serverAddress, data, func(reader io.Reader) ([]byte, error) {
		buf := make([]byte, responseSize)
		if _, err := io.ReadFull(reader, buf); err != nil {
			return nil, ErrReadResponseFailed
		}
		return buf, nil
	})
}

func (s *BetMessageSender) executeRequest(serverAddress string, data []byte, readResponse func(reader io.Reader) ([]byte, error)) ([]byte, error) {
	if err := s.connFactory.Connect(serverAddress); err != nil {
		return nil, ErrConnectionFailed
	}
	defer func() {
		_ = s.connFactory.Close()
		log.Infof("action: close_socket | result: success | resource: client_socket | client_id: %s", s.clientID)
	}()

	if err := writeAll(s.connFactory, data); err != nil {
		return nil, err
	}

	return readResponse(s.connFactory)
}

func writeAll(writer io.Writer, data []byte) error {
	totalSent := 0
	for totalSent < len(data) {
		sent, err := writer.Write(data[totalSent:])
		if err != nil {
			return ErrSendFailed
		}
		totalSent += sent
	}
	return nil
}

