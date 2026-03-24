import struct

BYTE_HIGH_SHIFT = 8
BYTE_LOW_MASK = 0xFF
NETWORK_ENDIAN_PREFIX = "!"

# Reusable struct format for [u8][u16] in network byte order (big-endian)
U8_U16_BE_FORMAT = f"{NETWORK_ENDIAN_PREFIX}BH"

def read_u16_be(data, idx):
    return (data[idx] << BYTE_HIGH_SHIFT) | (data[idx + 1] & BYTE_LOW_MASK)


def pack_u8_u16_be(first: int, second: int) -> bytes:
    return struct.pack(U8_U16_BE_FORMAT, first, second)


def unpack_u8_u16_be(data: bytes):
    return struct.unpack(U8_U16_BE_FORMAT, data)

