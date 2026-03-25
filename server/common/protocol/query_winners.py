from .message_types import MESSAGE_TYPE_QUERY_WINNERS_RESPONSE

QUERY_STATUS_PENDING = 0
QUERY_STATUS_READY = 1
MAX_UINT16 = 65535


def encode_query_winners_response(ready, count, winner_dnis):
    status = QUERY_STATUS_READY if ready else QUERY_STATUS_PENDING
    winners_csv = ','.join(winner_dnis).encode('utf-8')
    safe_count = max(0, min(count, MAX_UINT16))
    safe_csv_len = max(0, min(len(winners_csv), MAX_UINT16))
    winners_payload = winners_csv[:safe_csv_len]

    header = bytes([MESSAGE_TYPE_QUERY_WINNERS_RESPONSE, status])
    header += safe_count.to_bytes(2, byteorder='big')
    header += safe_csv_len.to_bytes(2, byteorder='big')
    return header + winners_payload

