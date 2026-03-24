package core

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/models"
)

var ErrAgencyFileEmpty = errors.New("agency file is empty")

func LoadAgencyBets(agencyID string) ([]*models.Bet, error) {
	path := fmt.Sprintf(".data/agency-%s.csv", agencyID)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	bets := make([]*models.Bet, 0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(record) != 5 {
			return nil, fmt.Errorf("invalid csv row length: %d", len(record))
		}

		bet, err := models.NewBet(agencyID, record[0], record[1], record[2], record[3], record[4])
		if err != nil {
			return nil, err
		}
		bets = append(bets, bet)
	}

	if len(bets) == 0 {
		return nil, ErrAgencyFileEmpty
	}

	return bets, nil
}

