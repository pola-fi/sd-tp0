package protocol

import (
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/models"
)

// Format: [agency_id:1][dni:8][fecha:10][len_nombre:1][nombre][len_apellido:1][apellido][len_numero:2][numero]
func (p *mixedSchemaProtocol) Serialize(bet *models.Bet) ([]byte, error) {
	msg := make([]byte, 0)

	msg = append(msg, bet.AgencyID[0])
	msg = append(msg, bet.Documento...)
	msg = append(msg, bet.Nacimiento...)

	msg = append(msg, byte(len(bet.Nombre)))
	msg = append(msg, bet.Nombre...)

	msg = append(msg, byte(len(bet.Apellido)))
	msg = append(msg, bet.Apellido...)

	numLen := len(bet.Numero)
	msg = append(msg, EncodeUint16BE(numLen)...)
	msg = append(msg, bet.Numero...)

	return msg, nil
}

func (p *mixedSchemaProtocol) SerializeBatch(bets []*models.Bet) ([]byte, error) {
	if len(bets) > MaxUint16Value {
		return nil, ErrBatchTooLarge
	}

	msg := make([]byte, 0)
	msg = append(msg, MessageTypeBatch)
	msg = append(msg, EncodeUint16BE(len(bets))...)

	for _, bet := range bets {
		payload, err := p.Serialize(bet)
		if err != nil {
			return nil, err
		}
		if len(payload) > MaxUint16Value {
			return nil, ErrBetPayloadTooLarge
		}

		msg = append(msg, EncodeUint16BE(len(payload))...)
		msg = append(msg, payload...)
	}

	return msg, nil
}

func (p *mixedSchemaProtocol) SerializeDone(agencyID string) ([]byte, error) {
	if len(agencyID) != 1 {
		return nil, ErrInvalidAgencyID
	}
	return []byte{MessageTypeDone, agencyID[0]}, nil
}

func (p *mixedSchemaProtocol) SerializeQueryWinners(agencyID string) ([]byte, error) {
	if len(agencyID) != 1 {
		return nil, ErrInvalidAgencyID
	}
	return []byte{MessageTypeQueryWinners, agencyID[0]}, nil
}












