from dataclasses import dataclass
from enum import IntEnum
import struct

from .endian import U8_U16_BE_FORMAT, pack_u8_u16_be, unpack_u8_u16_be

MAX_UINT16_VALUE = 65535


class BatchStatus(IntEnum):
    FAIL = 0
    SUCCESS = 1


@dataclass(frozen=True)
class BatchResponse:
    status: BatchStatus
    count: int


# Response frame: [status:1][count:2] in network byte order
RESPONSE_SIZE = struct.calcsize(U8_U16_BE_FORMAT)

def encode_batch_response(success: bool, count: int) -> bytes:
    status = BatchStatus.SUCCESS if success else BatchStatus.FAIL
    safe_count = max(0, min(count, MAX_UINT16_VALUE))
    return pack_u8_u16_be(int(status), safe_count)


def decode_batch_response(data: bytes) -> BatchResponse:
    if len(data) != RESPONSE_SIZE:
        raise ValueError(f"invalid batch response size: {len(data)}, expected {RESPONSE_SIZE}")

    status_raw, count = unpack_u8_u16_be(data)
    try:
        status = BatchStatus(status_raw)
    except ValueError as e:
        raise ValueError(f"invalid batch status: {status_raw}") from e

    return BatchResponse(status=status, count=count)

