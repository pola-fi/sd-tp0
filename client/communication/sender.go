package communication

import (
	"bufio"
	"errors"
	"io"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/connection"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/models"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/protocol"
)

// Errors
var (
	ErrConnectionFailed = errors.New("failed to connect to server")
	ErrSerializationFailed = errors.New("failed to serialize bet")
	ErrSendFailed = errors.New("failed to send data")
	ErrReadResponseFailed = errors.New("failed to read response")
)

type MessageSender interface {
	SendBet(bet *models.Bet, serverAddress string) (string, error)
}

type BetMessageSender struct {
	connFactory connection.Connection
	protocol     protocol.BetProtocol
}

func NewBetMessageSender() *BetMessageSender {
	return &BetMessageSender{
		connFactory: connection.NewTCPConnection(),
		protocol:     protocol.NewMixedSchemaProtocol(),
	}
}

func (s *BetMessageSender) SendBet(bet *models.Bet, serverAddress string) (string, error) {
	err := s.connFactory.Connect(serverAddress)
	if err != nil {
		return "", ErrConnectionFailed
	}
	defer s.connFactory.Close()

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
