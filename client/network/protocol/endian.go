package protocol

const Uint16ByteSize = 2

func EncodeUint16BE(v int) []byte {
	return []byte{byte(v >> 8), byte(v & 0xFF)}
}

func DecodeUint16BE(msb, lsb byte) int {
	return int(msb)<<8 | int(lsb)
}

