package common

import (
	"bufio"
	"errors"
	"io"
)

// Errors
var (
	ErrConnectionFailed = errors.New("failed to connect to server")
	ErrSerializationFailed = errors.New("failed to serialize bet")
	ErrSendFailed = errors.New("failed to send data")
	ErrReadResponseFailed = errors.New("failed to read response")
)

type MessageSender interface {
	SendBet(bet *Bet, serverAddress string) (string, error)
}

type BetMessageSender struct {
	connFactory Connection
	protocol    BetProtocol
	clientID    string
}

func NewBetMessageSender(clientID string) *BetMessageSender {
	return &BetMessageSender{
		connFactory: NewTCPConnection(),
		protocol:    NewMixedSchemaProtocol(),
		clientID:    clientID,
	}
}

func (s *BetMessageSender) SendBet(bet *Bet, serverAddress string) (string, error) {
	err := s.connFactory.Connect(serverAddress)
	if err != nil {
		return "", ErrConnectionFailed
	}
	defer func() {
		_ = s.connFactory.Close()
		log.Infof("action: close_socket | result: success | resource: client_socket | client_id: %s", s.clientID)
	}()

	data, err := s.protocol.Serialize(bet)
	if err != nil {
		return "", ErrSerializationFailed
	}

	totalSent := 0
	for totalSent < len(data) {
		sent, err := s.connFactory.Write(data[totalSent:])
		if err != nil {
			return "", ErrSendFailed
		}
		totalSent += sent
	}

	response, err := s.readResponse()
	if err != nil {
		return "", ErrReadResponseFailed
	}

	return response, nil
}

func (s *BetMessageSender) readResponse() (string, error) {
	reader := bufio.NewReader(s.connFactory)
	response, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	return response, nil
}
