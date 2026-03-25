package core

import "errors"

var (
	// Batch chunking errors
	ErrSingleBetExceedsMaxPacket = errors.New("single bet exceeds max packet size")
	ErrCouldNotBuildBatch        = errors.New("could not build batch")
)

