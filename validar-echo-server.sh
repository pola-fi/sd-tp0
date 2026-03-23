#!/bin/bash
set -euo pipefail

MENSAJE_TEST="HOLA"

RESPUESTA=$(docker run --rm \
    --network tp0_testing_net \
    busybox:latest \
    sh -c "echo '$MENSAJE_TEST' | nc -w 2 server 12345" 2>/dev/null || echo "")

if [ "$RESPUESTA" = "$MENSAJE_TEST" ]; then
    echo "action: test_echo_server | result: success"
    exit 0
else
    echo "action: test_echo_server | result: fail"
    exit 1
fi

