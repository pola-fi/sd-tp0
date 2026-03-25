package protocol

import "errors"

var (
	// Batch encoding errors
	ErrBatchTooLarge          = errors.New("batch too large")
	ErrBetPayloadTooLarge     = errors.New("single bet payload too large")
	ErrInvalidAgencyID        = errors.New("invalid agency id")
	ErrInvalidSerializedSize  = errors.New("invalid serialized batch size")

	// Batch response decoding errors
	ErrInvalidBatchResponseSize = errors.New("invalid batch response size")
	ErrInvalidBatchStatus       = errors.New("invalid batch status")

	// Query response decoding errors
	ErrInvalidQueryResponseSize   = errors.New("invalid query response size")
	ErrInvalidQueryMessageType    = errors.New("invalid query response message type")
	ErrInvalidQueryResponseStatus = errors.New("invalid query response status")
)

