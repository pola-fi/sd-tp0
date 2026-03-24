from ..utils import Bet
from .endian import read_u16_be
from .receiver import BATCH_COUNT_SIZE, BET_LENGTH_SIZE

# Single bet schema: [agency:1][dni:8][birthdate:10][name_len:1][name][lastname_len:1][lastname][number_len:2][number]
MIN_HEADER_SIZE = 20
AGENCY_OFFSET = 0
DNI_OFFSET = 1
DNI_SIZE = 8
BIRTHDATE_OFFSET = 9
BIRTHDATE_SIZE = 10
NAME_LEN_OFFSET = 19
NAME_START = 20


def deserialize_batch(data):
    if len(data) < BATCH_COUNT_SIZE:
        raise ValueError('missing batch count')

    batch_count = read_u16_be(data, 0)
    idx = BATCH_COUNT_SIZE
    bets = []

    for _ in range(batch_count):
        if idx + BET_LENGTH_SIZE > len(data):
            raise ValueError('missing bet length prefix')

        bet_length = read_u16_be(data, idx)
        idx += BET_LENGTH_SIZE
        end = idx + bet_length
        if end > len(data):
            raise ValueError('truncated bet payload')

        bets.append(deserialize_bet(data[idx:end]))
        idx = end

    if idx != len(data):
        raise ValueError('unexpected trailing bytes in batch message')

    return bets


def deserialize_bet(data):
    if len(data) < MIN_HEADER_SIZE:
        raise ValueError('bet payload too short')

    name_len = data[NAME_LEN_OFFSET]
    lastname_len_idx = NAME_START + name_len
    if len(data) < lastname_len_idx + 1:
        raise ValueError('missing lastname length')

    lastname_len = data[lastname_len_idx]
    number_len_idx = lastname_len_idx + 1 + lastname_len
    if len(data) < number_len_idx + 2:
        raise ValueError('missing number length')

    number_len = read_u16_be(data, number_len_idx)
    expected_len = number_len_idx + 2 + number_len
    if len(data) != expected_len:
        raise ValueError('invalid bet payload size')

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

    number_len = read_u16_be(data, lastname_end)
    number_start = lastname_end + 2
    number_end = number_start + number_len
    number = data[number_start:number_end].decode('utf-8')

    return Bet(agency_id, first_name, last_name, dni, birthdate, number)

