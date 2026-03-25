package network

import "errors"

var (
	ErrConnectionFailed     = errors.New("failed to connect to server")
	ErrSerializationFailed  = errors.New("failed to serialize batch")
	ErrSendFailed           = errors.New("failed to send data")
	ErrReadResponseFailed   = errors.New("failed to read response")
	ErrServerRejectedBatch  = errors.New("server rejected batch")
	ErrControlRequestFailed = errors.New("failed to send control request")
)

