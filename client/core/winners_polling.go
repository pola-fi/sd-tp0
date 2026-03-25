package core

import (
	"context"
	"time"
)

func (c *Client) waitForWinners(ctx context.Context) bool {
	for {
		if c.shouldExit(ctx) {
			return false
		}

		queryResp, err := c.sender.QueryWinners(c.config.AgencyID, c.config.ServerAddress)
		if c.handleError("consulta_ganadores", err) {
			return false
		}

		if queryResp.Ready {
			c.logWinnersReady(queryResp)
			return true
		}

		time.Sleep(200 * time.Millisecond)
	}
}

