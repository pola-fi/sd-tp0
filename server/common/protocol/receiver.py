from ..utils import Bet
from .endian import read_u16_be

MAX_PACKET_BYTES = 8192
RECV_BUFFER_SIZE = 1024

# Request frame: [batch_count:2][bet_len:2][bet_payload]...
BATCH_COUNT_SIZE = 2
BET_LENGTH_SIZE = 2


def recv_full_batch_message(client_sock):
    data = bytearray()
    expected_size = None

    while expected_size is None or len(data) < expected_size:
        chunk = client_sock.recv(RECV_BUFFER_SIZE)
        if not chunk:
            break
        data.extend(chunk)
        if len(data) > MAX_PACKET_BYTES:
            raise ValueError('batch message exceeds 8kB limit')
        expected_size = expected_batch_message_size(data)

    if expected_size is None:
        raise ValueError('incomplete batch message header')
    if len(data) < expected_size:
        raise ValueError('incomplete batch message body')

    return bytes(data[:expected_size])


def expected_batch_message_size(data):
    if len(data) < BATCH_COUNT_SIZE:
        return None

    batch_count = read_u16_be(data, 0)
    idx = BATCH_COUNT_SIZE
    for _ in range(batch_count):
        if len(data) < idx + BET_LENGTH_SIZE:
            return None
        bet_length = read_u16_be(data, idx)
        idx += BET_LENGTH_SIZE
        if len(data) < idx + bet_length:
            return None
        idx += bet_length

    return idx

