package protocol

import "github.com/7574-sistemas-distribuidos/docker-compose-init/client/models"

const (
	BatchCountFieldSize = Uint16ByteSize
	BetLengthFieldSize  = Uint16ByteSize
	BatchResponseSize   = 3
	QueryHeaderSize     = 6 // [type:1][status:1][count:2][dni_csv_len:2]
	BatchStatusFail     = byte(0)
	BatchStatusSuccess  = byte(1)
	QueryStatusPending  = byte(0)
	QueryStatusReady    = byte(1)
	MaxPacketBytes      = 8192
	MaxUint16Value      = 65535

	MessageTypeBatch                = byte(1)
	MessageTypeDone                 = byte(2)
	MessageTypeQueryWinners         = byte(3)
	MessageTypeQueryWinnersResponse = byte(4)
)

type BatchResponse struct {
	Success bool
	Count   int
}

type QueryWinnersResponse struct {
	Ready      bool
	Count      int
	WinnerDNIs []string
}

type BetProtocol interface {
	SerializeBatch(bets []*models.Bet) ([]byte, error)
	SerializeDone(agencyID string) ([]byte, error)
	SerializeQueryWinners(agencyID string) ([]byte, error)
	DeserializeBatchResponse(data []byte) (*BatchResponse, error)
	DeserializeQueryWinnersResponse(data []byte) (*QueryWinnersResponse, error)
}

