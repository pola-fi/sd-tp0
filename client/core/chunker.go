package core

import (
	"errors"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/models"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/network/protocol"
)

var (
	ErrSingleBetExceedsMaxPacket = errors.New("single bet exceeds max packet size")
	ErrCouldNotBuildBatch        = errors.New("could not build batch")
)

type Chunker struct {
	protocol  protocol.BetProtocol
	maxAmount int
}

func NewChunker(protocol protocol.BetProtocol, maxAmount int) *Chunker {
	return &Chunker{
		protocol:  protocol,
		maxAmount: maxAmount,
	}
}

func (bc *Chunker) NextBatchSize(bets []*models.Bet, start int) (int, error) {
	maxAmount := bc.maxAmount
	if maxAmount <= 0 {
		maxAmount = 1
	}

	remaining := len(bets) - start
	if remaining > maxAmount {
		remaining = maxAmount
	}

	currentBytes := protocol.BatchCountFieldSize
	batchSize := 0
	for batchSize < remaining {
		size, err := protocol.EstimateSingleBetBatchSize(bc.protocol, bets[start+batchSize])
		if err != nil {
			return 0, err
		}

		if currentBytes+size > protocol.MaxPacketBytes {
			if batchSize == 0 {
				return 0, ErrSingleBetExceedsMaxPacket
			}
			break
		}

		currentBytes += size
		batchSize++
	}

	if batchSize == 0 {
		return 0, ErrCouldNotBuildBatch
	}

	return batchSize, nil
}

