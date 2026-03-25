import logging

from .draw_state import DrawState
from .protocol.batch import encode_batch_response
from .protocol.bet import deserialize_batch
from .protocol.message_types import (
    MESSAGE_TYPE_BATCH,
    MESSAGE_TYPE_DONE,
    MESSAGE_TYPE_QUERY_WINNERS,
)
from .protocol.query_winners import encode_query_winners_response
from .utils import store_bets
from .winners_service import WinnersService


class BetHandler:
    def __init__(self, expected_agencies=5):
        self._draw_state = DrawState(expected_agencies=expected_agencies)
        self._winners_service = WinnersService()

    def process(self, raw_msg):
        if not raw_msg:
            raise ValueError('empty message')

        message_type = raw_msg[0]
        if message_type == MESSAGE_TYPE_BATCH:
            return self._process_batch(raw_msg[1:])
        if message_type == MESSAGE_TYPE_DONE:
            return self._process_done(raw_msg)
        if message_type == MESSAGE_TYPE_QUERY_WINNERS:
            return self._process_query_winners(raw_msg)
        raise ValueError(f'unknown message type: {message_type}')

    def _process_batch(self, payload):
        bets = deserialize_batch(payload)
        try:
            store_bets(bets)
        except Exception:
            logging.info(f'action: apuesta_recibida | result: fail | cantidad: {len(bets)}')
            return encode_batch_response(success=False, count=len(bets))

        self._draw_state.mark_agencies_seen(bets)
        logging.info(f'action: apuesta_recibida | result: success | cantidad: {len(bets)}')
        return encode_batch_response(success=True, count=len(bets))

    def _process_done(self, raw_msg):
        if len(raw_msg) != 2:
            return encode_batch_response(success=False, count=0)

        agency_id = chr(raw_msg[1])
        if self._draw_state.register_done(agency_id):
            winners_by_agency = self._winners_service.compute_winners_by_agency()
            self._draw_state.publish_winners(winners_by_agency)
            logging.info('action: sorteo | result: success')

        return encode_batch_response(success=True, count=0)

    def _process_query_winners(self, raw_msg):
        if len(raw_msg) != 2:
            return encode_query_winners_response(False, 0, [])

        agency_id = chr(raw_msg[1])
        ready, winners = self._draw_state.get_winners_if_ready(agency_id)
        return encode_query_winners_response(ready, len(winners), winners)
