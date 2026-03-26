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
	ID             string
	AgencyID       string
	ServerAddress  string
	LoopAmount     int
	LoopPeriod     time.Duration
	BatchMaxAmount int
}

// Client Entity that encapsulates how
type Client struct {
	config  ClientConfig
	bets    []*models.Bet
	sender  network.MessageSender
	chunker *Chunker
}

// NewClient Initializes a new client with configuration and pre-loaded bets.
func NewClient(config ClientConfig, bets []*models.Bet) *Client {
	proto := protocol.NewBetProtocol()
	return &Client{
		config:  config,
		bets:    bets,
		sender:  network.NewBetMessageSender(config.ID),
		chunker: NewChunker(proto, config.BatchMaxAmount),
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
	if !c.dispatchAllBatches(ctx) {
		return
	}

	if c.handleError("notify_done", c.sender.SendDone(c.config.AgencyID, c.config.ServerAddress)) {
		return
	}

	if !c.waitForWinners(ctx) {
		return
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

