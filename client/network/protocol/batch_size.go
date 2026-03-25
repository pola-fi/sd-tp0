package protocol

import (
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/models"
)

func EstimateSingleBetBatchSize(protocol BetProtocol, bet *models.Bet) (int, error) {
	serializedBatch, err := protocol.SerializeBatch([]*models.Bet{bet})
	if err != nil {
		return 0, err
	}

	batchHeaderSize := 1 + BatchCountFieldSize
	if len(serializedBatch) < batchHeaderSize {
		return 0, ErrInvalidSerializedSize
	}

	return len(serializedBatch) - batchHeaderSize, nil
}


