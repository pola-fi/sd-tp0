package core

import (
	"context"
	"time"

	"github.com/op/go-logging"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/models"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/network"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/network/protocol"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	AgencyID      string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	BatchMaxAmount int
}

// Client Entity that encapsulates how
type Client struct {
	config   ClientConfig
	bets     []*models.Bet
	sender   network.MessageSender
	proto    protocol.BetProtocol
	chunker  *Chunker
}

// NewClient Initializes a new client with configuration and pre-loaded bets.
func NewClient(config ClientConfig, bets []*models.Bet) *Client {
	proto := protocol.NewMixedSchemaProtocol()
	return &Client{
		config:   config,
		bets:     bets,
		sender:   network.NewBetMessageSender(),
		proto:    proto,
		chunker:  NewChunker(proto, config.BatchMaxAmount),
	}
}


// shouldExit Check if the client should exit
func (c *Client) shouldExit(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func (c *Client) handleError(action string, err error) bool {
	if err != nil {
		log.Errorf("action: %s | result: fail | error: %v", action, err)
		return true
	}
	return false
}

// StartClientLoop sends all bets for the agency in configured batches.
func (c *Client) StartClientLoop(ctx context.Context) {
	for start := 0; start < len(c.bets); {
		if c.shouldExit(ctx) {
			return
		}

		batchSize, err := c.chunker.NextBatchSize(c.bets, start)
		if c.handleError("build_batch", err) {
			return
		}

		response, err := c.sender.SendBatch(c.bets[start:start+batchSize], c.config.ServerAddress)
		if c.handleError("send_batch", err) {
			return
		}

		log.Infof("action: apuesta_enviada | result: success | cantidad: %v", response.Count)
		start += batchSize
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}


