package common

import (
	"context"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
}

type BetConfig struct {
	AgencyID      string
	Nombre        string
	Apellido      string
	Documento     string
	Nacimiento    string
	Numero        string
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	bet    *Bet
	sender MessageSender
}

// NewClient Initializes a new client receiving the configuration and bet
// as parameters
func NewClient(config ClientConfig, bet *Bet) *Client {
	client := &Client{
		config: config,
		bet:    bet,
		sender: NewBetMessageSender(),
	}
	return client
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

// handleError Handle errors and check if the client should exit
func (c *Client) handleError(ctx context.Context, action string, err error) bool {
	if err != nil {
		log.Errorf("action: %s | result: fail | error: %v", action, err)
		return true
	}
	return false
}

// waitLoopPeriod Wait the configured period between messages
func (c *Client) waitLoopPeriod(ctx context.Context) bool {
	select {
	case <-time.After(c.config.LoopPeriod):
		return true
	case <-ctx.Done():
		return false
	}
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop(ctx context.Context) {
	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
		if c.shouldExit(ctx) {
			return
		}

		response, err := c.sender.SendBet(c.bet, c.config.ServerAddress)
		if c.handleError(ctx, "send_bet", err) {
			return
		}

		log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v | response: %v",
			c.bet.GetDocument(),
			c.bet.GetNumber(),
			response,
		)

		// Wait a time between sending one message and the next one
		if !c.waitLoopPeriod(ctx) {
			return
		}

	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
