from .utils import Bet

BET_PROCESSED_ACK = b"ok\n"

RECV_BUFFER_SIZE = 1024
MIN_HEADER_SIZE = 20

AGENCY_OFFSET = 0
AGENCY_SIZE = 1
DNI_OFFSET = 1
DNI_SIZE = 8
BIRTHDATE_OFFSET = 9
BIRTHDATE_SIZE = 10
NAME_LEN_OFFSET = 19
NAME_START = 20

BYTE_HIGH_SHIFT = 8
BYTE_LOW_MASK = 0xFF


def recv_full_bet_message(client_sock):
    data = bytearray()
    expected_size = None

    while expected_size is None or len(data) < expected_size:
        chunk = client_sock.recv(RECV_BUFFER_SIZE)
        if not chunk:
            break
        data.extend(chunk)
        expected_size = expected_message_size(data)

    if expected_size is None:
        raise ValueError('incomplete bet message header')
    if len(data) < expected_size:
        raise ValueError('incomplete bet message body')

    return bytes(data[:expected_size])

# [agency:1][dni:8][birthdate:10][name_len:1][name][lastname_len:1][lastname][number_len:2][number]
def expected_message_size(data):
    if len(data) < MIN_HEADER_SIZE:
        return None

    name_len = data[NAME_LEN_OFFSET]
    lastname_len_idx = NAME_START + name_len
    if len(data) < lastname_len_idx + 1:
        return None

    lastname_len = data[lastname_len_idx]
    number_len_idx = lastname_len_idx + 1 + lastname_len
    if len(data) < number_len_idx + 2:
        return None

    number_len = (data[number_len_idx] << 8) | data[number_len_idx + 1]
    return number_len_idx + 2 + number_len


def deserialize_bet(data):
    agency_id = chr(data[AGENCY_OFFSET])
    dni = data[DNI_OFFSET:DNI_OFFSET + DNI_SIZE].decode('utf-8')
    birthdate = data[BIRTHDATE_OFFSET:BIRTHDATE_OFFSET + BIRTHDATE_SIZE].decode('utf-8')

    name_len = data[NAME_LEN_OFFSET]
    name_start = NAME_START
    name_end = name_start + name_len
    first_name = data[name_start:name_end].decode('utf-8')

    lastname_len = data[name_end]
    lastname_start = name_end + 1
    lastname_end = lastname_start + lastname_len
    last_name = data[lastname_start:lastname_end].decode('utf-8')

    number_len = (data[lastname_end] << BYTE_HIGH_SHIFT) | (data[lastname_end + 1] & BYTE_LOW_MASK)
    number_start = lastname_end + 2
    number_end = number_start + number_len
    number = data[number_start:number_end].decode('utf-8')

    return Bet(agency_id, first_name, last_name, dni, birthdate, number)
