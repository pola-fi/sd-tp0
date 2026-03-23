#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 2 ]]; then
  echo "Uso: $0 <archivo-salida> <cantidad-clientes>" >&2
  exit 1
fi

output_file="$1"
client_count="$2"

if ! [[ "$client_count" =~ ^[1-9][0-9]*$ ]]; then
  echo "Error: la cantidad de clientes debe ser un entero positivo" >&2
  exit 1
fi

{
  cat <<'YAML'
name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /app/main.py
    environment:
      - PYTHONUNBUFFERED=1
      - LOGGING_LEVEL=DEBUG
    volumes:
      - ./server/config.ini:/app/config.ini
    networks:
      - testing_net

YAML

  for i in $(seq 1 "$client_count"); do
    cat <<YAML
  client${i}:
    container_name: client${i}
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID=${i}
      - CLI_LOG_LEVEL=DEBUG
    volumes:
      - ./client/config.yaml:/app/config.yaml
    networks:
      - testing_net
    depends_on:
      - server

YAML
  done

  cat <<'YAML'
networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
YAML
} > "$output_file"

echo "Docker Compose generado en $output_file con $client_count cliente(s)."

