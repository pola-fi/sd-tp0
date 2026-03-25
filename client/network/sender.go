package network

import (
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/models"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/network/connection"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/network/protocol"
)

type MessageSender interface {
	SendBatch(bets []*models.Bet, serverAddress string) (*protocol.BatchResponse, error)
	SendDone(agencyID, serverAddress string) error
	QueryWinners(agencyID, serverAddress string) (*protocol.QueryWinnersResponse, error)
}

type BetMessageSender struct {
	connFactory connection.Connection
	proto       protocol.BetProtocol
}

func NewBetMessageSender() *BetMessageSender {
	return &BetMessageSender{
		connFactory: connection.NewTCPConnection(),
		proto:       protocol.NewBetProtocol(),
	}
}

func (s *BetMessageSender) SendBatch(bets []*models.Bet, serverAddress string) (*protocol.BatchResponse, error) {
	data, err := s.proto.SerializeBatch(bets)
	if err != nil {
		return nil, ErrSerializationFailed
	}

	respBytes, err := s.sendAndReceiveFixedResponse(serverAddress, data, protocol.BatchResponseSize)
	if err != nil {
		return nil, err
	}

	response, err := s.proto.DeserializeBatchResponse(respBytes)
	if err != nil {
		return nil, ErrReadResponseFailed
	}
	if !response.Success {
		return response, ErrServerRejectedBatch
	}

	return response, nil
}

func (s *BetMessageSender) SendDone(agencyID, serverAddress string) error {
	data, err := s.proto.SerializeDone(agencyID)
	if err != nil {
		return ErrControlRequestFailed
	}

	respBytes, err := s.sendAndReceiveFixedResponse(serverAddress, data, protocol.BatchResponseSize)
	if err != nil {
		return err
	}

	response, err := s.proto.DeserializeBatchResponse(respBytes)
	if err != nil {
		return ErrReadResponseFailed
	}
	if !response.Success {
		return ErrControlRequestFailed
	}
	return nil
}

func (s *BetMessageSender) QueryWinners(agencyID, serverAddress string) (*protocol.QueryWinnersResponse, error) {
	data, err := s.proto.SerializeQueryWinners(agencyID)
	if err != nil {
		return nil, ErrControlRequestFailed
	}

	respBytes, err := s.sendAndReceiveQueryResponse(serverAddress, data)
	if err != nil {
		return nil, err
	}

	response, err := s.proto.DeserializeQueryWinnersResponse(respBytes)
	if err != nil {
		return nil, ErrReadResponseFailed
	}
	return response, nil
}
