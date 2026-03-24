import logging

from .bet_protocol import deserialize_bet, BET_PROCESSED_ACK
from .utils import store_bets

class BetHandler:
    def process(self, raw_msg):
        bet = deserialize_bet(raw_msg)
        store_bets([bet])
        logging.info(
            f'action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}'
        )
        return BET_PROCESSED_ACK

