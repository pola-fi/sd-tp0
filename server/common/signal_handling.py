import logging
import signal
from functools import partial


def _handle_sigterm(server, _signum, _frame):
    logging.info('action: shutdown | result: in_progress | component: server | signal: SIGTERM')
    server.shutdown()


def register_sigterm_handler(server):
    signal.signal(signal.SIGTERM, partial(_handle_sigterm, server))
