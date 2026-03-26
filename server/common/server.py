import socket
import logging
from .exceptions import ServerShuttingDown


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._is_shutting_down = False

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        while not self._is_shutting_down:
            try:
                client_sock = self.__accept_new_connection()
                self.__handle_client_connection(client_sock)
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

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        peer_ip = None
        try:
            peer_ip = client_sock.getpeername()[0]
        except OSError:
            pass
        try:
            # TODO: Modify the receive to avoid short-reads
            msg = client_sock.recv(1024).rstrip().decode('utf-8')
            addr = client_sock.getpeername()
            logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {msg}')
            # TODO: Modify the send to avoid short-writes
            client_sock.send("{}\n".format(msg).encode('utf-8'))
        except OSError as e:
            if not self._is_shutting_down:
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
