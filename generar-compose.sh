#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 2 ]]; then
  echo "Uso: $0 <archivo-salida> <cantidad-clientes>" >&2
  exit 1
fi

output_file="$1"
client_count="$2"

if ! [[ "$client_count" =~ ^[0-9]+$ ]]; then
  echo "Error: la cantidad de clientes debe ser un entero no negativo" >&2
  exit 1
fi

NOMBRES=("Juan" "Maria" "Carlos" "Ana" "Luis" "Laura" "Pedro" "Sofia" "Diego" "Valentina")
APELLIDOS=("Garcia" "Rodriguez" "Lopez" "Martinez" "Gonzalez" "Perez" "Sanchez" "Ramirez" "Torres" "Diaz")

generar_dni() {
  echo $((10000000 + RANDOM % 90000000))
}

generar_fecha() {
  year=$((1960 + RANDOM % 41))
  month=$((1 + RANDOM % 12))
  day=$((1 + RANDOM % 28))
  printf "%04d-%02d-%02d" $year $month $day
}

generar_numero() {
  echo $((1000 + RANDOM % 9000))
}

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
    volumes:
      - ./server/config.ini:/app/config.ini
    networks:
      - testing_net

YAML

  for i in $(seq 1 "$client_count"); do
    agency_id=$i
    nombre_idx=$((RANDOM % ${#NOMBRES[@]}))
    apellido_idx=$((RANDOM % ${#APELLIDOS[@]}))
    nombre=${NOMBRES[$nombre_idx]}
    apellido=${APELLIDOS[$apellido_idx]}
    documento=$(generar_dni)
    nacimiento=$(generar_fecha)
    numero=$(generar_numero)

    cat <<YAML
  client${i}:
    container_name: client${i}
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID=${i}
      - CLI_AGENCY_ID=${agency_id}
      - CLI_BET_NOMBRE=${nombre}
      - CLI_BET_APELLIDO=${apellido}
      - CLI_BET_DOCUMENTO=${documento}
      - CLI_BET_NACIMIENTO=${nacimiento}
      - CLI_BET_NUMERO=${numero}
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

