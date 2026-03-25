package protocol

type mixedSchemaProtocol struct{}

func NewBetProtocol() BetProtocol {
	return &mixedSchemaProtocol{}
}
