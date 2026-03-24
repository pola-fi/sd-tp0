import logging

from .protocol.bet import deserialize_batch
from .protocol.batch import encode_batch_response
from .utils import store_bets

class BetHandler:
    def process(self, raw_msg):
        bets = deserialize_batch(raw_msg)
        try:
            store_bets(bets)
        except Exception:
            logging.info(f'action: apuesta_recibida | result: fail | cantidad: {len(bets)}')
            return encode_batch_response(success=False, count=len(bets))

        logging.info(f'action: apuesta_recibida | result: success | cantidad: {len(bets)}')
        return encode_batch_response(success=True, count=len(bets))

