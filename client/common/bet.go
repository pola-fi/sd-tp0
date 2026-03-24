package common

import (
	"errors"
)

type Bet struct {
	AgencyID    string
	Nombre      string
	Apellido    string
	Documento   string
	Nacimiento  string
	Numero      string
}

func NewBet(agencyID, nombre, apellido, documento, nacimiento, numero string) (*Bet, error) {
	
	if len(documento) != 8 {
		return nil, errors.New("DNI debe tener 8 dígitos")
	}
	if len(nacimiento) != 10 {
		return nil, errors.New("Fecha debe tener formato YYYY-MM-DD (10 caracteres)")
	}
	if len(agencyID) != 1 {
		return nil, errors.New("Agency ID debe tener 1 dígito")
	}
	if len(nombre) > 255 {
		return nil, errors.New("Nombre demasiado largo (máximo 255 caracteres)")
	}
	if len(apellido) > 255 {
		return nil, errors.New("Apellido demasiado largo (máximo 255 caracteres)")
	}
	if len(numero) > 65535 {
		return nil, errors.New("Número demasiado largo (máximo 65535 caracteres)")
	}

	return &Bet{
		AgencyID:   agencyID,
		Nombre:     nombre,
		Apellido:   apellido,
		Documento:  documento,
		Nacimiento: nacimiento,
		Numero:     numero,
	}, nil
}

func (b *Bet) GetDocument() string {
	return b.Documento
}

func (b *Bet) GetNumber() string {
	return b.Numero
}
