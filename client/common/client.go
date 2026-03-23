package common

import (
	"bufio"
	"context"
	"fmt"
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

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   connectionState
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop(ctx context.Context) {
	go func() {
		<-ctx.Done()
		c.conn.close()
	}()

	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
		if c.shouldExit(ctx) {
			return
		}

		// Create the connection the server in every loop iteration. Send an
		if c.handleError(ctx, "connect", c.conn.connect(c.config.ServerAddress)) {
			return
		}

		conn := c.conn.get()
		if conn == nil {
			if c.shouldExit(ctx) {
				return
			}
			log.Errorf("action: connect | result: fail | client_id: %v | error: nil connection", c.config.ID)
			return
		}

		// TODO: Modify the send to avoid short-write
		_, err := fmt.Fprintf(
			conn,
			"[CLIENT %v] Message N°%v\n",
			c.config.ID,
			msgID,
		)
		if c.handleError(ctx, "send_message", err) {
			c.conn.close()
			return
		}

		msg, err := bufio.NewReader(conn).ReadString('\n')
		c.conn.close()

		if c.handleError(ctx, "receive_message", err) {
			return
		}

		log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
			c.config.ID,
			msg,
		)

		// Wait a time between sending one message and the next one
		if !c.waitLoopPeriod(ctx) {
			return
		}

	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
