package protocol

import (
	"errors"
	"fmt"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/models"
)

const (
	BatchCountFieldSize = Uint16ByteSize
	BetLengthFieldSize  = Uint16ByteSize
	BatchResponseSize   = 3
	BatchStatusFail     = byte(0)
	BatchStatusSuccess  = byte(1)
	MaxPacketBytes      = 8192
	MaxUint16Value      = 65535
)

type BatchResponse struct {
	Success bool
	Count   int
}

type BetProtocol interface {
	SerializeBatch(bets []*models.Bet) ([]byte, error)
	DeserializeBatchResponse(data []byte) (*BatchResponse, error)
}

type MixedSchemaProtocol struct{}

func NewMixedSchemaProtocol() *MixedSchemaProtocol {
	return &MixedSchemaProtocol{}
}

// Format: [agency_id:1][dni:8][fecha:10][len_nombre:1][nombre][len_apellido:1][apellido][len_numero:2][numero]
func (p *MixedSchemaProtocol) Serialize(bet *models.Bet) ([]byte, error) {

	msg := make([]byte, 0)

	// Fixed fields
	msg = append(msg, bet.AgencyID[0])     // 1 byte agency_id
	msg = append(msg, bet.Documento...)    // 8 bytes dni
	msg = append(msg, bet.Nacimiento...)   // 10 bytes fecha

	// Variable fields with length prefix
	msg = append(msg, byte(len(bet.Nombre)))
	msg = append(msg, bet.Nombre...)

	msg = append(msg, byte(len(bet.Apellido)))
	msg = append(msg, bet.Apellido...)

	// 2 bytes for number length (big endian)
	numLen := len(bet.Numero)
	msg = append(msg, EncodeUint16BE(numLen)...)
	msg = append(msg, bet.Numero...)

	return msg, nil
}

func (p *MixedSchemaProtocol) SerializeBatch(bets []*models.Bet) ([]byte, error) {
	if len(bets) > MaxUint16Value {
		return nil, errors.New("batch too large")
	}

	msg := make([]byte, 0)
	msg = append(msg, EncodeUint16BE(len(bets))...)

	for _, bet := range bets {
		payload, err := p.Serialize(bet)
		if err != nil {
			return nil, err
		}
		if len(payload) > MaxUint16Value {
			return nil, errors.New("single bet payload too large")
		}

		msg = append(msg, EncodeUint16BE(len(payload))...)
		msg = append(msg, payload...)
	}

	return msg, nil
}

func EstimateSingleBetBatchSize(protocol BetProtocol, bet *models.Bet) (int, error) {
	serializedBatch, err := protocol.SerializeBatch([]*models.Bet{bet})
	if err != nil {
		return 0, err
	}
	if len(serializedBatch) < BatchCountFieldSize {
		return 0, errors.New("invalid serialized batch size")
	}
	return len(serializedBatch) - BatchCountFieldSize, nil
}

func (p *MixedSchemaProtocol) DeserializeBatchResponse(data []byte) (*BatchResponse, error) {
	if len(data) != BatchResponseSize {
		return nil, errors.New("invalid batch response size")
	}

	status := data[0]
	if status != BatchStatusSuccess && status != BatchStatusFail {
		return nil, fmt.Errorf("invalid batch status: %d", status)
	}

	count := DecodeUint16BE(data[1], data[2])
	return &BatchResponse{Success: status == BatchStatusSuccess, Count: count}, nil
}

