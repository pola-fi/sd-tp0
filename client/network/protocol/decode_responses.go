package protocol

import (
	"strings"
)

func (p *mixedSchemaProtocol) DeserializeBatchResponse(data []byte) (*BatchResponse, error) {
	if len(data) != BatchResponseSize {
		return nil, ErrInvalidBatchResponseSize
	}

	status := data[0]
	if status != BatchStatusSuccess && status != BatchStatusFail {
		return nil, ErrInvalidBatchStatus
	}

	count := DecodeUint16BE(data[1], data[2])
	return &BatchResponse{Success: status == BatchStatusSuccess, Count: count}, nil
}

func (p *mixedSchemaProtocol) DeserializeQueryWinnersResponse(data []byte) (*QueryWinnersResponse, error) {
	if len(data) < QueryHeaderSize {
		return nil, ErrInvalidQueryResponseSize
	}
	if data[0] != MessageTypeQueryWinnersResponse {
		return nil, ErrInvalidQueryMessageType
	}

	status := data[1]
	if status != QueryStatusPending && status != QueryStatusReady {
		return nil, ErrInvalidQueryResponseStatus
	}

	count := DecodeUint16BE(data[2], data[3])
	dniCSVLen := DecodeUint16BE(data[4], data[5])
	expectedSize := QueryHeaderSize + dniCSVLen
	if len(data) != expectedSize {
		return nil, ErrInvalidQueryResponseSize
	}

	winnerDNIs := make([]string, 0)
	if dniCSVLen > 0 {
		dniCSV := string(data[QueryHeaderSize:])
		parts := strings.Split(dniCSV, ",")
		winnerDNIs = make([]string, 0, len(parts))
		for _, dni := range parts {
			if dni != "" {
				winnerDNIs = append(winnerDNIs, dni)
			}
		}
	}

	return &QueryWinnersResponse{Ready: status == QueryStatusReady, Count: count, WinnerDNIs: winnerDNIs}, nil
}








