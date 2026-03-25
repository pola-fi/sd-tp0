package core

import "context"

func (c *Client) dispatchAllBatches(ctx context.Context) bool {
	for start := 0; start < len(c.bets); {
		if c.shouldExit(ctx) {
			return false
		}

		batchSize, err := c.chunker.NextBatchSize(c.bets, start)
		if c.handleError("build_batch", err) {
			return false
		}

		batch := c.bets[start : start+batchSize]
		response, err := c.sender.SendBatch(batch, c.config.ServerAddress)
		if c.handleError("send_batch", err) {
			return false
		}

		c.logBatchSent(batch, response.Count)
		start += batchSize
	}

	return true
}

