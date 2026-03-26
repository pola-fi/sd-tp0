import logging
import socket
import threading

from .exceptions import ServerShuttingDown
from .handler import BetHandler
from .protocol.endian import read_u16_be
from .protocol.message_types import (
    MESSAGE_TYPE_BATCH,
    MESSAGE_TYPE_DONE,
    MESSAGE_TYPE_QUERY_WINNERS,
)
from .protocol.query_winners import QUERY_STATUS_PENDING, QUERY_STATUS_READY
from .protocol.receiver import recv_full_message


class Server:
    def __init__(self, port, listen_backlog, expected_agencies=5, bet_handler=None):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._is_shutting_down = False
        self._bet_handler = bet_handler if bet_handler is not None else BetHandler(expected_agencies=expected_agencies)
        self._workers_lock = threading.Lock()
        self._workers = set()

    def run(self):
        """
        Server loop that accepts new connections and delegates each client
        to a dedicated worker thread.
        """
        while not self._is_shutting_down:
            try:
                client_sock = self.__accept_new_connection()
                self.__start_client_worker(client_sock)
            except ServerShuttingDown:
                break

    def shutdown(self):
        if self._is_shutting_down:
            return

        self._is_shutting_down = True
        try:
            self._server_socket.close()
            logging.info('action: close_socket | result: success | resource: server_socket')
        except OSError:
            pass

        self.__join_workers()

    def __start_client_worker(self, client_sock):
        worker = threading.Thread(target=self.__client_worker_entry, args=(client_sock,))
        with self._workers_lock:
            self._workers.add(worker)
        worker.start()

    def __client_worker_entry(self, client_sock):
        try:
            self.__handle_client_connection(client_sock)
        finally:
            current = threading.current_thread()
            with self._workers_lock:
                self._workers.discard(current)

    def __join_workers(self):
        with self._workers_lock:
            workers = list(self._workers)

        for worker in workers:
            worker.join(timeout=2)

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket.

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        peer_ip = None
        try:
            peer_ip = client_sock.getpeername()[0]
        except OSError:
            pass
        try:
            msg = recv_full_message(client_sock)
            addr = client_sock.getpeername()
            response = self._bet_handler.process(msg)
            self.__log_receive_success(addr[0], msg, response)
            client_sock.sendall(response)
        except OSError as e:
            if not self._is_shutting_down:
                logging.error(f"action: receive_message | result: fail | error: {e}")
        except (ValueError, UnicodeDecodeError) as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()
            if peer_ip is not None:
                logging.info(
                    f'action: close_socket | result: success | resource: client_socket | ip: {peer_ip}'
                )
            else:
                logging.info(
                    'action: close_socket | result: success | resource: client_socket'
                )

    def __log_receive_success(self, ip, msg, response):
        summary = self.__summarize_message(msg)
        if msg and msg[0] == MESSAGE_TYPE_QUERY_WINNERS:
            query_status = response[1] if len(response) >= 2 else None
            if query_status == QUERY_STATUS_PENDING:
                logging.debug(f'action: receive_message | result: success | ip: {ip} | {summary} | query_status: pending')
                return
            if query_status == QUERY_STATUS_READY:
                logging.info(f'action: receive_message | result: success | ip: {ip} | {summary} | query_status: ready')
                return

        logging.info(f'action: receive_message | result: success | ip: {ip} | {summary}')

    def __summarize_message(self, msg):
        if not msg:
            return 'msg_type: empty | bytes: 0'

        msg_type = msg[0]
        if msg_type == MESSAGE_TYPE_BATCH:
            batch_count = read_u16_be(msg, 1) if len(msg) >= 3 else 0
            return f'msg_type: batch | bytes: {len(msg)} | batch_count: {batch_count}'
        if msg_type == MESSAGE_TYPE_DONE:
            agency = chr(msg[1]) if len(msg) >= 2 else '?'
            return f'msg_type: done | bytes: {len(msg)} | agency_id: {agency}'
        if msg_type == MESSAGE_TYPE_QUERY_WINNERS:
            agency = chr(msg[1]) if len(msg) >= 2 else '?'
            return f'msg_type: query_winners | bytes: {len(msg)} | agency_id: {agency}'

        return f'msg_type: unknown({msg_type}) | bytes: {len(msg)}'

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        try:
            c, addr = self._server_socket.accept()
        except OSError:
            if self._is_shutting_down:
                raise ServerShuttingDown()
            raise
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
