package core

import (
	"strings"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/models"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/network/protocol"
)

func (c *Client) logBatchSent(batch []*models.Bet, count int) {
	dniStr, numeroStr := joinBetFields(batch)
	log.Infof("action: apuesta_enviada | result: success | cantidad: %v | dnis: [%s] | numeros: [%s]", count, dniStr, numeroStr)
}

func (c *Client) logWinnersReady(response *protocol.QueryWinnersResponse) {
	dniStr := strings.Join(response.WinnerDNIs, ",")
	log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %v | dnis: [%s]", response.Count, dniStr)
}

func joinBetFields(batch []*models.Bet) (string, string) {
	dnis := make([]string, len(batch))
	numeros := make([]string, len(batch))
	for i, bet := range batch {
		dnis[i] = bet.Documento
		numeros[i] = bet.Numero
	}
	return strings.Join(dnis, ","), strings.Join(numeros, ",")
}

