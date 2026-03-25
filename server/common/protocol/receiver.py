from .endian import read_u16_be
from .message_types import MESSAGE_TYPE_BATCH, MESSAGE_TYPE_DONE, MESSAGE_TYPE_QUERY_WINNERS

MAX_PACKET_BYTES = 8192
RECV_BUFFER_SIZE = 1024

# Batch request frame: [type:1][batch_count:2][bet_len:2][bet_payload]...
MESSAGE_TYPE_SIZE = 1
BATCH_COUNT_SIZE = 2
BET_LENGTH_SIZE = 2
CONTROL_MESSAGE_SIZE = 2


def recv_full_message(client_sock):
    data = bytearray()
    expected_size = None

    while expected_size is None or len(data) < expected_size:
        chunk = client_sock.recv(RECV_BUFFER_SIZE)
        if not chunk:
            break
        data.extend(chunk)
        if len(data) > MAX_PACKET_BYTES:
            raise ValueError('batch message exceeds 8kB limit')
        expected_size = expected_message_size(data)

    if expected_size is None:
        raise ValueError('incomplete message header')
    if len(data) < expected_size:
        raise ValueError('incomplete message body')

    return bytes(data[:expected_size])


def expected_message_size(data):
    if len(data) < MESSAGE_TYPE_SIZE:
        return None

    message_type = data[0]
    if message_type == MESSAGE_TYPE_BATCH:
        return expected_batch_message_size(data)
    if message_type in (MESSAGE_TYPE_DONE, MESSAGE_TYPE_QUERY_WINNERS):
        return CONTROL_MESSAGE_SIZE
    raise ValueError(f'unknown message type: {message_type}')


def expected_batch_message_size(data):
    if len(data) < MESSAGE_TYPE_SIZE + BATCH_COUNT_SIZE:
        return None

    batch_count = read_u16_be(data, MESSAGE_TYPE_SIZE)
    idx = MESSAGE_TYPE_SIZE + BATCH_COUNT_SIZE
    for _ in range(batch_count):
        if len(data) < idx + BET_LENGTH_SIZE:
            return None
        bet_length = read_u16_be(data, idx)
        idx += BET_LENGTH_SIZE
        if len(data) < idx + bet_length:
            return None
        idx += bet_length

    return idx
