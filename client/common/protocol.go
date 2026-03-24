package common

type BetProtocol interface {
	Serialize(bet *Bet) ([]byte, error)
}

type MixedSchemaProtocol struct{}

func NewMixedSchemaProtocol() *MixedSchemaProtocol {
	return &MixedSchemaProtocol{}
}

// Format: [agency_id:1][dni:8][fecha:10][len_nombre:1][nombre][len_apellido:1][apellido][len_numero:2][numero]
func (p *MixedSchemaProtocol) Serialize(bet *Bet) ([]byte, error) {
	
	msg := make([]byte, 0)
	
	// Fixed fields
	msg = append(msg, bet.AgencyID[0])     // 1 byte agency_id
	msg = append(msg, bet.Documento...)    // 8 bytes dni
	msg = append(msg, bet.Nacimiento...)   // 10 bytes fecha
	
	// Variable fields with length prefix
	msg = append(msg, byte(len(bet.Nombre)))
	msg = append(msg, bet.Nombre...)
	
	msg = append(msg, byte(len(bet.Apellido)))
	msg = append(msg, bet.Apellido...)
	
	// 2 bytes for number length (big endian)
	numLen := len(bet.Numero)
	msg = append(msg, byte(numLen>>8))
	msg = append(msg, byte(numLen&0xFF))
	msg = append(msg, bet.Numero...)
	
	return msg, nil
}
