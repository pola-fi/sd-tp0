package network

import (
	"errors"
	"io"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/models"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/network/connection"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/network/protocol"
)

// Errors
var (
	ErrConnectionFailed = errors.New("failed to connect to server")
	ErrSerializationFailed = errors.New("failed to serialize batch")
	ErrSendFailed = errors.New("failed to send data")
	ErrReadResponseFailed = errors.New("failed to read response")
	ErrServerRejectedBatch = errors.New("server rejected batch")
)

type MessageSender interface {
	SendBatch(bets []*models.Bet, serverAddress string) (*protocol.BatchResponse, error)
}

type BetMessageSender struct {
	connFactory connection.Connection
	proto     protocol.BetProtocol
}

func NewBetMessageSender() *BetMessageSender {
	return &BetMessageSender{
		connFactory: connection.NewTCPConnection(),
		proto:     protocol.NewMixedSchemaProtocol(),
	}
}

func (s *BetMessageSender) SendBatch(bets []*models.Bet, serverAddress string) (*protocol.BatchResponse, error) {
	err := s.connFactory.Connect(serverAddress)
	if err != nil {
		return nil, ErrConnectionFailed
	}
	defer s.connFactory.Close()

	data, err := s.proto.SerializeBatch(bets)
	if err != nil {
		return nil, ErrSerializationFailed
	}

	totalSent := 0
	for totalSent < len(data) {
		sent, err := s.connFactory.Write(data[totalSent:])
		if err != nil {
			return nil, ErrSendFailed
		}
		totalSent += sent
	}

	response, err := s.readResponse()
	if err != nil {
		return nil, ErrReadResponseFailed
	}
	if !response.Success {
		return response, ErrServerRejectedBatch
	}

	return response, nil
}

func (s *BetMessageSender) readResponse() (*protocol.BatchResponse, error) {
	buf := make([]byte, protocol.BatchResponseSize)
	if _, err := io.ReadFull(s.connFactory, buf); err != nil {
		return nil, err
	}
	return s.proto.DeserializeBatchResponse(buf)
}

