package common

import (
	"context"
	"time"
)

func (c *Client) shouldExit(ctx context.Context) bool {
	return ctx.Err() != nil
}

func (c *Client) handleError(ctx context.Context, action string, err error) bool {
	if err == nil {
		return false
	}
	if c.shouldExit(ctx) {
		return true
	}
	log.Errorf("action: %s | result: fail | client_id: %v | error: %v", action, c.config.ID, err)
	return true
}

func (c *Client) waitLoopPeriod(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	case <-time.After(c.config.LoopPeriod):
		return true
	}
}

